package user

import "context"

type User struct {
	ID int
	Name string
	Email string
}

type Repository interface {
	CreateTable(ctx context.Context) error
	Create(ctx context.Context , name , email string)(int, error)
	GetByID(ctx context.Context,id int) (User,error)	
}

// it's a new package , user organized by feature.
// User struct instead of string. Before, getUserById returned only a string
// The interface describes what the app needs, not how sql does it.
//ctx context.Context is the first parameter of every method. That's Go convention.It lets a request's timeout
// or cancellation reach all the way to db query

// Before

/*
createUsersTableMySQL / CreateUsersTablePostgres
createUserMySQL / CreateUserPostgres
getUserByIDMySQL / getuserByIDPostgres
*/

//Now
// CreateTable
// Create
// GetByID