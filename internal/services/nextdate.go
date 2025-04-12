// nextdate - describes function 'NextDate' - algorithm for finding  next date for task
package services

import (
	"errors"
	"strconv"
	"time"
	"unicode"

	"github.com/Ekvo/yandex-practicum-go-final-project/internal/model"
	"github.com/Ekvo/yandex-practicum-go-final-project/pkg/common"
)

// ErrServicesInvalidDate - if bad result after check the line contain time.Time
var ErrServicesInvalidDate = errors.New("invalid date")

// ErrServicesWrongRepeat - invalid format !_repeat_! string
// func nextDateByDay(now, taskDateStart time.Time, !_repeat_! string) (time.Time, error)
var ErrServicesWrongRepeat = errors.New("invalid repeat data")

// ErrServicUnexpectedBehavior - critical error of algorithm 'NextDate'
var ErrServicUnexpectedBehavior = errors.New("unexpected behavior")

// flags - for !_repeat_! string
const (
	day  = 'd' // [d number], example: d 56
	weak = 'w' // [w number(s)], example: w 1,2

	// [m number(s)] or [m number(s) number(s)], example: m 1,2 5,7
	//  m day(s)     or  m day(s)      month(s)
	month = 'm'
	year  = 'y' //[y], example: y
)

const (
	minDay = 1

	// limitation:
	// 1 - number of days,
	// 2 - max number of loop in func(s) - 'nextDateByWeek','nextDateByMonth' look down
	maxDay = 400

	monday = 1
	sunday = 7
)

// NextDate - main function to algorithm
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if len(repeat) == 0 {
		return "", ErrServicesWrongRepeat
	}
	taskDateStart, err := time.Parse(model.DateFormat, dstart)
	if err != nil {
		return "", ErrServicesInvalidDate
	}
	now = common.ReduceTimeToDay(now)
	var newDate time.Time
	switch ch := repeat[0]; ch {
	case day:
		newDate, err = nextDateByDay(now, taskDateStart, repeat)
	case year:
		newDate, err = nextDateByYear(now, taskDateStart, repeat)
	case weak:
		newDate, err = nextDateByWeek(now, taskDateStart, repeat)
	case month:
		newDate, err = nextDateByMonth(now, taskDateStart, repeat)
	default:
		err = ErrServicesWrongRepeat
	}
	if err != nil {
		return "", err
	}
	return newDate.UTC().Format(model.DateFormat), nil
}

const (
	daySeconds = 24 * 60 * 60
)

// nextDateByDay - all time to UNIX and work only with type int64
// for starters  - find count of days for adding
func nextDateByDay(now, taskDateStart time.Time, repeat string) (time.Time, error) {
	days, err := numberOfDays(repeat)
	if err != nil || days < minDay || days > maxDay {
		return time.Time{}, ErrServicesWrongRepeat
	}
	nowUnix := now.UTC().Unix()
	startUnix := taskDateStart.UTC().Unix()
	daysToSeconds := int64(days * daySeconds)
	newDate := time.Time{}
	if startUnix < nowUnix {
		segments := (nowUnix - startUnix) / daysToSeconds
		newDate = time.Unix(startUnix+(segments+1)*daysToSeconds, 0).UTC()
	} else {
		newDate = time.Unix(startUnix+daysToSeconds, 0).UTC()
	}
	return newDate, nil
}

// numberOfDays - find number of days
// rules: "number of days" is solo value
//
//	"number of days" always exists and is less than or equal to 400 and greater than 0
func numberOfDays(repeat string) (int, error) {
	n := len(repeat)
	if n < 3 || repeat[1] != ' ' {
		return 0, ErrServicesWrongRepeat
	}
	// example: d 45656 -> d - [0], space [1]
	days, err := strconv.Atoi(repeat[2:])
	if err != nil {
		return 0, ErrServicesWrongRepeat
	}
	return days, nil
}

// nextDateByYear - add solo or many year(s) to taskDateStart
func nextDateByYear(now, taskDateStart time.Time, repeat string) (time.Time, error) {
	if len(repeat) != 1 {
		return time.Time{}, ErrServicesWrongRepeat
	}
	newDate := taskDateStart.AddDate(1, 0, 0).UTC()
	nowYear := now.Year()
	startYear := taskDateStart.Year()
	if startYear < nowYear {
		years := nowYear - startYear
		newDate = taskDateStart.AddDate(years, 0, 0).UTC()
	}
	// if the month or day or (month and day) is less than "now"
	if !newDate.UTC().After(now.UTC()) {
		newDate = newDate.AddDate(1, 0, 0).UTC()
	}
	return newDate, nil
}

// nextDateByWeek - find new date by day(s) of week
// for starters - find number of days
func nextDateByWeek(now, taskDateStart time.Time, repeat string) (time.Time, error) {
	days, err := daysOfWeek(repeat)
	if err != nil {
		return time.Time{}, err
	}
	newDate := taskDateStart
	if newDate.Before(now) {
		newDate = now
	}
	for i := 0; i < maxDay; i++ {
		newDate = newDate.AddDate(0, 0, 1).UTC()
		if days[newDate.Weekday()] {
			return newDate, nil
		}
	}
	return time.Time{}, ErrServicUnexpectedBehavior
}

