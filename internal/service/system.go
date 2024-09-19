package service

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/0x0FACED/pdf-proto/pdf_service/gen"
	"github.com/0x0FACED/pdf-saver-api/internal/domain/models"
	"github.com/go-rod/rod/lib/proto"
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

	pdfData, err := s.visitPage(req.OriginalUrl)
	if err != nil {
		return nil, err
	}

	pdf = &models.PDF{
		Description: req.Description,
		Filename:    fmt.Sprintf("%d.pdf", time.Now().Unix()),
		Content:     pdfData,
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

func (s *PDFService) checkDailyLimit(userID int64) error {
	return s.mem.CheckDailyLimit(context.TODO(), userID)
}

func (s *PDFService) visitPage(url string) ([]byte, error) {
	page := s.rod.MustConnect().MustPage(url)

	page.MustWaitLoad().MustElement("body")

	// js для закрытия некоторых окон (например, cookie)
	_, err := page.Eval(`() => {
		let cookieBanner = document.querySelector('div[class*="cookie"]');
		if (cookieBanner) {
			cookieBanner.remove();
		}
	
		let modal = document.querySelector('.modal');
		if (modal) {
			modal.remove();
		}
	
		// Примеры для других возможных подсказок
		let popup = document.querySelector('.popup');
		if (popup) {
			popup.remove();
		}
	
		let notification = document.querySelector('.notification');
		if (notification) {
			notification.remove();
		}
	}`)
	if err != nil {
		log.Printf("Error closing pop up windows: %v", err)
	}

	var h float64
	var w float64
	var scale float64
	var mt float64
	var mb float64
	var ml float64
	var mr float64

	h = 11.69
	w = 8.27
	scale = 1
	mt = 0
	mb = 0
	ml = 0
	mr = 0

	// Устанавливаем опции для печати
	pdfStream, err := page.PDF(&proto.PagePrintToPDF{
		PrintBackground: true,
		Scale:           &scale,
		PaperWidth:      &w,
		PaperHeight:     &h,
		MarginTop:       &mt,
		MarginBottom:    &mb,
		MarginLeft:      &ml,
		MarginRight:     &mr,
	})
	if err != nil {
		return nil, err
	}

	// Чтение данных из StreamReader в []byte
	pdfBuffer, err := io.ReadAll(pdfStream)
	if err != nil {
		log.Fatal(err)
	}

	return pdfBuffer, nil

}
