-- name: GetTeamProfileByUserID :one
SELECT
  t.id AS team_id,
  t.name AS team_name,
  tm.role AS user_role,
  t.created_at AS team_created_at
  FROM teams t
  INNER JOIN team_memberships tm ON t.id = tm.team_id
  WHERE tm.user_id = $1
  LIMIT 1;

-- name: GetTeamByDomain :one
SELECT * FROM teams WHERE domain = $1;

-- name: CreateTeamMembership :one
INSERT INTO team_memberships (user_id, team_id, role)
VALUES ($1, $2, $3)
RETURNING *;
