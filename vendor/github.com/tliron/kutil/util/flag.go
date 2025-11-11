package util

import (
	flagpkg "flag"
	"os"
)

func SetFlagsFromEnvironment(prefix string) {
	flagpkg.VisitAll(func(flag *flagpkg.Flag) {
		if value, ok := os.LookupEnv(prefix + flag.Name); ok {
			flagpkg.Set(flag.Name, value)
		}
	})
}
