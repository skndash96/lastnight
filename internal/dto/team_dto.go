package dto

import "github.com/skndash96/lastnight-backend/internal/db"

type GetTeamsResponse struct {
	Data []db.GetTeamsByUserIDRow `json:"data"`
}

type JoinTeamResponse struct {
	Data *db.Team `json:"data"`
}
