
package server

import (
	"fmt"
	"math/rand"
	"net"
	// "regexp"
	"strconv"
	"time"

	"github.com/miekg/dns"

	. "model"
)

type MsgServer struct {
	
}

var insMsgServer *MsgServer = nil;

func GetMsgServer() *MsgServer {
	if(insMsgServer == nil) {
		insMsgServer = &MsgServer{};
	}
	return insMsgServer
}

func (c *MsgServer) createMsg(count int) *dns.Msg {
	// rr := new(dns.A);
	// rr.Hdr = dns.RR_Header{ Name: url, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 1 };
	// rr.A = net.ParseIP("127.0.0.1");

	pmsg := &dns.Msg {
		Question: make([]dns.Question, 1),
		Answer: make([]dns.RR, count),
	};
	pmsg.RecursionDesired = true;
	pmsg.RecursionAvailable = true;
	pmsg.Response = true;
	pmsg.Question[0] = dns.Question{ Name: "#.", Qtype: dns.TypeA, Qclass: dns.ClassINET };
	// pmsg.Answer[0] = rr;
	return pmsg;
}

func (c *MsgServer) createRR(url string) *dns.A {
	rr := new(dns.A);
	rr.Hdr = dns.RR_Header{ Name: url + ".", Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 1 };
	rr.A = net.ParseIP("127.0.0.1");
	return rr;
}

// message: Online
//    Question[0] #0.randNum.hashPublicKey
//    Question[1] domain
//    Question[2] rsa(hash(domain+randNum))
func (c *MsgServer) SendOnline() {
	// pmsg := &dns.Msg {
	// 	Question: make([]dns.Question, 1),
	// };
	strRandNum := strconv.FormatInt(int64(rand.Intn(0xfffffff)), 16);
	pmsg := c.createMsg(3);

	str := "#0." + strRandNum;
	// pmsg.Question[0] = dns.Question{ Name: str + ".", Qtype: dns.TypeA, Qclass: dns.ClassINET };
	pmsg.Answer[0] = c.createRR(str);

	domain := "a";
	// pmsg.Question[1] = dns.Question{ Name: domain + ".", Qtype: dns.TypeA, Qclass: dns.ClassINET };
	pmsg.Answer[1] = c.createRR(domain);

	rsaDomain := "b";
	// pmsg.Question[2] = dns.Question{ Name: rsaDomain + ".", Qtype: dns.TypeA, Qclass: dns.ClassINET };
	pmsg.Answer[2] = c.createRR(rsaDomain);

	// client := &dns.Client{ };
	
	// addr := net.JoinHostPort(GetComModel().BroadcastIp, GetComModel().BroadcastPort);
	
	go c.syncSendOnline(pmsg);
	// rst, _, err := client.Exchange(pmsg, addr);
	// fmt.Println(err);
	// fmt.Println(rst);
}

func (c *MsgServer) getSendAddr() string {
	md := GetComModel().ConfigMd;
	return net.JoinHostPort(md.ComClient.ManageIp, md.ComClient.ManagePort);
}

func (c *MsgServer) syncSendOnline(pmsg *dns.Msg) {
	// client := &dns.Client{ };
	client := &dns.Client{ Timeout: 1 * time.Millisecond };

	rst, _, err := client.Exchange(pmsg, c.getSendAddr());
	fmt.Println("111", err);
	fmt.Println("222", rst);
}

func (c *MsgServer) DealMsg(r *dns.Msg) {
	if(!r.Response) {
		return;
	}
	if(len(r.Answer) == 0) {
		return;
	}
}
