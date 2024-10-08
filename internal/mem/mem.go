package mem

import (
	"context"

	"github.com/0x0FACED/pdf-saver-api/internal/domain/models"
)

// В целом интерфейс для реализации через redis или memcached
// Через интерфейс, чтобы в теории можно было переключиться с redis на memcached
// Ну и соответствовать принципам SOLID отчасти)
type MemoryCacher interface {
	// сейвим pdf в редис, чтобы один юзер не мог постоянно дергать rod для получения pdf из ссылки
	// + добавим cooldown на такие запросы)))
	SavePDF(ctx context.Context, userID int64, pdf *models.PDF) error

	// Получить pdf по описанию (названию)
	GetPDF(ctx context.Context, userID int64, description string) (*models.PDF, error)

	DeletePDF(ctx context.Context, userID int64, description string) error

	DeleteAllPDF(ctx context.Context, userID int64) error

	DailyLimitChecker
}

type DailyLimitChecker interface {
	CheckDailyLimit(ctx context.Context, userID int64) error
}
