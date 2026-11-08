package rtc

import (
	"github.com/pion/webrtc/v3"
)

type Media struct {
	rtc
}

func (m *Media) CreateOffer() (*webrtc.SessionDescription, error) {
	if _, err := m.pc.AddTransceiverFromKind(webrtc.RTPCodecTypeAudio, webrtc.RTPTransceiverInit{Direction: webrtc.RTPTransceiverDirectionRecvonly}); err != nil {
		return nil, err
	}
	if _, err := m.pc.AddTransceiverFromKind(webrtc.RTPCodecTypeVideo, webrtc.RTPTransceiverInit{Direction: webrtc.RTPTransceiverDirectionRecvonly}); err != nil {
		return nil, err
	}
	return m.createOffer()
}

func (m *Media) AddTrack(track webrtc.TrackLocal) (*webrtc.RTPTransceiver, error) {
	return m.pc.AddTransceiverFromTrack(track)
}

func (m *Media) OnTrack(f func(*webrtc.TrackRemote, *webrtc.RTPReceiver)) {
	m.pc.OnTrack(f)
}

func NewMedia(conf *webrtc.Configuration, onCandidate func(*webrtc.ICECandidateInit), opts ...Option) (*Media, error) {
	m := &webrtc.MediaEngine{}
	if err := m.RegisterDefaultCodecs(); err != nil {
		return nil, err
	}
	r, err := newRTC(webrtc.NewAPI(webrtc.WithMediaEngine(m)), conf, onCandidate, opts...)
	if err != nil {
		return nil, err
	}
	return &Media{*r}, nil
}