// daysOfWeek - find unique day of week
// not repeat days
func daysOfWeek(repeat string) ([]bool, error) {
	n := len(repeat)
	if n < 3 || repeat[1] != ' ' {
		return nil, ErrServicesWrongRepeat
	}
	days := make([]bool, 8)
	for i := 2; i < n; i++ {
		ch := rune(repeat[i])
		if ch == ',' {
			if prev := repeat[i-1]; prev == ',' || prev == ' ' || i == n-1 {
				return nil, ErrServicesWrongRepeat
			}
			continue
		}
		if unicode.IsDigit(ch) {
			start := i
			for i++; i < n && unicode.IsDigit(rune(repeat[i])); i++ {
			}
			day, err := strconv.Atoi(repeat[start:i])
			if err != nil || day < monday || day > sunday || days[day] {
				return nil, ErrServicesWrongRepeat
			}
			days[day] = true
			i--
			continue
		}
		return nil, ErrServicesWrongRepeat
	}
	if days[sunday] {
		days[0] = true // time.Sunday = 0
	}
	return days, nil
}

const (
	maxDaysPerMonth     = 31
	lastDayMonth        = -1
	penultimateDayMonth = -2
)

// nextDateByMonth - find by (day of month) or (by month and day)
// for starters - get []any from  monthAndDay look down
//
// find by days           -> add day and check map on exist
//
// find by days and month -> start: do a day first in month and check month while [index]bool != true
// then -> find by days
func nextDateByMonth(now, taskDateStart time.Time, repeat string) (time.Time, error) {
	data, err := monthAndDay(repeat)
	if err != nil {
		return time.Time{}, err
	}
	days := data[0].(map[int]bool)
	var month []bool
	if len(data) == 2 {
		month = data[1].([]bool)
	}
	newDate := taskDateStart
	if !newDate.After(now) {
		newDate = now.AddDate(0, 0, 1).UTC()
	}
	for m, i := len(month), 0; i < maxDay; i++ {
		if m != 0 && !month[newDate.Month()] {
			newDate = newDate.AddDate(0, 1, 0).UTC()
			if newDate.Day() != 1 {
				newDate = common.BeginningOfMonth(newDate)
			}
			continue
		}
		if days[newDate.Day()] {
			return newDate, nil
		}
		if days[lastDayMonth] &&
			newDate.AddDate(0, 0, 1).UTC().Day() == 1 {
			return newDate, nil
		}
		if days[penultimateDayMonth] &&
			newDate.AddDate(0, 0, 2).UTC().Day() == 1 {
			return newDate, nil
		}
		newDate = newDate.AddDate(0, 0, 1).UTC()
	}
	return time.Time{}, ErrServicUnexpectedBehavior
}

// monthAndDay - find (number of days) map[int]bool -> (because may contain -1, -2)
//
//	or (number of days, number of month) -> ([]bool,map[int]bool)
//
// solo space  - create -> number of days
// twice space - create -> number of month
func monthAndDay(repeat string) ([]any, error) {
	n := len(repeat)
	if n < 3 || repeat[1] != ' ' {
		return nil, ErrServicesWrongRepeat
	}
	space := 1
	days := make(map[int]bool)
	months := make([]bool, time.December+1)
	for i := 2; i < n; i++ {
		ch := rune(repeat[i])
		if ch == ',' {
			if prev := repeat[i-1]; prev == ',' || prev == ' ' || i == n-1 {
				return nil, ErrServicesWrongRepeat
			}
			continue
		}
		if unicode.IsDigit(ch) || ch == '-' {
			if ch == '-' && unicode.IsDigit(rune(repeat[i-1])) { // m -1-1,31
				return nil, ErrServicesWrongRepeat
			}
			start := i
			for i++; i < n && unicode.IsDigit(rune(repeat[i])); i++ {
			}
			dayOrMonth, err := strconv.Atoi(repeat[start:i])
			if err != nil {
				return nil, ErrServicesWrongRepeat
			}
			if space == 1 {
				if d := common.Abs(dayOrMonth); d < minDay || d > maxDaysPerMonth ||
					dayOrMonth < penultimateDayMonth || days[dayOrMonth] {
					return nil, ErrServicesWrongRepeat
				}
				days[dayOrMonth] = true
			} else if space == 2 {
				month := time.Month(dayOrMonth)
				if month < time.January || month > time.December || months[month] {
					return nil, ErrServicesWrongRepeat
				}
				months[month] = true
			}
			i--
			continue
		}
		if ch == ' ' {
			space++
			if prev := repeat[i-1]; prev == ' ' || prev == ',' || space > 2 || i == n-1 {
				return nil, ErrServicesWrongRepeat
			}
			continue
		}
		return nil, ErrServicesWrongRepeat
	}
	if space == 1 {
		return []any{days}, nil
	}
	return []any{days, months}, nil
}
