package shortcode

import (
	"strings"
)

// Функция для перевода в 62 СС
func ToBase62(n int64) string {
	if n == 0 {
		return "0"
	}
	//Задаем символы для использования (всего 62 как и в 62 СС)
	charset := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	//strings.Builder работает как буфер — добавляет байты в конец без создания промежуточных строк.
	var sb strings.Builder
	sb.Grow(10)
	//Сбор строки в 62 СС для последующего разворота
	for n > 0 {
		remainder := n % 62
		sb.WriteByte(charset[remainder])
		n /= 62
	}
	//Возвращает готовую строку
	s := sb.String()
	//Создаем сегмент рун
	runes := []rune(s)
	//переворачиваем все символы
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	//возвращаем получившиюся строку
	return string(runes)
}
