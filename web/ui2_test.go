package web

import (
	"fmt"
	"testing"
)

func Test_queryJoin(t *testing.T) {
	type args struct {
		base string
		args []string
	}
	tests := []struct {
		args args
		want string
	}{
		{
			args: args{
				base: "/search",
				args: []string{
					"tags", "",
					"title", "",
					"text", "",
				},
			},
			want: "/search",
		},
		{
			args: args{
				base: "/search",
				args: []string{
					"tags", "foo",
					"title", "bar",
					"text", "baz",
				},
			},
			want: "/search?tags=foo&title=bar&text=baz",
		},
		{
			args: args{
				base: "/search",
				args: []string{
					"tags", "",
					"title", "",
					"text", "foo&bar?baz",
				},
			},
			want: "/search?text=foo%26bar%3Fbaz",
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("test %v", i), func(t *testing.T) {
			if got := queryJoin(tt.args.base, tt.args.args); got != tt.want {
				t.Errorf("queryJoin() = %v, want %v", got, tt.want)
			}
		})
	}
}
