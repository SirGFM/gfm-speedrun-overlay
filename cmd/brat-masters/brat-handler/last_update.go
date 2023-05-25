package brat_handler

import (
	"sync/atomic"
	"time"
)

// getLastUpdate retrieves the last updated time, in milliseconds.
func (ctx *serverContext) getLastUpdate() int64 {
	ctx.lastUpdate.m.RLock()
	ms := ctx.lastUpdate.t.Unix()
	ctx.lastUpdate.m.RUnlock()
	return ms * 1000
}

// setLastUpdate sets the update date, so watching clients may automatically refresh.
func (ctx *serverContext) setLastUpdate() {
	// (tl;dr: This is overkill)
	//
	// Ensure that only a single writer may ever acquire the write lock
	// (using atomic.CompareAndSwapUint32 to only allow a single goroutine
	// to enter the conditional at a time) and then lock it for writing.
	if atomic.CompareAndSwapUint32(&ctx.lastUpdate.tryWrite, 0, 1) {
		ctx.lastUpdate.m.Lock()
		ctx.lastUpdate.t = time.Now()
		ctx.lastUpdate.m.Unlock()
		atomic.StoreUint32(&ctx.lastUpdate.tryWrite, 0)
	}
}
