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
	GetUserByEmail(email string) (*model.User, error)
}

// UserService -
type UserService struct {
	userRepo UserRepository
}

// GetUserByEmail a single user by email
func (h *UserService) GetUserByEmail(email string, request rest.Request) (rest.Response, error) {
	user, err := h.userRepo.GetUserByEmail(email)
	if err != nil {
		return rest.APIErrResponse(err, http.StatusOK)
	}
	// send user json format response
	return rest.APIResponse(map[string]interface{}{
		"user": user,
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

	router := rest.FetchUserRouter(handler)
	lambda.Start(router)
}
