package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/bradfitz/http2"
	"github.com/xyproto/textgui"
)

const version_string = "http2check 0.1"

func msg(o *textgui.TextOutput, subject, msg string) {
	o.Println(fmt.Sprintf("%s%s%s %s", o.DarkGray("["), o.LightBlue(subject), o.DarkGray("]"), msg))
}

func main() {
	o := textgui.NewTextOutput(true, true)

	// Silence the http2 logging
	devnull, err := os.OpenFile("/dev/null", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		o.ErrExit("Could not open /dev/null for writing")
	}
	defer devnull.Close()
	log.SetOutput(devnull)

	// Flags

	version_help := "Show application name and version"
	quiet_help := "Don't write to standard out"

	version := flag.Bool("version", false, version_help)
	quiet := flag.Bool("q", false, quiet_help)

	flag.Usage = func() {
		fmt.Println()
		fmt.Println(version_string)
		fmt.Println("Check if webservers are using HTTP/2")
		fmt.Println()
		fmt.Println("Syntax: http2check [URI]")
		fmt.Println()
		fmt.Println("Possible flags:")
		fmt.Println("    --version                  " + version_help)
		fmt.Println("    --q                        " + quiet_help)
		fmt.Println("    --help                     This text")
		fmt.Println()
	}

	flag.Parse()

	// Use the flags and arguments

	o = textgui.NewTextOutput(true, !*quiet)

	args := flag.Args()

	if *version {
		o.Println(version_string)
		os.Exit(0)
	}

	// Default URL

	url := "https://http2.golang.org"
	if len(args) > 0 {
		url = args[0]
	}
	if !strings.Contains(url, "://") {
		url = "https://" + url
	}

	// Display the URL that is to be checked

	o.Println(o.DarkGray("GET") + " " + o.LightCyan(url))

	// GET over HTTP/2

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		o.ErrExit(err.Error())
	}
	rt := &http2.Transport{
		InsecureTLSDial: true,
	}
	res, err := rt.RoundTrip(req)
	if err != nil {
		errorMessage := strings.TrimSpace(err.Error())
		if errorMessage == "bad protocol:" {
			msg(o, "protocol", o.DarkRed("Not HTTP/2"))
		} else if errorMessage == "http2: unsupported scheme and no Fallback" {
			msg(o, "scheme", o.DarkRed("Unsupported, without fallback"))
		} else {
			o.ErrExit(errorMessage)
		}
		os.Exit(1)
	}

	// Final output

	msg(o, "protocol", o.White(res.Proto))
	msg(o, "status", o.White(res.Status))
}
