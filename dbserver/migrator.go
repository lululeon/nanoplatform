package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"strconv"
)

//go:embed migrations/*
var migrationsFS embed.FS

func getFileContents(migrationFile string) string {
	byteArr, errFile := migrationsFS.ReadFile(migrationFile)
	if errFile != nil {
		log.Fatal("Failed to read file!")
	}
	return string(byteArr)
}

func allFiles(efs *embed.FS) (files []string, err error) {
	startFromRoot := "."
	if err := fs.WalkDir(efs, startFromRoot, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		files = append(files, path)

		return nil
	}); err != nil {
		return nil, err
	}

	return files, nil
}

func latestMigration(files []string) string {
	var latest string
	for _, file := range files {
		if (len(file)) < 15 {
			log.Fatalf("Encountered file [%s] with bad name - cannot extract unique identifier for migration. Bailing.", file)
		}
		id, _ := strconv.Atoi(file[11:15])
		fmt.Printf("[%s] - [%d]\n", file, id)
	}

	return latest
}

func main() {
	files, errFindAllFiles := allFiles(&migrationsFS)
	if errFindAllFiles != nil {
		log.Fatal("can't list migration files!")
	}

	latestMigrationFilePath := latestMigration(files)
	sqlStr, pgUrl := Prep(latestMigrationFilePath)
	// for _, file := range files {
	// 	fmt.Println(file.Name(), file.IsDir())
	// }
}
