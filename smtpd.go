package main

import "fmt"
import (
	"database/sql"
	"log"
	"os"
	"os/user"
	"path"
	"runtime"
	"strings"

	"bitbucket.org/porkbonk/smtpd"
	"github.com/alexflint/go-arg"
	_ "github.com/mattn/go-sqlite3"
	"github.com/nochso/mlog"
	"path/filepath"
)

var db *sql.DB
var server *smtpd.Server
var dataDir = "./data"
var logDir = "./log"
var args struct {
	Verbose bool   `arg:"-v,help:enable detailed output"`
	Config  string `arg:"-c,help:path to a specific configuration file"`
	License bool   `arg:"-L,help:show software license"`
	Version bool   `arg:"-V,help:show version"`
}

var Version = ""
var BuildDate = ""

func main() {
	// Might exit early
	handleArgs()
	// Running as a daemon from now on
	prepareConfig()
	prepareLog()
	printVersion(true)
	runServer()
}

func handleArgs() {
	args.Config = "config.yaml"
	p, err := arg.NewParser(arg.Config{}, &args)
	if err != nil {
		fmt.Errorf("Error parsing CLI arguments: %s", err)
		os.Exit(-1)
	}
	err = p.Parse(os.Args[1:])
	if err == arg.ErrHelp {
		printVersion(false)
		p.WriteHelp(os.Stdout)
		os.Exit(0)
	}
	if err != nil {
		printVersion(false)
		p.Fail(err.Error())
	}
	if args.License {
		fmt.Println(string(MustAsset("LICENSE")))
		os.Exit(0)
	}
	if args.Version {
		printVersion(false)
		os.Exit(0)
	}
}

func runServer() {
	user, err := user.Current()
	if err == nil {
		mlog.Info("Running as user '%s' in %s", user.Name, getwd())
	}
	fp, err := filepath.Abs(args.Config)
	if err != nil {
		fp = args.Config
	}
	mlog.Info("Loaded configuration from %s", prettyPath(fp))
	mlog.Trace("%s", getConfigDiff())
	prepareDirs()
	db = openDatabase()
	defer db.Close()
	prepareDatabase(db)
	prepareCert()
	server := prepareServer()

	addr := fmt.Sprintf(":%d", cfg.Port)
	mlog.Info("Starting go-webmail server on %s", addr)
	err = server.ListenAndServe(addr)
	if err != nil {
		mlog.Fatalf("Error while listening/serving: %s", err)
	}
}

func getwd() string {
	wd, err := os.Getwd()
	if err != nil {
		return "unknown working directory"
	}
	return prettyPath(wd)
}

func prettyPath(path string) string {
	user, err := user.Current()
	if err == nil && strings.HasPrefix(path, user.HomeDir) {
		prefix := "~"
		if runtime.GOOS == "windows" {
			prefix = "%HOME%"
		}
		path = prefix + strings.TrimPrefix(path, user.HomeDir)
	}
	return path
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

func prepareLog() {
	mlog.DefaultFlags = log.Ldate | log.Ltime
	lvl := mlog.LevelInfo
	if args.Verbose {
		mlog.DefaultFlags = log.Ldate | log.Ltime | log.Lshortfile | log.Lmicroseconds
		lvl = mlog.LevelTrace
	}
	mlog.Start(lvl, path.Join(logDir, cfg.Log.Path))
}

func printVersion(log bool) {
	t := "go-webmail"
	if Version != "" {
		t += " " + Version
	}
	if BuildDate != "" {
		t += " " + BuildDate
	}
	t += "\n(C) 2016 Marcel Voigt - Released under the MIT license"
	if log {
		lines := strings.Split(t, "\n")
		mlog.Info("")
		for _, l := range lines {
			mlog.Info(l)
		}
		mlog.Info("")
	} else {
		fmt.Println(t + "\n")
	}
}
