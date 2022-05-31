package patching

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// FileMangerInterface defines the methods required to interact with the file system
type FileMangerInterface interface {
	GetDirectoryContents(directoryPath string) ([]string, error)
	CopyFileWithName(directoryPath, filePath, newName string) error
	FileExists(directoryPath, fileName string) bool
}

// matchType is used to identify what match type was used to match 2 file names
type matchType string

const (
	// MatchTypeExact means the ROM and config names (minus any tags i.e. (U), [!]) matched exactly
	MatchTypeExact matchType = "exact"
	// MatchTypeAlternate means a ROM's alternate name matched the config name
	MatchTypeAlternate matchType = "alternate"
	// MatchTypeFuzzy means a ROM's alternate name matched one of the configs alternate names
	MatchTypeFuzzy matchType = "fuzzy"
	// MatchTypeNone means no match was found
	MatchTypeNone matchType = "none"
)

// match records a match between a config file and a ROM, what
// type of match it was and whether it was existing or not i.e the config for the
// ROM already existed.
type match struct {
	configFile *Rom
	rom        *Rom
	matchType  matchType
	isExisting bool
}

// Patcher defines the dependencies in order to success patch a directory
type Patcher struct {
	fileManager FileMangerInterface
	commit      bool
}

// NewPatcher returns a new Patcher with the required dependencies
func NewPatcher(fileManager FileMangerInterface, commit bool) *Patcher {
	return &Patcher{
		fileManager: fileManager,
		commit:      commit,
	}
}

// PatchDirectory patches the given config directory with the ROMs in the given ROM directory.
// New files will only be created if the match type matches the match flag.
func (p *Patcher) PatchDirectory(configDirPath, romDirPath string, matchFlag matchType) error {
	// get a list of files from the config directory and the rom directory
	configDirFiles, err := p.fileManager.GetDirectoryContents(configDirPath)
	if err != nil {
		return err
	}
	romDirFiles, err := p.fileManager.GetDirectoryContents(romDirPath)
	if err != nil {
		return err
	}

	// filter out anything which does not look like a config file
	filteredConfigDirFiles := []string{}
	for _, file := range configDirFiles {
		if filepath.Ext(file) == ".cfg" {
			filteredConfigDirFiles = append(filteredConfigDirFiles, file)
		}
	}

	configFiles := make([]*Rom, len(filteredConfigDirFiles))
	for i, item := range filteredConfigDirFiles {
		configFiles[i] = NewRom(item)
	}

	roms := make([]*Rom, len(romDirFiles))
	for i, item := range romDirFiles {
		roms[i] = NewRom(item)
	}

	matches := p.matchRomSets(configFiles, roms)

	for _, match := range matches {
		if match.matchType != MatchTypeNone {
			if !match.isExisting && shouldInclude(match.matchType, matchFlag) {
				// make sure the file has not already been added (might have been added by
				// a previous match so we have to check)
				if p.fileManager.FileExists(configDirPath, match.rom.ConfigName()) {
					match.isExisting = true
					continue
				}
				// only do the file operations if --commit was specified. this gives the
				// users a chance to sanity check the log before changing any of their files
				if p.commit {
					p.fileManager.CopyFileWithName(configDirPath, match.configFile.FileName, match.rom.ConfigName())
				}
			}
		}
	}

	p.produceLog(len(roms), len(configFiles), romDirPath, configDirPath, matches, matchFlag)

	return nil
}

// matchRomSets will attempt to match a config file to one of the ROMs preferring exact matches,
// followed by alternate matches, then fuzzy matches.
func (p *Patcher) matchRomSets(configFiles, romSet []*Rom) []*match {
	matches := []*match{}
	for _, configFile := range configFiles {
	matchLoop:
		for _, rom := range romSet {
			// exactly match roms
			if configFile.Name == rom.Name {
				matches = append(matches, &match{
					configFile: configFile,
					rom:        rom,
					matchType:  MatchTypeExact,
					// check if the file names fully match (case-insensitive). This indicates that this match
					// is an existing match
					isExisting: strings.ToLower(configFile.FileName) == strings.ToLower(rom.ConfigName()),
				})
				continue matchLoop
			}

			// try to match rom on alternate name
			for _, romAlternateName := range rom.AlternateNames {
				if configFile.Name == romAlternateName {
					matches = append(matches, &match{
						configFile: configFile,
						rom:        rom,
						matchType:  MatchTypeAlternate,
					})
					continue matchLoop
				}
			}

			// if we still haven't matched anything, attempt to do a fuzzy match
			for _, romAlternateName := range rom.AlternateNames {
				for _, configAlternateName := range configFile.AlternateNames {
					if configAlternateName == romAlternateName {
						matches = append(matches, &match{
							configFile: configFile,
							rom:        rom,
							matchType:  MatchTypeFuzzy,
						})
						continue matchLoop
					}
				}
			}
		}
	}

	// record a no match for any ROMs which did not get matched
romSetLoop:
	for _, rom := range romSet {
		for _, match := range matches {
			if match.rom == rom && match.matchType != MatchTypeNone {
				continue romSetLoop
			}
		}
		matches = append(matches, &match{
			rom:       rom,
			matchType: MatchTypeNone,
		})
	}

	// record a no match for any config files which did not get matched
configSetLoop:
	for _, config := range configFiles {
		for _, match := range matches {
			if match.configFile == config && match.matchType != MatchTypeNone {
				continue configSetLoop
			}
		}
		matches = append(matches, &match{
			configFile: config,
			matchType:  MatchTypeNone,
		})
	}

	return matches
}

