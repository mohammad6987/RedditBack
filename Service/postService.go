package service

import (
	"context"
	"errors"
	"redditBack/model"
	"redditBack/repository"
)

type PostService struct {
	postRepo repository.PostRepository
	userRepo repository.UserRepository
}

func NewPostService(postRepo repository.PostRepository, userRepo repository.UserRepository) PostService {
	return PostService{postRepo: postRepo, userRepo: userRepo}
}

func (p *PostService) CreateNewPost(ctx context.Context, post *model.Post, username string) error {

	user, err := p.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return errors.New("Error in username")
	}
	post.UserID = user.ID
	return p.postRepo.Create(ctx, post)
}

func (p *PostService) EditPost(ctx context.Context, post *model.Post, username string) error {

	user, err := p.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return errors.New("Error in username")
	}
	tempPost, postErr := p.postRepo.FindByID(ctx, post.ID)
	if postErr != nil {
		return errors.New("post not found")
	}
	if tempPost.UserID != user.ID {
		return errors.New("unauthorized to edit post")
	}

	updatedPost := model.Post{
		ID:      tempPost.ID,
		Title:   post.Title,
		Content: post.Content,
		UserID:  tempPost.UserID, // Preserve original owner
	}
	return p.postRepo.Update(ctx, &updatedPost)

}

func (p *PostService) RemovePost(ctx context.Context, post *model.Post, username string) error {

	user, err := p.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return errors.New("Error in username")
	}
	tempPost, postErr := p.postRepo.FindByID(ctx, post.ID)
	if postErr != nil {
		return errors.New("post not found")
	}
	if tempPost.UserID != user.ID {
		return errors.New("unauthorized to edit post")
	}

	return p.postRepo.Delete(ctx, tempPost.ID)
}
