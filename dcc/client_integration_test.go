// +build integration

package dcc

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/fiorix/go-diameter/diam"
	"github.com/fiorix/go-diameter/diam/avp"
	"github.com/fiorix/go-diameter/diam/datatype"

	"github.com/fiorix/go-diameter/diam/dict"

	"github.com/skyfoxs/diameter-sample/dcc/dictionary"
)

func TestClientIntegration(t *testing.T) {
	dict.Default.Load(bytes.NewBufferString(dictionary.AppDictionary))
	dict.Default.Load(bytes.NewBufferString(dictionary.CreditControlDictionary))

	client := NewClient(DiameterConfig{
		URL:              "10.89.104.33:6553",
		OriginHost:       datatype.DiameterIdentity("jenkin13_OMR_TEST01"),
		OriginRealm:      datatype.DiameterIdentity("dtac.co.th"),
		DestinationHost:  datatype.DiameterIdentity("cbp241"),
		DestinationRealm: datatype.DiameterIdentity("www.huawei.com"),
		VendorID:         datatype.Unsigned32(0),
		ProductName:      datatype.UTF8String("omr"),
		FirmwareRevision: datatype.Unsigned32(1),
		WatchdogInterval: 3 * time.Second,
	})
	if err := client.Start(); err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	client.Init()

	request := &TestRequest{
		outCh: make(chan *diam.Message),
	}
	client.Serve(request)

	select {
	case err := <-client.ErrorNotify():
		t.Error(err)
	case m := <-request.ResponseNotify():
		fmt.Println(m)
	case <-time.After(time.Second):
		t.Error("server timeout")
	}

}

type TestRequest struct {
	outCh chan *diam.Message
}

func (r *TestRequest) Response(m *diam.Message) {
	r.outCh <- m
}

func (r *TestRequest) ResponseNotify() <-chan *diam.Message {
	return r.outCh
}

func (r *TestRequest) AVP() []*diam.AVP {
	return []*diam.AVP{
		diam.NewAVP(30951, avp.Mbit, 0, datatype.Integer64(625004290)),
		diam.NewAVP(avp.SubscriptionID, avp.Mbit, 0, &diam.GroupedAVP{
			AVP: []*diam.AVP{
				diam.NewAVP(avp.SubscriptionIDType, avp.Mbit, 0, datatype.Enumerated(0)),
				diam.NewAVP(avp.SubscriptionIDData, avp.Mbit, 0, datatype.UTF8String("66906300719")),
			},
		}),
		diam.NewAVP(avp.ServiceInformation, avp.Mbit, 0, &diam.GroupedAVP{
			AVP: []*diam.AVP{
				diam.NewAVP(21100, avp.Mbit, 0, &diam.GroupedAVP{
					AVP: []*diam.AVP{
						diam.NewAVP(20336, avp.Mbit, 0, datatype.UTF8String("66906300719")),
						diam.NewAVP(20340, avp.Mbit, 0, datatype.Unsigned32(5)),
						diam.NewAVP(20386, avp.Mbit, 0, datatype.Time(time.Now())),
					},
				}),
			},
		}),
		diam.NewAVP(avp.AuthApplicationID, avp.Mbit, 0, datatype.Unsigned32(4)),
		diam.NewAVP(avp.CCRequestType, avp.Mbit, 0, datatype.Enumerated(4)),
		diam.NewAVP(avp.ServiceContextID, avp.Mbit, 0, datatype.UTF8String("QuerySubinfo@huawei.com")),
		diam.NewAVP(avp.RequestedAction, avp.Mbit, 0, datatype.Enumerated(1)),
		diam.NewAVP(avp.EventTimestamp, avp.Mbit, 0, datatype.Time(time.Now())),
		diam.NewAVP(avp.ServiceIdentifier, avp.Mbit, 0, datatype.Unsigned32(0)),
		diam.NewAVP(avp.CCRequestNumber, avp.Mbit, 0, datatype.Unsigned32(0)),
	}
}
