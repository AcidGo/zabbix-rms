package catalog

import (
    "encoding/json"
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
    WorkspaceName   staring `json:"workspace_name"`
    WorkspaceId     int     `json:"workspace_id"`
}

type Workspaces []*Workspace

func (wss *Workspaces) FmtZbx() (string) {
    res := map[string]interface{}{
        "data": make([]map[string]interface{}, 0),
    }
    for _, val := range wss {
        t := make(map[string]interface{})
        t["{#TENANT_NAME}"] = val.TenantName
        t["{#TENANT_ID}"] = val.TenantId
        t["{#WORKSPACE_NAME}"] = val.WorkspaceName
        t["{#WORKSPACE_ID}"] = val.WorkspaceId
        res["data"] = append(res["data"], t)
    }

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

    wss := &Workspaces{}
    s := strings.Split(string(b), "\n")
    for _, val := range s {
        ws := &Workspace{}
        ss := strings.Split(val, workspaceSplitSymbol)
        if len(ss) != 4 {
            continue
        }
        ws.TenantName = ss[0]
        ws.TenantId, err = strconv.Atoi(s[1])
        if err != nil {
            continue
        }
        ws.WorkspaceName = s[2]
        ws.WorkspaceId, err = strconv.Atoi(s[3])
        if err != nil {
            continue
        }

        wss = append(wss, ws)
    }

    return wss, nil
}

type App struct {
    Name    string  `json:"app_name"`
}

type Apps []*App

func (apps *Apps) FmtZbx() (string) {
    res := map[string]interface{}{
        "data": make([]map[string]interface{}, 0),
    }
    for _, val := range apps {
        t := make(map[string]interface{})
        t["{#APP_NAME}"] = val.Name
        res["data"] = append(res["data"], t)
    }

    b, _ := json.Marshal(res)
    return string(b)
}

func GetApps(id int) (*Apps, error) {
    f, err := os.Open(appFilePath)
    if err != nil {
        return nil, err
    }
    defer f.Close()

    b, err := io.ReadAll(f)
    if err != nil {
        return nil, err
    }

    apps := &Apps{}
    s := strings.Split(string(b), "\n")
    for _, val := range s {
        app := &App{}
        ss := strings.Split(val, appSplitSymbol)
        if len(ss) != 2 {
            continue
        }
        if ss[0] != string(id) {
            continue
        }

        app.Name = s[1]

        apps = append(apps, app)
    }

    return apps, nil
}