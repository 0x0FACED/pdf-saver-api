package service

import (
	"context"
	"io"
	"log"

	"github.com/go-rod/rod/lib/proto"
)

const (
	JSClosePopUp = `() => {
		let cookieBanner = document.querySelector('div[class*="cookie"]');
		if (cookieBanner) {
			cookieBanner.remove();
		}
	
		let modal = document.querySelector('.modal');
		if (modal) {
			modal.remove();
		}

		let popup = document.querySelector('.popup');
		if (popup) {
			popup.remove();
		}
	
		let notification = document.querySelector('.notification');
		if (notification) {
			notification.remove();
		}
	}`
)

func (s *PDFService) checkDailyLimit(userID int64) error {
	return s.mem.CheckDailyLimit(context.TODO(), userID)
}

func (s *PDFService) visitPage(url string, scale float64) ([]byte, error) {
	page := s.rod.MustPage(url)

	page.MustWaitLoad().MustElement("body")

	// js для закрытия некоторых окон (например, cookie)
	_, err := page.Eval(JSClosePopUp)
	if err != nil {
		log.Printf("Error closing pop up windows: %v", err)
	}

	// TODO: перенести куда-то настройки стандартные
	var h float64
	var w float64
	var mt float64
	var mb float64
	var ml float64
	var mr float64

	h = 11.69
	w = 8.27
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
