package node

import (
	"bytes"
	"fmt"
)

var (
	StatusStarted = 1
	StatusStopped = 2
	StatusError   = 3
)

type Info struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Ip          string `json:"ip"`
	Port        int    `json:"port"`
	Token       string `json:"token"`
	ParentCode  string `json:"parentCode"`
	Status      int    `json:"status"`
	StatusError string `json:"statusError"`
}

func (this_ *Info) GetNodeStr() (str string) {
	return fmt.Sprintf("节点[%s]，IP[%s]，Port[%d]", this_.Name, this_.Ip, this_.Port)
}

func (this_ *Info) checkToken(token []byte) bool {
	nodeToken := []byte(this_.Token)
	if len(nodeToken) != len(token) {
		println(fmt.Sprintf(this_.GetNodeStr() + " Token check field"))
		return false
	}
	if !bytes.Contains(token, nodeToken) {
		println(fmt.Sprintf(this_.GetNodeStr() + " Token check field"))
		return false
	}
	return true
}

type PortForwarding struct {
	InPort  int    `json:"inPort"`
	InNode  string `json:"inNode"`
	OutPort int    `json:"outPort"`
	OutNode string `json:"outNode"`
}