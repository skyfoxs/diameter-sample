package dcc

import (
	"fmt"
	"io"
	"net"
	"time"

	"github.com/fiorix/go-diameter/diam"
	"github.com/fiorix/go-diameter/diam/avp"
	"github.com/fiorix/go-diameter/diam/datatype"
)

type DiameterClient struct {
	config  DiameterConfig
	conn    diam.Conn
	handler *diam.ServeMux

	errorCh   chan error
	ceaCh     chan *diam.Message
	dwaCh     chan *diam.Message
	dwAliveCh chan *diam.Message
	ccaCh     chan *diam.Message
	inCh      chan Request
}

type DiameterConfig struct {
	URL              string
	OriginHost       datatype.DiameterIdentity
	OriginRealm      datatype.DiameterIdentity
	DestinationHost  datatype.DiameterIdentity
	DestinationRealm datatype.DiameterIdentity
	VendorID         datatype.Unsigned32
	ProductName      datatype.UTF8String
	FirmwareRevision datatype.Unsigned32
	WatchdogInterval time.Duration
}

func (d *DiameterClient) ErrorNotify() <-chan error {
	return d.errorCh
}

func (d *DiameterClient) CERDoneNotify() <-chan *diam.Message {
	return d.ceaCh
}

func (d *DiameterClient) DWRDoneNotify() <-chan *diam.Message {
	return d.dwaCh
}

func (d *DiameterClient) WatchdogAliveNotify() <-chan *diam.Message {
	return d.dwAliveCh
}

func (d *DiameterClient) CCRDoneNotify() <-chan *diam.Message {
	return d.ccaCh
}

func NewClient(config DiameterConfig) *DiameterClient {
	client := &DiameterClient{
		config: config,

		errorCh:   make(chan error),
		ceaCh:     make(chan *diam.Message),
		dwaCh:     make(chan *diam.Message),
		dwAliveCh: make(chan *diam.Message),
		ccaCh:     make(chan *diam.Message),
		inCh:      make(chan Request, 10),
	}
	client.handler = diam.NewServeMux()
	client.handler.Handle("CEA", client.HandleCEA())
	client.handler.Handle("DWA", client.HandleDWA())
	client.handler.Handle("DWR", client.HandleDWR())
	client.handler.Handle("CCA", client.HandleCCA())

	return client
}

func (d *DiameterClient) Start() error {
	var err error
	d.conn, err = diam.Dial(d.config.URL, d.handler, nil)
	if err != nil {
		return err
	}
	return nil
}

func (d *DiameterClient) Init() {
	d.SendCER()

	<-d.CERDoneNotify()
	go d.LoopWatchdog()
	go d.Listen()
}

func (d *DiameterClient) LoopWatchdog() {
	for {
		d.SendDWR()

		d.dwAliveCh <- <-d.DWRDoneNotify()
		time.Sleep(d.config.WatchdogInterval)
	}
}

type Request interface {
	AVP() []*diam.AVP
	ResponseNotify() <-chan *diam.Message
	Response(*diam.Message)
}

func (d *DiameterClient) Listen() {
	for {
		request := <-d.inCh
		d.SendCCR(request.AVP())
		request.Response(<-d.CCRDoneNotify())
	}
}

func (d *DiameterClient) Serve(request Request) {
	d.inCh <- request
}

func (d *DiameterClient) SendCER() {
	m := diam.NewRequest(diam.CapabilitiesExchange, 0, nil)

	m.NewAVP(avp.OriginHost, avp.Mbit, 0, d.config.OriginHost)
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, d.config.OriginRealm)

	ip, _, _ := net.SplitHostPort(d.conn.LocalAddr().String())
	m.NewAVP(avp.HostIPAddress, avp.Mbit, 0, datatype.Address(net.ParseIP(ip)))
	m.NewAVP(avp.VendorID, avp.Mbit, 0, d.config.VendorID)
	m.NewAVP(avp.ProductName, 0, 0, d.config.ProductName)
	m.NewAVP(avp.OriginStateID, avp.Mbit, 0, datatype.Unsigned32(0))
	m.NewAVP(avp.SupportedVendorID, avp.Mbit, 0, datatype.Unsigned32(0))
	m.NewAVP(avp.AuthApplicationID, avp.Mbit, 0, datatype.Unsigned32(4))
	m.NewAVP(avp.AcctApplicationID, avp.Mbit, 0, datatype.Unsigned32(4))
	m.NewAVP(avp.FirmwareRevision, avp.Mbit, 0, d.config.FirmwareRevision)

	_, err := m.WriteTo(d.conn)
	if err != nil {
		d.errorCh <- err
	}
}

func (d *DiameterClient) HandleCEA() diam.HandlerFunc {
	return func(conn diam.Conn, m *diam.Message) {
		d.ceaCh <- m
	}
}

func (d *DiameterClient) SendDWR() {
	m := diam.NewRequest(diam.DeviceWatchdog, 0, nil)

	m.NewAVP(avp.OriginHost, avp.Mbit, 0, d.config.OriginHost)
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, d.config.OriginRealm)

	_, err := m.WriteTo(d.conn)
	if err != nil {
		d.errorCh <- err
	}
}

func (d *DiameterClient) HandleDWA() diam.HandlerFunc {
	return func(conn diam.Conn, m *diam.Message) {
		d.dwaCh <- m
	}
}

func (d *DiameterClient) HandleDWR() diam.HandlerFunc {
	return func(conn diam.Conn, m *diam.Message) {
		answerMessage := m.Answer(diam.Success)
		d.SendDWA(conn, answerMessage)
	}
}

func (d *DiameterClient) SendDWA(w io.Writer, m *diam.Message) {
	m.NewAVP(avp.OriginHost, avp.Mbit, 0, d.config.OriginHost)
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, d.config.OriginRealm)
	_, err := m.WriteTo(w)

	if err != nil {
		d.errorCh <- err
	}
}

func (d *DiameterClient) SendCCR(avps []*diam.AVP) {
	sessionID := fmt.Sprintf("dtac.co.th;OMR%s001", time.Now().Format("20060102150405000"))

	m := diam.NewRequest(diam.CreditControl, 4, nil)

	m.NewAVP(avp.SessionID, avp.Mbit, 0, datatype.UTF8String(sessionID))
	m.NewAVP(avp.DestinationHost, avp.Mbit, 0, d.config.DestinationHost)
	m.NewAVP(avp.DestinationRealm, avp.Mbit, 0, d.config.DestinationRealm)
	m.NewAVP(avp.OriginHost, avp.Mbit, 0, d.config.OriginHost)
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, d.config.OriginRealm)
	for _, avp := range avps {
		m.AddAVP(avp)
	}

	_, err := m.WriteTo(d.conn)
	if err != nil {
		d.errorCh <- err
	}
}

func (d *DiameterClient) HandleCCA() diam.HandlerFunc {
	return func(conn diam.Conn, m *diam.Message) {
		d.ccaCh <- m
	}
}
