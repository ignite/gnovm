package keeper

import (
	"context"

	corestore "cosmossdk.io/core/store"
	"cosmossdk.io/log"
	gnosdk "github.com/gnolang/gno/tm2/pkg/sdk"
	gnostore "github.com/gnolang/gno/tm2/pkg/store"
)

var _ gnostore.MultiStore = (*gnovmMultiStore)(nil)

// gnovmMultiStore is a wrapper around the Cosmos SDK store service that implements
// the gno store.MultiStore interface while restricting access to only the gnovm store.
type gnovmMultiStore struct {
	logger          log.Logger
	storeService    corestore.KVStoreService
	memStoreService corestore.MemoryStoreService
	storeKey        gnostore.StoreKey
	memStoreKey     gnostore.StoreKey
	ctx             gnosdk.Context
	sdkCtx          context.Context
	kvStore         corestore.KVStore // Cached regular KV store
	memStore        corestore.KVStore // Cached memory store
}

// NewGnovmMultiStore creates a new MultiStore wrapper that can be initialized
// without a context and will lazily initialize the store when needed.
func NewGnovmMultiStore(
	logger log.Logger,
	storeService corestore.KVStoreService,
	memStoreService corestore.MemoryStoreService,
	storeKey gnostore.StoreKey,
	memStoreKey gnostore.StoreKey,
) gnostore.MultiStore {
	return &gnovmMultiStore{
		logger:          logger,
		storeService:    storeService,
		storeKey:        storeKey,
		memStoreKey:     memStoreKey,
		memStoreService: memStoreService,
	}
}

// SetContext sets the contexts for the store wrapper. This allows lazy initialization.
func (ms *gnovmMultiStore) SetContext(ctx gnosdk.Context, sdkCtx context.Context) {
	ms.ctx = ctx.WithMultiStore(ms)
	ms.sdkCtx = sdkCtx
	// Always reset cached stores to ensure fresh state for each operation
	ms.kvStore = nil
	ms.memStore = nil
}

// GetStore implements types.MultiStore.
// It returns the gnovm store if the provided key matches our store key,
// otherwise it panics as per the interface contract.
func (ms *gnovmMultiStore) GetStore(key gnostore.StoreKey) gnostore.Store {
	if ms.sdkCtx == nil {
		panic("SDK context not set - call SetContext first")
	}

	var kvStore corestore.KVStore

	if key.Name() == ms.memStoreKey.Name() {
		// Always create fresh memory store instance to avoid cross-transaction contamination
		var memStore corestore.KVStore
		if ms.memStoreService != nil {
			func() {
				defer func() {
					if r := recover(); r != nil {
						// Memory store not available, fall back to regular store
						ms.logger.Info("Memory store not available, falling back to regular store", "error", r)
						memStore = nil
					}
				}()
				memStore = ms.memStoreService.OpenMemoryStore(ms.sdkCtx)
			}()
		}
		// If memory store service is not available or failed, fall back to regular store
		if memStore == nil {
			memStore = ms.storeService.OpenKVStore(ms.sdkCtx)
		}
		kvStore = memStore
	} else if key.Name() == ms.storeKey.Name() {
		// Always create fresh regular store instance to avoid cross-transaction contamination
		kvStore = ms.storeService.OpenKVStore(ms.sdkCtx)
	} else {
		panic("store not found: " + key.Name())
	}

	return &gnovmStore{
		logger:  ms.logger,
		kvStore: kvStore,
	}
}

// MultiCacheWrap implements types.MultiStore.
// Returns a cache-wrapped version of this MultiStore.
func (ms *gnovmMultiStore) MultiCacheWrap() gnostore.MultiStore {
	// Create a completely independent multistore instance for proper isolation
	cached := &gnovmMultiStore{
		logger:          ms.logger,
		storeService:    ms.storeService,
		memStoreService: ms.memStoreService,
		storeKey:        ms.storeKey,
		memStoreKey:     ms.memStoreKey,
		ctx:             ms.ctx,
		sdkCtx:          ms.sdkCtx,
		// Don't copy any cached stores - ensure complete isolation
		kvStore:  nil,
		memStore: nil,
	}
	return cached
}

