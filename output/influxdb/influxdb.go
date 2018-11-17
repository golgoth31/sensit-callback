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
package sensitinfluxdb

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/golgoth31/sensit-callback/config"
	"github.com/golgoth31/sensit-callback/sensitTypes"

	influxClient "github.com/influxdata/influxdb/client/v2"
)

var InfluxDBname = "sensit"
var InfluxDBretention = "3d"
var cfg = config.Config

// func init() {

// }
// queryDB convenience function to query the database
func queryDB(clnt influxClient.Client, cmd string) (res []influxClient.Result, err error) {
	q := influxClient.Query{
		Command:  cmd,
		Database: InfluxDBname,
	}
	if response, err := clnt.Query(q); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	} else {
		return res, err
	}
	return res, nil
}

func Write(outputChan chan []byte) {
	// config.InitConfig()
	write := true
	for write {
		o := <-outputChan
		log.Print("[DEBUG] Starting output module")
		// Make client
		clnt, err := influxClient.NewHTTPClient(influxClient.HTTPConfig{
			Addr: cfg.GetString("output.influxdb.host"),
		})
		if err != nil {
			fmt.Println("Error creating InfluxDB Client: ", err.Error())
		}
		defer clnt.Close()

		_, err = queryDB(clnt, fmt.Sprintf("CREATE DATABASE %s WITH DURATION %s", InfluxDBname, InfluxDBretention))
		if err != nil {
			log.Fatal(err)
		}

		var out sensittypes.SensitTempData
		err = json.Unmarshal(o, &out)
		log.Println(out)
		// Create a new point batch
		bp, _ := influxClient.NewBatchPoints(influxClient.BatchPointsConfig{
			Database:  "sensit",
			Precision: "ns",
		})

		// Create a point and add to batch
		tags := map[string]string{"sensit-mode": out.Mode}
		fields := map[string]interface{}{
			"temp":     out.Temp,
			"bat":      out.Bat,
			"humidity": out.Hum,
		}
		pt, err := influxClient.NewPoint("sensit", tags, fields, out.Time)
		if err != nil {
			fmt.Println("Error: ", err.Error())
		}
		bp.AddPoint(pt)

		// Write the batch
		clnt.Write(bp)
	}
}
