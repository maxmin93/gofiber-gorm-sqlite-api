package utils

import (
	"encoding/json"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

// ExtractAndCheckField extract json fields and Check if field exists
func ExtractAndCheckField(jsonStr []byte, field string) bool {
	var keys []string = ExtractFields(jsonStr)
	if keys != nil {
		return slices.Contains(keys, field)
	}
	return false
}

// ExtractFields extract json fields as string array
func ExtractFields(jsonStr []byte) []string {
	var f interface{}
	err := json.Unmarshal(jsonStr, &f)
	if err != nil {
		return nil
	}
	m := f.(map[string]interface{})
	return maps.Keys(m)
}

// CheckField check if field exists
func CheckField(keys []string, field string) bool {
	if keys != nil {
		return slices.Contains(keys, field)
	}
	return false
}

// ReplaceKeyOfMap replace key of map
func ReplaceKeyOfMap(m map[string]interface{}, oldKey string, newKey string) map[string]interface{} {
	r := make(map[string]interface{}, len(m))
	for k, v := range m {
		if k == oldKey {
			r[newKey] = v
		} else {
			r[k] = v
		}
	}
	return r
}

// ReplaceSliceByMap replace string in slice using mapTags
func ReplaceSliceByMap(tags []string, mapTags map[string]string) []string {
	var newTags []string
	for _, s := range tags {
		if t, exist := mapTags[s]; exist {
			newTags = append(newTags, t)
		}
	}
	return newTags
}

// MapS is a map with string keys and values.
type MapS map[string]string

// Reverse returns a new map with the keys and values swapped.
func (m MapS) Reverse() map[string]string {
	n := make(map[string]string, len(m))
	for k, v := range m {
		n[v] = k
	}
	return n
}

// MapT is a map with string keys and values.
type MapT map[string]interface{}

// Filter returns a new map with matched keys
func (m MapT) Filter(keys []string) map[string]interface{} {
	n := make(map[string]interface{}, len(m))
	for k, v := range m {
		if slices.Contains(keys, k) {
			n[k] = v
		}
	}
	return n
}
