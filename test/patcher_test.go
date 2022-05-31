package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/wamphlett/bezel-project-patcher/pkg/patching"
)

const (
	romDirectoryPath   = "C:\\Games\\N64"
	bezelDirectoryPath = "C:\\Retroarch\\config\\Mupen"
)

func TestCommonCasesWithFuzzyMatching(t *testing.T) {
	tt := map[string]struct {
		romDirContents           []string
		bezelDirContents         []string
		expectedBezelDirContents []string
	}{
		"matching rom": {
			romDirContents:   []string{"AeroFighters Assault (USA).n64"},
			bezelDirContents: []string{"AeroFighters Assault (USA).cfg"},
			expectedBezelDirContents: []string{
				"AeroFighters Assault (USA).cfg",
			},
		},
		"no matching rom": {
			romDirContents:   []string{"Goldeneye 007 (U) [!].n64"},
			bezelDirContents: []string{"AeroFighters Assault (USA).cfg"},
			expectedBezelDirContents: []string{
				"AeroFighters Assault (USA).cfg",
			},
		},
		"multiple mathcing roms": {
			romDirContents: []string{
				"Indy Racing 2000 (U).v64",
				"Indy Racing 2000 (U).n64",
			},
			bezelDirContents: []string{"Indy Racing 2000 (USA).cfg"},
			expectedBezelDirContents: []string{
				"Indy Racing 2000 (USA).cfg",
				"Indy Racing 2000 (U).cfg",
			},
		},
		// Apostrophes
		"rom with apostrophe but not in config": {
			romDirContents:   []string{"AeroFighter's Assault (USA).n64"},
			bezelDirContents: []string{"AeroFighters Assault (USA).cfg"},
			expectedBezelDirContents: []string{
				"AeroFighter's Assault (USA).cfg",
				"AeroFighters Assault (USA).cfg",
			},
		},
		"config with apostrophe but not in rom": {
			romDirContents:   []string{"AeroFighters Assault (USA).n64"},
			bezelDirContents: []string{"AeroFighter's Assault (USA).cfg"},
			expectedBezelDirContents: []string{
				"AeroFighters Assault (USA).cfg",
				"AeroFighter's Assault (USA).cfg",
			},
		},
		// Different Casing
		"matching names with different casing (different tags)": {
			romDirContents:   []string{"aerofighters assault (U) [!].n64"},
			bezelDirContents: []string{"AeroFighters Assault (USA).cfg"},
			expectedBezelDirContents: []string{
				"AeroFighters Assault (USA).cfg",
				"aerofighters assault (U) [!].cfg",
			},
		},
		"matching names with different casing (same tags)": {
			romDirContents:   []string{"aerofighters assault (USA).n64"},
			bezelDirContents: []string{"AeroFighters Assault (USA).cfg"},
			expectedBezelDirContents: []string{
				// only 1 should be present as the file system is not case sensitive
				"AeroFighters Assault (USA).cfg",
			},
		},
		// "The"
		"roms with 'The' moved to the end": {
			romDirContents:   []string{"The New Tetris (USA).n64"},
			bezelDirContents: []string{"New Tetris, The (USA).cfg"},
			expectedBezelDirContents: []string{
				"New Tetris, The (USA).cfg",
				"The New Tetris (USA).cfg",
			},
		},
		"roms with 'The' moved to the end and apostrophes": {
			romDirContents:   []string{"The Tony Hawk's Collection (USA) [!].n64"},
			bezelDirContents: []string{"Tony Hawk's Collection, The (USA).cfg"},
			expectedBezelDirContents: []string{
				"Tony Hawk's Collection, The (USA).cfg",
				"The Tony Hawk's Collection (USA) [!].cfg",
			},
		},
		"roms with 'The' moved to the end with an additional suffix": {
			romDirContents:   []string{"The Addams Family - Pugsley's Scavenger Hunt (USA).zip"},
			bezelDirContents: []string{"Addams Family, The - Pugsley's Scavenger Hunt (USA).cfg"},
			expectedBezelDirContents: []string{
				"Addams Family, The - Pugsley's Scavenger Hunt (USA).cfg",
				"The Addams Family - Pugsley's Scavenger Hunt (USA).cfg",
			},
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			// mock the contents of the directories
			manager := NewStubFileManager()
			manager.SetDirectoryContents(romDirectoryPath, tc.romDirContents)
			manager.SetDirectoryContents(bezelDirectoryPath, tc.bezelDirContents)

			// run the patcher
			patcher := patching.NewPatcher(manager, true)
			patcher.PatchDirectory(bezelDirectoryPath, romDirectoryPath, patching.MatchTypeFuzzy)

			// check the contents of the bezel directory to ensure the expected config files exist
			actualContents, err := manager.GetDirectoryContents(bezelDirectoryPath)
			require.NoError(t, err)

			assert.ElementsMatch(t, tc.expectedBezelDirContents, actualContents)
		})
	}
}

