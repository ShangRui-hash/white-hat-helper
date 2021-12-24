package models

import (
	"net/http"
	"web_app/pkg/hackflow"

	"github.com/Ullaakut/nmap"
)

type WebItem struct {
	URL         string   `json:"url"`
	Title       string   `json:"title"`
	Location    string   `json:"location"`
	StatusCode  int      `json:"status_code"`
	FingerPrint []string `json:"fingerprint"`
}

type WebDetail struct {
	WebItem
	WAFName        string                         `json:"waf_name"`
	RespBody       string                         `json:"resp_body"`
	RespHeader     http.Header                    `json:"resp_header"`
	ToolStatusList []ToolStatus                   `json:"tool_status_list"`
	Dirs           []hackflow.BruteForceURLResult `json:"dirs"`
}

type ToolStatus struct {
	Name   string `json:"name"`
	Status bool   `json:"status"`
}

type Dir struct {
	URL         string `json:"url"`
	StatusCode  int    `json:"status_code"`
	Title       string `json:"title"`
	Location    string `json:"location"`
	ContentSize string `json:"content_size"`
}

type HostListItem struct {
	IP          string `json:"ip"`
	OS          string `json:"os"`
	*PortInfo   `json:"port_info"`
	*WebInfo    `json:"web_info"`
	*DomainInfo `json:"domain_info"`
}

type Port struct {
	Port    int    `json:"port"`
	Service string `json:"service"`
	Status  string `json:"status"`
}

type PortDetail struct {
	Port
	Version  string        `json:"version"`
	Protocol string        `json:"protocol"`
	Script   []nmap.Script `json:"script"`
}

type PortInfo struct {
	Total    int    `json:"total"` //端口总数
	PortList []Port `json:"ports"` //部分端口
}

type WebInfo struct {
	Total   int       `json:"total"` //web服务总数
	WebList []WebItem `json:"webs"`  //部分web服务
}

type DomainInfo struct {
	Total      int      `json:"total"` //域名总数
	DomainList []string `json:"domains"`
}

type HostDetail struct {
	PortList []PortDetail `json:"ports"`
	WebList  []WebDetail  `json:"webs"`
}

type HostBaseInfo struct {
	IP         string   `json:"ip"`
	OS         string   `json:"os"`
	Company    string   `json:"company"`
	DomainList []string `json:"domain_list"`
}
