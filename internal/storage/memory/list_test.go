package memory

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestShouldReturnOnlyUpToDateEntries(t *testing.T) {
	type testCase struct {
		name    string
		entries Entries[string]
		now     time.Time
		want    Entries[string]
	}
	outdatedFirst := Entry[string]{Value: "first", Expiration: time.Now().Add(-time.Hour)}
	outdatedSecond := Entry[string]{Value: "first", Expiration: time.Now().Add(-time.Minute)}
	activeFirst := Entry[string]{Value: "first", Expiration: time.Now().Add(time.Hour)}
	activeSecond := Entry[string]{Value: "second", Expiration: time.Now().Add(time.Minute)}
	tests := []testCase{
		{
			name:    "should return empty list if all entries are outdated",
			entries: Entries[string]{outdatedFirst, outdatedSecond},
			now:     time.Now(),
			want:    Entries[string]{},
		},
		{
			name:    "should not affect slice if all entries are up to date",
			entries: Entries[string]{activeFirst, activeSecond},
			now:     time.Now(),
			want:    Entries[string]{activeFirst, activeSecond},
		},
		{
			name:    "should remove only outdated entries",
			entries: Entries[string]{activeFirst, outdatedFirst, activeSecond, outdatedSecond},
			now:     time.Now(),
			want:    Entries[string]{activeFirst, activeSecond},
		},
		{
			name:    "should operates on empty list",
			entries: Entries[string]{},
			now:     time.Now(),
			want:    Entries[string]{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.entries.Expire(tt.now))
		})
	}
}

func TestShouldProvideOnlyValues(t *testing.T) {
	type testCase struct {
		name    string
		entries Entries[string]
		want    []string
	}
	tests := []testCase{
		{
			name:    "should return empty list for empty entries",
			entries: Entries[string]{},
			want:    []string{},
		},
		{
			name:    "should return values for entries",
			entries: Entries[string]{{Value: "first", Expiration: time.Now()}},
			want:    []string{"first"},
		},
		{
			name:    "should return values for multiple entries",
			entries: Entries[string]{{Value: "first", Expiration: time.Now()}, {Value: "second", Expiration: time.Now()}},
			want:    []string{"first", "second"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.entries.Values())
		})
	}
}

func TestShouldAppendNewEntriesForKeyAndFlushThemAfterRead(t *testing.T) {
	type testCase struct {
		name     string
		key      string
		before   []string
		after    []string
		expected []string
	}
	tests := []testCase{
		{
			name:     "should append entries for new key",
			key:      "key",
			before:   []string{},
			after:    []string{"first", "second"},
			expected: []string{"first", "second"},
		},
		{
			name:     "should append entries for existing key",
			key:      "key",
			before:   []string{"first", "second"},
			after:    []string{"third", "forth"},
			expected: []string{"first", "second", "third", "forth"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewListStorage[string]()
			for _, value := range tt.before {
				s.Append(tt.key, value, time.Second)
			}
			for _, value := range tt.after {
				s.Append(tt.key, value, time.Second)
			}
			require.Equal(t, []string{tt.key}, s.GetKeys())
			require.Equal(t, tt.expected, s.FetchAndFlush(tt.key))
			require.Equal(t, []string{}, s.FetchAndFlush(tt.key))
		})
	}
}

func TestShouldProvideOnlyActiveData(t *testing.T) {
	type testCase struct {
		name     string
		key      string
		outdated []string
		active   []string
		expected []string
	}
	tests := []testCase{
		{
			name:     "should return empty list for empty storage",
			key:      "key",
			outdated: []string{},
			active:   []string{},
			expected: []string{},
		},
		{
			name:     "should return only active",
			key:      "key",
			outdated: []string{"first", "second"},
			active:   []string{"third"},
			expected: []string{"third"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewListStorage[string]()
			for _, value := range tt.outdated {
				s.Append(tt.key, value, -time.Second)
			}
			for _, value := range tt.active {
				s.Append(tt.key, value, time.Second)
			}
			require.Equal(t, tt.expected, s.FetchAndFlush(tt.key))
		})
	}
}

func TestShouldCleanEmptyLists(t *testing.T) {
	type testCase struct {
		name     string
		key      string
		outdated []string
		active   []string
		expected []string
	}
	tests := []testCase{
		{
			name:     "should return no keys for empty storage",
			key:      "key",
			outdated: []string{},
			active:   []string{},
			expected: []string{},
		},
		{
			name:     "should return key if there is list containing active entries",
			key:      "key",
			outdated: []string{"first", "second"},
			active:   []string{"third"},
			expected: []string{"key"},
		},
		{
			name:     "should return no keys after cleaning all outdated entries",
			key:      "key",
			outdated: []string{"first", "second"},
			active:   []string{},
			expected: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewListStorage[string]()
			for _, value := range tt.outdated {
				s.Append(tt.key, value, -time.Second)
			}
			for _, value := range tt.active {
				s.Append(tt.key, value, time.Second)
			}
			s.CleanOutdated(tt.key)
			require.Equal(t, tt.expected, s.GetKeys())
		})
	}
}
