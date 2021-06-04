package app

import (
    "bufio"
    "bytes"
    "fmt"
    "net/http"
    "io"
    "strconv"

    "github.com/AcidGo/zabbix-rms/pkg/catalog"
    "github.com/AcidGo/zabbix-rms/pkg/logger"
    "github.com/AcidGo/zabbix-rms/pkg/monitor"
    "github.com/AcidGo/zabbix-rms/pkg/zbxsend"
)

var (
    logging     *logger.ContextLogger
    zbxSend     *zbxsend.ZbxSend
)

func init() {
    logging = logger.FitContext("app")
}

func DiscoveryWorkspaceHd(w http.ResponseWriter, r *http.Request) {
    logging.Debugf("%s get request from %s", r.URL.Path, r.RemoteAddr)
    if r.Method != http.MethodGet {
        w.WriteHeader(http.StatusBadRequest)
        logging.Errorf("the request using %s method, expecting GET", r.Method)
        fmt.Fprintf(w, "please use GET method")
        return 
    }

    defer r.Body.Close()

    trs, err := catalog.GetWorkspaces()
    if err != nil {
        logging.Errorf("get an error while call catalog.GetWorkspaces: %v", err)
        w.WriteHeader(http.StatusBadRequest)
        fmt.Fprintf(w, "get an error: %v", err)
        return 
    }

    fmt.Fprintf(w, trs.FmtZbx())
}

func DiscoveryAppHd(w http.ResponseWriter, r *http.Request) {
    const queryTid = "tenant_id"
    const queryWid = "workspace_id"

    logging.Debugf("%s get request from %s", r.URL.Path, r.RemoteAddr)

    if r.Method != http.MethodGet {
        w.WriteHeader(http.StatusBadRequest)
        logging.Error("the request using %s method, expecting GET", r.Method)
        fmt.Fprintf(w, "please use GET method")
        return 
    }

    defer r.Body.Close()

    vals := r.URL.Query()
    logging.Debugf("url query values is %v", vals)

    val, ok := vals[queryTid]
    if !ok {
        w.WriteHeader(http.StatusBadRequest)
        fmt.Fprintf(w, "not found the arg %s", queryTid)
        return 
    }

    if len(val) != 1 {
        w.WriteHeader(http.StatusBadRequest)
        fmt.Fprintf(w, "the values of %s is not equal one", queryTid)
        return 
    }

    tenantId, err := strconv.Atoi(val[0])
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        fmt.Fprintf(w, "get an error: %v", err)
        return 
    }

    val, ok = vals[queryWid]
    if !ok {
        w.WriteHeader(http.StatusBadRequest)
        fmt.Fprintf(w, "not found the arg %s", queryWid)
        return 
    }

    if len(val) != 1 {
        w.WriteHeader(http.StatusBadRequest)
        fmt.Fprintf(w, "the values of %s is not equal one", queryWid)
        return 
    }

    workspaceId, err := strconv.Atoi(val[0])
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        fmt.Fprintf(w, "get an error: %v", err)
        return 
    }

    trs, err := catalog.GetApps(tenantId, workspaceId)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        fmt.Fprintf(w, "get an error: %v", err)
        return 
    }

    fmt.Fprintf(w, trs.FmtZbx())
}

func MonitorHd(w http.ResponseWriter, r *http.Request) {
    logging.Debugf("%s get request from %s", r.URL.Path, r.RemoteAddr)
    if r.Method != http.MethodPost {
        logging.Error("the request using %s method, expecting POST", r.Method)
        w.WriteHeader(http.StatusBadRequest)
        fmt.Fprintf(w, "please use POST method")
        return 
    }

    defer r.Body.Close()
    rd := bufio.NewReader(r.Body)
    bodyBytes, _ := io.ReadAll(rd)
    logging.Debugf("body of request: %s", string(bodyBytes))

    mp, err := monitor.Unpacket(bytes.NewReader(bodyBytes))
    if err != nil {
        logging.Error(err)
        w.WriteHeader(http.StatusBadRequest)
        fmt.Fprintf(w, "get an error while unpack body: %v", err)
        return 
    }

    // hardcode: gen an zabbix host name fmt
    zhost := fmt.Sprintf("rms_%d_%d", int(mp.TenantId), int(mp.WorkspaceId))
    // hardcode: gen an zabbix item key fmt
    zkey := fmt.Sprintf("rms.app[%s]", mp.AppName)
    err = zbxSend.Send(zhost, zkey, mp)
    if err != nil {
        logging.Error(err)
        w.WriteHeader(http.StatusBadRequest)
        fmt.Fprintf(w, "get an error while send msg to zabbix: %v", err)
        return 
    }

    fmt.Fprintf(w, "good")
}

func HttpSvr(l string, s *zbxsend.ZbxSend) (error) {
    zbxSend = s
    http.HandleFunc("/discovery_workspace", DiscoveryWorkspaceHd)
    http.HandleFunc("/discovery_app", DiscoveryAppHd)
    http.HandleFunc("/rms_monitor", MonitorHd)

    err := http.ListenAndServe(l, nil)
    return err
}