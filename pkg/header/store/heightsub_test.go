package store

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/celestiaorg/celestia-node/pkg/header"
)

func TestHeightSub(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	hs := newHeightSub[*header.DummyHeader]()

	// assert subscription returns nil for past heights
	{
		h := header.RandDummyHeader(t)
		h.Raw.Height = 100
		hs.SetHeight(99)
		hs.Pub(h)

		h, err := hs.Sub(ctx, 10)
		assert.ErrorIs(t, err, errElapsedHeight)
		assert.Nil(t, h)
	}

	// assert actual subscription works
	{
		go func() {
			// fixes flakiness on CI
			time.Sleep(time.Millisecond)

			h1 := header.RandDummyHeader(t)
			h1.Raw.Height = 101
			h2 := header.RandDummyHeader(t)
			h2.Raw.Height = 102
			hs.Pub(h1, h2)
		}()

		h, err := hs.Sub(ctx, 101)
		assert.NoError(t, err)
		assert.NotNil(t, h)
	}
}