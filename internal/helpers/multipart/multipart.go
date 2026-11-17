package multipart

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"mime/multipart"
)

// Form собранное multipart/form-data тело запроса вместе с заголовком Content-Type
type Form struct {
	Body        string
	ContentType string
}

// Builder собирает multipart/form-data тело запроса
type Builder struct {
	buf    bytes.Buffer
	writer *multipart.Writer
}

// NewBuilder создаёт пустой билдер multipart-формы
func NewBuilder() *Builder {
	b := &Builder{}
	b.writer = multipart.NewWriter(&b.buf)
	return b
}

// Field добавляет текстовое поле формы
func (b *Builder) Field(name, value string) *Builder {
	if err := b.writer.WriteField(name, value); err != nil {
		panic(fmt.Sprintf("Не удалось добавить поле %s в multipart-форму: %s", name, err.Error()))
	}
	return b
}

// File добавляет файл в форму под указанным именем поля
func (b *Builder) File(fieldName, fileName string, content []byte) *Builder {
	part, err := b.writer.CreateFormFile(fieldName, fileName)
	if err != nil {
		panic(fmt.Sprintf("Не удалось добавить файл %s в multipart-форму: %s", fileName, err.Error()))
	}
	if _, err = part.Write(content); err != nil {
		panic(fmt.Sprintf("Не удалось записать файл %s в multipart-форму: %s", fileName, err.Error()))
	}
	return b
}

// Build завершает форму и возвращает тело запроса с Content-Type, содержащим boundary
func (b *Builder) Build() Form {
	if err := b.writer.Close(); err != nil {
		panic(fmt.Sprintf("Не удалось завершить multipart-форму: %s", err.Error()))
	}
	return Form{Body: b.buf.String(), ContentType: b.writer.FormDataContentType()}
}

// GeneratePNG генерирует валидное изображение PNG минимального размера для загрузки в качестве фотографии
func GeneratePNG() []byte {
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	img.Set(0, 0, color.RGBA{R: 255, A: 255})

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		panic(fmt.Sprintf("Не удалось сгенерировать PNG: %s", err.Error()))
	}
	return buf.Bytes()
}
