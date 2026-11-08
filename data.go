package rtc

import "github.com/pion/webrtc/v3"

type Data struct {
	rtc
	onData func(*webrtc.DataChannel)
}

func (d *Data) CreateOffer() (*webrtc.SessionDescription, error) {
	dc, err := d.pc.CreateDataChannel("data", &webrtc.DataChannelInit{Ordered: &d.ordered})
	if err != nil {
		return nil, err
	}
	d.onData(dc)
	return d.createOffer()
}

func (d *Data) OnOffer(offer *webrtc.SessionDescription) (*webrtc.SessionDescription, error) {
	d.pc.OnDataChannel(d.onData)
	return d.rtc.OnOffer(offer)
}

func NewData(conf *webrtc.Configuration, onCandidate func(*webrtc.ICECandidateInit), onData func(*webrtc.DataChannel), opts ...Option) (*Data, error) {
	r, err := newRTC(webrtc.NewAPI(), conf, onCandidate, opts...)
	if err != nil {
		return nil, err
	}
	return &Data{*r, onData}, nil
}
