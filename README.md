## AIO [![GoDoc](https://godoc.org/github.com/itsmontoya/jsoon?status.svg)](https://godoc.org/github.com/itsmontoya/jsoon) ![Status](https://img.shields.io/badge/status-beta-yellow.svg)

AIO is an asynchronous io-manager in pure-go

## Benchmarks
```bash
## Go 1.7.4
# AIO Running wrk -c 20 http://172.16.0.201:1337/a
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   528.99us  359.63us  10.98ms   91.99%
    Req/Sec    19.86k     0.95k   22.17k    77.72%
  398809 requests in 10.10s, 3.50GB read
Requests/sec:  39488.34
Transfer/sec:    354.83MB

# Stdlib Running wrk -c 20 http://172.16.0.201:1337/a
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   801.86us  476.24us   6.19ms   76.31%
    Req/Sec    12.82k   648.63    13.77k    85.64%
  257670 requests in 10.10s, 2.26GB read
Requests/sec:  25513.25
Transfer/sec:    229.42MB
```
