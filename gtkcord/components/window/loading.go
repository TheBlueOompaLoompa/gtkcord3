package window

import (
	"github.com/TheBlueOompaLoompa/gtkcord3/gtkcord/components/animations"
	"github.com/gotk3/gotk3/gtk"
)

const LoadingTitle = "Connecting to Discord."

// NowLoading fades the internal stack view to show a spinning circle.
func NowLoading() {
	// Use a spinner:
	s, _ := animations.NewSpinner(75)

	// Use a custom header instead of the actual Header:
	h, _ := gtk.HeaderBarNew()
	h.SetTitle(LoadingTitle)
	h.SetShowCloseButton(true)
	h.ShowAll()

	// Set the loading animation:
	stackSet(Window.Main, "loading", s)
	stackSet(Window.Header.Main, "loading", h)
}
