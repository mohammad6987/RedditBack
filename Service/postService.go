package service

import (
	"context"
	"errors"
	"log"
	"redditBack/model"
	"redditBack/repository"
	"time"
)

type PostService struct {
	postRepo  repository.PostRepository
	userRepo  repository.UserRepository
	cacheRepo repository.CacheRepository
	voteRepo  repository.VoteRepository
}

func NewPostService(postRepo repository.PostRepository, userRepo repository.UserRepository, cacheRepo repository.CacheRepository, voteRepo repository.VoteRepository) PostService {
	return PostService{
		postRepo:  postRepo,
		userRepo:  userRepo,
		cacheRepo: cacheRepo,
		voteRepo:  voteRepo}
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
		UserID:  tempPost.UserID,
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
	err2 := p.voteRepo.Delete(ctx, tempPost.ID)
	if err2 != nil {
		log.Println("some erros in deleting votes, but igonre it")
	}
	return p.postRepo.Delete(ctx, tempPost.ID)

}

func (p *PostService) GetTopPosts(ctx context.Context, timeRange string) ([]*model.Post, error) {

	var startTime time.Time
	now := time.Now()

	cachedPosts, err := p.cacheRepo.GetTopPosts(ctx, timeRange)
	if err == nil && len(cachedPosts) > 0 {
		return cachedPosts, nil
	}
	switch timeRange {
	case "day":
		startTime = now.Add(-24 * time.Hour)
	case "week":
		startTime = now.Add(-7 * 24 * time.Hour)
	case "month":
		startTime = now.Add(-30 * 24 * time.Hour)
	case "all":
		startTime = time.Time{}
	default:
		return nil, errors.New("invalid time range")
	}

	posts, err := p.postRepo.FindTopPosts(ctx, startTime)
	if err != nil {
		return nil, err
	}
	if err := p.cacheRepo.CacheTopPosts(ctx, timeRange, posts); err != nil {
		log.Printf("Failed to cache posts: %v", err)
		return nil, err
	}

	return posts, nil
}
