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
	"github.com/jackc/pgx/v5/pgxpool"
)

type Migration struct {
	Id       int
	Filepath string
	Name     string
	Type     string
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

func getPendingMigrations(files []string, lastCommittedId int32) []Migration {
	var latestId = int(lastCommittedId)
	var migs []Migration

	// sheer insouciance
	lenIdentifier := len("nnnn")

	for _, file := range files {
		len := len(file)
		if len < lenIdentifier {
			log.Fatalf("Encountered file [%s] with bad name - cannot extract unique identifier for migration. Bailing.", file)
		}
		id, _ := strconv.Atoi(file[0:lenIdentifier])

		if id > latestId {
			_, filename := filepath.Split(file)
			migName := strings.Split(filename, ".")
			migType := "core"
			if strings.EqualFold(filepath.Ext(filename), ".json") {
				migType = "authz"
			}

			migs = append(migs, Migration{
				Id:       id,
				Filepath: file,
				Name:     migName[0][5:],
				Type:     migType,
			})
		}
	}

	return migs
}

func sqlForMigrationsRecord(mig Migration) string {
	if mig.Id == 0 || helpers.IsBlank(mig.Name) || helpers.IsBlank(mig.Hash) {
		log.Fatalf("Need valid id, name and hash for migration!")
	}

	return fmt.Sprintf("insert into migrations(id, name, mig_type, hash) values (%d, '%s', '%s', '%s');", mig.Id, mig.Name, mig.Type, mig.Hash)
}

func runInTransaction(ctx context.Context, tx pgx.Tx, sqlStrings []string, ch chan bool) {
	allOK := true

	for _, sql := range sqlStrings {
		_, err := tx.Exec(ctx, sql)
		if err != nil {
			fmt.Printf("SQL transaction error: %v\n", err)
			allOK = false
			break
		}
	}

	ch <- allOK
}

func getLatestCommittedMigrationId(ctx context.Context, pool *pgxpool.Pool, ch chan int32) {
	var id int32
	var name string
	var hash string

	rows, err := pool.Query(ctx, "select id, name, hash from public.migrations order by id desc limit 1")
	if err != nil {
		log.Fatalf("SQL query error: %v\n", err)
	}

	defer rows.Close()

	if rows.Next() {
		rows.Scan(&id, &name, &hash)
		fmt.Printf("🔒 Last migration: [%d][%s][%s]\n", id, name, hash)
		ch <- id
		return
	}

	ch <- 0
}

func runMigration(ctx context.Context, config *helpers.Config, pool *pgxpool.Pool) {
	fsys := os.DirFS(config.MigrationsDir)
	files, errFindAllFiles := allFiles(fsys)
	if errFindAllFiles != nil {
		log.Fatal("can't list migration files!")
	}

	ch := make(chan int32)
	chmig := make(chan bool)

	// prepare for migrations
	go getLatestCommittedMigrationId(ctx, pool, ch)
	lastCommittedId := <-ch
	migs := getPendingMigrations(files, lastCommittedId)

	var ok bool
	for _, mig := range migs {
		sqlTemplate := getFileContents(fmt.Sprintf("%s/%s", config.MigrationsDir, mig.Filepath))
		sqlStr := helpers.HydrateSQLTemplate(sqlTemplate, *config)

		mig.Hash = makeHash(sqlStr)
		sqlHashStore := sqlForMigrationsRecord(mig)

		fingerprint := fmt.Sprintf("Migration for [%d][%s][%s]", mig.Id, mig.Name, mig.Hash)

		tx, _ := pool.Begin(ctx)
		go runInTransaction(ctx, tx, []string{sqlStr, sqlHashStore}, chmig)
		ok = <-chmig

		if ok {
			fmt.Printf("✅ Committing: %s\n", fingerprint)
			tx.Commit(ctx)
		} else {
			fmt.Printf("❌ Failed - Rolling Back: %s\n", fingerprint)
			tx.Rollback(ctx)
			break //stop processing
		}
	}
}

func createMigration(
	ctx context.Context,
	config *helpers.Config,
	pool *pgxpool.Pool,
	migName string,
	migType string,
) {
	ch := make(chan int32)

	go getLatestCommittedMigrationId(ctx, pool, ch)
	lastCommittedId := <-ch

	nextId := int(lastCommittedId) + 1

	const allowed = "abcdefghijklmnopqrstuvwxyz0123456789-_"
	runes := []rune{}

	// toss any funky runes
	for _, c := range migName {
		nextChar := strings.ToLower(string(c))
		if strings.Contains(allowed, nextChar) {
			runes = append(runes[:], c)
		}
	}

	// identify regular sql migrations vs metadata migrations controlled by other subsystems
	ext := "sql"
	if migType != "core" {
		ext = "json"
	}
	cleanStr := strings.Join(strings.Split(string(runes), " "), "-")
	targetName := fmt.Sprintf("%04d-%s.%s", nextId, cleanStr, ext)
	targetPath := filepath.Join(config.MigrationsDir, targetName)

	emptyBytArray := []byte("")
	os.WriteFile(targetPath, emptyBytArray, 0644)
	fmt.Printf("✨ New migration file ready: %s\n", targetName)
}

func initialiseMigrations(ctx context.Context) {
	config := helpers.LoadConfig()
	sqlStr := "create table if not exists migrations (id int primary key, name text not null, mig_type text not null, hash text not null);"
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
	helpTxt := "Must provide a valid command to execute: init|create|migrate."
	createHelpTxt := "The create command requires a name for your migration."

	if len(os.Args) < 2 {
		fmt.Println(helpTxt)
		os.Exit(0)
	} else {
		config := helpers.LoadConfig()

		// need a pool for multi-statement queries
		pool, err := pgxpool.New(ctx, config.PgUrl)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
			os.Exit(1)
		}
		defer pool.Close()

		switch cmd := os.Args[1]; cmd {
		case "init":
			initialiseMigrations(ctx)
		case "migrate":
			runMigration(ctx, config, pool)
		case "create", "create-meta":
			if len(os.Args) < 3 {
				fmt.Println(createHelpTxt)
				break
			}
			nameArgs := strings.Join(os.Args[2:], "-")
			if cmd == "create" {
				createMigration(ctx, config, pool, nameArgs, "core")
			} else {
				createMigration(ctx, config, pool, nameArgs, "json")
			}
		default:
			fmt.Println(helpTxt)
		}
	}
	os.Exit(0)
}
