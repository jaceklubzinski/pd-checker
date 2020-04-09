package main

import (
	_ "github.com/jaceklubzinski/pd-checker/pkg/logformat"

	"github.com/jaceklubzinski/pd-checker/cmd"
)

// TODO:
// - retry
// - prometheus metrics
// - testy

func main() {
	cmd.Execute()
}
