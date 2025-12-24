package validators

import (
	"errors"
	"strconv"
)

// RangeParams represents validated range parameters
type RangeParams struct {
	Min *int
	Max *int
}

// ParseRangeParams validates and parses min/max query parameters
func ParseRangeParams(minStr, maxStr string) (*RangeParams, error) {
	params := &RangeParams{}
	
	if minStr != "" {
		min, err := strconv.Atoi(minStr)
		if err != nil {
			return nil, errors.New("invalid min parameter: must be a number")
		}
		if min < 0 {
			return nil, errors.New("min must be >= 0")
		}
		params.Min = &min
	}
	
	if maxStr != "" {
		max, err := strconv.Atoi(maxStr)
		if err != nil {
			return nil, errors.New("invalid max parameter: must be a number")
		}
		if max < 0 {
			return nil, errors.New("max must be >= 0")
		}
		params.Max = &max
	}
	
	// Validate range logic
	if params.Min != nil && params.Max != nil && *params.Min > *params.Max {
		return nil, errors.New("min must be <= max")
	}
	
	return params, nil
}

// GetMinMax returns min and max values, defaulting to 0 if nil
func (rp *RangeParams) GetMinMax() (min, max int) {
	if rp.Min != nil {
		min = *rp.Min
	}
	if rp.Max != nil {
		max = *rp.Max
	}
	return min, max
}
