package main

import (
	"fmt"
	// "time"
	// "regexp"
	"net"

	"github.com/miekg/dns"

	// util "util"
)

func main() {
	comTest := ComTest{};
	comTest.Run();
}

type ComTest struct{
	// db *sqlx.DB;
}

func (c *ComTest) Run(){
	fmt.Println("test run");

	// var dnsRR []dns.RR;
	rr := new(dns.A);
	rr.Hdr = dns.RR_Header{ Name: "#0.d811f8d.", Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 1 };
	rr.A = net.ParseIP("127.0.0.1");
	// dnsRR = append(dnsRR, rr);
	
	rr2 := new(dns.A);
	rr2.Hdr = dns.RR_Header{ Name: "#0.d811f8d.", Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 1 };
	rr2.A = net.ParseIP("127.0.0.1");
	
	rr3 := new(dns.A);
	rr3.Hdr = dns.RR_Header{ Name: "#0.d811f8d.", Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 1 };
	rr3.A = net.ParseIP("127.0.0.1");

	// arrQus := []dns.Question{};
	// qus := dns.Question { Name: "#0", Qtype: dns.TypeA, Qclass: dns.ClassINET };
	// arrQus = append(arrQus, qus);

	// r, _ := dns.NewRR("vcedit.lan 3600 IN A 127.0.0.1 ");

	pmsg := &dns.Msg {
		Question: make([]dns.Question, 1),
		Answer: make([]dns.RR, 3),
	};
	pmsg.Question[0] = dns.Question{ Name: "#.", Qtype: dns.TypeA, Qclass: dns.ClassINET };
	pmsg.Answer[0] = rr;
	pmsg.Answer[1] = rr2;
	pmsg.Answer[2] = rr3;
	// pmsg.Id = 0;
	pmsg.RecursionDesired = true;
	pmsg.RecursionAvailable = true;
	pmsg.Response = true;
	// pmsg.Opcode = 0;
	// pmsg.RecursionDesired = false;
	// pmsg.CheckingDisabled = false;
	// pmsg.Rcode = dns.RcodeSuccess;
	// pmsg.Question = arrQus;
	
	client := new(dns.Client);
	addr := net.JoinHostPort("255.255.255.255", "53");
	rst, _, err := client.Exchange(pmsg, addr);
	fmt.Println(err);
	fmt.Println(rst);
}
