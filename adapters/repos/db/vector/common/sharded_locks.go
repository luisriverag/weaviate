//                           _       _
// __      _____  __ ___   ___  __ _| |_ ___
// \ \ /\ / / _ \/ _` \ \ / / |/ _` | __/ _ \
//  \ V  V /  __/ (_| |\ V /| | (_| | ||  __/
//   \_/\_/ \___|\__,_| \_/ |_|\__,_|\__\___|
//
//  Copyright © 2016 - 2024 Weaviate B.V. All rights reserved.
//
//  CONTACT: hello@weaviate.io
//

package common

import "sync"

const DefaultShardedLocksCount = 512

type ShardedLocks struct {
	// sharded locks
	shards []sync.RWMutex
	// number of locks
	count uint64
}

func NewDefaultShardedLocks() *ShardedLocks {
	return NewShardedLocks(DefaultShardedLocksCount)
}

func NewShardedLocks(count uint64) *ShardedLocks {
	if count < 2 {
		count = 2
	}

	return &ShardedLocks{
		shards: make([]sync.RWMutex, count),
		count:  count,
	}
}

func (sl *ShardedLocks) LockAll() {
	for i := uint64(0); i < sl.count; i++ {
		sl.shards[i].Lock()
	}
}

func (sl *ShardedLocks) UnlockAll() {
	for i := int(sl.count) - 1; i >= 0; i-- {
		sl.shards[i].Unlock()
	}
}

func (sl *ShardedLocks) LockedAll(callback func()) {
	sl.LockAll()
	defer sl.UnlockAll()

	callback()
}

func (sl *ShardedLocks) Lock(id uint64) {
	sl.shards[id%sl.count].Lock()
}

func (sl *ShardedLocks) Unlock(id uint64) {
	sl.shards[id%sl.count].Unlock()
}

func (sl *ShardedLocks) Locked(id uint64, callback func()) {
	sl.Lock(id)
	defer sl.Unlock(id)

	callback()
}

func (sl *ShardedLocks) RLockAll() {
	for i := uint64(0); i < sl.count; i++ {
		sl.shards[i].RLock()
	}
}

func (sl *ShardedLocks) RUnlockAll() {
	for i := int(sl.count) - 1; i >= 0; i-- {
		sl.shards[i].RUnlock()
	}
}

func (sl *ShardedLocks) RLockedAll(callback func()) {
	sl.RLockAll()
	defer sl.RUnlockAll()

	callback()
}

func (sl *ShardedLocks) RLock(id uint64) {
	sl.shards[id%sl.count].RLock()
}

func (sl *ShardedLocks) RUnlock(id uint64) {
	sl.shards[id%sl.count].RUnlock()
}

func (sl *ShardedLocks) RLocked(id uint64, callback func()) {
	sl.RLock(id)
	defer sl.RUnlock(id)

	callback()
}
