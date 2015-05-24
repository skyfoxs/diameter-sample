package dcc

import (
	"net"
	"testing"
	"time"

	"github.com/fiorix/go-diameter/diam"
	"github.com/fiorix/go-diameter/diam/avp"
	"github.com/fiorix/go-diameter/diam/datatype"
	"github.com/fiorix/go-diameter/diam/diamtest"
)

type Server struct {
	*diamtest.Server
	conn diam.Conn

	errorCh chan error
	dwaCh   chan *diam.Message
}

func (s *Server) ErrorNotify() <-chan error {
	return s.errorCh
}

func (s *Server) ResponseWatchdogNotify() <-chan *diam.Message {
	return s.dwaCh
}

func NewTestServer() *Server {
	testServer := &Server{
		errorCh: make(chan error),
		dwaCh:   make(chan *diam.Message),
	}

	smux := diam.NewServeMux()
	smux.Handle("CER", testServer.HandleCER())
	smux.Handle("DWR", testServer.HandleDWR())
	smux.Handle("DWA", testServer.HandleDWA())

	testServer.Server = diamtest.NewServer(smux, nil)

	return testServer
}

func NewTestClient(address string) *DiameterClient {
	return NewClient(DiameterConfig{
		URL:              address,
		OriginHost:       datatype.DiameterIdentity("jenkin13_OMR_TEST01"),
		OriginRealm:      datatype.DiameterIdentity("dtac.co.th"),
		VendorID:         datatype.Unsigned32(0),
		ProductName:      datatype.UTF8String("omr"),
		FirmwareRevision: datatype.Unsigned32(1),
		WatchdogInterval: 100 * time.Millisecond,
	})
}

func TestClientRequestCER(t *testing.T) {
	server := NewTestServer()
	defer server.Close()

	client := NewTestClient(server.Address)
	if err := client.Start(); err != nil {
		t.Error(err)
	}

	client.SendCER()

	select {
	case err := <-server.ErrorNotify():
		t.Error(err)
	case err := <-client.ErrorNotify():
		t.Error(err)
	case <-client.CERDoneNotify():
	case <-time.After(time.Second):
		t.Error("timeout")
	}
}

func (s *Server) HandleCER() diam.HandlerFunc {
	return func(conn diam.Conn, m *diam.Message) {
		s.conn = conn
		answerMessage := m.Answer(diam.Success)
		s.SendCEA(answerMessage)
	}
}

func (s *Server) SendCEA(m *diam.Message) {
	m.NewAVP(avp.OriginHost, avp.Mbit, 0, datatype.DiameterIdentity("srv"))
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, datatype.DiameterIdentity("localhost"))
	m.NewAVP(avp.HostIPAddress, avp.Mbit, 0, datatype.Address(net.ParseIP("127.0.0.1")))
	m.NewAVP(avp.VendorID, avp.Mbit, 0, datatype.Unsigned32(99))
	m.NewAVP(avp.ProductName, avp.Mbit, 0, datatype.UTF8String("go-diameter"))
	_, err := m.WriteTo(s.conn)
	if err != nil {
		s.errorCh <- err
	}
}

func TestClientRequestDWR(t *testing.T) {
	server := NewTestServer()
	defer server.Close()

	client := NewTestClient(server.Address)
	if err := client.Start(); err != nil {
		t.Error(err)
	}

	client.SendDWR()

	select {
	case err := <-server.ErrorNotify():
		t.Error(err)
	case err := <-client.ErrorNotify():
		t.Error(err)
	case <-client.DWRDoneNotify():
	case <-time.After(time.Second):
		t.Error("timeout")
	}
}

func (s *Server) HandleDWR() diam.HandlerFunc {
	return func(conn diam.Conn, m *diam.Message) {
		s.conn = conn
		answerMessage := m.Answer(diam.Success)
		s.SendDWA(answerMessage)
	}
}

func (s *Server) SendDWA(m *diam.Message) {
	m.NewAVP(avp.OriginHost, avp.Mbit, 0, datatype.OctetString("srv"))
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, datatype.OctetString("localhost"))
	_, err := m.WriteTo(s.conn)
	if err != nil {
		s.errorCh <- err
	}
}

func TestClientHandshakeAndRequestWatchdog(t *testing.T) {
	server := NewTestServer()
	defer server.Close()

	client := NewTestClient(server.Address)
	if err := client.Start(); err != nil {
		t.Error(err)
	}

	client.Init()

	select {
	case err := <-server.ErrorNotify():
		t.Error(err)
	case err := <-client.ErrorNotify():
		t.Error(err)
	case <-client.DWRDoneNotify():
	case <-time.After(time.Second):
		t.Error("timeout")
	}
}

func TestBackgroundWatchdog(t *testing.T) {
	server := NewTestServer()
	defer server.Close()

	client := NewTestClient(server.Address)
	if err := client.Start(); err != nil {
		t.Error(err)
	}

	go client.LoopWatchdog()

	interval := 2
	for i := 0; i < interval; i++ {
		select {
		case err := <-server.ErrorNotify():
			t.Error(err)
		case err := <-client.ErrorNotify():
			t.Error(err)
		case <-client.WatchdogAliveNotify():
		case <-time.After(time.Second):
			t.Error("timeout")
		}
	}
}

func TestServerCallDWR(t *testing.T) {
	server := NewTestServer()
	defer server.Close()

	client := NewTestClient(server.Address)
	if err := client.Start(); err != nil {
		t.Fatal(err)
	}

	client.SendCER()

	select {
	case err := <-server.ErrorNotify():
		t.Fatal(err)
	case err := <-client.ErrorNotify():
		t.Fatal(err)
	case <-client.CERDoneNotify():
	case <-time.After(time.Second):
		t.Fatal("client timeout")
	}

	server.SendDWR()

	select {
	case err := <-server.ErrorNotify():
		t.Error(err)
	case err := <-client.ErrorNotify():
		t.Error(err)
	case <-server.ResponseWatchdogNotify():
	case <-time.After(time.Second):
		t.Error("server timeout")
	}
}

func (s *Server) SendDWR() {
	m := diam.NewRequest(watchdogExchange, 0, nil)

	m.NewAVP(avp.OriginHost, avp.Mbit, 0, datatype.DiameterIdentity("srv"))
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, datatype.DiameterIdentity("localhost"))

	_, err := m.WriteTo(s.conn)

	if err != nil {
		s.errorCh <- err
	}
}

func (s *Server) HandleDWA() diam.HandlerFunc {
	return func(conn diam.Conn, m *diam.Message) {
		s.dwaCh <- m
	}
}
