/*------------------------------------------------------------------------------
-- DATE:	       February 25, 2016
--
-- Source File:	 main.go
--
-- REVISIONS: 	March 1, 2016
--									Moved many functions out of this file into their own files.
--
-- DESIGNER:	   Marc Vouve
--
-- PROGRAMMER:	 Marc Vouve
--
--
-- INTERFACE:
--	func main()
--
-- NOTES: This file is the main file in the IPS suite.
------------------------------------------------------------------------------*/
package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/Unknwon/goconfig"
	"github.com/abh/geoip"
)

var configure *goconfig.ConfigFile
var wg *sync.WaitGroup
var stdOutMutex *sync.Mutex

/*-----------------------------------------------------------------------------
-- FUNCTION:    main
--
-- DATE:        March 04, 2016
--
-- REVISIONS:	  (none)
--
-- DESIGNER:		Marc Vouve
--
-- PROGRAMMER:	Marc Vouve
--
-- INTERFACE:		func main()
--
-- RETURNS: 		void
--
-- NOTES:			The main entry point for the ips.
------------------------------------------------------------------------------*/
func main() {
	var err error
	defer mainRecovery()
	logChan := make(chan logManager)
	var totalWorkers int

	if err != nil {
		fmt.Println("Need to configure geoip file location")
	}

	for _, fileName := range os.Args[1:] {
		go worker(fileName, logChan)
		totalWorkers++
	}
	observer(totalWorkers, logChan)

}

func observer(n int, c chan logManager) {
	geoIP, err := geoip.Open(geoIPFile)
	if err != nil {
		log.Fatalln(err)
	}
	logtrack := initManifest()
	for i := 0; i < n; i++ {
		logtrack.merge(<-c)
		fmt.Println("Completed Files: ", i+1, "/", n)

	}

	for _, ip := range logtrack.sortedEventKeys() {
		c, _ := geoIP.GetCountry(ip)
		if len(logtrack.Events[ip].Success) > 0 {
			fmt.Printf("%-16v: Failed to connect times %-5v Connected: %-5v times from: %-2v\n", ip, len(logtrack.Events[ip].Failure), len(logtrack.Events[ip].Success), c)
		}
	}
}

func worker(fileName string, c chan<- logManager) {
	log := initManifest()
	log.checkSecure(fileName)

	c <- log
}

func mainRecovery() {
	if r := recover(); r != nil {
		log.Fatalln("Program not invoked correctly")
	}
}
