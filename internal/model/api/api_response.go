package api

import (
	"ai-service/internal/util/exceptioncode"
	"ai-service/internal/util/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/pkg/errors"
)

type ApiResponse struct {
	Header    HeaderResponse `json:"header"`
	Data      interface{}    `json:"data"`
	Message   interface{}    `json:"message"`
	ErrorCode interface{}    `json:"errorCode"`
}

type ErrorResponse struct {
	HttpCode     int           `json:"-"`
	CodeMessage  string        `json:"codeMessage"`
	Message      interface{}   `json:"message"`
	Errors       interface{}   `json:"errors"`
	ErrorCode    int           `json:"errorCode"`
	ErrorMessage []interface{} `json:"errorMessage"`
}

func (e ErrorResponse) Error() string {
	return e.Message.(string)
}

func (e ErrorResponse) StatusCode() int {
	return e.HttpCode
}

type ErrorValidate struct {
	Key         string `json:"key"`
	CodeMessage string `json:"code"`
	Message     string `json:"message"`
}

type HttpError interface {
	Error() string
	StatusCode() int
}

// APIResponse ...
func GenerateResponse(statusCode int, body interface{}) (events.APIGatewayProxyResponse, error) {
	response := ApiResponse{
		Data:      body,
		Message:   "success",
		ErrorCode: nil,
	}
	responseBody, _ := json.Marshal(response)

	return events.APIGatewayProxyResponse{
		Body:       string(responseBody),
		StatusCode: statusCode,
	}, nil
}

// APIErrResponse ...
func ErrResponse(statusCode int, err error) (events.APIGatewayProxyResponse, error) {
	defaultMessage := exceptioncode.CodeBadRequest
	if statusCode == 401 {
		defaultMessage = exceptioncode.CodeUnauthorized
	} else if statusCode == 404 {
		defaultMessage = exceptioncode.CodeDataNotFound
	} else if statusCode == 405 {
		defaultMessage = exceptioncode.CodeMethodNotAllowed
	}

	errorMessage := ApiResponse{
		Message:   defaultMessage,
		ErrorCode: statusCode,
	}
	messageBody, _ := json.Marshal(errorMessage)

	return events.APIGatewayProxyResponse{
		Body:       string(messageBody),
		StatusCode: statusCode,
	}, nil
}

// APIServerError ...
func ServerError(statusCode int, err error) (events.APIGatewayProxyResponse, error) {

	printStackTrace(err)

	return events.APIGatewayProxyResponse{
		Body:       exceptioncode.CodeInternalServerError,
		StatusCode: http.StatusInternalServerError,
	}, err
}

type stacktracer interface {
	StackTrace() errors.StackTrace
}

type causer interface {
	Cause() error
}

func printStackTrace(err error) {

	var errStack errors.StackTrace

	for err != nil {
		// Find the earliest error.StackTrace
		if t, ok := err.(stacktracer); ok {
			errStack = t.StackTrace()
		}
		if c, ok := err.(causer); ok {
			err = c.Cause()
		} else {
			break
		}
	}
	if errStack != nil {
		fmt.Printf("%+v\n", errStack)
	} else {
		fmt.Printf("%+v\n", errors.WithStack(err))
	}
}

type HeaderResponse struct {
	ServerTimeMs  int64  `json:"serverTimeMs"`
	ProcessTimeMs int64  `json:"processTimeMs"`
	RequestId     string `json:"requestId"`
}

func (a *ApiResponse) ToJSON() ([]byte, error) {
	return json.Marshal(a)
}

func ConditionalResponse[T any](condition bool, a, b T) T {
	if condition {
		return a
	}
	return b
}

type Pagination struct {
	Page      int64 `json:"page,omitempty"`
	PageSize  int64 `json:"page_size,omitempty"`
	TotalPage int64 `json:"total_page,omitempty"`
	TotalData int64 `json:"total_data,omitempty"`
}

func newPagination(totalData, page, pageSize int64) Pagination {
	var totalPage int64
	if pageSize > 0 {
		totalPage = totalData / pageSize
	}

	if totalData%pageSize > 0 {
		totalPage++
	}

	return Pagination{
		Page:      page,
		PageSize:  pageSize,
		TotalPage: totalPage,
		TotalData: totalData,
	}
}

type HealthCheckResponse struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}
