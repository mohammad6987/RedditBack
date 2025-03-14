package repository

import (
	"context"
	"errors"
	"redditBack/model"

	"gorm.io/gorm"
)

type VoteRepository interface {
	Create(ctx context.Context, post *model.Vote) error
	FindByUserAndPost(ctx context.Context, userID uint, postID uint) (*model.Vote, error)
	Update(ctx context.Context, vote *model.Vote) error
	Delete(ctx context.Context, id uint) error
}

type VoteRepositoryImp struct {
	db *gorm.DB
}

func NewVoteRepository(db *gorm.DB) VoteRepositoryImp {
	return VoteRepositoryImp{db: db}
}

func (r *VoteRepositoryImp) Create(ctx context.Context, vote *model.Vote) error {
	return r.db.WithContext(ctx).Create(vote).Error
}

func (r *VoteRepositoryImp) FindByUserAndPost(ctx context.Context, userID uint, postID uint) (*model.Vote, error) {
	var vote model.Vote
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND post_id = ?", userID, postID).
		First(&vote).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &vote, err
}
func (r *VoteRepositoryImp) Update(ctx context.Context, vote *model.Vote) error {
	return r.db.WithContext(ctx).Save(vote).Error
}

func (r *VoteRepositoryImp) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&model.Post{}, id)
	if result.RowsAffected == 0 {
		return errors.New("No post record with this ID!")
	}
	return result.Error
}
