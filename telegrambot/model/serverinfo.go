package model

import (
	"fmt"
	"time"
)

type ServerInfoMessage struct {
	Level      string
	EventName  string
	Content    string
	serverID   int32
	vendorID   int32
	vendorName string
	serverName string
	time       time.Time
}

func NewServerInfoMessage(level string, serverID int32, serverName string, vendorID int32, vendorName string) *ServerInfoMessage {
	msg := new(ServerInfoMessage)
	msg.Level = level
	msg.serverID = serverID
	msg.vendorID = vendorID
	msg.serverName = serverName
	msg.vendorName = vendorName
	return msg
}

func (sim *ServerInfoMessage) String() string {
	return fmt.Sprintf(
		//"服務ID: %v \n"+
		//	"代理ID: %v \n"+
		"代理名稱: %v \n"+
			"服務名稱: %v \n"+
			"告警時間: %v \n"+
			"告警等級: %v \n"+
			"告警項目: %v \n"+
			"訊息: %v \n",
		sim.vendorName, sim.serverName, sim.time.String(), sim.Level, sim.EventName, sim.Content)
}

func (sim *ServerInfoMessage) SetServerID(id int32) {
	sim.serverID = id
}
func (sim *ServerInfoMessage) SetServerName(name string) {
	sim.serverName = name
}
func (sim *ServerInfoMessage) SetTime() {
	sim.time = time.Now()
}
func (sim *ServerInfoMessage) SetVendorID(id int32) {
	sim.vendorID = id
}
func (sim *ServerInfoMessage) SetVendorName(name string) {
	sim.vendorName = name
}
