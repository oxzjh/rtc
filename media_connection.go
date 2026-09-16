package rtc

import "github.com/pion/webrtc/v3"

type MediaConnection struct {
	*connection
}

func (c *MediaConnection) CreateOffer() (*webrtc.SessionDescription, error) {
	if _, err := c.pc.AddTransceiverFromKind(webrtc.RTPCodecTypeAudio, webrtc.RTPTransceiverInit{Direction: webrtc.RTPTransceiverDirectionRecvonly}); err != nil {
		return nil, err
	}
	if _, err := c.pc.AddTransceiverFromKind(webrtc.RTPCodecTypeVideo, webrtc.RTPTransceiverInit{Direction: webrtc.RTPTransceiverDirectionRecvonly}); err != nil {
		return nil, err
	}
	return c.createOffer()
}

func (c *MediaConnection) AddTrack(track webrtc.TrackLocal) (*webrtc.RTPTransceiver, error) {
	return c.pc.AddTransceiverFromTrack(track)
}

func (c *MediaConnection) OnTrack(f func(*webrtc.TrackRemote, *webrtc.RTPReceiver)) {
	c.pc.OnTrack(f)
}

func NewMediaConnection(conf *webrtc.Configuration, onCandidate func(*webrtc.ICECandidateInit)) (*MediaConnection, error) {
	m := &webrtc.MediaEngine{}
	if err := m.RegisterDefaultCodecs(); err != nil {
		return nil, err
	}
	connection, err := newConnection(webrtc.NewAPI(webrtc.WithMediaEngine(m)), conf, onCandidate)
	if err != nil {
		return nil, err
	}
	return &MediaConnection{connection}, nil
}
