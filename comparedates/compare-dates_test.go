package comparedates

import (
	"testing"
	"time"
)

func TestInTimespan(t *testing.T) {
	cases := []struct {
		start time.Time
		end   time.Time
		check time.Time
		want  bool
	}{
		{time.Date(2010, time.January, 3, 0, 0, 0, 0, time.UTC), time.Date(2011, time.January, 3, 0, 0, 0, 0, time.UTC), time.Date(2010, time.March, 3, 0, 0, 0, 0, time.UTC), true},
		{time.Date(2010, time.January, 3, 0, 0, 0, 0, time.UTC), time.Date(2011, time.January, 3, 0, 0, 0, 0, time.UTC), time.Date(2005, time.March, 3, 0, 0, 0, 0, time.UTC), false},
		{time.Date(2010, time.January, 3, 0, 0, 0, 0, time.UTC), time.Date(2011, time.January, 3, 0, 0, 0, 0, time.UTC), time.Date(2015, time.March, 3, 0, 0, 0, 0, time.UTC), false},
		{time.Date(2010, time.January, 3, 10, 45, 30, 0, time.UTC), time.Date(2010, time.January, 3, 10, 45, 45, 0, time.UTC), time.Date(2010, time.January, 3, 10, 45, 40, 0, time.UTC), true},
		{time.Date(2010, time.January, 3, 10, 45, 30, 0, time.UTC), time.Date(2010, time.January, 3, 10, 45, 45, 0, time.UTC), time.Date(2010, time.January, 3, 10, 45, 46, 0, time.UTC), false},
		{time.Date(2010, time.January, 3, 10, 45, 30, 0, time.UTC), time.Date(2010, time.January, 3, 10, 45, 45, 0, time.UTC), time.Date(2010, time.January, 3, 10, 45, 29, 0, time.UTC), false},
	}
	for _, c := range cases {
		got := InTimespan(c.start, c.end, c.check)
		if got != c.want {
			t.Errorf("InTimespan(%v, %v, %v) == %v, want %v", c.start, c.end, c.check, got, c.want)
		}
	}
}

func TestParseTimeStr(t *testing.T) {
	cases := []struct {
		str  string
		want time.Time
	}{
		{"2010-01-01 10:34:50", time.Date(2010, time.January, 1, 10, 34, 50, 0, time.UTC)},
		{"2010-01-01", time.Date(2010, time.January, 1, 0, 0, 0, 0, time.UTC)},
		{"2010-03-11", time.Date(2010, time.March, 11, 0, 0, 0, 0, time.UTC)},
		{"2010-01-01 11:02:10", time.Date(2010, time.January, 1, 11, 02, 10, 0, time.UTC)},
	}
	for _, c := range cases {
		got := ParseTimeStr(c.str)
		if got != c.want {
			t.Errorf("ParseTime(%v) == %v, want %v", c.str, got, c.want)
		}
	}
}

func TestParseTimeNs(t *testing.T) {
	cases := []struct {
		str  string
		want time.Time
	}{
		{"1445444040096", time.Date(2015, time.October, 21, 18, 14, 00, 96, time.Local)},
	}
	for _, c := range cases {
		got := ParseTimeNs(c.str)
		if got != c.want {
			t.Errorf("ParseTime(%v) == %v, want %v", c.str, got, c.want)
		}
	}
}
