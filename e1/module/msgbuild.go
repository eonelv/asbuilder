package module

import (
	"reflect"
	"fmt"
	. "e1/core"
	"e1/utils"
	"e1/mgr"
	"os/exec"
	"e1/log"
	"e1/cfg"
)

const (
	QUERY uint16 = 1
	BUILD uint16 = 2
	BUILD_SP uint16 = 3
)

type MsgBuild struct {
	Action uint16
	Num byte
	PData []byte
}

type Project struct {
	ID uint64
	Name NAME_STRING
	Version NAME_STRING
	Builder NAME_STRING
}

type BuildInfo struct {
	ID uint64
	Result uint16
}

func (this *MsgBuild) GetNetBytes() ([]byte, bool) {
	return GenNetBytes(uint16(CMD_BUILD), reflect.ValueOf(this))
}

func (this *MsgBuild) CreateByBytes(bytes []byte) (bool, int) {
	return utils.Byte2Struct(reflect.ValueOf(this), bytes)
}

func (this *MsgBuild) Process(p interface{}) {
	puser, ok := p.(*User)
	if !ok {
		return
	}
	switch this.Action{
	case QUERY:
		this.query(puser)
	case BUILD:
		this.build(puser, cfg.ServerCfg[cfg.BUILD_CMD])
	case BUILD_SP:
		this.build(puser, cfg.ServerCfg[cfg.BUILD_CMD_SP])
	}
}

func (this *MsgBuild) build(user *User, cmds string) {
	var err error
	project := &Project{}
	utils.Byte2Struct(reflect.ValueOf(project), this.PData)
	rows, err1 := mgr.DBMgr.PreQuery("select pname_en, pvname_en, isBuilding from t_vb_project where id = ?", project.ID)
	if err1 != nil {
		return
	}
	isBuilding := rows[0].GetBoolean("isBuilding")
	if isBuilding {
		return
	}
	projectName := rows[0].GetString("pname_en")
	buildName := rows[0].GetString("pvname_en")

	cmds =  cfg.ServerCfg[cfg.SERVER_HOME] + "build/" +projectName + "/" + buildName + "/" + cmds;

	fmt.Println(cmds)
	_, err = mgr.DBMgr.PreExecute("update t_vb_project set isBuilding = 1 where id = ?", project.ID)

	defer func() {
		_, err = mgr.DBMgr.PreExecute("update t_vb_project set isBuilding = 0 where id = ?", project.ID)
		if err != nil {
			return
		}
	}()

	if err != nil {
		fmt.Println(err)
		return
	}

	msgReturn := &MsgBuild{}
	msgReturn.Action = this.Action
	msgReturn.Num = this.Num

	msgBuildInfo := &BuildInfo{}
	msgBuildInfo.ID = project.ID
	msgBuildInfo.Result = 1

	tempData, ok := utils.Struct2Bytes(reflect.ValueOf(msgBuildInfo))
	if !ok {
		return
	}
	msgReturn.PData = tempData
	UserMgr.BroadcastMessage(msgReturn)
	//执行编译
	cmd := exec.Command(cmds, "", "")

	var bytes []byte
	bytes, err = cmd.Output()
	if err == nil {
		log.LoggerSys.Println(string(bytes))
		fmt.Println(string(bytes))
	} else {
		log.LoggerSys.Println(err, string(bytes))
	}
	msgBuildInfo.Result = 2
	tempData, ok = utils.Struct2Bytes(reflect.ValueOf(msgBuildInfo))
	if !ok {
		return
	}
	msgReturn.PData = tempData
	UserMgr.BroadcastMessage(msgReturn)
	fmt.Println("编译完成", cmds)
}

func (this *MsgBuild) query(user *User) {
	rows, err := mgr.DBMgr.PreQuery("select id, pname, pvname, isBuilding from t_vb_project")
	if err != nil {
		fmt.Println(err)
		return
	}
	this.Num = byte(len(rows))
	var totalData []byte = []byte{}
	for _,v := range rows {
		project := &Project{}
		project.ID = v.GetUint64("id")
		utils.CopyArray(reflect.ValueOf(&project.Name), []byte(v.GetString("pname")))
		utils.CopyArray(reflect.ValueOf(&project.Version), []byte(v.GetString("pvname")))
		if v.GetBoolean("isBuilding") {
			utils.CopyArray(reflect.ValueOf(&project.Builder), []byte(user.ID))
		}
		data,_ := utils.Struct2Bytes(reflect.ValueOf(project))
		totalData = append(totalData, data ...)
	}

	this.PData = totalData
	user.Sender.Send(this)
}

