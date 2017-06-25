package aozoraconv

import (
	"fmt"
)

// ParseLine parse single line and returns values
func ParseLine(s string, m, k, t *int, uni, uni2 *int32) error {
	var err error
	if _, err = fmt.Sscanf(s, "%d-%02X%02X	U+%X+%X	", m, k, t, uni, uni2); err == nil {
		return nil
	} else if _, err = fmt.Sscanf(s, "%d-%02X%02X	U+%X	", m, k, t, uni); err == nil {
		return nil
	} else if _, err = fmt.Sscanf(s, "%d-%02X%02X		", m, k, t); err == nil {
		return nil
	}
	return fmt.Errorf("could not parse %q; %v", s, err)
}