func TestCommonCasesWithAlternateMatching(t *testing.T) {
	tt := map[string]struct {
		romDirContents           []string
		bezelDirContents         []string
		expectedBezelDirContents []string
	}{
		"matching rom": {
			romDirContents:   []string{"AeroFighters Assault (USA).n64"},
			bezelDirContents: []string{"AeroFighters Assault (USA).cfg"},
			expectedBezelDirContents: []string{
				"AeroFighters Assault (USA).cfg",
			},
		},
		"no matching rom": {
			romDirContents:   []string{"Goldeneye 007 (U) [!].n64"},
			bezelDirContents: []string{"AeroFighters Assault (USA).cfg"},
			expectedBezelDirContents: []string{
				"AeroFighters Assault (USA).cfg",
			},
		},
		"multiple mathcing roms": {
			romDirContents: []string{
				"Indy Racing 2000 (U).v64",
				"Indy Racing 2000 (U).n64",
			},
			bezelDirContents: []string{"Indy Racing 2000 (USA).cfg"},
			expectedBezelDirContents: []string{
				"Indy Racing 2000 (USA).cfg",
				"Indy Racing 2000 (U).cfg",
			},
		},
		// Apostrophes
		"rom with apostrophe but not in config": {
			romDirContents:   []string{"AeroFighter's Assault (USA).n64"},
			bezelDirContents: []string{"AeroFighters Assault (USA).cfg"},
			expectedBezelDirContents: []string{
				"AeroFighter's Assault (USA).cfg",
				"AeroFighters Assault (USA).cfg",
			},
		},
		"config with apostrophe but not in rom": {
			romDirContents:   []string{"AeroFighters Assault (USA).n64"},
			bezelDirContents: []string{"AeroFighter's Assault (USA).cfg"},
			expectedBezelDirContents: []string{
				"AeroFighter's Assault (USA).cfg",
			},
		},
		// Different Casing
		"matching names with different casing (different tags)": {
			romDirContents:   []string{"aerofighters assault (U) [!].n64"},
			bezelDirContents: []string{"AeroFighters Assault (USA).cfg"},
			expectedBezelDirContents: []string{
				"AeroFighters Assault (USA).cfg",
				"aerofighters assault (U) [!].cfg",
			},
		},
		"matching names with different casing (same tags)": {
			romDirContents:   []string{"aerofighters assault (USA).n64"},
			bezelDirContents: []string{"AeroFighters Assault (USA).cfg"},
			expectedBezelDirContents: []string{
				// only 1 should be present as the file system is not case sensitive
				"AeroFighters Assault (USA).cfg",
			},
		},
		// "The"
		"roms with 'The' moved to the end": {
			romDirContents:   []string{"The New Tetris (USA).n64"},
			bezelDirContents: []string{"New Tetris, The (USA).cfg"},
			expectedBezelDirContents: []string{
				"New Tetris, The (USA).cfg",
				"The New Tetris (USA).cfg",
			},
		},
		"roms with 'The' moved to the end and apostrophes": {
			romDirContents:   []string{"The Tony Hawk's Collection (USA) [!].n64"},
			bezelDirContents: []string{"Tony Hawk's Collection, The (USA).cfg"},
			expectedBezelDirContents: []string{
				"Tony Hawk's Collection, The (USA).cfg",
				"The Tony Hawk's Collection (USA) [!].cfg",
			},
		},
		"roms with 'The' moved to the end with an additional suffix": {
			romDirContents:   []string{"The Addams Family - Pugsley's Scavenger Hunt (USA).zip"},
			bezelDirContents: []string{"Addams Family, The - Pugsley's Scavenger Hunt (USA).cfg"},
			expectedBezelDirContents: []string{
				"Addams Family, The - Pugsley's Scavenger Hunt (USA).cfg",
				"The Addams Family - Pugsley's Scavenger Hunt (USA).cfg",
			},
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			// mock the contents of the directories
			manager := NewStubFileManager()
			manager.SetDirectoryContents(romDirectoryPath, tc.romDirContents)
			manager.SetDirectoryContents(bezelDirectoryPath, tc.bezelDirContents)

			// run the patcher
			patcher := patching.NewPatcher(manager, true)
			patcher.PatchDirectory(bezelDirectoryPath, romDirectoryPath, patching.MatchTypeAlternate)

			// check the contents of the bezel directory to ensure the expected config files exist
			actualContents, err := manager.GetDirectoryContents(bezelDirectoryPath)
			require.NoError(t, err)

			assert.ElementsMatch(t, tc.expectedBezelDirContents, actualContents)
		})
	}
}

