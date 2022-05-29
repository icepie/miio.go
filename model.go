package miio

import (
	"fmt"
)

type PropParam struct {
	Did   string `json:"did"`
	Siid  int    `json:"siid"`
	Piid  int    `json:"piid"`
	Value any    `json:"value,omitempty"`
}

func (p *PropParam) SetDid(did string) *PropParam {
	if did == "" {
		p.Did = fmt.Sprintf("%d-%d", p.Siid, p.Piid)
	} else {
		p.Did = did
	}
	return p
}

type PropParams []PropParam

type ActionParam struct {
	Did  string `json:"did"`
	Siid int    `json:"siid"`
	Aiid int    `json:"aiid"`
	In   []any  `json:"in,omitempty"`
	Out  []any  `json:"out,omitempty"`
}

func (a *ActionParam) SetDid(did string) *ActionParam {
	if did == "" {
		a.Did = fmt.Sprintf("%d-%d", a.Siid, a.Aiid)
	} else {
		a.Did = did
	}
	return a
}

type PropRet struct {
	PropParam
	Code    int `json:"code"`
	ExeTime int `json:"exe_time"`
}

type PropRets []PropRet

type PropParamsReq struct {
	Params PropParams `json:"params"`
	Method string     `json:"method,omitempty"`
}

// type ActionRet struct {
// 	ActionParam
// 	Miid        int `json:"miid"`
// 	Code        int `json:"code"`
// 	ExeTime     int `json:"exe_time"`
// 	WithLatency int `json:"withLatency"`
// }

// type ActionParamReq struct {
// 	Param ActionParam `json:"params"`
// }

type Ret struct {
	ID    int `json:"id"`
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
	Result  any `json:"result"`
	ExeTime int `json:"exe_time"`
}

type Info struct {
	Life      int    `json:"life"`
	Model     string `json:"model"`
	Token     string `json:"token"`
	Ipflag    int    `json:"ipflag"`
	MiioVer   string `json:"miio_ver"`
	HwVer     string `json:"hw_ver"`
	Mmfree    int    `json:"mmfree"`
	FwVer     string `json:"fw_ver"`
	Mac       string `json:"mac"`
	WifiFwVer string `json:"wifi_fw_ver"`
	Ap        struct {
		Ssid    string `json:"ssid"`
		Bssid   string `json:"bssid"`
		Rssi    int    `json:"rssi"`
		Primary int    `json:"primary"`
	} `json:"ap"`
	Netif struct {
		LocalIP string `json:"localIp"`
		Mask    string `json:"mask"`
		Gw      string `json:"gw"`
	} `json:"netif"`
	PartnerID string `json:"partner_id"`
}
