package terminal

//
// TreePrefix
//

type TreePrefix []bool

func (self TreePrefix) Print(indent int, last bool) {
	PrintIndent(indent)

	for _, element := range self {
		if element {
			Print("  ")
		} else {
			Print("│ ")
		}
	}

	if last {
		Print("└─")
	} else {
		Print("├─")
	}
}
