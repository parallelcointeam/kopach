package main

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/parallelcointeam/kopach/defs"
	"github.com/parallelcointeam/kopach/global"
	"github.com/parallelcointeam/pod/btcjson"
)

const (
	showHelpMessage = "Specify -h to show available options"
)

// commandUsage display the usage for a specific command.
func commandUsage(method string) {
	appName := filepath.Base(os.Args[0])
	appName = strings.TrimSuffix(appName, filepath.Ext(appName))
	usage, err := btcjson.MethodUsageText(method)
	if err != nil {
		// This should never happen since the method was already checked before calling this function, but be safe.
		fmt.Fprintln(os.Stderr, "Failed to obtain command usage:", err)
		return
	}
	fmt.Fprintln(os.Stderr, "\n", appName, Version(), "\nUsage:")
	fmt.Fprintf(os.Stderr, "  %s\n", usage)
}

// usage displays the general usage when the help flag is not displayed and and an invalid command was specified.  The commandUsage function is used instead when a valid command was specified.
func usage(errorMessage string) {
	appName := filepath.Base(os.Args[0])
	appName = strings.TrimSuffix(appName, filepath.Ext(appName))
	fmt.Fprintln(os.Stderr, errorMessage)
	fmt.Fprintln(os.Stderr, appName, Version())
	fmt.Fprintln(os.Stderr, "\n", appName, Version(), "\nUsage:")
	fmt.Fprintf(os.Stderr, "  %s [OPTIONS] \n\n",
		appName)
	fmt.Fprintln(os.Stderr, showHelpMessage)
}

func main() {
	cfg, _, err := loadConfig()
	if err != nil {
		os.Exit(1)
	}
	url := defs.ParseURL(cfg.URL)
	global.Endpoints = append(global.Endpoints, url)
	switch {
	case !cfg.OneOnly:
		p := url.Port + 1
		if p == 0 {
			p = 11048
			url.Port = p
		}
		for i := p; i < p+8; i++ {
			global.Endpoints = append(global.Endpoints, defs.URL{
				Username: url.Username,
				Password: url.Password,
				Protocol: url.Protocol,
				Address:  url.Address,
				Port:     i,
			})
		}

		fallthrough
	case len(cfg.OtherPorts) > 0:
		for i := range cfg.OtherPorts {
			global.Endpoints = append(global.Endpoints, defs.URL{
				Username: url.Username,
				Password: url.Password,
				Protocol: url.Protocol,
				Address:  url.Address,
				Port:     cfg.OtherPorts[i],
			})
		}
	}
	r := global.Endpoints
	var e []defs.URL
	for i := range r {
		if _, err := net.Dial("tcp", fmt.Sprintf("%s:%d", r[i].Address, r[i].Port)); err == nil {
			e = append(e, r[i])
		}
	}
	global.Endpoints = e
	for i := range global.Endpoints {
		fmt.Println(global.Endpoints[i].String())
	}
}
