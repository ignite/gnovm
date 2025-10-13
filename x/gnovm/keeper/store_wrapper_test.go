package keeper

import (
	"context"
	"testing"

	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"github.com/stretchr/testify/require"
	"gotest.tools/v3/assert"

	"log/slog"

	bft "github.com/gnolang/gno/tm2/pkg/bft/types"
	gnosdk "github.com/gnolang/gno/tm2/pkg/sdk"
	gnostore "github.com/gnolang/gno/tm2/pkg/store"
	"github.com/gnolang/gno/tm2/pkg/store/types"
)

// mockKVStore implements store.KVStore for testing
type mockKVStore struct {
	data map[string][]byte
}

func newMockKVStore() *mockKVStore {
	return &mockKVStore{
		data: make(map[string][]byte),
	}
}

func (m *mockKVStore) Get(key []byte) ([]byte, error) {
	return m.data[string(key)], nil
}

func (m *mockKVStore) Has(key []byte) (bool, error) {
	_, exists := m.data[string(key)]
	return exists, nil
}

func (m *mockKVStore) Set(key, value []byte) error {
	m.data[string(key)] = value
	return nil
}

func (m *mockKVStore) Delete(key []byte) error {
	delete(m.data, string(key))
	return nil
}

func (m *mockKVStore) Iterator(start, end []byte) (store.Iterator, error) {
	return newMockIterator(m.data, start, end, false), nil
}

func (m *mockKVStore) ReverseIterator(start, end []byte) (store.Iterator, error) {
	return newMockIterator(m.data, start, end, true), nil
}

// mockIterator implements store.Iterator for testing
type mockIterator struct {
	data    map[string][]byte
	keys    []string
	index   int
	start   []byte
	end     []byte
	reverse bool
}

func newMockIterator(data map[string][]byte, start, end []byte, reverse bool) *mockIterator {
	keys := make([]string, 0, len(data))
	for k := range data {
		// Simple key filtering for test purposes
		if start != nil && string(k) < string(start) {
			continue
		}
		if end != nil && string(k) >= string(end) {
			continue
		}
		keys = append(keys, k)
	}

	if reverse {
		// Reverse the keys for reverse iteration
		for i, j := 0, len(keys)-1; i < j; i, j = i+1, j-1 {
			keys[i], keys[j] = keys[j], keys[i]
		}
	}

	return &mockIterator{
		data:    data,
		keys:    keys,
		index:   -1,
		start:   start,
		end:     end,
		reverse: reverse,
	}
}

func (m *mockIterator) Domain() (start, end []byte) {
	return m.start, m.end
}

func (m *mockIterator) Valid() bool {
	return m.index >= 0 && m.index < len(m.keys)
}

func (m *mockIterator) Next() {
	m.index++
}

func (m *mockIterator) Key() []byte {
	if !m.Valid() {
		return nil
	}
	return []byte(m.keys[m.index])
}

func (m *mockIterator) Value() []byte {
	if !m.Valid() {
		return nil
	}
	return m.data[m.keys[m.index]]
}

func (m *mockIterator) Error() error {
	return nil
}

func (m *mockIterator) Close() error {
	return nil
}

// mockStoreService implements store.KVStoreService for testing
type mockStoreService struct {
	store *mockKVStore
}

func newMockStoreService() *mockStoreService {
	return &mockStoreService{
		store: newMockKVStore(),
	}
}

func (m *mockStoreService) OpenKVStore(ctx context.Context) store.KVStore {
	return m.store
}

func (m *mockStoreService) OpenMemoryStore(ctx context.Context) store.KVStore {
	return m.store
}

func TestNewGnovmMultiStore(t *testing.T) {
	logger := log.NewNopLogger()
	storeService := newMockStoreService()
	memStoreService := newMockStoreService()
	storeKey := gnostore.NewStoreKey("test")
	memStoreKey := gnostore.NewStoreKey("test-mem")
	gnoCtx := gnosdk.NewContext(gnosdk.RunTxModeCheck, nil, &bft.Header{ChainID: "test-chain"}, slog.Default())

	multiStore := NewGnovmMultiStore(logger, storeService, memStoreService, storeKey, memStoreKey, gnoCtx, context.Background())

	require.NotNil(t, multiStore)
	assert.Assert(t, multiStore != nil)
}

