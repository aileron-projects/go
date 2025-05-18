package zcron

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Crontab is a cron scheduler.
type Crontab struct {
	second  uint64
	minute  uint64
	hour    uint64
	day     uint64
	month   uint64
	week    uint64
	timeNow func() time.Time
}

func (c *Crontab) valid() bool {
	now := c.timeNow()
	var day, week int
	var month time.Month
	for i := 0; i < 10*366; i++ { // Check 10 years.
		now = now.Add(24 * time.Hour)
		month = now.Month()
		day = now.Day()
		week = int(now.Weekday())
		if c.day&(1<<day) > 0 && c.week&(1<<week) > 0 && c.month&(1<<month) > 0 {
			return true
		}
	}
	return false
}

func (c *Crontab) WithTimeFunc(timeNow func() time.Time) {
	if timeNow == nil {
		return
	}
	c.timeNow = timeNow
}

// Now returns the current time returned from the internal clock.
// Use [Crontab.WithTimeFunc] to replace the internal clock.
func (c *Crontab) Now() time.Time {
	return c.timeNow()
}

// Next returns the next cron scheduled time.
// It internally calls [Crontab.NextAfter] with the current time
// returned from the internal clock.
// Use [Crontab.WithTimeFunc] when replacing the internal clock.
func (c *Crontab) Next() time.Time {
	return c.NextAfter(c.timeNow())
}

// NexAfter returns the next cron scheduled time after t.
func (c *Crontab) NextAfter(t time.Time) time.Time {
	now := t
	loc := now.Location()
	hour, min, sec := now.Clock()

	var year, day, week int
	var month time.Month

	for {
		year, month, day = now.Date()
		week = int(now.Weekday())

		if c.month&(1<<month) == 0 { // No schedule this month.
			month += 1 // Move to the next month.
			if month > 12 {
				month = 1
				year += 1
			}
			now = time.Date(year, month, 1, 0, 0, 0, 0, loc)
			hour, min, sec = -1, -1, -1
			continue
		}
		if c.day&(1<<day) == 0 || c.week&(1<<week) == 0 { // No schedule this day and day of week.
			now = now.Add(24 * time.Hour) // Move to the next day.
			hour, min, sec = -1, -1, -1
			continue
		}

		s := nextTime(c.second, sec, 59) // Check next scheduled second.

		m := min
		if min == -1 || !(s > sec) || c.minute&(1<<min) == 0 {
			m = nextTime(c.minute, min, 59)
		}

		h := hour
		if hour == -1 || !(m > min || (m == min && s > sec)) || c.hour&(1<<hour) == 0 {
			h = nextTime(c.hour, hour, 23)
			if !(h > hour) {
				now = now.Add(24 * time.Hour)
				hour, min, sec = -1, -1, -1
				continue
			}
		}

		return time.Date(year, month, day, h, m, s, 0, loc)
	}
}

// nextTime returns next schedule time.
// max should be
//   - 23 for hours
//   - 59 for minutes
//   - 59 for seconds
func nextTime(targets uint64, now int, max int) int {
	for i := now + 1; i <= max; i++ {
		if targets&(1<<i) > 0 {
			return i
		}
	}
	for i := 0; i <= now; i++ {
		if targets&(1<<i) > 0 {
			return i
		}
	}
	return 0
}

// ParseError reports cron parse error.
type ParseError struct {
	Err   error
	What  string
	Value string
}

func (e *ParseError) Unwrap() error {
	return e.Err
}

func (e *ParseError) Error() string {
	msg := "ztime/zcron: parse error. invalid " + e.What
	if e.Value != "" {
		msg += "(got:" + e.Value + ")"
	}
	if e.Err != nil {
		msg += "[" + e.Err.Error() + "]"
	}
	return msg
}

func (e *ParseError) Is(target error) bool {
	for target != nil {
		ee, ok := target.(*ParseError)
		if ok {
			return e.What == ee.What
		}
		target = errors.Unwrap(target)
	}
	return false
}

