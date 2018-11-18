package quic

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"math/big"
	"strings"

	"github.com/lucas-clemente/quic-go"
)

type Server struct {
	preHandlers []func(ResponseWriter, *Request)
	r, mr       map[string]func(ResponseWriter, *Request)
	addr        string
}

func NewServer(addr string) *Server {
	s := &Server{}
	s.r = make(map[string]func(ResponseWriter, *Request))
	s.mr = make(map[string]func(ResponseWriter, *Request))
	return s
}
func (s *Server) ListenAndServe() error {
	l, e := quic.ListenAddr(s.addr, GenConfig(), nil)
	if e != nil {
		return e
	}
	for {
		sess, e := l.Accept()
		if e != nil {
			return e
		}
		go s.ServeSess(sess)
	}
}
func (s *Server) ServeSess(sess quic.Session) {
	defer sess.Close()
	stream, e := sess.AcceptStream()
	if e != nil {
		return
	}
	reader := bufio.NewReader(stream)
	line, _, e := reader.ReadLine()
	if e != nil {
		return
	}
	r := &Request{}
	w := ResponseWriter{c: stream}
	e = json.Unmarshal(line, r)
	if e != nil {
		stream.Write([]byte("protocal invalid"))
		sess.Close()
		return
	}
	for _, v := range s.preHandlers {
		v(w, r)
	}
	url := strings.Split(r.Url, "?")[0]
	if h, ok := s.r[url]; ok {
		h(w, r)
	} else if k, ok := hasPreffixInMap(s.mr, url); ok {
		s.mr[k](w, r)
	} else {
		w.ReturnErr("404 not found")
	}
}
func (mainServer *Server) HandleFunc(url string, f func(ResponseWriter, *Request)) {
	mainServer.r[url] = f
}
func (s *Server) HandleMultiReqs(url string, f func(ResponseWriter, *Request)) {
	s.mr[url] = f
}
func GenConfig() *tls.Config {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}
	return &tls.Config{Certificates: []tls.Certificate{tlsCert}}
}
func hasPreffixInMap(m map[string]func(ResponseWriter, *Request), p string) (string, bool) {
	for k, _ := range m {
		if len(p) >= len(k) && k == p[:len(k)] {
			return k, true
		}
	}
	return "", false
}
