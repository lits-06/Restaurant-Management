// Package usecase contains usecase implementations for report service.
package usecase

import (
	"context"
	"fmt"
	"math"
	"restaurant-management/services/report-service/internal/domain"
	"restaurant-management/services/report-service/internal/repository"
	"time"
)

// OrderServiceClient defines methods to call Order service
type OrderServiceClient interface {
	GetOrdersByDateRange(ctx context.Context, fromDate, toDate time.Time) (interface{}, error)
	GetOrderByID(ctx context.Context, orderID string) (interface{}, error)
}

// PaymentServiceClient defines methods to call Payment service
type PaymentServiceClient interface {
	GetPaymentsByDateRange(ctx context.Context, fromDate, toDate time.Time) (interface{}, error)
	GetPaymentByID(ctx context.Context, paymentID string) (interface{}, error)
}

// InventoryServiceClient defines methods to call Inventory service
type InventoryServiceClient interface {
	GetInventorySnapshot(ctx context.Context) (interface{}, error)
	GetLowStockItems(ctx context.Context) (interface{}, error)
}

// UserServiceClient defines methods to call User service
type UserServiceClient interface {
	GetUserByID(ctx context.Context, userID string) (interface{}, error)
	ListUsers(ctx context.Context) (interface{}, error)
}

// MenuServiceClient defines methods to call Menu service
type MenuServiceClient interface {
	GetMenuItemByID(ctx context.Context, itemID string) (interface{}, error)
	ListMenuItems(ctx context.Context) (interface{}, error)
}

// ReportUseCase implements business logic for report operations.
type ReportUseCase struct {
	repo                  repository.ReportRepository
	orderClient           OrderServiceClient
	paymentClient         PaymentServiceClient
	inventoryClient       InventoryServiceClient
	userClient            UserServiceClient
	menuClient            MenuServiceClient
}

// NewReportUseCase creates a new report use case.
func NewReportUseCase(
	repo repository.ReportRepository,
	orderClient OrderServiceClient,
	paymentClient PaymentServiceClient,
	inventoryClient InventoryServiceClient,
	userClient UserServiceClient,
	menuClient MenuServiceClient,
) *ReportUseCase {
	return &ReportUseCase{
		repo:            repo,
		orderClient:     orderClient,
		paymentClient:   paymentClient,
		inventoryClient: inventoryClient,
		userClient:      userClient,
		menuClient:      menuClient,
	}
}

// GetSalesReport generates sales report for a date range
func (uc *ReportUseCase) GetSalesReport(ctx context.Context, dateRange *domain.DateRange) (*domain.SalesReport, error) {
	if err := dateRange.ValidatePeriod(); err != nil {
		return nil, err
	}

	// Fetch payments from Payment service
	payments, err := uc.paymentClient.GetPaymentsByDateRange(ctx, dateRange.FromDate, dateRange.ToDate)
	if err != nil {
		// Continue with nil if service unavailable - for demo purposes
		payments = nil
	}

	report := &domain.SalesReport{
		DailyBreakdown:   []domain.DailySales{},
		PaymentBreakdown: []domain.PaymentMethodSales{},
	}

	// Mock calculation (in production, aggregate real payment data)
	if payments != nil {
		_ = payments
	}

	report.TotalSales = 45000000.0
	report.TotalTax = 4500000.0
	report.TotalDiscount = 2250000.0
	report.NetRevenue = 47250000.0
	report.TotalOrders = 150
	report.AverageOrderValue = 300000.0

	// Daily breakdown (7 days)
	days := []string{"2026-04-03", "2026-04-04", "2026-04-05", "2026-04-06", "2026-04-07", "2026-04-08", "2026-04-09"}
	salesByDay := []float64{6000000, 6500000, 7000000, 6200000, 6800000, 7000000, 5500000}
	ordersByDay := []int32{20, 22, 25, 21, 23, 24, 15}

	for i := range days {
		report.DailyBreakdown = append(report.DailyBreakdown, domain.DailySales{
			Date:   days[i],
			Sales:  salesByDay[i],
			Orders: ordersByDay[i],
		})
	}

	// Payment method breakdown
	report.PaymentBreakdown = []domain.PaymentMethodSales{
		{Method: "CASH", Amount: 20000000, Count: 70},
		{Method: "CREDIT_CARD", Amount: 15000000, Count: 50},
		{Method: "DEBIT_CARD", Amount: 7000000, Count: 20},
		{Method: "MOBILE_WALLET", Amount: 3000000, Count: 10},
	}

	_ = uc.repo.SaveSalesData(ctx, report)
	return report, nil
}

