package grpc

import (
	"context"
	"errors"

	"restaurant-management/services/staff-service/internal/domain"
	"restaurant-management/services/staff-service/internal/usecase"

	pb "restaurant-management/proto/staff"

	"go.uber.org/zap"

	"restaurant-management/shared/pkg/logger"
)

// StaffHandler handles gRPC requests for staff service.
type StaffHandler struct {
	pb.UnimplementedStaffServiceServer
	staffUseCase *usecase.StaffUseCase
}

// NewStaffHandler creates a new StaffHandler.
func NewStaffHandler(staffUseCase *usecase.StaffUseCase) *StaffHandler {
	return &StaffHandler{
		staffUseCase: staffUseCase,
	}
}

// CreateStaff creates a new staff member.
func (h *StaffHandler) CreateStaff(ctx context.Context, req *pb.CreateStaffRequest) (*pb.CreateStaffResponse, error) {
	if req.Name == "" {
		return &pb.CreateStaffResponse{Success: false, Message: "name is required"}, nil
	}
	if req.Role == "" {
		return &pb.CreateStaffResponse{Success: false, Message: "role is required"}, nil
	}
	if req.Contact == "" {
		return &pb.CreateStaffResponse{Success: false, Message: "contact is required"}, nil
	}

	staff, err := h.staffUseCase.CreateStaff(ctx, req.Name, req.Role, req.Contact, req.Avatar)
	if err != nil {
		logger.Error("Failed to create staff", zap.Error(err))
		return &pb.CreateStaffResponse{Success: false, Message: err.Error()}, nil
	}

	return &pb.CreateStaffResponse{
		Staff:   convertStaffToProto(staff),
		Success: true,
		Message: "staff created successfully",
	}, nil
}

// GetStaff retrieves a staff member by ID.
func (h *StaffHandler) GetStaff(ctx context.Context, req *pb.GetStaffRequest) (*pb.GetStaffResponse, error) {
	if req.StaffId == "" {
		return &pb.GetStaffResponse{Success: false, Message: "staff_id is required"}, nil
	}

	staff, err := h.staffUseCase.GetStaff(ctx, req.StaffId)
	if err != nil {
		logger.Error("Failed to get staff", zap.Error(err))
		if errors.Is(err, domain.ErrStaffNotFound) {
			return &pb.GetStaffResponse{Success: false, Message: "staff not found"}, nil
		}
		return &pb.GetStaffResponse{Success: false, Message: err.Error()}, nil
	}

	return &pb.GetStaffResponse{Staff: convertStaffToProto(staff), Success: true, Message: "staff retrieved successfully"}, nil
}

// UpdateStaff updates staff information.
func (h *StaffHandler) UpdateStaff(ctx context.Context, req *pb.UpdateStaffRequest) (*pb.UpdateStaffResponse, error) {
	if req.StaffId == "" {
		return &pb.UpdateStaffResponse{Success: false, Message: "staff_id is required"}, nil
	}

	staff, err := h.staffUseCase.UpdateStaff(ctx, req.StaffId, req.Name, req.Role, req.Contact, req.Avatar)
	if err != nil {
		logger.Error("Failed to update staff", zap.Error(err))
		if errors.Is(err, domain.ErrStaffNotFound) {
			return &pb.UpdateStaffResponse{Success: false, Message: "staff not found"}, nil
		}
		return &pb.UpdateStaffResponse{Success: false, Message: err.Error()}, nil
	}

	return &pb.UpdateStaffResponse{Staff: convertStaffToProto(staff), Success: true, Message: "staff updated successfully"}, nil
}

// DeleteStaff deletes a staff member.
func (h *StaffHandler) DeleteStaff(ctx context.Context, req *pb.DeleteStaffRequest) (*pb.DeleteStaffResponse, error) {
	if req.StaffId == "" {
		return &pb.DeleteStaffResponse{Success: false, Message: "staff_id is required"}, nil
	}

	if err := h.staffUseCase.DeleteStaff(ctx, req.StaffId); err != nil {
		logger.Error("Failed to delete staff", zap.Error(err))
		if errors.Is(err, domain.ErrStaffNotFound) {
			return &pb.DeleteStaffResponse{Success: false, Message: "staff not found"}, nil
		}
		return &pb.DeleteStaffResponse{Success: false, Message: err.Error()}, nil
	}

	return &pb.DeleteStaffResponse{Success: true, Message: "staff deleted successfully"}, nil
}

// ListStaff retrieves staff members with pagination and keyword search.
func (h *StaffHandler) ListStaff(ctx context.Context, req *pb.ListStaffRequest) (*pb.ListStaffResponse, error) {
	page := int(req.Page)
	if page <= 0 {
		page = 1
	}
	pageSize := int(req.PageSize)
	if pageSize <= 0 {
		pageSize = 20
	}

	staffMembers, total, err := h.staffUseCase.ListStaff(ctx, page, pageSize, req.Keyword)
	if err != nil {
		logger.Error("Failed to list staff", zap.Error(err))
		return &pb.ListStaffResponse{Success: false, Message: err.Error()}, nil
	}

	protoStaff := make([]*pb.Staff, len(staffMembers))
	for i, staff := range staffMembers {
		protoStaff[i] = convertStaffToProto(staff)
	}

	return &pb.ListStaffResponse{
		Staff:    protoStaff,
		Total:    int32(total),
		Page:     int32(page),
		PageSize: int32(pageSize),
		Success:  true,
		Message:  "staff retrieved successfully",
	}, nil
}

func convertStaffToProto(staff *domain.Staff) *pb.Staff {
	if staff == nil {
		return nil
	}

	return &pb.Staff{
		StaffId: staff.StaffID,
		Name:    staff.Name,
		Role:    staff.Role,
		Contact: staff.Contact,
		Avatar:  staff.Avatar,
	}
}
