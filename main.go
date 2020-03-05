package main

import (
	"fmt"
	"os"
)

func main() {
	integrationKey := os.Getenv("PAGERDUTY_INTEGRATION_KEY")

	fmt.Println(integrationKey)
}
