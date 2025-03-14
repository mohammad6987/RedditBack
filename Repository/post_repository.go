package repository

import (
	"context"
	"errors"
	"redditBack/model"

	"gorm.io/gorm"
)

type PostRepository interface {
	Create(ctx context.Context, post *model.Post) error
	FindByID(ctx context.Context, id uint) (*model.Post, error)
	Update(ctx context.Context, post *model.Post) error
	Delete(ctx context.Context, id uint) error
	UpdateScore(ctx context.Context, postID uint, scoreDelta int) error
}

type PostRepositoryImpl struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) PostRepositoryImpl {
	return PostRepositoryImpl{db: db}
}

func (r *PostRepositoryImpl) Create(ctx context.Context, post *model.Post) error {
	return r.db.WithContext(ctx).Create(post).Error
}

func (r *PostRepositoryImpl) FindByID(ctx context.Context, id uint) (*model.Post, error) {
	var post model.Post
	err := r.db.WithContext(ctx).First(&post, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &post, err
}

func (r *PostRepositoryImpl) Update(ctx context.Context, post *model.Post) error {
	result := r.db.WithContext(ctx).Model(&model.Post{}).
		Where("ID = ?", post.ID).
		Updates(post)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("No post record with this ID!")
	}

	return nil
}

func (r *PostRepositoryImpl) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&model.Post{}, id)
	if result.RowsAffected == 0 {
		return errors.New("No post record with this ID!")
	}
	return result.Error
}

func (r *PostRepositoryImpl) UpdateScore(ctx context.Context, postID uint, scoreDelta int) error {
	result := r.db.WithContext(ctx).
		Model(&model.Post{}).
		Where("id = ?", postID).
		Update("cached_score", gorm.Expr("cached_score + ?", scoreDelta))

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("post not found")
	}

	return nil
}
