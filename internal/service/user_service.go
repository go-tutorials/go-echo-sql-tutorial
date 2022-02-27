package service

import (
	"context"
	"database/sql"
	"fmt"
	q "github.com/core-go/sql"
	"reflect"
	"strings"

	. "go-service/internal/model"
)

type UserService interface {
	All(ctx context.Context) (*[]User, error)
	Load(ctx context.Context, id string) (*User, error)
	Insert(ctx context.Context, user *User) (int64, error)
	Update(ctx context.Context, user *User) (int64, error)
	Patch(ctx context.Context, user map[string]interface{}) (int64, error)
	Delete(ctx context.Context, id string) (int64, error)
}

type userService struct {
	DB *sql.DB
}

func NewUserService(db *sql.DB) UserService {
	return &userService{DB: db}
}

func (s *userService) All(ctx context.Context) (*[]User, error) {
	query := "select id, username, email, phone, date_of_birth from users"
	rows, err := s.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	var users []User
	for rows.Next() {
		var user User
		err = rows.Scan(&user.Id, &user.Username, &user.Phone, &user.Email, &user.DateOfBirth)
		users = append(users, user)
	}
	return &users, nil
}

func (s *userService) Load(ctx context.Context, id string) (*User, error) {
	var user User
	query := "select id, username, email, phone, date_of_birth from users where id = $1"
	err := s.DB.QueryRowContext(ctx, query, id).Scan(&user.Id, &user.Username, &user.Email, &user.Phone, &user.DateOfBirth)
	if err != nil {
		errMsg := err.Error()
		if strings.Compare(fmt.Sprintf(errMsg), "0 row(s) returned") == 0 {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &user, nil
}

func (s *userService) Insert(ctx context.Context, user *User) (int64, error) {
	query := "insert into users (id, username, email, phone, date_of_birth) values ($1, $2, $3, $4, $5)"
	stmt, er0 := s.DB.Prepare(query)
	if er0 != nil {
		return -1, nil
	}
	res, er1 := stmt.ExecContext(ctx, user.Id, user.Username, user.Email, user.Phone, user.DateOfBirth)
	if er1 != nil {
		return -1, nil
	}
	return res.RowsAffected()
}

func (s *userService) Update(ctx context.Context, user *User) (int64, error) {
	query := "update users set username = $1, email = $2, phone = $3, date_of_birth = $4 where id = $5"
	stmt, er0 := s.DB.Prepare(query)
	if er0 != nil {
		return -1, nil
	}
	res, er1 := stmt.ExecContext(ctx, user.Username, user.Email, user.Phone, user.DateOfBirth, user.Id)
	if er1 != nil {
		return -1, er1
	}
	return res.RowsAffected()
}

func (s *userService) Patch(ctx context.Context, user map[string]interface{}) (int64, error) {
	userType := reflect.TypeOf(User{})
	jsonColumnMap := q.MakeJsonColumnMap(userType)
	colMap := q.JSONToColumns(user, jsonColumnMap)
	keys, _ := q.FindPrimaryKeys(userType)
	query, args := q.BuildToPatch("users", colMap, keys, q.BuildDollarParam)
	res, err := s.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return -1, err
	}
	return res.RowsAffected()
}

func (s *userService) Delete(ctx context.Context, id string) (int64, error) {
	query := "delete from users where id = $1"
	stmt, er0 := s.DB.Prepare(query)
	if er0 != nil {
		return -1, nil
	}
	res, er1 := stmt.ExecContext(ctx, id)
	if er1 != nil {
		return -1, er1
	}
	return res.RowsAffected()
}
