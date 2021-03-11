package model

import "time"

type FscStatus struct {
	FSCId                   string        `json:"fsc_id"`
	Site                    string        `json:"site"`
	RunningTime             time.Duration `json:"running_time,omitempty"`
	CurrentFileConnections  int64         `json:"current_file_connections"`
	CurrentAdminConnections int64         `json:"current_admin_connections"`
}
