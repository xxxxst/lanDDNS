
package control

import (
	// "fmt"
)
type DomainGroup struct {
	MapPreciseIpMatch map[string] []string;			// ip => doamin
	MapPreciseDomainMatch map[string] string;		// domain => ip
	// ArrFuzzyMatch []DomainMatch;
	MapFuzzyIpMatch map[string] []*DomainMatch;		// ip => Match
	MapFuzzyDomainMatch map[string] *DomainMatch;	// domain => Match
}

func CreateDomainGroup() *DomainGroup {
	md := DomainGroup{};
	md.MapPreciseIpMatch = make(map[string] []string);
	md.MapPreciseDomainMatch = make(map[string] string);
	// md.ArrFuzzyMatch = []DomainMatch{};
	md.MapFuzzyIpMatch = make(map[string] []*DomainMatch);
	md.MapFuzzyDomainMatch = make(map[string] *DomainMatch);
	return &md;
}

func (c *DomainGroup) setDomain(ip string, domainTmpl string) {
	arr, ok := c.MapPreciseIpMatch[ip];
	if(!ok) {
		arr = []string{};
	}
	arr = append(arr, domainTmpl);
	c.MapPreciseIpMatch[ip] = arr;
}

func (c *DomainGroup)setMatch(ip string, mat *DomainMatch) {
	arr, ok := c.MapFuzzyIpMatch[ip];
	if(!ok) {
		arr = []*DomainMatch{};
	}
	arr = append(arr, mat);
	c.MapFuzzyIpMatch[ip] = arr;
}

func (c *DomainGroup) AddDomainTmpl(ip string, domainTmpl string) {
	// c.RemoveIp(ip);
	c.RemoveDomainTmpl(domainTmpl);

	mat := CreateDomainMatch(ip, domainTmpl);
	if(!mat.IsFuzzyMatch()) {
		c.MapPreciseDomainMatch[domainTmpl] = ip;
		c.setDomain(ip, domainTmpl);
		// c.MapPreciseIpMatch[ip] = domainTmpl;
		return;
	}
	// c.ArrFuzzyMatch = append(c.ArrFuzzyMatch, mat);
	// c.MapFuzzyIpMatch[ip] = &mat;
	c.setMatch(ip, &mat);
	c.MapFuzzyDomainMatch[domainTmpl] = &mat;
}

func (c *DomainGroup) RemoveIp(ip string) {
	arrDomain, ok := c.MapPreciseIpMatch[ip];
	if(ok) {
		delete(c.MapPreciseIpMatch, ip);
		for i:=0; i<len(arrDomain); i++ {
			delete(c.MapPreciseDomainMatch, arrDomain[i]);
		}
		return;
	}
	arrMat, ok2 := c.MapFuzzyIpMatch[ip];
	if(ok2) {
		delete(c.MapFuzzyIpMatch, ip);
		for i:=0; i<len(arrMat); i++ {
			delete(c.MapFuzzyDomainMatch, arrMat[i].GetDomainTmpl());
		}
		// delete(c.MapFuzzyDomainMatch, mat.GetDomainTmpl());
	}
}

func (c *DomainGroup) RemoveDomainTmpl(domain string) {
	ip, ok := c.MapPreciseDomainMatch[domain];
	if(ok) {
		// delete(c.MapPreciseIpMatch, ip);
		delete(c.MapPreciseDomainMatch, domain);

		arr, subOk := c.MapPreciseIpMatch[ip];
		if(subOk) {
			for i:=0; i<len(arr); i++ {
				if(arr[i] == domain) {
					arr[i] = arr[len(arr) - 1];
					arr[len(arr) - 1] = "";
					arr = arr[:len(arr)-1];
					if(len(arr) == 0) {
						delete(c.MapPreciseIpMatch, ip);
					} else {
						c.MapPreciseIpMatch[ip] = arr;
					}
					break;
				}
			}
		}
		return;
	}
	mat, ok2 := c.MapFuzzyDomainMatch[domain];
	if(ok2) {
		// delete(c.MapFuzzyIpMatch, mat.GetIp());
		delete(c.MapFuzzyDomainMatch, domain);

		arr, subOk := c.MapFuzzyIpMatch[ip];
		if(subOk) {
			for i:=0; i<len(arr); i++ {
				if(arr[i].GetDomainTmpl() == domain) {
					arr[i] = arr[len(arr) - 1];
					arr[len(arr) - 1] = nil;
					arr = arr[:len(arr)-1];
					if(len(arr) == 0) {
						delete(c.MapFuzzyIpMatch, mat.GetIp());
					} else {
						c.MapFuzzyIpMatch[ip] = arr;
					}
					break;
				}
			}
		}
	}
}

func (c *DomainGroup) Test(domain string) string {
	ip, ok := c.MapPreciseDomainMatch[domain];
	if(ok) {
		return ip;
	}

	// arr := &c.ArrFuzzyMatch;
	// for i:=0; i < len(c.ArrFuzzyMatch); i++ {
	// 	ip = c.ArrFuzzyMatch[i].Test(domain);
	// 	if(ip != "") {
	// 		return ip;
	// 	}
	// }
	for key := range c.MapFuzzyDomainMatch {
		it := c.MapFuzzyDomainMatch[key];
		ip = it.Test(domain);
		if(ip != "") {
			return ip;
		}
	}

	return "";
}
