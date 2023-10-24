package parser

import (
	"reflect"
	"testing"
)

func Test_extractUniqHTMLLinks(t *testing.T) {
	tests := []struct {
		name    string
		html    string
		want    map[string]struct{}
		wantErr bool
	}{
		{
			name: "MultipleLinksInSourceHTML",
			html: `<a href="https://www.google.com">Google</a>
<a href="https://www.google.com">Google</a>
<a href="https://www.google.com">Google</a>
<a href="https://google.com">Google</a>`,
			want: map[string]struct{}{
				"https://www.google.com": {},
				"https://google.com":     {},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := extractUniqHTMLLinks(tt.html)
			if (err != nil) != tt.wantErr {
				t.Errorf("extractUniqHTMLLinks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extractUniqHTMLLinks() got = %v, want %v", got, tt.want)
			}
		})
	}
}
