package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestIsBusinessDay(t *testing.T) {
	tests := []struct {
		name string
		date string
		want bool
	}{
		{
			name: "should return false for a Saturday",
			date: "2020-09-05",
			want: false,
		},
		{
			name: "should return false for a Sunday",
			date: "2020-09-06",
			want: false,
		},
		{
			name: "should return false for a holiday",
			date: "2020-09-07",
			want: false,
		},
		{
			name: "should return true for a business day",
			date: "2020-09-08",
			want: true,
		},
		{
			name: "should return true for a business day",
			date: "2020-09-09",
			want: true,
		},
		{
			name: "should return true for a business day",
			date: "2020-09-10",
			want: true,
		},
		{
			name: "should return true for a business day",
			date: "2020-09-11",
			want: true,
		},
		{
			name: "should return true for a business day",
			date: "2020-09-12",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			date, _ := time.Parse("2006-01-02", tt.date)
			assert.Equal(t, tt.want, IsBusinessDay(date))
		})
	}
}

func TestNextBusinessDay(t *testing.T) {
	tests := []struct {
		name      string
		date      string
		want      string
		inclusive bool
	}{
		{
			name:      "should return the next business day",
			date:      "2020-09-05",
			want:      "2020-09-08",
			inclusive: false,
		},
		{
			name:      "should return the next business day",
			date:      "2020-09-06",
			want:      "2020-09-08",
			inclusive: false,
		},
		{
			name:      "should return the next business day",
			date:      "2020-09-07",
			want:      "2020-09-08",
			inclusive: false,
		},
		{
			name:      "should return the next business day",
			date:      "2020-09-08",
			want:      "2020-09-09",
			inclusive: false,
		},
		{
			name:      "should return the next business day",
			date:      "2020-09-09",
			want:      "2020-09-10",
			inclusive: false,
		},
		{
			name:      "should return the next business day",
			date:      "2020-09-10",
			want:      "2020-09-11",
			inclusive: false,
		},
		{
			name:      "should return the next business day",
			date:      "2020-09-11",
			want:      "2020-09-14",
			inclusive: false,
		},
		{
			name:      "should return the next business day",
			date:      "2020-09-12",
			want:      "2020-09-14",
			inclusive: false,
		},
		{
			name:      "should return the next business day",
			date:      "2020-09-13",
			want:      "2020-09-14",
			inclusive: false,
		},
		{
			name:      "should return the next business day",
			date:      "2020-09-14",
			want:      "2020-09-15",
			inclusive: false,
		},
		{
			name:      "should return the next business day",
			date:      "2020-09-15",
			want:      "2020-09-16",
			inclusive: false,
		},
		{
			name:      "should return the next business day",
			date:      "2020-09-16",
			want:      "2020-09-17",
			inclusive: false,
		},
		{
			name:      "should return the next business day",
			date:      "2023-12-24",
			want:      "2023-12-26",
			inclusive: false,
		},
		{
			name:      "should return the next business day",
			date:      "2023-12-29",
			want:      "2023-12-29",
			inclusive: true,
		},
		{
			name:      "should return the next business day",
			date:      "2023-12-29",
			want:      "2024-01-02",
			inclusive: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			date, _ := time.Parse("2006-01-02", tt.date)
			want, _ := time.Parse("2006-01-02", tt.want)
			assert.Equal(t, want, NextBusinessDay(date, tt.inclusive))
		})
	}
}
