package health

import (
	"github.com/Wiry21/LoadBalancer/pkg/domain"
	"errors"
	log "github.com/sirupsen/logrus"
	"net"
	"time"
)

type HealthChecker struct {
	servers []*domain.Server

	period int
}

// NewChecker will create a new HealthChecker.
func NewChecker(_conf *domain.Config, servers []*domain.Server) (*HealthChecker, error) {
	if len(servers) == 0 {
		return nil, errors.New("A server list expected, gotten an empty list")
	}
	return &HealthChecker{
		servers: servers,
	}, nil
}

// Start keeps looping indefinitely try to check the health of every server
// the caller is responsible for creating the goroutine when this should run
func (hc *HealthChecker) Start() {
	log.Info("Starting the health checker...")
	ticker := time.NewTicker(time.Second * 2)
	defer ticker.Stop()
	for {
		select {
		case _ = <-ticker.C:
			for _, server := range hc.servers {
				go checkHealth(server)
			}
		}
	}
}

// changes the liveness of the server (either from live to dead or the other way around)
func checkHealth(server *domain.Server) {
	_, err := net.DialTimeout("tcp", server.Url.Host, time.Second*5)
	if err != nil {
		log.Errorf("Could not connect to the server at '%s'", server.Url.Host)
		old := server.SetLiveness(false)
		if old {
			//atomic.StoreInt64(&server.Count, 0)
			log.Warnf("Transitioning server '%s' from Live to Unavailable state", server.Url.Host)
		}
		return
	}
	old := server.SetLiveness(true)
	if !old {
		log.Infof("Transitioning server '%s' from Unavailable to Live state", server.Url.Host)
	}
	log.Infof("Count requests at server '%s' = '%d'", server.Url.Host, server.Count)
}
