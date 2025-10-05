// Package parser
//
// This package provides functions to parse Quadlet files or whole directories.
package parser

// QuadletDirectory Represent a directory with Quadlets
type QuadletDirectory struct {
	Quadlets    map[string]Quadlet // Map about Quadlets in the directory
	DisabledQSR []string           // Globally disabled Quadlet Syntax Rules
}

// Quadlet Represent a Quadlet file including its dropins
type Quadlet struct {
	Name        string                       // Name of the Quadlet file
	References  []string                     // Which other Quadlet it reference (e.g.: Network=foo.network)
	DisabledQSR []string                     // Disabled Quadlet Syntax Rules
	Properties  map[string][]QuadletProperty // Properties for each sections in the file
	Dropins     []Dropin                     // Properties in drop in directories, list in precedence: latest is the effective
	Header      []string                     // The comment lines at the beginning of the file, ignored in dropins
}

// QuadletProperty Key-value pair in Quadlet file
type QuadletProperty struct {
	Property string // Name of the property
	Value    string // Value of the property
}

// Dropin directory to override settings from Quadlet files
type Dropin struct {
	Directory  string                       // Name of directory where dropin put
	FileName   string                       // Name of the file in dropin directory
	Properties map[string][]QuadletProperty // Content of the file
}
