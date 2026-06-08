// Package repository contains in-memory implementations.
package repository

import (
	"context"
	"fmt"
	"restaurant-management/services/report-service/internal/domain"
	"sync"
	"time"
)

// InMemoryReportRepository implements ReportRepository using in-memory storage.
type InMemoryReportRepository struct {
	mu                sync.RWMutex
	salesData         map[string]*domain.SalesReport
	orderData         map[string]*domain.OrderReport
	staffPerformance  map[string]*domain.StaffPerformance
	popularItems      map[string][]domain.PopularItem
}

// NewInMemoryReportRepository creates a new in-memory report repository.
func NewInMemoryReportRepository() *InMemoryReportRepository {
	return &InMemoryReportRepository{
		salesData:        make(map[string]*domain.SalesReport),
		orderData:        make(map[string]*domain.OrderReport),
		staffPerformance: make(map[string]*domain.StaffPerformance),
		popularItems:     make(map[string][]domain.PopularItem),
	}
}

// SaveSalesData stores sales data for analysis
func (r *InMemoryReportRepository) SaveSalesData(ctx context.Context, salesData *domain.SalesReport) error {
	if salesData == nil {
		return domain.ErrUnableToGenerateReport
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	key := fmt.Sprintf("sales_%d", time.Now().Unix())
	r.salesData[key] = salesData

	return nil
}

// GetSalesData retrieves sales data for a date range
func (r *InMemoryReportRepository) GetSalesData(ctx context.Context, dateRange *domain.DateRange) (*domain.SalesReport, error) {
	if err := dateRange.ValidatePeriod(); err != nil {
		return nil, err
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	// For in-memory storage, return the most recent sales data
	// In production, would filter by date range
	var latestReport *domain.SalesReport
	for _, report := range r.salesData {
		if latestReport == nil {
			latestReport = report
		}
	}

	if latestReport == nil {
		return nil, domain.ErrNoDataAvailable
	}

	return latestReport, nil
}

// SaveOrderData stores order analytics data
func (r *InMemoryReportRepository) SaveOrderData(ctx context.Context, orderData *domain.OrderReport) error {
	if orderData == nil {
		return domain.ErrUnableToGenerateReport
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	key := fmt.Sprintf("orders_%d", time.Now().Unix())
	r.orderData[key] = orderData

	return nil
}

// GetOrderData retrieves order analytics
func (r *InMemoryReportRepository) GetOrderData(ctx context.Context, dateRange *domain.DateRange) (*domain.OrderReport, error) {
	if err := dateRange.ValidatePeriod(); err != nil {
		return nil, err
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	// Return most recent order data
	var latestReport *domain.OrderReport
	for _, report := range r.orderData {
		if latestReport == nil {
			latestReport = report
		}
	}

	if latestReport == nil {
		return nil, domain.ErrNoDataAvailable
	}

	return latestReport, nil
}

// SaveStaffPerformance stores staff metrics
func (r *InMemoryReportRepository) SaveStaffPerformance(ctx context.Context, staffID string, performance *domain.StaffPerformance) error {
	if staffID == "" || performance == nil {
		return domain.ErrUnableToGenerateReport
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	key := fmt.Sprintf("staff_%s_%d", staffID, time.Now().Unix())
	r.staffPerformance[key] = performance

	return nil
}

// GetStaffPerformance retrieves staff metrics
func (r *InMemoryReportRepository) GetStaffPerformance(ctx context.Context, staffID string, dateRange *domain.DateRange) (*domain.StaffPerformance, error) {
	if staffID == "" {
		return nil, domain.ErrUnableToGenerateReport
	}

	if err := dateRange.ValidatePeriod(); err != nil {
		return nil, err
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	// Find staff performance by ID
	for key, perf := range r.staffPerformance {
		if perf.StaffID == staffID {
			return perf, nil
		}
		_ = key
	}

	return nil, domain.ErrNoDataAvailable
}

// ListAllStaffPerformance retrieves all staff metrics for a period
func (r *InMemoryReportRepository) ListAllStaffPerformance(ctx context.Context, dateRange *domain.DateRange) ([]domain.StaffPerformance, error) {
	if err := dateRange.ValidatePeriod(); err != nil {
		return nil, err
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []domain.StaffPerformance
	seen := make(map[string]bool)

	for _, perf := range r.staffPerformance {
		if !seen[perf.StaffID] {
			result = append(result, *perf)
			seen[perf.StaffID] = true
		}
	}

	if len(result) == 0 {
		return nil, domain.ErrNoDataAvailable
	}

	return result, nil
}

// SavePopularItems stores popular items data
func (r *InMemoryReportRepository) SavePopularItems(ctx context.Context, items []domain.PopularItem) error {
	if len(items) == 0 {
		return domain.ErrUnableToGenerateReport
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	key := fmt.Sprintf("popular_%d", time.Now().Unix())
	r.popularItems[key] = items

	return nil
}

// GetPopularItems retrieves top N popular items
func (r *InMemoryReportRepository) GetPopularItems(ctx context.Context, topN int32, dateRange *domain.DateRange) ([]domain.PopularItem, error) {
	if topN <= 0 {
		topN = 10
	}

	if err := dateRange.ValidatePeriod(); err != nil {
		return nil, err
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	// Return most recent popular items list
	for _, items := range r.popularItems {
		if int32(len(items)) > 0 {
			if int32(len(items)) < topN {
				return items, nil
			}
			return items[:topN], nil
		}
	}

	return nil, domain.ErrNoDataAvailable
}
