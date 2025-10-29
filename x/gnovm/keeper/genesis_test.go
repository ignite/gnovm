package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ignite/gnovm/x/gnovm/types"

	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),
	}

	f := initFixture(t)
	err := f.keeper.InitGenesis(f.ctx, genesisState)
	require.NoError(t, err)
	got, err := f.keeper.ExportGenesis(f.ctx)
	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, genesisState.Params, got.Params)
}

func TestGenesisStateExport(t *testing.T) {
	tests := []struct {
		name        string
		setupState  func(f *fixture)
		validateGen func(t *testing.T, gen *types.GenesisState)
	}{
		{
			name: "empty state exports successfully",
			setupState: func(f *fixture) {
				// no setup needed
			},
			validateGen: func(t *testing.T, gen *types.GenesisState) {
				require.NotNil(t, gen)
				require.NotNil(t, gen.State)
			},
		},
		{
			name: "state with data exports all key-value pairs",
			setupState: func(f *fixture) {
				// Initialize genesis first
				err := f.keeper.InitGenesis(f.ctx, types.GenesisState{
					Params: types.DefaultParams(),
				})
				require.NoError(t, err)

				// Add some test data to the store
				sdkCtx := f.ctx.(sdk.Context)
				store := f.storeService.OpenKVStore(sdkCtx)
				testData := map[string]string{
					"test/key1": "value1",
					"test/key2": "value2",
					"test/key3": "value3",
				}
				for k, v := range testData {
					err := store.Set([]byte(k), []byte(v))
					require.NoError(t, err)
				}
			},
			validateGen: func(t *testing.T, gen *types.GenesisState) {
				require.NotNil(t, gen)
				require.NotNil(t, gen.State)
				require.NotEmpty(t, gen.State)

				// Verify test data is present
				found := make(map[string]bool)
				for _, kv := range gen.State {
					key := string(kv.Key)
					if key == "test/key1" || key == "test/key2" || key == "test/key3" {
						found[key] = true
					}
				}
				require.True(t, found["test/key1"], "test/key1 should be exported")
				require.True(t, found["test/key2"], "test/key2 should be exported")
				require.True(t, found["test/key3"], "test/key3 should be exported")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := initFixture(t)
			tt.setupState(f)

			gen, err := f.keeper.ExportGenesis(f.ctx)
			require.NoError(t, err)
			tt.validateGen(t, gen)
		})
	}
}

func TestGenesisStateImport(t *testing.T) {
	tests := []struct {
		name      string
		genState  types.GenesisState
		validate  func(t *testing.T, f *fixture)
		expectErr bool
	}{
		{
			name: "import empty state",
			genState: types.GenesisState{
				Params: types.DefaultParams(),
				State:  []types.KVPair{},
			},
			validate: func(t *testing.T, f *fixture) {
				// Verify initialization succeeded
				sdkCtx := f.ctx.(sdk.Context)
				store := f.storeService.OpenKVStore(sdkCtx)
				iter, err := store.Iterator(nil, nil)
				require.NoError(t, err)
				defer iter.Close()

				// Should have some data from initialization
				require.True(t, iter.Valid())
			},
			expectErr: false,
		},
		{
			name: "import state with key-value pairs",
			genState: types.GenesisState{
				Params: types.DefaultParams(),
				State: []types.KVPair{
					{Key: []byte("imported/key1"), Value: []byte("imported_value1")},
					{Key: []byte("imported/key2"), Value: []byte("imported_value2")},
					{Key: []byte("imported/key3"), Value: []byte("imported_value3")},
				},
			},
			validate: func(t *testing.T, f *fixture) {
				sdkCtx := f.ctx.(sdk.Context)
				store := f.storeService.OpenKVStore(sdkCtx)

				// Verify imported data
				val1, err := store.Get([]byte("imported/key1"))
				require.NoError(t, err)
				require.Equal(t, []byte("imported_value1"), val1)

				val2, err := store.Get([]byte("imported/key2"))
				require.NoError(t, err)
				require.Equal(t, []byte("imported_value2"), val2)

				val3, err := store.Get([]byte("imported/key3"))
				require.NoError(t, err)
				require.Equal(t, []byte("imported_value3"), val3)
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := initFixture(t)

			err := f.keeper.InitGenesis(f.ctx, tt.genState)
			if tt.expectErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			if tt.validate != nil {
				tt.validate(t, f)
			}
		})
	}
}

func TestGenesisRoundTrip(t *testing.T) {
	f := initFixture(t)

	// Initialize with default genesis
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),
		State: []types.KVPair{
			{Key: []byte("roundtrip/key1"), Value: []byte("roundtrip_value1")},
			{Key: []byte("roundtrip/key2"), Value: []byte("roundtrip_value2")},
		},
	}

	err := f.keeper.InitGenesis(f.ctx, genesisState)
	require.NoError(t, err)

	// Add additional data after initialization
	sdkCtx := f.ctx.(sdk.Context)
	store := f.storeService.OpenKVStore(sdkCtx)
	err = store.Set([]byte("additional/key"), []byte("additional_value"))
	require.NoError(t, err)

	// Export genesis
	exported, err := f.keeper.ExportGenesis(f.ctx)
	require.NoError(t, err)
	require.NotNil(t, exported)

	// Verify all keys are present in exported state
	keyMap := make(map[string][]byte)
	for _, kv := range exported.State {
		keyMap[string(kv.Key)] = kv.Value
	}

	// Check original imported keys
	require.Contains(t, keyMap, "roundtrip/key1")
	require.Equal(t, []byte("roundtrip_value1"), keyMap["roundtrip/key1"])
	require.Contains(t, keyMap, "roundtrip/key2")
	require.Equal(t, []byte("roundtrip_value2"), keyMap["roundtrip/key2"])

	// Check additional key
	require.Contains(t, keyMap, "additional/key")
	require.Equal(t, []byte("additional_value"), keyMap["additional/key"])

	// Create a new fixture and import the exported state
	f2 := initFixture(t)
	err = f2.keeper.InitGenesis(f2.ctx, *exported)
	require.NoError(t, err)

	// Verify data in the new fixture
	sdkCtx2 := f2.ctx.(sdk.Context)
	store2 := f2.storeService.OpenKVStore(sdkCtx2)
	val1, err := store2.Get([]byte("roundtrip/key1"))
	require.NoError(t, err)
	require.Equal(t, []byte("roundtrip_value1"), val1)

	val2, err := store2.Get([]byte("roundtrip/key2"))
	require.NoError(t, err)
	require.Equal(t, []byte("roundtrip_value2"), val2)

	valAdditional, err := store2.Get([]byte("additional/key"))
	require.NoError(t, err)
	require.Equal(t, []byte("additional_value"), valAdditional)
}
