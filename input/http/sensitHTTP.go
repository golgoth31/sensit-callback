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
package sensithttp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func GetMessage() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			log.Print("[WARN] page not found")
			http.NotFound(w, r)
			return
		}
		if r.Header["Authorization"] != nil {
			if r.Header["Authorization"][0] == fmt.Sprintf("bearer %s", bearerToken) {
				body, err := ioutil.ReadAll(r.Body)
				LogOnError(err, "Can't read request body")
				err = json.Unmarshal(body, &sensitData)
				LogOnError(err, "Can't parse request body")
				log.Printf("[INFO] %v", sensitData)
			} else {
				log.Printf("[WARN] Bad bearer: %v", r.Header["Authorization"][0])
			}
		} else {
			log.Print("[WARN] No authorization given")
			http.NotFound(w, r)
			return
		}
	})
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", callbackPort), nil))
}