// Parse parses the given cron expression and returns Crontab.
// The syntax follows
//
//	TZ=UTC * * * * * *
//	|      | | | | | |
//	|      | | | | | |- Day of week
//	|      | | | | |--- Month
//	|      | | | |----- Day of month
//	|      | | |------- Hour
//	|      | |--------- Minute
//	|      |----------- Second (Optional)
//	|------------------ Timezone (Optional)
//
//	Field name   | Mandatory  | Values          | Special characters
//	----------   | ---------- | --------------  | -------------------
//	Timezone     | No         | Timezone name   |
//	Second       | No         | 0-59            | * / , -
//	Minute       | Yes        | 0-59            | * / , -
//	Hours        | Yes        | 0-23            | * / , -
//	Day of month | Yes        | 1-31            | * / , -
//	Month        | Yes        | 1-12 or JAN-DEC | * / , -
//	Day of week  | Yes        | 0-6 or SUN-SAT  | * / , -
//
// Note that the "Day of month" and the "Day of week" are
// evaluated with AND condition.
//
// Following aliases are defined for convenience.
//
//	Alias name   | Alias value | Usage example         |
//	----------   | ----------- | --------------------  |
//	CRON_TZ      | TZ          | CRON_TZ=UTC 0 0 * * * |
//	@monthly     | 0 0 1 * *   | TZ=UTC @monthly       |
//	@weekly      | 0 0 * * 0   | TZ=UTC @weekly        |
//	@daily       | 0 0 * * *   | TZ=UTC @daily         |
//	@hourly      | 0 * * * *   | TZ=UTC @hourly        |
//	@sunday      | 0 0 * * 0   | TZ=UTC @sunday        |
//	@monday      | 0 0 * * 1   | TZ=UTC @monday        |
//	@tuesday     | 0 0 * * 2   | TZ=UTC @tuesday       |
//	@wednesday   | 0 0 * * 3   | TZ=UTC @wednesday     |
//	@thursday    | 0 0 * * 4   | TZ=UTC @thursday      |
//	@friday      | 0 0 * * 5   | TZ=UTC @friday        |
//	@saturday    | 0 0 * * 6   | TZ=UTC @saturday      |
//
// In addition, "@every <Duration>" expression can be used.
// Format of the duration must follow the [time.ParseDuration].
//
// Example of "@every" expression.
//
//	Duration  | Resolved Cron        | Notes                 |
//	--------- | -------------------- | --------------------- |
//	-1s       | ERROR                | Duration must be >=0  |
//	0s        | ERROR                | Duration must be >=0  |
//	1s        | */1 * * * * *        |                       |
//	1m        | 0 */1 * * * *        |                       |
//	1h        | 0 0 */1 * * *        |                       |
//	61s       | */1 */1 * * *        |                       |
//	15m30s    | */30 */15 * * * *    |                       |
//	65m30s    | */30 */5 */1 * * *   |                       |
//	1h30m     | 0 */30 */1 * * *     |                       |
//	23h59m59s | */59 */59 */23 * * * |                       |
//	24h       | ERROR                | Duration must be <24h |
//
// See the references.
//   - https://en.wikipedia.org/wiki/Cron
//   - https://crontab.guru/
//   - https://crontab.cronhub.io/
func Parse(crontab string) (*Crontab, error) {
	crontab = strings.Trim(crontab, " \n\r\t\f,")
	fields := strings.Fields(replaceAlias(crontab))

	loc := time.Local // Default location.
	if strings.HasPrefix(fields[0], "TZ=") {
		parsedLoc, err := time.LoadLocation(strings.TrimPrefix(fields[0], "TZ="))
		if err != nil {
			return nil, &ParseError{Err: err, What: "timezone"}
		}
		loc = parsedLoc
		fields = fields[1:]
	}

	c := &Crontab{
		timeNow: func() time.Time { return time.Now().In(loc) },
	}
	switch len(fields) {
	case 5:
		// Valid number of fields.
		fields = append([]string{"0"}, fields...) // Add the "second" field.
	case 6:
		// Valid number of fields.
	default:
		return nil, &ParseError{What: "number of fields"}
	}
	var ok bool
	if c.second, ok = parseValue(fields[0], 0, 59); !ok {
		return nil, &ParseError{What: "second", Value: fields[0]}
	}
	if c.minute, ok = parseValue(fields[1], 0, 59); !ok {
		return nil, &ParseError{What: "minute", Value: fields[1]}
	}
	if c.hour, ok = parseValue(fields[2], 0, 23); !ok {
		return nil, &ParseError{What: "hour", Value: fields[2]}
	}
	if c.day, ok = parseValue(fields[3], 1, 31); !ok {
		return nil, &ParseError{What: "day of month", Value: fields[3]}
	}
	if c.month, ok = parseValue(normalizeMonth(fields[4]), 1, 12); !ok {
		return nil, &ParseError{What: "month", Value: fields[4]}
	}
	if c.week, ok = parseValue(normalizeWeek(fields[5]), 0, 6); !ok {
		return nil, &ParseError{What: "day of week", Value: fields[5]}
	}
	if !c.valid() {
		return nil, &ParseError{What: "scheduling (unschedulable)"}
	}
	return c, nil
}

