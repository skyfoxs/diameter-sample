package dcc

import (
	"net"

	"github.com/fiorix/go-diameter/diam"
	"github.com/fiorix/go-diameter/diam/avp"
	"github.com/fiorix/go-diameter/diam/datatype"
)

type DiameterClient struct {
	URL     string
	conn    diam.Conn
	handler *diam.ServeMux

	errorCh chan error
	cerCh   chan *diam.Message
}

func (d *DiameterClient) ErrorNotify() <-chan error {
	return d.errorCh
}

func (d *DiameterClient) CERDoneNotify() <-chan *diam.Message {
	return d.cerCh
}

func NewClient(address string) *DiameterClient {
	client := &DiameterClient{
		URL: address,

		errorCh: make(chan error),
		cerCh:   make(chan *diam.Message),
	}
	client.handler = diam.NewServeMux()
	client.handler.Handle("CEA", client.HandleCEA())

	return client
}

func (d *DiameterClient) Start() {
	var err error
	d.conn, err = diam.Dial(d.URL, d.handler, nil)
	if err != nil {
		d.errorCh <- err
	}
}

func (d *DiameterClient) SendCER() {
	var (
		CapabilitiesExchange uint32 = 257
		OriginHost                  = datatype.DiameterIdentity("jenkin13_OMR_TEST01")
		OriginRealm                 = datatype.DiameterIdentity("dtac.co.th")
		vendorID                    = datatype.Unsigned32(0)
		productName                 = datatype.UTF8String("omr")
		FirmwareRevision            = datatype.Unsigned32(1)
	)

	m := diam.NewRequest(CapabilitiesExchange, 0, nil)

	m.NewAVP(avp.OriginHost, avp.Mbit, 0, OriginHost)
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, OriginRealm)

	ip, _, _ := net.SplitHostPort(d.conn.LocalAddr().String())
	m.NewAVP(avp.HostIPAddress, avp.Mbit, 0, datatype.Address(net.ParseIP(ip)))
	m.NewAVP(avp.VendorID, avp.Mbit, 0, vendorID)
	m.NewAVP(avp.ProductName, 0, 0, productName)
	m.NewAVP(avp.OriginStateID, avp.Mbit, 0, datatype.Unsigned32(0))
	m.NewAVP(avp.SupportedVendorID, avp.Mbit, 0, datatype.Unsigned32(0))
	m.NewAVP(avp.AuthApplicationID, avp.Mbit, 0, datatype.Unsigned32(4))
	m.NewAVP(avp.AcctApplicationID, avp.Mbit, 0, datatype.Unsigned32(4))
	m.NewAVP(avp.FirmwareRevision, avp.Mbit, 0, FirmwareRevision)

	_, err := m.WriteTo(d.conn)
	if err != nil {
		d.errorCh <- err
	}
}

func (d *DiameterClient) HandleCEA() diam.HandlerFunc {
	return func(conn diam.Conn, m *diam.Message) {
		d.cerCh <- m
	}
}
