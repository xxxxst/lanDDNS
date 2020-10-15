
package control

import (
	// "fmt"
	"strings"
)
type DomainMatch struct {
	ip string;
	domainTmpl string;
	
	isFuzzyMatch bool;
	arrUrl []string;
	fuzzyFirst bool;
	fuzzyLast bool;
}

func CreateDomainMatch(ip string, domainTmpl string) DomainMatch {
	md := DomainMatch{};
	md.SetTemplate(ip, domainTmpl);
	// md.isFuzzyMatch = false;
	// md.fuzzyFirst = false;
	// md.fuzzyLast = false;
	// md.domainTmpl = "";
	// md.arrUrl = []string{};
	return md;
}

func (c *DomainMatch) IsFuzzyMatch() bool {
	return c.isFuzzyMatch;
}

func (c *DomainMatch) GetIp() string {
	return c.ip;
}

func (c *DomainMatch) GetDomainTmpl() string {
	return c.domainTmpl;
}

func (c *DomainMatch) SetTemplate(ip string, domainTmpl string) {
	c.ip = ip;
	c.domainTmpl = domainTmpl;
	
	c.isFuzzyMatch = false;
	c.fuzzyFirst = false;
	c.fuzzyLast = false;
	c.arrUrl = []string{};

	if(len(domainTmpl) == 0) {
		return;
	}

	funSplit := func(c rune) bool {
		return (c == '*');
	}

	// arr := strings.Split(domainTmpl, "*");
	arr := strings.FieldsFunc(domainTmpl, funSplit);
	if(len(arr) <=1 ) {
		return;
	}

	c.isFuzzyMatch = true;
	c.fuzzyFirst = (domainTmpl[0] == '*');
	c.fuzzyLast = (domainTmpl[len(domainTmpl)-1] == '*');

	c.arrUrl = arr;
	// fmt.Println(len(arr), arr);
}

func (c *DomainMatch) Test(domain string) string {
	// fmt.Println("aaa:", c.isFuzzyMatch, c.ip, domain);
	if(!c.isFuzzyMatch) {
		if(domain == c.domainTmpl) {
			return c.ip;
		}
		return "";
	}

	if(len(c.arrUrl) == 0) {
		return c.ip;
	}
	if(len(domain) == 0) {
		return "";
	}

	lastIdx := len(c.arrUrl) - 1;
	matchIdx := 0;
	for i:=0;i<len(c.arrUrl);i++ {
		idx := strings.Index(domain[matchIdx:], c.arrUrl[i]);
		idx = idx + matchIdx;
		matchIdx = idx + len(c.arrUrl[i]);
		if(idx<0) {
			return "";
		}
		if(!c.fuzzyFirst && i == 0) {
			if(idx != 0) {
				return "";
			}
		}
		if(!c.fuzzyLast && i == lastIdx) {
			if(idx + len(c.arrUrl[i]) != len(domain)) {
				return "";
			}
		}
		if(i != lastIdx && matchIdx >= len(domain)) {
			return "";
		}
	}

	return c.ip;
}
