package main

import (
	"fmt"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
	"github.com/goinggo/mapstructure"
)

//topics:1.升级容器 2.安装容器 3.启停，删除容器 4.配置容器 5.容器状态查询 6.容器状态上报
const (
	CMD_CON_INSTALL    = "CMD_CON_INSTALL"
	CMD_STATUS_QUERY   = "CMD_STATUS_QUERY" //安装和升级状态查看
	CMD_CON_START      = "CMD_CON_START"
	CMD_CON_STOP       = "CMD_CON_STOP"
	CMD_CON_REMOVE     = "CMD_CON_REMOVE"
	CMD_CON_SET_CONFIG = "CMD_CON_SET_CONFIG"
	CMD_CON_GET_CONFIG = "CMD_CON_GET_CONFIG"
	CMD_CON_STATUS     = "CMD_CON_STATUS"
	REP_CON_STATUS     = "REP_CON_STATUS"  //data主题
	EVENT_CON_ALARM    = "EVENT_CON_ALARM" //data主题
	CMD_CON_UPGRADE    = "CMD_CON_UPGRADE"
	REP_JOB_RESULT     = "REP_JOB_RESULT" //安装和升级结果上报，data主题
	CMD_CON_LOG        = "CMD_CON_LOG"
)
const (
	EDGE_CMD   = "/v1/appName/service/command"
	EDGE_REPLY = "/v1/appName/service/reply"
	EDGE_DATA  = "/v1/appName/service/data"
)

var waitchan = make(chan bool)

type ContainerAPI struct {
	LoggingClient logger.LoggingClient
}

type ServiceDataParam struct {
	Cmd      string      `json:"cmd"`
	DeviceId string      `json:"deviceId"`
	Data     interface{} `json:"data"` //todo,数据以物模型规范的数据格式上报json
}
type ServiceData struct {
	Type      string           `json:"type,omitempty"`
	Mid       int64            `json:"mid,omitempty"`
	DeviceId  string           `json:"deviceId,omitempty"`
	Timestamp int64            `json:"timestamp,omitempty"`
	Expire    int              `json:"expire,omitempty"`
	Param     ServiceDataParam `json:"param,omitempty"`
}

type ServiceCMDParam struct {
	Cmd   string      `json:"cmd,omitempty"`
	Paras interface{} `json:"paras,omitempty"` //todo,数据以物模型规范的数据格式上报json
}
type ServiceCMD struct {
	Type      string          `json:"type,omitempty"`
	Mid       int64           `json:"mid,omitempty"`
	DeviceId  string          `json:"deviceId,omitempty"`
	Timestamp int64           `json:"timestamp,omitempty"`
	Expire    int             `json:"expire,omitempty"`
	App       string          `json:"app,omitempty"`
	Param     ServiceCMDParam `json:"param,omitempty"`
}

type ServiceReplyParam struct {
	Cmd   string      `json:"cmd,omitempty"`
	Paras interface{} `json:"paras,omitempty"` //todo,数据以物模型规范的数据格式上报json
}
type ServiceReply struct {
	Type      string            `json:"type,omitempty"`
	Mid       int64             `json:"mid,omitempty"`
	DeviceId  string            `json:"deviceId,omitempty"`
	Timestamp int64             `json:"timestamp,omitempty"`
	Code      int               `json:"code,omitempty"`
	Msg       string            `json:"msg,omitempty"`
	Param     ServiceReplyParam `json:"param,omitempty"`
}

///////////////////////
//install
type ConInstallCmd struct {
	Container string `json:"container,omitempty"`
	Image     string `json:"image,omitempty"`
	// WithAPP   string   `json:"with_app,omitempty"`
	CfgCpu  CfgCpu_struct  `json:"cfg_cpu,omitempty"`
	CfgMem  CfgMem_struct  `json:"cfg_mem,omitempty"`
	CfgDisk CfgDisk_struct `json:"cfg_disk,omitempty"`
	Port    string         `json:"port,omitempty"`
	Mount   []string       `json:"mount,omitempty"`
	Dev     []string       `json:"dev,omitempty"`
	JobId   int            `json:"job_id,omitempty"`
	Policy  int            `json:"policy,omitempty"`
}
type ConInstallReply struct {
}

//start,stop,remove

//config-set,get

//upgrade
//log

