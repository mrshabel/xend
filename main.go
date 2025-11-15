package main

// Xend is a local file-server that allows serving of a directory over http

import (
	"compress/gzip"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

// responseWriter extends the http.ResponseWriter to get additional insights on the request
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func (r *responseWriter) Write(b []byte) (int, error) {
	// record response size. data could be chunked so size is progressively updated
	size, err := r.ResponseWriter.Write(b)
	r.size += size
	return size, err
}

func (r *responseWriter) WriteHeader(statusCode int) {
	// get response status code and write on connection
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

// compressionResponseWriter extends the http.ResponseWriter to enable compression for text-based files
type compressionResponseWriter struct {
	io.Writer
	http.ResponseWriter

	// flag to determine if the compression header has been set already or not
	isCompressed bool
}

func (w *compressionResponseWriter) WriteHeader(statusCode int) {
	if !w.isCompressed {
		// set compression header
		w.Header().Set("Content-Encoding", "gzip")
		w.isCompressed = true
	}
	// writer status
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *compressionResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func main() {
	// entry flags
	host := flag.String("host", "localhost", "HTTP network host to listen on")
	port := flag.Int("port", 8000, "HTTP network port to listen on")
	dir := flag.String("dir", ".", "Root directory to serve files from")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "xend - A local file-server\n\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: xend [options]\n\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Example: xend -host 0.0.0.0 -port 9000 -dir ./public\n\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Options:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	addr := fmt.Sprintf("%s:%d", *host, *port)
	srv := &http.Server{
		Addr: addr,
	}

	done := make(chan struct{})
	fileServer := http.FileServer(http.Dir(*dir))
	handler := composeMiddlewares(http.StripPrefix("/", fileServer))
	http.Handle("/", handler)

	// graceful shutdown in a goroutine
	go func() {
		ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		defer stop()
		<-ctx.Done()

		log.Println("shutting down gracefully, press Ctrl+C again to force")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("Server forced to shutdown with error: %v\n", err)
		}
		close(done)
	}()

	log.Printf("Starting server on %s, serving %q\n", addr, *dir)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("HTTP server error: %v", err)
	}
	// wait for shutdown signal
	<-done
	log.Println("Server shutdown complete")
}

// composeMiddlewares chains the app middlewares in their right order
func composeMiddlewares(handler http.Handler) http.Handler {
	return loggingMiddleware(
		securityMiddleware(
			compressionMiddleware(handler),
		),
	)
}

// log request details
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get request duration and remote host
		start := time.Now()
		res := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(res, r)
		duration := time.Since(start)

		log.Printf(`%s - "%s %s %s" %d %d %s`,
			r.RemoteAddr, r.Method, r.RequestURI, r.Proto, res.statusCode, res.size, duration)
	})
}

// gzip-based compression of assets
func compressionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// verify that client supports gzip compression
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		// set headers: https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Encoding
		w.Header().Set("Content-Encoding", "gzip")

		compressor := gzip.NewWriter(w)
		defer compressor.Close()

		// perform compression
		res := &compressionResponseWriter{Writer: compressor, ResponseWriter: w}
		// proceed to next handler
		next.ServeHTTP(res, r)
	})
}

func securityMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// verify that target path isnt hidden
		if checkHiddenDir(r.URL.Path) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 page not found"))
			return
		}

		next.ServeHTTP(w, r)
	})
}

// checkHiddenDir verifies if the target directory is a well-known hidden directory
func checkHiddenDir(path string) bool {
	// remove prefix / if present since url paths are os-agnostic
	path = strings.TrimPrefix(path, "/")

	parts := strings.SplitSeq(path, string(filepath.Separator))
	for part := range parts {
		if strings.HasPrefix(part, ".") {
			return true
		}
	}
	return false
}
