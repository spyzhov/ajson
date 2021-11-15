package tokens

import (
	"fmt"
	"math"
	"strconv"
)

// Number is a temporary token
type Number struct {
	Alias string
	Value float64
	IsInt bool
}

var _ Token = (*Number)(nil)

func NewNumber(value string) (*Number, error) {
	float, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return nil, fmt.Errorf("value %q can't be converted to float64: %w", value, err)
	}
	return &Number{
		Alias: value,
		Value: float,
		IsInt: math.Mod(float, 1) == 0,
	}, nil
}

func (o *Number) String() string {
	if o == nil {
		return "<nil>"
	}
	return o.Alias
}