func TestCommonCasesWithExactMatching(t *testing.T) {
	tt := map[string]struct {
		romDirContents           []string
		bezelDirContents         []string
		expectedBezelDirContents []string
	}{
		"matching rom": {
			romDirContents:   []string{"AeroFighters Assault (USA).n64"},
			bezelDirContents: []string{"AeroFighters Assault (USA).cfg"},
			expectedBezelDirContents: []string{
				"AeroFighters Assault (USA).cfg",
			},
		},
		"no matching rom": {
			romDirContents:   []string{"Goldeneye 007 (U) [!].n64"},
			bezelDirContents: []string{"AeroFighters Assault (USA).cfg"},
			expectedBezelDirContents: []string{
				"AeroFighters Assault (USA).cfg",
			},
		},
		"multiple mathcing roms": {
			romDirContents: []string{
				"Indy Racing 2000 (U).v64",
				"Indy Racing 2000 (U).n64",
			},
			bezelDirContents: []string{"Indy Racing 2000 (USA).cfg"},
			expectedBezelDirContents: []string{
				"Indy Racing 2000 (USA).cfg",
				"Indy Racing 2000 (U).cfg",
			},
		},
		// Apostrophes
		"rom with apostrophe but not in config": {
			romDirContents:   []string{"AeroFighter's Assault (USA).n64"},
			bezelDirContents: []string{"AeroFighters Assault (USA).cfg"},
			expectedBezelDirContents: []string{
				"AeroFighters Assault (USA).cfg",
			},
		},
		"config with apostrophe but not in rom": {
			romDirContents:   []string{"AeroFighters Assault (USA).n64"},
			bezelDirContents: []string{"AeroFighter's Assault (USA).cfg"},
			expectedBezelDirContents: []string{
				"AeroFighter's Assault (USA).cfg",
			},
		},
		// Different Casing
		"matching names with different casing (different tags)": {
			romDirContents:   []string{"aerofighters assault (U) [!].n64"},
			bezelDirContents: []string{"AeroFighters Assault (USA).cfg"},
			expectedBezelDirContents: []string{
				"AeroFighters Assault (USA).cfg",
				"aerofighters assault (U) [!].cfg",
			},
		},
		"matching names with different casing (same tags)": {
			romDirContents:   []string{"aerofighters assault (USA).n64"},
			bezelDirContents: []string{"AeroFighters Assault (USA).cfg"},
			expectedBezelDirContents: []string{
				// only 1 should be present as the file system is not case sensitive
				"AeroFighters Assault (USA).cfg",
			},
		},
		// "The"
		"roms with 'The' moved to the end": {
			romDirContents:   []string{"The New Tetris (USA).n64"},
			bezelDirContents: []string{"New Tetris, The (USA).cfg"},
			expectedBezelDirContents: []string{
				"New Tetris, The (USA).cfg",
			},
		},
		"roms with 'The' moved to the end and apostrophes": {
			romDirContents:   []string{"The Tony Hawk's Collection (USA) [!].n64"},
			bezelDirContents: []string{"Tony Hawk's Collection, The (USA).cfg"},
			expectedBezelDirContents: []string{
				"Tony Hawk's Collection, The (USA).cfg",
			},
		},
		"roms with 'The' moved to the end with an additional suffix": {
			romDirContents:   []string{"The Addams Family - Pugsley's Scavenger Hunt (USA).zip"},
			bezelDirContents: []string{"Addams Family, The - Pugsley's Scavenger Hunt (USA).cfg"},
			expectedBezelDirContents: []string{
				"Addams Family, The - Pugsley's Scavenger Hunt (USA).cfg",
			},
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			// mock the contents of the directories
			manager := NewStubFileManager()
			manager.SetDirectoryContents(romDirectoryPath, tc.romDirContents)
			manager.SetDirectoryContents(bezelDirectoryPath, tc.bezelDirContents)

			// run the patcher
			patcher := patching.NewPatcher(manager, true)
			patcher.PatchDirectory(bezelDirectoryPath, romDirectoryPath, patching.MatchTypeExact)

			// check the contents of the bezel directory to ensure the expected config files exist
			actualContents, err := manager.GetDirectoryContents(bezelDirectoryPath)
			require.NoError(t, err)

			assert.ElementsMatch(t, tc.expectedBezelDirContents, actualContents)
		})
	}
}

