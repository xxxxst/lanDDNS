
package control

import (
	"fmt"
	"regexp"
	"strings"

	"gopkg.in/ini.v1"

	. "model"
	util "util"
)

type ConfigCtl struct {
	RegDomainFormat *regexp.Regexp;
	RegComment *regexp.Regexp;
	RegSpace *regexp.Regexp;
}

var insConfigCtl *ConfigCtl = nil;

func GetConfigCtl() *ConfigCtl {
	if(insConfigCtl == nil) {
		ins := &ConfigCtl{};
		
		ins.RegDomainFormat = regexp.MustCompile("(#.*)|([\\s]+)");
		ins.RegComment = regexp.MustCompile("#.*");
		ins.RegSpace = regexp.MustCompile("[\\s]+");

		insConfigCtl = ins;
	}
	return insConfigCtl;
}

func (c *ConfigCtl) LoadConfig(path string) *ComConfig {
	md := CreateComConfig();

	if(util.FileExists(path)) {
		ini.MapTo(md, path);
	} else {
		cfg := ini.Empty();
		cfg.ReflectFrom(md);
		cfg.SaveTo(path);
	}
	// cfg, err := ini.Load(path);
	// if(err == nil) {
	// 	cfg.MapTo(md);
	// } else {
	// 	cfg.SaveTo(path);
	// }

	return md;
}

func (c *ConfigCtl) formatHost(text string) *DomainGroup {
	// regComment := regexp.MustCompile("#.*");
	// regSpace := regexp.MustCompile("[\\s]+");
	// regDomain := regexp.MustCompile("^(?:\\s*)([^\\s#]+)(?:\\s+([^\\s#]+))+");

	// rst := make(map[string] string);

	group := CreateDomainGroup();

	text = strings.Replace(text, "\r\n", "\n", -1);
	arr := strings.Split(text, "\n");
	for i := 0; i < len(arr); i++ {
		str := arr[i];
		idx := strings.Index(str, "#");
		if(idx >= 0) {
			str = str[0:idx];
		}
		if(len(str) == 0) {
			continue;
		}

		str = c.RegDomainFormat.ReplaceAllString(str, " ");
		// str = c.RegComment.ReplaceAllString(str, "");
		// str = c.RegSpace.ReplaceAllString(str, " ");
		str = strings.Trim(str, " ");

		arr2 := strings.Split(str, " ");

		// if(regComment.MatchString(arr[i])) {
		// 	fmt.Println("Comment: ", arr[i]);
		// 	continue;
		// }

		// arr2 := regDomain.FindAllStringSubmatch(arr[i], -1);
		if(len(arr2) < 2) {
			continue;
		}
		// fmt.Println("rst: ", len(arr2), arr2);
		strIp := strings.ToLower(arr2[0]);
		for j:=1; j < len(arr2); j++ {
			domain := arr2[j];
			group.AddDomainTmpl(strIp, domain);
			// mat := CreateDomainMatch(strIp, domain);
			// fmt.Println("--rst:" + strIp + " - " + domain + " - ");
			// fmt.Println("--rst:-" + arr2[j] + "-");
		}
	}
	// fmt.Println(group.MapPreciseDomainMatch);
	// fmt.Println(group.ArrFuzzyMatch);

	// fmt.Println("  -", group.Test("test.lan"));
	// fmt.Println("  -", group.Test("www.test.lan"));
	// fmt.Println("  -", group.Test("aa.test.lan"));
	// fmt.Println("  -", group.Test("aa.test1.bb.lan"));
	// fmt.Println("  -", group.Test("aa.test1.lan.aa"));

	return group;
}

func (c *ConfigCtl) LoadHostStatic(path string) *DomainGroup {
	text := util.ReadFileString(path);
	return c.formatHost(text);
}

func (c *ConfigCtl) SaveHostDynamic(path string, group *DomainGroup) {
	str := "";
	endl := "\r\n";
	str += "" + endl;
	str += "# This file is generated automatically by the program, do not modify it!" + endl;
	str += "" + endl;
	for ip := range group.MapPreciseIpMatch {
		line := ip;
		arr := group.MapPreciseIpMatch[ip];
		for i:=0; i<len(arr); i++ {
			line += " " + arr[i];
		}
		str += line + endl;
	}
	str += "" + endl;
	for ip := range group.MapFuzzyIpMatch {
		line := ip;
		arr := group.MapFuzzyIpMatch[ip];
		for i:=0; i<len(arr); i++ {
			line += " " + arr[i].GetDomainTmpl();
		}
		str += line + endl;
	}
	util.SaveFileString(path, str);
}

// func (c *ConfigCtl) formatMac(text string) *DomainGroup {
// 	text := util.ReadFileString(path);
// 	return c.formatHost(text);
// }

func (c *ConfigCtl) LoadHostMac(path string) *DomainGroup {
	text := util.ReadFileString(path);
	return c.formatHost(text);
}

func testConfigCtl() {
	fmt.Print("");
}
