package web

type errorDetail string

const (
	InternalServerErrorString       errorDetail = "server error"
	TokenExpectedErrorString        errorDetail = "token expected"
	TokenInvalidErrorString         errorDetail = "invalid token"
	InvalidSubjectErrorString       errorDetail = "invalid subject"
	InvalidBasicAuthFormErrorString errorDetail = "invalid basic auth form"
	InvalidLoginPasswordErrorString errorDetail = "invalid login or password"
)

type ErrorResponse struct {
	Message interface{} `json:"message"`
}

type NotFoundError struct {
	Message string
}

func (NotFoundError) Error() string {
	return "not found error"
}

type AlreadyExistError struct {
	Message string
}

func (AlreadyExistError) Error() string {
	return "already exist error"
}

type ValidationError struct {
	Message string
}

func (ValidationError) Error() string {
	return "validation error"
}

type InternalServerError struct {
}

func (InternalServerError) Error() string {
	return "internal server error"
}

type BadRequestError struct {
	Message string
}

func (BadRequestError) Error() string {
	return "bad request error"
}

type UnAuthorizedError struct {
}

func (UnAuthorizedError) Error() string {
	return "unauthorized error"
}

type ForbiddenError struct {
}

func (ForbiddenError) Error() string {
	return "forbidden error"
}
