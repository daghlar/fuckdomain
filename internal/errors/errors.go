package errors

import (
	"fmt"
	"runtime"
	"sync"
)

type ErrorType string

const (
	ErrorTypeDNS        ErrorType = "DNS"
	ErrorTypeHTTP       ErrorType = "HTTP"
	ErrorTypeConfig     ErrorType = "CONFIG"
	ErrorTypeValidation ErrorType = "VALIDATION"
	ErrorTypeIO         ErrorType = "IO"
	ErrorTypeNetwork    ErrorType = "NETWORK"
	ErrorTypeTimeout    ErrorType = "TIMEOUT"
	ErrorTypeRateLimit  ErrorType = "RATE_LIMIT"
	ErrorTypeUnknown    ErrorType = "UNKNOWN"
)

type AppError struct {
	Type      ErrorType
	Message   string
	Details   map[string]interface{}
	Err       error
	File      string
	Line      int
	Timestamp string
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Type, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Type, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func (e *AppError) WithDetails(details map[string]interface{}) *AppError {
	e.Details = details
	return e
}

func (e *AppError) WithError(err error) *AppError {
	e.Err = err
	return e
}

func NewError(errorType ErrorType, message string) *AppError {
	_, file, line, _ := runtime.Caller(1)
	return &AppError{
		Type:      errorType,
		Message:   message,
		Details:   make(map[string]interface{}),
		File:      file,
		Line:      line,
		Timestamp: fmt.Sprintf("%d", runtime.Nano()),
	}
}

func NewErrorWithError(errorType ErrorType, message string, err error) *AppError {
	_, file, line, _ := runtime.Caller(1)
	return &AppError{
		Type:      errorType,
		Message:   message,
		Details:   make(map[string]interface{}),
		Err:       err,
		File:      file,
		Line:      line,
		Timestamp: fmt.Sprintf("%d", runtime.Nano()),
	}
}

func WrapError(err error, message string) *AppError {
	if appErr, ok := err.(*AppError); ok {
		appErr.Message = message + ": " + appErr.Message
		return appErr
	}

	_, file, line, _ := runtime.Caller(1)
	return &AppError{
		Type:      ErrorTypeUnknown,
		Message:   message,
		Details:   make(map[string]interface{}),
		Err:       err,
		File:      file,
		Line:      line,
		Timestamp: fmt.Sprintf("%d", runtime.Nano()),
	}
}

func IsErrorType(err error, errorType ErrorType) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Type == errorType
	}
	return false
}

func GetErrorType(err error) ErrorType {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Type
	}
	return ErrorTypeUnknown
}

func GetErrorDetails(err error) map[string]interface{} {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Details
	}
	return nil
}

type ErrorCollector struct {
	errors []*AppError
	mu     sync.RWMutex
}

func NewErrorCollector() *ErrorCollector {
	return &ErrorCollector{
		errors: make([]*AppError, 0),
	}
}

func (ec *ErrorCollector) Add(err error) {
	ec.mu.Lock()
	defer ec.mu.Unlock()

	if appErr, ok := err.(*AppError); ok {
		ec.errors = append(ec.errors, appErr)
	} else {
		ec.errors = append(ec.errors, NewErrorWithError(ErrorTypeUnknown, "Unknown error", err))
	}
}

func (ec *ErrorCollector) AddError(errorType ErrorType, message string) {
	ec.mu.Lock()
	defer ec.mu.Unlock()
	ec.errors = append(ec.errors, NewError(errorType, message))
}

func (ec *ErrorCollector) AddErrorWithError(errorType ErrorType, message string, err error) {
	ec.mu.Lock()
	defer ec.mu.Unlock()
	ec.errors = append(ec.errors, NewErrorWithError(errorType, message, err))
}

func (ec *ErrorCollector) GetErrors() []*AppError {
	ec.mu.RLock()
	defer ec.mu.RUnlock()
	return append([]*AppError(nil), ec.errors...)
}

func (ec *ErrorCollector) GetErrorsByType(errorType ErrorType) []*AppError {
	ec.mu.RLock()
	defer ec.mu.RUnlock()

	var filtered []*AppError
	for _, err := range ec.errors {
		if err.Type == errorType {
			filtered = append(filtered, err)
		}
	}
	return filtered
}

func (ec *ErrorCollector) Count() int {
	ec.mu.RLock()
	defer ec.mu.RUnlock()
	return len(ec.errors)
}

func (ec *ErrorCollector) CountByType(errorType ErrorType) int {
	ec.mu.RLock()
	defer ec.mu.RUnlock()

	count := 0
	for _, err := range ec.errors {
		if err.Type == errorType {
			count++
		}
	}
	return count
}

func (ec *ErrorCollector) Clear() {
	ec.mu.Lock()
	defer ec.mu.Unlock()
	ec.errors = make([]*AppError, 0)
}

func (ec *ErrorCollector) HasErrors() bool {
	ec.mu.RLock()
	defer ec.mu.RUnlock()
	return len(ec.errors) > 0
}

func (ec *ErrorCollector) PrintSummary() {
	ec.mu.RLock()
	defer ec.mu.RUnlock()

	if len(ec.errors) == 0 {
		return
	}

	fmt.Println("\nError Summary:")
	fmt.Println("==============")

	typeCount := make(map[ErrorType]int)
	for _, err := range ec.errors {
		typeCount[err.Type]++
	}

	for errorType, count := range typeCount {
		fmt.Printf("%s: %d\n", errorType, count)
	}

	fmt.Printf("Total: %d\n", len(ec.errors))
}

func (ec *ErrorCollector) PrintDetailed() {
	ec.mu.RLock()
	defer ec.mu.RUnlock()

	if len(ec.errors) == 0 {
		return
	}

	fmt.Println("\nDetailed Errors:")
	fmt.Println("================")

	for i, err := range ec.errors {
		fmt.Printf("%d. [%s] %s\n", i+1, err.Type, err.Message)
		if err.Err != nil {
			fmt.Printf("   Caused by: %v\n", err.Err)
		}
		if len(err.Details) > 0 {
			fmt.Printf("   Details: %v\n", err.Details)
		}
		fmt.Printf("   Location: %s:%d\n", err.File, err.Line)
		fmt.Println()
	}
}
