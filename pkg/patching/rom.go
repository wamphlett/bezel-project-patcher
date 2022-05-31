package patching

import (
	"path/filepath"
	"strings"
)

// Rom is used to hold information about a file
type Rom struct {
	FileName       string
	Name           string
	AlternateNames []string
}

// NewRom builds a new ROM and works out all the alternate names
func NewRom(fileName string) *Rom {
	baseName := getBaseName(fileName)
	return &Rom{
		FileName:       fileName,
		Name:           baseName,
		AlternateNames: getAlternateNames(baseName, true),
	}
}

// ConfigName returns the config name that should be used
func (r *Rom) ConfigName() string {
	return strings.TrimSuffix(r.FileName, filepath.Ext(r.FileName)) + ".cfg"
}

// getBaseName returns the name of the file minus any extensions and tags i.e. (U), [!] etc
func getBaseName(fileName string) string {
	// strip any file extensions or tags {
	for _, substr := range []string{".", "[", "("} {
		if i := strings.Index(fileName, substr); i > 0 {
			fileName = fileName[:i]
		}
	}

	return strings.ToLower(strings.TrimSpace(fileName))
}

// getAlternateNames works out all the possible alternate names for the given name
func getAlternateNames(name string, recursive bool) []string {
	name = strings.ToLower(name)
	alternates := []string{name}

	// the no-intro ROM naming convention moves the "the" to the end of the ROM name but
	// before any suffix's. e.g. "the simpsons - ultimate" would become "simpsons, the - ultimate"
	if strings.HasPrefix(name, "the ") || strings.Contains(name, ", the") {
		var correctedName string
		if strings.HasPrefix(name, "the ") {
			// some roms have a suffix and the "the" should be placed before this suffix
			nameParts := strings.Split(name[4:], " - ")
			nameParts[0] += ", the"
			correctedName = strings.Join(nameParts, " - ")
		} else {
			i := strings.Index(name, ", the")
			correctedName = "the " + name[:i] + name[i+5:]
		}

		alternates = append(alternates, correctedName)
		if recursive {
			alternates = append(alternates, getAlternateNames(correctedName, false)...)
		}
	}

	// get the name with any apostrophes removed
	if strings.Index(name, "'") >= 0 {
		alternates = append(alternates, strings.ReplaceAll(name, "'", ""))
	}

	return uniqueItems(alternates)
}

// uniqueItems returns only the unique items from a slice
func uniqueItems(slice []string) []string {
	uniqueSlice := []string{}
	for _, item := range slice {
		if !containsItem(uniqueSlice, item) {
			uniqueSlice = append(uniqueSlice, item)
		}
	}
	return uniqueSlice
}

// containsItem returns true if the string is present in the given slice
func containsItem(haystack []string, needle string) bool {
	for _, i := range haystack {
		if i == needle {
			return true
		}
	}
	return false
}
