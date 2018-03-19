package utils

func KeysFromUintMap(m map[uint]EmptyStruct) []uint {
	keys := make([]uint, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// MergeStringMap merges m2 into m1, overriding any values in m1
func MergeStringMap(m1 map[string]interface{}, m2 map[string]interface{}) {
	for k, v := range m2 {
		m1[k] = v
	}
}
