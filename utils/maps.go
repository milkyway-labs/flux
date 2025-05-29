package utils

import (
	"fmt"
	"reflect"
)

// CopyMap recursively copies values from src map to dst map.
func CopyMap(dst, src map[string]any) {
	for k, v := range src {
		d, exists := dst[k]
		if !exists {
			dst[k] = v
			continue
		}

		// We only care about maps, if the value is not a map, we just overwrite it
		dstVal := reflect.ValueOf(d)
		srcVal := reflect.ValueOf(v)
		if dstVal.Kind() != reflect.Map || srcVal.Kind() != reflect.Map {
			dst[k] = v
			continue
		}

		// Construct maps from reflect.Value
		dstMap := map[string]any{}
		for k2, v2 := range dstVal.Seq2() {
			dstMap[fmt.Sprint(k2.Interface())] = v2.Interface()
		}
		srcMap := map[string]any{}
		for k2, v2 := range srcVal.Seq2() {
			srcMap[fmt.Sprint(k2.Interface())] = v2.Interface()
		}

		CopyMap(dstMap, srcMap)
		dst[k] = dstMap
	}
}
