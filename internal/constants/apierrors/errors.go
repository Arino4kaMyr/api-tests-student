package apierrors

// Коды ошибок (поле error) и сообщения (поле message) из ответов API
const (
	CodeUnauthorized  = "unauthorized"
	CodeInvalidReq    = "invalid_request"
	CodeNotFound      = "not_found"
	CodeForbidden     = "forbidden"
	CodeAlreadyExists = "already_exists"

	MsgAuthHeaderRequired = "Authorization header required"
	MsgInvalidToken       = "Invalid or expired token" //nolint:gosec // текст ошибки API, не секрет
)

// InvalidToken заведомо некорректный токен авторизации
const InvalidToken = "incorrect_token"
