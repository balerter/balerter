package loki

import "fmt"

func directionValidate(v string) error {
	if v != "" && v != "forward" && v != "backward" {
		return fmt.Errorf("option Direction support only values: 'forward' and 'backward'")
	}

	return nil
}
