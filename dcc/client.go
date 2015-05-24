package dcc

import (
	"net"
	"time"

	"github.com/fiorix/go-diameter/diam"
	"github.com/fiorix/go-diameter/diam/avp"
	"github.com/fiorix/go-diameter/diam/datatype"
)

var (
	CapabilitiesExchange uint32 = 257
	watchdogExchange     uint32 = 280
)

type DiameterClient struct {
	config  DiameterConfig
	conn    diam.Conn
	handler *diam.ServeMux

	errorCh   chan error
	ceaCh     chan *diam.Message
	dwaCh     chan *diam.Message
	dwAliveCh chan *diam.Message
}

type DiameterConfig struct {
	URL              string
	OriginHost       datatype.DiameterIdentity
	OriginRealm      datatype.DiameterIdentity
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

func NewClient(config DiameterConfig) *DiameterClient {
	client := &DiameterClient{
		config: config,

		errorCh:   make(chan error),
		ceaCh:     make(chan *diam.Message),
		dwaCh:     make(chan *diam.Message),
		dwAliveCh: make(chan *diam.Message),
	}
	client.handler = diam.NewServeMux()
	client.handler.Handle("CEA", client.HandleCEA())
	client.handler.Handle("DWA", client.HandleDWA())

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
}

func (d *DiameterClient) LoopWatchdog() {
	for {
		d.SendDWR()

		d.dwAliveCh <- <-d.DWRDoneNotify()
		time.Sleep(d.config.WatchdogInterval)
	}
}

func (d *DiameterClient) SendCER() {
	m := diam.NewRequest(CapabilitiesExchange, 0, nil)

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
	m := diam.NewRequest(watchdogExchange, 0, nil)

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
