.PHONY: all clean

all: format test build

test:
	go test -v .

format:
	gofmt -w ./controllers ./dao ./middleware  ./main.go

build:
	# 设置交叉编译参数:
	# GOOS为目标编译系统, mac os则为 "darwin", window系列则为 "windows"
	GOOS="linux" GOARCH="amd64" go build  -o search_linux

clean:
	go clean -i