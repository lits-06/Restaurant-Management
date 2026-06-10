package grpc

import (
	"context"
	"errors"

	"restaurant-management/services/user-service/internal/domain"
	"restaurant-management/services/user-service/internal/usecase"

	pb "restaurant-management/proto/user"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserHandler struct {
	pb.UnimplementedUserServiceServer
	userUseCase *usecase.UserUseCase
}

func NewUserHandler(userUseCase *usecase.UserUseCase) *UserHandler {
	return &UserHandler{userUseCase: userUseCase}
}

func (h *UserHandler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}
	if req.Username == "" {
		return nil, status.Error(codes.InvalidArgument, "username is required")
	}
	if req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	user, err := h.userUseCase.CreateUser(ctx, req.Email, req.Username, req.FullName, req.Phone, req.Password, protoRolesToDomain(req.Roles))
	if err != nil {
		return nil, handleError(err)
	}
	return &pb.CreateUserResponse{User: domainUserToProto(user), Success: true, Message: "User created successfully"}, nil
}

func (h *UserHandler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	user, err := h.userUseCase.GetUser(ctx, req.UserId)
	if err != nil {
		return nil, handleError(err)
	}
	return &pb.GetUserResponse{User: domainUserToProto(user), Success: true, Message: "User retrieved successfully"}, nil
}

func (h *UserHandler) GetUserByEmail(ctx context.Context, req *pb.GetUserByEmailRequest) (*pb.GetUserResponse, error) {
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}
	user, err := h.userUseCase.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, handleError(err)
	}
	return &pb.GetUserResponse{User: domainUserToProto(user), Success: true, Message: "User retrieved successfully"}, nil
}

func (h *UserHandler) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	user, err := h.userUseCase.UpdateUser(ctx, req.UserId, req.Email, req.Username, req.FullName, req.Phone, protoStatusToDomain(req.Status))
	if err != nil {
		return nil, handleError(err)
	}
	return &pb.UpdateUserResponse{User: domainUserToProto(user), Success: true, Message: "User updated successfully"}, nil
}

func (h *UserHandler) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	if err := h.userUseCase.DeleteUser(ctx, req.UserId); err != nil {
		return nil, handleError(err)
	}
	return &pb.DeleteUserResponse{Success: true, Message: "User deleted successfully"}, nil
}

func (h *UserHandler) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	page := int(req.Page)
	if page <= 0 {
		page = 1
	}
	pageSize := int(req.PageSize)
	if pageSize <= 0 {
		pageSize = 20
	}

	users, total, err := h.userUseCase.ListUsers(ctx, page, pageSize, protoStatusToDomain(req.Status), protoRoleToDomain(req.Role), req.Keyword)
	if err != nil {
		return nil, handleError(err)
	}

	protoUsers := make([]*pb.User, len(users))
	for i, u := range users {
		protoUsers[i] = domainUserToProto(u)
	}
	return &pb.ListUsersResponse{
		Users: protoUsers, Total: int32(total),
		Page: int32(page), PageSize: int32(pageSize),
		Success: true, Message: "Users retrieved successfully",
	}, nil
}

func (h *UserHandler) AssignRole(ctx context.Context, req *pb.AssignRoleRequest) (*pb.AssignRoleResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	if len(req.Roles) == 0 {
		return nil, status.Error(codes.InvalidArgument, "at least one role is required")
	}
	if err := h.userUseCase.AssignRoles(ctx, req.UserId, protoRolesToDomain(req.Roles)); err != nil {
		return nil, handleError(err)
	}
	return &pb.AssignRoleResponse{Success: true, Message: "Roles assigned successfully"}, nil
}

func (h *UserHandler) GetUserRoles(ctx context.Context, req *pb.GetUserRolesRequest) (*pb.GetUserRolesResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	roles, err := h.userUseCase.GetUserRoles(ctx, req.UserId)
	if err != nil {
		return nil, handleError(err)
	}
	return &pb.GetUserRolesResponse{Roles: domainRolesToProto(roles), Success: true, Message: "User roles retrieved successfully"}, nil
}

func (h *UserHandler) VerifyCredentials(ctx context.Context, req *pb.VerifyCredentialsRequest) (*pb.VerifyCredentialsResponse, error) {
	if req.Email == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "email and password are required")
	}
	user, err := h.userUseCase.VerifyCredentials(ctx, req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidCredentials):
			return &pb.VerifyCredentialsResponse{Success: false, Message: "invalid credentials"}, nil
		case errors.Is(err, domain.ErrAccountSuspended):
			return &pb.VerifyCredentialsResponse{Success: false, Message: "account is suspended"}, nil
		default:
			return nil, status.Error(codes.Internal, "internal server error")
		}
	}
	return &pb.VerifyCredentialsResponse{
		Success: true,
		Message: "credentials verified",
		UserId:  user.ID,
		Email:   user.Email,
		Roles:   domainRolesToProto(user.Roles),
	}, nil
}

