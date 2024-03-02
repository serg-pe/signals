package queue

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPush(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name             string
		pushSequence     []int
		expectedSequence []int
	}{
		{
			name:             "push to empty queue",
			pushSequence:     []int{0},
			expectedSequence: []int{0},
		},
		{
			name:             "push multiple values",
			pushSequence:     []int{0, 1, 2, 3, 4},
			expectedSequence: []int{0, 1, 2, 3, 4},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			q := New[int]()

			for _, val := range tc.pushSequence {
				q.Push(val)
			}

			actualNode := q.head
			expectedCounter := 0
			expectedLen := len(tc.expectedSequence)
			for actualNode != nil && expectedCounter < expectedLen {
				assert.Equal(t, tc.expectedSequence[expectedCounter], actualNode.data)
				actualNode = actualNode.next
				expectedCounter++
			}
			assert.Equal(t, expectedCounter, expectedLen)
			assert.Nil(t, actualNode)
		})
	}
}

func TestPop(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		pushSequence   []int
		popCount       int
		expectedResult int
		expectedStatus bool
	}{
		{
			name:           "pop from empty queue",
			pushSequence:   []int{},
			popCount:       1,
			expectedResult: 0,
			expectedStatus: false,
		},
		{
			name:           "pop from empty queue 2 times",
			pushSequence:   []int{},
			popCount:       2,
			expectedResult: 0,
			expectedStatus: false,
		},
		{
			name:           "pop from filled queue",
			pushSequence:   []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			popCount:       5,
			expectedResult: 5,
			expectedStatus: true,
		},
		{
			name:           "pop last",
			pushSequence:   []int{1, 2, 3},
			popCount:       3,
			expectedResult: 3,
			expectedStatus: true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			q := New[int]()

			for _, val := range tc.pushSequence {
				q.Push(val)
			}

			var (
				actualResult int
				actualStatus bool
			)
			for i := 0; i < tc.popCount; i++ {
				actualResult, actualStatus = q.Pop()
			}
			assert.Equal(t, tc.expectedResult, actualResult)
			assert.Equal(t, tc.expectedStatus, actualStatus)
		})
	}
}
