package core

import (
	"net"
	"fmt"
)

type TCPSender struct {
	conn net.TCPConn
	dataChan chan []byte
	exit chan bool
}

func CreateTCPSender(conn *net.TCPConn) *TCPSender {
	sender := &TCPSender{*conn, make(chan []byte), make(chan bool, 1)}
	if sender == nil {
		return nil
	}
	return sender
}

func (sender *TCPSender) Send(msg NetMsg) {
	bytes, ok := msg.GetNetBytes()
	if !ok {
		return
	}
	sender.SendBytes(bytes)
}

func (sender *TCPSender) SendBytes(bytes []byte) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	sender.dataChan <- bytes
}

func (sender *TCPSender) send(datas []byte)  bool{
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	n, err := sender.conn.Write(datas)
	if err != nil {
		fmt.Println("TcpSender send error:", n, "reason:", err)
		sender.conn.CloseWrite()
		return false
	}
	return true
}

func (sender *TCPSender) Start() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	for {
		select {
		case data := <- sender.dataChan:
			if !sender.send(data) {
				return
			}
		case <-sender.exit:
			close(sender.dataChan)
			for data := range sender.dataChan {
				sender.send(data)
			}
			sender.conn.CloseWrite()
		}
	}
}

func (sender *TCPSender) Close() {
	defer func() {
		if x := recover(); x != nil {
			fmt.Println("TcpSender Close failed", x)
		}
	}()
	sender.exit <- true
}
