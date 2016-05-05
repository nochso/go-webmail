package main

import "fmt"
import (
	"bitbucket.org/porkbonk/smtpd"
	"database/sql"
	"github.com/alexflint/go-arg"
	_ "github.com/mattn/go-sqlite3"
	"github.com/nochso/mlog"
	"log"
	"os"
	"os/user"
	"path"
	"runtime"
	"strings"
)

var db *sql.DB
var server *smtpd.Server
var dataDir = "./data"
var logDir = "./log"
var args struct {
	Verbose bool   `arg:"-v,help:enable detailed output"`
	Config  string `arg:"-c,help:path to a specific configuration file"`
}

var Version = ""
var BuildDate = ""

func main() {
	args.Config = "config.yaml"
	arg.MustParse(&args)
	prepareConfig()
	mlog.DefaultFlags = log.Ldate | log.Ltime
	lvl := mlog.LevelInfo
	if args.Verbose {
		mlog.DefaultFlags = log.Ldate | log.Ltime | log.Lshortfile | log.Lmicroseconds
		lvl = mlog.LevelTrace
	}
	mlog.Start(lvl, path.Join(logDir, "smtpd.log"))
	printVersion()
	user, err := user.Current()
	if err == nil {
		mlog.Info("Running as user '%s' in %s", user.Name, getwd(user))
	}
	mlog.Info("Loaded configuration from %s", args.Config)
	mlog.Trace("%s", getConfigDiff())
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

func getwd(user *user.User) string {
	wd, err := os.Getwd()
	if err != nil {
		wd = "unknown working directory"
	}
	if strings.HasPrefix(wd, user.HomeDir) {
		prefix := "~"
		if runtime.GOOS == "windows" {
			prefix = "%HOME%"
		}
		wd = prefix + strings.TrimPrefix(wd, user.HomeDir)
	}
	return wd
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
	version := "go-webmail (personal MDA+webmail)"
	if Version != "" {
		version += " " + Version
	}
	if BuildDate != "" {
		version = fmt.Sprintf("%s (built %s)", version, BuildDate)
	}
	mlog.Info(strings.Repeat("-", len(version)))
	mlog.Info(version)
	mlog.Info("Copyright (c) 2016 Marcel Voigt")
	mlog.Info("License: MIT")
}
