package util

import (
	"errors"
	"regexp"
	"time"
)

func DateParse(date string) (t time.Time, err error) {
	t, err = time.Parse("2006-01-02", date)
	if err != nil {
		return
	}
	return
}

func DateFormat(t time.Time) string {
	return t.Format("2006-01-02")
}

func DateTimeFormat(t time.Time) string {
	return t.Format("2006-01-02 15:04")
}

func UnixMilliSec(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}

// TimeFormats default time formats will be parsed as
var TimeFormats = []string{"1/2/2006", "1/2/2006 15:4:5", "2006", "2006-1", "2006/1/2", "2006-1-2", "2006-1-2 15", "2006-1-2 15:4", "2006-1-2 15:4:5", "1-2", "15:4:5", "15:4", "15", "15:4:5 Jan 2, 2006 MST", "2006-01-02 15:04:05.999999999 -0700 MST", "2006-01-02T15:04:05-07:00"}
var hasTimeRegexp = regexp.MustCompile(`(\s+|^\s*)\d{1,2}((:\d{1,2})*|((:\d{1,2}){2}\.(\d{3}|\d{6}|\d{9})))\s*$`) // match 15:04:05, 15:04:05.000, 15:04:05.000000 15, 2017-01-01 15:04, etc
var onlyTimeRegexp = regexp.MustCompile(`^\s*\d{1,2}((:\d{1,2})*|((:\d{1,2}){2}\.(\d{3}|\d{6}|\d{9})))\s*$`)      // match 15:04:05, 15, 15:04:05.000, 15:04:05.000000, etc

// Parse parse string to time
func ParseTime(str string) (t time.Time, err error) {
	now := time.Now()
	var (
		setCurrentTime  bool
		parseTime       []int
		currentLocation = time.Now().Location()
		currentTime     = []int{now.Nanosecond(), now.Second(), now.Minute(), now.Hour(), now.Day(), int(now.Month()), now.Year()}
		onlyTimeInStr   = true
	)

	hasTimeInStr := hasTimeRegexp.MatchString(str) // match 15:04:05, 15
	onlyTimeInStr = hasTimeInStr && onlyTimeInStr && onlyTimeRegexp.MatchString(str)
	if t, err = parseWithFormat(str); err == nil {
		location := t.Location()
		if location == nil || location.String() == "UTC" {
			location = currentLocation
		}

		parseTime = []int{t.Nanosecond(), t.Second(), t.Minute(), t.Hour(), t.Day(), int(t.Month()), t.Year()}

		for i, v := range parseTime {
			// Don't reset hour, minute, second if current time str including time
			if hasTimeInStr && i <= 3 {
				continue
			}

			// If value is zero, replace it with current time
			if v == 0 {
				if setCurrentTime {
					parseTime[i] = currentTime[i]
				}
			} else {
				setCurrentTime = true
			}

			// if current time only includes time, should change day, month to current time
			if onlyTimeInStr {
				if i == 4 || i == 5 {
					parseTime[i] = currentTime[i]
					continue
				}
			}
		}

		t = time.Date(parseTime[6], time.Month(parseTime[5]), parseTime[4], parseTime[3], parseTime[2], parseTime[1], parseTime[0], location)
		currentTime = []int{t.Nanosecond(), t.Second(), t.Minute(), t.Hour(), t.Day(), int(t.Month()), t.Year()}
	}
	return
}

func parseWithFormat(str string) (t time.Time, err error) {
	for _, format := range TimeFormats {
		t, err = time.Parse(format, str)
		if err == nil {
			return
		}
	}
	err = errors.New("Can't parse string as time: " + str)
	return
}
