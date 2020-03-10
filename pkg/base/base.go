package base

import "time"

func RepeatFunction(a func()) {
	ticker := time.NewTicker(60 * time.Second)
	for _ = range ticker.C {
		a()
	}
}
