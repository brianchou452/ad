# Dcard Backend Intern Assignment

## Introduction

使用 Golang 設計並且實作一個簡化的廣告投放服務,該服務應該有兩個 API,一個用於產生廣告,一個用於列出廣告。每個廣告都有它出現的條件(例如跟據使用者的年齡),產生廣告的 API 用來產生與設定條件。投放廣告的 API 就要跟據條件列出符合使用條件的廣告

## 如何使用

TBD

### 如何進行開發

### 如何進行測試

### 如何進行壓力測試

## 專案架構

## 用到的資源

- Grafana Dashboard: [k6-load-testing-results_rev3.json](https://github.com/luketn/docker-k6-grafana-influxdb/blob/master/dashboards/k6-load-testing-results_rev3.json)

```bash
     execution: local
        script: /scripts/main.js
        output: -

     scenarios: (100.00%) 2 scenarios, 515 max VUs, 1m1s max duration (incl. graceful stop):
              * addAd: 100 iterations for each of 15 VUs (maxDuration: 1m0s, exec: addAd, gracefulStop: 1s)
              * getAd: 500 looping VUs for 30s (exec: getAd, gracefulStop: 1s)


     ✓ Get AD API status is 200
     ✓ Admin API status is 200

     checks.........................: 100.00% ✓ 375461       ✗ 0     
     data_received..................: 52 MB   1.7 MB/s
     data_sent......................: 56 MB   1.9 MB/s
     http_req_blocked...............: avg=21.27µs  min=606ns    med=1.78µs  max=54.54ms  p(90)=2.74µs   p(95)=3.25µs 
     http_req_connecting............: avg=2.35µs   min=0s       med=0s      max=24.24ms  p(90)=0s       p(95)=0s     
     http_req_duration..............: avg=38.87ms  min=172.28µs med=35.74ms max=198.26ms p(90)=67.44ms  p(95)=81.21ms
       { expected_response:true }...: avg=38.87ms  min=172.28µs med=35.74ms max=198.26ms p(90)=67.44ms  p(95)=81.21ms
     http_req_failed................: 0.00%   ✓ 0            ✗ 375461
     http_req_receiving.............: avg=1.26ms   min=6.05µs   med=18.09µs max=100.27ms p(90)=147.14µs p(95)=3.39ms 
     http_req_sending...............: avg=159.15µs min=3.25µs   med=8.49µs  max=69.73ms  p(90)=15.45µs  p(95)=71.33µs
     http_req_tls_handshaking.......: avg=0s       min=0s       med=0s      max=0s       p(90)=0s       p(95)=0s     
     http_req_waiting...............: avg=37.45ms  min=151.71µs med=35.48ms max=163.03ms p(90)=63.48ms  p(95)=73.5ms 
     http_reqs......................: 375461  12507.181051/s
     iteration_duration.............: avg=39.94ms  min=254.63µs med=36.55ms max=218.29ms p(90)=69.22ms  p(95)=83.53ms
     iterations.....................: 375461  12507.181051/s
     vus............................: 500     min=500        max=515 
     vus_max........................: 515     min=515        max=515 


running (0m30.0s), 000/515 VUs, 375461 complete and 0 interrupted iterations
addAd ✓ [======================================] 15 VUs   0m02.2s/1m0s  1500/1500 iters, 100 per VU
getAd ✓ [======================================] 500 VUs  30s   
```
