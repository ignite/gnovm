package types

import "github.com/gnolang/gno/tm2/pkg/std"

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	if err := gs.Params.Validate(); err != nil {
		return err
	}
	for _, p := range gs.Packages {
		if err := p.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func (p Package) Validate() error {
	return p.ToMemPackage().ValidateBasic()
}

func (p Package) ToMemPackage() *std.MemPackage {
	mp := &std.MemPackage{
		Name: p.Name,
		Path: p.Path,
	}
	for _, f := range p.Files {
		mp.Files = append(mp.Files, &std.MemFile{
			Name: f.Name,
			Body: f.Body,
		})
	}
	return mp
}
