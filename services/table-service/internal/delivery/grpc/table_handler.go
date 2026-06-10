package grpc

import (
	"context"
	"errors"

	pb "restaurant-management/proto/table"
	"restaurant-management/services/table-service/internal/domain"
	"restaurant-management/services/table-service/internal/usecase"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// TableHandler handles gRPC requests for the table service.
type TableHandler struct {
	pb.UnimplementedTableServiceServer
	uc *usecase.TableUseCase
}

func NewTableHandler(uc *usecase.TableUseCase) *TableHandler {
	return &TableHandler{uc: uc}
}

func (h *TableHandler) CreateTable(ctx context.Context, req *pb.CreateTableRequest) (*pb.CreateTableResponse, error) {
	if req.TableNumber <= 0 {
		return nil, status.Error(codes.InvalidArgument, "table_number must be a positive integer")
	}
	if req.Capacity <= 0 {
		return nil, status.Error(codes.InvalidArgument, "capacity must be positive")
	}

	table, err := h.uc.CreateTable(ctx, int(req.TableNumber), int(req.Capacity))
	if err != nil {
		return nil, toGRPCError(err)
	}
	return &pb.CreateTableResponse{Table: toProto(table), Success: true, Message: "table created"}, nil
}

func (h *TableHandler) GetTable(ctx context.Context, req *pb.GetTableRequest) (*pb.GetTableResponse, error) {
	if req.TableId == "" {
		return nil, status.Error(codes.InvalidArgument, "table_id is required")
	}
	table, err := h.uc.GetTable(ctx, req.TableId)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return &pb.GetTableResponse{Table: toProto(table), Success: true, Message: "ok"}, nil
}

func (h *TableHandler) UpdateTable(ctx context.Context, req *pb.UpdateTableRequest) (*pb.UpdateTableResponse, error) {
	if req.TableId == "" {
		return nil, status.Error(codes.InvalidArgument, "table_id is required")
	}
	table, err := h.uc.UpdateTable(ctx, req.TableId, int(req.TableNumber), int(req.Capacity))
	if err != nil {
		return nil, toGRPCError(err)
	}
	return &pb.UpdateTableResponse{Table: toProto(table), Success: true, Message: "table updated"}, nil
}

func (h *TableHandler) DeleteTable(ctx context.Context, req *pb.DeleteTableRequest) (*pb.DeleteTableResponse, error) {
	if req.TableId == "" {
		return nil, status.Error(codes.InvalidArgument, "table_id is required")
	}
	if err := h.uc.DeleteTable(ctx, req.TableId); err != nil {
		return nil, toGRPCError(err)
	}
	return &pb.DeleteTableResponse{Success: true, Message: "table deleted"}, nil
}

func (h *TableHandler) ListTables(ctx context.Context, req *pb.ListTablesRequest) (*pb.ListTablesResponse, error) {
	page := int(req.Page)
	if page <= 0 {
		page = 1
	}
	pageSize := int(req.PageSize)
	if pageSize <= 0 {
		pageSize = 10
	}

	tables, total, err := h.uc.ListTables(ctx, page, pageSize, protoStatusToDomain(req.Status))
	if err != nil {
		return nil, toGRPCError(err)
	}

	protoTables := make([]*pb.Table, len(tables))
	for i, t := range tables {
		protoTables[i] = toProto(t)
	}
	return &pb.ListTablesResponse{
		Tables:   protoTables,
		Total:    int32(total),
		Page:     int32(page),
		PageSize: int32(pageSize),
		Success:  true,
		Message:  "ok",
	}, nil
}

func (h *TableHandler) UpdateTableStatus(ctx context.Context, req *pb.UpdateTableStatusRequest) (*pb.UpdateTableStatusResponse, error) {
	if req.TableId == "" {
		return nil, status.Error(codes.InvalidArgument, "table_id is required")
	}
	tableStatus := protoStatusToDomain(req.Status)
	if tableStatus == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid status")
	}
	table, err := h.uc.UpdateTableStatus(ctx, req.TableId, tableStatus)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return &pb.UpdateTableStatusResponse{Table: toProto(table), Success: true, Message: "status updated"}, nil
}

func (h *TableHandler) GetAvailableTables(ctx context.Context, req *pb.GetAvailableTablesRequest) (*pb.GetAvailableTablesResponse, error) {
	tables, err := h.uc.GetAvailableTables(ctx, int(req.MinCapacity))
	if err != nil {
		return nil, toGRPCError(err)
	}
	protoTables := make([]*pb.Table, len(tables))
	for i, t := range tables {
		protoTables[i] = toProto(t)
	}
	return &pb.GetAvailableTablesResponse{Tables: protoTables, Success: true, Message: "ok"}, nil
}

// --- helpers ---

func toProto(t *domain.Table) *pb.Table {
	if t == nil {
		return nil
	}
	return &pb.Table{
		TableId:     t.ID,
		TableNumber: int32(t.TableNumber),
		Capacity:    int32(t.Capacity),
		Status:      domainStatusToProto(t.Status),
		CreatedAt:   timestamppb.New(t.CreatedAt),
		UpdatedAt:   timestamppb.New(t.UpdatedAt),
	}
}

func domainStatusToProto(s domain.TableStatus) pb.TableStatus {
	switch s {
	case domain.StatusAvailable:
		return pb.TableStatus_STATUS_AVAILABLE
	case domain.StatusCleaning:
		return pb.TableStatus_STATUS_CLEANING
	case domain.StatusOutOfService:
		return pb.TableStatus_STATUS_OUT_OF_SERVICE
	default:
		return pb.TableStatus_STATUS_UNKNOWN
	}
}

func protoStatusToDomain(s pb.TableStatus) domain.TableStatus {
	switch s {
	case pb.TableStatus_STATUS_AVAILABLE:
		return domain.StatusAvailable
	case pb.TableStatus_STATUS_CLEANING:
		return domain.StatusCleaning
	case pb.TableStatus_STATUS_OUT_OF_SERVICE:
		return domain.StatusOutOfService
	default:
		return ""
	}
}

func toGRPCError(err error) error {
	switch {
	case errors.Is(err, domain.ErrTableNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, domain.ErrTableNumberAlreadyExists):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, domain.ErrTableNotAvailable):
		return status.Error(codes.FailedPrecondition, err.Error())
	case errors.Is(err, domain.ErrInvalidTableNumber),
		errors.Is(err, domain.ErrInvalidCapacity),
		errors.Is(err, domain.ErrCapacityTooLarge),
		errors.Is(err, domain.ErrInvalidStatus),
		errors.Is(err, domain.ErrInvalidTableID):
		return status.Error(codes.InvalidArgument, err.Error())
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
