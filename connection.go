package rtc

import "github.com/pion/webrtc/v3"

type connection struct {
	pc *webrtc.PeerConnection
}

func (c *connection) createOffer() (*webrtc.SessionDescription, error) {
	offer, err := c.pc.CreateOffer(nil)
	if err != nil {
		return nil, err
	}
	if err = c.pc.SetLocalDescription(offer); err != nil {
		return nil, err
	}
	return &offer, nil
}

func (c *connection) OnOffer(offer *webrtc.SessionDescription) (*webrtc.SessionDescription, error) {
	if err := c.pc.SetRemoteDescription(*offer); err != nil {
		return nil, err
	}
	answer, err := c.pc.CreateAnswer(nil)
	if err != nil {
		return nil, err
	}
	if err = c.pc.SetLocalDescription(answer); err != nil {
		return nil, err
	}
	return &answer, nil
}

func (c *connection) OnAnswer(answer *webrtc.SessionDescription) error {
	return c.pc.SetRemoteDescription(*answer)
}

func (c *connection) OnCandidate(candidate *webrtc.ICECandidateInit) error {
	return c.pc.AddICECandidate(*candidate)
}

func (c *connection) OnConnectionStateChange(f func(webrtc.PeerConnectionState)) {
	c.pc.OnConnectionStateChange(f)
}

func (c *connection) Close() error {
	c.pc.OnConnectionStateChange(nil)
	return c.pc.Close()
}

func newConnection(api *webrtc.API, conf *webrtc.Configuration, onCandidate func(*webrtc.ICECandidateInit)) (*connection, error) {
	pc, err := api.NewPeerConnection(*conf)
	if err != nil {
		return nil, err
	}
	pc.OnICECandidate(func(candidate *webrtc.ICECandidate) {
		if candidate != nil {
			c := candidate.ToJSON()
			onCandidate(&c)
		} else {
			pc.OnICECandidate(nil)
		}
	})
	return &connection{pc}, nil
}
