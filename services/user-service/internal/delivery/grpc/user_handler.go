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

// UserHandler handles gRPC requests for user service
type UserHandler struct {
	pb.UnimplementedUserServiceServer
	userUseCase *usecase.UserUseCase
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(userUseCase *usecase.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
	}
}

// CreateUser creates a new user
func (h *UserHandler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	// Validate request
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}
	if req.Username == "" {
		return nil, status.Error(codes.InvalidArgument, "username is required")
	}
	if req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	// Convert proto roles to domain roles
	roles := protoRolesToDomain(req.Roles)

	// Create user
	user, err := h.userUseCase.CreateUser(ctx, req.Email, req.Username, req.FullName, req.Phone, req.Password, roles)
	if err != nil {
		return nil, handleError(err)
	}

	return &pb.CreateUserResponse{
		User:    domainUserToProto(user),
		Success: true,
		Message: "User created successfully",
	}, nil
}

// GetUser retrieves a user by ID
func (h *UserHandler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	user, err := h.userUseCase.GetUser(ctx, req.UserId)
	if err != nil {
		return nil, handleError(err)
	}

	return &pb.GetUserResponse{
		User:    domainUserToProto(user),
		Success: true,
		Message: "User retrieved successfully",
	}, nil
}

// UpdateUser updates user information
func (h *UserHandler) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	// Convert proto status to domain status
	status := protoStatusToDomain(req.Status)

	user, err := h.userUseCase.UpdateUser(ctx, req.UserId, req.Email, req.Username, req.FullName, req.Phone, status)
	if err != nil {
		return nil, handleError(err)
	}

	return &pb.UpdateUserResponse{
		User:    domainUserToProto(user),
		Success: true,
		Message: "User updated successfully",
	}, nil
}

// DeleteUser deletes a user
func (h *UserHandler) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	if err := h.userUseCase.DeleteUser(ctx, req.UserId); err != nil {
		return nil, handleError(err)
	}

	return &pb.DeleteUserResponse{
		Success: true,
		Message: "User deleted successfully",
	}, nil
}

// ListUsers retrieves users with filters and pagination
func (h *UserHandler) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	// Default pagination
	page := int(req.Page)
	if page <= 0 {
		page = 1
	}
	pageSize := int(req.PageSize)
	if pageSize <= 0 {
		pageSize = 10
	}

	// Convert proto filters to domain
	status := protoStatusToDomain(req.Status)
	role := protoRoleToDomain(req.Role)

	users, total, err := h.userUseCase.ListUsers(ctx, page, pageSize, status, role)
	if err != nil {
		return nil, handleError(err)
	}

	// Convert users to proto
	protoUsers := make([]*pb.User, len(users))
	for i, user := range users {
		protoUsers[i] = domainUserToProto(user)
	}

	return &pb.ListUsersResponse{
		Users:    protoUsers,
		Total:    int32(total),
		Page:     int32(page),
		PageSize: int32(pageSize),
		Success:  true,
		Message:  "Users retrieved successfully",
	}, nil
}

// AssignRole assigns roles to a user
func (h *UserHandler) AssignRole(ctx context.Context, req *pb.AssignRoleRequest) (*pb.AssignRoleResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	if len(req.Roles) == 0 {
		return nil, status.Error(codes.InvalidArgument, "at least one role is required")
	}

	// Convert proto roles to domain
	roles := protoRolesToDomain(req.Roles)

	if err := h.userUseCase.AssignRoles(ctx, req.UserId, roles); err != nil {
		return nil, handleError(err)
	}

	return &pb.AssignRoleResponse{
		Success: true,
		Message: "Roles assigned successfully",
	}, nil
}

// GetUserRoles retrieves the roles of a user
func (h *UserHandler) GetUserRoles(ctx context.Context, req *pb.GetUserRolesRequest) (*pb.GetUserRolesResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	roles, err := h.userUseCase.GetUserRoles(ctx, req.UserId)
	if err != nil {
		return nil, handleError(err)
	}

	// Convert domain roles to proto
	protoRoles := domainRolesToProto(roles)

	return &pb.GetUserRolesResponse{
		Roles:   protoRoles,
		Success: true,
		Message: "User roles retrieved successfully",
	}, nil
}

