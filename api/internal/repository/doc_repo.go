package repository

import (
	"context"

	"github.com/skndash96/lastnight-backend/internal/db"
)

type docRepository struct {
	q *db.Queries
}

func NewDocRepository(d db.DBTX) *docRepository {
	return &docRepository{
		q: db.New(d),
	}
}

func (r *docRepository) GetOrCreateDoc(ctx context.Context, key, mimeType string, size int64, sha256 string) (*db.GetOrCreateDocRow, error) {
	doc, err := r.q.GetOrCreateDoc(ctx, db.GetOrCreateDocParams{
		StorageKey:   key,
		FileMimeType: mimeType,
		FileSize:     size,
		FileSha256:   sha256,
	})
	if err != nil {
		return nil, NewRepoError(err, RepoErrInternal, "Failed to create doc")
	}
	return &doc, nil
}

func (r *docRepository) GetDoc(ctx context.Context, id int32) (*db.Doc, error) {
	ref, err := r.q.GetDoc(ctx, id)
	if err != nil {
		return nil, NewRepoError(err, RepoErrInternal, "Failed to get doc")
	}
	return &ref, nil
}

func (r *docRepository) CreateDocRef(ctx context.Context, docID int32, teamID int32, userID int32, name string) (*db.DocRef, error) {
	ref, err := r.q.CreateDocRef(ctx, db.CreateDocRefParams{
		FileName: name,
		DocID:    docID,
		TeamID:   teamID,
		UserID:   userID,
	})
	if err != nil {
		return nil, NewRepoError(err, RepoErrInternal, "Failed to create doc reference")
	}
	return &ref, nil
}

func (r *docRepository) CreateDocRefTag(ctx context.Context, docRefID int32, keyID int32, valueID int32) error {
	err := r.q.CreateDocRefTags(ctx, db.CreateDocRefTagsParams{
		DocRefID: docRefID,
		KeyID:    keyID,
		ValueID:  valueID,
	})

	if err != nil {
		return NewRepoError(err, RepoErrInternal, "Failed to create doc tag")
	}

	return nil
}

func (r *docRepository) DeleteAllDocRefTags(ctx context.Context, tx db.DBTX, docRefID int32) error {
	err := r.q.DeleteAllDocRefTags(ctx, docRefID)

	if err != nil {
		return NewRepoError(err, RepoErrInternal, "Failed to delete doc ref tags")
	}

	return nil
}
