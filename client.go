package quic

import (
	"bufio"
	"crypto/tls"
	"encoding/json"

	"github.com/lucas-clemente/quic-go"
)

// Post post data to server
func Post(addr, url string, data interface{}) (*Response, error) {
	conf := &tls.Config{}
	conf.InsecureSkipVerify = true
	s, e := quic.DialAddr(addr, conf, nil)
	if e != nil {
		return nil, e
	}
	c, e := s.OpenStreamSync()
	if e != nil {
		return nil, e
	}
	req := &Request{}
	req.Url = url
	req.Body = data
	b, e := json.Marshal(req)
	if e != nil {
		return nil, e
	}
	_, e = c.Write(b)
	if e != nil {
		return nil, e
	}
	_, e = c.Write([]byte("\n"))
	if e != nil {
		return nil, e
	}
	line, _, e := bufio.NewReader(c).ReadLine()
	if e != nil {
		return nil, e
	}
	rp := Response{}
	e = json.Unmarshal(line, &rp)
	if e != nil {
		return nil, e
	}
	return &rp, nil
}

// Get Post without body
func Get(addr, url string) (*Response, error) {
	conf := &tls.Config{}
	conf.InsecureSkipVerify = true
	s, e := quic.DialAddr(addr, conf, nil)
	if e != nil {
		return nil, e
	}
	c, e := s.OpenStreamSync()
	if e != nil {
		return nil, e
	}
	req := &Request{}
	req.Url = url
	b, e := json.Marshal(req)
	if e != nil {
		return nil, e
	}
	_, e = c.Write(b)
	if e != nil {
		return nil, e
	}
	_, e = c.Write([]byte("\n"))
	if e != nil {
		return nil, e
	}
	line, _, e := bufio.NewReader(c).ReadLine()
	if e != nil {
		return nil, e
	}
	rp := Response{}
	e = json.Unmarshal(line, &rp)
	if e != nil {
		return nil, e
	}
	return &rp, nil
}
