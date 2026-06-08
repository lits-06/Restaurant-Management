// Package domain contains domain error definitions for the report service.
package domain

import "errors"

var (
	// ErrReportNotFound is returned when a report cannot be found
	ErrReportNotFound = errors.New("report not found")

	// ErrInvalidDateRange is returned when date range is invalid (from > to)
	ErrInvalidDateRange = errors.New("invalid date range: from date cannot be after to date")

	// ErrMissingDateRange is returned when required dates are missing
	ErrMissingDateRange = errors.New("date range is required")

	// ErrEmptyReportType is returned when report type is empty
	ErrEmptyReportType = errors.New("report type is required")

	// ErrInvalidReportType is returned for unsupported report types
	ErrInvalidReportType = errors.New("invalid report type")

	// ErrInvalidExportFormat is returned for unsupported export formats
	ErrInvalidExportFormat = errors.New("invalid export format")

	// ErrNoDataAvailable is returned when no data found for the requested period
	ErrNoDataAvailable = errors.New("no data available for the requested period")

	// ErrUnableToGenerateReport is returned when report generation fails
	ErrUnableToGenerateReport = errors.New("unable to generate report")
)
