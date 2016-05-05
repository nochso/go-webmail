package main

import "fmt"
import (
	"bitbucket.org/porkbonk/smtpd"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/nochso/mlog"
	"log"
	"os"
	"os/user"
	"path"
	"strings"
)

var db *sql.DB
var server *smtpd.Server
var dataDir = "./data"
var logDir = "./log"

var Version = ""
var BuildDate = ""

func main() {
	prepareConfig()
	mlog.DefaultFlags = log.Ldate | log.Ltime | log.Lmicroseconds
	mlog.Start(mlog.LevelTrace, path.Join(logDir, "smtpd.log"))
	printVersion()
	user, err := user.Current()
	if err == nil {
		mlog.Info("Running as user '%s'", user.Name)
	}
	mlog.Info("Loaded configuration:\n%s", getConfigDiff())
	prepareDirs()
	db = openDatabase()
	defer db.Close()
	prepareDatabase(db)
	prepareCert()
	server := prepareServer()

	addr := fmt.Sprintf(":%d", cfg.Port)
	mlog.Info("Starting smtpd server on %s", addr)
	err = server.ListenAndServe(addr)
	if err != nil {
		mlog.Fatalf("Error while listening/serving: %s", err)
	}
}

func prepareDirs() {
	if _, err := os.Stat(dataDir); err != nil {
		err := os.MkdirAll(dataDir, 0760)
		if err != nil {
			mlog.Fatalf("Unable to create data folder '%s': %s", dataDir, err)
		}
	}
	if _, err := os.Stat(logDir); err != nil {
		err := os.MkdirAll(logDir, 0760)
		if err != nil {
			mlog.Fatalf("Unable to create log folder '%s': %s", dataDir, err)
		}
	}
}

func printVersion() {
	version := "noch.so smtpd"
	if Version != "" {
		version = fmt.Sprintf("%s %s", version, Version)
	}
	if BuildDate != "" {
		version = fmt.Sprintf("%s (built %s)", version, BuildDate)
	}
	mlog.Info(strings.Repeat("-", len(version)))
	mlog.Info(version)
}
