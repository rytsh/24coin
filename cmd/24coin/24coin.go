package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/rytsh/24coin/internal/api"
	"github.com/rytsh/24coin/internal/common"
	"github.com/rytsh/24coin/internal/trader"
)

const helpText = `24coin [OPTIONS]
24coin price viewer

Options:
  -P, --port [5000]
    Port number
  -H, --host [0.0.0.0]
    Host address
  -v, --version
    Show version number
  -h, --help
    Show help
`

func usage() {
	fmt.Println(helpText)
	os.Exit(0)
}

var flagVersion bool

func exit(code int) {
	os.Exit(code)
}

func flagParse() {
	flag.Usage = usage

	flag.BoolVar(&flagVersion, "v", false, "")
	flag.BoolVar(&flagVersion, "version", false, "")

	flag.StringVar(&common.Settings.UI.Host, "H", "0.0.0.0", "")
	flag.StringVar(&common.Settings.UI.Host, "host", "0.0.0.0", "")
	flag.IntVar(&common.Settings.UI.Port, "P", 5000, "")
	flag.IntVar(&common.Settings.UI.Port, "port", 5000, "")

	flag.Parse()

	if common.Settings.UI.Port > 65535 {
		exit(11)
	}

	if flagVersion {
		fmt.Println(common.Version)
		os.Exit(0)
	}

	// additional args
	// return flag.Args()
}

func main() {
	flagParse()
	exit := make(chan int)
	// common.SignalCheck(exit, api.Close)
	common.SignalCheck(exit, api.Close, trader.Close)

	// get realtime data
	trader.WebSocketConnect()

	// serve API
	api.Serve()

	os.Exit(<-exit)
}
