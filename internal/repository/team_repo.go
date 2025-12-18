package repository

import (
	"context"

	"github.com/skndash96/lastnight-backend/internal/db"
)

type TeamRepository struct {
	q *db.Queries
}

func NewTeamRepository(d db.DBTX) *TeamRepository {
	return &TeamRepository{
		q: db.New(d),
	}
}

func (r *TeamRepository) GetTeamsByUserID(ctx context.Context, id int32) ([]db.GetTeamsByUserIDRow, error) {
	return r.q.GetTeamsByUserID(ctx, id)
}

func (r *TeamRepository) GetTeamByDomain(ctx context.Context, d string) (db.Team, error) {
	return r.q.GetTeamByDomain(ctx, d)
}

func (r *TeamRepository) GetTeamMembershipByUserID(ctx context.Context, user_id, team_id int32) (db.TeamMembership, error) {
	return r.q.GetTeamMembershipByUserID(ctx, db.GetTeamMembershipByUserIDParams{
		UserID: user_id,
		TeamID: team_id,
	})
}

func (r *TeamRepository) CreateTeamMembership(ctx context.Context, user_id, team_id int32, role db.TeamUserRole) (db.TeamMembership, error) {
	return r.q.CreateTeamMembership(ctx, db.CreateTeamMembershipParams{
		UserID: user_id,
		TeamID: team_id,
		Role:   role,
	})
}