// replaceAlias replaces aliases.
func replaceAlias(exp string) string {
	if strings.Count(exp, "@") > 1 {
		return exp // Invalid expression. Return as-is.
	}
	exp = strings.ReplaceAll(exp, "CRON_TZ=", "TZ=")
	val := strings.ToUpper(exp)
	if before, after, found := strings.Cut(val, "@EVERY"); found {
		d, err := time.ParseDuration(strings.TrimSpace(strings.ToLower(after)))
		if err != nil || d <= 0 || d >= 24*time.Hour {
			return exp // Invalid expression. Return as-is.
		}
		s, m, h := "*", "*", "*"
		if hours := int64(d / time.Hour); hours > 0 {
			d -= time.Duration(hours) * time.Hour
			h = "*/" + strconv.FormatInt(hours, 10)
			m = "0"
			s = "0"
		}
		if minutes := int64(d / time.Minute); minutes > 0 {
			d -= time.Duration(minutes) * time.Minute
			m = "*/" + strconv.FormatInt(minutes, 10)
			s = "0"
		}
		if seconds := int64(d / time.Second); seconds > 0 {
			s = "*/" + strconv.FormatInt(seconds, 10)
		}
		return before + fmt.Sprintf("%s %s %s * * *", s, m, h)
	}
	repl := map[string]string{
		"@MONTHLY":   "0 0 1 * *",
		"@WEEKLY":    "0 0 * * 0",
		"@DAILY":     "0 0 * * *",
		"@HOURLY":    "0 * * * *",
		"@SUNDAY":    "0 0 * * 0",
		"@MONDAY":    "0 0 * * 1",
		"@TUESDAY":   "0 0 * * 2",
		"@WEDNESDAY": "0 0 * * 3",
		"@THURSDAY":  "0 0 * * 4",
		"@FRIDAY":    "0 0 * * 5",
		"@SATURDAY":  "0 0 * * 6",
	}
	for k, v := range repl {
		if strings.Contains(val, k) {
			val = strings.ReplaceAll(val, k, v)
			return val
		}
	}
	return exp
}

// normalizeMonth returns normalized month expression.
// normalizeMonth returns "1"-"12" by replacing "JAN"-"DEC".
// It treats all alphabets case insensitive.
func normalizeMonth(exp string) string {
	if regexp.MustCompile(`[a-zA-Z]{4,}`).MatchString(exp) {
		return exp // Invalid expression. Return as-is.
	}
	val := strings.ToUpper(exp)
	repl := map[string]string{
		"JAN": "1",
		"FEB": "2",
		"MAR": "3",
		"APR": "4",
		"MAY": "5",
		"JUN": "6",
		"JUL": "7",
		"AUG": "8",
		"SEP": "9",
		"OCT": "10",
		"NOV": "11",
		"DEC": "12",
	}
	for k, v := range repl {
		val = strings.ReplaceAll(val, k, v)
	}
	if regexp.MustCompile(`[a-zA-Z]`).MatchString(val) {
		return exp // Invalid expression. Return as-is.
	}
	return val
}

