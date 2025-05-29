package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopyMap(t *testing.T) {
	dst := map[string]any{
		"key1": map[string]any{
			"key2": "val2",
		},
		"key3": 3,
		"key4": map[string]any{
			"key5": map[string]any{
				"key6": 100.2,
				"key7": "val7",
			},
		},
	}
	src := map[string]any{
		"key1": 100,
		"key4": map[string]any{
			"key5": map[string]any{
				"key6": "val6",
			},
			"key8": "val8",
		},
	}

	CopyMap(dst, src)

	require.Equal(t, map[string]any{
		"key1": 100,
		"key3": 3,
		"key4": map[string]any{
			"key5": map[string]any{
				"key6": "val6",
				"key7": "val7",
			},
			"key8": "val8",
		},
	}, dst)
}
