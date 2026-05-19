package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/Bayan2019/rbk-it-school-hw-6/internal/auth"
)

type contextKey string

const (
	userContextKey contextKey = "user"
	errorsKey      contextKey = "request_errors"
	RequestIDKey   contextKey = "request_id"
)

// type Roles string

// const (
// 	RolesAdmin Roles = "admin"
// 	RolesUser  Roles = "user"
// )

type UserContext struct {
	ID    int64      `json:"id"`
	Email string     `json:"email"`
	Role  auth.Roles `json:"role"`
}

type ErrorList struct {
	Errors []error
}

/// Middlewares ///
/// Middlewares ///
/// Middlewares ///

// 3. Create a Middleware to initialize and process the error slice
func ErrorCollectorMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Initialize the error list container
		errList := &ErrorList{Errors: make([]error, 0)}

		// Inject the pointer into the request context
		ctx := context.WithValue(r.Context(), errorsKey, errList)

		// Pass the new context down the chain
		next.ServeHTTP(w, r.WithContext(ctx))

		// AFTER the handlers run, process collected errors (like Gin middleware)
		if len(errList.Errors) > 0 {
			// Handle your errors centrally here (e.g., log them, write JSON)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	})
}

// 4. Helper function equivalent to Gin's c.Error(err)
func AppendError(r *http.Request, err error) {
	if err == nil {
		return
	}
	if errList, ok := r.Context().Value(errorsKey).(*ErrorList); ok {
		errList.Errors = append(errList.Errors, err)
	}
}

// 5. Helper function equivalent to Gin's c.Errors
func GetErrors(r *http.Request) *ErrorList {
	if errList, ok := r.Context().Value(errorsKey).(*ErrorList); ok {
		return errList
	}
	return nil
}

func (el *ErrorList) String() string {
	if el == nil || len(el.Errors) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("Errors:\n")

	for _, err := range el.Errors {
		if err != nil {
			// Formats exactly like Gin's standard format: Error #01: msg
			fmt.Fprintf(&sb, "Error: %s\n", err.Error())
		}
	}
	return sb.String()
}

/// User ///

func withUser(ctx context.Context, user UserContext) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

func UserFromContext(ctx context.Context) (UserContext, error) {
	user, ok := ctx.Value(userContextKey).(UserContext)
	if !ok {
		return UserContext{}, errors.New("user not found in context")
	}

	return user, nil
}