func TestGnovmMultiStore_GetStore(t *testing.T) {
	logger := log.NewNopLogger()
	storeService := newMockStoreService()
	memStoreService := newMockStoreService()
	storeKey := gnostore.NewStoreKey("test")
	memStoreKey := gnostore.NewStoreKey("test-mem")
	gnoCtx := gnosdk.NewContext(gnosdk.RunTxModeCheck, nil, &bft.Header{ChainID: "test-chain"}, slog.Default())

	multiStore := NewGnovmMultiStore(logger, storeService, memStoreService, storeKey, memStoreKey, gnoCtx, context.Background())

	t.Run("valid store key", func(t *testing.T) {
		store := multiStore.GetStore(storeKey)
		require.NotNil(t, store)
	})

	t.Run("invalid store key panics", func(t *testing.T) {
		invalidKey := gnostore.NewStoreKey("invalid")

		require.Panics(t, func() {
			multiStore.GetStore(invalidKey)
		})
	})
}

func TestGnovmMultiStore_MultiCacheWrap(t *testing.T) {
	logger := log.NewNopLogger()
	storeService := newMockStoreService()
	memStoreService := newMockStoreService()
	storeKey := gnostore.NewStoreKey("test")
	memStoreKey := gnostore.NewStoreKey("test-mem")
	gnoCtx := gnosdk.NewContext(gnosdk.RunTxModeCheck, nil, &bft.Header{ChainID: "test-chain"}, slog.Default())

	multiStore := NewGnovmMultiStore(logger, storeService, memStoreService, storeKey, memStoreKey, gnoCtx, context.Background())

	cachedStore := multiStore.MultiCacheWrap()
	require.NotNil(t, cachedStore)

	// For our simple implementation, it should return the same instance
	assert.Equal(t, multiStore, cachedStore)
}

func TestGnovmMultiStore_MultiWrite(t *testing.T) {
	logger := log.NewNopLogger()
	storeService := newMockStoreService()
	memStoreService := newMockStoreService()
	storeKey := gnostore.NewStoreKey("test")
	memStoreKey := gnostore.NewStoreKey("test-mem")
	gnoCtx := gnosdk.NewContext(gnosdk.RunTxModeCheck, nil, &bft.Header{ChainID: "test-chain"}, slog.Default())

	multiStore := NewGnovmMultiStore(logger, storeService, memStoreService, storeKey, memStoreKey, gnoCtx, context.Background())

	// Should not panic
	require.NotPanics(t, func() {
		multiStore.MultiWrite()
	})
}

func TestGnovmStore_BasicOperations(t *testing.T) {
	logger := log.NewNopLogger()
	storeService := newMockStoreService()
	memStoreService := newMockStoreService()
	storeKey := gnostore.NewStoreKey("test")
	memStoreKey := gnostore.NewStoreKey("test-mem")
	gnoCtx := gnosdk.NewContext(gnosdk.RunTxModeCheck, nil, &bft.Header{ChainID: "test-chain"}, slog.Default())

	multiStore := NewGnovmMultiStore(logger, storeService, memStoreService, storeKey, memStoreKey, gnoCtx, context.Background())
	store := multiStore.GetStore(storeKey)

	t.Run("set and get", func(t *testing.T) {
		key := []byte("testkey")
		value := []byte("testvalue")

		// Initially should not exist
		assert.Assert(t, !store.Has(key))
		assert.Assert(t, store.Get(key) == nil)

		// Set the value
		store.Set(key, value)

		// Should now exist and return correct value
		assert.Assert(t, store.Has(key))
		assert.DeepEqual(t, value, store.Get(key))
	})

	t.Run("delete", func(t *testing.T) {
		key := []byte("deletekey")
		value := []byte("deletevalue")

		// Set and verify
		store.Set(key, value)
		assert.Assert(t, store.Has(key))

		// Delete and verify
		store.Delete(key)
		assert.Assert(t, !store.Has(key))
		assert.Assert(t, store.Get(key) == nil)
	})
}

