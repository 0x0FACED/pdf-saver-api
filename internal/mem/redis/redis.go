package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/0x0FACED/pdf-saver-api/config"
	"github.com/0x0FACED/pdf-saver-api/internal/domain/models"
	"github.com/0x0FACED/pdf-saver-api/internal/mem"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	r *redis.Client

	cfg config.MemCacheConfig
}

func New(cfg config.MemCacheConfig) mem.MemoryCacher {
	client := redis.NewClient(
		&redis.Options{
			Addr: cfg.Host + ":" + cfg.Port,
		},
	)

	return &RedisClient{
		r:   client,
		cfg: cfg,
	}
}

func (r *RedisClient) SavePDF(ctx context.Context, userID int64, pdf *models.PDF) error {
	pdfKey := fmt.Sprintf("user:%d:pdf:%s", userID, pdf.Description)

	if err := r.r.Set(ctx, pdfKey, pdf, 0).Err(); err != nil {
		return err
	}

	listKey := fmt.Sprintf("user:%d:pdfs", userID)
	if err := r.r.LPush(ctx, listKey, pdf.Description).Err(); err != nil {
		return err
	}

	if err := r.r.LTrim(ctx, listKey, 0, 4).Err(); err != nil {
		return err
	}

	return nil
}

func (r *RedisClient) GetPDF(ctx context.Context, userID int64, description string) (*models.PDF, error) {
	pdfKey := fmt.Sprintf("user:%d:pdf:%s", userID, description)

	val, err := r.r.Get(ctx, pdfKey).Bytes()
	if err == redis.Nil {
		return nil, fmt.Errorf("PDF not found")
	} else if err != nil {
		return nil, err
	}

	var pdf models.PDF
	if err := json.Unmarshal(val, &pdf); err != nil {
		return nil, err
	}

	return &pdf, nil
}

func (r *RedisClient) DeletePDF(ctx context.Context, userID int64, description string) error {
	pdfKey := fmt.Sprintf("user:%d:pdf:%s", userID, description)
	listKey := fmt.Sprintf("user:%d:pdfs", userID)

	if err := r.r.Del(ctx, pdfKey).Err(); err != nil {
		return err
	}
	if err := r.r.LRem(ctx, listKey, 1, description).Err(); err != nil {
		return err
	}

	return nil
}

func (r *RedisClient) DeleteAllPDF(ctx context.Context, userID int64) error {
	listKey := fmt.Sprintf("user:%d:pdfs", userID)

	descriptions, err := r.r.LRange(context.Background(), listKey, 0, -1).Result()
	if err != nil {
		return err
	}

	for _, description := range descriptions {
		pdfKey := fmt.Sprintf("user:%d:pdf:%s", userID, description)
		if err := r.r.Del(context.Background(), pdfKey).Err(); err != nil {
			return err
		}
	}

	if err := r.r.Del(context.Background(), listKey).Err(); err != nil {
		return err
	}

	return nil
}

func (r *RedisClient) CheckDailyLimit(ctx context.Context, userID int64) error {
	limitKey := fmt.Sprintf("user:%d:pdf:daily_limit", userID)
	count, _ := r.r.Get(ctx, limitKey).Int()

	if count >= 15 {
		return fmt.Errorf("daily limit exceeded")
	}

	// Если первый запрос, установим 24-часовой лимит
	if count == 0 {
		r.r.Set(ctx, limitKey, 1, 24*time.Hour)
	} else {
		r.r.Incr(ctx, limitKey)
	}

	return nil
}
