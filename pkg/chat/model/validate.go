package model

import (
	"fmt"
	"slices"
)

// oneOf checks that at most one of the fields in a protobuf oneof group is
// set. The map is keyed by JSON field name.
func oneOf(group string, set map[string]bool) error {
	var found []string
	for name, ok := range set {
		if ok {
			found = append(found, name)
		}
	}
	if len(found) > 1 {
		slices.Sort(found)
		return fmt.Errorf("only one of the %s fields may be set, got %v", group, found)
	}
	return nil
}
