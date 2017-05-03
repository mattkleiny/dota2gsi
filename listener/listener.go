// Copyright © 2017 Matthew Kleinschafer
// 
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the “Software”), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
// 
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
// 
// THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package listener

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

// Starts a HTTP server listening on the given port for game status updates
// Calls the given callback with the provided game state synchronously
// TODO: make this async/use go routines instead
func StartListener(port int, callback func(state *GameState)) {
	// create a handler for processing incoming requests and forwarding it to our callback
	handler := func(writer http.ResponseWriter, request *http.Request) {
		if request.Body == nil {
			return // invalid request
		}

		// decode the state packet
		state := new(GameState)
		decoder := json.NewDecoder(request.Body)
		err := decoder.Decode(state)
		if err != nil {
			log.Fatal(err)
		}

		callback(state)
	}

	// simple mux; don't like using global state on the http package
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)

	// start the http server
	addr := ":" + strconv.Itoa(port)
	err := http.ListenAndServe(addr, mux)
	if err != nil {
		log.Fatal(err)
	}
}