// produceLog write a log file to config directory to give a detailed description of what the patching did
func (p *Patcher) produceLog(romCount, configCount int, romDirPath, configPath string, matches []*match, matchFlag matchType) {
	romsWithoutConfig := []string{}
	configWithoutRoms := []string{}
	createdFiles := map[matchType][]string{}
	for _, t := range []matchType{MatchTypeExact, MatchTypeAlternate, MatchTypeFuzzy} {
		createdFiles[t] = []string{}
	}

	for _, match := range matches {
		if match.matchType == MatchTypeNone {
			if match.configFile != nil {
				configWithoutRoms = append(configWithoutRoms, match.configFile.FileName)
			} else {
				romsWithoutConfig = append(romsWithoutConfig, match.rom.FileName)
			}
		} else {
			if match.isExisting {
				continue
			}
			if strings.ToLower(match.rom.ConfigName()) != strings.ToLower(match.configFile.FileName) {
				createdFiles[match.matchType] = append(createdFiles[match.matchType], fmt.Sprintf("%s -> %s copied from: %s", match.rom.FileName, match.rom.ConfigName(), match.configFile.FileName))
			}
		}
	}

	log := ""
	if !p.commit {
		log = "[DRY]\n\n"
	}
	log += fmt.Sprintf("Found %d config files in: %s\nFound %d roms in: %s\n\n", configCount, configPath, romCount, romDirPath)
	log += fmt.Sprintf("Missing ROMs: %d\nMissing config: %d\n\n", len(configWithoutRoms), len(romsWithoutConfig))

	totalFileCount := len(createdFiles[MatchTypeExact]) + len(createdFiles[MatchTypeAlternate]) + len(createdFiles[MatchTypeFuzzy])
	skippedFileCount := 0
	if !shouldInclude(MatchTypeExact, matchFlag) {
		skippedFileCount += len(createdFiles[MatchTypeExact])
	}
	if !shouldInclude(MatchTypeAlternate, matchFlag) {
		skippedFileCount += len(createdFiles[MatchTypeAlternate])
	}
	if !shouldInclude(MatchTypeFuzzy, matchFlag) {
		skippedFileCount += len(createdFiles[MatchTypeFuzzy])
	}
	createdFilesCount := totalFileCount - skippedFileCount

	log += fmt.Sprintf("Created %d new files\n", createdFilesCount)
	log += fmt.Sprintf("Skipped %d new files\n\n", skippedFileCount)

	if len(configWithoutRoms) > 0 {
		sortAlphabetical(configWithoutRoms)
		log += fmt.Sprintf("CONFIG WITH MISSING ROMS\n%s\n\n", strings.Join(configWithoutRoms, "\n"))
	}

	if len(romsWithoutConfig) > 0 {
		sortAlphabetical(romsWithoutConfig)
		log += fmt.Sprintf("ROMS WITH MISSING CONFIG\n%s\n\n", strings.Join(romsWithoutConfig, "\n"))
	}

	if len(createdFiles[MatchTypeExact]) > 0 {
		sortAlphabetical(createdFiles[MatchTypeExact])
		log += fmt.Sprintf("NEW FILES (EXACT MATCHES)%s\n%s\n\n", skipped(MatchTypeExact, matchFlag), strings.Join(createdFiles[MatchTypeExact], "\n"))
	}

	if len(createdFiles[MatchTypeAlternate]) > 0 {
		sortAlphabetical(createdFiles[MatchTypeAlternate])
		log += fmt.Sprintf("NEW FILES (GOOD MATCHES)%s\n%s\n\n", skipped(MatchTypeAlternate, matchFlag), strings.Join(createdFiles[MatchTypeAlternate], "\n"))
	}

	if len(createdFiles[MatchTypeFuzzy]) > 0 {
		sortAlphabetical(createdFiles[MatchTypeFuzzy])
		log += fmt.Sprintf("NEW FILES (FUZZY MATCHES)%s\n%s\n\n", skipped(MatchTypeFuzzy, matchFlag), strings.Join(createdFiles[MatchTypeFuzzy], "\n"))
	}

	// write the log to a file and swallow any errors
	if err := p.writeLogToFile(configPath, log); err != nil {
		fmt.Errorf("failed to write log file: %s", err.Error())
	}
}

// writeLogToFile write the log to a log file in the config directory
func (p *Patcher) writeLogToFile(configDirPath, log string) error {
	logName := fmt.Sprintf("patch-log.%d.log", time.Now().Unix())
	return os.WriteFile(filepath.Join(configDirPath, logName), []byte(log), 0644)
}

// sortAlphabetical sorts a slice alphabetically
func sortAlphabetical(s []string) {
	sort.Sort(sort.StringSlice(s))
}

func skipped(matchType matchType, matchFlag matchType) string {
	if !shouldInclude(matchType, matchFlag) {
		return " [SKIPPED]"
	}
	return ""
}

// shouldInclude returns whether the current match type should be included in
// processing based on the given match flag value
func shouldInclude(matchType, matchFlag matchType) bool {
	if matchType == MatchTypeNone {
		return false
	}
	switch matchFlag {
	case MatchTypeExact:
		if matchType == MatchTypeExact {
			return true
		}
	case MatchTypeAlternate:
		if matchType == MatchTypeExact || matchType == MatchTypeAlternate {
			return true
		}
	case MatchTypeFuzzy:
		return true
	}
	return false
}
