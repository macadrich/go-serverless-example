package model

import (
	"crypto/rand"
	"magdb/dynamodb"
	"math/big"

	uuid "github.com/satori/go.uuid"
	"gitlab.com/azitra/magdb"
)

var maxRand = big.NewInt(1<<63 - 1)

// UserRepository stores and fetches user
type UserRepository struct {
	datastore magdb.Datastore
}

func randomID() (int64, error) {
	// Get a random number within the range [0, 1<<63-1)
	n, err := rand.Int(rand.Reader, maxRand)
	if err != nil {
		return 0, err
	}
	// Don't assign 0.
	return n.Int64() + 1, nil
}

// NewUserRepository instance
func NewUserRepository(dbs dynamodb.Datastore) *UserRepository {
	return &UserRepository{datastore: dbs}
}

// InsertUser add user in database
func (userRepo *UserRepository) InsertUser(user *User) error {
	id := uuid.NewV4()
	user.ID = id.String()
	return userRepo.datastore.Add(user)
}

// GetUserByID get user by id
func (userRepo *UserRepository) GetUserByID(id string) (*User, error) {
	var user *User
	searchByID := dynamodb.ID(id)
	attributes := dynamodb.FieldAttributes("id")
	if err := userRepo.datastore.Get(searchByID, attributes, &user); err != nil {
		return nil, err
	}
	return user, nil
}

// ListUsers list view all users in databases
func (userRepo *UserRepository) ListUsers() (*[]User, error) {
	var users *[]User
	if err := userRepo.datastore.List(&users); err != nil {
		return nil, err
	}
	return users, nil
}

// GetUserByEmail get user by email
func (userRepo *UserRepository) GetUserByEmail(email string) (*User, error) {
	var user *User
	searchByEmail := dynamodb.Email(email)
	attributes := dynamodb.FieldAttributes("email-index", "email")
	if err := userRepo.datastore.Get(searchByEmail, attributes, &user); err != nil {
		return nil, err
	}
	return user, nil
}
