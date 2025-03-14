package repository

import (
	"context"
	"errors"
	"redditBack/model"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type VoteRepository interface {
	Create(ctx context.Context, post *model.Vote) error
	FindByID(ctx context.Context, id uint) (*model.Vote, error)
	Update(ctx context.Context, vote *model.Vote) error
	Delete(ctx context.Context, id uint) error
}

type VoteRepositoryImp struct {
	db  *gorm.DB
	rdb *redis.Client
}

func newVoteRepository(db *gorm.DB) *VoteRepositoryImp {
	return &VoteRepositoryImp{db: db}
}

func (r *VoteRepositoryImp) Create(ctx context.Context, vote *model.Vote) error {
	return r.db.WithContext(ctx).Create(vote).Error
}

func (r *VoteRepositoryImp) FindByID(ctx context.Context, id uint) (*model.Vote, error) {
	var vote model.Vote
	err := r.db.WithContext(ctx).First(&vote, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &vote, err
}

func (r *VoteRepositoryImp) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&model.Post{}, id)
	if result.RowsAffected == 0 {
		return errors.New("No post record with this ID!")
	}
	return result.Error
}
