package auth_v1

import (
	"context"
	authv1 "github.com/SeiFlow-3P2/auth_service/pkg/proto/v1"
	verfic "github.com/SeiFlow-3P2/auth_service/pkg/utils/verifications"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"strconv"
)

type serverAPI struct {
	auth Auth
	authv1.UnimplementedAuthServiceServer
}
type Auth interface {
	SingUpByEmail(ctx context.Context, name string, email string, password []byte, telegramID uint) (userID uuid.UUID, accessToken string, refreshToken string, message string, err error)

	SingUpByOauth(
		ctx context.Context,
		provider string,
		oauthToken string,
		telegramID string,
	) (userID uuid.UUID, accessToken string, refreshToken string, message string, err error)

	LoginByEmail(
		ctx context.Context,
		email string,
		password []byte,
	) (userID uuid.UUID, accessToken string, refreshToken string, message string, err error)

	LoginByOauth(
		ctx context.Context,
		provider string,
		oauthToken string,
	) (userID uuid.UUID, accessToken string, refreshToken string, message string, err error)

	RefreshToken(
		ctx context.Context,
		RefreshToken string,
	) (accessToken string, refreshToken string, err error)

	Logout(ctx context.Context,
		userID uuid.UUID) (err error)
	UserInfo(ctx context.Context, userID uuid.UUID) (
		id string,
		telegramId uint,
		username string,
		email string,
		photoUrl string,
		createdAt string,
		updatedAt string, err error)
	HealthCheck(ctx context.Context) (status string, err error)
}

func Register(gRPCServer *grpc.Server, auth Auth) {
	authv1.RegisterAuthServiceServer(gRPCServer, &serverAPI{auth: auth})
}
func (s *serverAPI) logout(ctx context.Context, in *authv1.LogoutRequest) (*emptypb.Empty, error) {
	userId, err := uuid.Parse(in.GetUserId()) //TODO Дождаться обновления протофайла и сравнить, работает ли?
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user id")
	}

	err = s.auth.Logout(ctx, userId)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to logout")
	}

	return &emptypb.Empty{}, nil
}
func (s *serverAPI) refreshToken(ctx context.Context, in *authv1.RefreshTokenRequest) (*authv1.RefreshTokenResponse, error) {
	if in.GetRefreshToken() == "" {
		return nil, status.Error(codes.InvalidArgument, "no refresh token")
	}
	accessToken, refreshToken, err := s.auth.RefreshToken(ctx, in.GetRefreshToken())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to refresh token"+err.Error())
	}
	return &authv1.RefreshTokenResponse{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}
func (s *serverAPI) Login(ctx context.Context, in *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	if in.GetOauth() != nil {
		panic("not implemented")
	}
	if in.GetEmail() != nil {

		valid, err := verfic.VerifyEmail(in.GetEmail().Email)
		if err != nil || !valid {
			return nil, status.Error(codes.InvalidArgument, "invalid email")
		}

		if in.GetEmail().Password == "" {
			return nil, status.Error(codes.InvalidArgument, "No password")

		}
		userID, accessToken, refreshToken, message, err := s.auth.LoginByEmail(ctx, in.GetEmail().Email, []byte(in.GetEmail().Password))
		if err != nil {
			return nil, status.Error(codes.Internal, "failed to login")
		}
		return &authv1.LoginResponse{UserId: userID.String(), AccessToken: accessToken, RefreshToken: refreshToken, Message: message}, nil
	}
	return nil, status.Error(codes.InvalidArgument, "invalid auth")
}
func (s *serverAPI) SignUp(ctx context.Context, in *authv1.SignUpRequest) (*authv1.SignUpResponse, error) {
	oAuth := in.GetOauth()
	emailAuth := in.GetEmail()

	if emailAuth != nil {

		res, _ := verfic.VerifyEmail(emailAuth.Email)
		if !res {
			return nil, status.Error(codes.InvalidArgument, "invalid email")
		}
		if emailAuth.Password == "" {
			return nil, status.Error(codes.InvalidArgument, "No password")
		}
		userID, accessToken, refreshToken, message, err := s.auth.SingUpByEmail(ctx, emailAuth.Username, emailAuth.Email, emailAuth.Password, emailAuth.TelegramId.Value)
		if err != nil {
			return nil, status.Error(codes.Internal, "failed to sing up")
		}
		return &authv1.SignUpResponse{UserId: userID.String(), AccessToken: accessToken, RefreshToken: refreshToken, Message: message}, nil

	} else if oAuth != nil {
		userID, accessToken, refreshToken, message, err := s.auth.SingUpByOauth(ctx, oAuth.Provider, oAuth.OauthToken, oAuth.TelegramId.Value)
		if err != nil {
			return nil, status.Error(codes.Internal, "failed to sing up")
		}
		return &authv1.SignUpResponse{UserId: userID.String(), AccessToken: accessToken, RefreshToken: refreshToken, Message: message}, nil
	} else {
		return nil, status.Error(codes.InvalidArgument, "invalid auth")
	}
}

func (s *serverAPI) GetUserInfo(
	ctx context.Context,
	in *authv1.GetUserInfoRequest,
) (*authv1.GetUserInfoResponse, error) {

	if in.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "No user id")
	}
	userID, err := uuid.Parse(in.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user id")
	}
	id, telegramId, username, email, photoUrl, createdAt, updatedAt, err := s.auth.UserInfo(ctx, userID)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get user info")
	}
	return &authv1.GetUserInfoResponse{User: &authv1.UserInfo{Id: id, TelegramId: &wrappers.StringValue{Value: strconv.Itoa(int(telegramId))}, Subscription: false, Username: username, Email: email,
		PhotoUrl:  &wrappers.StringValue{Value: photoUrl},
		CreatedAt: createdAt, UpdatedAt: updatedAt},
	}, nil
}
func (s *serverAPI) HealthCheck(ctx context.Context, in *emptypb.Empty) (*authv1.HealthCheckResponse, error) {
	stat, err := s.auth.HealthCheck(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get HealthCheck")
	}
	return &authv1.HealthCheckResponse{Status: stat}, nil
}
