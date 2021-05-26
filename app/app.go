package app

import (
    "fmt"
    "net/http"

    "github.com/AcidGo/zabbix-rms/pkg/monitor"
)

var (
    workspaceIp     string
    zbxSender       *sender.ZbxSender
)

func monitorHanlder(w http.ResponseWriter, r *http.Request) {
    var err error

    if r.Method != http.MethodPost {
        w.WriteHeader(http.StatusBadRequest)
        fmt.Fprintf(w, "please use POST method")
        return 
    }

    defer r.Body.Close()

    msg, err := monitor.NewRMSAlert(workspaceIp)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        fmt.Fprintf(w, "gen rms alert struct is failed: %v", err)
        return 
    }

    err = json.NewDecoder(r.Body).Decode(msg)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        fmt.Fprintf(w, "decoding struct from request's body is failed: %v", err)
        return 
    }

    err = zbxSender.Send(msg.GetSlot(), msg.GetKey(), msg.GetVal())
}