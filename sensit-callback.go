// Copyright © 2018 David Sabatie <david.sabatie@notrenet.com>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"log"
	"os"

	"github.com/golgoth31/sensit-callback/config"
	"github.com/golgoth31/sensit-callback/input/aws"

	"github.com/golgoth31/sensit-callback/output"
	"github.com/golgoth31/sensit-callback/payload"
	"github.com/hashicorp/logutils"
)

var cfg = config.Config

var LogFilter *logutils.LevelFilter

func init() {
	LogFilter = &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "INFO", "WARN", "ERROR", "CRIT"},
		MinLevel: logutils.LogLevel("INFO"),
		Writer:   os.Stderr,
	}
}

func main() {
	// init config
	// config.InitConfig()

	// domolib.FailOnError(err, "Failed to parse config file")

	// Change log level to configured level
	LogFilter.SetMinLevel(logutils.LogLevel(cfg.GetString("log.level")))
	log.SetOutput(LogFilter)
	log.Print("[DEBUG] Starting to listen for sensit callback")

	// start payload
	go sensitpayload.Decode(config.PayloadChan, config.OutputChan)

	// write output
	sensitoutput.Start()

	// listen sqs
	sensitsqs.GetMessage(cfg.GetString("input.mode.readqURL"), config.PayloadChan)
}
