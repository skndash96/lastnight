-- name: GetTeamsByUserID :many
SELECT
  t.id AS team_id,
  t.name AS team_name,
  t.domain AS team_domain,
  tm.id AS membership_id,
  tm.role AS user_role,
  tm.joined_at AS user_joined_at
  FROM teams t
  INNER JOIN team_memberships tm ON t.id = tm.team_id
  WHERE tm.user_id = $1;

-- name: GetTeamByDomain :one
SELECT * FROM teams WHERE domain = $1;

-- name: GetTeamMembershipByUserID :one
SELECT * FROM team_memberships WHERE user_id = $1 AND team_id = $2;

-- name: CreateTeamMembership :one
INSERT INTO team_memberships (user_id, team_id, role)
VALUES ($1, $2, $3)
RETURNING *;
