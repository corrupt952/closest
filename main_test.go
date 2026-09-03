package main

import (
	"reflect"
	"testing"
)

func TestWithImplicitSearch(t *testing.T) {
	testCases := []struct {
		name string
		args []string
		want []string
	}{
		{
			name: "bare pattern gets search inserted",
			args: []string{"closest", "config.yaml"},
			want: []string{"closest", "search", "config.yaml"},
		},
		{
			name: "flag before pattern gets search inserted",
			args: []string{"closest", "-a", ".envrc"},
			want: []string{"closest", "search", "-a", ".envrc"},
		},
		{
			name: "explicit search is left alone",
			args: []string{"closest", "search", "config.yaml"},
			want: []string{"closest", "search", "config.yaml"},
		},
		{
			name: "version is left alone",
			args: []string{"closest", "version"},
			want: []string{"closest", "version"},
		},
		{
			name: "no args is left alone",
			args: []string{"closest"},
			want: []string{"closest"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := withImplicitSearch(tc.args)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("withImplicitSearch(%v) = %v, want %v", tc.args, got, tc.want)
			}
		})
	}
}
