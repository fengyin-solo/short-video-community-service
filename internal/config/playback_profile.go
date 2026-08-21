package config

type PlaybackProfile struct {
	Device string
	Codecs map[string]bool
	MaxK   int
}

func LoadPlaybackProfile(device string) *PlaybackProfile {
	if device == "" {
		return &PlaybackProfile{Device: "unknown", MaxK: 4}
	}
	return &PlaybackProfile{Device: device, Codecs: map[string]bool{"h264": true}, MaxK: 2}
}
