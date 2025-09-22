package dns

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/miekg/dns"
)

type Resolver struct {
	timeout time.Duration
	client  *dns.Client
}

func NewResolver(timeoutSeconds int) *Resolver {
	timeout := time.Duration(timeoutSeconds) * time.Second
	client := &dns.Client{
		Timeout: timeout,
	}

	return &Resolver{
		timeout: timeout,
		client:  client,
	}
}

func (r *Resolver) Resolve(domain string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	done := make(chan string, 1)
	errChan := make(chan error, 1)

	go func() {
		ip, err := r.resolveA(domain)
		if err != nil {
			errChan <- err
			return
		}
		done <- ip
	}()

	select {
	case ip := <-done:
		return ip, nil
	case err := <-errChan:
		return "", err
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

func (r *Resolver) resolveA(domain string) (string, error) {
	msg := new(dns.Msg)
	msg.Id = dns.Id()
	msg.RecursionDesired = true
	msg.Question = make([]dns.Question, 1)
	msg.Question[0] = dns.Question{
		Name:   dns.Fqdn(domain),
		Qtype:  dns.TypeA,
		Qclass: dns.ClassINET,
	}

	servers := []string{"8.8.8.8:53", "1.1.1.1:53", "8.8.4.4:53"}

	for _, server := range servers {
		response, _, err := r.client.Exchange(msg, server)
		if err != nil {
			continue
		}

		if response.Rcode != dns.RcodeSuccess {
			continue
		}

		for _, answer := range response.Answer {
			if aRecord, ok := answer.(*dns.A); ok {
				return aRecord.A.String(), nil
			}
		}
	}

	return "", fmt.Errorf("no A record found for %s", domain)
}

func (r *Resolver) ResolveCNAME(domain string) (string, error) {
	msg := new(dns.Msg)
	msg.Id = dns.Id()
	msg.RecursionDesired = true
	msg.Question = make([]dns.Question, 1)
	msg.Question[0] = dns.Question{
		Name:   dns.Fqdn(domain),
		Qtype:  dns.TypeCNAME,
		Qclass: dns.ClassINET,
	}

	servers := []string{"8.8.8.8:53", "1.1.1.1:53", "8.8.4.4:53"}

	for _, server := range servers {
		response, _, err := r.client.Exchange(msg, server)
		if err != nil {
			continue
		}

		if response.Rcode != dns.RcodeSuccess {
			continue
		}

		for _, answer := range response.Answer {
			if cnameRecord, ok := answer.(*dns.CNAME); ok {
				return cnameRecord.Target, nil
			}
		}
	}

	return "", fmt.Errorf("no CNAME record found for %s", domain)
}

func (r *Resolver) ResolveMX(domain string) ([]string, error) {
	msg := new(dns.Msg)
	msg.Id = dns.Id()
	msg.RecursionDesired = true
	msg.Question = make([]dns.Question, 1)
	msg.Question[0] = dns.Question{
		Name:   dns.Fqdn(domain),
		Qtype:  dns.TypeMX,
		Qclass: dns.ClassINET,
	}

	var mxRecords []string
	servers := []string{"8.8.8.8:53", "1.1.1.1:53", "8.8.4.4:53"}

	for _, server := range servers {
		response, _, err := r.client.Exchange(msg, server)
		if err != nil {
			continue
		}

		if response.Rcode != dns.RcodeSuccess {
			continue
		}

		for _, answer := range response.Answer {
			if mxRecord, ok := answer.(*dns.MX); ok {
				mxRecords = append(mxRecords, mxRecord.Mx)
			}
		}
		break
	}

	if len(mxRecords) == 0 {
		return nil, fmt.Errorf("no MX records found for %s", domain)
	}

	return mxRecords, nil
}

func (r *Resolver) ResolveTXT(domain string) ([]string, error) {
	msg := new(dns.Msg)
	msg.Id = dns.Id()
	msg.RecursionDesired = true
	msg.Question = make([]dns.Question, 1)
	msg.Question[0] = dns.Question{
		Name:   dns.Fqdn(domain),
		Qtype:  dns.TypeTXT,
		Qclass: dns.ClassINET,
	}

	var txtRecords []string
	servers := []string{"8.8.8.8:53", "1.1.1.1:53", "8.8.4.4:53"}

	for _, server := range servers {
		response, _, err := r.client.Exchange(msg, server)
		if err != nil {
			continue
		}

		if response.Rcode != dns.RcodeSuccess {
			continue
		}

		for _, answer := range response.Answer {
			if txtRecord, ok := answer.(*dns.TXT); ok {
				for _, txt := range txtRecord.Txt {
					txtRecords = append(txtRecords, txt)
				}
			}
		}
		break
	}

	if len(txtRecords) == 0 {
		return nil, fmt.Errorf("no TXT records found for %s", domain)
	}

	return txtRecords, nil
}

func (r *Resolver) IsValidDomain(domain string) bool {
	_, err := net.LookupHost(domain)
	return err == nil
}
