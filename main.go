package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/miekg/dns"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)

	ns := os.Args[1]
	host := os.Args[2]

	m := new(dns.Msg)
	m.RecursionDesired = true
	m.SetQuestion(dns.Fqdn(host), dns.TypeAAAA)

	c := new(dns.Client)
	c.Net = "tcp-tls"
	in, rtt, err := c.Exchange(m, ns+":853")
	if err != nil {
		log.Println(err)
	}

	fmt.Println(in)
	fmt.Println(rtt)

	doh(m, ns)
}

func doh(m *dns.Msg, ns string) {

	client := &http.Client{}

	packedmsg, err := m.Pack()
	if err != nil {
		log.Println(err)
	}

	q := base64.RawURLEncoding.EncodeToString(packedmsg)

	u, err := url.Parse("https://" + ns + "/dns-query?dns=" + q)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(u)

	req := &http.Request{
		Method:     "GET",
		URL:        u,
		Proto:      "HTTP/2",
		ProtoMajor: 2,
		ProtoMinor: 0,
		Header: map[string][]string{
			"accept": {"application/dns-message"},
		},
		Host: u.Hostname(),
		Body: nil,
	}

	t := time.Now()
	resp1, err := client.Do(req)
	t2 := time.Since(t)
	fmt.Println(t2)
	if err != nil {
		log.Println(err)
	}
	if resp1 == nil {
		return
	}

	defer resp1.Body.Close()

	rb, err := ioutil.ReadAll(resp1.Body)
	if err != nil {
		log.Println(err)
	}

	msg2 := new(dns.Msg)
	msg2.Unpack(rb)
	fmt.Println(msg2)
	fmt.Println(string(rb))

	t3 := time.Now()
	resp, err := client.Get(u.String())
	t4 := time.Since(t3)
	fmt.Println(t4)
	if err != nil {
		log.Println(err)
	}
	if resp == nil {
		return
	}

	defer resp.Body.Close()

	r, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	msg3 := new(dns.Msg)
	msg3.Unpack(r)
	fmt.Println(msg2)
	fmt.Println(string(r))
}
