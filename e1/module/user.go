package module

import (
	"fmt"
	. "e1/core"
	. "e1/mgr"
	"e1/mgr"
)

const (
	USER_STATUS_INIT int16 = 0
	USER_STATUS_ONLINE int16 = 100
	USER_STATUS_OFFLINE int16 = 200
)

type User struct {
	ID ObjectID
	recvChan chan *Command
	innerChan chan *Command
	netChan chan *Command
	Sender *TCPSender
	Status int16
}

func CreateUser(id ObjectID) *User{
	user := &User{}
	user.ID = id
	user.recvChan = make(chan *Command)
	user.innerChan = make(chan *Command)
	RegisterChan(id, user.innerChan)

	user.Status = USER_STATUS_INIT
	go startRecv(user)

	return user
}

func startRecv(user *User) {
	for {
		select {
		case msg:=<-user.recvChan:
			if msg == nil && user.Status == USER_STATUS_ONLINE{
				return
			}
			user.processClientMessage(msg)
		case msg:= <-user.innerChan:
			if msg == nil && user.Status == USER_STATUS_ONLINE{
				return
			}
			user.processInnerMessage(msg)
		}
	}
}

func (user *User) processClientMessage(msg *Command) {
	if msg == nil {
		return
	}
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("User processClientMsg failed:", err, " cmd:", msg.Cmd)
		}
	}()
	netMsg := CreateMessage(msg)
	netMsg.Process(user)
}

func (user *User) processInnerMessage(msg *Command) {
	 if msg == nil {
		 return
	 }
	switch msg.Cmd{
	case CMD_USER_LOGIN:
		user.userLogin(msg)
	case CMD_USER_OFFLINE:
		user.userOffline(msg)
	}
}

func (user *User) userLogin(msg *Command) {
	user.Status = USER_STATUS_ONLINE
	user.netChan = msg.RetChan
	user.Sender = msg.OtherInfo.(*TCPSender)

	msg.RetChan = user.recvChan
	user.netChan <- msg
}

func (user *User) userOffline(msg *Command) {
	if msg.RetChan != user.netChan {
		return
	}
	mgr.UnRegisterChan(user.ID)
	close(user.recvChan)
	close(user.innerChan)
	user.Sender.Close()
	user.Status = USER_STATUS_ONLINE
}
