
package server

import (
	"fmt"
	"net"
	"regexp"

	"github.com/miekg/dns"
)

type DnsServer struct {
	ttl uint32;
}

var insDnsServer *DnsServer = nil;

func msgAcceptFunc(dh dns.Header) dns.MsgAcceptAction {
	_QR := (uint16(1) << 15);
	if isResponse := dh.Bits&_QR != 0; isResponse {
		// return dns.MsgIgnore
		return dns.MsgAccept
	}

	// Don't allow dynamic updates, because then the sections can contain a whole bunch of RRs.
	opcode := int(dh.Bits>>11) & 0xF
	if opcode != dns.OpcodeQuery && opcode != dns.OpcodeNotify {
		return dns.MsgRejectNotImplemented
	}

	if dh.Qdcount != 1 {
		return dns.MsgReject
	}
	// NOTIFY requests can have a SOA in the ANSWER section. See RFC 1996 Section 3.7 and 3.11.
	if dh.Ancount > 1 {
		return dns.MsgReject
	}
	// IXFR request could have one SOA RR in the NS section. See RFC 1995, section 3.
	if dh.Nscount > 1 {
		return dns.MsgReject
	}
	if dh.Arcount > 2 {
		return dns.MsgReject
	}
	return dns.MsgAccept
}

func GetDnsServer() *DnsServer {
	if(insDnsServer == nil) {
		insDnsServer = &DnsServer{};
		insDnsServer.ttl = 60;
	}
	return insDnsServer
}

// func (c *DnsServer) Init() {
// 	dns.DefaultMsgAcceptFunc = msgAcceptFunc;
// }

func (c *DnsServer) findAddr(domain string) string {
	reg, _ := regexp.Compile(".*vcedit\\.lan");
	if reg.Match([]byte(domain)) {
		return "127.0.0.1";
	} else {
		return "";
	}
}

// response dns type A
func (c *DnsServer) WriteTypeA(w dns.ResponseWriter, r *dns.Msg, addr string) {
	name := r.Question[0].Name;

	m := new(dns.Msg);
	m.SetReply(r);
	m.Authoritative = true;

	var dnsRR []dns.RR;
	// for _, address := range addresses {
	rr := new(dns.A);
	fmt.Println(name);
	fmt.Println(c);
	rr.Hdr = dns.RR_Header{ Name: name, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: c.ttl };
	rr.A = net.ParseIP(addr);
	dnsRR = append(dnsRR, rr);
	// }
	m.Answer = dnsRR;
	w.WriteMsg(m);
}

// unsupport dns type
func (c *DnsServer) HandleUnsupportType(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg);
	m.Authoritative = true;
	m.SetRcode(r, dns.RcodeNameError);
	w.WriteMsg(m);
}

// func (s *DnsServer) HandleMatched(w dns.ResponseWriter, qtype uint16) {

// }

func (c *DnsServer) queryDns(pmsg *dns.Msg) *dns.Msg {
	// m := &dns.Msg{
	// 	Question: make([]dns.Question, 1),
	// }
	// m.Question[0] = dns.Question{domain, qtype, dns.ClassCHAOS}
	client := new(dns.Client);
	addr := net.JoinHostPort("114.114.114.114", "53");
	rst, _, _ := client.Exchange(pmsg, addr);
	return rst;
}

// unmatched dns
func (c *DnsServer) HandleUnmatched(w dns.ResponseWriter, r *dns.Msg) {
	// m := new(dns.Msg);
	// m.SetRcode(r, dns.RcodeServerFailure);
	// w.WriteMsg(m);
	rst := c.queryDns(r);
	if(rst != nil) {
		w.WriteMsg(rst);
	} else {
		c.HandleUnsupportType(w, r);
	}
}

