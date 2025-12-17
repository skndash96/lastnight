package repository

import (
	"context"

	"github.com/skndash96/lastnight-backend/internal/db"
)

type TeamRepository interface {
	GetTeamProfileByUserID(ctx context.Context, id int32) (db.GetTeamProfileByUserIDRow, error)
	GetTeamByDomain(ctx context.Context, d string) (db.Team, error)
	CreateTeamMembership(ctx context.Context, user_id, team_id int32, role db.TeamUserRole) (db.TeamMembership, error)
}

type teamRepository struct {
	q *db.Queries
}

func NewTeamRepository(d db.DBTX) TeamRepository {
	return &teamRepository{
		q: db.New(d),
	}
}

func (r *teamRepository) GetTeamProfileByUserID(ctx context.Context, id int32) (db.GetTeamProfileByUserIDRow, error) {
	return r.q.GetTeamProfileByUserID(ctx, id)
}

func (r *teamRepository) GetTeamByDomain(ctx context.Context, d string) (db.Team, error) {
	return r.q.GetTeamByDomain(ctx, d)
}

func (r *teamRepository) CreateTeamMembership(ctx context.Context, user_id, team_id int32, role db.TeamUserRole) (db.TeamMembership, error) {
	return r.q.CreateTeamMembership(ctx, db.CreateTeamMembershipParams{
		UserID: user_id,
		TeamID: team_id,
		Role:   role,
	})
}