// normalizeWeek returns normalized week expression.
// normalizeWeek returns "0"-"6" by replacing "SUN"-"SAT".
// It treats all alphabets case insensitive.
func normalizeWeek(exp string) string {
	if regexp.MustCompile(`[a-zA-Z]{4,}`).MatchString(exp) {
		return exp // Invalid expression. Return as-is.
	}
	val := strings.ToUpper(exp)
	repl := map[string]string{
		"SUN": "0",
		"MON": "1",
		"TUE": "2",
		"WED": "3",
		"THU": "4",
		"FRI": "5",
		"SAT": "6",
	}
	for k, v := range repl {
		val = strings.ReplaceAll(val, k, v)
	}
	if regexp.MustCompile(`[a-zA-Z]`).MatchString(val) {
		return exp // Invalid expression. Return as-is.
	}
	return val
}

// parseValue parses each fields of cron expression.
// Allowed expressions are listed below.
// Other expression results in error, or false at second returned value.
// It is always be failure if the given min and max is min>max,
// Both min and max MUST be zero or grater than zero.
//
//   - Wildcard           : "*"
//   - Number             : "5"
//   - Range              : "10-20"
//   - Wildcard with step : "*/3"
//   - Number with step   : "5/3"
//   - Range with step    : "10-20/3"
func parseValue(exp string, min, max int) (cron uint64, ok bool) {
	if min > max || min < 0 || max < 0 { // This is not allowed.
		return 0, false
	}

	result := uint64(0)
	for _, e := range strings.Split(exp, ",") {
		fields := strings.Split(e, "/")
		switch len(fields) {
		case 1: // Format 'wildcard', 'number' or 'range'.
			ini, end, ok := parseRange(fields[0], min, max)
			if !ok {
				return 0, false
			}
			result |= toBitArray(ini, end, 1)

		case 2: // Format 'wildcard with step', 'number with step' or 'range with step'.
			step, err := strconv.Atoi(fields[1]) // Parse step.
			if err != nil {
				return 0, false
			}
			if step <= 0 {
				return 0, false // Zero or negative step is not supported.
			}
			ini, end, ok := parseRange(fields[0], min, max)
			if !ok {
				return 0, false
			}
			if ini == end {
				end = max
			}
			result |= toBitArray(ini, end, step)

		default: // Unsupported format.
			return 0, false
		}
	}
	return result, true
}

// toBitArray returns bit array that represents
// value range. Returned value has flags between
// min and max by step.
func toBitArray(min, max, step int) uint64 {
	v := uint64(0)
	for i := min; i <= max; i += step {
		v |= 1 << i
	}
	return v
}

// parseRange returns the value range of the given expression.
// Allowed formats are wildcard, number and range.
// It is always be failure if the given min and max is min>max,
// Both min and max MUST be zero or grater than zero.
//
//   - Wildcard : "*"      returns min,max and true
//   - Number   : "5"      returns 5,5 and true
//   - Range    : "10-20"  returns 10,20 and true
//   - Others   :          returns 0,0 and false
func parseRange(exp string, min, max int) (int, int, bool) {
	if min > max || min < 0 || max < 0 { // This is not allowed.
		return 0, 0, false
	}

	if exp == "*" {
		return min, max, true
	}

	if !strings.Contains(exp, "-") {
		val, err := strconv.Atoi(exp) // exps is just a umber.
		if err != nil {
			return 0, 0, false
		}
		if val < min || val > max {
			return 0, 0, false
		}
		return val, val, true
	}

	fields := strings.Split(exp, "-")
	if len(fields) != 2 {
		return 0, 0, false
	}
	ini, err := strconv.Atoi(fields[0])
	if err != nil {
		return 0, 0, false
	}
	end, err := strconv.Atoi(fields[1])
	if err != nil {
		return 0, 0, false
	}
	if ini > end || ini < min || end > max {
		return 0, 0, false
	}

	return ini, end, true
}
