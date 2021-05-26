package monitor

import (
    "fmt"
)

type RMSAlert struct {
    rangeIP         string
    Time            uint    `json:"time"`
    Level           uint    `json:"level"`
    TenantName      string  `json:"tenantName"`
    TenantId        uint    `json:"tenantId"`
    WorkspaceName   string  `json:"workspaceName"`
    WorkspaceId     uint    `json:"workspaceId"`
    App             string  `json:"app"`
    DsId            string  `json:"dsId"`
    MonitorItem     string  `json:"monitorItem"`
    Url             string  `json:"url"`
}

func NewRMSAlert(ip string) (*RMSAlert, error) {
    if ip == ""{
        return nil, fmt.Errorf("the ip of rms alert is necessary")
    }

    return &RMSAlert{rangeIP: ip}, nil
}

func (r *RMSAlert) GetSlot() (string) {
    return fmt.Sprintf("rms_%s_%s", r.WorkspaceId, r.rangeIP)
}

func (r *RMSAlert) GetKey() (string) {
    return fmt.Sprintf("rms.%s.%s[%s]", )
}