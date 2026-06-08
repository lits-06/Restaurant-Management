package grpc

import (
	"context"
	"errors"
	"time"

	"restaurant-management/proto/payment"
	"restaurant-management/services/payment-service/internal/domain"
	"restaurant-management/services/payment-service/internal/usecase"
	"restaurant-management/shared/pkg/logger"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// PaymentHandler handles gRPC requests for payment service
type PaymentHandler struct {
	payment.UnimplementedPaymentServiceServer
	paymentUseCase *usecase.PaymentUseCase
	logger         logger.Logger
}

// NewPaymentHandler creates a new payment handler
func NewPaymentHandler(paymentUseCase *usecase.PaymentUseCase, logger logger.Logger) *PaymentHandler {
	return &PaymentHandler{
		paymentUseCase: paymentUseCase,
		logger:         logger,
	}
}

// CreatePayment creates a new payment
func (h *PaymentHandler) CreatePayment(ctx context.Context, req *payment.CreatePaymentRequest) (*payment.CreatePaymentResponse, error) {
	h.logger.Info("CreatePayment request", "order_id", req.OrderId, "amount", req.Amount)

	if req.OrderId == "" {
		return &payment.CreatePaymentResponse{
			Success: false,
			Message: "order_id is required",
		}, nil
	}

	if req.Amount <= 0 {
		return &payment.CreatePaymentResponse{
			Success: false,
			Message: "amount must be greater than 0",
		}, nil
	}

	if req.Method == payment.PaymentMethod_METHOD_UNKNOWN {
		return &payment.CreatePaymentResponse{
			Success: false,
			Message: "payment method is required",
		}, nil
	}

	// Convert proto method to domain method
	domainMethod := convertProtoMethodToDomain(req.Method)

	// Create payment
	pmt, err := h.paymentUseCase.CreatePayment(
		ctx,
		req.OrderId,
		req.Amount,
		req.Tip,
		domainMethod,
		req.CustomerName,
		req.Notes,
	)

	if err != nil {
		h.logger.Error("Failed to create payment", "error", err)
		return &payment.CreatePaymentResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &payment.CreatePaymentResponse{
		Payment: convertPaymentToProto(pmt),
		Success: true,
		Message: "payment created successfully",
	}, nil
}

// GetPayment retrieves a payment by ID
func (h *PaymentHandler) GetPayment(ctx context.Context, req *payment.GetPaymentRequest) (*payment.GetPaymentResponse, error) {
	h.logger.Info("GetPayment request", "payment_id", req.PaymentId)

	if req.PaymentId == "" {
		return &payment.GetPaymentResponse{
			Success: false,
			Message: "payment_id is required",
		}, nil
	}

	pmt, err := h.paymentUseCase.GetPayment(ctx, req.PaymentId)
	if err != nil {
		if errors.Is(err, domain.ErrPaymentNotFound) {
			return &payment.GetPaymentResponse{
				Success: false,
				Message: "payment not found",
			}, nil
		}
		h.logger.Error("Failed to get payment", "error", err)
		return &payment.GetPaymentResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &payment.GetPaymentResponse{
		Payment: convertPaymentToProto(pmt),
		Success: true,
		Message: "payment retrieved successfully",
	}, nil
}

// ProcessPayment processes a payment
func (h *PaymentHandler) ProcessPayment(ctx context.Context, req *payment.ProcessPaymentRequest) (*payment.ProcessPaymentResponse, error) {
	h.logger.Info("ProcessPayment request", "payment_id", req.PaymentId)

	if req.PaymentId == "" {
		return &payment.ProcessPaymentResponse{
			Success: false,
			Message: "payment_id is required",
		}, nil
	}

	if req.TransactionId == "" {
		return &payment.ProcessPaymentResponse{
			Success: false,
			Message: "transaction_id is required",
		}, nil
	}

	pmt, err := h.paymentUseCase.ProcessPayment(ctx, req.PaymentId, req.TransactionId)
	if err != nil {
		if errors.Is(err, domain.ErrPaymentNotFound) {
			return &payment.ProcessPaymentResponse{
				Success: false,
				Message: "payment not found",
			}, nil
		}
		h.logger.Error("Failed to process payment", "error", err)
		return &payment.ProcessPaymentResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &payment.ProcessPaymentResponse{
		Payment: convertPaymentToProto(pmt),
		Success: true,
		Message: "payment processed successfully",
	}, nil
}

// RefundPayment processes a refund
func (h *PaymentHandler) RefundPayment(ctx context.Context, req *payment.RefundPaymentRequest) (*payment.RefundPaymentResponse, error) {
	h.logger.Info("RefundPayment request", "payment_id", req.PaymentId, "amount", req.Amount)

	if req.PaymentId == "" {
		return &payment.RefundPaymentResponse{
			Success: false,
			Message: "payment_id is required",
		}, nil
	}

	if req.Amount <= 0 {
		return &payment.RefundPaymentResponse{
			Success: false,
			Message: "refund amount must be greater than 0",
		}, nil
	}

	pmt, err := h.paymentUseCase.RefundPayment(ctx, req.PaymentId, req.Amount, req.Reason)
	if err != nil {
		if errors.Is(err, domain.ErrPaymentNotFound) {
			return &payment.RefundPaymentResponse{
				Success: false,
				Message: "payment not found",
			}, nil
		}
		h.logger.Error("Failed to refund payment", "error", err)
		return &payment.RefundPaymentResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &payment.RefundPaymentResponse{
		Payment: convertPaymentToProto(pmt),
		Success: true,
		Message: "payment refunded successfully",
	}, nil
}

// ListPayments retrieves payments with pagination and filters
func (h *PaymentHandler) ListPayments(ctx context.Context, req *payment.ListPaymentsRequest) (*payment.ListPaymentsResponse, error) {
	h.logger.Info("ListPayments request", "page", req.Page, "page_size", req.PageSize)

	page := int(req.Page)
	pageSize := int(req.PageSize)

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	// Convert proto enums to domain enums
	domainStatus := convertProtoStatusToDomain(req.Status)
	domainMethod := convertProtoMethodToDomain(req.Method)

	// Convert timestamps
	var fromTime, toTime *time.Time
	if req.FromDate != nil {
		t := req.FromDate.AsTime()
		fromTime = &t
	}
	if req.ToDate != nil {
		t := req.ToDate.AsTime()
		toTime = &t
	}

	payments, total, err := h.paymentUseCase.ListPayments(ctx, page, pageSize, domainStatus, domainMethod, fromTime, toTime)
	if err != nil {
		h.logger.Error("Failed to list payments", "error", err)
		return &payment.ListPaymentsResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	protoPayments := make([]*payment.Payment, len(payments))
	for i, pmt := range payments {
		protoPayments[i] = convertPaymentToProto(pmt)
	}

	return &payment.ListPaymentsResponse{
		Payments: protoPayments,
		Total:    int32(total),
		Page:     int32(page),
		PageSize: int32(pageSize),
		Success:  true,
		Message:  "payments retrieved successfully",
	}, nil
}

// GetPaymentsByOrder retrieves payments for an order
func (h *PaymentHandler) GetPaymentsByOrder(ctx context.Context, req *payment.GetPaymentsByOrderRequest) (*payment.GetPaymentsByOrderResponse, error) {
	h.logger.Info("GetPaymentsByOrder request", "order_id", req.OrderId)

	if req.OrderId == "" {
		return &payment.GetPaymentsByOrderResponse{
			Success: false,
			Message: "order_id is required",
		}, nil
	}

	payments, err := h.paymentUseCase.GetPaymentsByOrder(ctx, req.OrderId)
	if err != nil {
		h.logger.Error("Failed to get payments by order", "error", err)
		return &payment.GetPaymentsByOrderResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	protoPayments := make([]*payment.Payment, len(payments))
	for i, pmt := range payments {
		protoPayments[i] = convertPaymentToProto(pmt)
	}

	return &payment.GetPaymentsByOrderResponse{
		Payments: protoPayments,
		Success:  true,
		Message:  "payments retrieved successfully",
	}, nil
}

// Helper functions

func convertPaymentToProto(pmt *domain.Payment) *payment.Payment {
	if pmt == nil {
		return nil
	}

	return &payment.Payment{
		PaymentId:     pmt.PaymentID,
		OrderId:       pmt.OrderID,
		Amount:        pmt.Amount,
		Tip:           pmt.Tip,
		Total:         pmt.Total,
		Method:        convertDomainMethodToProto(pmt.Method),
		Status:        convertDomainStatusToProto(pmt.Status),
		TransactionId: pmt.TransactionID,
		CustomerName:  pmt.CustomerName,
		Notes:         pmt.Notes,
		CreatedAt:     timestamppb.New(pmt.CreatedAt),
		UpdatedAt:     timestamppb.New(pmt.UpdatedAt),
	}
}

func convertDomainMethodToProto(method domain.PaymentMethod) payment.PaymentMethod {
	switch method {
	case domain.MethodCash:
		return payment.PaymentMethod_METHOD_CASH
	case domain.MethodCreditCard:
		return payment.PaymentMethod_METHOD_CREDIT_CARD
	case domain.MethodDebitCard:
		return payment.PaymentMethod_METHOD_DEBIT_CARD
	case domain.MethodMobileWallet:
		return payment.PaymentMethod_METHOD_MOBILE_WALLET
	case domain.MethodBankTransfer:
		return payment.PaymentMethod_METHOD_BANK_TRANSFER
	default:
		return payment.PaymentMethod_METHOD_UNKNOWN
	}
}

func convertProtoMethodToDomain(method payment.PaymentMethod) domain.PaymentMethod {
	switch method {
	case payment.PaymentMethod_METHOD_CASH:
		return domain.MethodCash
	case payment.PaymentMethod_METHOD_CREDIT_CARD:
		return domain.MethodCreditCard
	case payment.PaymentMethod_METHOD_DEBIT_CARD:
		return domain.MethodDebitCard
	case payment.PaymentMethod_METHOD_MOBILE_WALLET:
		return domain.MethodMobileWallet
	case payment.PaymentMethod_METHOD_BANK_TRANSFER:
		return domain.MethodBankTransfer
	default:
		return domain.MethodUnknown
	}
}

func convertDomainStatusToProto(status domain.PaymentStatus) payment.PaymentStatus {
	switch status {
	case domain.StatusPending:
		return payment.PaymentStatus_STATUS_PENDING
	case domain.StatusProcessing:
		return payment.PaymentStatus_STATUS_PROCESSING
	case domain.StatusCompleted:
		return payment.PaymentStatus_STATUS_COMPLETED
	case domain.StatusFailed:
		return payment.PaymentStatus_STATUS_FAILED
	case domain.StatusRefunded:
		return payment.PaymentStatus_STATUS_REFUNDED
	case domain.StatusPartiallyRefunded:
		return payment.PaymentStatus_STATUS_PARTIALLY_REFUNDED
	default:
		return payment.PaymentStatus_STATUS_UNKNOWN
	}
}

func convertProtoStatusToDomain(status payment.PaymentStatus) domain.PaymentStatus {
	switch status {
	case payment.PaymentStatus_STATUS_PENDING:
		return domain.StatusPending
	case payment.PaymentStatus_STATUS_PROCESSING:
		return domain.StatusProcessing
	case payment.PaymentStatus_STATUS_COMPLETED:
		return domain.StatusCompleted
	case payment.PaymentStatus_STATUS_FAILED:
		return domain.StatusFailed
	case payment.PaymentStatus_STATUS_REFUNDED:
		return domain.StatusRefunded
	case payment.PaymentStatus_STATUS_PARTIALLY_REFUNDED:
		return domain.StatusPartiallyRefunded
	default:
		return domain.StatusUnknown
	}
}
