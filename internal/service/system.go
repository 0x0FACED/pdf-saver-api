package service

import (
	"context"
	"fmt"
	"time"

	"github.com/0x0FACED/pdf-proto/pdf_service/gen"
	"github.com/0x0FACED/pdf-saver-api/internal/domain/models"
)

func (s *PDFService) ConvertToPDF(ctx context.Context, req *gen.ConvertToPDFRequest) (*gen.ConvertToPDFResponse, error) {
	if err := s.checkDailyLimit(req.UserId); err != nil {
		return nil, err
	}

	pdf, err := s.mem.GetPDF(ctx, req.UserId, req.Description)
	if err == nil {
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
	if err := s.mem.SavePDF(ctx, req.UserId, pdf); err != nil {
		return nil, err
	}

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
