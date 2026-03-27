// Package jsonfields provides utilities for filtering JSON output by field names.
package jsonfields

import (
	"encoding/json"
	"fmt"
)

// FilterFields marshals v to JSON then returns only the requested top-level keys.
// For slices, applies filtering to each element. If fields is nil/empty, returns full JSON.
func FilterFields(v any, fields []string) (json.RawMessage, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("marshal: %w", err)
	}
	if len(fields) == 0 {
		return data, nil
	}

	// Detect if top-level is array or object
	if len(data) > 0 && data[0] == '[' {
		var items []json.RawMessage
		if err := json.Unmarshal(data, &items); err != nil {
			return nil, err
		}
		filtered := make([]json.RawMessage, len(items))
		for i, item := range items {
			filtered[i], err = filterObject(item, fields)
			if err != nil {
				return nil, err
			}
		}
		return json.Marshal(filtered)
	}
	return filterObject(data, fields)
}

func filterObject(data json.RawMessage, fields []string) (json.RawMessage, error) {
	var m map[string]json.RawMessage
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("unmarshal object: %w", err)
	}
	result := make(map[string]json.RawMessage, len(fields))
	for _, f := range fields {
		if v, ok := m[f]; ok {
			result[f] = v
		}
	}
	out, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("marshal filtered: %w", err)
	}
	return out, nil
}
