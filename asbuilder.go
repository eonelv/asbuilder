package main

import (
	"net"
	"os"
	. "e1/core"
	. "e1/module"
	"e1/mgr"
	"e1/cfg"
	_ "e1/cfg"
	"runtime"
	. "e1/log"
	_ "e1/log"
	"fmt"
)

func main() {
	Start()
}

func Start() {

	runtime.GOMAXPROCS(runtime.NumCPU())
	cfgOK, cfgErr := cfg.LoadCfg()
	LogInfo("")
	LogInfo("---------------------------------")
	if !cfgOK {
		LogInfo("Load server config error.", cfgErr)
		os.Exit(100)
	}
	LogInfo("Load server config success.")
	if !mgr.CreateDBMgr(cfg.GetServerHome() + "/" + cfg.GetDBName()) {
		LogError("Connect dataBase error")
		os.Exit(101)
	}
	mgr.CreateChanMgr()
	CreateUserMgr()
	sysChan := make(chan *Command)
	mgr.RegisterChan(SYSTEM_CHAN_ID, sysChan)
	go processTCP()
	LogInfo("Server bootup success.")
	for {
		select {
		case msg := <-sysChan:
			LogInfo("main recv msg:", msg.Cmd)
			if msg.Cmd == CMD_SYSTEM_MAIN_CLOSE {
				return
			}
		}
	}
}

func checkError(err error){
	if err != nil {
		LogError(err)
		os.Exit(0)
	}
}

func processTCP() {
	service := fmt.Sprintf(":%d",  cfg.GetServerPort())
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






