package array

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		data        []stored[int]
		at          int
		expected    int
		expectedErr error
	}{
		{
			name:        "get first",
			data:        []stored[int]{{true, 10}, {true, 20}, {true, 30}},
			at:          0,
			expected:    10,
			expectedErr: nil,
		},
		{
			name:        "get middle",
			data:        []stored[int]{{true, 10}, {true, 20}, {true, 30}},
			at:          1,
			expected:    20,
			expectedErr: nil,
		},
		{
			name:        "get last",
			data:        []stored[int]{{true, 10}, {true, 20}, {true, 30}},
			at:          2,
			expected:    30,
			expectedErr: nil,
		},
		{
			name:        "get unused",
			data:        []stored[int]{{false, 20}},
			at:          0,
			expected:    0,
			expectedErr: errAlreadyRemoved,
		},
		{
			name:        "get from minus range",
			data:        []stored[int]{{true, 10}, {true, 20}, {true, 30}},
			at:          -1,
			expected:    0,
			expectedErr: errNegativeRangeNotSupported,
		},
		{
			name:        "out of range",
			data:        []stored[int]{{true, 10}, {true, 20}, {true, 30}},
			at:          3,
			expected:    0,
			expectedErr: errOutOfRange,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			s := ArrayStorage[int]{}
			s.capacity = len(tc.data)
			s.length = s.capacity
			s.storage = make([]stored[int], s.capacity)
			copy(s.storage, tc.data)

			actual, err := s.Get(tc.at)
			if err == nil {
				assert.NoError(t, err)
			} else {
				assert.ErrorIs(t, err, tc.expectedErr)
			}

			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestAdd(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		capacity    int
		data        []int
		expected    []stored[int]
		expectedCap int
		expectedLen int
	}{
		{
			name:     "add without reallocation cap == length",
			data:     []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			capacity: 10,
			expected: []stored[int]{
				{true, 1}, {true, 2}, {true, 3}, {true, 4},
				{true, 5}, {true, 6}, {true, 7}, {true, 8},
				{true, 9}, {true, 10},
			},
			expectedCap: 10,
			expectedLen: 10,
		},
		{
			name:     "add with reallocation",
			capacity: 10,
			data:     []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11},
			expected: []stored[int]{
				{true, 1}, {true, 2}, {true, 3}, {true, 4},
				{true, 5}, {true, 6}, {true, 7}, {true, 8},
				{true, 9}, {true, 10}, {true, 11},
				{false, 0}, {false, 0}, {false, 0},
				{false, 0}, {false, 0}, {false, 0}, {false, 0},
				{false, 0}, {false, 0},
			},
			expectedCap: 20,
			expectedLen: 11,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			defer func() {
				err := recover()
				assert.Nil(t, err)
			}()

			s := New[int](tc.capacity)

			assert.Equal(t, tc.capacity, cap(s.storage))

			for _, item := range tc.data {
				s.Add(item)
			}

			assert.Equal(t, tc.expectedCap, cap(s.storage))
			assert.Equal(t, s.capacity, tc.expectedCap)
			assert.Equal(t, s.length, tc.expectedLen)

			for id, val := range s.storage {
				assert.Equal(t, tc.expected[id], val)
			}
		})
	}
}

func TestRemove(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		data        []stored[int]
		capacity    int
		length      int
		removeAt    int
		expectedErr error
	}{
		{
			name:        "remove used",
			data:        []stored[int]{{true, 0}, {true, 1}, {true, 2}},
			capacity:    3,
			length:      3,
			removeAt:    1,
			expectedErr: nil,
		},
		{
			name:        "remove out of range",
			data:        []stored[int]{{true, 0}, {true, 1}, {true, 2}},
			capacity:    5,
			length:      3,
			removeAt:    4,
			expectedErr: errOutOfRange,
		},
		{
			name:        "remove already removed",
			data:        []stored[int]{{true, 0}, {false, 1}, {true, 2}},
			capacity:    3,
			length:      3,
			removeAt:    1,
			expectedErr: errAlreadyRemoved,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			s := ArrayStorage[int]{
				capacity: tc.capacity,
				length:   tc.length,
			}
			s.storage = make([]stored[int], s.capacity)
			copy(s.storage, tc.data)

			err := s.Remove(tc.removeAt)
			if tc.expectedErr == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tc.expectedErr)
			}
		})
	}
}

