package web

type errorDetail string

const (
	UnauthorizedError   errorDetail = "authorisation error"
	InternalServerError errorDetail = "server error"
	ForbiddenError      errorDetail = "forbidden"
)

type ErrorResponse struct {
	Message errorDetail
}
