package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

const (
	DescriptionAddFlag    = "Add current directory to bookmarks"
	DescriptionRemoveFlag = "Remove current directory from bookmarks"
)

type flags struct {
	add    *bool
	remove *bool
}

type bookmark struct {
	Alias string `json:"alias"`
	Path  string `json:"path"`
}

func readFlags() *flags {
	addFlag := flag.Bool("add", false, DescriptionAddFlag)
	removeFlag := flag.Bool("remove", false, DescriptionRemoveFlag)
	flag.Parse()

	return &flags{
		add:    addFlag,
		remove: removeFlag,
	}
}

func createConfigFile(path string, initialData []byte) (ok bool) {
	cfgDir := filepath.Dir(path)
	err := createConfigDir(cfgDir)
	if err != nil {
		return false
	}

	err = os.WriteFile(path, initialData, 0600)
	return err == nil
}

func createConfigDir(dir string) error {
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Error getting home directory:", err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal("Failed to get current directory", err)
	}

	flags := readFlags()
	if *flags.add {
		fmt.Println("Add flag is set")
	}
	if *flags.remove {
		fmt.Println("Remove flag is set")
	}

	newBookmark := bookmark{
		Alias: "bruno",
		Path:  cwd,
	}

	jsonData, err := json.MarshalIndent(newBookmark, "", "  ")
	if err != nil {
		log.Fatal("Error marshaling JSON:", err)
	}

	cfgPath := filepath.Join(homeDir, ".config", "bookmarker", "bookmarker.json")
	cfgFile, err := os.ReadFile(cfgPath)
	if err != nil {
		log.Println("Config file does not exist, creating it at:", cfgPath)
		if ok := createConfigFile(cfgPath, jsonData); !ok {
			log.Fatal("Failed to create configuration file")
		}
	}

	cfgFile, _ = os.ReadFile(cfgPath)

	fmt.Println("File content:\n", string(cfgFile))

	fmt.Println(cwd)
}
