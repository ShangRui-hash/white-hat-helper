package models

import "github.com/Ullaakut/nmap"

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

type WebItem struct {
	URL        string `json:"url"`
	Title      string `json:"title"`
	StatusCode int    `json:"status_code"`
}

type WebDetail struct {
	WebItem
	RespHeader  map[string]interface{} `json:"resp_header"`
	RespBody    string                 `json:"resp_body"`
	WAFName     string                 `json:"waf_name"`
	FingerPrint []string               `json:"fingerprint"`
	Dirs        []Dir                  `json:"dirs"`
}

type Dir struct {
	URL         string `json:"url"`
	StatusCode  int    `json:"status_code"`
	Title       string `json:"title"`
	ContentSize string `json:"content_size"`
}
type HostListItem struct {
	IP         string    `json:"ip"`
	OS         string    `json:"os"`
	DomainList []string  `json:"domain_list"`
	PortList   []Port    `json:"ports"`
	WebList    []WebItem `json:"webs"`
}

type HostDetail struct {
	IP         string       `json:"ip"`
	OS         string       `json:"os"`
	Company    string       `json:"company"`
	DomainList []string     `json:"domain_list"`
	PortList   []PortDetail `json:"ports"`
	WebList    []WebDetail  `json:"webs"`
}
