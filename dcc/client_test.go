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
	smux.Handle("CCR", testServer.HandleCCR())

	testServer.Server = diamtest.NewServer(smux, nil)

	return testServer
}

func NewTestClient(address string) *diameterClient {
	return NewClient(DiameterConfig{
		URL:              address,
		OriginHost:       datatype.DiameterIdentity("client"),
		OriginRealm:      datatype.DiameterIdentity("localhost"),
		DestinationHost:  datatype.DiameterIdentity("srv"),
		DestinationRealm: datatype.DiameterIdentity("localhost"),
		VendorID:         datatype.Unsigned32(0),
		ProductName:      datatype.UTF8String("go-diameter"),
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
	defer client.Close()

	client.sendCER()

	select {
	case err := <-server.ErrorNotify():
		t.Error(err)
	case err := <-client.ErrorNotify():
		t.Error(err)
	case <-client.cerDoneNotify():
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
	defer client.Close()

	client.sendDWR()

	select {
	case err := <-server.ErrorNotify():
		t.Error(err)
	case err := <-client.ErrorNotify():
		t.Error(err)
	case <-client.dwrDoneNotify():
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
	defer client.Close()

	client.Init()

	select {
	case err := <-server.ErrorNotify():
		t.Error(err)
	case err := <-client.ErrorNotify():
		t.Error(err)
	case <-client.dwrDoneNotify():
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
	defer client.Close()

	go client.loopWatchdog()

	interval := 2
	for i := 0; i < interval; i++ {
		select {
		case err := <-server.ErrorNotify():
			t.Error(err)
		case err := <-client.ErrorNotify():
			t.Error(err)
		case <-client.watchdogAliveNotify():
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
	defer client.Close()

	client.sendCER()

	select {
	case err := <-server.ErrorNotify():
		t.Fatal(err)
	case err := <-client.ErrorNotify():
		t.Fatal(err)
	case <-client.cerDoneNotify():
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
	m := diam.NewRequest(diam.DeviceWatchdog, 0, nil)

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

func TestClientCallCCR(t *testing.T) {
	server := NewTestServer()
	defer server.Close()

	client := NewTestClient(server.Address)
	if err := client.Start(); err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	client.sendCCR([]*diam.AVP{})

	select {
	case err := <-server.ErrorNotify():
		t.Error(err)
	case err := <-client.ErrorNotify():
		t.Error(err)
	case <-client.ccrDoneNotify():
	case <-time.After(time.Second):
		t.Error("server timeout")
	}
}

func (s *Server) HandleCCR() diam.HandlerFunc {
	return func(conn diam.Conn, m *diam.Message) {
		s.conn = conn
		answerMessage := m.Answer(diam.Success)
		s.SendCCA(answerMessage)
	}
}

func (s *Server) SendCCA(m *diam.Message) {
	m.NewAVP(avp.OriginHost, avp.Mbit, 0, datatype.DiameterIdentity("srv"))
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, datatype.DiameterIdentity("localhost"))

	_, err := m.WriteTo(s.conn)
	if err != nil {
		s.errorCh <- err
	}
}

func TestClientRunBackgroundCCR(t *testing.T) {
	server := NewTestServer()
	defer server.Close()

	client := NewTestClient(server.Address)
	if err := client.Start(); err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	client.Init()

	request := &mockRequest{
		outCh: make(chan *diam.Message),
	}
	anotherRequest := &mockRequest{
		outCh: make(chan *diam.Message),
	}
	client.Serve(request)
	client.Serve(anotherRequest)

	interval := 2
	for i := 0; i < interval; i++ {
		select {
		case err := <-server.ErrorNotify():
			t.Error(err)
		case err := <-client.ErrorNotify():
			t.Error(err)
		case <-request.ResponseNotify():
		case <-anotherRequest.ResponseNotify():
		case <-time.After(time.Second):
			t.Error("server timeout")
		}
	}
}

type mockRequest struct {
	outCh chan *diam.Message
}

func (r *mockRequest) Response(m *diam.Message) {
	r.outCh <- m
}

func (r *mockRequest) ResponseNotify() <-chan *diam.Message {
	return r.outCh
}

func (r *mockRequest) AVP() []*diam.AVP {
	return []*diam.AVP{}
}
