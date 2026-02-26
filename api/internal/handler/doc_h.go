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

// @Summary Update DocRef Tags
// @Description Replace tags of a document reference.
// @Tags Document
// @Accept json
// @Param teamID path string true "Team ID"
// @Param docID path string true "Document ID"
// @Param doc_ref_tags_request body dto.UpdateDocRefTagsBody true "Replace tags request"
// @Produce json
// @Success 200
// @Failure default {object} dto.ErrorResponse
// @Router /api/teams/{teamID}/refs/{docRefID}/tags [put]
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
