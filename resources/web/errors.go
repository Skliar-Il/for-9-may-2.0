package web

type errorDetail string

const (
	UnauthorizedError    errorDetail = "authorisation error"
	InternalServerError  errorDetail = "server error"
	TokenExpectedError   errorDetail = "token expected"
	TokenInvalidError    errorDetail = "invalid token"
	InvalidSubjectError  errorDetail = "invalid subject"
	InvalidBasicAuthForm errorDetail = "invalid basic auth form"
	InvalidLoginPassword errorDetail = "invalid login or password"
)

type ErrorResponse struct {
	Message errorDetail
}
