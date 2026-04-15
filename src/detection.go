package main

import (
	"fmt"
	"net"
	"regexp"
	"slices"
	"strconv"
	"time"
)

type BadIP struct {
	IP   string
	Port int
}

var (
	mainConfig  Config
	mainBlocked Blocked
	OS          string
	iplist      = make(map[string][]int)
	hotlist     = make(map[string][]time.Time)
	isRunning   = true

	rateLimitRe   = regexp.MustCompile(`([0-9]+\.[0-9]+\.[0-9]+\.[0-9]+):([0-9]+)\s+was\s+blocked\s+for\s+exceeding\s+rate\s+limits$`)
	splitPacketRe = regexp.MustCompile(`^([0-9]+\.[0-9]+\.[0-9]+\.[0-9]+):([0-9]+)\s+tried\s+to\s+send\s+split\s+packet`)
	badRconRe     = regexp.MustCompile(`(?i)Bad\sRcon:\s(?:.*)\sfrom\s\"([0-9]+\.[0-9]+\.[0-9]+\.[0-9]+):([0-9]+)\"`)
)

func regex(incoming string, re *regexp.Regexp) *BadIP {

	if matches := re.FindStringSubmatch(incoming); matches != nil {
		port, _ := strconv.Atoi(matches[2])
		return &BadIP{IP: matches[1], Port: port}
	}

	return nil
}

func processUDPData(data string) {
	var ret *BadIP
	// A2S abuse / Reflected DDoS
	if ret = regex(data, rateLimitRe); ret != nil && !slices.Contains(mainBlocked.IPs, ret.IP) {
		if existingPorts, exists := iplist[ret.IP]; exists {
			iplist[ret.IP] = append(existingPorts, ret.Port)
			counted := countPorts(iplist[ret.IP], ret.Port)

			if counted >= 3 {
				fmt.Printf("Found a naughty IP address: %s (rate limit, repeating port was: %d)\n", ret.IP, ret.Port)
				blockIP(ret.IP)
				delete(hotlist, ret.IP)
				delete(iplist, ret.IP)
			} else if counted <= 1 && len(iplist[ret.IP]) > 12 {
				fmt.Printf("Popping: %s (most likely valid person spamming browser refresh)\n", ret.IP)
				future := time.Now().UTC().Add(240 * time.Second)
				timestamp := future.Unix()

				if existingTimes, exists := hotlist[ret.IP]; exists {
					hotlist[ret.IP] = append(existingTimes, time.Unix(timestamp, 0))
					timenow := time.Now().UTC()
					hotlist[ret.IP] = slices.DeleteFunc(hotlist[ret.IP], func(t time.Time) bool {
						return t.Before(timenow)
					})

					if len(hotlist[ret.IP]) > 3 {
						fmt.Printf("Found a naughty IP address: %s (rate limit, port blasting)\n", ret.IP)
						delete(hotlist, ret.IP)
						blockIP(ret.IP)
					}
				} else {
					hotlist[ret.IP] = []time.Time{time.Unix(timestamp, 0)}
				}
				delete(iplist, ret.IP)
			}
		} else {
			iplist[ret.IP] = []int{ret.Port}
		}
	}
	// Split packet
	if ret = regex(data, splitPacketRe); ret != nil && !slices.Contains(mainBlocked.IPs, ret.IP) {
		fmt.Printf("Found extremely naughty IP address: %s (SPLIT PACKET!)\n", ret.IP)
		blockIP(ret.IP)
	}
	// Bad RCON attempt
	if ret = regex(data, badRconRe); ret != nil && !slices.Contains(mainBlocked.IPs, ret.IP) {
		fmt.Printf("Found semi-naughty IP address: %s (bad RCON)\n", ret.IP)
		blockIP(ret.IP)
	}
}

func logReceiver(host string, port int) {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP(host), Port: port})

	if err != nil {
		fmt.Printf("Error: Failed to start UDP server: %v\n", err)
		return
	}

	defer conn.Close()

	fmt.Printf("GSDDDOSS log receiver started, listening on UDP socket at \"%s:%d\".\n", host, port)
	// Buffer is 128KB in size
	buf := make([]byte, 1<<17)

	for isRunning {
		conn.SetReadDeadline(time.Now().Add(1 * time.Second))
		n, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				continue
			}
			fmt.Printf("Error reading UDP: %v\n", err)
			continue
		}

		if n > 4 {
			data := string(buf[4:n])
			processUDPData(data)
		}
	}
}
