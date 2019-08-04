package main

import (
	"flag"
	"github.com/cnlh/nps/vender/github.com/astaxie/beego"
	"github.com/cnlh/nps/vender/github.com/astaxie/beego/logs"
	"lib/common"
	"lib/crypt"
	"lib/daemon"
	"lib/file"
	"lib/install"
	"lib/server"
	"lib/server/connection"
	"lib/server/test"
	"lib/server/tool"
	"lib/version"
	"log"
	"os"
	"path/filepath"
	_ "web/routers"
)



func main() {
	flag.Parse()//参数初始化
	beego.LoadAppConfig("ini", filepath.Join(common.GetRunPath(), "conf", "nps.conf"))
	if len(os.Args) > 1 {//使程序执行不同的操作
		switch os.Args[1] {
		//执行测试
		case "test":
			test.TestServerConfig()
			log.Println("test ok, no error")
			return
		//程序正常执行
		case "start", "restart", "stop", "status", "reload"://执行，重启，停止，重装
			daemon.InitDaemon("nps", .GetRunPath(), .GetTmpPath())
		//安装到系统文件下
		case "install":
			install.InstallNps()
			return
		}
	}
	//日志等级
	if level = beego.AppConfig.String("log_level"); level == "" {
		level = "7"
	}
	logs.Reset()//重置日志配置
	//日志收集
	logs.EnableFuncCallDepth(true)
	logs.SetLogFuncCallDepth(3)

	//标准日志
	if *logType == "stdout" {
		logs.SetLogger(logs.AdapterConsole, `{"level":`+level+`,"color":true}`)
	} else {
		logs.SetLogger(logs.AdapterFile, `{"level":`+level+`,"filename":"`+beego.AppConfig.String("log_path")+`","daily":false,"maxlines":100000,"color":true}`)
	}


	task := &file.Tunnel{
		Mode: "webServer",
	}

	//客户端连接接口测试
	bridgePort, err := beego.AppConfig.Int("bridge_port")
	if err != nil {
		logs.Error("Getting bridge_port error", err)
		os.Exit(0)
	}
	//版本号声明
	logs.Info("the version of server is %s ,allow client version to be %s", version.VERSION, version.GetVersion())
	//配置文件内端口冲突检测
	connection.InitConnectionService()
	//加密路径申明
	crypt.InitTls(filepath.Join(.GetRunPath(), "conf", "server.pem"), filepath.Join(.GetRunPath(), "conf", "server.key"))
	//初始化所有端口
	tool.InitAllowPort()
	//日志打印
	tool.StartSystemInfo()
	//根据客户端通信接口建立新的总server
	server.StartNewServer(bridgePort, task, beego.AppConfig.String("bridge_type"))
}