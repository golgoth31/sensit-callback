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
package sensitpayload

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"

	"forge.notrenet.com/domosense/sensit-callback/sensitTypes"

	"github.com/icza/bitio"
)

var downlinkMode = false

const batteryConst = 2.7

// Mode activated
// 0: alert
// 1: temperature/humidity
// 2: light
// 3: door
// 4: vibration
// 5: magnet
var Mode uint64

// Light mesured in Lux
var Light float64

// Temperature mesured in °C
var Temperature float64

// Humidity mesured in %
var Humidity float64

// Battery level in V
var Battery float64

// FrameType ...
// 0: periodic
// 1: double click
// 2: alert
// 3: mode change
var FrameType uint64

// UplinkPeriod is time between 2 messages
// 0: 10m
// 1: 1h
// 2: 6h
// 3: 1d
var UplinkPeriod uint64

// Firmware version
var Firmware string

// Uplink data
// Uplink byte 3
var lightMask uint64
var lightValue uint64
var magnetMax uint64
var spare uint64
var magnetState uint64
var temperatureLsb uint64

// Uplink byte 4
var humidityRaw uint64
var firmwareMajor uint64
var firmwareMinor uint64

// AlertCounter for door
var AlertCounter uint64

// Downlink data
// Downlink byte 1
var sendingPeriodMsb uint64
var tempAlertLow uint64

// Downlink byte 2
var sendingPeriodLSB uint64
var tempAlertHigh uint64

// Downlink byte 3
var lightAlertLow uint64

// Downlink byte 4
var lightAlertHigh uint64

// Downlink byte 5
var accTransientThr uint64

// Downlink byte 6
var accTransientCount uint64

// Downlink byte 7
var disableHighPassFilter uint64
var enableDetectionZ uint64
var enableDetectionY uint64
var enableDetectionX uint64
var dstaRate uint64
var mods uint64

// Downlink byte 8
var reservedZero uint64
var doorThresholdAlert uint64

var DecodeData = true

