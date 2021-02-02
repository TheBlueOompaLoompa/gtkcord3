package keyring

import (
	"github.com/TheBlueOompaLoompa/gtkcord3/internal/log"
	"github.com/zalando/go-keyring"
)

func Get() string {
	k, err := keyring.Get("gtkcord", "token")
	if err != nil {
		log.Errorln("[non-fatal] Failed to get Gtkcord token from keyring")
	}

	if k == "" {
		log.Infoln("Keyring token is empty.")
	}

	return k
}

func Set(token string) {
	if err := keyring.Set("gtkcord", "token", token); err != nil {
		log.Errorln("[non-fatal] Failed to set Gtkcord token to keyring")
	}
}

func Delete() {
	if err := keyring.Delete("gtkcord", "token"); err != nil {
		log.Errorln("[non-fatal] Failed to delete Gtkcord token from keyring")
	}
}
