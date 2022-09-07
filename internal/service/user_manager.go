package service

import (
	"context"
	"errors"
	"fmt"
	"test/internal/broker"
	"test/internal/data"
	"test/internal/dto"
	"test/internal/storage"
)

// UserManager ...
type UserManager struct {
	dataStorage  storage.Storage
	cacheStorage storage.Cache
	broker       broker.MessageBroker
}

func New(dataStorage storage.Storage, cacheStorage storage.Cache, broker broker.MessageBroker) *UserManager {
	return &UserManager{
		dataStorage:  dataStorage,
		cacheStorage: cacheStorage,
		broker:       broker,
	}
}

func (s *UserManager) AddUser(ctx context.Context, user dto.User) error {
	u := data.User{
		ID:       user.ID,
		UserName: user.UserName,
	}

	if err := s.dataStorage.AddUser(ctx, u); err != nil {
		fmt.Println("UserManager AddUser s.dataStorage.AddUser", err)
		return err
	}

	if err := s.broker.Produce(u); err != nil {
		fmt.Println("UserManager AddUser s.broker.Produce", err)
		return err
	}

	fmt.Println("User added")

	return nil
}

func (s *UserManager) FetchAllUsers(ctx context.Context) ([]dto.User, error) {
	var usersList []data.User

	isCached, err := s.cacheStorage.IsExists(ctx)
	if err != nil {
		fmt.Println("UserManager FetchAllUsers s.cacheStorage.IsExists", err)
		return nil, err
	}

	if isCached {
		usersList, err = s.cacheStorage.GetAllUsers(ctx)
		if err != nil {
			fmt.Println("UserManager FetchAllUsers s.cacheStorage.GetAllUsers", err)
			return nil, err
		}
	} else {
		usersList, err = s.dataStorage.FetchAllUsers(ctx)
		if err != nil {
			fmt.Println("UserManager FetchAllUsers s.dataStorage.FetchAllUsers", err)
			return nil, err
		}

		if err := s.cacheStorage.SetAllUsers(ctx, usersList); err != nil {
			fmt.Println("UserManager FetchAllUsers s.cacheStorage.SetAllUsers", err)
			return nil, err
		}
	}

	convertedList, err := convert(usersList)
	if err != nil {
		fmt.Println("UserManager FetchAllUsers convert", err)
		return nil, err
	}

	fmt.Println("Users fetched")

	return convertedList, nil
}

func (s *UserManager) DeleteUser(ctx context.Context, id int) error {
	if err := s.dataStorage.DeleteUser(ctx, id); err != nil {
		fmt.Println("UserManager DeleteUser s.dataStorage.DeleteUser", err)
		return err
	}

	fmt.Println("User deleted")

	return nil
}

func convert(l []data.User) ([]dto.User, error) {
	usersList := []dto.User{}

	if len(l) == 0 {
		return nil, errors.New("UserManager convert list is empty")
	}

	for _, v := range l {
		usersList = append(usersList, dto.User{
			ID:       v.ID,
			UserName: v.UserName,
		})
	}

	return usersList, nil
}
