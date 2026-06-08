package grpc

import (
	"context"
	"errors"

	pb "restaurant-management/proto/table"
	"restaurant-management/services/table-service/internal/domain"
	"restaurant-management/services/table-service/internal/repository"
	"restaurant-management/services/table-service/internal/usecase"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// TableHandler handles gRPC requests for table service
type TableHandler struct {
	pb.UnimplementedTableServiceServer
	tableUseCase       *usecase.TableUseCase
	reservationUseCase *usecase.ReservationUseCase
}

// NewTableHandler creates a new TableHandler
func NewTableHandler(tableUseCase *usecase.TableUseCase, reservationUseCase *usecase.ReservationUseCase) *TableHandler {
	return &TableHandler{
		tableUseCase:       tableUseCase,
		reservationUseCase: reservationUseCase,
	}
}

// CreateTable creates a new table
func (h *TableHandler) CreateTable(ctx context.Context, req *pb.CreateTableRequest) (*pb.CreateTableResponse, error) {
	// Validate request
	if req.TableNumber == "" {
		return nil, status.Error(codes.InvalidArgument, "table_number is required")
	}
	if req.Capacity <= 0 {
		return nil, status.Error(codes.InvalidArgument, "capacity must be positive")
	}
	if req.Location == "" {
		return nil, status.Error(codes.InvalidArgument, "location is required")
	}

	// Create table
	table, err := h.tableUseCase.CreateTable(ctx, req.TableNumber, int(req.Capacity), req.Location)
	if err != nil {
		return nil, handleError(err)
	}

	return &pb.CreateTableResponse{
		Table:   domainTableToProto(table),
		Success: true,
		Message: "Table created successfully",
	}, nil
}

// GetTable retrieves a table by ID
func (h *TableHandler) GetTable(ctx context.Context, req *pb.GetTableRequest) (*pb.GetTableResponse, error) {
	if req.TableId == "" {
		return nil, status.Error(codes.InvalidArgument, "table_id is required")
	}

	table, err := h.tableUseCase.GetTable(ctx, req.TableId)
	if err != nil {
		return nil, handleError(err)
	}

	return &pb.GetTableResponse{
		Table:   domainTableToProto(table),
		Success: true,
		Message: "Table retrieved successfully",
	}, nil
}

// UpdateTable updates table information
func (h *TableHandler) UpdateTable(ctx context.Context, req *pb.UpdateTableRequest) (*pb.UpdateTableResponse, error) {
	if req.TableId == "" {
		return nil, status.Error(codes.InvalidArgument, "table_id is required")
	}

	table, err := h.tableUseCase.UpdateTable(ctx, req.TableId, req.TableNumber, int(req.Capacity), req.Location)
	if err != nil {
		return nil, handleError(err)
	}

	return &pb.UpdateTableResponse{
		Table:   domainTableToProto(table),
		Success: true,
		Message: "Table updated successfully",
	}, nil
}

// DeleteTable deletes a table
func (h *TableHandler) DeleteTable(ctx context.Context, req *pb.DeleteTableRequest) (*pb.DeleteTableResponse, error) {
	if req.TableId == "" {
		return nil, status.Error(codes.InvalidArgument, "table_id is required")
	}

	if err := h.tableUseCase.DeleteTable(ctx, req.TableId); err != nil {
		return nil, handleError(err)
	}

	return &pb.DeleteTableResponse{
		Success: true,
		Message: "Table deleted successfully",
	}, nil
}

// ListTables retrieves tables with filters and pagination
func (h *TableHandler) ListTables(ctx context.Context, req *pb.ListTablesRequest) (*pb.ListTablesResponse, error) {
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
	tableStatus := protoStatusToDomain(req.Status)

	tables, total, err := h.tableUseCase.ListTables(ctx, page, pageSize, tableStatus, req.Location)
	if err != nil {
		return nil, handleError(err)
	}

	// Convert tables to proto
	protoTables := make([]*pb.Table, len(tables))
	for i, table := range tables {
		protoTables[i] = domainTableToProto(table)
	}

	return &pb.ListTablesResponse{
		Tables:   protoTables,
		Total:    int32(total),
		Page:     int32(page),
		PageSize: int32(pageSize),
		Success:  true,
		Message:  "Tables retrieved successfully",
	}, nil
}

// UpdateTableStatus updates the status of a table
func (h *TableHandler) UpdateTableStatus(ctx context.Context, req *pb.UpdateTableStatusRequest) (*pb.UpdateTableStatusResponse, error) {
	if req.TableId == "" {
		return nil, status.Error(codes.InvalidArgument, "table_id is required")
	}

	// Convert proto status to domain
	tableStatus := protoStatusToDomain(req.Status)
	if tableStatus == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid status")
	}

	table, err := h.tableUseCase.UpdateTableStatus(ctx, req.TableId, tableStatus, req.OrderId)
	if err != nil {
		return nil, handleError(err)
	}

	return &pb.UpdateTableStatusResponse{
		Table:   domainTableToProto(table),
		Success: true,
		Message: "Table status updated successfully",
	}, nil
}

