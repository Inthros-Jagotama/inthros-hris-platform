package numbering

import (
	"testing"
	"time"
)

func TestFormatTemplate(t *testing.T) {
	now := time.Date(2026, time.August, 14, 0, 0, 0, 0, time.UTC)

	cases := []struct {
		name     string
		template string
		sequence int
		want     string
	}{
		{
			name:     "sequence with padding, month roman, year",
			template: "SK/{sequence:3}/HRIS/{month_roman}/{year}",
			sequence: 7,
			want:     "SK/007/HRIS/VIII/2026",
		},
		{
			name:     "plain sequence and two digit year",
			template: "CTR-{sequence}-{yy}{month}",
			sequence: 42,
			want:     "CTR-42-2608",
		},
		{
			name:     "sequence padding wider than value",
			template: "{sequence:5}",
			sequence: 3,
			want:     "00003",
		},
		{
			name:     "unknown token left literal",
			template: "{prefix}-{sequence:2}",
			sequence: 1,
			want:     "{prefix}-01",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := FormatTemplate(tc.template, tc.sequence, now)
			if got != tc.want {
				t.Errorf("FormatTemplate(%q, %d) = %q, want %q", tc.template, tc.sequence, got, tc.want)
			}
		})
	}
}

func TestResetKeyFor(t *testing.T) {
	now := time.Date(2026, time.August, 14, 0, 0, 0, 0, time.UTC)

	cases := []struct {
		period string
		want   string
	}{
		{ResetPeriodYearly, "2026"},
		{ResetPeriodMonthly, "2026-08"},
		{ResetPeriodNever, ""},
	}

	for _, tc := range cases {
		t.Run(tc.period, func(t *testing.T) {
			got := ResetKeyFor(tc.period, now)
			if got != tc.want {
				t.Errorf("ResetKeyFor(%q) = %q, want %q", tc.period, got, tc.want)
			}
		})
	}
}
