package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"redditBack/model"
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheRepository interface {
	CacheTopPosts(ctx context.Context, timeRange string, posts []*model.Post) error
	GetTopPosts(ctx context.Context, timeRange string) ([]*model.Post, error)
	InvalidatePostRanking(ctx context.Context) error
	CachePost(ctx context.Context, post *model.Post) error
	GetPost(ctx context.Context, postID uint) (*model.Post, error)
	InvalidateToken(ctx context.Context, token string, expiration time.Duration) error
	IsTokenInvalid(ctx context.Context, token string) (bool, error)
}

type RedisCacheRepository struct {
	client *redis.Client
}

func NewRedisCacheRepository(client *redis.Client) RedisCacheRepository {
	return RedisCacheRepository{client: client}
}

func (r *RedisCacheRepository) CacheTopPosts(ctx context.Context, timeRange string, posts []*model.Post) error {
	pipe := r.client.TxPipeline()

	rankingKey := fmt.Sprintf("posts:ranking:%s", timeRange)

	postsKey := "posts:details"

	for _, post := range posts {
		pipe.ZAdd(ctx, rankingKey, redis.Z{
			Score:  float64(post.CachedScore),
			Member: post.ID,
		})

		postJson, _ := json.Marshal(post)
		pipe.HSet(ctx, postsKey, fmt.Sprintf("%d", post.ID), postJson)
	}

	expiration := getExpiration(timeRange)
	pipe.Expire(ctx, rankingKey, expiration)
	pipe.Expire(ctx, postsKey, 24*time.Hour)

	_, err := pipe.Exec(ctx)
	return err
}

func (r *RedisCacheRepository) GetTopPosts(ctx context.Context, timeRange string) ([]*model.Post, error) {
	rankingKey := fmt.Sprintf("posts:ranking:%s", timeRange)
	postsKey := "posts:details"

	ids, err := r.client.ZRevRange(ctx, rankingKey, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	var posts []*model.Post
	for _, idStr := range ids {
		postJson, err := r.client.HGet(ctx, postsKey, idStr).Result()
		if err != nil {
			continue
		}

		var post model.Post
		if err := json.Unmarshal([]byte(postJson), &post); err == nil {
			posts = append(posts, &post)
		}
	}

	return posts, nil
}

func (r *RedisCacheRepository) InvalidatePostRanking(ctx context.Context) error {

	keys := []string{
		"posts:ranking:day",
		"posts:ranking:week",
		"posts:ranking:month",
	}
	return r.client.Del(ctx, keys...).Err()
}

func (r *RedisCacheRepository) CachePost(ctx context.Context, post *model.Post) error {
	postJson, err := json.Marshal(post)
	if err != nil {
		return err
	}

	return r.client.HSet(ctx, "posts:details",
		fmt.Sprintf("%d", post.ID),
		postJson,
	).Err()
}

func (r *RedisCacheRepository) GetPost(ctx context.Context, postID uint) (*model.Post, error) {
	postJson, err := r.client.HGet(ctx, "posts:details", fmt.Sprintf("%d", postID)).Result()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var post model.Post
	err = json.Unmarshal([]byte(postJson), &post)
	return &post, err
}

func getExpiration(timeRange string) time.Duration {
	switch timeRange {
	case "day":
		return 24 * time.Hour
	case "week":
		return 7 * 24 * time.Hour
	case "month":
		return 30 * 24 * time.Hour
	default:
		return 24 * time.Hour
	}

}

func (r *RedisCacheRepository) InvalidateToken(ctx context.Context, token string, expiration time.Duration) error {
	err := r.client.Set(ctx, "invalid_tokens:"+token, "1", expiration).Err()
	log.Print(err)
	return err
}

func (r *RedisCacheRepository) IsTokenInvalid(ctx context.Context, token string) (bool, error) {
	exists, err := r.client.Exists(ctx, "invalid_tokens:"+token).Result()
	return exists > 0, err
}
