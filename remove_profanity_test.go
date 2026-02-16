package main

import "testing"

func Test_removeProfanity(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		msg  string
		want string
	}{
		{
			name: "no profanity",
			msg:  "hello world!",
			want: "hello world!",
		},
		{
			name: "one profanity",
			msg:  "hello fornax world!",
			want: "hello **** world!",
		},
		{
			name: "profanity with exclam",
			msg:  "hello kerfuffle!",
			want: "hello kerfuffle!",
		},
		{
			name: "profanity at end",
			msg:  "hello world, sharbert",
			want: "hello world, ****",
		},
		{
			name: "uppercase profanity",
			msg:  "KERFUFFLE world!",
			want: "**** world!",
		},
		{
			name: "profanity only",
			msg:  "fornax",
			want: "****",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := removeProfanity(tt.msg)
			if got != tt.want {
				t.Errorf("removeProfanity() = %v, want %v", got, tt.want)
			}
		})
	}
}
