

### 简单使用方法：


./search: -c  "RFID:3d" -start "20180925" -end "20180926" -dir /gopath/src/gobible/logmanager/cli



#### 交叉编译

GOOS=linux GOARCH=amd64 go build  -o search_linux
