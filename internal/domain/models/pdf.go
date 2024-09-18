package models

type PDF struct {
	// Ориг ссылка на страницу, с которой делаем pdf
	OriginalURL string
	// Название pdf (мб будет пустое, тогда сервак даст название случайное из UUID)
	Description string
	// Содержание pdf файла
	Content []byte
}
