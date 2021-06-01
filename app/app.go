package app

import (
    "net/http"

    "github.com/AcidGo/zabbix-rms/pkg/catalog"
    "github.com/AcidGo/zabbix-rms/pkg/monitor"
    "github.com/AcidGo/zabbix-rms/pkg/zbxsend"
)

var (
    zbxSend     *zbxsend.ZbxSend
)

func DiscoveryWorkspaceHd(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        w.WriteHeader(http.StatusBadRequest)
        fmt.Fprintf(w, "please use GET method")
        return 
    }

    defer r.Body.Close()

    trs, err := catalog.GetWorkspaces()
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        fmt.Fprintf(w, "get an error: %v", err)
        return 
    }

    fmt.Fprintf(w, trs.FmtZbx())
}

func DiscoveryAppHd(w http.ResponseWriter, r *http.Request) {
    const queryArg = "tenant_id"
    if r.Method != http.MethodGet {
        w.WriteHeader(http.StatusBadRequest)
        fmt.Fprintf(w, "please use GET method")
        return 
    }

    defer r.Body.Close()

    vals := r.URL.Query()
    val, ok := vals[queryArg]
    if !ok {
        w.WriteHeader(http.StatusBadRequest)
        fmt.Fprintf(w, "not found the arg %s", queryArg)
        return 
    }

    if len(val) != 1 {
        w.WriteHeader(http.StatusBadRequest)
        fmt.Fprintf(w, "the values of %s is not equal one", queryArg)
        return 
    }

    tenantId, err := strconv.Atoi(val[0])
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        fmt.Fprintf(w, "get an error: %v", err)
        return 
    }

    trs, err := catalog.GetApps(tenantId)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        fmt.Fprintf(w, "get an error: %v", err)
        return 
    }

    fmt.Fprintf(w, trs.FmtZbx())
}

func MonitorHd(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        w.WriteHeader(http.StatusBadRequest)
        fmt.Fprintf(w, "please use POST method")
        return 
    }

    defer r.Body.Close()

    mp, err := monitor.Unpakc(r.Body)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        fmt.Fprintf(w, "get an error while unpack body: %v", err)
        return 
    }

    // hardcode: gen an zabbix host name fmt
    zhost := fmt.Sprintf("rms_%d_%d", mp.TenantId, mp.WorkspaceId)
    // hardcode: gen an zabbix item key fmt
    zkey := fmt.Sprintf("rms.app[%s]", mp.AppName)
    err = zbxSend.Send(zhost, zkey, mp)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        fmt.Fprintf(w, "get an error while send msg to zabbix: %v", err)
        return 
    }

    fmt.Fprintf(w, "good")
}

func HttpSvr(l string, s *zbxsend.ZbxSend) (error) {
    zbxSend = s
    http.HandleFunc("/discover_workspace", DiscoveryWorkspaceHd)
    http.HandleFunc("/discover_app", DiscoveryAppHd)
    http.HandleFunc("/rms_monitor", MonitorHd)

    err := http.ListenAndServe(l, nil)
    return err
}