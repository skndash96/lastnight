package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/skndash96/lastnight-backend/internal/db"
	"github.com/skndash96/lastnight-backend/internal/service"
)

// this handler need not be swagger documented because it is a private handler for workers
type internalHandler struct {
	docSrv *service.DocumentService
}

func NewInternalHandler(docSrv *service.DocumentService) *internalHandler {
	return &internalHandler{
		docSrv: docSrv,
	}
}

func (h *internalHandler) GetDoc(c echo.Context) error {
	docID := c.Param("docID")
	if docID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "docID is required")
	}

	docIDInt, err := strconv.Atoi(docID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid docRefID")
	}

	doc, err := h.docSrv.GetDoc(c.Request().Context(), int32(docIDInt))
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, doc)

	return nil
}

func (h *internalHandler) UpdateDocProcStatus(c echo.Context) error {
	docID := c.Param("docID")
	if docID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "docID is required")
	}

	docIDInt, err := strconv.Atoi(docID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid docID")
	}

	var v struct {
		Status db.DocProcStatus `json:"status"`
	}

	if err := c.Bind(&v); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	switch v.Status {
	case db.DocProcStatusPending, db.DocProcStatusCompleted, db.DocProcStatusFailed:
		// pass
	default:
		return echo.NewHTTPError(http.StatusBadRequest, "invalid status")
	}

	err = h.docSrv.UpdateDocStatus(c.Request().Context(), int32(docIDInt), v.Status)
	if err != nil {
		return err
	}

	c.NoContent(http.StatusOK)

	return nil
}