// Helper functions

// domainUserToProto converts domain.User to pb.User
func domainUserToProto(user *domain.User) *pb.User {
	if user == nil {
		return nil
	}

	return &pb.User{
		UserId:    user.ID,
		Email:     user.Email,
		Username:  user.Username,
		FullName:  user.FullName,
		Phone:     user.Phone,
		Roles:     domainRolesToProto(user.Roles),
		Status:    domainStatusToProto(user.Status),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}
}

// domainRolesToProto converts domain roles to proto roles
func domainRolesToProto(roles []domain.UserRole) []pb.UserRole {
	protoRoles := make([]pb.UserRole, len(roles))
	for i, role := range roles {
		protoRoles[i] = domainRoleToProto(role)
	}
	return protoRoles
}

// domainRoleToProto converts a single domain role to proto role
func domainRoleToProto(role domain.UserRole) pb.UserRole {
	switch role {
	case domain.RoleAdmin:
		return pb.UserRole_ROLE_ADMIN
	case domain.RoleManager:
		return pb.UserRole_ROLE_MANAGER
	case domain.RoleWaiter:
		return pb.UserRole_ROLE_WAITER
	case domain.RoleChef:
		return pb.UserRole_ROLE_CHEF
	case domain.RoleCashier:
		return pb.UserRole_ROLE_CASHIER
	default:
		return pb.UserRole_ROLE_UNKNOWN
	}
}

// protoRoleToDomain converts proto role to domain role
func protoRoleToDomain(role pb.UserRole) domain.UserRole {
	switch role {
	case pb.UserRole_ROLE_ADMIN:
		return domain.RoleAdmin
	case pb.UserRole_ROLE_MANAGER:
		return domain.RoleManager
	case pb.UserRole_ROLE_WAITER:
		return domain.RoleWaiter
	case pb.UserRole_ROLE_CHEF:
		return domain.RoleChef
	case pb.UserRole_ROLE_CASHIER:
		return domain.RoleCashier
	default:
		return ""
	}
}

// protoRolesToDomain converts proto roles to domain roles
func protoRolesToDomain(protoRoles []pb.UserRole) []domain.UserRole {
	roles := make([]domain.UserRole, 0, len(protoRoles))
	for _, protoRole := range protoRoles {
		if role := protoRoleToDomain(protoRole); role != "" {
			roles = append(roles, role)
		}
	}
	return roles
}

// domainStatusToProto converts domain status to proto status
func domainStatusToProto(status domain.UserStatus) pb.UserStatus {
	switch status {
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

// protoStatusToDomain converts proto status to domain status
func protoStatusToDomain(status pb.UserStatus) domain.UserStatus {
	switch status {
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

// handleError converts domain errors to gRPC status codes
func handleError(err error) error {
	if err == nil {
		return nil
	}

	// Check for specific domain errors
	switch {
	case errors.Is(err, domain.ErrUserNotFound):
		return status.Error(codes.NotFound, "user not found")
	case errors.Is(err, domain.ErrUserAlreadyExists):
		return status.Error(codes.AlreadyExists, "user already exists")
	case errors.Is(err, domain.ErrEmailAlreadyExists):
		return status.Error(codes.AlreadyExists, "email already exists")
	case errors.Is(err, domain.ErrUsernameAlreadyExists):
		return status.Error(codes.AlreadyExists, "username already exists")
	case errors.Is(err, domain.ErrInvalidEmail),
		errors.Is(err, domain.ErrInvalidEmailFormat),
		errors.Is(err, domain.ErrInvalidUsername),
		errors.Is(err, domain.ErrUsernameTooShort),
		errors.Is(err, domain.ErrInvalidFullName),
		errors.Is(err, domain.ErrInvalidPassword),
		errors.Is(err, domain.ErrInvalidRole),
		errors.Is(err, domain.ErrNoRolesAssigned):
		return status.Error(codes.InvalidArgument, err.Error())
	default:
		// Internal server error for unknown errors
		return status.Error(codes.Internal, "internal server error")
	}
}
