// Package grpc contains gRPC handler implementations for the report service.
package grpc

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pb "restaurant-management/proto/report"
	"restaurant-management/services/report-service/internal/domain"
	"restaurant-management/services/report-service/internal/usecase"
	"time"
)

// ReportHandler implements the ReportService gRPC handler.
type ReportHandler struct {
	pb.UnimplementedReportServiceServer
	useCase *usecase.ReportUseCase
}

// NewReportHandler creates a new report gRPC handler.
func NewReportHandler(useCase *usecase.ReportUseCase) *ReportHandler {
	return &ReportHandler{
		useCase: useCase,
	}
}

// GetSalesReport implements GetSalesReport RPC.
func (h *ReportHandler) GetSalesReport(ctx context.Context, req *pb.GetSalesReportRequest) (*pb.GetSalesReportResponse, error) {
	var fromDate, toDate time.Time

	if req.FromDate != nil {
		fromDate = req.FromDate.AsTime()
	}
	if req.ToDate != nil {
		toDate = req.ToDate.AsTime()
	}

	// Default to current day if not specified
	if fromDate.IsZero() || toDate.IsZero() {
		now := time.Now()
		fromDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		toDate = fromDate.Add(24 * time.Hour)
	}

	dateRange := &domain.DateRange{
		FromDate: fromDate,
		ToDate:   toDate,
	}

	report, err := h.useCase.GetSalesReport(ctx, dateRange)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	salesProto := convertSalesReportToProto(report)
	return &pb.GetSalesReportResponse{
		Report:  salesProto,
		Success: true,
		Message: "Sales report generated successfully",
	}, nil
}

