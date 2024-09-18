package redis

import (
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

func (r *RedisClient) SavePDF(pdf *models.PDF) error { panic("impl me") }

func (r *RedisClient) GetPDF(desc string) (*models.PDF, error) { panic("impl me") }
