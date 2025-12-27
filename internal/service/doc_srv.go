package service

import (
	"context"
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/skndash96/lastnight-backend/internal/db"
	"github.com/skndash96/lastnight-backend/internal/provider"
	"github.com/skndash96/lastnight-backend/internal/queue"
	"github.com/skndash96/lastnight-backend/internal/repository"
)

type DocumentService struct {
	pool           *pgxpool.Pool
	uploadProvider provider.UploadProvider
	ingestionQ     *queue.IngestionQ
}

func NewDocumentService(uploadProvider provider.UploadProvider, pool *pgxpool.Pool, ingestionQ *queue.IngestionQ) *DocumentService {
	return &DocumentService{
		pool:           pool,
		uploadProvider: uploadProvider,
		ingestionQ:     ingestionQ,
	}
}

type PresignUploadResult struct {
	Url    *url.URL
	Fields map[string]string
}

func (s *DocumentService) GetDoc(ctx context.Context, id int32) (*db.Doc, error) {
	docRepo := repository.NewDocRepository(s.pool)
	ref, err := docRepo.GetDoc(ctx, id)
	if err != nil {
		return nil, NewSrvError(err, SrvErrInternal, "Failed to get doc reference")
	}
	return ref, nil
}

func (s *DocumentService) PresignUpload(ctx context.Context, teamID int32, name, mimeType string, size int64) (*PresignUploadResult, error) {
	// Allow only application/* mime types
	// TODO: Support image/* mime types while combining images to PDF if done so
	if !strings.HasPrefix(mimeType, "application/") {
		return nil, NewSrvError(nil, SrvErrInvalidInput, "Upload presign failed: Only application/* mime types are allowed.")
	}

	key := generateTmpObjectKey(teamID, name)
	url, fields, err := s.uploadProvider.PresignUpload(ctx, key, name, mimeType, size)
	if err != nil {
		// TODO: Handle error properly
		return nil, NewSrvError(err, SrvErrInternal, fmt.Sprintf("failed to presign upload for file %s", name))
	}

	return &PresignUploadResult{
		Url:    url,
		Fields: fields,
	}, nil
}

func (s *DocumentService) CommitUpload(ctx context.Context, teamID, userID int32, tmpKey, name, mime string, tags [][]int32) error {
	info, err := s.uploadProvider.GetUploadInfo(ctx, tmpKey)
	if err != nil {
		return NewSrvError(err, SrvErrInternal, fmt.Sprintf("failed to get upload info for %s", tmpKey))
	}
	// TODO: Validate blob

	newKey := convertTmpKey(tmpKey)
	if newKey == "" {
		return NewSrvError(nil, SrvErrInvalidInput, "Upload completion failed: invalid key")
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return NewSrvError(err, SrvErrInternal, "failed to begin transaction")
	}
	defer tx.Rollback(ctx)

	docRepo := repository.NewDocRepository(tx)

	// (hash, size) duplication check happens here
	doc, err := docRepo.GetOrCreateDoc(ctx, newKey, mime, info.Size, info.SHA256)
	if err != nil {
		return NewSrvError(err, SrvErrInternal, fmt.Sprintf("failed to create doc for %s", newKey))
	}

	// TODO: Add constraint UNIQUE (doc_id, team_id, user_id, name)
	// or idempotency in this route
	docRef, err := docRepo.CreateDocRef(ctx, doc.ID, teamID, userID, name)
	if err != nil {
		return NewSrvError(err, SrvErrInternal, fmt.Sprintf("failed to create doc reference for %s", newKey))
	}

	for _, tag := range tags {
		if err := docRepo.CreateDocRefTag(ctx, docRef.ID, tag[0], tag[1]); err != nil {
			return NewSrvError(err, SrvErrInternal, fmt.Sprintf("failed to create doc tag for %s", newKey))
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return NewSrvError(err, SrvErrInternal, "failed to complete doc")
	}

	if doc.Created == false {
		fmt.Printf("Doc already exists, using duplicate: %s\n", doc.StorageKey)
		if err := s.uploadProvider.DeleteObject(ctx, tmpKey); err != nil {
			// it's not a fatal error, so we can continue
			fmt.Printf("failed to delete duplicate upload %s: %v\n", tmpKey, err)
		}

		return nil
	}

	err = s.uploadProvider.MoveObject(ctx, newKey, tmpKey)
	if err != nil {
		// TODO: Retry and Failure handling
		return NewSrvError(err, SrvErrInternal, fmt.Sprintf("failed to move upload from %s to %s", tmpKey, newKey))
	}

	job := &queue.IngestionJob{
		ID:    doc.ID,
		RefID: docRef.ID,
	}

	if err := s.ingestionQ.Enqueue(ctx, job); err != nil {
		// TODO: Retry and Failure handling
		return NewSrvError(err, SrvErrInternal, fmt.Sprintf("failed to start ingest job for %s", doc.StorageKey))
	}

	return nil
}

func (s *DocumentService) UpdateDocStatus(ctx context.Context, docID int32, status db.DocProcStatus) error {
	docRepo := repository.NewDocRepository(s.pool)

	if err := docRepo.UpdateDocStatus(ctx, docID, status); err != nil {
		return NewSrvError(err, SrvErrInternal, "failed to update document status")
	}

	return nil
}

func (s *DocumentService) UpdateDocRefTags(ctx context.Context, teamID, userID, docRefID int32, tags [][]int32) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return NewSrvError(err, SrvErrInternal, "failed to start transaction")
	}
	defer tx.Rollback(ctx)

	docRepo := repository.NewDocRepository(tx)

	if err := docRepo.DeleteAllDocRefTags(ctx, tx, docRefID); err != nil {
		return NewSrvError(err, SrvErrInternal, fmt.Sprintf("failed to delete doc tags for %d", docRefID))
	}

	for _, tag := range tags {
		if err := docRepo.CreateDocRefTag(ctx, docRefID, tag[0], tag[1]); err != nil {
			return NewSrvError(err, SrvErrInternal, fmt.Sprintf("failed to create doc tag for %d", docRefID))
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return NewSrvError(err, SrvErrInternal, "failed to complete doc")
	}

	return nil
}

func generateTmpObjectKey(teamID int32, originalName string) string {
	ext := path.Ext(originalName)
	if len(ext) > 10 {
		ext = ""
	}

	id := uuid.NewString()

	return fmt.Sprintf("tmp/team_%d/uploads/%s%s", teamID, id, strings.ToLower(ext))
}

func convertTmpKey(key string) string {
	if strings.HasPrefix(key, "tmp/") {
		return strings.Replace(key, "tmp/", "", 1)
	}
	return ""
}
