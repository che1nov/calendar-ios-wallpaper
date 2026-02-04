package domain

type BackgroundStyle string

const (
	BgPlain    BackgroundStyle = "plain"
	BgGradient BackgroundStyle = "gradient"
	BgNoise    BackgroundStyle = "noise"
	BgIOS      BackgroundStyle = "ios"
)

func ParseBackgroundStyle(v string) BackgroundStyle {
	switch BackgroundStyle(v) {
	case BgPlain, BgGradient, BgNoise, BgIOS:
		return BackgroundStyle(v)
	default:
		return BgIOS
	}
}
