package dto

type ContentType int8

const (
	// Unknown Текст или ничего
	Unknown ContentType = iota
	// Photo фотография
	Photo
	// Video Видео
	Video
	// VideoNote Кружок
	VideoNote
	// Voice Голосовое сообщение
	Voice
	// Document документ
	Document
	// Audio mp3
	Audio
)
