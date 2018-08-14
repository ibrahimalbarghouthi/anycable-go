// +build darwin,mrb linux,mrb

package rpc_cached

import (
	"log"
	"testing"

	"github.com/anycable/anycable-go/mrb"

	"github.com/stretchr/testify/assert"
)

var (
	cache *MCache
)

func init() {
	var err error
	cache, err = NewMCache(mrb.DefaultEngine())

	if err != nil {
		log.Fatalf("Failed to initialize mruby cache: %s", err)
	}
}

func TestMAction(t *testing.T) {
	maction, err := NewMAction(
		cache,
		"BenchmarkChannel",
		`
		def echo(data)
			transmit response: data
		end
		`,
	)

	assert.Nil(t, err)

	cache.engine.VM.FullGC()

	origObjects := cache.engine.VM.LiveObjectCount()

	res, err := maction.Perform("{\"action\":\"echo\",\"text\":\"hello\"}")

	assert.Nil(t, err)

	newObjects := cache.engine.VM.LiveObjectCount()

	if origObjects != newObjects {
		t.Fatalf("Object count was not what was expected after action call: %d %d", origObjects, newObjects)
	}

	identifier := "{\\\"channel\\\":\\\"BenchmarkChannel\\\"}"

	assert.Empty(t, res.Streams)
	assert.False(t, res.Disconnect)
	assert.False(t, res.StopAllStreams)
	assert.Equal(t, []string{"{\"identifier\":\"" + identifier + "\",\"message\":{\"response\":{\"action\":\"echo\",\"text\":\"hello\"}}}"}, res.Transmissions)
}