package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	fmt.Println("Starting retrosheet-go")

	teams := []Team{}

	teamFileName := ""
	rosterFiles := []string{}
	eventFiles := []string{}

	// load the data
	entries, err := ListFilesRecursively("data", true)
	if err != nil {
		fmt.Println("unable to read data directory:", err)
		os.Exit(1)
	}
	for _, entry := range entries {
		if strings.HasSuffix(entry, ".EVA") || strings.HasSuffix(entry, ".EVN") {
			eventFiles = append(eventFiles, entry)
			continue
		}
		if strings.HasSuffix(entry, ".ROS") {
			rosterFiles = append(rosterFiles, entry)
			continue
		}
		if strings.Contains(entry, "TEAM") {
			teamFileName = entry
			continue
		}
	}

	fmt.Println("teamFile =", teamFileName)
	fmt.Println("rosterFiles =", rosterFiles)
	fmt.Println("eventFiles =", eventFiles)

	// process team file
	teamFile, err := os.Open(teamFileName)
	if err != nil {
		fmt.Println("error opening team file:", err)
	}

	scanner := bufio.NewScanner(teamFile)
	for scanner.Scan() {
		split := strings.Split(scanner.Text(), ",")
		teams = append(teams, Team{
			Code:     split[0],
			League:   split[1],
			Location: split[2],
			Name:     split[3],
		})
	}
	for _, team := range teams {
		fmt.Println(team)
	}
}

type Team struct {
	Location string
	Name     string
	Code     string
	League   string
}

func (t Team) String() string {
	return fmt.Sprintf("%s %s %s %s", t.Code, t.League, t.Location, t.Name)
}

// ListFilesRecursively walks the directory tree rooted at 'root'
// and returns a slice of strings containing the full paths of all files found.
func ListFilesRecursively(root string, firstLevelOnly bool) ([]string, error) {
	var filePaths []string // Slice to store the full paths of files

	// filepath.WalkDir is a more efficient and modern alternative to filepath.Walk
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		// If an error occurs while accessing a path, return it to stop the walk.
		if err != nil {
			return err
		}

		// Check if the current entry is a regular file.
		if d.Type().IsRegular() {
			filePaths = append(filePaths, path)
		}

		// skip walking nested directories
		if d.IsDir() && firstLevelOnly && path != root {
			return fs.SkipDir
		}

		// Return nil to continue the walk.
		return nil
	})

	// If there was an error during the walk (other than those handled within the walk function),
	// return the error.
	if err != nil {
		return nil, fmt.Errorf("error walking the directory tree: %w", err)
	}

	// Return the slice of file paths.
	return filePaths, nil
}
