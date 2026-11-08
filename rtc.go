package rtc

import (
	"time"

	"github.com/pion/webrtc/v3"
)

type rtc struct {
	pc          *webrtc.PeerConnection
	timeout     time.Duration
	ordered     bool
	onConnected func()
	onTimeout   func()
	onClosed    func()
}

func (r *rtc) createOffer() (*webrtc.SessionDescription, error) {
	offer, err := r.pc.CreateOffer(nil)
	if err != nil {
		return nil, err
	}
	if err = r.pc.SetLocalDescription(offer); err != nil {
		return nil, err
	}
	return &offer, nil
}

func (c *rtc) OnOffer(offer *webrtc.SessionDescription) (*webrtc.SessionDescription, error) {
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

func (r *rtc) OnAnswer(answer *webrtc.SessionDescription) error {
	return r.pc.SetRemoteDescription(*answer)
}

func (r *rtc) OnCandidate(candidate *webrtc.ICECandidateInit) error {
	return r.pc.AddICECandidate(*candidate)
}

func (r *rtc) Close() error {
	return r.pc.Close()
}

func newRTC(api *webrtc.API, conf *webrtc.Configuration, onCandidate func(*webrtc.ICECandidateInit), opts ...Option) (*rtc, error) {
	pc, err := api.NewPeerConnection(*conf)
	if err != nil {
		return nil, err
	}
	r := &rtc{pc: pc, timeout: 30 * time.Second}
	for _, o := range opts {
		o(r)
	}
	pc.OnICECandidate(func(candidate *webrtc.ICECandidate) {
		if candidate != nil {
			c := candidate.ToJSON()
			onCandidate(&c)
		} else {
			pc.OnICECandidate(nil)
		}
	})
	var timer *time.Timer
	timer = time.AfterFunc(r.timeout, func() {
		timer = nil
		if r.onTimeout != nil {
			r.onTimeout()
		}
	})
	pc.OnConnectionStateChange(func(state webrtc.PeerConnectionState) {
		switch state {
		case webrtc.PeerConnectionStateConnected:
			if timer != nil {
				timer.Stop()
				timer = nil
			}
			if r.onConnected != nil {
				r.onConnected()
			}
		case webrtc.PeerConnectionStateDisconnected, webrtc.PeerConnectionStateFailed, webrtc.PeerConnectionStateClosed:
			pc.OnICECandidate(nil)
			pc.OnConnectionStateChange(nil)
			if timer != nil {
				timer.Stop()
				timer = nil
			}
			if r.onClosed != nil {
				r.onClosed()
			}
		}
	})
	return r, nil
}
