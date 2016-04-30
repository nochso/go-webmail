package main

import "fmt"
import (
	"bitbucket.org/chrj/smtpd"
	"database/sql"
	"github.com/jbrodriguez/mlog"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"os/user"
	"path"
)

var db *sql.DB
var server *smtpd.Server
var host = "noch.so,loggle.tv"
var smtpPort = 25
var dataDir = "./data"
var logDir = "./log"

var Version = ""
var BuildDate = ""

func main() {
	mlog.Start(mlog.LevelInfo, path.Join(logDir, "smtpd.log"))
	mlog.Info("-----------------------------------------------")
	printVersion()
	user, err := user.Current()
	if err == nil {
		mlog.Info("Running as user '%s'", user.Name)
	}
	prepareDirs()
	db = openDatabase()
	defer db.Close()
	prepareDatabase(db)
	prepareCert()
	server := prepareServer()
	mlog.Info("Starting smtpd server")
	server.ListenAndServe(fmt.Sprintf(":%d", smtpPort))
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
	mlog.Info(version)
}
