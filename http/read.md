
### 支持 热重启 ！

打开两个终端

一个终端输入：ps -ef|grep 应用名

一个终端输入请求：curl "http://127.0.0.1:8080/xxxx"

热升级

kill -HUP 进程 ID

打开一个终端输入请求：curl "http://127.0.0.1:8080/xxxx"



// pkill xxxx  会发送 kill -TERM 信号
// 会稍等10秒钟后才会退出！ps aux 可以验证




### 可以使用另外一个进程进行热加载！ ./hot 进程！ 实验成功！