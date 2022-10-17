package websitecategorization

import (
	"fmt"
	"net/url"
	"strings"
)

// Option adds parameters to the query.
type Option func(v url.Values)

var _ = []Option{
	OptionOutputFormat("JSON"),
	OptionMinConfidence(0.55),
	OptionOrder("IAB"),
}

// OptionOutputFormat sets Response output format JSON | XML | CSV. Default: JSON.
func OptionOutputFormat(outputFormat string) Option {
	return func(v url.Values) {
		v.Set("outputFormat", strings.ToUpper(outputFormat))
	}
}

// OptionMinConfidence sets The minimum confidence for the predictions. The higher this value the fewer
// false-positive results will be returned. Acceptable values: 0.00 - 1.00. Default: 0.55.
func OptionMinConfidence(value float64) Option {
	return func(v url.Values) {
		v.Set("minConfidence", fmt.Sprintf("%f", value))
	}
}

// OptionOrder sets the categories output order (for GetAllCategories functions only).
// ABC - output categories ordered alphabetically by the name field. IAB - output categories ordered by the id field.
// Acceptable values: ABC | IAB. Default: IAB.
func OptionOrder(order string) Option {
	return func(v url.Values) {
		v.Set("order", strings.ToUpper(order))
	}
}