// GetAvailableTables retrieves available tables
func (h *TableHandler) GetAvailableTables(ctx context.Context, req *pb.GetAvailableTablesRequest) (*pb.GetAvailableTablesResponse, error) {
	minCapacity := int(req.MinCapacity)
	if minCapacity < 0 {
		minCapacity = 0
	}

	tables, err := h.tableUseCase.GetAvailableTables(ctx, minCapacity, req.Location)
	if err != nil {
		return nil, handleError(err)
	}

	// Convert tables to proto
	protoTables := make([]*pb.Table, len(tables))
	for i, table := range tables {
		protoTables[i] = domainTableToProto(table)
	}

	return &pb.GetAvailableTablesResponse{
		Tables:  protoTables,
		Success: true,
		Message: "Available tables retrieved successfully",
	}, nil
}

// CreateReservation creates a new table reservation.
func (h *TableHandler) CreateReservation(ctx context.Context, req *pb.CreateReservationRequest) (*pb.CreateReservationResponse, error) {
	if req.TableId == "" {
		return nil, status.Error(codes.InvalidArgument, "table_id is required")
	}
	if req.CustomerName == "" {
		return nil, status.Error(codes.InvalidArgument, "customer_name is required")
	}
	if req.StartTime == nil || req.EndTime == nil {
		return nil, status.Error(codes.InvalidArgument, "start_time and end_time are required")
	}

	items := make([]domain.ReservationItem, 0, len(req.Items))
	for _, item := range req.Items {
		items = append(items, domain.ReservationItem{
			MenuItemID: item.MenuItemId,
			Quantity:   int(item.Quantity),
			Note:       item.Note,
		})
	}

	reservation, err := h.reservationUseCase.CreateReservation(
		ctx,
		req.TableId,
		req.CustomerName,
		req.CustomerPhone,
		req.StartTime.AsTime(),
		req.EndTime.AsTime(),
		req.Notes,
		items,
	)
	if err != nil {
		return nil, handleError(err)
	}

	return &pb.CreateReservationResponse{
		Reservation: domainReservationToProto(reservation),
		Success:     true,
		Message:     "Reservation created successfully",
	}, nil
}

// GetReservation retrieves a reservation by ID.
func (h *TableHandler) GetReservation(ctx context.Context, req *pb.GetReservationRequest) (*pb.GetReservationResponse, error) {
	if req.ReservationId == "" {
		return nil, status.Error(codes.InvalidArgument, "reservation_id is required")
	}

	reservation, err := h.reservationUseCase.GetReservation(ctx, req.ReservationId)
	if err != nil {
		return nil, handleError(err)
	}

	return &pb.GetReservationResponse{
		Reservation: domainReservationToProto(reservation),
		Success:     true,
		Message:     "Reservation retrieved successfully",
	}, nil
}

