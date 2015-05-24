package dcc

import (
	"fmt"
	"io"
	"net"
	"testing"

	"github.com/fiorix/go-diameter/diam"
	"github.com/fiorix/go-diameter/diam/avp"
	"github.com/fiorix/go-diameter/diam/datatype"
	"github.com/fiorix/go-diameter/diam/diamtest"
)

type Server struct {
	*diamtest.Server

	errorCh chan error
}

func (s *Server) ErrorNotify() <-chan error {
	return s.errorCh
}

func NewTestServer() *Server {
	testServer := &Server{
		errorCh: make(chan error),
	}

	smux := diam.NewServeMux()
	smux.Handle("CER", testServer.HandleCER())
	smux.Handle("DWR", testServer.HandleDWR())

	testServer.Server = diamtest.NewServer(smux, nil)

	return testServer
}

func TestClientRequestCER(t *testing.T) {
	server := NewTestServer()
	defer server.Close()

	client := NewClient(server.Address)
	if err := client.Start(); err != nil {
		t.Error(err)
	}

	client.SendCER()

	select {
	case err := <-server.ErrorNotify():
		t.Error(err)
	case err := <-client.ErrorNotify():
		t.Error(err)
	case m := <-client.CERDoneNotify():
		fmt.Println(m)
	}
}

func (s *Server) HandleCER() diam.HandlerFunc {
	return func(conn diam.Conn, m *diam.Message) {
		answerMessage := m.Answer(diam.Success)
		s.SendCEA(conn, answerMessage)
	}
}

func (s *Server) SendCEA(w io.Writer, m *diam.Message) {
	m.NewAVP(avp.OriginHost, avp.Mbit, 0, datatype.DiameterIdentity("srv"))
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, datatype.DiameterIdentity("localhost"))
	m.NewAVP(avp.HostIPAddress, avp.Mbit, 0, datatype.Address(net.ParseIP("127.0.0.1")))
	m.NewAVP(avp.VendorID, avp.Mbit, 0, datatype.Unsigned32(99))
	m.NewAVP(avp.ProductName, avp.Mbit, 0, datatype.UTF8String("go-diameter"))
	_, err := m.WriteTo(w)
	if err != nil {
		s.errorCh <- err
	}
}

func TestClientRequestDWR(t *testing.T) {
	server := NewTestServer()
	defer server.Close()

	client := NewClient(server.Address)
	if err := client.Start(); err != nil {
		t.Error(err)
	}

	client.SendDWR()

	select {
	case err := <-server.ErrorNotify():
		t.Error(err)
	case err := <-client.ErrorNotify():
		t.Error(err)
	case m := <-client.DWRDoneNotify():
		fmt.Println(m)
	}
}

func (s *Server) HandleDWR() diam.HandlerFunc {
	return func(conn diam.Conn, m *diam.Message) {
		answerMessage := m.Answer(diam.Success)
		s.SendDWA(conn, answerMessage)
	}
}

func (s *Server) SendDWA(w io.Writer, m *diam.Message) {
	m.NewAVP(avp.OriginHost, avp.Mbit, 0, datatype.OctetString("srv"))
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, datatype.OctetString("localhost"))

	_, err := m.WriteTo(w)
	if err != nil {
		s.errorCh <- err
	}
}
