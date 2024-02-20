package main

import (
	"flag"
	"fmt"
	"github.com/mkch/wol"
	probing "github.com/prometheus-community/pro-bing"
	"log"
	"math"
	"os"
	"regexp"
	"strings"
	"time"
)

func isValidAddress(mac string) bool {
	if len(mac) == 0 {
		return false
	}
	// Regular expression to match MAC address format
	macRegex := regexp.MustCompile(`^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$`)
	return macRegex.MatchString(mac)
}

func isHostUp(host string, pings int, timeOut int) bool {
	pinger, err := probing.NewPinger(host)
	if err != nil {
		panic(err)
	}

	pinger.Timeout = time.Duration(timeOut) * time.Second
	pinger.Count = pings
	if err = pinger.Run(); err != nil {
		// fail to execute the ping, happens if network interface is down or other issues
		log.Println("Failed to execute the ping: ", err)
		return false
	}

	// considering the network up if the packet loss is less than 50%
	s := pinger.Statistics()
	return s.PacketLoss < 50.0
}

// display app usage and then quit
func usage() {
	fmt.Println("Usage: [flags] mac ip/hostname\n  " +
		"-a int\n\tAttempts to wake the host (default 4)\n  -e int\n\t" +
		"Exponential wait between retries in seconds (default 5)\n  -t int\n\t" +
		"Ping Timeout in seconds (default 2)\n  -w int\n\t" +
		"Seconds to wait between sending the wake command and pinging (default 60)")
	os.Exit(0)
}

func main() {
	flag.Usage = usage

	wait := flag.Int("w", 60, "Seconds to wait between sending the wake command and pinging")
	attempts := flag.Float64("a", 4, "Attempts to wake the host")
	expWaitTime := flag.Float64("e", 5, "Exponential wait between retries in seconds")
	timeOut := flag.Int("t", 2, "Ping Timeout in seconds")
	flag.Parse()

	// trim possible spaces surrounding mac and ip
	macAddr := strings.TrimSpace(flag.Arg(0))
	ip := strings.TrimSpace(flag.Arg(1))
	// make a simple check on provided arguments
	if len(flag.Args()) != 2 || len(ip) == 0 || !isValidAddress(macAddr) {
		usage()
	}

	success := false
	isFirstRun := true
	for i := 0.0; i < *attempts; i++ {
		// avoids unnecessary waiting if it is the last iteration and unsuccessful
		if !isFirstRun {
			// wait a certain amount of seconds between executions, it grows exponentially (5, 10, 20...)
			s := time.Duration(*expWaitTime * math.Pow(2, i))
			log.Printf("The host is not up, I will try again to wake it up in %d seconds\n", s)
			time.Sleep(s * time.Second)
		}

		isFirstRun = false

		// try to wake up the host
		if err := wol.Wake(macAddr); err != nil {
			// sending packet generate an error, retrying without waiting
			log.Printf("Unable to send wol to %s\n%s\nretrying\n", macAddr, err)
			isFirstRun = true
			continue
		} else {
			// packet sent without errors
			log.Printf("Magic packet sent, waiting %d seconds before pinging\n", *wait)
		}

		// wait for the host finishing to boot
		time.Sleep(time.Duration(*wait) * time.Second)

		//try to ping the host, if more than half of the packets arrives the host is considered up
		log.Printf("Sendig 4 pings to %s\n", ip)
		if isHostUp(ip, 4, *timeOut) {
			success = true
			// exit from loop
			break
		}
	}

	if success {
		log.Printf("The host %s is up\n", ip)
	} else {
		log.Printf("Failed to wake up %s\n", macAddr)
	}
}
