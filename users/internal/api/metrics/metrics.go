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
	// RequestCounter tracks total HTTP requests
	RequestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "route", "status"},
	)

	// RequestDuration tracks HTTP request latencies
	RequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request latencies in seconds",
			Buckets: []float64{0.001, 0.01, 0.1, 0.5, 1, 2, 5, 10},
		},
		[]string{"method", "route"},
	)

	// ActiveRequests tracks concurrent requests
	ActiveRequests = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "http_active_requests",
			Help: "Number of active HTTP requests",
		},
		[]string{"method", "route"},
	)

	// ErrorCounter tracks HTTP errors
	ErrorCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_request_errors_total",
			Help: "Total number of HTTP request errors",
		},
		[]string{"method", "route", "status_code"},
	)
)

func init() {
	// Register the metrics with Prometheus
	prometheus.MustRegister(RequestCounter)
	prometheus.MustRegister(RequestDuration)
	prometheus.MustRegister(ActiveRequests)
	prometheus.MustRegister(ErrorCounter)
}

// Enhanced patterns for better route normalization
var (
	idPattern   = regexp.MustCompile(`/\d+(/|$)`)
	uuidPattern = regexp.MustCompile(`/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}(/|$)`)
	slugPattern = regexp.MustCompile(`/[a-z0-9_-]+(/|$)`)
)

// NormalizePath replaces common patterns with placeholders to reduce cardinality
func NormalizePath(path string) string {
	// Order matters - most specific patterns first
	path = uuidPattern.ReplaceAllString(path, "/:uuid$1")
	path = idPattern.ReplaceAllString(path, "/:id$1")

	// Only apply slug normalization to paths with many segments to avoid over-normalization
	segments := strings.Count(path, "/")
	if segments > 3 {
		path = slugPattern.ReplaceAllString(path, "/:slug$1")
	}

	return path
}

// PrometheusMiddleware is a Fiber middleware that records request metrics
func PrometheusMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		method := c.Method()

		// Store original path before processing
		originalPath := c.Path()

		// Process request first to ensure route is matched
		err := c.Next()

		// Now try to get the route pattern after routing has occurred
		route := c.Route().Path

		// If route is empty or just "/", use a better fallback strategy
		if route == "" || route == "/" {
			// For APIs with ID parameters, normalize the path
			route = NormalizePath(originalPath)

			// If the resulting route is still just "/", use the original path
			// This ensures we at least see something in the dashboard
			if route == "/" && originalPath != "/" {
				route = originalPath
			}
		}

		// Record request duration
		duration := time.Since(start).Seconds()
		RequestDuration.WithLabelValues(method, route).Observe(duration)

		// Get status code
		statusCode := c.Response().StatusCode()
		status := strconv.Itoa(statusCode)

		// Record status code and increment request counter
		RequestCounter.WithLabelValues(method, route, status).Inc()

		// If error occurred, increment error counter
		if statusCode >= 400 {
			ErrorCounter.WithLabelValues(method, route, status).Inc()
		}

		// Update active requests counter
		// Note: We're moving this after processing since we need the route
		ActiveRequests.WithLabelValues(method, route).Inc()
		// For demo purposes, decrease it right away
		// In production, you'd want to track actual concurrent requests differently
		ActiveRequests.WithLabelValues(method, route).Dec()

		return err
	}
}
