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
package sensittypes

import "time"

type CallbackData struct {
	Device    string    `json:"device"`
	Time      string    `json:"time"`
	Data      string    `json:"data"`
	MessageID string    `json:"messageId"`
	Ack       string    `json:"ack"`
	Timestamp time.Time `json:",omitempty"`
}

type SensitTempData struct {
	Device string    `json:""`
	Time   time.Time `json:""`
	Temp   float64   `json:""`
	Bat    float64   `json:""`
	Hum    float64   `json:""`
	Mode   string    `json:""`
}

// SensitioData structure (callbackType = sensitio)
type SensitioData struct {
	MessageID        string `json:"messageID,omitempty"`
	ReceivedAt       string `json:",omitempty"`
	DeviceID         string `json:"deviceId,omitempty"`
	SigfoxID         string `json:"sigfoxId,omitempty"`
	DeviceName       string `json:",omitempty"`
	Payload          string `json:",omitempty"`
	BatteryIndicator string `json:",omitempty"`
	Value            string `json:",omitempty"`
	Sensor           string `json:",omitempty"`
	State            string `json:",omitempty"`
	Events           string `json:",omitempty"`
}

// SigfoxDataBidir structure (callbackType = sigfox)
type SigfoxDataBidir struct {
	Device    string `json:",omitempty"`
	Data      string `json:",omitempty"`
	SeqNumber string `json:",omitempty"`
	Ack       string `json:",omitempty"`
}

// SigfoxDataAdvanced ...
type SigfoxDataAdvanced struct {
	Device           string `json:",omitempty"`
	Time             string `json:",omitempty"`
	MessageID        string `json:"messageID,omitempty"`
	Data             string `json:",omitempty"`
	Lqi              string `json:",omitempty"`
	FixedLat         string `json:",omitempty"`
	FixedLng         string `json:",omitempty"`
	OperatorName     string `json:",omitempty"`
	CountryCode      string `json:",omitempty"`
	ComputedLocation SigfoxComputedLocation
}

// SigfoxComputedLocation ...
type SigfoxComputedLocation struct {
	Lat    string `json:",omitempty"`
	Lng    string `json:",omitempty"`
	Radius string `json:",omitempty"`
	Source string `json:",omitempty"`
	Status string `json:",omitempty"`
}
