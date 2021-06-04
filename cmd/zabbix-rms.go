package main

import (
    "flag"
    "fmt"
    "os"
    "path/filepath"

    "github.com/AcidGo/zabbix-rms/app"
    "github.com/AcidGo/zabbix-rms/pkg/logger"
    "github.com/AcidGo/zabbix-rms/pkg/zbxsend"
)

var (
    // flag vars
    httplsn     string
    zbxHost     string
    zbxPort     int
    lName       string
    lDir        string
    lLevel      string

    // runtime vars
    logging     *logger.ContextLogger
    zbxSend     *zbxsend.ZbxSend

    // app info
    AppName             string
    AppAuthor           string
    AppVersion          string
    AppGitCommitHash    string
    AppBuildTime        string
    AppGoVersion        string
)

func init() {
    flag.StringVar(&httplsn, "l", ":6768", "http server listening")
    flag.StringVar(&zbxHost, "zh", "", "zabbix server/proxy address")
    flag.IntVar(&zbxPort, "zp", 10051, "zabbix server/proxy port")
    flag.StringVar(&lName, "lname", "zabbix-rms.log", "logger file name")
    flag.StringVar(&lDir, "ldir", "/tmp", "logger file dir placing")
    flag.StringVar(&lLevel, "llevel", "info", "logging level setting")
    flag.Usage = flagUsage
    flag.Parse()

    if !IsDir(lDir) {
        logging.Fatalf("the log dir %s is not a dir not no exists", lDir)
    }
    logPath := filepath.Join(lDir, lName)
    err := logger.LogFileSetting(logPath)
    if err != nil {
        logging.Fatal(err)
    }
    // logger.ReportCallerSetting(true)
    err = logger.LogLevelSetting(lLevel)
    if err != nil {
        logging.Fatal(err)
    }
}

func main() {
    var err error

    zbxSend, err := zbxsend.NewZbxSend(zbxHost, zbxPort)
    if err != nil {
        logging.Fatal(err)
    }

    err = app.HttpSvr(httplsn, zbxSend)
    if err != nil {
        logging.Fatal(err)
    }
}

func flagUsage() {
    usageMsg := fmt.Sprintf(`App: %s
Version: %s
Author: %s
GitCommit: %s
BuildTime: %s
GoVersion: %s
Options:
`, AppName, AppVersion, AppAuthor, AppGitCommitHash, AppBuildTime, AppGoVersion)

    fmt.Fprintf(os.Stderr, usageMsg)
    flag.PrintDefaults()
}

func IsDir(path string) (bool) {
    s, err := os.Stat(path)
    if err != nil {
        return false
    }
    return s.IsDir()
}