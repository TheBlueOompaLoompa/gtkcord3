package variables

import "github.com/gotk3/gotk3/gtk"

var (
	WindowWidth  = 720
	WindowHeight = 900

	EmbedAvatarSize = 24
	EmbedMaxWidth   = 300
	EmbedImgHeight  = 300 // max
	EmbedMargin     = 8

	AvatarSize    = 42 // gtk.ICON_SIZE_DND
	AvatarPadding = 10

	// used as fallback, the settings one overrides this
	MaxMessageWidth = 750

	SpinnerSize  = 56
	ChannelWidth = 240

	InputIconSize = gtk.ICON_SIZE_LARGE_TOOLBAR
)
