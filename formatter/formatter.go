package formatter

// Formatter defines the interface for a formatter that takes
// a string and attempts to reformat/colorize it. Formatters
// return the formatted string and a boolean indicating whether
// formatting succeeded.
type Formatter interface {
	Format(string) (string, bool)
}

// List defines a slice of Formatters.
type List []Formatter
