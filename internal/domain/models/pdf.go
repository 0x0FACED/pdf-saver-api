package models

type PDF struct {
	// Имя файла
	Filename string `json:"filename"`
	// Название pdf (мб будет пустое, тогда сервак даст название случайное из UUID)
	Description string `json:"description"`
	// Содержание pdf файла
	Content []byte `json:"data"`
}
