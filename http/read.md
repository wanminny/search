
### 支持 热重启 ！

打开两个终端

一个终端输入：ps -ef|grep 应用名

一个终端输入请求：curl "http://127.0.0.1:8080/xxxx"

热升级

kill -HUP 进程 ID

打开一个终端输入请求：curl "http://127.0.0.1:8080/xxxx"



