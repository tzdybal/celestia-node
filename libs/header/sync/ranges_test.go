package sync

import (
	"github.com/celestiaorg/celestia-node/libs/header/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAddParallel(t *testing.T) {
	var pending ranges[*test.DummyHeader]

	n := 500
	suite := test.NewTestSuite(t)
	headers := suite.GenDummyHeaders(1000)
	for i := 0; i < n; i++ {
		go func() {
			pending.Add(headers[i])
		}()
	}

	last := uint64(0)
	for _, r := range pending.ranges {
		assert.Greater(t, r.start, last)
		last = r.start
	}
}

func TestBeforeOnRektRanges(t *testing.T) {
	var pending ranges[*test.DummyHeader]

	suite := test.NewTestSuite(t)
	headers := suite.GenDummyHeaders(10)

	pending.Add(headers[0])
	pending.Add(headers[1])
	pending.Add(headers[2])
	pending.Add(headers[4])
	pending.ranges = append(pending.ranges, newRange(headers[3]))

	// first "sync"
	r, ok := pending.FirstRangeWithin(1, 10)
	require.True(t, ok)
	r.Before(10)

	// second "sync"
	r, ok = pending.FirstRangeWithin(1, 6)
	r.Before(6)
}
