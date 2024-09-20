package service

import (
	"github.com/0x0FACED/pdf-proto/pdf_service/gen"
	"github.com/0x0FACED/pdf-saver-api/config"
	"github.com/0x0FACED/pdf-saver-api/internal/logger"
	"github.com/0x0FACED/pdf-saver-api/internal/mem"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

type PDFService struct {
	gen.UnimplementedPDFServiceServer

	logger *logger.ZapLogger
	cfg    config.PRFServiceConfig
	// Еще реализовать механизм очереди через каналы (брокер не нужен)

	rod *rod.Browser
	mem mem.MemoryCacher
}

func New(logger *logger.ZapLogger, cfg config.PRFServiceConfig, mem mem.MemoryCacher) *PDFService {
	path, _ := launcher.LookPath()
	u := launcher.New().Bin(path).MustLaunch()
	r := rod.New()

	r.MustConnect().ControlURL(u)

	return &PDFService{
		logger: logger,
		cfg:    cfg,
		rod:    r,
		mem:    mem,
	}
}
