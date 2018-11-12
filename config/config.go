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
package config

import (
	"fmt"

	sensittypes "forge.notrenet.com/domosense/sensit-callback/sensitTypes"
	"github.com/spf13/viper"
)

var Config = viper.New()
var PayloadChan chan sensittypes.CallbackData
var OutputChan chan []byte

func InitConfig() {
	Config.SetConfigName("config")                 // name of config file (without extension)
	Config.AddConfigPath("/etc/sensit-callback/")  // path to look for the config file in
	Config.AddConfigPath("$HOME/.sensit-callback") // call multiple times to add many search paths
	Config.AddConfigPath(".")                      // optionally look for config in the working directory
	err := Config.ReadInConfig()                   // Find and read the config file
	if err != nil {                                // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s", err))
	}
	PayloadChan = make(chan sensittypes.CallbackData)
	OutputChan = make(chan []byte)
}
