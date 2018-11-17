package sensitoutput

import (
	"fmt"
	"log"

	"github.com/golgoth31/sensit-callback/config"
	"github.com/golgoth31/sensit-callback/output/influxdb"
)

var cfg = config.Config

func Start() {
	for k := range cfg.GetStringMap("output") {
		log.Printf("[DEBUG] Loading output module: %s", k)
		switch k {
		case "influxdb":
			if cfg.GetBool(fmt.Sprintf("output.%s.enabled", k)) {
				go sensitinfluxdb.Write(config.OutputChan)
			}
		case "default":
			panic("unknow module")
		}
	}
}
