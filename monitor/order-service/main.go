package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/net/context"
	"log"
	"order-service/pkg/client"
	"time"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"path"},
	)
	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path"},
	)

	activeConnections = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "active_connections",
			Help: "Number of active connections",
		},
	)
)

func init() {
	prometheus.MustRegister(httpRequestsTotal, httpRequestDuration)
}

func middleware(c *gin.Context) {
	httpRequestsTotal.WithLabelValues(c.Request.URL.Path).Inc()
	timer := prometheus.NewTimer(httpRequestDuration.WithLabelValues(c.Request.URL.Path))
	c.Next()
	timer.ObserveDuration()
}

func InitTracer(serviceName string) (*sdktrace.TracerProvider, error) {
	// Create the Jaeger exporter.
	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint("http://localhost:14268/api/traces")))
	if err != nil {
		return nil, err
	}

	// Create a new tracer provider with the Jaeger exporter.
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
		)),
	)
	otel.SetTracerProvider(tp)
	return tp, nil
}

var tracer = otel.Tracer("order-service")

func main() {
	tp, err := InitTracer("order-service")
	if err != nil {
		log.Fatalf("failed to initialize TracerProvider: %v", err)
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	router.Use(middleware)
	router.POST("/order", func(c *gin.Context) {
		ctx, span := tracer.Start(c.Request.Context(), "order")
		defer span.End()
		spanCtx := trace.SpanContextFromContext(ctx)
		fmt.Println(spanCtx.TraceID().String())
		GetProduct(ctx, "1")
		time.Sleep(5 * time.Second)
		c.JSON(200, gin.H{
			"message": "order created",
		})
	})
	if err := router.Run(":8089"); err != nil {
		panic(err)
	}
}

func GetProduct(ctx context.Context, productID string) (string, error) {
	ctx, span := tracer.Start(ctx, "GetProduct")
	defer span.End()
	return client.Get(ctx, "http://localhost:8087/product/"+productID)
}
