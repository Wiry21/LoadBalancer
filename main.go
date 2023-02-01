package main

import (
	"LoadBalancer/pkg/domain"
	"LoadBalancer/pkg/health"
	"LoadBalancer/pkg/strategy"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

	"LoadBalancer/pkg/config"
)

var (
	port       = flag.Int("port", 8080, "where to start balancer")
	configFile = flag.String("config-path", "", "The config file")
)

type Balancer struct {
	// Config is the configuration loaded from a config file
	Config *config.Config
	// ServerList will contain a mapping between matcher and replicas
	ServerList map[string]*config.ServerList
}

func NewBalancer(conf *config.Config) *Balancer {
	// prevent multiple or invalid matchers before creating the server
	serverMap := make(map[string]*config.ServerList, 0)

	for _, service := range conf.Services {
		servers := make([]*domain.Server, 0)
		for _, replica := range service.Replicas {
			ur, err := url.Parse(replica.Url)
			if err != nil {
				log.Fatal(err)
			}
			proxy := httputil.NewSingleHostReverseProxy(ur)
			servers = append(servers, &domain.Server{
				Url:      ur,
				Proxy:    proxy,
				Metadata: replica.Metadata,
			})
		}
		checker, err := health.NewChecker(nil, servers)
		if err != nil {
			log.Fatal(err)
		}
		serverMap[service.Matcher] = &config.ServerList{
			Servers:  servers,
			Name:     service.Name,
			Strategy: strategy.LoadStrategy(service.Strategy),
			Hc:       checker,
		}
	}
	// start all the health checkers for all provided matchers
	for _, sl := range serverMap {
		go sl.Hc.Start()
	}
	return &Balancer{
		Config:     conf,
		ServerList: serverMap,
	}

}

// Looks for the first server list that matches the reqPath (i.e. matcher)
// Will return an error if no matcher have been found.
func (b *Balancer) findServiceList(reqPath string) (*config.ServerList, error) {
	for matcher, s := range b.ServerList {
		if strings.HasPrefix(reqPath, matcher) {
			return s, nil
		}
	}
	return nil, fmt.Errorf("could not find a matcher for url: '%s'", reqPath)
}

func (b *Balancer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	// We need to support per service forwarding, i.e. this method should
	// read the request path, say host:port/service/rest/of/url this should
	// be load balanced against service named "service" and url will be
	// "host{i}:port{i}/rest/of/url
	sl, err := b.findServiceList(req.URL.Path)
	if err != nil {
		log.Error(err)
		res.WriteHeader(http.StatusNotFound)
		return
	}
	next, err := sl.Strategy.Next(sl.Servers)
	if err != nil {
		log.Error(err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	next.Proxy.ModifyResponse = func(res *http.Response) error {
		if res.StatusCode != 200 {
			return errors.New("not 200 status code from the host")
		}
		return nil
	}
	next.Proxy.ErrorHandler = func(res http.ResponseWriter, req *http.Request, err error) {
		fmt.Println(err)
		b.ServeHTTP(res, req)
	}
	next.Forward(res, req)
}

func main() {
	flag.Parse()
	file, err := os.Open(*configFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	conf, err := config.LoadConfig(file)
	if err != nil {
		log.Fatal(err)
	}

	balancer := NewBalancer(conf)
	server := http.Server{
		Addr:    fmt.Sprintf(":%d", *port),
		Handler: balancer,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
