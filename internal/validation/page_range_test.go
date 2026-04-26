package validation

import (
	"reflect"
	"testing"
)

func TestParsePageRangeValid(t *testing.T) {
	tests := []struct {
		name string
		spec string
		want []int
	}{
		{name: "single page", spec: "1", want: []int{1}},
		{name: "simple range", spec: "1-5", want: []int{1, 2, 3, 4, 5}},
		{name: "list", spec: "1,3,5", want: []int{1, 3, 5}},
		{name: "mixed", spec: "1-3,7,10-12", want: []int{1, 2, 3, 7, 10, 11, 12}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParsePageRange(tt.spec)
			if err != nil {
				t.Fatalf("ParsePageRange() error = %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("ParsePageRange() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParsePageRangeInvalid(t *testing.T) {
	tests := []string{
		"0",
		"-1",
		"5-1",
		"abc",
		"1,,2",
		"1-",
		"-5",
	}

	for _, spec := range tests {
		t.Run(spec, func(t *testing.T) {
			if _, err := ParsePageRange(spec); err == nil {
				t.Fatalf("expected error for %q", spec)
			}
		})
	}
}