func TestGnovmStore_Iterator(t *testing.T) {
	logger := log.NewNopLogger()
	storeService := newMockStoreService()
	memStoreService := newMockStoreService()
	storeKey := gnostore.NewStoreKey("test")
	memStoreKey := gnostore.NewStoreKey("test-mem")
	gnoCtx := gnosdk.NewContext(gnosdk.RunTxModeCheck, nil, &bft.Header{ChainID: "test-chain"}, slog.Default())

	multiStore := NewGnovmMultiStore(logger, storeService, memStoreService, storeKey, memStoreKey, gnoCtx, context.Background())
	store := multiStore.GetStore(storeKey)

	// Set up test data
	testData := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}

	for k, v := range testData {
		store.Set([]byte(k), []byte(v))
	}

	t.Run("forward iterator", func(t *testing.T) {
		iter := store.Iterator(nil, nil)
		defer iter.Close()

		count := 0
		// Start iteration
		iter.Next()
		for iter.Valid() {
			key := iter.Key()
			value := iter.Value()

			require.NotNil(t, key)
			require.NotNil(t, value)

			expectedValue, exists := testData[string(key)]
			require.True(t, exists)
			assert.DeepEqual(t, []byte(expectedValue), value)

			count++
			iter.Next()
		}

		assert.Equal(t, len(testData), count)
	})

	t.Run("reverse iterator", func(t *testing.T) {
		iter := store.ReverseIterator(nil, nil)
		defer iter.Close()

		count := 0
		// Start iteration
		iter.Next()
		for iter.Valid() {
			key := iter.Key()
			value := iter.Value()

			require.NotNil(t, key)
			require.NotNil(t, value)

			expectedValue, exists := testData[string(key)]
			require.True(t, exists)
			assert.DeepEqual(t, []byte(expectedValue), value)

			count++
			iter.Next()
		}

		assert.Equal(t, len(testData), count)
	})
}

func TestGnovmStore_CacheWrap(t *testing.T) {
	logger := log.NewNopLogger()
	storeService := newMockStoreService()
	memStoreService := newMockStoreService()
	storeKey := gnostore.NewStoreKey("test")
	memStoreKey := gnostore.NewStoreKey("test-mem")
	gnoCtx := gnosdk.NewContext(gnosdk.RunTxModeCheck, nil, &bft.Header{ChainID: "test-chain"}, slog.Default())

	multiStore := NewGnovmMultiStore(logger, storeService, memStoreService, storeKey, memStoreKey, gnoCtx, context.Background())
	store := multiStore.GetStore(storeKey)

	cachedStore := store.CacheWrap()
	require.NotNil(t, cachedStore)

	// For our simple implementation, it should return the same instance
	assert.Equal(t, store, cachedStore)
}

func TestGnovmStore_Write(t *testing.T) {
	logger := log.NewNopLogger()
	storeService := newMockStoreService()
	memStoreService := newMockStoreService()
	storeKey := gnostore.NewStoreKey("test")
	memStoreKey := gnostore.NewStoreKey("test-mem")
	gnoCtx := gnosdk.NewContext(gnosdk.RunTxModeCheck, nil, &bft.Header{ChainID: "test-chain"}, slog.Default())

	multiStore := NewGnovmMultiStore(logger, storeService, memStoreService, storeKey, memStoreKey, gnoCtx, context.Background())
	store := multiStore.GetStore(storeKey)

	// Should not panic
	require.NotPanics(t, func() {
		store.Write()
	})
}

func TestKeeper_NewMultiStore(t *testing.T) {
	// This test would require setting up a full keeper, which is complex
	// For now, we test the core functionality through the direct constructor
	logger := log.NewNopLogger()
	storeService := newMockStoreService()
	memStoreService := newMockStoreService()
	storeKey := gnostore.NewStoreKey("test")
	memStoreKey := gnostore.NewStoreKey("test-mem")
	gnoCtx := gnosdk.NewContext(gnosdk.RunTxModeCheck, nil, &bft.Header{ChainID: "test-chain"}, slog.Default())

	multiStore := NewGnovmMultiStore(logger, storeService, memStoreService, storeKey, memStoreKey, gnoCtx, context.Background())

	// Verify it implements the MultiStore interface
	var _ types.MultiStore = multiStore

	// Verify basic functionality
	store := multiStore.GetStore(storeKey)
	require.NotNil(t, store)

	// Verify it implements the Store interface
	var _ types.Store = store
}
