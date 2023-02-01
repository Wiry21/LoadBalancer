# LoadBalancer
Распределитель нагрузки

Задача:
Написать веб-сервис, распределяющий нагрузку между N веб сервисами равномерно. Если один из веб сервисов перестает отвечать, то исключить его из массива веб сервисов и распределять нагрузку между N-1 сервисами. 
Нагрузкой для каждого сервиса является количество запросов, находящихся у него в обработке. Относительно этого числа распределение должно быть равномерным. 

На выходе у вас должен получиться HTTP прокси сервер, который равномерно распределит нагрузку между N сервисами. 


Требования к ЯП: -

Требования к развертке: Docker

Требования к конфигурации: список веб сервисов с их адресами прописать в отдельном конфигурационном файле. 


Шаги тестирования: 

1)Поднять минимум 5 сервисов таргетов 

2)Поднять сервис распределитель 

3)Запустить скрипт спама запросами 

4)Каждые 10 секунд каждый из сервисов-таргетов должен логировать текущее количество запросов в обработке 

5)На 100-ой секунде отключить один из сервисов-таргетов 

6)На 200-ой включить его обратно 

По итогу тестирования лог файлы таргетов должны свидетельствовать о равномерном распределении нагрузки, в момент отключения таргета n, таргеты N-1 должны разделить между собой нагрузку. В момент включения таргета n, таргеты N должны освободить часть нагрузки под поднявшийся таргет.
