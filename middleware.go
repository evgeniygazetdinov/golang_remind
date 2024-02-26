package main

import (
	"log"
	"net"
	"net/http"
	"time"

	_ "sql-trainer-server/docs" // Этот импорт будет создан автоматически
)

// Обертка для ResponseWriter чтобы отслеживать статус и размер ответа
type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
	size       int64
}

func (w *responseWriterWrapper) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *responseWriterWrapper) Write(b []byte) (int, error) {
	size, err := w.ResponseWriter.Write(b)
	w.size += int64(size)
	return size, err
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wrapper := &responseWriterWrapper{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		startTime := time.Now()

		// Разделяем IP и порт клиента
		clientIP := r.RemoteAddr
		if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
			clientIP = host
		}

		// Логируем детали запроса
		log.Printf(
			"\n→ %s %s\n  Client IP: %s\n  Headers: %v",
			r.Method,
			r.RequestURI,
			clientIP,
			r.Header,
		)

		next.ServeHTTP(wrapper, r)

		duration := time.Since(startTime)
		log.Printf(
			"\n← %s %s\n  Status: %d\n  Duration: %v\n  Size: %d bytes\n",
			r.Method,
			r.RequestURI,
			wrapper.statusCode,
			duration,
			wrapper.size,
		)
	})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Разрешаем запросы с любого источника
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// Разрешаем все необходимые заголовки
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, X-CSRF-Token")

		// Разрешаем кредентиалы (если нужно)
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Обработка preflight запросов
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
