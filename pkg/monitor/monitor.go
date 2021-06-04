package monitor

import (
    "encoding/json"
    "io"
)

type MPacket struct {
    // Time            uint64  `json:"time"`
    Level           int     `json:"level"`
    Content         string  `json:"content"`
    // WorkspaceName   string  `json:"workspaceName"`
    TenantId        float32 `json:"tenantId,string"`
    WorkspaceId     float32 `json:"workspaceId,string"`
    AppName         string  `json:"app"`
    // DsId            string  `json:"dsId"`
    MonitorItem     string  `json:"monitorItem"`
    // Url             string  `json:"url"`
}

func Unpacket(r io.Reader) (*MPacket, error) {
    target := &MPacket{}
    err := json.NewDecoder(r).Decode(target)
    if err != nil {
        return nil, err
    }

    return target, nil
}