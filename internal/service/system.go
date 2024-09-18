package service

import (
	"context"
	"log"

	"github.com/0x0FACED/pdf-proto/pdf_service/gen"
)

func (s *PDFService) ConvertToPDF(context.Context, *gen.ConvertToPDFRequest) (*gen.ConvertToPDFResponse, error) {
	// Сначала проверяем наличие в памяти (memcache | redis)
	// Если там пусто, посещаем и возвращаем
	// Логика посещения html страницы и конвертации ее в PDF
	// И отправки ее юзеру через тг
}

func (s *PDFService) GetSavedPDF(context.Context, *gen.GetSavedPDFRequest) (*gen.GetSavedPDFResponse, error) {
	// Будет редис (или memcache) для хранения последних сохраненных юзером ссылок
	// в редисе будет лежать сжатый контент
	// по запросу проверяем редис, если есть - возвращаем
	// Иначе же вернем пустой ответ + ошибку, мол контент не найден
}

// Посещаем страницу и собираем с нее контент в виде []byte
func (s *PDFService) visitPage(url string) ([]byte, error) {
	page := s.rod.MustPage(url)

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

}

// Здесь будем сод
func (s *PDFService) bytesToPDF(bytePDF []byte) error {

}
