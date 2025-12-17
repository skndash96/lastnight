package dto

import "github.com/skndash96/lastnight-backend/internal/db"

type GetTeamResponse struct {
	Data *db.Team `json:"data"`
}

type JoinTeamResponse struct {
	Data *db.Team `json:"data"`
}