// GetInventoryReport implements GetInventoryReport RPC.
func (h *ReportHandler) GetInventoryReport(ctx context.Context, req *pb.GetInventoryReportRequest) (*pb.GetInventoryReportResponse, error) {
	report, err := h.useCase.GetInventoryReport(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	invProto := convertInventoryReportToProto(report)
	return &pb.GetInventoryReportResponse{
		Report:  invProto,
		Success: true,
		Message: "Inventory report generated successfully",
	}, nil
}

// GetOrderReport implements GetOrderReport RPC.
func (h *ReportHandler) GetOrderReport(ctx context.Context, req *pb.GetOrderReportRequest) (*pb.GetOrderReportResponse, error) {
	var fromDate, toDate time.Time

	if req.FromDate != nil {
		fromDate = req.FromDate.AsTime()
	}
	if req.ToDate != nil {
		toDate = req.ToDate.AsTime()
	}

	if fromDate.IsZero() || toDate.IsZero() {
		now := time.Now()
		fromDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		toDate = fromDate.Add(24 * time.Hour)
	}

	dateRange := &domain.DateRange{
		FromDate: fromDate,
		ToDate:   toDate,
	}

	report, err := h.useCase.GetOrderReport(ctx, dateRange)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	orderProto := convertOrderReportToProto(report)
	return &pb.GetOrderReportResponse{
		Report:  orderProto,
		Success: true,
		Message: "Order report generated successfully",
	}, nil
}

// GetStaffPerformanceReport implements GetStaffPerformanceReport RPC.
func (h *ReportHandler) GetStaffPerformanceReport(ctx context.Context, req *pb.GetStaffPerformanceReportRequest) (*pb.GetStaffPerformanceReportResponse, error) {
	var fromDate, toDate time.Time

	if req.FromDate != nil {
		fromDate = req.FromDate.AsTime()
	}
	if req.ToDate != nil {
		toDate = req.ToDate.AsTime()
	}

	if fromDate.IsZero() || toDate.IsZero() {
		now := time.Now()
		fromDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		toDate = fromDate.Add(24 * time.Hour)
	}

	dateRange := &domain.DateRange{
		FromDate: fromDate,
		ToDate:   toDate,
	}

	performances, err := h.useCase.GetStaffPerformanceReport(ctx, req.StaffId, dateRange)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var protoPerformances []*pb.StaffPerformance
	for _, perf := range performances {
		protoPerformances = append(protoPerformances, &pb.StaffPerformance{
			StaffId:             perf.StaffID,
			StaffName:           perf.StaffName,
			OrdersHandled:       perf.OrdersHandled,
			TotalSales:          perf.TotalSales,
			AverageOrderValue:   perf.AverageOrderValue,
			CustomerComplaints:  perf.CustomerComplaints,
		})
	}

	return &pb.GetStaffPerformanceReportResponse{
		Performances: protoPerformances,
		Success:      true,
		Message:      "Staff performance report generated successfully",
	}, nil
}

// GetPopularItemsReport implements GetPopularItemsReport RPC.
func (h *ReportHandler) GetPopularItemsReport(ctx context.Context, req *pb.GetPopularItemsReportRequest) (*pb.GetPopularItemsReportResponse, error) {
	var fromDate, toDate time.Time

	if req.FromDate != nil {
		fromDate = req.FromDate.AsTime()
	}
	if req.ToDate != nil {
		toDate = req.ToDate.AsTime()
	}

	if fromDate.IsZero() || toDate.IsZero() {
		now := time.Now()
		fromDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		toDate = fromDate.Add(24 * time.Hour)
	}

	dateRange := &domain.DateRange{
		FromDate: fromDate,
		ToDate:   toDate,
	}

	items, err := h.useCase.GetPopularItemsReport(ctx, req.TopN, dateRange)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var protoItems []*pb.PopularItem
	for _, item := range items {
		protoItems = append(protoItems, &pb.PopularItem{
			ItemId:        item.ItemID,
			ItemName:      item.ItemName,
			QuantitySold:  item.QuantitySold,
			Revenue:       item.Revenue,
			Category:      item.Category,
		})
	}

	return &pb.GetPopularItemsReportResponse{
		Items:   protoItems,
		Success: true,
		Message: "Popular items report generated successfully",
	}, nil
}

// GetRevenueAnalytics implements GetRevenueAnalytics RPC.
func (h *ReportHandler) GetRevenueAnalytics(ctx context.Context, req *pb.GetRevenueAnalyticsRequest) (*pb.GetRevenueAnalyticsResponse, error) {
	var fromDate, toDate time.Time

	if req.FromDate != nil {
		fromDate = req.FromDate.AsTime()
	}
	if req.ToDate != nil {
		toDate = req.ToDate.AsTime()
	}

	if fromDate.IsZero() || toDate.IsZero() {
		now := time.Now()
		fromDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		toDate = fromDate.Add(24 * time.Hour)
	}

	dateRange := &domain.DateRange{
		FromDate: fromDate,
		ToDate:   toDate,
	}

	analytics, err := h.useCase.GetRevenueAnalytics(ctx, dateRange)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	analyticsProto := convertRevenueAnalyticsToProto(analytics)
	return &pb.GetRevenueAnalyticsResponse{
		Analytics: analyticsProto,
		Success:   true,
		Message:   "Revenue analytics generated successfully",
	}, nil
}

// ExportReport implements ExportReport RPC.
func (h *ReportHandler) ExportReport(ctx context.Context, req *pb.ExportReportRequest) (*pb.ExportReportResponse, error) {
	var fromDate, toDate time.Time

	if req.FromDate != nil {
		fromDate = req.FromDate.AsTime()
	}
	if req.ToDate != nil {
		toDate = req.ToDate.AsTime()
	}

	if fromDate.IsZero() || toDate.IsZero() {
		now := time.Now()
		fromDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		toDate = fromDate.Add(24 * time.Hour)
	}

	dateRange := &domain.DateRange{
		FromDate: fromDate,
		ToDate:   toDate,
	}

	format := domain.ExportFormat(req.Format)
	fileData, fileName, contentType, err := h.useCase.ExportReport(ctx, req.ReportType, format, dateRange)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.ExportReportResponse{
		FileData:    fileData,
		FileName:    fileName,
		ContentType: contentType,
		Success:     true,
		Message:     "Report exported successfully",
	}, nil
}

// Helper conversion functions

func convertSalesReportToProto(report *domain.SalesReport) *pb.SalesReport {
	var dailyProto []*pb.DailySales
	for _, daily := range report.DailyBreakdown {
		dailyProto = append(dailyProto, &pb.DailySales{
			Date:   daily.Date,
			Sales:  daily.Sales,
			Orders: daily.Orders,
		})
	}

	var paymentProto []*pb.PaymentMethodSales
	for _, payment := range report.PaymentBreakdown {
		paymentProto = append(paymentProto, &pb.PaymentMethodSales{
			Method: payment.Method,
			Amount: payment.Amount,
			Count:  payment.Count,
		})
	}

	return &pb.SalesReport{
		TotalSales:        report.TotalSales,
		TotalTax:          report.TotalTax,
		TotalDiscount:     report.TotalDiscount,
		NetRevenue:        report.NetRevenue,
		TotalOrders:       report.TotalOrders,
		AverageOrderValue: report.AverageOrderValue,
		DailyBreakdown:    dailyProto,
		PaymentBreakdown:  paymentProto,
	}
}

func convertInventoryReportToProto(report *domain.InventoryReport) *pb.InventoryReport {
	var itemsProto []*pb.IngredientStock
	for _, item := range report.Items {
		itemsProto = append(itemsProto, &pb.IngredientStock{
			IngredientId:  item.IngredientID,
			Name:          item.Name,
			CurrentStock:  item.CurrentStock,
			MinimumStock:  item.MinimumStock,
			UnitCost:      item.UnitCost,
			TotalValue:    item.TotalValue,
			IsLowStock:    item.IsLowStock,
		})
	}

	return &pb.InventoryReport{
		TotalInventoryValue: report.TotalInventoryValue,
		LowStockItemsCount:  report.LowStockItemsCount,
		Items:               itemsProto,
		TotalWasteValue:     report.TotalWasteValue,
	}
}

func convertOrderReportToProto(report *domain.OrderReport) *pb.OrderReport {
	var hourlyProto []*pb.HourlyOrders
	for _, hourly := range report.HourlyBreakdown {
		hourlyProto = append(hourlyProto, &pb.HourlyOrders{
			Hour:       hourly.Hour,
			OrderCount: hourly.OrderCount,
		})
	}

	return &pb.OrderReport{
		TotalOrders:        report.TotalOrders,
		CompletedOrders:    report.CompletedOrders,
		CancelledOrders:    report.CancelledOrders,
		CompletionRate:     report.CompletionRate,
		AveragePreparationTime: report.AvgPreparationTime,
		HourlyBreakdown:    hourlyProto,
	}
}

func convertRevenueAnalyticsToProto(analytics *domain.RevenueAnalytics) *pb.RevenueAnalytics {
	var trendProto []*pb.DailySales
	for _, trend := range analytics.TrendData {
		trendProto = append(trendProto, &pb.DailySales{
			Date:   trend.Date,
			Sales:  trend.Sales,
			Orders: trend.Orders,
		})
	}

	return &pb.RevenueAnalytics{
		TotalRevenue: analytics.TotalRevenue,
		GrowthRate:   analytics.GrowthRate,
		TrendData:    trendProto,
		PeakDay:      analytics.PeakDay,
		PeakHour:     analytics.PeakHour,
	}
}
