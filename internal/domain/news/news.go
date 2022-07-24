package news

import (
	"github.com/AscaroLabs/go-news/internal/domain/files"
	"github.com/AscaroLabs/go-news/internal/domain/tags"
)

type News struct {
	// id повости
	Id string `json:"id,omitempty"`
	// Название
	Title string `json:"title,omitempty"`
	// Автор
	Author string `json:"author,omitempty"`
	// Активность (true - активна, false - черновик)
	Active bool `json:"active,omitempty"`
	// Дата начала активности
	ActiveFrom int64 `json:"activeFrom,omitempty"`
	// Текстовое описание
	Text string `json:"text,omitempty"`
	// Текствое описание (для визуального редактора)
	TextJson string `json:"textJson,omitempty"`
	// Идентификатор пользователя
	UserId string `json:"userId,omitempty"`
	// Список тегов
	Tags []*tags.Tag `json:"tags,omitempty"`
	// Список прикрепленных файлов
	FilesInfo []*files.File `json:"filesInfo,omitempty"`
	// Важное
	IsImportant bool `json:"isImportant,omitempty"`
}
