package module

import (
	"net"
	. "e1/core"
	. "e1/utils"
	. "e1/mgr"
)

import (
	"fmt"
	"time"
	"io"
	"reflect"
)

const (
	STR_AUTH_REQ    = "<policy-file-request/>"
	STR_AUTH_RETURN = "<?xml version=\"1.0\"?><cross-domain-policy><site-control permitted-cross-domain-policies='all'/><allow-access-from domain=\"*\" to-ports=\"*\"/></cross-domain-policy>"
)

type TCPClient struct {
	ID ObjectID
	Conn *net.TCPConn
	Sender *TCPSender
	dataChan     chan *Command // 自身接收用的channel
	userChan chan *Command // 跟自己相关的User的接收channel
	isLogin bool
}

func ProcessRecv(client *TCPClient) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)    //这里的err其实就是panic传入的内容
		}
	}()
	conn := client.Conn
	defer conn.CloseWrite()
	defer func() {
		if client.dataChan != nil {
			close(client.dataChan)
		}
	}()
	defer client.close()
	for {
		headerBytes := make([]byte, HEADER_LENGTH)
		_, err := io.ReadFull(conn, headerBytes)
		if err != nil {
			fmt.Println("Read Data Error", client.ID, err.Error())
			break
		}

		if headerBytes[0] == STR_AUTH_REQ[0] && !client.isLogin {
			tempbuf := make([]byte, len(STR_AUTH_REQ)-int(HEADER_LENGTH))
			_, err = io.ReadFull(conn, tempbuf)
			if err != nil {
				//log
				fmt.Println("HandleUserConnect read rest auth req err", err)
				return
			}
			headerBytes = append(headerBytes, tempbuf...)
			authReq := string(headerBytes)
			if authReq == STR_AUTH_REQ {
				conn.Write(append([]byte(STR_AUTH_RETURN), 0))
				fmt.Println("send auth over.")
			} else {
				//log
				fmt.Println("recv wrong auth req:", authReq)
			}

			return
		}
		header := &PackHeader{}
		Byte2Struct(reflect.ValueOf(header), headerBytes)

		bodyBytes := make([]byte, header.Length-HEADER_LENGTH)
		_, err = io.ReadFull(conn, bodyBytes)
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		client.processClientMessage(header, bodyBytes)
	}
	fmt.Println("TCPClient cant receive data no more")
}

func (client *TCPClient)processClientMessage(header *PackHeader, datas []byte) {
	if !client.isLogin {
		client.processLogin(header, datas)
	} else {
		client.routMsgToUser(header, datas)
	}
}

func (client *TCPClient) processLogin(header *PackHeader, datas []byte) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	type MSG_CONNECT struct {
		Connet byte
	}

	client.isLogin = true
	go client.Sender.Start()

	systemChan := GetChanByID(SYSTEM_USER_CHAN_ID)
	client.dataChan = make(chan *Command)
	msgSend := &Command{CMD_SYSTEM_USER_LOGIN, client.ID, client.dataChan, nil}
	msgSend.OtherInfo = client.Sender

	select {
	case systemChan <- msgSend:
	case <-time.After(5 * time.Second):
		fmt.Println("loginUserToGame put user channel failed:", CMD_SYSTEM_USER_LOGIN)
	}
	client.waitLoginReturn()
}

func (client *TCPClient) waitLoginReturn() bool {
	msg := <-client.dataChan
	if msg.RetChan == nil {
		return false
	}
	client.userChan = msg.RetChan

	return true
}

// 将消息路由到玩家处理
func (client *TCPClient) routMsgToUser(header *PackHeader, data []byte) bool {
	msg := &Command{header.Cmd, data, nil, nil}

	select {
	case client.userChan <- msg:
	case <-time.After(5 * time.Second):
		fmt.Println("routMsgToUser put user channel failed:", client.ID, UserMgr.users[client.ID].Status)
		return false
	}

	return true
}

func (client *TCPClient) close() {
	if !client.isLogin {
		return
	}
	userInnerChan := GetChanByID(client.ID)
	closeMsg := &Command{CMD_SYSTEM_USER_OFFLINE, nil, client.dataChan, nil}
	client.isLogin = false

	select {
	case userInnerChan <- closeMsg:
	case <-time.After(5 * time.Second):
		fmt.Println("sendOffline put user channel failed:", client.ID)
		return
	}
	return
}
