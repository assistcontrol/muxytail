// Convert from caddy's structured JSON to combined log format
package caddy

import (
	"encoding/json"
	"fmt"
	"regexp"
)

// Convert is the public interface to convert its string argument
// into a formatted and colorized string. If the line is not a
// caddy JSON log line, it is returned unmodified. The formatted
// string is returned along with a boolean flag for whether it's
// been converted.
func Convert(in string) (string, bool) {
	// Try to unmarshal the caddy JSON log entry into a
	// caddyLog struct
	var cl *caddyLog
	if err := json.Unmarshal([]byte(in), &cl); err != nil {
		// Not a valid log file
		return in, false
	}

	return format(cl), true
}

// regexp whiteSpaceRE is a regexp matching spaces so that they
// can be collapsed later.
var whiteSpaceRE = regexp.MustCompile(`\s+`)

// format returns the formatted log. It relies on caddyLog struct
// fields doing their magic in String() methods.
func format(cl *caddyLog) string {
	bracketL := colorize.Bracket("[")
	bracketR := colorize.Bracket("]")

	msg := fmt.Sprintf("%s %s%s%s %s (%s) %s %s %s %s%s%s",
		cl.Req.Remote,
		bracketL, cl.TS, bracketR,
		cl.URL(),
		cl.Status,
		cl.Req.Method,
		cl.Req.Proto,
		cl.Req.Headers.Referer,
		bracketL, cl.Req.Headers.UA, bracketR,
	)

	return whiteSpaceRE.ReplaceAllString(msg, " ") // Collapse whitespace
}