// Decode paylaod
func Decode(pc chan sensittypes.CallbackData, o chan []byte) {
	log.Print("[DEBUG] Starting decode module")
	for DecodeData {
		p := <-pc
		decoded, err := hex.DecodeString(p.Data)
		payloadIoReader := bytes.NewBuffer(decoded)
		payloadBitReader := bitio.NewReader(payloadIoReader)

		// byte 1
		batteryLevelMsb, err := payloadBitReader.ReadBits(1)
		FrameType, err = payloadBitReader.ReadBits(2)
		UplinkPeriod, err = payloadBitReader.ReadBits(2)
		Mode, err = payloadBitReader.ReadBits(3)

		// byte 2
		temperatureMsb, err := payloadBitReader.ReadBits(4)
		batteryLevelLsb, err := payloadBitReader.ReadBits(4)

		// byte 3
		switch {
		case Mode == 2:
			lightMask, err = payloadBitReader.ReadBits(2)
			lightValue, err = payloadBitReader.ReadBits(6)
		case Mode == 3:
			magnetMax, err = payloadBitReader.ReadBits(8)
		default:
			spare, err = payloadBitReader.ReadBits(1)
			magnetState, err = payloadBitReader.ReadBits(1)
			temperatureLsb, err = payloadBitReader.ReadBits(6)
		}

		// byte 4
		switch {
		case Mode == 0:
			firmwareMajor, err = payloadBitReader.ReadBits(4)
			firmwareMinor, err = payloadBitReader.ReadBits(4)
			Firmware = fmt.Sprintf("%d.%d", firmwareMajor, firmwareMinor)
		case Mode == 1:
			humidityRaw, err = payloadBitReader.ReadBits(8)
		default:
			AlertCounter, err = payloadBitReader.ReadBits(8)
		}

		// Downlink byte 1
		// test the first bit of the downlink payload: if "EOF", no need to go further in decoding
		sendingPeriodMsb, err = payloadBitReader.ReadBits(1)
		if err == nil {
			log.Print("[DEBUG] Downlink requested")
			downlinkMode = true
		} else if fmt.Sprint(err) == "EOF" {
			downlinkMode = false
		} else {
			log.Printf("[CRIT] %v", err)
		}
		if downlinkMode {
			// Downlink byte 1
			tempAlertLow, err = payloadBitReader.ReadBits(7)

			// Downlink byte 2
			sendingPeriodLSB, err = payloadBitReader.ReadBits(1)
			tempAlertHigh, err = payloadBitReader.ReadBits(7)

			// Downlink byte 3
			lightAlertLow, err = payloadBitReader.ReadBits(8)

			// Downlink byte 4
			lightAlertHigh, err = payloadBitReader.ReadBits(8)

			// Downlink byte 5
			accTransientThr, err = payloadBitReader.ReadBits(8)

			// Downlink byte 6
			accTransientCount, err = payloadBitReader.ReadBits(8)

			// Downlink byte 7
			disableHighPassFilter, err = payloadBitReader.ReadBits(1)
			enableDetectionZ, err = payloadBitReader.ReadBits(1)
			enableDetectionY, err = payloadBitReader.ReadBits(1)
			enableDetectionX, err = payloadBitReader.ReadBits(1)
			dstaRate, err = payloadBitReader.ReadBits(2)
			mods, err = payloadBitReader.ReadBits(2)

			// Downlink byte 8
			reservedZero, err = payloadBitReader.ReadBits(1)
			doorThresholdAlert, err = payloadBitReader.ReadBits(7)

			// fmt.Printf("temp => %v", sendingPeriodMsb)
			// fmt.Printf("humidity => %v", tempAlertLow)
		}

		// build battery data
		var batteryByte bytes.Buffer
		batteryBitWriter := bitio.NewWriter(&batteryByte)
		batteryBitWriter.WriteBits(0, 3)
		batteryBitWriter.WriteBits(batteryLevelMsb, 1)
		batteryBitWriter.WriteBits(batteryLevelLsb, 4)
		batteryBitWriter.Close()
		batteryRaw, _ := binary.Uvarint(batteryByte.Bytes())
		Battery = float64(batteryRaw)*0.05 + batteryConst

		// Print common information
		fmt.Printf("Bat => %v\n", Battery)
		fmt.Printf("FrameType => %v\n", FrameType)
		fmt.Printf("UplinkPeriod => %v\n", UplinkPeriod)
		fmt.Printf("Mode => %v\n", Mode)
		var out []byte
		switch Mode {
		// case 0:
		// 	fmt.Printf("Spare => %v", spare)
		// 	fmt.Printf("Magnet => %v", magnetState)
		case 1:
			var temperatureByte bytes.Buffer
			temperatureBitWriter := bitio.NewWriter(&temperatureByte)
			temperatureBitWriter.WriteBits(0, 6)
			temperatureBitWriter.WriteBits(temperatureMsb, 4)
			temperatureBitWriter.WriteBits(temperatureLsb, 6)
			temperatureBitWriter.Close()
			Temperature = float64(float64(binary.BigEndian.Uint16(temperatureByte.Bytes())-200) / 8)
			Humidity = float64(humidityRaw) / 2
			// if err != nil {
			// 	log.Print(err)
			// }
			fmt.Printf("tempbytes => %v", binary.BigEndian.Uint16(temperatureByte.Bytes()))
			fmt.Printf("temp => %v\n", Temperature)
			fmt.Printf("humidity => %v\n", Humidity)
			outdata := sensittypes.SensitTempData{
				Device: p.Device,
				Time:   p.Timestamp,
				Temp:   Temperature,
				Bat:    Battery,
				Hum:    Humidity,
				Mode:   "temperature",
			}
			out, err = json.Marshal(outdata)
		case 2:
			switch lightMask {
			case 0:
				Light = float64(lightValue) / 96
			case 1:
				Light = float64(lightValue) * 8 / 96
			case 2:
				Light = float64(lightValue) * 64 / 96
			case 3:
				Light = float64(lightValue) * 1024 / 96
			}
			// fmt.Printf("Light => %v", Light)
			// case 3:
			// 	fmt.Printf("Magnet => %v", magnetMax)
		}
		// p <- "done"
		o <- out
	}
}
