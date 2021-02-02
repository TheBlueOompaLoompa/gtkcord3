package main

import (
	"os"
	"strconv"

	"github.com/TheBlueOompaLoompa/gtkcord3/gtkcord/components/window"
	"github.com/TheBlueOompaLoompa/gtkcord3/gtkcord/variables"
	"github.com/TheBlueOompaLoompa/gtkcord3/internal/log"
)

func LoadEnvs() {
	if css := os.Getenv("GTKCORD_CUSTOM_CSS"); css != "" {
		window.CustomCSS = css
	}

	if w, _ := strconv.Atoi(os.Getenv("GTKCORD_MSGWIDTH")); w > 100 { // min 100
		variables.MaxMessageWidth = w
	}

	if os.Getenv("GTKCORD_QUIET") == "0" {
		log.Quiet = false
		profile = true
	}
}
