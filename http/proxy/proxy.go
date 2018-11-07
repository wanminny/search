package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

/*
// 1. 直接透传 ；采用将请求继续发给真正处理的机器； eg :

```
var httpClient = http.Client{}
func transferDirect(w http.ResponseWriter,r *http.Request){

	req := http.NewRequest(method,url,bytes.NewBuffer(string body))
	req.Header = r.Header
	resp,err := httpClient.Do(req)
}
```

2. 视频流 播放 【简单处理】
 	// 视频流；文件流等
	//f,err := os.Open("")

	//w.Header.Set("Content-Type","video/mp4")
	//http.ServeContent(f)
}

3. 使用代理 ProxyTransfer【比如有些 raw流是无法使用透传的;比如文件上传？视频上传】
//
*/

func ProxyTransfer(w http.ResponseWriter,r * http.Request) {

	// 1. http://cc.com/a/b/c to ==>  http://127.0.0.1:12345/a/b/c
	// 2. 某些跨域情况跨域使用这个做简单处理 ;
	//后端真实的服务地址
	httpURl := "http://127.0.0.1:12345"
	url, _ := url.Parse(httpURl)
	proxy := httputil.NewSingleHostReverseProxy(url)

	proxy.ServeHTTP(w, r)
}