// GetInventoryReport generates inventory snapshot report
func (uc *ReportUseCase) GetInventoryReport(ctx context.Context) (*domain.InventoryReport, error) {
	// Fetch inventory from Inventory service
	inventory, err := uc.inventoryClient.GetInventorySnapshot(ctx)
	if err != nil {
		inventory = nil
	}

	report := &domain.InventoryReport{
		Items: []domain.IngredientStock{},
	}

	// Mock inventory data (in production, aggregate real inventory)
	if inventory != nil {
		_ = inventory
	}

	report.Items = []domain.IngredientStock{
		{
			IngredientID: "ING001",
			Name:         "Rice",
			CurrentStock: 500,
			MinimumStock: 100,
			UnitCost:     50000,
			TotalValue:   25000000,
			IsLowStock:   false,
		},
		{
			IngredientID: "ING002",
			Name:         "Fish",
			CurrentStock: 150,
			MinimumStock: 50,
			UnitCost:     200000,
			TotalValue:   30000000,
			IsLowStock:   false,
		},
		{
			IngredientID: "ING003",
			Name:         "Vegetables",
			CurrentStock: 30,
			MinimumStock: 100,
			UnitCost:     100000,
			TotalValue:   3000000,
			IsLowStock:   true,
		},
	}

	report.TotalInventoryValue = 58000000.0
	report.LowStockItemsCount = 1
	report.TotalWasteValue = 500000.0

	_ = uc.repo.SaveOrderData(ctx, nil)
	return report, nil
}

// GetOrderReport generates order analytics report
func (uc *ReportUseCase) GetOrderReport(ctx context.Context, dateRange *domain.DateRange) (*domain.OrderReport, error) {
	if err := dateRange.ValidatePeriod(); err != nil {
		return nil, err
	}

	// Fetch orders from Order service
	orders, err := uc.orderClient.GetOrdersByDateRange(ctx, dateRange.FromDate, dateRange.ToDate)
	if err != nil {
		orders = nil
	}

	report := &domain.OrderReport{
		HourlyBreakdown: []domain.HourlyOrders{},
	}

	// Mock order data (in production, aggregate real orders)
	if orders != nil {
		_ = orders
	}

	report.TotalOrders = 150
	report.CompletedOrders = 140
	report.CancelledOrders = 10
	report.CompletionRate = 93.3
	report.AvgPreparationTime = 25.5

	// Hourly breakdown
	for hour := 11; hour <= 22; hour++ {
		count := int32(math.Sin(float64(hour-11)/3) * 20)
		if count < 0 {
			count = 0
		}
		if count > 25 {
			count = 25
		}
		report.HourlyBreakdown = append(report.HourlyBreakdown, domain.HourlyOrders{
			Hour:       int32(hour),
			OrderCount: count,
		})
	}

	_ = uc.repo.SaveOrderData(ctx, report)
	return report, nil
}

// GetStaffPerformanceReport generates staff performance metrics
func (uc *ReportUseCase) GetStaffPerformanceReport(ctx context.Context, staffID string, dateRange *domain.DateRange) ([]domain.StaffPerformance, error) {
	if err := dateRange.ValidatePeriod(); err != nil {
		return nil, err
	}

	// Fetch staff from User service
	staffList, err := uc.userClient.ListUsers(ctx)
	if err != nil {
		staffList = nil
	}

	var performances []domain.StaffPerformance

	if staffID != "" {
		// Single staff performance
		performance := domain.StaffPerformance{
			StaffID:            staffID,
			StaffName:          fmt.Sprintf("Staff %s", staffID),
			OrdersHandled:      45,
			TotalSales:         12000000,
			AverageOrderValue:  266667,
			CustomerComplaints: 2,
		}
		_ = uc.repo.SaveStaffPerformance(ctx, staffID, &performance)
		performances = append(performances, performance)
	} else {
		// All staff performance
		if staffList != nil {
			_ = staffList
		}

		staffMembers := []struct {
			ID   string
			Name string
		}{
			{"STAFF001", "Nguyễn Văn A"},
			{"STAFF002", "Trần Thị B"},
			{"STAFF003", "Phạm Văn C"},
		}

		for _, staff := range staffMembers {
			performance := domain.StaffPerformance{
				StaffID:            staff.ID,
				StaffName:          staff.Name,
				OrdersHandled:      int32(40 + getRandomInt(10)),
				TotalSales:         float64(10000000 + getRandomInt(5)*1000000),
				AverageOrderValue:  250000,
				CustomerComplaints: int32(getRandomInt(3)),
			}
			_ = uc.repo.SaveStaffPerformance(ctx, staff.ID, &performance)
			performances = append(performances, performance)
		}
	}

	return performances, nil
}

