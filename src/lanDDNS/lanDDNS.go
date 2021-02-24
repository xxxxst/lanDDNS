package lanDDNS

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"os/signal"
    "path/filepath"
	// "regexp"
	"runtime"
	"syscall"
	"time"

	// "github.com/miekg/dns"
	
	. "control"
	tag "tag"
	. "server"
	. "model"
	util "util"
)

type LanDDNS struct {
	// server *Server;
}

var ins *LanDDNS

func GetApp() *LanDDNS {
	if ins == nil {
		ins = new(LanDDNS)
	}
	return ins
}

// func (c *LanDDNS) listenAsync(isIpv6 bool) {
// 	str := "";
// 	if(isIpv6) {
// 		str = "6";
// 	}

// 	// check port
// 	udpAddr, err := net.ResolveUDPAddr("udp" + str, ":53")
// 	if err != nil {
// 		// fmt.Printf("run error: %v\n", err)
// 		return;
// 	}

// 	// 监听UDP端口
// 	p, err := net.ListenUDP("udp" + str, udpAddr)
// 	if err != nil {
// 		// fmt.Printf("run error: %v\n", err)
// 		return;
// 	}

// 	// tcp端口
// 	tcpAddr, err := net.ResolveTCPAddr("tcp" + str, ":53")
// 	if err != nil {
// 		// fmt.Printf("run error: %v\n", err)
// 		return;
// 	}
// 	l, err := net.ListenTCP("tcp", tcpAddr)

// 	// 协程启动DNS Server
// 	go dns.ActivateAndServe(l, p, GetDnsServer())
// }

func (c *LanDDNS) Run() {
	runtime.GOMAXPROCS(runtime.NumCPU());
	rand.Seed(time.Now().Unix())

	// c.server, _ = NewServer()
	
	comMd := GetComModel();
	comMd.Version = "v1.0.0";
	comMd.IsDebug = tag.GetTag().Debug;
	fmt.Println("LanDDNS " + comMd.Version);

	exeDir, _ := filepath.Abs(filepath.Dir(os.Args[0]));
	rootDir := exeDir + "/";
	if(comMd.IsDebug) {
		workDir, _ := os.Getwd();
		rootDir = workDir + "/";
	}
	comMd.RootDir = rootDir;
	comMd.ConfigDir = rootDir + "data/";

	os.Mkdir(comMd.ConfigDir, os.ModePerm);
	
	comMd.ConfigMd = GetConfigCtl().LoadConfig(comMd.ConfigDir + "config.ini");

	if(comMd.ConfigMd.Server.LogMac) {
		util.SaveFileString(rootDir + "mac.log", "");
	} else {
		os.Remove(rootDir + "mac.log");
	}

	GetConfigListenServer().Run();
	GetMacArpListenServer().Run();

	// c.WaitExit();
	// return;

	// GetDnsServer().Init();

	// run server
	dnsIp := comMd.ConfigMd.Server.DnsIp;
	if(dnsIp == "0.0.0.0") {
		dnsIp = "";
	}
	dnsPort := comMd.ConfigMd.Server.DnsPort;
	dnsAddr := net.JoinHostPort(dnsIp, dnsPort);
	enableDnsServer := (dnsPort != "" && dnsPort != "0");
	if(enableDnsServer) {
		GetDnsServer().ListenAsync(dnsAddr);
	}
	
	// run manage server
	mngIp := comMd.ConfigMd.ComClient.ManageIp;
	if(mngIp == "0.0.0.0") {
		mngIp = "";
	}
	mngPort := comMd.ConfigMd.ComClient.ManagePort;
	mngAddr := net.JoinHostPort(mngIp, mngPort);
	enableMngServer := (mngPort != "" && mngPort != "0");
	if(enableMngServer && (!enableDnsServer || mngAddr != dnsAddr)) {
		GetDnsServer().ListenAsync(mngAddr);
	}

	// // 检测UDP端口
	// udpAddr, err := net.ResolveUDPAddr("udp", ":53")
	// if err != nil {
	// 	fmt.Printf("run error: %v\n", err)
	// }
	// // 监听UDP端口
	// p, err := net.ListenUDP("udp", udpAddr)
	// if err != nil {
	// 	fmt.Printf("run error: %v\n", err)
	// }

	// // tcp端口
	// tcpAddr, err := net.ResolveTCPAddr("tcp", ":53")
	// if err != nil {
	// 	fmt.Printf("run error: %v\n", err)
	// }
	// l, err := net.ListenTCP("tcp", tcpAddr)

	// // 协程启动DNS Server
	// go dns.ActivateAndServe(l, p, server)

	// // ipv6
	// // 检测UDP端口
	// udp6Addr, err := net.ResolveUDPAddr("udp6", ":53")
	// if err != nil {
	// 	fmt.Printf("run error: %v\n", err)
	// }
	// // 监听UDP端口
	// p6, err := net.ListenUDP("udp6", udp6Addr)
	// if err != nil {
	// 	fmt.Printf("run error: %v\n", err)
	// }

	// // tcp端口
	// tcp6Addr, err := net.ResolveTCPAddr("tcp6", ":53")
	// if err != nil {
	// 	fmt.Printf("run error: %v\n", err)
	// }
	// l6, err := net.ListenTCP("tcp6", tcp6Addr)

	// // 协程启动DNS Server
	// go dns.ActivateAndServe(l6, p6, server)

	fmt.Println("----------------------------------");
	fmt.Println("server start");
	fmt.Println("----------------------------------");
	
	// if(enableMngServer) {
	// 	time.Sleep(time.Second);
	// 	GetMsgServer().SendOnline();
	// }

	c.WaitExit();
}

func (c *LanDDNS) WaitExit() {
	// wait exit
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	for {
		s := <-sig
		switch s {
		default: {
			fmt.Println("exiting")
			return
		}
		}
	}
}

func log() {

}