//cup字段：
type CfgCpu_struct struct {
	Cpus      int    `json:"cpus,omitempty"`
	Frequency int    `json:"frequency,omitempty"`
	Cache     int    `json:"cache,omitempty"`
	Arch      string `json:"arch,omitempty"`
	CpuLmt    int    `json:"cpu_lmt,omitempty"`
}

//mem字段：
type CfgMem_struct struct {
	Phy    int `json:"phy,omitempty"`
	Virt   int `json:"virt,omitempty"`
	MemLmt int `json:"mem_lmt,omitempty"`
}

//disk字段：
type CfgDisk_struct struct {
	Disk    int `json:"disk,omitempty"`
	DiskLmt int `json:"disk_lmt,omitempty"`
}

//define a function for the default message handler
var ServiceDataHandler MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

func (ca *ContainerAPI) FUNC_CMD_CON_INSTALL(servicecmd *ServiceCMD) *ServiceReply {

	var servicereply = new(ServiceReply)
	var coninstallcmd = new(ConInstallCmd)

	if err := mapstructure.Decode(servicecmd.Param.Paras, coninstallcmd); err != nil {
		fmt.Println(err)
	}
	fmt.Println(coninstallcmd)

	ca.LoggingClient.Debug(coninstallcmd.Image)
	ca.LoggingClient.Debug(coninstallcmd.Container)
	ca.LoggingClient.Debug(coninstallcmd.CfgCpu.Arch)
	//安装容器

	servicereply.Mid = servicecmd.Mid
	servicereply.Timestamp = time.Now().Unix()
	servicereply.Type = servicecmd.Type
	servicereply.DeviceId = servicecmd.DeviceId
	servicereply.Code, servicereply.Msg = gd.conInstall(coninstallcmd.Container, coninstallcmd.Image)

	return servicereply
}

// func (ca *ContainerAPI) FUNC_CMD_STATUS_QUERY() {

// }
func (ca *ContainerAPI) FUNC_CMD_CON_START(servicecmd *ServiceCMD) *ServiceReply {
	var servicereply = new(ServiceReply)
	var coninstallcmd = new(ConInstallCmd)

	if err := mapstructure.Decode(servicecmd.Param.Paras, coninstallcmd); err != nil {
		fmt.Println(err)
	}
	servicereply.Mid = servicecmd.Mid
	servicereply.Timestamp = time.Now().Unix()
	servicereply.Type = servicecmd.Type
	servicereply.DeviceId = servicecmd.DeviceId
	servicereply.Code, servicereply.Msg = gd.conStart(coninstallcmd.Container)

	return servicereply
}
func (ca *ContainerAPI) FUNC_CMD_CON_STOP(servicecmd *ServiceCMD) *ServiceReply {
	var servicereply = new(ServiceReply)
	var coninstallcmd = new(ConInstallCmd)

	if err := mapstructure.Decode(servicecmd.Param.Paras, coninstallcmd); err != nil {
		fmt.Println(err)
	}
	servicereply.Mid = servicecmd.Mid
	servicereply.Timestamp = time.Now().Unix()
	servicereply.Type = servicecmd.Type
	servicereply.DeviceId = servicecmd.DeviceId
	servicereply.Code, servicereply.Msg = gd.conStop(coninstallcmd.Container)

	return servicereply
}
func (ca *ContainerAPI) FUNC_CMD_CON_REMOVE(servicecmd *ServiceCMD) *ServiceReply {
	var servicereply = new(ServiceReply)
	var coninstallcmd = new(ConInstallCmd)

	if err := mapstructure.Decode(servicecmd.Param.Paras, coninstallcmd); err != nil {
		fmt.Println(err)
	}
	servicereply.Mid = servicecmd.Mid
	servicereply.Timestamp = time.Now().Unix()
	servicereply.Type = servicecmd.Type
	servicereply.DeviceId = servicecmd.DeviceId
	servicereply.Code, servicereply.Msg = gd.conRemove(coninstallcmd.Container)

	return servicereply
}

func (ca *ContainerAPI) FUNC_CMD_CON_SET_CONFIG() {

}
func (ca *ContainerAPI) FUNC_CMD_CON_GET_CONFIG() {

}
func (ca *ContainerAPI) FUNC_CMD_CON_STATUS() {

}
func (ca *ContainerAPI) FUNC_CMD_CON_UPGRADE() {

}
func (ca *ContainerAPI) FUNC_CMD_CON_LOG() {

}