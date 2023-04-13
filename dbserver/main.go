package main

import (
	"context"
	"dbserver/helpers"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"strconv"

	"github.com/jackc/pgx/v5"
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
	var latestId int
	var latestMigrationFile string

	// sheer insouciance
	lenPrefix := len("migrations/")
	lenIdentifier := len("nnnn")

	for _, file := range files {
		len := len(file)
		if len < 15 {
			log.Fatalf("Encountered file [%s] with bad name - cannot extract unique identifier for migration. Bailing.", file)
		}
		id, _ := strconv.Atoi(file[lenPrefix:(lenPrefix + lenIdentifier)])
		if id > latestId {
			latestId = id
			latestMigrationFile = file
		}
	}

	return latestMigrationFile
}

func runMigration(ctx context.Context, pgUrl string, sqlStr string) {
	fmt.Println(sqlStr)
	conn, err := pgx.Connect(context.Background(), pgUrl)

	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	defer conn.Close(ctx)

	tx, _ := conn.Begin(ctx)
	statusText, errTx := tx.Exec(ctx, sqlStr)
	if errTx != nil {
		fmt.Printf("SQL transaction error: %v\nStatus text:%s\n", errTx, statusText)
		fmt.Println("*** Rolling back!! ***")
		tx.Rollback(ctx)
	} else {
		fmt.Println("*** Committing... ***")
		tx.Commit(ctx)
	}
}

func main() {
	files, errFindAllFiles := allFiles(&migrationsFS)
	if errFindAllFiles != nil {
		log.Fatal("can't list migration files!")
	}

	latestMigrationFilePath := latestMigration(files)
	sqlTemplate := getFileContents(latestMigrationFilePath)
	sqlStr, pgUrl := helpers.Prep(sqlTemplate)

	runMigration(context.Background(), pgUrl, sqlStr)
}
