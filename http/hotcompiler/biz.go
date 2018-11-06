package hotcompiler

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"log"
)

type config struct {
	//执行的app名字，默认当前目录文字
	AppName string `yaml:"appname"`
	//指定ouput执行的程序路径
	Output string `yaml:"output"`
	//需要追加监听的文件后缀名字，默认是'.go'，
	WatchExts []string `yaml:"watch_exts"`
	//需要追加监听的目录，默认是当前文件夹，
	WatchPaths []string `yaml:"watch_paths"`
	//执行时的额外参数
	CmdArgs []string `yaml:"cmd_args"`
	//执行时追加的环境变量
	Envs []string `yaml:"envs"`
	//vendor 目录下的文件是否也监听
	VendorWatch bool `yaml:"vendor_watch"`
	//不需要监听的目录
	ExcludedPaths []string `yaml:"excluded_paths"`
	//需要编译的包或文件,优先使用-p参数
	BuildPkg string `yaml:"build_pkg"`
	//在go build 时期接收的-tags参数
	BuildTags string `yaml:"build_tags"`
}


var state sync.Mutex
var cmd *exec.Cmd
var started = make(chan bool)
//var cfg *config


func Autobuild(files []string,currpath string) {
	state.Lock()
	defer state.Unlock()

	log.Println("Start building...\n")

	//currpath, _ = os.Getwd()

	os.Chdir(currpath)

	cmdName := "go"

	var err error

	args := []string{"build"}
	args = append(args, "-o", "http")   // "http" 是 	//硬编码
	//if cfg.BuildTags != "" {
	//	args = append(args, "-tags", cfg.BuildTags)
	//}
	args = append(args, files...)

	log.Println(cmdName,args)

	bcmd := exec.Command(cmdName, args...)
	bcmd.Env = append(os.Environ(), "GOGC=off")
	bcmd.Stdout = os.Stdout
	bcmd.Stderr = os.Stderr
	err = bcmd.Run()

	if err != nil {
		log.Fatal("============== Build failed ===================\n")
		return
	}
	log.Println("Build was successful\n")

	//硬编码
	Restart("http")
}

func Kill() {
	defer func() {
		if e := recover(); e != nil {
			fmt.Println("Kill.recover -> ", e)
		}
	}()

	log.Println("kill running process")
	if cmd != nil && cmd.Process != nil {
		log.Println("killer...")
		err := cmd.Process.Kill()
		if err != nil {
			fmt.Println("Kill -> ", err)
		}
	}
}

func Restart(appname string) {
	// 原理是先杀死 appname 然后在启动 【主要循环多次也是ok的；第一次killer是进不去的！】
	Kill()
	go Start(appname)
}

func Start(appname string) {

	log.Printf("Restarting %s ...\n", appname)
	if strings.Index(appname, "./") == -1 {
		appname = "./" + appname
	}

	cmd = exec.Command(appname)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Args = append([]string{appname})
	//cmd.Args = append([]string{appname}, cfg.CmdArgs...)
	//cmd.Env = append(os.Environ(), cfg.Envs...)
	cmd.Env = append(os.Environ())

	go cmd.Run()
	log.Printf("%s is running...\n", appname)

	started <- true

}