func (h *UserHandler) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	if req.OldPassword == "" {
		return nil, status.Error(codes.InvalidArgument, "old_password is required")
	}
	if req.NewPassword == "" {
		return nil, status.Error(codes.InvalidArgument, "new_password is required")
	}
	if err := h.userUseCase.ChangePassword(ctx, req.UserId, req.OldPassword, req.NewPassword); err != nil {
		return nil, handleError(err)
	}
	return &pb.ChangePasswordResponse{Success: true, Message: "Password changed successfully"}, nil
}

// ── conversion helpers ────────────────────────────────────────

func domainUserToProto(u *domain.User) *pb.User {
	if u == nil {
		return nil
	}
	return &pb.User{
		UserId:    u.ID,
		Email:     u.Email,
		Username:  u.Username,
		FullName:  u.FullName,
		Phone:     u.Phone,
		Roles:     domainRolesToProto(u.Roles),
		Status:    domainStatusToProto(u.Status),
		CreatedAt: timestamppb.New(u.CreatedAt),
		UpdatedAt: timestamppb.New(u.UpdatedAt),
	}
}

func domainRolesToProto(roles []domain.UserRole) []pb.UserRole {
	out := make([]pb.UserRole, len(roles))
	for i, r := range roles {
		out[i] = domainRoleToProto(r)
	}
	return out
}

func domainRoleToProto(r domain.UserRole) pb.UserRole {
	switch r {
	case domain.RoleUser:
		return pb.UserRole_ROLE_USER
	case domain.RoleManager:
		return pb.UserRole_ROLE_MANAGER
	case domain.RoleChef:
		return pb.UserRole_ROLE_CHEF
	case domain.RoleWaiter:
		return pb.UserRole_ROLE_WAITER
	case domain.RoleAdmin:
		return pb.UserRole_ROLE_ADMIN
	default:
		return pb.UserRole_ROLE_UNKNOWN
	}
}

func protoRoleToDomain(r pb.UserRole) domain.UserRole {
	switch r {
	case pb.UserRole_ROLE_USER:
		return domain.RoleUser
	case pb.UserRole_ROLE_MANAGER:
		return domain.RoleManager
	case pb.UserRole_ROLE_CHEF:
		return domain.RoleChef
	case pb.UserRole_ROLE_WAITER:
		return domain.RoleWaiter
	case pb.UserRole_ROLE_ADMIN:
		return domain.RoleAdmin
	default:
		return ""
	}
}

func protoRolesToDomain(protoRoles []pb.UserRole) []domain.UserRole {
	roles := make([]domain.UserRole, 0, len(protoRoles))
	for _, r := range protoRoles {
		if d := protoRoleToDomain(r); d != "" {
			roles = append(roles, d)
		}
	}
	return roles
}

func domainStatusToProto(s domain.UserStatus) pb.UserStatus {
	switch s {
	case domain.StatusActive:
		return pb.UserStatus_STATUS_ACTIVE
	case domain.StatusInactive:
		return pb.UserStatus_STATUS_INACTIVE
	case domain.StatusSuspended:
		return pb.UserStatus_STATUS_SUSPENDED
	default:
		return pb.UserStatus_STATUS_UNKNOWN
	}
}

func protoStatusToDomain(s pb.UserStatus) domain.UserStatus {
	switch s {
	case pb.UserStatus_STATUS_ACTIVE:
		return domain.StatusActive
	case pb.UserStatus_STATUS_INACTIVE:
		return domain.StatusInactive
	case pb.UserStatus_STATUS_SUSPENDED:
		return domain.StatusSuspended
	default:
		return ""
	}
}

func handleError(err error) error {
	if err == nil {
		return nil
	}
	switch {
	case errors.Is(err, domain.ErrUserNotFound):
		return status.Error(codes.NotFound, "user not found")
	case errors.Is(err, domain.ErrEmailAlreadyExists):
		return status.Error(codes.AlreadyExists, "email already exists")
	case errors.Is(err, domain.ErrUsernameAlreadyExists):
		return status.Error(codes.AlreadyExists, "username already exists")
	case errors.Is(err, domain.ErrUserAlreadyExists):
		return status.Error(codes.AlreadyExists, "user already exists")
	case errors.Is(err, domain.ErrInvalidEmail),
		errors.Is(err, domain.ErrInvalidEmailFormat),
		errors.Is(err, domain.ErrInvalidUsername),
		errors.Is(err, domain.ErrUsernameTooShort),
		errors.Is(err, domain.ErrInvalidFullName),
		errors.Is(err, domain.ErrInvalidPassword),
		errors.Is(err, domain.ErrInvalidRole),
		errors.Is(err, domain.ErrNoRolesAssigned):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, domain.ErrInvalidCredentials):
		return status.Error(codes.Unauthenticated, "invalid credentials")
	case errors.Is(err, domain.ErrAccountSuspended):
		return status.Error(codes.PermissionDenied, "account is suspended")
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
