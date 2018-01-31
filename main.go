// Copyright 2017, the project authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package main

import (
	"bitbucket.org/mattklein/dota2gsi/dota2"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

var portFlag = flag.Int("port", 4000, "The port to listen on for game state information")

func main() {
	parseFlags()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	address := "localhost:" + strconv.Itoa(*portFlag)
	updates := dota2.ListenForUpdates(address)

	running := true
	for running {
		select {
		case update := <-updates:
			fmt.Println(update)

		case <-signals:
			running = false
		}
	}
}

func parseFlags() {
	flag.Parse()

	if *portFlag == 0 {
		flag.Usage()
		log.Fatal("A valid port was expecte")
	}
}
