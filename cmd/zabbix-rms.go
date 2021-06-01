package main

import (
    "flag"

    "github.com/AcidGo/zabbix-rms/pkg/zbxsend"
)

var (
    // flag vars
    httplsn     string
    zbxHost     string
    zbxPort     int

    // runtime vars
    zbxSend     *zbxsend.ZbxSend
)

func init() {
    var err error

    flag.StringVar(&httplsn, "l", ":6768", "http server listening")
    flag.StringVar(&zbxHost, "zh", "", "zabbix server/proxy address")
    flag.IntVar(&zbxPort, "zp", 10051, "zabbix server/proxy port")
    flag.Parse()
}

func main() {
    var err error

    zbxSend, err := zbxsend.NewZbxSend(zbxHost, zbxPort)
    if err != nil {
        log.Fatal(err)
    }

    err = app.HttpSvr(httplsn, zbxSend)
    if err != nil {
        log.Fatal(err)
    }
}