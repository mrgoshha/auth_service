package model

type ResponseError struct {
	// Описание ошибки
	Message string `json:"message"`
	// Идентификатор запроса. Предназначен для более быстрого поиска проблем.
	RequestId string `json:"request_id,omitempty"`
}
