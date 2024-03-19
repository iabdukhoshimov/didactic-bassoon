package service

import (
	"context"
	"encoding/json"

	pbUser "gitlab.com/tour/generated/user_service"
	"gitlab.com/tour/internal/core/repository"
	"gitlab.com/tour/internal/core/repository/psql/sqlc"
	"gitlab.com/tour/internal/pkg/logger"
	"gitlab.com/tour/internal/pkg/security"
	"gitlab.com/tour/internal/pkg/serializer"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gopkg.in/guregu/null.v4/zero"
)

type UserService struct {
	pbUser.UnimplementedUserServiceServer
	store  repository.Store
	issuer security.IssuerInterface
}

func NewUserService(store repository.Store, issuer security.IssuerInterface) *UserService {
	return &UserService{
		store:  store,
		issuer: issuer,
	}
}

func (s *UserService) CreateUser(ctx context.Context, req *pbUser.CreateUserRequest) (resp *pbUser.CreatedResponse, err error) {
	var (
		dbReq sqlc.CreateUserParams
	)

	dbReq.Email = req.Email
	dbReq.FirstName = req.FirstName
	dbReq.LastName = req.LastName
	dbReq.Role = sqlc.Roles(req.Role)
	dbReq.OrganizationID = zero.StringFrom(req.OrganizationId)

	existingUser, err := s.store.GetUserByEmail(ctx, req.Email)
	if err == nil && existingUser == (sqlc.GetUserByEmailRow{}) {
		logger.Log.Error("!!!Get User By Email this Email Already exists  " + err.Error())
		return nil, status.Error(codes.AlreadyExists, "internal error")
	}

	hashedPassword, err := security.HashPassword(req.Password)
	if err != nil {
		logger.Log.Error("!!!CreateUser---> while hashing password " + err.Error())
		return nil, err
	}
	dbReq.Password = hashedPassword

	createdUserID, err := s.store.CreateUser(ctx, dbReq)
	if err != nil {
		logger.Log.Error("! while creating New User  " + err.Error())
		return nil, err
	}

	// orqReq := sqlc.CreateOrganizationEmployeeParams{
	// 	OrgID:  zero.StringFrom(req.OrganizationId),
	// 	UserID: createdUserID,
	// }

	// _, err = s.store.CreateOrganizationEmployee(ctx, orqReq)
	// if err != nil {
	// 	logger.Log.Error("! while creating Organization Employee " + err.Error())
	// 	return nil, err
	// }

	return &pbUser.CreatedResponse{
		Id: createdUserID,
	}, nil
}

func (s *UserService) GetUserById(ctx context.Context, req *pbUser.GetById) (*pbUser.User, error) {
	var respUser pbUser.User

	user, err := s.store.GetUser(ctx, req.Id)
	if err != nil {
		logger.Log.Error("!!!GetUser---> while getting admin " + err.Error())
		return nil, err
	}

	if err := serializer.MarshalUnMarshal(user, &respUser); err != nil {
		return nil, err
	}

	return &respUser, nil
}

func (s *UserService) GetUsers(ctx context.Context, req *pbUser.GetAllRequest) (*pbUser.UserList, error) {
	var respUsers pbUser.UserList

	users, err := s.store.GetUsers(ctx)
	if err != nil {
		logger.Log.Error("!!!GetAllUsers---> while getting admin " + err.Error())
		return nil, err
	}

	if err := serializer.MarshalUnMarshal(users, &respUsers.Users); err != nil {
		return nil, err
	}

	return &respUsers, nil
}

func (s *UserService) UpdateUserById(ctx context.Context, req *pbUser.UpdateUser) (*emptypb.Empty, error) {
	var (
		dbReq sqlc.UpdateUserParams
	)

	bytes, err := json.Marshal(req)
	if err != nil {
		logger.Log.Error("!!!UpdateUserById---> while marshiling to bytes " + err.Error())
		return nil, err
	}

	err = json.Unmarshal(bytes, &dbReq)
	if err != nil {
		logger.Log.Error("!!!UpdateUserById---> while unmarshiling to bytes " + err.Error())
		return nil, err
	}

	err = s.store.UpdateUser(ctx, dbReq)
	if err != nil {
		logger.Log.Error("!!!UpdateUserById---> while updating admin " + err.Error())
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *UserService) DeleteUser(ctx context.Context, req *pbUser.GetById) (*emptypb.Empty, error) {
	err := s.store.DeleteUser(ctx, req.Id)
	if err != nil {

		logger.Log.Error("There is an error while deleting admin at service level ==> " + err.Error())
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
