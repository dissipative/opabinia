package renderer

import (
	"reflect"
	"testing"
)

func Test_processLink(t *testing.T) {
	tests := []struct {
		name string
		link []byte
		want []byte
	}{
		{
			name: "MarkdownLinkRelative",
			link: []byte("page.md"),
			want: []byte("/page.html"),
		},
		{
			name: "MarkdownLinkExternal",
			link: []byte("https://go.dev/page.md"),
			want: []byte("https://go.dev/page.md"),
		},
		{
			name: "OtherLinkExternal",
			link: []byte("https://go.dev"),
			want: []byte("https://go.dev"),
		},
		{
			name: "EmptyLink",
			link: []byte(""),
			want: []byte(""),
		},
		{
			name: "RelativeLink",
			link: []byte("/relative"),
			want: []byte("/relative"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := processLink(tt.link); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("processLink() = %s, want %s", got, tt.want)
			}
		})
	}
}
