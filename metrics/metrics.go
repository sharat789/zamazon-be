package metrics

import (
	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	// Same metric declarations as before
	RequestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"route", "status"}, // Removed method to reduce cardinality
	)

	RequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request latencies in seconds",
			Buckets: []float64{0.001, 0.01, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
		},
		[]string{"route"},
	)

	ActiveRequests = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "http_active_requests",
			Help: "Number of active HTTP requests",
		},
		[]string{"route"},
	)

	ErrorCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_request_errors_total",
			Help: "Total number of HTTP request errors",
		},
		[]string{"route", "status_code"},
	)
)

func init() {
	prometheus.MustRegister(RequestCounter, RequestDuration, ActiveRequests, ErrorCounter)
}

// Complete mapping of all known routes based on your microservices
var knownRoutes = map[string]string{
	// Auth service routes
	"/auth/hash-password":     "/auth/hash-password",
	"/auth/verify-password":   "/auth/verify-password",
	"/auth/generate-token":    "/auth/generate-token",
	"/auth/verify-token":      "/auth/verify-token",
	"/auth/authorize-by-role": "/auth/authorize-by-role",
	"/auth/generate-code":     "/auth/generate-code",

	// Catalog service routes
	"/health":     "/health",
	"/products":   "/products",
	"/categories": "/categories",

	// Transactions service routes
	"/buyer/health":   "/buyer/health",
	"/buyer/verify":   "/buyer/verify",
	"/buyer/checkout": "/buyer/checkout",
	"/buyer/orders":   "/buyer/orders",

	// Users service routes
	"/users/register":   "/users/register",
	"/users/login":      "/users/login",
	"/users/health":     "/users/health",
	"/users/verifyUser": "/users/verifyUser",
	"/users/verify":     "/users/verify",
	"/users/profile":    "/users/profile",
	"/users/cart":       "/users/cart",
	"/users/order":      "/users/order",
}

// Path parameter patterns for normalization
var (
	idPattern         = regexp.MustCompile(`/(\d+|[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12})(/|$)`)
	productIdPattern  = regexp.MustCompile(`/cart/([^/]+)(/|$)`)
	orderIdPattern    = regexp.MustCompile(`/order/([^/]+)(/|$)`)
	categoryIdPattern = regexp.MustCompile(`/categories/([^/]+)(/|$)`)
	productIdPattern2 = regexp.MustCompile(`/products/([^/]+)(/|$)`)
)

// NormalizePath normalizes the request path to reduce cardinality
func NormalizePath(path string) string {
	// Skip normalizing the metrics endpoint itself
	if path == "/metrics" {
		return "/metrics"
	}

	// Try to match exact known routes first
	if normalizedPath, exists := knownRoutes[path]; exists {
		return normalizedPath
	}

	// Normalize paths with IDs for specific service patterns
	path = idPattern.ReplaceAllString(path, "/:id$2")
	path = productIdPattern.ReplaceAllString(path, "/cart/:productId$2")
	path = orderIdPattern.ReplaceAllString(path, "/order/:id$2")
	path = categoryIdPattern.ReplaceAllString(path, "/categories/:id$2")
	path = productIdPattern2.ReplaceAllString(path, "/products/:id$2")

	// Service-specific normalization
	if strings.HasPrefix(path, "/users/cart/") {
		return "/users/cart/:productId"
	}
	if strings.HasPrefix(path, "/users/order/") {
		return "/users/order/:id"
	}
	if strings.HasPrefix(path, "/buyer/order/") {
		return "/buyer/order/:id"
	}
	if strings.HasPrefix(path, "/products/") {
		return "/products/:id"
	}
	if strings.HasPrefix(path, "/categories/") {
		return "/categories/:id"
	}

	// Final check against known route prefixes
	for knownPath := range knownRoutes {
		if strings.HasPrefix(path, knownPath+"/") {
			return knownPath + "/:param"
		}
	}

	return path
}

// PrometheusMiddleware records request metrics
func PrometheusMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Skip metrics endpoint completely
		if c.Path() == "/metrics" {
			return c.Next()
		}

		// Get normalized route path
		route := NormalizePath(c.Path())

		// Record active requests
		ActiveRequests.WithLabelValues(route).Inc()
		defer ActiveRequests.WithLabelValues(route).Dec()

		// Time the request execution
		start := time.Now()
		err := c.Next()
		duration := time.Since(start).Seconds()

		// Record request duration
		RequestDuration.WithLabelValues(route).Observe(duration)

		// Record request count by status
		statusCode := c.Response().StatusCode()
		status := strconv.Itoa(statusCode)
		RequestCounter.WithLabelValues(route, status).Inc()

		// Record errors separately
		if statusCode >= 400 {
			ErrorCounter.WithLabelValues(route, strconv.Itoa(statusCode)).Inc()
		}

		return err
	}
}
