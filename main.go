package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/miekg/dns"
)

func main() {

	dns.HandleFunc(".", func(w dns.ResponseWriter, r *dns.Msg) {

		m := new(dns.Msg)
		log.Printf("Q. %s", r.Question[0].Name)
		if r.Question[0].Qtype != dns.TypeA {
			return
		}
		m.SetReply(r)
		cl := dns.Client{}
		a := dns.Msg{}
		var (
			ok  bool
			ans *dns.A
		)
		n := r.Question[0].Name
		t := dns.TypeA
		for ans == nil {
			a.SetQuestion(n, t)
			res, _, err := cl.Exchange(&a, "8.8.8.8:53")
			if err != nil || len(res.Answer) == 0 {
				return
			}

			ans, ok = res.Answer[0].(*dns.A)
			if !ok {
				cn, cok := res.Answer[0].(*dns.CNAME)
				if !cok {
					return
				}
				n = cn.Target
				t = dns.TypeCNAME
			}
		}
		rr := &dns.A{
			Hdr: dns.RR_Header{Name: r.Question[0].Name, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 0},
			A:   ans.A,
		}

		m.Answer = append(m.Answer, rr)
		w.WriteMsg(m)
	})

	go func() {
		server := &dns.Server{Addr: "[::]:53", Net: "tcp", TsigSecret: nil, ReusePort: false}
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()
	go func() {
		server := &dns.Server{Addr: "[::]:53", Net: "udp", TsigSecret: nil, ReusePort: false}
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}

	}()
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	s := <-sig
	log.Printf("Signal (%s) received, stopping\n", s)
}
