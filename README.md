一个基于gin+gorm+beego/validation+redis搭建的简单blog-web后端，使用jwt来做身份验证,藉此了解学习gin框架。
计划作为个人博客后端。

简单的使用ab压测能有100左右的QPS(在自己的渣机上)比起java SpringBoot+Mybatis快了一倍以上。（同等数据量，同等测试参数（1000次请求100同时））（数据库redis均在远程阿里云服务器）在达到5000时，自己的学生机上的mysql和redis承受不住了。。（到mysql和redis压力顶点看日志的话差不多每一个查询花费70ms-80ms）

```
Server Software:
Server Hostname:        127.0.0.1
Server Port:            8000

Document Path:          /api/v1/tags?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6IkpTIiwicGFzc3dvcmQiOiIxMjM0NTYiLCJleHAiOjE1ODQ3MTgzMzQsImlzcyI6Imdpbi1ibG9nIn0.Iu73L9-ryGbHorCeeg8HzaE0YWiMC4leNuz6gLG1DZE
Document Length:        321 bytes

Concurrency Level:      100
Time taken for tests:   10.973 seconds
Complete requests:      1000
Failed requests:        0
Total transferred:      445000 bytes
HTML transferred:       321000 bytes
Requests per second:    91.13 [#/sec] (mean)
Time per request:       1097.279 [ms] (mean)
Time per request:       10.973 [ms] (mean, across all concurrent requests)
Transfer rate:          39.60 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    1   0.5      1       2
Processing:   135  764 1496.7    240    8744
Waiting:      135  763 1496.9    240    8744
Total:        136  765 1496.7    240    8745

Percentage of the requests served within a certain time (ms)
  50%    240
  66%    449
  75%    477
  80%    677
  90%   1443
  95%   4042
  98%   8098
  99%   8431
 100%   8745 (longest request)
```



