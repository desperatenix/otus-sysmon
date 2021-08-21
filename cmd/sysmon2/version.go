package main

import (
"flag"
"fmt"
)

var (
	buildDate = "UNKNOWN"
	gitHash   = "UNKNOWN"
	osSys = "UNKNOWN"
	arch = "UNKNOWN"
)

func printVersion() {
	fmt.Printf("System monitor built for %s arch %s on %s git %s\n", osSys, arch, buildDate, gitHash)
}

func isVersionCommand() bool {
	for _, name := range flag.Args() {
		if name == "version" {
			return true
		}
	}
	return false
}
