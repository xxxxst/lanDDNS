
package control

import (
	// "fmt"
	// "strings"
)
type MacCtl struct {
	DomainMacGroup	*DomainGroup;			// domain => mac
	// MapDomainToMac map[string] string;		// domain => mac
	MapMacToIp map[string] string;			// mac => ip
}

func CreateMacCtl() *MacCtl {
	md := &MacCtl{};
	md.DomainMacGroup = CreateDomainGroup();
	// md.MapDomainToMac = make(map[string] string);
	md.MapMacToIp = make(map[string] string);
	return md;
}

func (c *MacCtl) AddDomainTmpl(mac string, domain string) {
	// c.MapDomainToMac[domain] = mac;
	c.DomainMacGroup.AddDomainTmpl(mac, domain);
}

func (c *MacCtl) RemoveDomainTmpl(domain string) {
	c.DomainMacGroup.RemoveDomainTmpl(domain);
	// delete(c.MapDomainToMac, domain);
}

func (c *MacCtl) SetMacIp(mac string, ip string) {
	c.MapMacToIp[mac] = ip;
}

func (c *MacCtl) UpdateGroup(group *DomainGroup) {
	c.DomainMacGroup = group;

	arrDelMac := []string{};
	for key := range c.MapMacToIp {
		if _,ok := group.MapPreciseIpMatch[key]; !ok {
			arrDelMac = append(arrDelMac,key);
			continue;
		}
		if _,ok := group.MapFuzzyIpMatch[key]; !ok {
			arrDelMac = append(arrDelMac,key);
			continue;
		}
	}
	for i:=0; i < len(arrDelMac); i++ {
		delete(c.MapMacToIp, arrDelMac[i]);
	}
	// delete(c.MapDomainToMac, domain);
}

func (c *MacCtl) Test(domain string) string {
	mac := c.DomainMacGroup.Test(domain);
	if(mac == "") {
		return "";
	}

	ip, ok := c.MapMacToIp[mac];
	if(ok) {
		return ip;
	}
	return "";
}
