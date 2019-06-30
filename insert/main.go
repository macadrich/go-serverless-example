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
	InsertUser(user *model.User) error
}

// UserService -
type UserService struct {
	userRepo UserRepository
}

// InsertUser -
func (us *UserService) InsertUser(request rest.Request) (rest.Response, error) {
	var user *model.User

	if err := rest.ParseRequestBody(request, &user); err != nil {
		return rest.APIErrResponse(err, http.StatusBadRequest)
	}

	if err := us.userRepo.InsertUser(user); err != nil {
		return rest.APIErrResponse(err, http.StatusInternalServerError)
	}

	return rest.APIResponse(map[string]string{
		"id":      user.ID,
		"message": "user created",
		"status":  "success",
	}, http.StatusCreated)
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

	router := rest.RegisterRouter(handler)
	lambda.Start(router)
}
