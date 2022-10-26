package caddy

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/assistcontrol/muxytail/color"
	"github.com/mileusna/useragent"
)

//
// CADDY LOG DATA STRUCTURE
//

// struct caddyLog is holds a single caddy log entry.
type caddyLog struct {
	Req struct {
		Headers struct {
			Referer caddyReferer   `json:"Referer"`
			UA      caddyUserAgent `json:"User-Agent"`
		} `json:"headers"`
		Host   string        `json:"host"`
		Method string        `json:"method"`
		Proto  string        `json:"proto"`
		Remote caddyRemoteIP `json:"remote_ip"`
		URI    string        `json:"uri"`
	} `json:"request"`
	Status caddyStatus    `json:"status"`
	TS     caddyTimeStamp `json:"ts"`
}

// caddyLog.URL constructs a request path (foo.com/bar.html)
func (cL *caddyLog) URL() string {
	return blue(cL.Req.Host + cL.Req.URI)
}

//
// CADDY LOG FIELD TYPES
//

// ---
// caddyReferer is the client-supplied referer URL
type caddyReferer []string

// caddyReferer.String returns the longest referer.
func (referers caddyReferer) String() string {
	return longest(referers)
}

// ---
// caddyRemoteIP is the client IP.
type caddyRemoteIP string

// caddyRemoteIP.String does a reverse lookup on the IP
// if possible.
func (cR caddyRemoteIP) String() string {
	remote := string(cR)
	rev, err := net.LookupAddr(remote)
	if err != nil || len(rev) == 0 {
		return blue(remote)
	}

	// Trim trailing .
	remote = strings.TrimSuffix(rev[0], ".")
	return yellow(remote)
}

// ---
// caddyStatus is the HTTP status code (200).
type caddyStatus int

// caddyStatus.String colorizes the status code.
func (statusCode caddyStatus) String() string {
	status := strconv.Itoa(int(statusCode))

	if statusCode >= 200 && statusCode < 300 {
		return green(status)
	} else if statusCode >= 400 {
		return red(status)
	} else {
		return yellow(status)
	}
}

// ---
// caddyTimeStamp is the timestamp of the request
type caddyTimeStamp float64

// timeFormat is the format string for displayed date/times.
const timeFormat = "2/Jan 15:04:05"

// caddyTimeStamp.String formats the TS using const timeFormat.
func (cTS caddyTimeStamp) String() string {
	return time.Unix(int64(cTS), 0).Format(timeFormat)
}

// ---
// caddyUserAgent holds the client UA
type caddyUserAgent []string

// caddyUserAgent.String parses the longest UA in the form
// "browserName browserMajor | OSName OSMajor"
func (cUA caddyUserAgent) String() string {
	chosenUA := longest(cUA)

	if len(chosenUA) == 0 {
		return "-"
	}

	ua := useragent.Parse(chosenUA)

	var bot string
	if ua.Bot {
		bot = "BOT:"
	}

	s := fmt.Sprintln(bot, ua.Name, major(ua.Version), "│", ua.OS, major(ua.OSVersion))
	s = strings.TrimSpace(s)
	s = strings.TrimSuffix(s, " │") // If OS is empty
	return s
}

//
// HELPER FUNCTIONS
//

func blue(s string) string   { return color.ColorizeString("Blue", s) }
func green(s string) string  { return color.ColorizeString("Green", s) }
func yellow(s string) string { return color.ColorizeString("Yellow", s) }
func red(s string) string    { return color.ColorizeString("Red", s) }

// major returns the major only from a semver version string
func major(s string) string {
	major, _, _ := strings.Cut(s, ".")
	return major
}

// longest returns the longest string in a string slice
func longest(strs []string) string {
	var chosen string

	for _, s := range strs {
		if len(s) > len(chosen) {
			chosen = s
		}
	}

	return chosen
}
