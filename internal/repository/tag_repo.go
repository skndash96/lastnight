package repository

import (
	"context"

	"github.com/skndash96/lastnight-backend/internal/db"
)

type TagRepo struct {
	q *db.Queries
}

func NewTagRepo(d db.DBTX) *TagRepo {
	return &TagRepo{
		q: db.New(d),
	}
}

func (r *TagRepo) CreateTag(ctx context.Context, teamID int32, tagName string, dataType db.TagDataType) (db.Tag, error) {
	return r.q.CreateTag(ctx, db.CreateTagParams{
		TeamID:   teamID,
		Name:     tagName,
		DataType: dataType,
	})
}

func (r *TagRepo) CreateTagValue(ctx context.Context, tagID int32, value string) (db.TagValue, error) {
	return r.q.CreateTagValue(ctx, db.CreateTagValueParams{
		TagID: tagID,
		Value: value,
	})
}

func (r *TagRepo) DeleteTag(ctx context.Context, tagID int32) (db.Tag, error) {
	return r.q.DeleteTag(ctx, tagID)
}

func (r *TagRepo) DeleteTagValue(ctx context.Context, tagValueID int32) (db.TagValue, error) {
	return r.q.DeleteTagValue(ctx, tagValueID)
}

func (r *TagRepo) ListTags(ctx context.Context, teamID int32) ([]db.Tag, error) {
	return r.q.ListTags(ctx, teamID)
}

func (r *TagRepo) ListTagValues(ctx context.Context, tagID int32) ([]db.TagValue, error) {
	return r.q.ListTagValues(ctx, tagID)
}

func (r *TagRepo) UpdateTag(ctx context.Context, tagID int32, tagName string) (db.Tag, error) {
	return r.q.UpdateTag(ctx, db.UpdateTagParams{
		ID:   tagID,
		Name: tagName,
	})
}
