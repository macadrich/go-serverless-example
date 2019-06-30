package main

import (
	"log"
	"magdb/server-example/model"
	rest "magdb/server-example/rest"
	"net/http"
	"os"

	"gitlab.com/azitra/magdb"

	"github.com/aws/aws-lambda-go/lambda"
)

// UserRepository -
type UserRepository interface {
	GetUserByID(id string) (*model.User, error)
	ListUsers() (*[]model.User, error)
}

// UserService -
type UserService struct {
	userRepo UserRepository
}

// GetUserByID a single user by id
func (h *UserService) GetUserByID(id string, request rest.Request) (rest.Response, error) {
	user, err := h.userRepo.GetUserByID(id)
	if err != nil {
		return rest.APIErrResponse(err, http.StatusNotFound)
	}

	return rest.APIResponse(map[string]interface{}{
		"user": user,
	}, http.StatusOK)
}

// ListUsers list users
func (h *UserService) ListUsers(request rest.Request) (rest.Response, error) {
	users, err := h.userRepo.ListUsers()
	if err != nil {
		return rest.APIErrResponse(err, http.StatusNotFound)
	}

	return rest.APIResponse(map[string]interface{}{
		"users": users,
	}, http.StatusOK)
}

func main() {
	connection := magdb.NewMagDB(os.Getenv("REGION"), os.Getenv("USER_TABLE"))
	db, err := connection.InitDynamoDBConnection()
	if err != nil {
		log.Panic(err)
	}
	repository := model.NewUserRepository(db)
	handler := &UserService{
		userRepo: repository,
	}

	router := rest.UserRouter(handler)
	lambda.Start(router)
}
