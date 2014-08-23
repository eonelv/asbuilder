package module

import (
	"reflect"
	"fmt"
	. "e1/core"
	"e1/utils"
	"e1/mgr"
	"e1/cfg"
	"os/exec"
	_ "e1/log"
	. "e1/log"
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
		this.build(puser, cfg.GetCmd())
	case BUILD_SP:
		this.build(puser, cfg.GetCmdSP())
	}
}

func (this *MsgBuild) build(user *User, cmds string) {
	project := &Project{}
	utils.Byte2Struct(reflect.ValueOf(project), this.PData)
	rows, err := mgr.DBMgr.PreQuery("select pname_en, pvname_en, isBuilding from t_vb_project where id = ?", project.ID)
	if err != nil {
		return
	}
	isBuilding := rows[0].GetBoolean("isBuilding")
	if isBuilding {
		return
	}

	projectName := rows[0].GetString("pname_en")
	buildName := rows[0].GetString("pvname_en")

	var cmds2 string = ""
	if cmds == cfg.GetCmdSP() {
		cmds2 = cfg.GetServerHome() + "/build/" +projectName + "/" + buildName + "/" + cfg.GetCmdSPInner()
	}
	cmds =  cfg.GetServerHome() + "/build/" +projectName + "/" + buildName + "/" + cmds

	go execBuild(cmds, cmds2, project, user, this)

}

func execBuild(cmds string, cmds2 string, project *Project, user *User, msgBuild *MsgBuild) {
	var err error
	_, err = mgr.DBMgr.PreExecute("update t_vb_project set isBuilding = 1, builder=? where id = ?", string(user.ID), project.ID)

	defer func() {
		_, err = mgr.DBMgr.PreExecute("update t_vb_project set isBuilding = 0, builder=? where id = ?", "", project.ID)
		if err != nil {
			return
		}
	}()

	if err != nil {
		fmt.Println(err)
		return
	}

	msgReturn := &MsgBuild{}
	msgReturn.Action = msgBuild.Action
	msgReturn.Num = msgBuild.Num

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
	if cmds2 != "" {
		LogInfo("Compile sp", cmds2)
		cmd2 := exec.Command(cmds2, "", "")

		bytes2, err2 := cmd2.Output()
		if err2 == nil {
			LogInfo(string(bytes2))
		} else {
			LogInfo(err, string(bytes2))
		}
	}

	LogInfo("Compile sp", cmds)
	cmd := exec.Command(cmds, "", "")

	var bytes []byte
	bytes, err = cmd.Output()
	if err == nil {
		LogInfo(string(bytes))
	} else {
		LogInfo(err, string(bytes))
	}

	msgBuildInfo.Result = 2
	tempData, ok = utils.Struct2Bytes(reflect.ValueOf(msgBuildInfo))
	if !ok {
		return
	}
	msgReturn.PData = tempData

	defer UserMgr.BroadcastMessage(msgReturn)
	LogInfo("编译完成", cmds)
}

func (this *MsgBuild) query(user *User) {
	rows, err := mgr.DBMgr.PreQuery("select id, pname, pvname, isBuilding, builder from t_vb_project")
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
		isBuilding := v.GetBoolean("isBuilding")
		if isBuilding {
			utils.CopyArray(reflect.ValueOf(&project.Builder), []byte(v.GetString("builder")))
		}
		data,_ := utils.Struct2Bytes(reflect.ValueOf(project))
		totalData = append(totalData, data ...)
	}

	this.PData = totalData
	user.Sender.Send(this)
}

