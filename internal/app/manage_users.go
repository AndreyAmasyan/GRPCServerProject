package app

import (
	"context"
	"errors"
	"fmt"
	"test/internal/dto"
	"test/internal/service"
	api "test/pkg/grpc"
)

// GRPCServer ...
type GRPCServer struct {
	processor service.Processor
}

func New(p service.Processor) *GRPCServer {
	return &GRPCServer{
		processor: p,
	}
}

func (s *GRPCServer) AddUser(ctx context.Context, req *api.AddUserRequest) (*api.AddUserResponse, error) {
	reqUser := req.GetUser()
	p := s.processor

	dtoUser := dto.User{
		ID:       int(reqUser.GetUid()),
		UserName: reqUser.GetName(),
	}

	if err := p.AddUser(ctx, dtoUser); err != nil {
		fmt.Println("GRPCServer AddUser p.AddUser", err)
		return nil, err
	}

	return &api.AddUserResponse{
		User: reqUser,
	}, nil
}

func (s *GRPCServer) FetchAllUsers(ctx context.Context, req *api.Empty) (*api.FetchAllUsersResponse, error) {
	p := s.processor

	usersList, err := p.FetchAllUsers(ctx)
	if err != nil {
		fmt.Println("GRPCServer FetchAllUsers p.FetchAllUsers", err)
		return nil, err
	}

	convertedList, err := convert(usersList)
	if err != nil {
		fmt.Println("GRPCServer FetchAllUsers convert", err)
		return nil, err
	}

	return &api.FetchAllUsersResponse{
		User: convertedList,
	}, nil
}

func (s *GRPCServer) DeleteUser(ctx context.Context, req *api.DeleteUserRequest) (*api.DeleteUserResponse, error) {
	p := s.processor

	if err := p.DeleteUser(ctx, int(req.GetUid())); err != nil {
		fmt.Println("GRPCServer DeleteUser p.DeleteUser", err)
		return nil, err
	}

	return &api.DeleteUserResponse{
		Uid: req.GetUid(),
	}, nil
}

func convert(l []dto.User) ([]*api.User, error) {
	usersList := []*api.User{}

	if len(l) == 0 {
		return nil, errors.New("users list is empty")
	}

	for _, v := range l {
		usersList = append(usersList, &api.User{
			Uid:  int32(v.ID),
			Name: v.UserName,
		})
	}

	return usersList, nil
}
