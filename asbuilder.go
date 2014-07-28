package main

import (
	"net"
	"os"
	. "e1/core"
	. "e1/module"
	"e1/mgr"
	"fmt"
	"e1/cfg"
	"runtime"
	"e1/log"
)

func main() {
	Start()
}

func Start() {
	logger ,err := log.CreateLogger("goserver.log")
	if err != nil{
		fmt.Println("Create syslog file error.")
		return
	}
	log.LoggerSys = logger
	runtime.GOMAXPROCS(runtime.NumCPU())
	logger.Println("")
	logger.Println("---------------------------------")
	logger.Printf("服务器为%d核", runtime.NumCPU())
	if !cfg.LoadCfg() {
		log.LoggerSys.Print("Load server config error.")
		os.Exit(100)
	}
	if !mgr.CreateDBMgr(cfg.ServerCfg[cfg.SERVER_HOME] + "/" + cfg.ServerCfg[cfg.DB_NAME]) {
		log.LoggerSys.Print("Connect dataBase error")
		os.Exit(101)
	}
	mgr.CreateChanMgr()
	CreateUserMgr()
	sysChan := make(chan *Command)
	mgr.RegisterChan(SYSTEM_CHAN_ID, sysChan)
	go processTCP()
	logger.Println("服务器启动成功")
	for {
		select {
		case msg := <-sysChan:
			log.LoggerSys.Print("main recv msg:", msg.Cmd)
			if msg.Cmd == CMD_SYSTEM_MAIN_CLOSE {
				return
			}
		}
	}
}

func checkError(err error){
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}

func processTCP() {
	service := ":" + cfg.ServerCfg[cfg.SERVER_PORT]
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)
	for {
		conn,err := listener.AcceptTCP()
		if err != nil {
			continue
		}
		go processConnect(conn)
	}
}

func processConnect(conn *net.TCPConn){
	client := &TCPClient{}
	objID := conn.RemoteAddr().String()
	client.ID = ObjectID(objID)
	client.Conn = conn
	client.Sender = CreateTCPSender(conn)
	go ProcessRecv(client)
}