// ListReservations retrieves reservations with filters and pagination.
func (h *TableHandler) ListReservations(ctx context.Context, req *pb.ListReservationsRequest) (*pb.ListReservationsResponse, error) {
	page := int(req.Page)
	if page <= 0 {
		page = 1
	}
	pageSize := int(req.PageSize)
	if pageSize <= 0 {
		pageSize = 10
	}

	filters := repository.ReservationListFilters{
		Page:     page,
		PageSize: pageSize,
		TableID:  req.TableId,
		Status:   protoReservationStatusToDomain(req.Status),
	}
	if req.FromTime != nil {
		filters.FromTime = req.FromTime.AsTime()
	}
	if req.ToTime != nil {
		filters.ToTime = req.ToTime.AsTime()
	}

	reservations, total, err := h.reservationUseCase.ListReservations(ctx, filters)
	if err != nil {
		return nil, handleError(err)
	}

	protoReservations := make([]*pb.Reservation, len(reservations))
	for i, reservation := range reservations {
		protoReservations[i] = domainReservationToProto(reservation)
	}

	return &pb.ListReservationsResponse{
		Reservations: protoReservations,
		Total:        int32(total),
		Page:         int32(page),
		PageSize:     int32(pageSize),
		Success:      true,
		Message:      "Reservations retrieved successfully",
	}, nil
}

// CancelReservation cancels a reservation by ID.
func (h *TableHandler) CancelReservation(ctx context.Context, req *pb.CancelReservationRequest) (*pb.CancelReservationResponse, error) {
	if req.ReservationId == "" {
		return nil, status.Error(codes.InvalidArgument, "reservation_id is required")
	}

	reservation, err := h.reservationUseCase.CancelReservation(ctx, req.ReservationId)
	if err != nil {
		return nil, handleError(err)
	}

	return &pb.CancelReservationResponse{
		Reservation: domainReservationToProto(reservation),
		Success:     true,
		Message:     "Reservation cancelled successfully",
	}, nil
}

// Helper functions

// domainTableToProto converts domain.Table to pb.Table
func domainTableToProto(table *domain.Table) *pb.Table {
	if table == nil {
		return nil
	}

	return &pb.Table{
		TableId:        table.ID,
		TableNumber:    table.TableNumber,
		Capacity:       int32(table.Capacity),
		Status:         domainStatusToProto(table.Status),
		Location:       table.Location,
		CurrentOrderId: table.CurrentOrderID,
		CreatedAt:      timestamppb.New(table.CreatedAt),
		UpdatedAt:      timestamppb.New(table.UpdatedAt),
	}
}

func domainReservationToProto(reservation *domain.Reservation) *pb.Reservation {
	if reservation == nil {
		return nil
	}

	items := make([]*pb.ReservationItem, len(reservation.Items))
	for i, item := range reservation.Items {
		items[i] = &pb.ReservationItem{
			MenuItemId: item.MenuItemID,
			Quantity:   int32(item.Quantity),
			Note:       item.Note,
		}
	}

	return &pb.Reservation{
		ReservationId: reservation.ID,
		TableId:       reservation.TableID,
		CustomerName:  reservation.CustomerName,
		CustomerPhone: reservation.CustomerPhone,
		Notes:         reservation.Notes,
		Status:        domainReservationStatusToProto(reservation.Status),
		StartTime:     timestamppb.New(reservation.StartTime),
		EndTime:       timestamppb.New(reservation.EndTime),
		Items:         items,
		CreatedAt:     timestamppb.New(reservation.CreatedAt),
		UpdatedAt:     timestamppb.New(reservation.UpdatedAt),
	}
}

// domainStatusToProto converts domain status to proto status
func domainStatusToProto(status domain.TableStatus) pb.TableStatus {
	switch status {
	case domain.StatusAvailable:
		return pb.TableStatus_STATUS_AVAILABLE
	case domain.StatusOccupied:
		return pb.TableStatus_STATUS_OCCUPIED
	case domain.StatusReserved:
		return pb.TableStatus_STATUS_RESERVED
	case domain.StatusCleaning:
		return pb.TableStatus_STATUS_CLEANING
	case domain.StatusOutOfService:
		return pb.TableStatus_STATUS_OUT_OF_SERVICE
	default:
		return pb.TableStatus_STATUS_UNKNOWN
	}
}

