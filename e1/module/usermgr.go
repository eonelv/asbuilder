package module

import (
	"e1/core"
	"e1/mgr"
	"time"
	. "e1/log"
)

var UserMgr UserManager
type UserManager struct {
	users map[core.ObjectID] *User
	systemChan chan *core.Command
}

func CreateUserMgr() bool {
	UserMgr = UserManager{}
	UserMgr.systemChan = make(chan *core.Command)
	UserMgr.users = make(map[core.ObjectID]*User)
	go startRecv(&UserMgr)
	return true
}

func startRecv(userMgr *UserManager) {
	mgr.RegisterChan(core.SYSTEM_USER_CHAN_ID, UserMgr.systemChan)
	defer mgr.UnRegisterChan(core.SYSTEM_USER_CHAN_ID)
	for {
		select {
		case msg := <-UserMgr.systemChan:
			userMgr.processMsg(msg)
		}
	}
}

func (this *UserManager) processMsg(msg *core.Command) {
	switch msg.Cmd{
	case core.CMD_SYSTEM_USER_LOGIN:
		this.processUserLogin(msg)
	case core.CMD_SYSTEM_BROADCAST:
		this.processBroadCast(msg.Message.(core.NetMsg))
	}
}

func (this *UserManager) processUserLogin(msg *core.Command) {
	id := msg.Message.(core.ObjectID)
	u := CreateUser(id)

	this.users[id] = u

	select {
	case u.innerChan <- msg:
	case <-time.After(10 * time.Second):
		LogError("new user busy :", id)
		return
	}
}

func (this *UserManager)processBroadCast(msg core.NetMsg) {
	for _, u := range this.users {
		u.Sender.Send(msg)
		LogError("BroadcastMessage to User," ,u.ID, u.Status)
	}
}

func (this *UserManager)BroadcastMessage(msg core.NetMsg) {
	cmd := &core.Command{core.CMD_SYSTEM_BROADCAST, msg, nil, nil}
	select {
	case this.systemChan <- cmd:
	case <- time.After(20*time.Second):
			return
	}
}


