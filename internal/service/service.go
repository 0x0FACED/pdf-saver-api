package service

import (
	"github.com/0x0FACED/pdf-proto/pdf_service/gen"
	"github.com/0x0FACED/pdf-saver-api/config"
	"github.com/0x0FACED/pdf-saver-api/internal/logger"
)

type PDFService struct {
	gen.UnimplementedPDFServiceServer

	logger *logger.ZapLogger
	cfg    config.PRFServiceConfig
}

func New(logger *logger.ZapLogger, cfg config.PRFServiceConfig) *PDFService {
	return &PDFService{
		logger: logger,
		cfg:    cfg,
	}
}