func TestApplyTo(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		sequence   []stored[int]
		filterFunc func(entry int) bool
		apply      func(entry int) int
		expected   []int
	}{
		{
			name: "apply to odd",
			sequence: []stored[int]{
				{true, 1}, {true, 2}, {true, 3}, {true, 4},
				{true, 5}, {true, 6}, {true, 7}, {true, 8},
				{true, 9}, {true, 10},
			},
			filterFunc: func(entry int) bool {
				return entry%2 == 0
			},
			apply: func(entry int) int {
				return entry * entry
			},
			expected: []int{1, 4, 3, 16, 5, 36, 7, 64, 9, 100},
		},
		{
			name:     "apply greater with unused",
			sequence: []stored[int]{{true, 1}, {false, 2}, {false, 3}, {true, 4}, {true, 4}, {true, 6}},
			filterFunc: func(entry int) bool {
				return entry != 4
			},
			apply: func(entry int) int {
				return entry + 1
			},
			expected: []int{2, 2, 3, 4, 4, 7},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			s := ArrayStorage[int]{
				length:   len(tc.sequence),
				capacity: cap(tc.sequence),
				storage:  make([]stored[int], cap(tc.sequence)),
			}
			copy(s.storage, tc.sequence)

			s.ApplyTo(tc.filterFunc, tc.apply)

			assert.Equal(t, len(tc.sequence), len(s.storage))
			assert.Len(t, tc.sequence, s.length)
			assert.Equal(t, cap(tc.sequence), cap(s.storage))
			assert.Equal(t, cap(s.storage), s.capacity)

			for id, val := range s.storage {
				assert.Equal(t, tc.expected[id], val.data)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		sequence    []stored[int]
		id          int
		update      func(entry int) int
		expected    []int
		expectedErr error
	}{
		{
			name:        "update used",
			sequence:    []stored[int]{{true, 1}, {true, 2}, {true, 3}},
			id:          0,
			update:      func(entry int) int { return entry * 10 },
			expected:    []int{10, 2, 3},
			expectedErr: nil,
		},
		{
			name:        "update unused",
			sequence:    []stored[int]{{true, 1}, {false, 2}, {true, 3}},
			id:          1,
			update:      func(entry int) int { return entry * 10 },
			expected:    []int{1, 2, 3},
			expectedErr: errAlreadyRemoved,
		},
		{
			name:        "update out of range",
			sequence:    []stored[int]{{true, 1}, {false, 2}, {true, 3}},
			id:          10,
			update:      func(entry int) int { return entry * 10 },
			expected:    []int{1, 2, 3},
			expectedErr: errOutOfRange,
		},
		{
			name:        "update in negative range",
			sequence:    []stored[int]{{true, 1}, {false, 2}, {true, 3}},
			id:          -1,
			update:      func(entry int) int { return entry * 10 },
			expected:    []int{1, 2, 3},
			expectedErr: errNegativeRangeNotSupported,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			s := ArrayStorage[int]{
				length:   len(tc.sequence),
				capacity: cap(tc.sequence),
				storage:  make([]stored[int], cap(tc.sequence)),
			}
			copy(s.storage, tc.sequence)

			actualErr := s.Update(tc.id, tc.update)

			assert.Equal(t, len(tc.sequence), len(s.storage))
			assert.Len(t, tc.sequence, s.length)
			assert.Equal(t, cap(tc.sequence), cap(s.storage))
			assert.Equal(t, cap(s.storage), s.capacity)

			if tc.expectedErr == nil {
				assert.NoError(t, actualErr)
			} else {
				assert.ErrorIs(t, actualErr, tc.expectedErr)
			}

			for id, val := range s.storage {
				assert.Equal(t, tc.expected[id], val.data)
			}
		})
	}
}
