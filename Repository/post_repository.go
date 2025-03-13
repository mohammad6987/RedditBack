package repository

import (
	"context"
	"redditBack/model"
)

type PostRepository interface {
	Create(ctx context.Context, post *model.Post) error
	FindByID(ctx context.Context, id uint) (*model.Post, error)
	Update(ctx context.Context, post *model.Post) error
	Delete(ctx context.Context, id uint) error
}