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

package dota2

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

// Starts a HTTP server listening on the given port for game status updates
// Calls the given callback with the current game state, synchronously
// TODO: use go routines and a channel instead
func StartListener(port int, callback func(state *GameState)) error {
	// create a handler for processing incoming requests and forwarding it to our callback
	handler := func(writer http.ResponseWriter, request *http.Request) {
		// decode the state packet and forward it to the callback
		state := new(GameState)
		decoder := json.NewDecoder(request.Body)
		err := decoder.Decode(state)
		if err != nil {
			log.Print(err)
		}
		callback(state)
	}

	// start the http server
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)
	addr := ":" + strconv.Itoa(port)
	return http.ListenAndServe(addr, mux)
}
