package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"
)

type Blocked struct {
	IPs []string `json:"ips"`
}

func blockedSave() error {
	file, err := os.Create("blocked.json")

	if err != nil {
		return err
	}

	defer file.Close()
	json.NewEncoder(file).Encode(mainBlocked)

	return nil
}

// Loads the blocked IPs
func blockedLoad() (*Blocked, error) {
	file, err := os.Open("blocked.json")

	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Warning: Blocked list file not found. (A default one will be created.)")
			defaultBlocked := Blocked{
				IPs: []string{},
			}
			if err := blockedSave(); err != nil {
				return nil, err
			}
			return &defaultBlocked, nil
		}

		return nil, err
	}

	defer file.Close()

	var data Blocked

	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return nil, err
	}

	return &data, nil
}

func blockIP(ip string) {
	if ip == "127.0.0.1" || ip == "::1" {
		fmt.Println("Warning: Not blocking localhost IP address.")
		return
	}

	if slices.Contains(mainBlocked.IPs, ip) {
		return
	}

	mainBlocked.IPs = append(mainBlocked.IPs, ip)

	if err := blockedSave(); err != nil {
		fmt.Printf("Error: Failed to save blocked list: %v\n", err)
		return
	}

	var err error
	ruleName := "GSDDDOSS blocked"

	if mainConfig.CommandAddBlock != "" {
		cmdStr := strings.Replace(mainConfig.CommandAddBlock, "{ip}", ip, 1)
		err = exec.Command("cmd", "/c", cmdStr).Run()
	} else if OS == "linux" || OS == "linux2" {
		err = exec.Command("iptables", "-A", "INPUT", "-s", ip, "-j", "DROP", "-m", "comment", "--comment", "GSDDDOSS blocked").Run()
	} else {
		if mainConfig.WindowsRuleGrouped {
			err = exec.Command("netsh",
				"advfirewall",
				"firewall",
				"set",
				"rule",
				"name="+ruleName,
				"dir=in",
				"new",
				"remoteip="+strings.Join(mainBlocked.IPs, ",")).Run()

			if err != nil {
				err = exec.Command("netsh",
					"advfirewall",
					"firewall",
					"add",
					"rule",
					"name="+ruleName,
					"dir=in",
					"action=block",
					"description=IP addresses caught by GSDDDOSS tool.",
					"enable=yes",
					"profile=any",
					"interface=any",
					"remoteip="+strings.Join(mainBlocked.IPs, ",")).Run()
			}

		} else {
			err = exec.Command("netsh",
				"advfirewall",
				"firewall",
				"add",
				"rule",
				"name="+ruleName,
				"dir=in",
				"action=block",
				"description=IP address caught by GSDDDOSS tool.",
				"enable=yes",
				"profile=any",
				"interface=any",
				"remoteip="+ip).Run()
		}
	}

	if err != nil {
		fmt.Printf("Error: Failed to block IP address \"%s\": %v\n", ip, err)
	} else {
		fmt.Printf("Blocked IP address \"%s\".\n", ip)
	}
}
