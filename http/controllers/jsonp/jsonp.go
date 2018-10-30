package jsonp

import (
	"net/http"
	"fmt"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
)

//纯粹的显示页面 没有演示作用；需要通过html页面来处理; [没有用处]
func Jpage(w http.ResponseWriter, r *http.Request,params httprouter.Params)  {
	const indexHtml = `<!DOCTYPE html>
<html>
<head><title>Go JSONP Server</title></head>
<body>
<button id="btn">Click to get HTTP header via JSONP</button>
<pre id="result"></pre>
<script>
'use strict';
var btn = document.getElementById("btn");
var result = document.getElementById("result");

function myCallback(acptlang) {
  result.innerHTML = JSON.stringify(acptlang, null, 2);
}

function jsonp() {
  result.innerHTML = "Loading ...";
  var tag = document.createElement("script");
  tag.src = "/jsonp?callback=myCallback";
  document.querySelector("head").appendChild(tag);
}

btn.addEventListener("click", jsonp);
</script>
</body>
</html>`
	fmt.Fprintf(w, indexHtml)
	//io.WriteString(w,indexHtml)
	//w.Write()
	return
}


// jsonp
func JsonpHandler(w http.ResponseWriter, r *http.Request,params httprouter.Params) {
	callbackName := r.URL.Query().Get("callback")
	if callbackName == "" {
		fmt.Fprintf(w, "Please give callback name in query string")
		return
	}

	// 组装自己的逻辑参数【需要和前端约定！因为 是前端自己执行自己的回调函数;】
	b, err := json.Marshal(r.Header)
	if err != nil {
		fmt.Fprintf(w, "json encode error")
		return
	}

	w.Header().Set("Content-Type", "application/javascript")
	//log.Println(callbackName,string(b))

	// 将结果放回个前端;前端自己调用！是一个mycallback(params) 即callbackName(b) 的形式 前端自己处理！
	fmt.Fprintf(w, "%s(%s);", callbackName, b)
}
