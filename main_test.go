package main

import (
	"testing"
)

func Test_reduce(t *testing.T) {
	type args struct {
		dx int
		dy int
	}
	tests := []struct {
		name  string
		args  args
		want  int
		want1 int
	}{
		{"test1", args{6, 2}, 3, 1},
		{"test2", args{6, 4}, 3, 2},
		{"test3", args{22, 6}, 11, 3},
		{"test4", args{30, 10}, 3, 1},
		{"test5", args{9, 6}, 3, 2},
		{"test6", args{21, 18}, 7, 6},
		{"test7", args{18, -15}, 6, -5},
		{"test8", args{0, -2}, 0, -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := reduce(tt.args.dx, tt.args.dy)
			if got != tt.want {
				t.Errorf("reduce() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("reduce() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
