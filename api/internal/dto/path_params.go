package dto

type TeamPathParams struct {
	TeamID int32 `param:"teamID" validate:"required"`
}

type DocRefPathParams struct {
	TeamPathParams
	DocRefID int32 `param:"docRefID" validate:"required"`
}

type TagPathParams struct {
	TeamPathParams
	TagID int32 `param:"tagID" validate:"required"`
}

type TagValuePathParams struct {
	TagPathParams
	TagValueID int32 `param:"tagValueID" validate:"required"`
}
