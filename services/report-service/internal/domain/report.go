// Package domain contains core business entities for the report service.
package domain

import (
	"time"
)

// ReportPeriod represents the time period for a report
type ReportPeriod int32

const (
	PeriodUnknown ReportPeriod = 0
	PeriodDaily   ReportPeriod = 1
	PeriodWeekly  ReportPeriod = 2
	PeriodMonthly ReportPeriod = 3
	PeriodYearly  ReportPeriod = 4
	PeriodCustom  ReportPeriod = 5
)

// ExportFormat represents the export file format
type ExportFormat int32

const (
	FormatUnknown ExportFormat = 0
	FormatPDF     ExportFormat = 1
	FormatExcel   ExportFormat = 2
	FormatCSV     ExportFormat = 3
	FormatJSON    ExportFormat = 4
)

// SalesReport contains sales analytics data
type SalesReport struct {
	TotalSales         float64
	TotalTax           float64
	TotalDiscount      float64
	NetRevenue         float64
	TotalOrders        int32
	AverageOrderValue  float64
	DailyBreakdown     []DailySales
	PaymentBreakdown   []PaymentMethodSales
}

// DailySales represents sales for a specific day
type DailySales struct {
	Date   string
	Sales  float64
	Orders int32
}

// PaymentMethodSales represents sales by payment method
type PaymentMethodSales struct {
	Method string
	Amount float64
	Count  int32
}

// InventoryReport contains inventory snapshot
type InventoryReport struct {
	TotalInventoryValue float64
	LowStockItemsCount  int32
	Items               []IngredientStock
	TotalWasteValue     float64
}

// IngredientStock represents stock level of an ingredient
type IngredientStock struct {
	IngredientID string
	Name         string
	CurrentStock float64
	MinimumStock float64
	UnitCost     float64
	TotalValue   float64
	IsLowStock   bool
}

// OrderReport contains order analytics
type OrderReport struct {
	TotalOrders         int32
	CompletedOrders     int32
	CancelledOrders     int32
	CompletionRate      float64
	AvgPreparationTime  float64
	HourlyBreakdown     []HourlyOrders
}

// HourlyOrders represents order count by hour
type HourlyOrders struct {
	Hour       int32
	OrderCount int32
}

// StaffPerformance represents individual staff metrics
type StaffPerformance struct {
	StaffID              string
	StaffName            string
	OrdersHandled        int32
	TotalSales           float64
	AverageOrderValue    float64
	CustomerComplaints   int32
}

// PopularItem represents a popular menu item
type PopularItem struct {
	ItemID       string
	ItemName     string
	QuantitySold int32
	Revenue      float64
	Category     string
}

// RevenueAnalytics contains revenue trend analysis
type RevenueAnalytics struct {
	TotalRevenue float64
	GrowthRate   float64
	TrendData    []DailySales
	PeakDay      string
	PeakHour     string
}

// DateRange represents a time range for reports
type DateRange struct {
	FromDate time.Time
	ToDate   time.Time
}

// ValidatePeriod validates report period dates
func (dr *DateRange) ValidatePeriod() error {
	if dr.FromDate.After(dr.ToDate) {
		return ErrInvalidDateRange
	}
	if dr.FromDate.Equal(dr.ToDate) && dr.FromDate.IsZero() {
		return ErrMissingDateRange
	}
	return nil
}
