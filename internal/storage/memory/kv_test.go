package memory

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestShouldStoreAndRetrieveDataFromMemoryByKey(t *testing.T) {
	type testCase struct {
		name   string
		set    map[string]string
		get    map[string]bool
		exists map[string]bool
	}
	tests := []testCase{
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kv := NewKVStorage[string]()
			for key, value := range tt.set {
				kv.Set(key, value)
			}
			for key, value := range tt.get {
				result, _ := kv.Get(key)
				require.Equal(t, value, result)
			}
			for key, value := range tt.exists {
				_, result := kv.Get(key)
				require.Equal(t, value, result)
			}
		})
	}
}
