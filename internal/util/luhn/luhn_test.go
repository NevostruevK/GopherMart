package luhn

import "testing"

func TestValid(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "first ok 9278923470",
			args: args{"9278923470"},
			want: true,
		},
		{
			name: "second ok 12345678903",
			args: args{"12345678903"},
			want: true,
		},
		{
			name: "third ok 346436439",
			args: args{"346436439"},
			want: true,
		},
		{
			name: "fourth ok small 34",
			args: args{"34"},
			want: true,
		},
		{
			name: "fifth ok the smallest 0",
			args: args{"0"},
			want: true,
		},
		{
			name: "fourth ok big 3412467861523417256376417634216437124376147",
			args: args{"3412467861523417256376417634216437124376147"},
			want: true,
		},
		{
			name: "error empty string",
			args: args{""},
			want: false,
		},
		{
			name: "error not a number 346abcdef436439",
			args: args{"346abcdef436439"},
			want: false,
		},
		{
			name: "error wrong number 346436438",
			args: args{"346436438"},
			want: false,
		},
		{
			name: "error wrong big 3412467861523417256376417634216437124376142",
			args: args{"3412467861523417256376417634216437124376142"},
			want: false,
		},
		{
			name: "error wrong the smallest 3",
			args: args{"3"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Valid([]byte(tt.args.s)); got != tt.want {
				t.Errorf("Valid() = %v, want %v", got, tt.want)
			}
		})
	}
}
