package server

import (
	"log"
	"net"

	"github.com/0x0FACED/pdf-proto/pdf_service/gen"
	"github.com/0x0FACED/pdf-saver-api/config"
	"github.com/0x0FACED/pdf-saver-api/internal/logger"
	"github.com/0x0FACED/pdf-saver-api/internal/mem/redis"
	"github.com/0x0FACED/pdf-saver-api/internal/service"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type server struct {
	service *service.PDFService
	logger  *logger.ZapLogger
}

func New(service *service.PDFService, logger *logger.ZapLogger) *server {
	return &server{
		service: service,
		logger:  logger,
	}
}

const maxMsgSize = 30 * 1024 * 1024 // 30 МБ

func Start() error {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalln("Failed to load cfg: ", err)
		return err
	}
	logger := logger.New()
	logger.Debug("cfg: ", zap.Any("cfg", cfg.PDF))
	lis, err := net.Listen("tcp", cfg.PDF.Host+":"+cfg.PDF.Port)
	if err != nil {
		logger.Fatal("Cannot Dial()", zap.Error(err))
		return err
	}
	mem := redis.New(cfg.MemCache)

	// grpc сервер создаем
	grpcServer := grpc.NewServer(grpc.MaxRecvMsgSize(maxMsgSize), grpc.MaxSendMsgSize(maxMsgSize))

	// создаем объект сервиса
	service := service.New(logger, cfg.PDF, mem)

	// создаем объект НАШЕГО сервера (по сути оболочка для сервиса)
	srv := New(service, logger)

	// регистрируем сервис
	gen.RegisterPDFServiceServer(grpcServer, srv.service)

	// стартует прослушивание по адресу cfg.PDF.Host+":"+cfg.PDF.Post
	return grpcServer.Serve(lis)
}
