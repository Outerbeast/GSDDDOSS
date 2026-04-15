package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

func run() error {
	OS = runtime.GOOS

	if !isAdmin() {
		fmt.Println("Error: This script must be run with administrative / super user privileges.")
		fmt.Println("(This is necessary to block IP addresses with your system's firewall.)")

		return os.ErrPermission
	}
	// cmdline args
	var (
		cliHost        string
		cliPort        int
		cliGrouped     bool
		cliFirewallCmd string
	)

	flag.StringVar(&cliHost, "host", "", "Listener address")
	flag.IntVar(&cliPort, "port", 0, "Listener port")
	flag.BoolVar(&cliGrouped, "grouped-rules", false, "Windows: group IPs into single firewall rule")
	flag.StringVar(&cliFirewallCmd, "firewall-cmd", "", "Custom firewall command (use {ip} for IP address)")
	flag.Parse()

	loadedConfig, err := configLoad()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	mainConfig = *loadedConfig

	loadedBlocked, err := blockedLoad()

	if err != nil {
		return fmt.Errorf("failed to load blocked list: %w", err)
	}

	mainBlocked = *loadedBlocked

	if cliHost != "" {
		mainConfig.ListenerAddr = cliHost
	}

	if cliPort != 0 {
		mainConfig.ListenerPort = cliPort
	}

	if cliGrouped {
		mainConfig.WindowsRuleGrouped = true
	}

	if cliFirewallCmd != "" {
		mainConfig.CommandAddBlock = cliFirewallCmd
	}
	// !-TODO-!: Holdover from Python, technically not needed? - determine if this stays for future purpose
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\nProgram exiting gracefully.")
		isRunning = false
	}()

	logReceiver(mainConfig.ListenerAddr, mainConfig.ListenerPort)

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
