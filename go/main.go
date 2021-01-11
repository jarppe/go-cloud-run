package main

import (
	"bytes"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.SetPrefix("go-cloud-run: ")
	log.Print("Server starting...")

	s := NewServer()

	s.GET("/ping", s.ping())
	s.GET("/hello/:name", s.hello(), middleware.Gzip())
	s.GET("/", s.index())
	s.GET("/turbo/", s.turbo())
	s.GET("/turbo/:page", s.turbo())
	s.GET("/*", s.notFound())

	s.Up()
	log.Print("Server ready")
	WaitForSignal()
	s.Down()

	log.Printf("Server terminated")
	os.Exit(0)
}


func (s *Server) ping() func(c echo.Context) error {
	return func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	}
}

func (s *Server) hello() func(c echo.Context) error {
	return func(c echo.Context) error {
		name := c.Param("name")
		message := c.QueryParam("message")
		if message == "" {
			message = "Greetings"
		}
		return c.Render(http.StatusOK, "hello.html", map[string]string{
			"Message": message,
			"Name":    name,
		})
	}
}

func (s *Server) index() echo.HandlerFunc {
	var buf bytes.Buffer
	if err := s.ExecuteTemplate(&buf, "index.html", nil); err != nil {
		log.Fatalf("Can't generate index.html: %v", err)
	}
	return s.assets.Handler(buf.Bytes(), "text/html; charset=utf-8")
}

func (s *Server) notFound() func(c echo.Context) error {
	return func(c echo.Context) error {
		return c.Render(http.StatusNotFound, "404.html", nil)
	}
}

func (s *Server) turbo() echo.HandlerFunc {
	return func(c echo.Context) error {
		page := "/turbo/" + c.Param("page")
		return c.Render(http.StatusOK, "turbo.html", map[string]interface{}{
			"PageName": page,
		})
	}
}

func WaitForSignal() {
	sigCh := make(chan os.Signal)
	signal.Notify(sigCh,
		syscall.SIGINT,
		syscall.SIGHUP,
		syscall.SIGTERM,
		syscall.SIGUSR2)
	sig := <-sigCh
	signal.Reset(sig)
	log.Printf("Got signal %q, terminating...", sig)
}
