package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/skndash96/lastnight-backend/internal/service"
)

type internalHandler struct {
	uploadSrv *service.UploadService
}

func NewInternalHandler(uploadSrv *service.UploadService) *internalHandler {
	return &internalHandler{
		uploadSrv: uploadSrv,
	}
}

// this handler need not be documented because it is a private handler for workers
func (h *internalHandler) GetUploadRef(c echo.Context) error {
	uploadRefID := c.Param("uploadRefID")
	if uploadRefID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "uploadRefID is required")
	}

	uploadRefIDInt, err := strconv.Atoi(uploadRefID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid uploadRefID")
	}

	ref, err := h.uploadSrv.GetUploadRef(c.Request().Context(), int32(uploadRefIDInt))
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, ref)

	return nil
}
