package hackflow

import (
	"log"

	"github.com/miekg/dns"
)

//dnsClient 客户端
type dnsClient struct {
	serverAddr string
}

//NewDnsClient 构造函数
func NewDNSClient(serverAddr string) *dnsClient {
	return &dnsClient{
		serverAddr: serverAddr,
	}
}

//LookupA 查询A记录
func (c *dnsClient) LookupA(fqdn string) ([]string, error) {
	var m dns.Msg
	m.SetQuestion(dns.Fqdn(fqdn), dns.TypeA)
	in, err := dns.Exchange(&m, c.serverAddr)
	if err != nil {
		log.Println("dns.Exchange failed,err:", err)
		return nil, err
	}
	ips := make([]string, 0, len(in.Answer))
	for i := range in.Answer {
		if a, ok := in.Answer[i].(*dns.A); ok {
			ips = append(ips, a.A.String())
		}
	}
	return ips, nil
}

//LookupCNAME 查询CNAME记录
func (c *dnsClient) LookupCNAME(fqdn string) ([]string, error) {
	var m dns.Msg
	m.SetQuestion(dns.Fqdn(fqdn), dns.TypeCNAME)
	in, err := dns.Exchange(&m, c.serverAddr)
	if err != nil {
		log.Println("dns.Exchange failed,err:", err)
		return nil, err
	}
	fqdnList := make([]string, 0, len(in.Answer))
	for i := range in.Answer {
		if c, ok := in.Answer[i].(*dns.CNAME); ok {
			fqdnList = append(fqdnList, c.Target)
		}
	}
	return fqdnList, nil
}

//Lookup 根据域名查询其对应的ip地址
func (c *dnsClient) Lookup(fqdn string) ([]string, error) {
	for {
		cnames, err := c.LookupCNAME(fqdn)
		if err != nil {
			log.Println("c.LookupCNAME failed,err:", err)
			return nil, nil
		}
		if len(cnames) == 0 {
			break
		}
		fqdn = cnames[0]
	}
	return c.LookupA(fqdn)
}