// listen proc
func (c *DnsServer) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	if(r == nil) {
		return;
	}

	// message response
	if(r.Response) {
		if(len(r.Answer) == 0) {
			return;
		}

		GetMsgServer().DealMsg(r);

		// fmt.Println("bbb:", r.Answer[0].String());
		return;
	}

	// question
	if(len(r.Question) == 0) {
		c.HandleUnsupportType(w, r);
		return;
	}
	name := r.Question[0].Name;
	qtype := r.Question[0].Qtype;
	fmt.Println("aaa:", name, qtype);

	// 对于A记录解析请求
	// if qtype == dns.TypeA {

	domain := name
	if(len(domain) <= 1) {
		c.HandleUnsupportType(w, r);
		return;
	}

	// 去掉结尾的.
	if(domain[len(domain)-1] == '.') {
		domain = domain[0 : len(domain)-1]
	}

	// self message
	if(domain[0] == '#') {
		w.Close();
		return;
	}

	// find address
	addr := c.findAddr(domain);

	// not matched
	if(addr == "") {
		c.HandleUnmatched(w, r);
		return;
	}

	// matched
	switch(qtype) {
	case dns.TypeA: {
		// A
		c.WriteTypeA(w, r, addr);
		break;
	}
	default: {
		// Unmatched
		c.HandleUnmatched(w, r);
		break;
	}
	}



	// // 返回固定DNS解析
	// reg, _ := regexp.Compile(".*vcedit\\.lan")
	// if reg.Match([]byte(domain)) {
	// 	fmt.Println("bbb:", domain)
	// 	addresses = []string{"127.0.0.1"}
	// } else {
	// 	s.HandleUnmatched(w, r)
	// 	return
	// }

	// // 构建返回信息
	// m := new(dns.Msg)
	// m.SetReply(r)
	// m.Authoritative = true

	// var dnsRR []dns.RR
	// for _, address := range addresses {
	// 	rr := new(dns.A)
	// 	rr.Hdr = dns.RR_Header{Name: name, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: uint32(ttl)}
	// 	rr.A = net.ParseIP(address)
	// 	dnsRR = append(dnsRR, rr)
	// }
	// m.Answer = dnsRR
	// w.WriteMsg(m)



	// } else {
	// 	fmt.Println("query type not support", qtype)
	// 	s.HandleUnmatched(w, r)
	// }
}

func (c *DnsServer) ListenAsync(port string) error {
	err := c.logicListenAsync(port, true);
	if(err != nil) {
		fmt.Println("ListenAsync err: ", err);
		return err;
	}

	err = c.logicListenAsync(port, false);
	if(err != nil) {
		fmt.Println("ListenAsync err: ", err);
		return err;
	}

	return nil;
}

func (c *DnsServer) logicListenAsync(port string, isIpv6 bool) error {
	str := "";
	if(isIpv6) {
		str = "6";
	}

	// check port
	udpAddr, err := net.ResolveUDPAddr("udp" + str, port)
	if err != nil {
		// fmt.Printf("run error: %v\n", err)
		return err;
	}

	// 监听UDP端口
	p, err := net.ListenUDP("udp" + str, udpAddr)
	if err != nil {
		// fmt.Printf("run error: %v\n", err)
		return err;
	}

	// tcp端口
	tcpAddr, err := net.ResolveTCPAddr("tcp" + str, port)
	if err != nil {
		// fmt.Printf("run error: %v\n", err)
		return err;
	}
	l, err := net.ListenTCP("tcp", tcpAddr)

	// 协程启动DNS Server
	// go dns.ActivateAndServe(l, p, c)
	go c.runListenAsync(l, p);

	return nil;
}

func (c *DnsServer) runListenAsync(l *net.TCPListener, p *net.UDPConn) {
	// err := dns.ActivateAndServe(l, p, c);

	server := &dns.Server{Listener: l, PacketConn: p, Handler: c, MsgAcceptFunc: msgAcceptFunc};
	err := server.ActivateAndServe()

	if(err != nil) {
		fmt.Println("runListenAsync err: ", err);
	}
}
