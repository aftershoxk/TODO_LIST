package next_date

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const DateFormat = "20060102"

func afterNow(date, now time.Time) bool {
	return date.After(now)
}

func strToInt(arr []string) ([]int, error) {
	result := make([]int, 0, len(arr))
	for i := 0; i < len(arr); i++ {
		n, err := strconv.Atoi(arr[i])
		if err != nil {
			return nil, fmt.Errorf("cannot convert")
		}
		result = append(result, n)
	}
	return result, nil
}

func strToIntWeek(arr []string) ([]int, error) {
	result := make([]int, 0, len(arr))
	for i := 0; i < len(arr); i++ {
		n, err := strconv.Atoi(arr[i])
		if n == 7 {
			n = 0
		}
		if err != nil || n < 0 || n > 6 {
			return nil, fmt.Errorf("invalid weekday")
		}
		result = append(result, n)
	}
	return result, nil
}

func IsInArray(x int, arr []int) bool {
	for i := 0; i < len(arr); i++ {
		if x == arr[i] {
			return true
		}
	}
	return false
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("input the correct info")
	}
	date, err := time.Parse(DateFormat, dstart)
	if err != nil {
		return "", err
	}
	info := strings.Fields(repeat)

	switch info[0] {
	case "d":
		if len(info) != 2 {
			return "", fmt.Errorf("incorrect format")
		}
		interval, err := strconv.Atoi(info[1])
		if err != nil {
			return "", err
		}
		if interval < 1 || interval > 400 {
			return "", fmt.Errorf("incorrect format")
		}

		for {
			date = date.AddDate(0, 0, interval)
			if afterNow(date, now) {
				break
			}
		}

	case "y":
		if len(info) != 1 {
			return "", fmt.Errorf("incorrect format")
		}

		for {
			date = date.AddDate(1, 0, 0)
			if afterNow(date, now) {
				break
			}
		}
	case "w":
		if len(info) != 2 {
			return "", fmt.Errorf("incorrect format")
		}

		parsed := strings.Split(info[1], ",")
		allowed, err := strToIntWeek(parsed)
		if err != nil {
			return "", err
		}

		for {
			date = date.AddDate(0, 0, 1)
			if !afterNow(date, now) {
				continue
			}
			wd := int(date.Weekday())
			if IsInArray(wd, allowed) {
				break
			}

		}
	case "m":
		if len(info) < 2 || len(info) > 3 {
			return "", fmt.Errorf("incorrect format")
		}
		var d [32]bool
		var m [13]bool
		lastDayAllowed := false
		prevLastDayAllowed := false
		days := strings.Split(info[1], ",")
		daysAllowed, err := strToInt(days)
		if err != nil {
			return "", fmt.Errorf("incorrect format")
		}
		if len(info) == 3 {
			months := strings.Split(info[2], ",")
			monthsAllowed, err := strToInt(months)
			if err != nil {
				return "", fmt.Errorf("incorrect format")
			}
			for _, v := range monthsAllowed {
				if v >= 1 && v <= 12 {
					m[v] = true
				} else {
					return "", fmt.Errorf("incorrect format")
				}
			}
		}
		for _, v := range daysAllowed {
			switch {
			case v >= 1 && v <= 31:
				d[v] = true
			case v == -1:
				lastDayAllowed = true
			case v == -2:
				prevLastDayAllowed = true
			default:
				return "", fmt.Errorf("incorrect format")
			}
		}
		for {
			date = date.AddDate(0, 0, 1)
			var monthOK bool
			if !afterNow(date, now) {
				continue
			}
			if len(info) == 2 {
				monthOK = true
			} else if m[int(date.Month())] {
				monthOK = true
			}
			currentDay := date.Day()
			lastDayOfMonth := time.Date(
				date.Year(),
				date.Month()+1,
				0,
				0, 0, 0, 0,
				date.Location(),
			).Day()
			var dayOK bool
			if d[currentDay] {
				dayOK = true
			}
			if lastDayAllowed && currentDay == lastDayOfMonth {
				dayOK = true
			}
			if prevLastDayAllowed && currentDay == lastDayOfMonth-1 {
				dayOK = true
			}
			if !monthOK || !dayOK {
				continue
			}
			break
		}
	default:
		return "", fmt.Errorf("incorrect format")
	}
	return date.Format(DateFormat), nil
}
