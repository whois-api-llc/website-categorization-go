package websitecategorization

import (
	"net/url"
	"reflect"
	"testing"
)

// TestOptions tests the Options functions.
func TestOptions(t *testing.T) {
	tests := []struct {
		name   string
		values url.Values
		option Option
		want   string
	}{
		{
			name:   "outputFormat",
			values: url.Values{},
			option: OptionOutputFormat("JSON"),
			want:   "outputFormat=JSON",
		},
		{
			name:   "minConfidence",
			values: url.Values{},
			option: OptionMinConfidence(0.55),
			want:   "minConfidence=0.550000",
		},
		{
			name:   "order",
			values: url.Values{},
			option: OptionOrder("ABC"),
			want:   "order=ABC",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.option(tt.values)
			if got := tt.values.Encode(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Option() = %v, want %v", got, tt.want)
			}
		})
	}
}