// GetPopularItemsReport generates popular items report
func (uc *ReportUseCase) GetPopularItemsReport(ctx context.Context, topN int32, dateRange *domain.DateRange) ([]domain.PopularItem, error) {
	if topN <= 0 {
		topN = 10
	}

	if err := dateRange.ValidatePeriod(); err != nil {
		return nil, err
	}

	// Fetch menu items from Menu service
	menuItems, err := uc.menuClient.ListMenuItems(ctx)
	if err != nil {
		menuItems = nil
	}

	var items []domain.PopularItem

	// Mock popular items data
	if menuItems != nil {
		_ = menuItems
	}

	mockItems := []domain.PopularItem{
		{ItemID: "ITEM001", ItemName: "Cơm Tấm", QuantitySold: 250, Revenue: 7500000, Category: "Rice"},
		{ItemID: "ITEM002", ItemName: "Phở Bò", QuantitySold: 200, Revenue: 8000000, Category: "Noodles"},
		{ItemID: "ITEM003", ItemName: "Gà Nước Mắm", QuantitySold: 180, Revenue: 6300000, Category: "Chicken"},
		{ItemID: "ITEM004", ItemName: "Canh Chua Cá", QuantitySold: 150, Revenue: 4500000, Category: "Soup"},
		{ItemID: "ITEM005", ItemName: "Nem Rán", QuantitySold: 320, Revenue: 6400000, Category: "Appetizers"},
	}

	if int32(len(mockItems)) > topN {
		items = mockItems[:topN]
	} else {
		items = mockItems
	}

	_ = uc.repo.SavePopularItems(ctx, items)
	return items, nil
}

// GetRevenueAnalytics generates revenue trend analysis
func (uc *ReportUseCase) GetRevenueAnalytics(ctx context.Context, dateRange *domain.DateRange) (*domain.RevenueAnalytics, error) {
	if err := dateRange.ValidatePeriod(); err != nil {
		return nil, err
	}

	// Get sales report for trend data
	salesReport, err := uc.GetSalesReport(ctx, dateRange)
	if err != nil {
		return nil, err
	}

	analytics := &domain.RevenueAnalytics{
		TotalRevenue: salesReport.NetRevenue,
		GrowthRate:   12.5,
		TrendData:    salesReport.DailyBreakdown,
		PeakDay:      "2026-04-08",
		PeakHour:     "12:00",
	}

	return analytics, nil
}

// ExportReport exports report in specified format
func (uc *ReportUseCase) ExportReport(ctx context.Context, reportType string, format domain.ExportFormat, dateRange *domain.DateRange) ([]byte, string, string, error) {
	if reportType == "" {
		return nil, "", "", domain.ErrEmptyReportType
	}

	if err := dateRange.ValidatePeriod(); err != nil {
		return nil, "", "", err
	}

	var contentType string
	var filename string

	switch format {
	case domain.FormatPDF:
		contentType = "application/pdf"
		filename = fmt.Sprintf("report_%s_%d.pdf", reportType, time.Now().Unix())
	case domain.FormatExcel:
		contentType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
		filename = fmt.Sprintf("report_%s_%d.xlsx", reportType, time.Now().Unix())
	case domain.FormatCSV:
		contentType = "text/csv"
		filename = fmt.Sprintf("report_%s_%d.csv", reportType, time.Now().Unix())
	case domain.FormatJSON:
		contentType = "application/json"
		filename = fmt.Sprintf("report_%s_%d.json", reportType, time.Now().Unix())
	default:
		return nil, "", "", domain.ErrInvalidExportFormat
	}

	// Mock export data (in production, generate real file)
	mockData := []byte(fmt.Sprintf(`{"type":"%s","format":"%v","generated_at":"%s"}`, reportType, format, time.Now()))

	return mockData, filename, contentType, nil
}

// Helper function to generate pseudo-random numbers
func getRandomInt(max int) int {
	return (int(time.Now().UnixNano()) % max) + 1
}
