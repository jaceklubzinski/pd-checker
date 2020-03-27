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
	layoutISO := "2006-01-02T15:04:05Z"
	converted, err := time.Parse(layoutISO, value)
	CheckErr(err)
	return converted
}

func LastTillNowDuration(start string) time.Duration {
	return time.Now().Sub(StringToDate(start))
}

func NewStartDate(start string, timer string) string {
	startDate := StringToDate(start)
	timerDuration, err := time.ParseDuration(timer)
	CheckErr(err)
	return startDate.Add(timerDuration).String()
}
