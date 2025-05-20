package auth_v1

import (
	"context"
	authv1 "github.com/SeiFlow-3P2/auth_service/pkg/proto/v1"
	verfic "github.com/SeiFlow-3P2/auth_service/pkg/utils/verifications"
	"github.com/golang/protobuf/ptypes/wrappers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type serverAPI struct {
	auth Auth
	authv1.UnimplementedAuthServiceServer
}
type Auth interface {
	SingUpByEmail(
		ctx context.Context,
		email string,
		password string,
		telegramID string,
	) (userID string, accessToken string, refreshToken string, message string, err error)

	SingUpByOauth(
		ctx context.Context,
		provider string,
		oauthToken string,
		telegramID string,
	) (userID string, accessToken string, refreshToken string, message string, err error)

	LoginByEmail(
		ctx context.Context,
		email string,
		password string,
		telegramID string,
	) (userID string, accessToken string, refreshToken string, message string, err error)

	LoginByOauth(
		ctx context.Context,
		provider string,
		oauthToken string,
	) (userID string, accessToken string, refreshToken string, message string, err error)

	RefreshToken(
		ctx context.Context,
		RefreshToken string,
	) (accessToken string, refreshToken string, err error)

	Logout(ctx context.Context,
		RefreshToken string)
	GetUserInfo(ctx context.Context, userID string) (
		id string,
		telegramId string,
		username string,
		email string,
		photoUrl string,
		subscription bool,
		createdAt string,
		updatedAt string, err error)
	HealthCheck(ctx context.Context) (status string, err error)
}

func Register(gRPCServer *grpc.Server, auth Auth) {
	authv1.RegisterAuthServiceServer(gRPCServer, &serverAPI{auth: auth})
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
		userID, accessToken, refreshToken, message, err := s.auth.SingUpByEmail(ctx, emailAuth.Email, emailAuth.Password, emailAuth.TelegramId.Value)
		if err != nil {
			return nil, status.Error(codes.Internal, "failed to sing up")
		}
		return &authv1.SignUpResponse{UserId: userID, AccessToken: accessToken, RefreshToken: refreshToken, Message: message}, nil

	} else if oAuth != nil {
		userID, accessToken, refreshToken, message, err := s.auth.SingUpByOauth(ctx, oAuth.Provider, oAuth.OauthToken, oAuth.TelegramId.Value)
		if err != nil {
			return nil, status.Error(codes.Internal, "failed to sing up")
		}
		return &authv1.SignUpResponse{UserId: userID, AccessToken: accessToken, RefreshToken: refreshToken, Message: message}, nil
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
	id, telegramId, username, email, photoUrl, subscription, createdAt, updatedAt, err := s.auth.GetUserInfo(ctx, in.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get user info")
	}
	return &authv1.GetUserInfoResponse{User: &authv1.UserInfo{Id: id, TelegramId: &wrappers.StringValue{Value: telegramId}, Username: username, Email: email,
		PhotoUrl: &wrappers.StringValue{Value: photoUrl}, Subscription: subscription,
		CreatedAt: createdAt, UpdatedAt: updatedAt}}, nil
}

func (s *serverAPI) HealthCheck(ctx context.Context, in *emptypb.Empty) (*authv1.HealthCheckResponse, error) {
	stat, err := s.auth.HealthCheck(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get HealthCheck")
	}
	return &authv1.HealthCheckResponse{Status: stat}, nil
}
