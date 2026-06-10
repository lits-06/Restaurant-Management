package grpc

import (
	"context"
	"errors"

	"restaurant-management/services/schedule-service/internal/domain"
	"restaurant-management/services/schedule-service/internal/usecase"

	pb "restaurant-management/proto/schedule"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ScheduleHandler struct {
	pb.UnimplementedScheduleServiceServer
	uc *usecase.ScheduleUseCase
}

func NewScheduleHandler(uc *usecase.ScheduleUseCase) *ScheduleHandler {
	return &ScheduleHandler{uc: uc}
}

func (h *ScheduleHandler) CreateShift(ctx context.Context, req *pb.CreateShiftRequest) (*pb.CreateShiftResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	if req.Date == "" {
		return nil, status.Error(codes.InvalidArgument, "date is required")
	}
	if req.StartTime == "" || req.EndTime == "" {
		return nil, status.Error(codes.InvalidArgument, "start_time and end_time are required")
	}
	if req.Role == "" {
		return nil, status.Error(codes.InvalidArgument, "role is required")
	}

	shift, err := h.uc.CreateShift(ctx, req.UserId, req.Date, req.StartTime, req.EndTime, req.Role, req.Notes, req.CreatedBy)
	if err != nil {
		return nil, handleError(err)
	}
	return &pb.CreateShiftResponse{Shift: toProto(shift), Success: true, Message: "Shift created"}, nil
}

func (h *ScheduleHandler) GetShift(ctx context.Context, req *pb.GetShiftRequest) (*pb.GetShiftResponse, error) {
	if req.ShiftId == "" {
		return nil, status.Error(codes.InvalidArgument, "shift_id is required")
	}
	shift, err := h.uc.GetShift(ctx, req.ShiftId)
	if err != nil {
		return nil, handleError(err)
	}
	return &pb.GetShiftResponse{Shift: toProto(shift), Success: true, Message: "OK"}, nil
}

func (h *ScheduleHandler) UpdateShift(ctx context.Context, req *pb.UpdateShiftRequest) (*pb.UpdateShiftResponse, error) {
	if req.ShiftId == "" {
		return nil, status.Error(codes.InvalidArgument, "shift_id is required")
	}
	shift, err := h.uc.UpdateShift(ctx, req.ShiftId, req.Date, req.StartTime, req.EndTime, req.Notes)
	if err != nil {
		return nil, handleError(err)
	}
	return &pb.UpdateShiftResponse{Shift: toProto(shift), Success: true, Message: "Shift updated"}, nil
}

func (h *ScheduleHandler) DeleteShift(ctx context.Context, req *pb.DeleteShiftRequest) (*pb.DeleteShiftResponse, error) {
	if req.ShiftId == "" {
		return nil, status.Error(codes.InvalidArgument, "shift_id is required")
	}
	if err := h.uc.DeleteShift(ctx, req.ShiftId); err != nil {
		return nil, handleError(err)
	}
	return &pb.DeleteShiftResponse{Success: true, Message: "Shift deleted"}, nil
}

func (h *ScheduleHandler) ListShifts(ctx context.Context, req *pb.ListShiftsRequest) (*pb.ListShiftsResponse, error) {
	shifts, total, err := h.uc.ListShifts(ctx, req.Month, req.UserId, req.Role, int(req.Page), int(req.PageSize))
	if err != nil {
		return nil, handleError(err)
	}

	pbShifts := make([]*pb.Shift, len(shifts))
	for i, s := range shifts {
		pbShifts[i] = toProto(s)
	}
	return &pb.ListShiftsResponse{
		Shifts:   pbShifts,
		Total:    int32(total),
		Page:     req.Page,
		PageSize: req.PageSize,
		Success:  true,
		Message:  "OK",
	}, nil
}

func toProto(s *domain.Shift) *pb.Shift {
	return &pb.Shift{
		ShiftId:   s.ShiftID,
		UserId:    s.UserID,
		Date:      s.Date,
		StartTime: s.StartTime,
		EndTime:   s.EndTime,
		Role:      s.Role,
		Notes:     s.Notes,
		CreatedBy: s.CreatedBy,
		CreatedAt: timestamppb.New(s.CreatedAt),
		UpdatedAt: timestamppb.New(s.UpdatedAt),
	}
}

func handleError(err error) error {
	if errors.Is(err, domain.ErrShiftNotFound) {
		return status.Error(codes.NotFound, err.Error())
	}
	if errors.Is(err, domain.ErrUserIDRequired) ||
		errors.Is(err, domain.ErrInvalidDate) ||
		errors.Is(err, domain.ErrInvalidTime) ||
		errors.Is(err, domain.ErrEndTimeBeforeStart) ||
		errors.Is(err, domain.ErrInvalidRole) {
		return status.Error(codes.InvalidArgument, err.Error())
	}
	return status.Error(codes.Internal, "internal server error")
}
