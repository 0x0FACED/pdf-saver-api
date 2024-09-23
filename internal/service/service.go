package service

import (
	"github.com/0x0FACED/pdf-proto/pdf_service/gen"
	"github.com/0x0FACED/pdf-saver-api/config"
	"github.com/0x0FACED/pdf-saver-api/internal/logger"
	"github.com/0x0FACED/pdf-saver-api/internal/mem"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"go.uber.org/zap"
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
	logger.Info("Creating PDFService...")
	browserPath, ex := launcher.LookPath()
	if !ex {
		logger.Info("Cant find Chrome, exiting...")
		// tmp
		panic("exiting")
	}
	u := launcher.New().Bin(browserPath).MustLaunch()

	r := rod.New().ControlURL(u).MustConnect()

	logger.Info("Connected to Rod", zap.String("browserPath", browserPath), zap.String("url", u))

	return &PDFService{
		logger: logger,
		cfg:    cfg,
		rod:    r,
		mem:    mem,
	}
}
