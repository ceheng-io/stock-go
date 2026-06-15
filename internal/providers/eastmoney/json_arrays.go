package eastmoney

import (
	"bytes"
	"encoding/json"
)

func decodeBoardDynamicArray(raw json.RawMessage) ([]boardDynamicItem, error) {
	if !isJSONArrayPayload(raw) {
		return []boardDynamicItem{}, nil
	}
	var items []boardDynamicItem
	if err := json.Unmarshal(raw, &items); err != nil {
		return nil, err
	}
	return items, nil
}

func decodeStringArray(raw json.RawMessage) ([]string, error) {
	if !isJSONArrayPayload(raw) {
		return []string{}, nil
	}
	var items []string
	if err := json.Unmarshal(raw, &items); err != nil {
		return nil, err
	}
	return items, nil
}

func decodeMapArray(raw json.RawMessage) ([]map[string]any, error) {
	if !isJSONArrayPayload(raw) {
		return []map[string]any{}, nil
	}
	var items []map[string]any
	if err := json.Unmarshal(raw, &items); err != nil {
		return nil, err
	}
	return items, nil
}

func isJSONArrayPayload(raw json.RawMessage) bool {
	raw = bytes.TrimSpace(raw)
	return len(raw) > 0 && raw[0] == '['
}