// MultiWrite implements types.MultiStore.
// Flushes any cached writes to the underlying store.
// Since we're using the store service directly, this is a no-op.
func (ms *gnovmMultiStore) MultiWrite() {
	// No-op as the store service handles writes directly
}

var _ gnostore.Store = (*gnovmStore)(nil)

// gnovmStore implements the gno Store interface using the Cosmos SDK KVStore.
type gnovmStore struct {
	logger  log.Logger
	kvStore corestore.KVStore
}

// Get implements types.Store.
func (s *gnovmStore) Get(key []byte) []byte {
	value, err := s.kvStore.Get(key)
	if err != nil {
		s.logger.Error("failed to get value from store", "key", string(key), "error", err)
		return nil
	}
	return value
}

// Has implements types.Store.
func (s *gnovmStore) Has(key []byte) bool {
	has, err := s.kvStore.Has(key)
	if err != nil {
		s.logger.Error("failed to check if key exists in store", "key", string(key), "error", err)
		return false
	}
	return has
}

// Set implements types.Store.
func (s *gnovmStore) Set(key, value []byte) {
	if err := s.kvStore.Set(key, value); err != nil {
		s.logger.Error("failed to set value in store", "key", string(key), "error", err)
	}
}

// Delete implements types.Store.
func (s *gnovmStore) Delete(key []byte) {
	if err := s.kvStore.Delete(key); err != nil {
		s.logger.Error("failed to delete key from store", "key", string(key), "error", err)
	}
}

// Iterator implements types.Store.
func (s *gnovmStore) Iterator(start, end []byte) gnostore.Iterator {
	iter, err := s.kvStore.Iterator(start, end)
	if err != nil {
		s.logger.Error("failed to create iterator", "start", string(start), "end", string(end), "error", err)
		return &emptyIterator{}
	}
	return &gnovmIterator{logger: s.logger, iter: iter}
}

// ReverseIterator implements types.Store.
func (s *gnovmStore) ReverseIterator(start, end []byte) gnostore.Iterator {
	iter, err := s.kvStore.ReverseIterator(start, end)
	if err != nil {
		s.logger.Error("failed to create reverse iterator", "start", string(start), "end", string(end), "error", err)
		return &emptyIterator{}
	}
	return &gnovmIterator{logger: s.logger, iter: iter}
}

// CacheWrap implements types.Store.
// Returns a cache-wrapped version of this store.
func (s *gnovmStore) CacheWrap() gnostore.Store {
	// Return the same store instance since the underlying KVStore
	// handles caching and isolation at the Cosmos SDK level
	return s
}

// Write implements types.Store.
// Flushes any cached writes to the underlying store.
func (s *gnovmStore) Write() {
	// No-op as the store service handles writes directly
}

var _ gnostore.Iterator = (*gnovmIterator)(nil)

// gnovmIterator wraps the Cosmos SDK iterator to implement the gno Iterator interface.
type gnovmIterator struct {
	logger log.Logger
	iter   corestore.Iterator
}

// Domain implements types.Iterator.
func (it *gnovmIterator) Domain() (start, end []byte) {
	return it.iter.Domain()
}

// Valid implements types.Iterator.
func (it *gnovmIterator) Valid() bool {
	return it.iter.Valid()
}

// Next implements types.Iterator.
func (it *gnovmIterator) Next() {
	it.iter.Next()
}

// Key implements types.Iterator.
func (it *gnovmIterator) Key() []byte {
	return it.iter.Key()
}

// Value implements types.Iterator.
func (it *gnovmIterator) Value() []byte {
	return it.iter.Value()
}

// Error implements types.Iterator.
func (it *gnovmIterator) Error() error {
	return it.iter.Error()
}

// Close implements types.Iterator.
func (it *gnovmIterator) Close() error {
	return it.iter.Close()
}

var _ gnostore.Iterator = (*emptyIterator)(nil)

// emptyIterator is a no-op iterator returned when iterator creation fails.
type emptyIterator struct{}

func (it *emptyIterator) Domain() (start, end []byte) { return nil, nil }
func (it *emptyIterator) Valid() bool                 { return false }
func (it *emptyIterator) Next()                       {}
func (it *emptyIterator) Key() []byte                 { return nil }
func (it *emptyIterator) Value() []byte               { return nil }
func (it *emptyIterator) Error() error                { return nil }
func (it *emptyIterator) Close() error                { return nil }
