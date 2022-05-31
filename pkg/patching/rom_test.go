package patching

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRomParse(t *testing.T) {
	tt := map[string]struct {
		fileName    string
		expectedRom *Rom
	}{
		"standard rom name": {
			fileName: "Tony Hawk's Collection, The (USA) [!].n64",
			expectedRom: &Rom{
				FileName: "Tony Hawk's Collection, The (USA) [!].n64",
				Name:     "tony hawk's collection, the",
				AlternateNames: []string{
					"the tony hawk's collection",
					"the tony hawks collection",

					"tony hawk's collection, the",
					"tony hawks collection, the",
				},
			},
		},
		"no-intro rom name": {
			fileName: "Tony Hawk's Collection, The (USA) [!].n64",
			expectedRom: &Rom{
				FileName: "Tony Hawk's Collection, The (USA) [!].n64",
				Name:     "tony hawk's collection, the",
				AlternateNames: []string{
					"the tony hawk's collection",
					"the tony hawks collection",

					"tony hawk's collection, the",
					"tony hawks collection, the",
				},
			},
		},
		"standard rom name with suffix": {
			fileName: "The Tony Hawk's Collection - Ultimate Edition (USA) [!].n64",
			expectedRom: &Rom{
				FileName: "The Tony Hawk's Collection - Ultimate Edition (USA) [!].n64",
				Name:     "the tony hawk's collection - ultimate edition",
				AlternateNames: []string{
					"the tony hawk's collection - ultimate edition",
					"the tony hawks collection - ultimate edition",

					"tony hawk's collection, the - ultimate edition",
					"tony hawks collection, the - ultimate edition",
				},
			},
		},
		"no-intro rom name with suffix": {
			fileName: "Tony Hawk's Collection, The - Ultimate Edition (USA) [!].n64",
			expectedRom: &Rom{
				FileName: "Tony Hawk's Collection, The - Ultimate Edition (USA) [!].n64",
				Name:     "tony hawk's collection, the - ultimate edition",
				AlternateNames: []string{
					"the tony hawk's collection - ultimate edition",
					"the tony hawks collection - ultimate edition",

					"tony hawk's collection, the - ultimate edition",
					"tony hawks collection, the - ultimate edition",
				},
			},
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			rom := NewRom(tc.fileName)
			assert.Equal(t, tc.expectedRom.FileName, rom.FileName)
			assert.Equal(t, tc.expectedRom.Name, rom.Name)
			assert.ElementsMatch(t, tc.expectedRom.AlternateNames, rom.AlternateNames)
		})
	}

}

func TestSameRomWithDifferntNameFormatMatchAlternates(t *testing.T) {
	rom1 := NewRom("The Tony Hawk's Collection (USA) [!].n64")
	rom2 := NewRom("Tony Hawk's Collection, The (USA) [!].n64")
	assert.ElementsMatch(t, rom1.AlternateNames, rom2.AlternateNames)
}
