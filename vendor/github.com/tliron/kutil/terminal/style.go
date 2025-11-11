package terminal

import (
	"strings"
)

var StdoutStylist = NewStylist(false)
var StderrStylist = NewStylist(false)

//
// Stylist
//

type Stylist struct {
	Colorize bool
}

func NewStylist(colorize bool) *Stylist {
	return &Stylist{colorize}
}

// ([Colorizer] signature)
func (self *Stylist) Heading(name string) string {
	if self.Colorize {
		return ColorGreen(strings.ToUpper(name))
	} else {
		return strings.ToUpper(name)
	}
}

// ([Colorizer] signature)
func (self *Stylist) Path(name string) string {
	if self.Colorize {
		return ColorCyan(name)
	} else {
		return name
	}
}

// ([Colorizer] signature)
func (self *Stylist) Name(name string) string {
	if self.Colorize {
		return ColorBlue(name)
	} else {
		return name
	}
}

// ([Colorizer] signature)
func (self *Stylist) TypeName(name string) string {
	if self.Colorize {
		return ColorMagenta(name)
	} else {
		return name
	}
}

// ([Colorizer] signature)
func (self *Stylist) Value(name string) string {
	if self.Colorize {
		return ColorYellow(name)
	} else {
		return name
	}
}

// ([Colorizer] signature)
func (self *Stylist) Error(name string) string {
	if self.Colorize {
		return ColorRed(name)
	} else {
		return name
	}
}
