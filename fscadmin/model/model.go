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

type FscStatus struct {
	Status                  Status        `json:"status"`
	FSCId                   string        `json:"fsc_id"`
	Site                    string        `json:"site"`
	RunningTime             time.Duration `json:"running_time,omitempty"`
	CurrentFileConnections  int64         `json:"current_file_connections"`
	CurrentAdminConnections int64         `json:"current_admin_connections"`
}
