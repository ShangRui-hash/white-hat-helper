package models

//CompanyStat 企业的统计信息
type CompanyStat struct {
	HostCount    uint32 `json:"host_count"`
	DomainCount  uint32 `json:"domain_count"`
	WebSiteCount uint32 `json:"web_site_count"`
}
