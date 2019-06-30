package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

// ResponseError -
type ResponseError map[string]error

// Request alias for api gateway request
type Request events.APIGatewayProxyRequest

// Response alias for api gateway response
type Response events.APIGatewayProxyResponse

// APIResponse wrapped around the api gateway proxy response
func APIResponse(data interface{}, code int) (Response, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return APIErrResponse(fmt.Errorf("%s", err.Error()), http.StatusNotAcceptable)
	}
	return Response{
		Headers: map[string]string{
			"Access-Control-Allow-Origin": "*",
			"Content-Type":                "application/json",
		},
		Body:       string(body),
		StatusCode: code,
	}, nil
}

// APIErrResponse returns an error in a specified format
func APIErrResponse(err error, code int) (Response, error) {
	data := map[string]string{
		"error": err.Error(),
	}
	body, err := json.Marshal(data)
	if err != nil {
		return APIErrResponse(fmt.Errorf("%s", err.Error()), http.StatusNotAcceptable)
	}
	return Response{
		Headers: map[string]string{
			"Access-Control-Allow-Origin": "*",
			"Content-Type":                "application/json",
		},
		Body:       string(body),
		StatusCode: code,
	}, err
}

// ParseRequestBody parses the json to a given struct pointer
func ParseRequestBody(request Request, obj interface{}) error {
	return json.Unmarshal([]byte(request.Body), &obj)
}

// RequestHandleFunc is an alias for an api gateway request signature
type RequestHandleFunc func(request Request) (Response, error)

// RegisterHandler represents a register user Lambda handler
type RegisterHandler interface {
	InsertUser(request Request) (Response, error)
}

// UserHandler represents a user Lambda handler
type UserHandler interface {
	GetUserByID(id string, request Request) (Response, error)
	ListUsers(request Request) (Response, error)
}

// FetchUserHandler represents a get user by email Lambda handler
type FetchUserHandler interface {
	GetUserByEmail(email string, request Request) (Response, error)
}

// RegisterRouter routes restful endpoints to the correct method
func RegisterRouter(h RegisterHandler) RequestHandleFunc {
	return func(request Request) (Response, error) {
		switch request.HTTPMethod {
		case "POST":
			return h.InsertUser(request)
		default:
			return APIErrResponse(errors.New("method not allowed"), http.StatusMethodNotAllowed)
		}
	}
}

// UserRouter routes restful endpoints to the correct method
func UserRouter(h UserHandler) RequestHandleFunc {
	return func(request Request) (Response, error) {
		switch request.HTTPMethod {
		case "GET":
			id := request.PathParameters["id"]
			if id != "" {
				return h.GetUserByID(id, request)
			}
			return h.ListUsers(request)
		default:
			return APIErrResponse(errors.New("method not allowed"), http.StatusMethodNotAllowed)
		}
	}
}

// FetchUserRouter routes user info endpoints
func FetchUserRouter(h FetchUserHandler) RequestHandleFunc {
	return func(request Request) (Response, error) {
		switch request.HTTPMethod {
		case "GET":
			email := request.PathParameters["email"]
			if email != "" {
				return h.GetUserByEmail(email, request)
			}
			return APIErrResponse(errors.New("method not allowed"), http.StatusMethodNotAllowed)
		default:
			return APIErrResponse(errors.New("method not allowed"), http.StatusMethodNotAllowed)
		}
	}
}
