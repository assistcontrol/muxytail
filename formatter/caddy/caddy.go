// Convert from caddy's structured JSON to combined log format
package caddy

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"

	"github.com/assistcontrol/muxytail/color"
	"github.com/assistcontrol/muxytail/config"
)

// whiteSpaceRE is used to collapse consecutive whitespace into
// a single space
var whiteSpaceRE = regexp.MustCompile(`\s+`)

// struct colorizer holds the colorizer functions generated from
// the passed caddy config
type colorizer struct {
	Bracket                            color.Colorizer
	Host                               color.Colorizer
	StatusOK, StatusError, StatusOther color.Colorizer
	URL                                color.Colorizer
}

// Format parses a JSON-formatted caddy log entry, formats it
// into a legible style, and colorizes it. It returns the formatted
// line and a boolean indicating whether formatting was successful.
// If bool is false, the input was not a well-formed caddy JSON log
// entry.
func (clr *colorizer) Format(in string) (string, bool) {
	cLog, err := parse(in)
	if err != nil {
		return in, false
	}

	out := clr.formatLog(cLog)
	out = whiteSpaceRE.ReplaceAllString(out, " ") // collapse whitespace

	return out, true
}

// formatLog takes a parsed caddyLog struct, formats it, and
// colorizes it.
func (clr *colorizer) formatLog(cLog *caddyLog) string {
	bracketL := clr.Bracket("[")
	bracketR := clr.Bracket("]")

	s := fmt.Sprintf("%s %s%s%s %s (%s) %s %s %s %s%s%s",
		clr.Host(cLog.Req.Remote),
		bracketL, cLog.TS, bracketR,
		clr.URL(cLog.URL()),
		clr.colorizeStatus(cLog.Status),
		cLog.Req.Method,
		cLog.Req.Proto,
		cLog.Req.Headers.Referer,
		bracketL, cLog.Req.Headers.UA, bracketR,
	)

	return s
}

// colorizeStatus colorizes the supplied HTTP status code
// based on its numeric range.
func (clr *colorizer) colorizeStatus(statusCode caddyStatus) string {
	status := strconv.Itoa(int(statusCode))

	if statusCode >= 200 && statusCode < 300 {
		return clr.StatusOK(status)
	} else if statusCode >= 400 {
		return clr.StatusError(status)
	} else {
		return clr.StatusOther(status)
	}
}

// New takes a config.CaddyConfig struct specifying color strings,
// and returns a *colorizer struct of colorization functions that
// is capable of parsing and formatting caddy JSON log entries.
func New(conf config.CaddyConfig) *colorizer {
	c := &colorizer{
		Bracket:     color.GenerateColorizer(conf.Bracket),
		Host:        color.GenerateColorizer(conf.Host),
		StatusOK:    color.GenerateColorizer(conf.StatusOK),
		StatusError: color.GenerateColorizer(conf.StatusError),
		StatusOther: color.GenerateColorizer(conf.StatusOther),
		URL:         color.GenerateColorizer(conf.URL),
	}

	return c
}

// parse attempts to unmarshal a log entry into a caddyLog struct.
// It returns a pointer to that new struct and an error indicating
// whether parsing was successful.
func parse(in string) (*caddyLog, error) {
	// If unmarshalling fails, input wasn't a caddy JSON log entry
	var cl *caddyLog
	err := json.Unmarshal([]byte(in), &cl)

	return cl, err
}
