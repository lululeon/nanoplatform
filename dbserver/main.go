package main

import (
	"context"
	"crypto/sha1"
	"dbserver/helpers"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
)

type Migration struct {
	Id       int
	Filepath string
	Name     string
	Hash     string
}

func getFileContents(migrationFile string) string {
	byteArr, errFile := os.ReadFile(migrationFile)
	if errFile != nil {
		log.Fatal("Failed to read file!")
	}
	return string(byteArr)
}

func allFiles(filesys fs.FS) (files []string, err error) {
	startFromRoot := "."
	if err := fs.WalkDir(filesys, startFromRoot, func(path string, d fs.DirEntry, err error) error {
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

func makeHash(str string) string {
	byteArr := []byte(str)

	// could use bcrypt or similar for beefier sec but we are not fussed atm:
	hasher := sha1.New()
	_, hashErr := hasher.Write(byteArr)
	if hashErr != nil {
		log.Fatal("can't hash the migration!")
	}

	// c/shouldve stored the integer hash, but for some reason I'd fancied storing the hex representation of the checksum and now the tbl is built that way, so:
	return fmt.Sprintf("%x", hasher.Sum(nil))
}

func latestMigration(files []string) Migration {
	var latestId int = 0
	var latestMigrationFile string

	// sheer insouciance
	lenIdentifier := len("nnnn")

	for _, file := range files {
		len := len(file)
		if len < lenIdentifier {
			log.Fatalf("Encountered file [%s] with bad name - cannot extract unique identifier for migration. Bailing.", file)
		}
		id, _ := strconv.Atoi(file[0:lenIdentifier])
		if id > latestId {
			latestId = id
			latestMigrationFile = file
		}
	}

	_, filename := filepath.Split(latestMigrationFile)
	migName := strings.Split(filename, ".")

	return Migration{
		Id:       latestId,
		Filepath: latestMigrationFile,
		Name:     migName[0][5:],
	}
}

func sqlForMigrationsRecord(mig Migration) string {
	if mig.Id == 0 || helpers.IsBlank(mig.Name) || helpers.IsBlank(mig.Hash) {
		log.Fatalf("Need valid id, name and hash for migration!")
	}

	return fmt.Sprintf("insert into migrations(id, name, hash) values (%d, '%s', '%s');", mig.Id, mig.Name, mig.Hash)
}

func runInTransaction(ctx context.Context, tx pgx.Tx, sqlStrings []string) bool {
	allOK := true

	for _, sql := range sqlStrings {
		_, err := tx.Exec(ctx, sql)
		if err != nil {
			fmt.Printf("SQL transaction error: %v\n", err)
			allOK = false
			break
		}
	}

	return allOK
}

func runMigration(ctx context.Context) {
	config := helpers.LoadConfig()

	fsys := os.DirFS(config.MigrationsDir)
	files, errFindAllFiles := allFiles(fsys)
	if errFindAllFiles != nil {
		log.Fatal("can't list migration files!")
	}

	mig := latestMigration(files)
	sqlTemplate := getFileContents(fmt.Sprintf("%s/%s", config.MigrationsDir, mig.Filepath))

	sqlStr := helpers.HydrateSQLTemplate(sqlTemplate, *config)
	mig.Hash = makeHash(sqlStr)
	sqlHashStore := sqlForMigrationsRecord(mig)
	fingerprint := fmt.Sprintf("migration for [%d][%s][%s]", mig.Id, mig.Name, mig.Hash)

	conn, err := pgx.Connect(ctx, config.PgUrl)

	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	defer conn.Close(ctx)

	tx, _ := conn.Begin(ctx)
	ok := runInTransaction(ctx, tx, []string{sqlStr, sqlHashStore})

	if ok {
		fmt.Printf("✅ Committing: %s\n", fingerprint)
		tx.Commit(ctx)
	} else {
		fmt.Printf("❌ Failed - Rolling Back: %s\n", fingerprint)
		tx.Rollback(ctx)
	}
}

func initialiseMigrations(ctx context.Context) {
	config := helpers.LoadConfig()
	sqlStr := "create table if not exists migrations (id int primary key, name text not null, hash text not null);"
	conn, err := pgx.Connect(ctx, config.PgUrl)

	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	defer conn.Close(ctx)
	statusText, errTx := conn.Exec(ctx, sqlStr)
	if errTx != nil {
		fmt.Printf("SQL transaction error: %v\nStatus text:%s\n", errTx, statusText)
	} else {
		fmt.Println("Migrations table ready.")
	}
}

func main() {
	ctx := context.Background()
	helpTxt := "Must provide a valid command to execute: init|create|migrate"

	if len(os.Args) != 2 {
		fmt.Println(helpTxt)
		os.Exit(0)
	} else {
		switch cmd := os.Args[1]; cmd {
		case "init":
			initialiseMigrations(ctx)
		case "migrate":
			runMigration(ctx)
		default:
			fmt.Println(helpTxt)
		}
	}
	os.Exit(0)
}
