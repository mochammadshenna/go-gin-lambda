package httphelper

import (
	"ai-service/internal/model/api"
	"ai-service/internal/state"
	"ai-service/internal/util/exceptioncode"
	"ai-service/internal/util/helper"
	"ai-service/internal/util/json"
	"ai-service/internal/util/logger"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/schema"
)

var Decoder = schema.NewDecoder()

func init() {
	Decoder.RegisterConverter([]string{}, convertStringCommaSeparated)
}

func convertStringCommaSeparated(value string) reflect.Value {
	// if empty value should return []string{} instead of []string{""}
	if value == "" {
		return reflect.ValueOf([]string{})
	}
	return reflect.ValueOf(strings.Split(value, ","))
}

func ReadQueryParam(request *http.Request, result interface{}) error {
	err := Decoder.Decode(result, request.URL.Query())
	if err != nil {
		return err
	}
	return nil
}

func Close(ctx context.Context, body io.Closer) {
	if body != nil {
		err := body.Close()
		if err != nil {
			logger.Errorf(ctx, "got an error while closing the body; err=%+v", err)
			return
		}
	}
}

func Read(request *http.Request, result interface{}) error {
	defer Close(request.Context(), request.Body)

	if request.Method == http.MethodPost || request.Method == http.MethodPut || request.Method == http.MethodPatch {
		var b bytes.Buffer
		_, err := io.Copy(&b, request.Body)
		if err != nil {
			logger.Error(request.Context(), strings.Replace(fmt.Sprintf("request body: %+v", b.String()), "\u0026", "", 1), err)
			return api.ErrorResponse{
				Message: exceptioncode.CodeInvalidRequest,
			}
		}

		jsonDecoder := json.NewDecoder(&b)
		err = jsonDecoder.Decode(&result)
		if err != nil && err != io.EOF {
			logger.Error(request.Context(), strings.Replace(fmt.Sprintf("request body: %+v", b.String()), "\u0026", "", 1), err)
			return api.ErrorResponse{
				CodeMessage: exceptioncode.CodeInvalidRequest,
				Message:     err.Error(),
			}
		}

		return nil
	}

	err := Decoder.Decode(result, request.URL.Query())
	if err != nil {
		logger.Error(request.Context(), strings.Replace(fmt.Sprintf("request params: %s", request.URL.RawQuery), "\u0026", ", ", -1), err)
		return parseError(err)
	}

	return nil
}

func Write(ctx context.Context, writer http.ResponseWriter, data interface{}) {
	writer.WriteHeader(http.StatusOK)
	response := api.ApiResponse{
		Header:    getHeader(writer),
		Data:      data,
		Message:   "success",
		ErrorCode: nil,
	}
	write(ctx, writer, response)
}

func WriteSuccessWithStatusCode(ctx context.Context, writer http.ResponseWriter, data interface{}, message string, statusCode int) {
	defaultStatusCode := http.StatusOK // Change this from http.StatusCreated to http.StatusOK
	if statusCode > 200 && statusCode <= 299 {
		defaultStatusCode = statusCode
	}
	writer.WriteHeader(defaultStatusCode)
	response := api.ApiResponse{
		Header:    getHeader(writer), // Add this line to include the header
		Data:      data,
		Message:   message,
		ErrorCode: nil,
	}
	write(ctx, writer, response)
}

func WriteError(ctx context.Context, writer http.ResponseWriter, errorResponse error) {
	httpError, ok := errorResponse.(api.HttpError)
	if ok {
		WriteErrorWithStatusCode(ctx, writer, httpError.Error(), httpError.StatusCode())
		return
	}
	WriteErrorWithStatusCode(ctx, writer, errorResponse.Error(), http.StatusBadRequest)
}

func WriteErrorWithStatusCode(ctx context.Context, writer http.ResponseWriter, errMessage string, statusCode int) {
	defaultStatusCode := http.StatusBadRequest
	if statusCode >= 400 && statusCode <= 599 {
		defaultStatusCode = statusCode
	}
	writer.WriteHeader(defaultStatusCode)
	response := api.ApiResponse{
		Header:    getHeader(writer),
		Message:   errMessage,
		ErrorCode: defaultStatusCode,
	}
	write(ctx, writer, response)
}

func WriteErrorWithData(ctx context.Context, writer http.ResponseWriter, errorResponse error, data interface{}) {
	httpError, ok := errorResponse.(api.HttpError)
	if ok {
		WriteErrorWithStatusCode(ctx, writer, httpError.Error(), httpError.StatusCode())
		return
	}
	writer.WriteHeader(http.StatusBadRequest)
	response := api.ApiResponse{
		Message: httpError.Error(),
	}
	write(ctx, writer, response)
}

func write(_ context.Context, writer http.ResponseWriter, response api.ApiResponse) {
	data, err := response.ToJSON()
	helper.PanicOnError(err)
	writer.Write(data)
}

func parseError(err error) error {
	errors := []api.ErrorValidate{}
	new := err.(schema.MultiError)
	for i, a := range new {
		errors = append(errors, api.ErrorValidate{
			Key:         i,
			CodeMessage: "VALIDATION",
			Message:     a.Error(),
		})
	}
	return api.ErrorResponse{
		CodeMessage: exceptioncode.CodeInvalidValidation,
		Message:     "validation error",
		Errors:      errors,
	}
}

func CompareParamWithClaimToken(param string, claimId int64) bool {
	paramId, _ := strconv.ParseInt(param, 10, 64)
	return paramId == claimId
}

func GetHostString(request *http.Request) string {
	return fmt.Sprintf("%s%s", request.URL.Host, request.URL)
}

// GetPlatform will return request header `Platform` or `Platform-Type`.
// Header `Platform` will be prioritized
func GetPlatform(request *http.Request) string {
	platformType := request.Header.Get(state.HttpHeaders().PlatformType)
	platform := request.Header.Get(state.HttpHeaders().Platform)

	if platform == "" {
		platform = platformType
	}

	return platform
}

func getHeader(writer http.ResponseWriter) api.HeaderResponse {
	headerResponse := api.HeaderResponse{
		ServerTimeMs: time.Now().Unix(),
		RequestId:    writer.Header().Get(state.HttpHeaders().RequestId),
	}

	startTimeHeader := writer.Header().Get("X-Start-Time")
	if len(startTimeHeader) > 0 {
		startTime, _ := strconv.ParseInt(startTimeHeader, 10, 64)
		headerResponse.ProcessTimeMs = time.Since(time.Unix(0, startTime)).Milliseconds()
	}

	return headerResponse
}

func WriteErrorAuthorization(ctx context.Context, writer http.ResponseWriter, data interface{}) {
	writer.WriteHeader(http.StatusOK)

	message := exceptioncode.CodeUnauthorized
	if writer.Header().Get("X-Custom-Message") != "" {
		message = writer.Header().Get("X-Custom-Message")
	}
	response := api.ApiResponse{
		Header:    getHeader(writer),
		Data:      data,
		Message:   message,
		ErrorCode: http.StatusUnauthorized,
	}
	write(ctx, writer, response)
}

func WriteErrorClaimToken(ctx context.Context, writer http.ResponseWriter) {
	logger.Error(ctx, api.ErrorResponse{Message: exceptioncode.CodeTokenInvalid, ErrorCode: http.StatusUnauthorized})
	WriteErrorWithStatusCode(ctx, writer, exceptioncode.CodeTokenInvalid, http.StatusUnauthorized)
}
