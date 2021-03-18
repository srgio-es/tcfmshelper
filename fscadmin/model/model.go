package model

import "time"

type Status string

const (
	STATUS_OK Status = "OK"
	STATUS_KO Status = "KO"
)

type Error struct {
	Status  Status `json:"status"`
	Message string `json:"message"`
}

type version struct {
	Version   string `json:"version"`
	BuildDate string `json:"build_date,omitempty"`
}

type FSCVersion struct {
	FmsServerCache     version `json:"fms_server_cache,omitempty"`
	FmsUtil            version `json:"fms_util,omitempty"`
	FscJavaClientProxy version `json:"fsc_java_client_proxy,omitempty"`
}

type FscStatus struct {
	Status                  Status        `json:"status"`
	FSCId                   string        `json:"fsc_id,omitempty"`
	Site                    string        `json:"site,omitempty"`
	RunningTime             time.Duration `json:"running_time,omitempty"`
	CurrentFileConnections  int64         `json:"current_file_connections,omitempty"`
	CurrentAdminConnections int64         `json:"current_admin_connections,omitempty"`
}
