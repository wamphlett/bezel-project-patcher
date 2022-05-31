package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/wamphlett/bezel-project-patcher/pkg/files"
	"github.com/wamphlett/bezel-project-patcher/pkg/patching"
)

var (
	commit        *bool
	fuzzyMatching *bool
	exactOnly     *bool
)

func init() {
	commit = flag.Bool("commit", false, "commit will write the new config files")
	exactOnly = flag.Bool("exact-only", false, "matching will only include exact matches")
	fuzzyMatching = flag.Bool("fuzzy", false, "matching will include fuzzy matches")
}

func main() {
	flag.Parse()
	if len(flag.Args()) != 2 {
		fmt.Println("expected 2 arguments. example: bezel-project-patcher <path-to-config-directory> <path-to-rom-directory>")
		return
	}

	if *exactOnly && *fuzzyMatching {
		fmt.Println("cannot use both fuzzy matching (--fuzzy) and exact only matching (--exact-only) together")
		return
	}

	if !*commit {
		fmt.Println("DRY RUN ONLY. No files will be modified.")
	}

	configDirectory := flag.Arg(0)
	romDirectory := flag.Arg(1)

	matchFlag := patching.MatchTypeAlternate
	if *exactOnly {
		matchFlag = patching.MatchTypeExact
	} else if *fuzzyMatching {
		matchFlag = patching.MatchTypeFuzzy
	}

	fileManager := files.FileManager{}
	patcher := patching.NewPatcher(&fileManager, *commit)

	if err := patcher.PatchDirectory(configDirectory, romDirectory, matchFlag); err != nil {
		fmt.Errorf("failed to successfully patch directory: %s", err.Error())
		os.Exit(1)
	}

	if !*commit {
		fmt.Println("Patch finished but no files were modified. It is strongly recommended to check logs before committing the changes.")
		fmt.Printf("Run 'bezel-project-patcher --commit %s %s' to commit the changes\n", configDirectory, romDirectory)
		return
	}

	fmt.Printf("Successfully patched config directory %s. See the log file for more information.", configDirectory)
}
