package constraint

import "encoding/json"

// RawConstraint represents rules that we can use to define the audience of the segment.
// In other words, the audience in the segment is defined by a set of constraints.
type RawConstraint struct {
	Name       string          `json:"name"`
	Parameters json.RawMessage `json:"parameters"`
}
