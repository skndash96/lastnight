package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/skndash96/lastnight-backend/internal/dto"
	"github.com/skndash96/lastnight-backend/internal/service"
)

type tagHandler struct {
	tagSrv *service.TagService
}

func NewTagHandler(s *service.TagService) *tagHandler {
	return &tagHandler{
		tagSrv: s,
	}
}

// GetTags retrieves the tags of a team.
// @Summary Get Tags
// @Tags Tag
// @Description Get the tags of a team
// @Param teamID path string true "Team ID"
// @Produce json
// @Success 200 {object} dto.GetTagsResponse
// @Failure default {object} dto.ErrorResponse
// @Router /api/teams/{teamID}/tags [get]
func (h *tagHandler) ListTags(c echo.Context) error {
	v := new(dto.GetTagsRequest)
	if err := c.Bind(v); err != nil {
		return err
	}

	tags, err := h.tagSrv.ListTags(c.Request().Context(), v.TeamID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.GetTagsResponse{
		Data: tags,
	})
}

// GetTagValues retrieves the values of a tag.
// @Summary Get Tag Values
// @Tags Tag
// @Description Get the values of a tag
// @Param teamID path string true "Team ID"
// @Param tagID path string true "Tag ID"
// @Produce json
// @Success 200 {object} dto.GetTagValuesResponse
// @Failure default {object} dto.ErrorResponse
// @Router /api/teams/{teamID}/tags/{tagID}/values [get]
func (h *tagHandler) ListTagValues(c echo.Context) error {
	v := new(dto.GetTagValuesRequest)
	if err := c.Bind(v); err != nil {
		return err
	}

	values, err := h.tagSrv.ListTagValues(c.Request().Context(), v.TagID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.GetTagValuesResponse{
		Data: values,
	})
}

// @Summary New Tag
// @Tags Tag
// @Description Create a new tag
// @Param teamID path string true "Team ID"
// @Param tag body dto.CreateTagBody true "Tag"
// @Produce json
// @Success 201 {object} dto.CreateTagResponse
// @Failure default {object} dto.ErrorResponse
// @Router /api/teams/{teamID}/tags [post]
func (h *tagHandler) CreateTag(c echo.Context) error {
	v := new(dto.CreateTagRequest)
	if err := c.Bind(v); err != nil {
		return err
	}

	if err := c.Validate(v); err != nil {
		return err
	}

	tag, err := h.tagSrv.CreateTag(c.Request().Context(), v.TeamID, v.Name, v.DataType)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, dto.CreateTagResponse{
		Data: tag,
	})
}

// @Summary Update Tag
// @Tags Tag
// @Description Update a tag
// @Param teamID path string true "Team ID"
// @Param tagID path string true "Tag ID"
// @Param tag body dto.UpdateTagBody true "Tag"
// @Produce json
// @Success 200 {object} dto.UpdateTagResponse
// @Failure default {object} dto.ErrorResponse
// @Router /api/teams/{teamID}/tags/{tagID} [put]
func (h *tagHandler) UpdateTag(c echo.Context) error {
	v := new(dto.UpdateTagRequest)
	if err := c.Bind(v); err != nil {
		return err
	}

	if err := c.Validate(v); err != nil {
		return err
	}

	tag, err := h.tagSrv.UpdateTag(c.Request().Context(), v.TagID, v.Name)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.UpdateTagResponse{
		Data: tag,
	})
}

// @Summary Delete Tag
// @Tags Tag
// @Description Delete a tag
// @Param teamID path string true "Team ID"
// @Param tagID path string true "Tag ID"
// @Produce json
// @Success 200 {object} dto.DeleteTagResponse
// @Failure default {object} dto.ErrorResponse
// @Router /api/teams/{teamID}/tags/{tagID} [delete]
func (h *tagHandler) DeleteTag(c echo.Context) error {
	v := new(dto.DeleteTagRequest)

	if err := c.Bind(v); err != nil {
		return err
	}

	if err := c.Validate(v); err != nil {
		return err
	}

	tag, err := h.tagSrv.DeleteTag(c.Request().Context(), v.TagID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.DeleteTagResponse{
		Data: tag,
	})
}

// @Summary Create Tag Value
// @Tags Tag
// @Description Create a new tag value
// @Param teamID path string true "Team ID"
// @Param tagID path string true "Tag ID"
// @Param value body dto.CreateTagValueBody true "Value"
// @Produce json
// @Success 201 {object} dto.CreateTagValueResponse
// @Failure default {object} dto.ErrorResponse
// @Router /api/teams/{teamID}/tags/{tagID}/values [post]
func (h *tagHandler) CreateTagValue(c echo.Context) error {
	v := new(dto.CreateTagValueRequest)

	if err := c.Bind(v); err != nil {
		return err
	}

	if err := c.Validate(v); err != nil {
		return err
	}

	value, err := h.tagSrv.CreateTagValue(c.Request().Context(), v.TagID, v.Value)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, dto.CreateTagValueResponse{
		Data: value,
	})
}

// @Summary Delete Tag Value
// @Tags Tag
// @Description Delete a tag value
// @Param teamID path string true "Team ID"
// @Param tagID path string true "Tag ID"
// @Param valueID path string true "Value ID"
// @Produce json
// @Success 200 {object} dto.DeleteTagValueResponse
// @Failure default {object} dto.ErrorResponse
// @Router /api/teams/{teamID}/tags/{tagID}/values/{valueID} [delete]
func (h *tagHandler) DeleteTagValue(c echo.Context) error {
	v := new(dto.DeleteTagValueRequest)

	if err := c.Bind(v); err != nil {
		return err
	}

	if err := c.Validate(v); err != nil {
		return err
	}

	value, err := h.tagSrv.DeleteTagValue(c.Request().Context(), v.TagValueID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.DeleteTagValueResponse{
		Data: value,
	})
}
