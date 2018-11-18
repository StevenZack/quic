package quic

import (
	"encoding/json"

	quic "github.com/lucas-clemente/quic-go"
)

type ResponseWriter struct {
	c quic.Stream
}

type Request struct {
	Url  string
	Body interface{}
}
type Response struct {
	Status, Info string
	Body         interface{}
}

// ReturnErr status:"ERR", info:info
func (r *ResponseWriter) ReturnErr(info string) error {
	rp := Response{}
	rp.Status = "ERR"
	rp.Info = info
	b, e := json.Marshal(rp)
	if e != nil {
		return e
	}
	_, e = r.c.Write(b)
	if e != nil {
		return e
	}
	_, e = r.c.Write([]byte("\n"))
	if e != nil {
		return e
	}
	return nil
}

// Close stream.Close
func (r *ResponseWriter) Close() error {
	return r.c.Close()
}

// ReturnInfo status:"OK", info:info
func (r *ResponseWriter) ReturnInfo(info string) error {
	rp := Response{}
	rp.Status = "OK"
	rp.Info = info
	b, e := json.Marshal(rp)
	if e != nil {
		return e
	}
	_, e = r.c.Write(b)
	if e != nil {
		return e
	}
	_, e = r.c.Write([]byte("\n"))
	if e != nil {
		return e
	}
	return nil
}

// ReturnData status:"OK", body:data
func (r *ResponseWriter) ReturnData(data interface{}) error {
	rp := Response{}
	rp.Status = "OK"
	rp.Body = data
	b, e := json.Marshal(rp)
	if e != nil {
		return e
	}
	_, e = r.c.Write(b)
	if e != nil {
		return e
	}
	_, e = r.c.Write([]byte("\n"))
	if e != nil {
		return e
	}
	return nil
}
