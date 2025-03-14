package service

import (
	"context"
	"errors"
	"fmt"
	"redditBack/model"
	"redditBack/repository"
)

var (
	ErrInvalidVoteValue = errors.New("vote value must be 1 or -1")
	ErrSelfVote         = errors.New("cannot vote on your own post")
)

type VoteService struct {
	voteRepo  repository.VoteRepository
	postRepo  repository.PostRepository
	userRepo  repository.UserRepository
	cacheRepo repository.CacheRepository
}

func NewVoteService(voteRepo repository.VoteRepository, postRepo repository.PostRepository,
	userRepo repository.UserRepository, cacheRepo repository.CacheRepository) VoteService {
	return VoteService{
		voteRepo:  voteRepo,
		postRepo:  postRepo,
		userRepo:  userRepo,
		cacheRepo: cacheRepo,
	}
}

func (s *VoteService) VotePost(ctx context.Context, postID uint, username string, voteValue int) error {

	if voteValue != 1 && voteValue != -1 && voteValue != 0 {
		return ErrInvalidVoteValue
	}
	fmt.Printf("postID : %u", postID)
	post, err := s.postRepo.FindByID(ctx, postID)
	if err != nil || post == nil {
		return fmt.Errorf("post not found")
	}
	user, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return fmt.Errorf("can't find user with username %s", username)
	}
	fmt.Print(post.UserID)
	if post.UserID == user.ID {
		return ErrSelfVote
	}

	existingVote, err := s.voteRepo.FindByUserAndPost(ctx, user.ID, postID)
	var voteDelta int

	if err == nil && existingVote != nil {
		voteDelta = voteValue - existingVote.VoteValue
		existingVote.VoteValue = voteValue
		err = s.voteRepo.Update(ctx, existingVote)
	} else {

		voteDelta = voteValue
		fmt.Printf("%u %s %s", user.ID, postID, voteValue)
		newVote := &model.Vote{
			UserID:    user.ID,
			PostID:    postID,
			VoteValue: voteValue,

		}
		err = s.voteRepo.Create(ctx, newVote)
	}

	if err != nil {
		return fmt.Errorf("failed to process vote: %w", err)
	}

	err = s.postRepo.UpdateScore(ctx, postID, voteDelta)
	if err != nil {
		return fmt.Errorf("failed to update post score: %w", err)
	}

	s.cacheRepo.InvalidatePostRanking(ctx)
	return nil
}