func domainReservationStatusToProto(status domain.ReservationStatus) pb.ReservationStatus {
	switch status {
	case domain.ReservationStatusReserved:
		return pb.ReservationStatus_RESERVATION_STATUS_RESERVED
	case domain.ReservationStatusCancelled:
		return pb.ReservationStatus_RESERVATION_STATUS_CANCELLED
	case domain.ReservationStatusCompleted:
		return pb.ReservationStatus_RESERVATION_STATUS_COMPLETED
	default:
		return pb.ReservationStatus_RESERVATION_STATUS_UNKNOWN
	}
}

// protoStatusToDomain converts proto status to domain status
func protoStatusToDomain(status pb.TableStatus) domain.TableStatus {
	switch status {
	case pb.TableStatus_STATUS_AVAILABLE:
		return domain.StatusAvailable
	case pb.TableStatus_STATUS_OCCUPIED:
		return domain.StatusOccupied
	case pb.TableStatus_STATUS_RESERVED:
		return domain.StatusReserved
	case pb.TableStatus_STATUS_CLEANING:
		return domain.StatusCleaning
	case pb.TableStatus_STATUS_OUT_OF_SERVICE:
		return domain.StatusOutOfService
	default:
		return ""
	}
}

func protoReservationStatusToDomain(status pb.ReservationStatus) domain.ReservationStatus {
	switch status {
	case pb.ReservationStatus_RESERVATION_STATUS_RESERVED:
		return domain.ReservationStatusReserved
	case pb.ReservationStatus_RESERVATION_STATUS_CANCELLED:
		return domain.ReservationStatusCancelled
	case pb.ReservationStatus_RESERVATION_STATUS_COMPLETED:
		return domain.ReservationStatusCompleted
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
	case errors.Is(err, domain.ErrTableNotFound):
		return status.Error(codes.NotFound, "table not found")
	case errors.Is(err, domain.ErrTableAlreadyExists):
		return status.Error(codes.AlreadyExists, "table already exists")
	case errors.Is(err, domain.ErrTableNumberAlreadyExists):
		return status.Error(codes.AlreadyExists, "table number already exists")
	case errors.Is(err, domain.ErrTableNotAvailable):
		return status.Error(codes.FailedPrecondition, "table is not available")
	case errors.Is(err, domain.ErrTableOutOfService):
		return status.Error(codes.FailedPrecondition, "table is out of service")
	case errors.Is(err, domain.ErrTableAlreadyOccupied):
		return status.Error(codes.FailedPrecondition, "table is already occupied")
	case errors.Is(err, domain.ErrReservationNotFound):
		return status.Error(codes.NotFound, "reservation not found")
	case errors.Is(err, domain.ErrReservationConflict):
		return status.Error(codes.FailedPrecondition, "reservation time conflict")
	case errors.Is(err, domain.ErrReservationAlreadyCancelled):
		return status.Error(codes.FailedPrecondition, "reservation already cancelled")
	case errors.Is(err, domain.ErrReservationAlreadyCompleted):
		return status.Error(codes.FailedPrecondition, "reservation already completed")
	case errors.Is(err, domain.ErrInvalidTableNumber),
		errors.Is(err, domain.ErrInvalidCapacity),
		errors.Is(err, domain.ErrCapacityTooLarge),
		errors.Is(err, domain.ErrInvalidLocation),
		errors.Is(err, domain.ErrInvalidStatus),
		errors.Is(err, domain.ErrInvalidOrderID),
		errors.Is(err, domain.ErrInvalidTableID),
		errors.Is(err, domain.ErrInvalidReservationTime),
		errors.Is(err, domain.ErrInvalidReservationStatus),
		errors.Is(err, domain.ErrInvalidReservationItem),
		errors.Is(err, domain.ErrInvalidReservationCustomer):
		return status.Error(codes.InvalidArgument, err.Error())
	default:
		// Internal server error for unknown errors
		return status.Error(codes.Internal, "internal server error")
	}
}
