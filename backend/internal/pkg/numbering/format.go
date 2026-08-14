package numbering

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var romanMonths = [...]string{"I", "II", "III", "IV", "V", "VI", "VII", "VIII", "IX", "X", "XI", "XII"}

var sequenceTokenRe = regexp.MustCompile(`\{sequence(?::(\d+))?\}`)

// FormatTemplate substitutes numbering tokens in template with concrete
// values for the given sequence number and point in time. Unknown tokens
// are left as literal text.
func FormatTemplate(template string, sequence int, now time.Time) string {
	result := sequenceTokenRe.ReplaceAllStringFunc(template, func(match string) string {
		sub := sequenceTokenRe.FindStringSubmatch(match)
		if sub[1] == "" {
			return strconv.Itoa(sequence)
		}
		width, err := strconv.Atoi(sub[1])
		if err != nil {
			return strconv.Itoa(sequence)
		}
		return fmt.Sprintf("%0*d", width, sequence)
	})

	replacer := strings.NewReplacer(
		"{year}", strconv.Itoa(now.Year()),
		"{yy}", fmt.Sprintf("%02d", now.Year()%100),
		"{month}", fmt.Sprintf("%02d", int(now.Month())),
		"{month_roman}", romanMonths[int(now.Month())-1],
	)
	return replacer.Replace(result)
}

// ResetKeyFor returns the period key used to detect when a sequence should
// reset: the year for "yearly", "YYYY-MM" for "monthly", and a constant
// empty string for "never" (which therefore never resets).
func ResetKeyFor(period string, now time.Time) string {
	switch period {
	case ResetPeriodYearly:
		return strconv.Itoa(now.Year())
	case ResetPeriodMonthly:
		return fmt.Sprintf("%04d-%02d", now.Year(), int(now.Month()))
	default:
		return ""
	}
}
