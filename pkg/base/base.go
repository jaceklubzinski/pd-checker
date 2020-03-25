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

func StringToDate(value string) time.Time {
	layoutISO := "2006-01-02 15:04:05"
	converted, err := time.Parse(layoutISO, value)
	CheckErr(err)
	return converted
}
