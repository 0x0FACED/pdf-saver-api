package service

import (
	"context"

	"github.com/0x0FACED/pdf-proto/pdf_service/gen"
)

func (s *PDFService) ConvertToPDF(context.Context, *gen.ConvertToPDFRequest) (*gen.ConvertToPDFResponse, error) {
	// Логика посещения html страницы и конвертации ее в PDF
	// И отправки ее юзеру через тг
}

func (s *PDFService) GetSavedPDF(context.Context, *gen.GetSavedPDFRequest) (*gen.GetSavedPDFResponse, error) {
	// Будет редис (или memcache) для хранения последних сохраненных юзером ссылок
	// в редисе будет лежать сжатый контент
	// по запросу проверяем редис, если есть - возвращаем
	// Иначе же вернем пустой ответ + ошибку, мол контент не найден
}
