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

func newPostService(postRepo repository.PostRepository, userRepo repository.UserRepository) PostService {
	return PostService{postRepo: postRepo, userRepo: userRepo}
}

func (p *PostService) createNewPost(ctx context.Context, post *model.Post) error {
	usernameVal := ctx.Value("user_id")
	username, ok := usernameVal.(string)
	if !ok {
		return errors.New("invalid username type in context")
	}
	user, err := p.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return errors.New("Error in username")
	}
	post.UserID = user.ID
	return p.postRepo.Create(ctx, post)
}

func (p *PostService) editPost(ctx context.Context, post *model.Post) error {
	usernameVal := ctx.Value("user_id")
	username, ok := usernameVal.(string)
	if !ok {
		return errors.New("invalid username type in context")
	}
	user, err := p.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return errors.New("Error in username")
	}
	tempPost, postErr := p.postRepo.FindByID(ctx, post.ID)
	if postErr != nil {
		return errors.New("Post with this ID doesn't exist!")
	}
	if tempPost.UserID != user.ID {
		return errors.New("you don't have authority to edit this code!")
	}

	updatedPost := model.Post{
        ID:     tempPost.ID,
        Title:   post.Title,
        Content: post.Content,
        UserID:  tempPost.UserID, // Preserve original owner
    }
	return p.postRepo.Update(ctx , &updatedPost)

	
}

func (p *PostService) removePost(ctx context.Context, post *model.Post) error {
	usernameVal := ctx.Value("user_id")
	username, ok := usernameVal.(string)
	if !ok {
		return errors.New("invalid username type in context")
	}
	user, err := p.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return errors.New("Error in username")
	}
	tempPost, postErr := p.postRepo.FindByID(ctx, post.ID)
	if postErr != nil {
		return errors.New("Post with this ID doesn't exist!")
	}
	if tempPost.UserID != user.ID {
		return errors.New("you don't have authority to remove this post!")
	}

	return p.postRepo.Delete(ctx ,tempPost.ID )
}
