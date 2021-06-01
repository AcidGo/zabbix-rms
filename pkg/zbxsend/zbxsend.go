package zbxsend

import (
    "encoding/json"
    "github.com/blacked/go-zabbix"
)

type ZbxSend struct {
    host    string
    port    int
    s       *zabbix.Sender
}

func NewZbxSend(host string, port int) (*ZbxSend, error) {
    return &ZbxSend{
        host:   host, 
        port:   port,
        s:      zabbix.NewSender(host, port),
    }, nil
}

func (zs *ZbxSend) Send(zhost, key string, value interface{}) (error) {
    b, err := json.Marshal(value)
    if err != nil {
        return err
    }

    metrics := make([]*zabbix.Metric, 1)
    metrics[0] = zabbix.NewMetric(
        host,
        key,
        string(b),
    )

    packet := zabbix.NewPacket(metrics)
    _, err := zs.s.Send(packet)

    return err
}