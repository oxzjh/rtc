package rtc

import "github.com/pion/webrtc/v3"

type DataConnection struct {
	*connection
	ordered bool
	onData  func(*webrtc.DataChannel)
}

func (c *DataConnection) CreateOffer() (*webrtc.SessionDescription, error) {
	dc, err := c.pc.CreateDataChannel("data", &webrtc.DataChannelInit{Ordered: &c.ordered})
	if err != nil {
		return nil, err
	}
	c.onData(dc)
	return c.createOffer()
}

func (c *DataConnection) OnOffer(offer *webrtc.SessionDescription) (*webrtc.SessionDescription, error) {
	c.pc.OnDataChannel(c.onData)
	return c.connection.OnOffer(offer)
}

func NewDataConnection(conf *webrtc.Configuration, onCandidate func(*webrtc.ICECandidateInit), ordered bool, onData func(*webrtc.DataChannel)) (*DataConnection, error) {
	connection, err := newConnection(webrtc.NewAPI(), conf, onCandidate)
	if err != nil {
		return nil, err
	}
	return &DataConnection{connection, ordered, onData}, nil
}
