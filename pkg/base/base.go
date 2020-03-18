package base

import "time"

//RepeatFunction repeat given function every n seconds
func RepeatFunction(a func()) {
	ticker := time.NewTicker(60 * time.Second)
	for range ticker.C {
		a()
	}
}

//CheckErr return error
func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}
