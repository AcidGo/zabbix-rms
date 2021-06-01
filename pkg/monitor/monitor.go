package monitor

type MPacket struct {
    // Time            uint64  `json:"time"`
    Level           int     `json:"level"`
    Content         string  `json:"content"`
    // WorkspaceName   string  `json:"workspaceName"`
    TenantId        int     `json:"tenantId"`
    WorkspaceId     int     `json:"workspaceId"`
    AppName         string  `json:"app"`
    // DsId            string  `json:"dsId"`
    MonitorItem     string  `json:"monitorItem"`
    // Url             string  `json:"url"`
}

func Unpakc(r io.Reader) (*MPacket, error) {
    target := &MPacket{}
    err := json.NewDecoder(r).Decode(target)
    if err != nil {
        return nil, err
    }

    return target, nil
}