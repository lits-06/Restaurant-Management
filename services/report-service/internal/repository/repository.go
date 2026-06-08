// Package repository contains repository interfaces and implementations.
package repository

import (
	"context"
	"restaurant-management/services/report-service/internal/domain"
)

// ReportRepository defines operations for report data access.
type ReportRepository interface {
	// SaveSalesData stores sales data for analysis
	SaveSalesData(ctx context.Context, salesData *domain.SalesReport) error

	// GetSalesData retrieves sales data for a date range
	GetSalesData(ctx context.Context, dateRange *domain.DateRange) (*domain.SalesReport, error)

	// SaveOrderData stores order analytics data
	SaveOrderData(ctx context.Context, orderData *domain.OrderReport) error

	// GetOrderData retrieves order analytics
	GetOrderData(ctx context.Context, dateRange *domain.DateRange) (*domain.OrderReport, error)

	// SaveStaffPerformance stores staff metrics
	SaveStaffPerformance(ctx context.Context, staffID string, performance *domain.StaffPerformance) error

	// GetStaffPerformance retrieves staff metrics
	GetStaffPerformance(ctx context.Context, staffID string, dateRange *domain.DateRange) (*domain.StaffPerformance, error)

	// ListAllStaffPerformance retrieves all staff metrics for a period
	ListAllStaffPerformance(ctx context.Context, dateRange *domain.DateRange) ([]domain.StaffPerformance, error)

	// SavePopularItems stores popular items data
	SavePopularItems(ctx context.Context, items []domain.PopularItem) error

	// GetPopularItems retrieves top N popular items
	GetPopularItems(ctx context.Context, topN int32, dateRange *domain.DateRange) ([]domain.PopularItem, error)
}
