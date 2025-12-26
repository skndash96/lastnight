package dto

import (
	"github.com/skndash96/lastnight-backend/internal/db"
)

// TODO: Refactor DTO so that it does NOT contain any database-specific types

// ------ body ------
type Tag struct {
	KeyID   int32 `json:"key_id" validate:"required"`
	ValueID int32 `json:"value_id" validate:"required"`
}

type UpdateFiltersBody struct {
	Filters []Tag `json:"filters" validate:"required"`
}

type CreateTagKeyBody struct {
	Name     string         `json:"name" validate:"required,min=2,max=100"`
	DataType db.TagDataType `json:"data_type" validate:"required,min=2,max=100"`
}

type UpdateTagKeyBody struct {
	Name string `json:"name" validate:"required,min=2,max=100"`
	// data type NOT allowed
}

type CreateTagValueBody struct {
	Value string `json:"value" validate:"required,min=2,max=100"`
}

// ------ request ------
type ListFiltersRequest struct {
	TeamPathParams
}

type UpdateFiltersRequest struct {
	TeamPathParams
	UpdateFiltersBody
}

type CreateTagKeyRequest struct {
	TeamPathParams
	CreateTagKeyBody
}

type UpdateTagKeyRequest struct {
	TagPathParams
	UpdateTagKeyBody
}

type DeleteTagKeyRequest struct {
	TagPathParams
}

type GetTagValuesRequest struct {
	TagPathParams
}

type CreateTagValueRequest struct {
	TagPathParams
	CreateTagValueBody
}

type DeleteTagValueRequest struct {
	TagValuePathParams
}

// ------ response ------
type ListFiltersResponse struct {
	Data []db.Tag `json:"data"`
}

type CreateTagKeyResponse struct {
	Data *db.TagKey `json:"data"`
}

type UpdateTagKeyResponse struct {
	Data *db.TagKey `json:"data"`
}

type DeleteTagKeyResponse struct {
	Data *db.TagKey `json:"data"`
}

type CreateTagValueResponse struct {
	Data *db.TagValue `json:"data"`
}

type UpdateTagValueResponse struct {
	Data *db.TagValue `json:"data"`
}

type DeleteTagValueResponse struct {
	Data *db.TagValue `json:"data"`
}
