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

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/xeusalmighty/dota2gsi/dota2"
)

// Command line flags
var (
	portFlag = flag.Int("port", 4000, "The port to listen on for game state information")
)

// Entry point for the application
func main() {
	parseFlags()

	// listen for system signals
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// listen for game state updates
	states := dota2.ListenForUpdates(*portFlag)

	// run forever
	running := true
	for running {
		select {
		case state := <-states:
			// handle game state updates
			fmt.Println(state)

		case <-signals:
			// handle system signals
			running = false
		}
	}
}

// Parses command line flags/arguments
func parseFlags() {
	flag.Parse()

	if *portFlag == 0 {
		flag.Usage()
		log.Fatal("A valid port was expecte")
	}
}
