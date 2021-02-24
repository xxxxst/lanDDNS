
package server

import (
	"fmt"
	"sync"
	"time"
	
	"github.com/rjeczalik/notify"

	. "control"
	. "model"
	util "util"
)

type ConfigListenServer struct {
	
}

var insConfigListenServer *ConfigListenServer = nil;

func GetConfigListenServer() *ConfigListenServer {
	if(insConfigListenServer == nil) {
		insConfigListenServer = &ConfigListenServer{};
	}
	return insConfigListenServer
}

func (c *ConfigListenServer) Run() {
	md := GetComModel();
	cfgCtl := GetConfigCtl();
	mainCtl := GetMainCtl();

	// fmt.Println(md.ConfigDir + "config.ini");
	// rst := cfgCtl.LoadConfig(md.ConfigDir + "config.ini");

	// hostStaticPath := md.ConfigDir + "host.static.txt";
	// mainCtl.DomainGroupStatic = cfgCtl.LoadHostStatic(hostStaticPath);
	c.updateHostStatic();
	c.updateHostMac();

	// fmt.Println(mainCtl.DomainMacCtl.DomainMacGroup);
	// fmt.Println(mainCtl.DomainMacCtl.MapMacToIp);

	hostDynamicPath := md.ConfigDir + "host.dynamic.txt";
	cfgCtl.SaveHostDynamic(hostDynamicPath, mainCtl.DomainGroupDynamic);

	if(md.ConfigMd.Server.UseStaticHost) {
		path := md.ConfigDir + "host.static.txt";
		if(!util.FileExists(path)) {
			endl := "\r\n";
			str := "" + endl +
				"# [ip] [domain] [domain] ..." + endl +
				"# domain support '*', example:" + endl +
				"# 127.0.0.1 domain.lan www.domain.lan *.domain.lan" + endl;
			util.SaveFileString(path, str);
		}
		c.watch(md.ConfigDir, "host.static.txt", func() {
			c.updateHostStatic();
		});
	}

	if(md.ConfigMd.Server.MacIp != "") {
		path := md.ConfigDir + "host.mac.txt";
		if(!util.FileExists(path)) {
			endl := "\r\n";
			str := "" + endl +
				"# [mac] [domain] [domain] ..." + endl +
				"# domain support '*', example:" + endl +
				"# 00:00:00:aa:bb:11 domain.lan www.domain.lan *.domain.lan" + endl;
			util.SaveFileString(path, str);
		}
		c.watch(md.ConfigDir, "host.mac.txt", func() {
			c.updateHostMac();
		});
	}
	
	// md.ConfigMd = rst;
}

func (c *ConfigListenServer) updateHostStatic() {
	md := GetComModel();
	cfgCtl := GetConfigCtl();
	mainCtl := GetMainCtl();
	mainCtl.DomainGroupStatic = cfgCtl.LoadHostStatic(md.ConfigDir + "host.static.txt");
}

func (c *ConfigListenServer) updateHostMac() {
	md := GetComModel();
	cfgCtl := GetConfigCtl();
	mainCtl := GetMainCtl();

	group := cfgCtl.LoadHostMac(md.ConfigDir + "host.mac.txt");
	mainCtl.DomainMacCtl.UpdateGroup(group);
}

func (c *ConfigListenServer) watch(path string, fileName string, cb func()) {
	flagNeedUpdate := false;
	var lock sync.Mutex;

	ch := make(chan notify.EventInfo, 1);
	err := notify.Watch(path, ch, notify.Create, notify.Write);
	if(err != nil) {
		// fmt.Println(err);
		return;
	}

	go (func() {
		defer notify.Stop(ch);
		for {
			select {
			case ev := <-ch: {
				fname := util.GetFileName(ev.Path());
				if fname != fileName {
					break;
				}
				lock.Lock();
				flagNeedUpdate = true;
				lock.Unlock();
			}
			}
		}
	})()

	ticker := time.NewTicker(time.Second)
	go (func() {
		defer ticker.Stop();
		for {
			select {
			case <-ticker.C: {
				lock.Lock();
				if !flagNeedUpdate {
					lock.Unlock();
					break;
				}
				flagNeedUpdate = false;
				lock.Unlock();

				cb();
			}
			}
		}
	})()
}

func testConfigListenServer() {
	fmt.Print("");
}
