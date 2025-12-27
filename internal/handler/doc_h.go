package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/skndash96/lastnight-backend/internal/auth"
	"github.com/skndash96/lastnight-backend/internal/dto"
	"github.com/skndash96/lastnight-backend/internal/service"
)

type docHandler struct {
	docSrv *service.DocumentService
}

func NewDocHandler(docSrv *service.DocumentService) *docHandler {
	return &docHandler{
		docSrv: docSrv,
	}
}

// @Summary Create pre-signed request
// @Description Create a pre-signed request for uploading files to S3 via POST policy
// @Tags Document
// @Accept json
// @Param teamID path string true "Team ID"
// @Param upload_request body dto.PresignUploadBody true "Presign request"
// @Produce json
// @Success 201 {object} dto.PresignUploadResponse
// @Failure default {object} dto.ErrorResponse
// @Router /api/teams/{teamID}/uploads/presign [post]
func (h *docHandler) PresignUpload(c echo.Context) error {
	session, ok := auth.GetSession(c)
	if !ok {
		return echo.ErrUnauthorized
	}

	v := dto.PresignUploadRequest{}
	if err := c.Bind(&v); err != nil {
		return err
	}

	if err := c.Validate(&v); err != nil {
		return err
	}

	result, err := h.docSrv.PresignUpload(c.Request().Context(), session.TeamID, v.Name, v.MimeType, v.Size)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, &dto.PresignUploadResponse{
		Url:    result.Url.String(),
		Fields: result.Fields,
	})

	return nil
}

// @Summary Commit upload
// @Description Call this route after client-side uploading to the bucket via POST policy. Processes uploaded file.
// @Tags Document
// @Accept json
// @Param teamID path string true "Team ID"
// @Param upload_request body dto.CommitUploadBody true "Commit upload request"
// @Produce json
// @Success 200
// @Failure default {object} dto.ErrorResponse
// @Router /api/teams/{teamID}/uploads/commit [post]
func (h *docHandler) CommitUpload(c echo.Context) error {
	session, ok := auth.GetSession(c)
	if !ok {
		return echo.ErrUnauthorized
	}

	v := dto.CommitUploadRequest{}
	if err := c.Bind(&v); err != nil {
		return err
	}

	if err := c.Validate(&v); err != nil {
		return err
	}

	tags := make([][]int32, len(v.Tags))
	for i, tag := range v.Tags {
		tags[i] = []int32{tag.KeyID, tag.ValueID}
	}

	err := h.docSrv.CommitUpload(c.Request().Context(), session.TeamID, session.UserID, v.Key, v.Name, v.MimeType, tags)
	if err != nil {
		return err
	}

	c.NoContent(http.StatusCreated)

	return nil
}

// @Summary Update DocRef Tags
// @Description Replace tags of a document reference.
// @Tags Document
// @Accept json
// @Param teamID path string true "Team ID"
// @Param uploadID path string true "Upload ID"
// @Param upload_request body dto.UpdateDocRefBody true "Replace tags request"
// @Produce json
// @Success 200
// @Failure default {object} dto.ErrorResponse
// @Router /api/teams/{teamID}/upload-refs/{docRefID}/tags [put]
func (h *docHandler) UpdateDocRefTags(c echo.Context) error {
	session, ok := auth.GetSession(c)
	if !ok {
		return echo.ErrUnauthorized
	}

	v := dto.UpdateDocRefTagsRequest{}
	if err := c.Bind(&v); err != nil {
		return err
	}

	if err := c.Validate(&v); err != nil {
		return err
	}

	tags := make([][]int32, len(v.Tags))
	for i, tag := range v.Tags {
		tags[i] = []int32{tag.KeyID, tag.ValueID}
	}

	err := h.docSrv.UpdateDocRefTags(c.Request().Context(), session.TeamID, session.UserID, v.DocRefID, tags)
	if err != nil {
		return err
	}

	c.NoContent(http.StatusCreated)

	return nil
}