func TestDirectoryWhichHasAlreadyBeenMapped(t *testing.T) {
	// mock the contents of the directories
	manager := NewStubFileManager()
	manager.SetDirectoryContents(romDirectoryPath, []string{"The New Tetris (USA).n64"})
	manager.SetDirectoryContents(bezelDirectoryPath, []string{"The New Tetris (USA).cfg", "The New Tetris (U).cfg"})

	// run the patcher
	patcher := patching.NewPatcher(manager, true)
	patcher.PatchDirectory(bezelDirectoryPath, romDirectoryPath, patching.MatchTypeFuzzy)

	// check the contents of the bezel directory to ensure the expected config files exist
	actualContents, err := manager.GetDirectoryContents(bezelDirectoryPath)
	require.NoError(t, err)

	assert.ElementsMatch(t, []string{"The New Tetris (USA).cfg", "The New Tetris (U).cfg"}, actualContents)
}

func TestDirectoryWhichHasTheSameRomWithMultipleExtensions(t *testing.T) {
	// mock the contents of the directories
	manager := NewStubFileManager()
	manager.SetDirectoryContents(romDirectoryPath, []string{"The New Tetris (USA).n64", "The New Tetris (USA).z64"})
	manager.SetDirectoryContents(bezelDirectoryPath, []string{"The New Tetris (U).cfg"})

	// run the patcher
	patcher := patching.NewPatcher(manager, true)
	patcher.PatchDirectory(bezelDirectoryPath, romDirectoryPath, patching.MatchTypeFuzzy)

	// check the contents of the bezel directory to ensure the expected config files exist
	actualContents, err := manager.GetDirectoryContents(bezelDirectoryPath)
	require.NoError(t, err)

	assert.ElementsMatch(t, []string{"The New Tetris (USA).cfg", "The New Tetris (U).cfg"}, actualContents)
}
