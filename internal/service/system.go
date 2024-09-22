package service

import (
	"context"
	"fmt"
	"time"

	"github.com/0x0FACED/pdf-proto/pdf_service/gen"
	"github.com/0x0FACED/pdf-saver-api/internal/domain/models"
	"go.uber.org/zap"
)

func (s *PDFService) ConvertToPDF(ctx context.Context, req *gen.ConvertToPDFRequest) (*gen.ConvertToPDFResponse, error) {
	s.logger.Debug("Received req for save pdf", zap.Any("req", req))
	if err := s.checkDailyLimit(req.UserId); err != nil {
		return nil, err
	}

	pdf, err := s.mem.GetPDF(ctx, req.UserId, req.Description)
	if err == nil {
		s.logger.Debug("PDF Found in Redis", zap.Any("pdf", pdf))
		return &gen.ConvertToPDFResponse{
			PdfData:  pdf.Content,
			Filename: pdf.Filename,
		}, nil
	}

	// Получаем пдф в виде байтов (только контент)
	pdfData, err := s.visitPage(req.OriginalUrl, req.Scale)
	if err != nil {
		return nil, err
	}
	s.logger.Debug("Visited page", zap.String("req", req.OriginalUrl))

	// Сжимаем с помощью gzip наш массив байтов
	// выгодно с точки зрения хранения на сервере в редисе
	compressedData, err := compressPDF(pdfData)
	if err != nil {
		return nil, err
	}

	// В mem мы же сейвим структуру PDF вместе с неймом и описанием
	// Далее может анмаршалить
	pdf = &models.PDF{
		Description: req.Description,
		Filename:    fmt.Sprintf("%d.pdf", time.Now().Unix()),
		Content:     compressedData,
	}

	/*if err := s.mem.SavePDF(ctx, req.UserId, pdf); err != nil {
		s.logger.Debug("Cant save to Redis", zap.Error(err))
		return nil, err
	}*/

	s.logger.Debug("Finished saving", zap.String("pdf", pdf.Filename))
	return &gen.ConvertToPDFResponse{
		PdfData:  pdf.Content,
		Filename: pdf.Filename,
	}, nil
}

func (s *PDFService) GetSavedPDF(ctx context.Context, req *gen.GetSavedPDFRequest) (*gen.GetSavedPDFResponse, error) {
	pdf, err := s.mem.GetPDF(ctx, req.UserId, req.Description)
	if err != nil {
		return nil, fmt.Errorf("PDF not found")
	}

	return &gen.GetSavedPDFResponse{
		PdfData:  pdf.Content,
		Filename: pdf.Filename,
	}, nil
}

func (s *PDFService) DeletePDF(ctx context.Context, req *gen.DeletePDFRequest) (*gen.DeletePDFResponse, error) {
	if err := s.mem.DeletePDF(ctx, req.UserId, req.Description); err != nil {
		return nil, err
	}

	return &gen.DeletePDFResponse{
		Message: "PDF deleted successfully",
	}, nil
}

func (s *PDFService) DeleteAllPDF(ctx context.Context, req *gen.DeleteAllPDFRequest) (*gen.DeleteAllPDFResponse, error) {
	if err := s.mem.DeleteAllPDF(ctx, req.UserId); err != nil {
		return nil, err
	}

	return &gen.DeleteAllPDFResponse{
		Message: "All PDFs deleted successfully",
	}, nil
}
