package models

import (
	"errors"
	"fmt"
	"net"
	"regexp"
)

type ParamAddTask struct {
	baseParam
	CompanyID int      `json:"company_id" binding:"required"`
	Name      string   `json:"name" binding:"required"`
	ScanArea  []string `json:"scan_area" binding:"required"`
}

func (p *ParamAddTask) Validate() error {
	if len(p.ScanArea) == 0 {
		return errors.New("scan area is empty")
	}
	for i := range p.ScanArea {
		if IsIP(p.ScanArea[i]) {
			continue
		} else if IsCIDR(p.ScanArea[i]) {
			continue
		} else if IsDomain(p.ScanArea[i]) {
			continue
		} else {
			return fmt.Errorf("scan area::%s is invalid", p.ScanArea[i])
		}
	}
	return nil
}

func IsIP(ip string) bool {
	return net.ParseIP(ip) != nil
}

func IsCIDR(cids string) bool {
	_, _, err := net.ParseCIDR(cids)
	return err == nil
}

func IsDomain(domain string) bool {
	domainReg := regexp.MustCompile(`^([a-zA-Z0-9]([a-zA-Z0-9\\-]{0,61}[a-zA-Z0-9])?\\.)+[a-zA-Z]{2,6}(/)`)
	return domainReg.MatchString(domain)
}
