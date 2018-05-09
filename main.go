package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/xyproto/term"
	"golang.org/x/net/http2"
)

const version_string = "http2check 0.6"

// Message with an optional additional string that will appear in paranthesis
func msg(o *term.TextOutput, subject, msg string, extra ...string) {
	if len(extra) == 0 {
		o.Println(fmt.Sprintf("%s%s%s %s", o.DarkGray("["), o.LightBlue(subject), o.DarkGray("]"), msg))
	} else {
		o.Println(fmt.Sprintf("%s%s%s %s (%s)", o.DarkGray("["), o.LightBlue(subject), o.DarkGray("]"), msg, extra[0]))
	}
}

// We have an IPv6 addr where the URL needs to be changed from https://something to [something]:443
func fixIPv6(url string) string {
	port := ""
	if strings.HasPrefix(url, "http://") && !strings.HasSuffix(url, ":80") {
		url = url[7:]
		port = ":80"
	}
	if strings.HasPrefix(url, "https://") && !strings.HasSuffix(url, ":443") {
		url = url[8:]
		port = ":443"
	}
	return "[" + url + "]" + port
}

func main() {
	o := term.NewTextOutput(true, true)

	// Silence the http2 logging
	devnull, err := os.OpenFile(os.DevNull, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
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
		fmt.Println("Check if a given webserver is using HTTP/2")
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

	// Create a new terminal output struct (for colored text)
	o = term.NewTextOutput(runtime.GOOS != "windows", !*quiet)

	// Check if the version flag was given
	if *version {
		o.Println(version_string)
		os.Exit(0)
	}

	// Retrieve the commandline arguments
	args := flag.Args()

	// The default URL
	url := "https://http2.golang.org"
	if len(args) > 0 {
		url = args[0]
	}
	ipaddr := net.ParseIP(url)
	if ipaddr.DefaultMask() == nil {
		// Not a valid IPv4 address
		// Check if it's likely to be IPv6.

		// TODO: Find a better way to detect this
		if strings.Contains(url, "::") {
			url = fixIPv6(url)
		}
	}
	if !strings.Contains(url, "://") {
		url = "https://" + url
	}

	/*
	 * Enumerate the interfaces and strip strings like "%eth0",
	 * because they are parsed incorrectly by Go, with errors like:
	 * parse [ff02::1%!e(MISSING)th0]:443: invalid URL escape "%!e(MISSING)t"
	 */
	interfaces, err := net.Interfaces()
	if err != nil {
		o.ErrExit(err.Error())
	}
	for _, iface := range interfaces {
		// TODO: Find the final % and check if it is followed by an iface, instead
		iName := "%" + iface.Name
		if strings.Contains(url, iName) {
			o.Println(o.DarkGray("ignoring \"" + iName + "\""))
			url = strings.Replace(url, iName, "", -1)
			break
		}
	}

	// Display the URL that is about be checked
	o.Println(o.DarkGray("GET") + " " + o.LightCyan(url))

	// GET over HTTP/2
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		if strings.HasSuffix(err.Error(), "hexadecimal escape in host") {
			url = fixIPv6(url)
		} else {
			o.ErrExit(err.Error())
		}
	}
	tlsconf := &tls.Config{InsecureSkipVerify: true}
	rt := &http2.Transport{TLSClientConfig: tlsconf}
	res, err := rt.RoundTrip(req)
	if err != nil {
		// Pick up typical problems with IPv6 addresses
		// TODO: Find an exact way to do this instead
		if strings.Contains(err.Error(), "too many colons") {
			url = fixIPv6(url)
			o.Println(o.LightYellow("IPv6") + " " + o.DarkGray(url))
			req, err = http.NewRequest("GET", url, nil)
			if err != nil {
				o.ErrExit(err.Error())
			}
			res, err = rt.RoundTrip(req)
		}
		if err != nil {
			// Better looking error messages
			errorMessage := strings.TrimSpace(err.Error())
			if errorMessage == "bad protocol:" {
				msg(o, "protocol", o.DarkRed("Not HTTP/2"))
			} else if errorMessage == "http2: unsupported scheme and no Fallback" {
				msg(o, "HTTP/2", o.DarkRed("Not supported"))
			} else if strings.HasPrefix(errorMessage, "dial tcp") && strings.HasSuffix(errorMessage, ": connection refused") {
				msg(o, "host", o.DarkRed("Down"), errorMessage)
			} else if strings.HasPrefix(errorMessage, "tls: oversized record received with length ") {
				msg(o, "protocol", o.DarkRed("No HTTPS support"), errorMessage)
			} else if strings.HasPrefix(errorMessage, "http2: unexpected ALPN protocol") {
				msg(o, "protocol", o.DarkRed("Not HTTP/2"))
			} else if strings.HasPrefix(errorMessage, "dial tcp: lookup") {
				msg(o, "host", o.DarkRed("Down"), "host not found")
			} else {
				o.ErrExit(errorMessage)
			}
			os.Exit(1)
		}
	}

	// The final output
	msg(o, "protocol", o.White(res.Proto))
	msg(o, "status", o.White(res.Status))
}
