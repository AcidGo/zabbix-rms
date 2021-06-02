package catalog

import (
    "encoding/json"
    "io"
    "os"
    "strconv"
    "strings"
)

const (
    workspaceFilePath       = "ws.map"
    workspaceSplitSymbol    = "@@"
    appFilePath             = "app.map"
    appSplitSymbol          = "@@"
)

type Workspace struct {
    TenantName      string  `json:"tenant_name"`
    TenantId        int     `json:"tenant_id"`
    WorkspaceName   string  `json:"workspace_name"`
    WorkspaceId     int     `json:"workspace_id"`
}

type Workspaces struct {
    ws []*Workspace
}

func (wss *Workspaces) FmtZbx() (string) {
    res := map[string]interface{}{}
    var mSlice []map[string]interface{}
    for _, val := range wss.ws {
        t := make(map[string]interface{})
        t["{#TENANT_NAME}"] = val.TenantName
        t["{#TENANT_ID}"] = val.TenantId
        t["{#WORKSPACE_NAME}"] = val.WorkspaceName
        t["{#WORKSPACE_ID}"] = val.WorkspaceId
        mSlice = append(mSlice, t)
    }
    res["data"] = mSlice

    b, _ := json.Marshal(res)
    return string(b)
}

func GetWorkspaces() (*Workspaces, error) {
    f, err := os.Open(workspaceFilePath)
    if err != nil {
        return nil, err
    }
    defer f.Close()

    b, err := io.ReadAll(f)
    if err != nil {
        return nil, err
    }

    wss := &Workspaces{ws: make([]*Workspace, 0)}
    s := strings.Split(string(b), "\n")
    for _, val := range s {
        ws := &Workspace{}
        ss := strings.Split(val, workspaceSplitSymbol)
        if len(ss) != 4 {
            continue
        }
        ws.TenantName = ss[0]
        ws.TenantId, err = strconv.Atoi(ss[1])
        if err != nil {
            continue
        }
        ws.WorkspaceName = ss[2]
        ws.WorkspaceId, err = strconv.Atoi(ss[3])
        if err != nil {
            continue
        }

        wss.ws = append(wss.ws, ws)
    }

    return wss, nil
}

type App struct {
    Name    string  `json:"app_name"`
}

type Apps struct {
    apps    []*App
}

func (apps *Apps) FmtZbx() (string) {
    res := map[string]interface{}{}
    var mSlice []map[string]interface{}
    for _, val := range apps.apps {
        t := make(map[string]interface{})
        t["{#APP_NAME}"] = val.Name
        mSlice = append(mSlice, t)
    }
    res["data"] = mSlice

    b, _ := json.Marshal(res)
    return string(b)
}

func GetApps(tId, wId int) (*Apps, error) {
    f, err := os.Open(appFilePath)
    if err != nil {
        return nil, err
    }
    defer f.Close()

    b, err := io.ReadAll(f)
    if err != nil {
        return nil, err
    }

    apps := &Apps{apps: make([]*App, 0)}
    s := strings.Split(string(b), "\n")
    for _, val := range s {
        ss := strings.Split(val, appSplitSymbol)
        if len(ss) != 3 {
            continue
        }
        if ss[0] != strconv.Itoa(tId) || ss[1] != strconv.Itoa(wId) {
            continue
        }

        app := &App{Name: ss[2]}
        apps.apps = append(apps.apps, app)
    }

    return apps, nil
}