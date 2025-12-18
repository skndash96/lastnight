package service

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/skndash96/lastnight-backend/internal/db"
	"github.com/skndash96/lastnight-backend/internal/helpers"
	"github.com/skndash96/lastnight-backend/internal/repository"
)

type TagService struct {
	db *pgxpool.Pool
}

func NewTagService(p *pgxpool.Pool) *TagService {
	return &TagService{
		db: p,
	}
}

func (s *TagService) ListTags(ctx context.Context, teamID int32) ([]db.Tag, error) {
	tagRepo := repository.NewTagRepo(s.db)
	tags, err := tagRepo.ListTags(ctx, teamID)
	if err != nil {
		return nil, err
	}
	return tags, nil
}

func (s *TagService) ListTagValues(ctx context.Context, tagID int32) ([]db.TagValue, error) {
	tagValueRepo := repository.NewTagRepo(s.db)
	tagValues, err := tagValueRepo.ListTagValues(ctx, tagID)
	if err != nil {
		return nil, err
	}
	return tagValues, nil
}

func (s *TagService) CreateTag(ctx context.Context, teamID int32, name string, dataType db.TagDataType) (*db.Tag, error) {
	tagRepo := repository.NewTagRepo(s.db)
	tag, err := tagRepo.CreateTag(ctx, teamID, name, dataType)
	if err != nil {
		if helpers.IsUniqueViolation(err) {
			return nil, NewSrvError(err, SrvErrInvalidInput, "tag value already exists")
		}
		return nil, err
	}
	return &tag, nil
}

func (s *TagService) CreateTagValue(ctx context.Context, tagID int32, value string) (*db.TagValue, error) {
	tagValueRepo := repository.NewTagRepo(s.db)
	tagValue, err := tagValueRepo.CreateTagValue(ctx, tagID, value)
	if err != nil {
		if helpers.IsUniqueViolation(err) {
			return nil, NewSrvError(err, SrvErrInvalidInput, "tag value already exists")
		}
		return nil, err
	}
	return &tagValue, nil
}

func (s *TagService) UpdateTag(ctx context.Context, tagID int32, name string) (*db.Tag, error) {
	tagRepo := repository.NewTagRepo(s.db)
	tag, err := tagRepo.UpdateTag(ctx, tagID, name)
	if err != nil {
		if helpers.IsUniqueViolation(err) {
			return nil, NewSrvError(err, SrvErrInvalidInput, "tag value already exists")
		}
		return nil, err
	}
	return &tag, nil
}

func (s *TagService) DeleteTag(ctx context.Context, tagID int32) (*db.Tag, error) {
	tagRepo := repository.NewTagRepo(s.db)
	tag, err := tagRepo.DeleteTag(ctx, tagID)
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

func (s *TagService) DeleteTagValue(ctx context.Context, tagValueID int32) (*db.TagValue, error) {
	tagValueRepo := repository.NewTagRepo(s.db)
	tagValue, err := tagValueRepo.DeleteTagValue(ctx, tagValueID)
	if err != nil {
		return nil, err
	}
	return &tagValue, nil
}
