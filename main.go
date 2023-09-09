package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/logging"
)

var (
	logger    *logging.Logger
	projectID string
)

func main() {
	projectID = os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		projectID = "your-project-id"
	}

	h := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		// AddSource: true,
		Level: slog.LevelDebug,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.SourceKey {
				a.Key = "logging.googleapis.com/sourceLocation"
			}

			if a.Key == slog.MessageKey {
				a.Key = "message"
			}

			if a.Key == slog.LevelKey {
				a.Key = "severity"
			}

			return a
		},
	})
	slog.SetDefault(slog.New(h))

	ctx := context.Background()

	client, err := logging.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// NOTE: ローカル環境ではここで8秒程度かかるが、クラウド環境では数ミリ秒で終わる
	logger = client.Logger(
		"my-log",
		logging.RedirectAsJSON(os.Stdout),
	)

	slog.Info("starting server...")
	http.HandleFunc("/", handler)

	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		slog.Info(fmt.Sprintf("defaulting to port %s", port), slog.String("port", port))
	}

	// Start HTTP server.
	slog.Info(fmt.Sprintf("listening on port %s", port), slog.String("port", port))

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		slog.Error(err.Error())
		return
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	if strings.Contains(r.URL.Path, "favicon.ico") {
		http.NotFound(w, r)
		return
	}

	for i := 0; i < 3; i++ {
		le, err := logger.ToLogEntry(logging.Entry{
			HTTPRequest: &logging.HTTPRequest{
				Request: r,
			},
		}, projectID)
		if err != nil {
			log.Fatalf("Failed to create LogEntry: %v", err)
		}

		slog.Info(fmt.Sprintf("%s hello, cloud.google.com/go/logging with slog!", strconv.Itoa(i)),
			slog.String("name", "John"),
			slog.Int("age", 30),
			slog.String("timestamp", le.Timestamp.String()),
			slog.Any("logging.googleapis.com/labels", le.Labels),
			slog.String("logging.googleapis.com/insertId", le.InsertId),
			slog.String("logging.googleapis.com/spanId", le.SpanId),
			slog.String("logging.googleapis.com/trace", le.Trace),
			slog.Bool("logging.googleapis.com/trace_sampled", le.TraceSampled),
		)

		slog.Warn(fmt.Sprintf("%s warn, cloud.google.com/go/logging with slog!", strconv.Itoa(i)),
			slog.String("name", "Paul"),
			slog.Int("age", 27),
			slog.String("timestamp", le.Timestamp.String()),
			slog.Any("logging.googleapis.com/labels", le.Labels),
			slog.String("logging.googleapis.com/insertId", le.InsertId),
			slog.String("logging.googleapis.com/spanId", le.SpanId),
			slog.String("logging.googleapis.com/trace", le.Trace),
			slog.Bool("logging.googleapis.com/trace_sampled", le.TraceSampled),
		)

		slog.Error(fmt.Sprintf("%s error, cloud.google.com/go/logging with slog!", strconv.Itoa(i)),
			slog.String("name", "George"),
			slog.Int("age", 33),
			slog.String("timestamp", le.Timestamp.String()),
			slog.Any("logging.googleapis.com/labels", le.Labels),
			slog.String("logging.googleapis.com/insertId", le.InsertId),
			slog.String("logging.googleapis.com/spanId", le.SpanId),
			slog.String("logging.googleapis.com/trace", le.Trace),
			slog.Bool("logging.googleapis.com/trace_sampled", le.TraceSampled),
		)
	}

	name := os.Getenv("NAME")
	if name == "" {
		name = "World"
	}
	fmt.Fprintf(w, "[%s] Hello %s!\n", time.Now(), name)
}
