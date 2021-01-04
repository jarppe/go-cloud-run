package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/jarppe/go-cloud-run/assets"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type Server struct {
	// db, queue, secrets, etc...
	echo     *echo.Echo
	assets   *assets.Assets
	renderer templates
}

func NewServer() *Server {
	e := echo.New()
	e.HideBanner = true

	assetsPath, found := os.LookupEnv("ASSETS")
	if !found {
		panic(`missing environment value: "ASSETS"`)
	}
	assetsContext := assets.NewAssetsContext("assets", assetsPath)
	renderer := initTemplates(assetsContext)

	mode, modeSet := os.LookupEnv("MODE")
	if !modeSet {
		mode = "production"
	}
	e.Debug = mode != "production"
	e.Renderer = renderer
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())
	e.Use(assetsContext.Middleware())

	server := &Server{
		echo:     e,
		assets:   assetsContext,
		renderer: renderer,
	}

	server.routing(e)
	return server
}

func (s *Server) Start() {
	go func() {
		host, hostSet := os.LookupEnv("HOST")
		if !hostSet {
			host = "0.0.0.0"
		}
		port, portSet := os.LookupEnv("PORT")
		if !portSet {
			port = "8080"
		}

		log.Printf("Server listenig at %s:%s", host, port)
		if err := s.echo.Start(host + ":" + port); err != nil {
			log.Printf("HTTP Server termination cause: %v (%[1]T)", err)
		}
	}()
}

func (s *Server) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.echo.Shutdown(ctx); err != nil {
		log.Printf("HTTP Server shutdown cause: %v (%[1]T)", err)
	}
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
	if err := s.renderer.ExecuteTemplate(&buf, "index.html", nil); err != nil {
		log.Fatalf("Can't generate index.html: %v", err)
	}
	return s.assets.Handler(buf.Bytes(), "text/html; charset=utf-8")
}

func (s *Server) notFound() func(c echo.Context) error {
	return func(c echo.Context) error {
		return c.Render(http.StatusNotFound, "404.html", nil)
	}
}

func (s *Server) routing(e *echo.Echo) {
	e.GET("/ping", s.ping())
	e.GET("/hello/:name", s.hello(), middleware.Gzip())
	e.GET("/", s.index())
	e.GET("/turbo/", s.turbo())
	e.GET("/turbo/:page", s.turbo())
	// e.GET("/*", s.notFound())
}

func (s *Server) turbo() echo.HandlerFunc {
	return func(c echo.Context) error {
		page := "/turbo/" + c.Param("page")
		return c.Render(http.StatusOK, "turbo.html", map[string]interface{}{
			"PageName": page,
		})
	}
}

type templates struct {
	*template.Template
}

func (t templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.ExecuteTemplate(w, name, data)
}

func initTemplates(assets *assets.Assets) templates {
	t := template.New("")
	t.Funcs(template.FuncMap{
		"toString": toString,
		"link":     link,
		"asset":    assets.AssetRef,
	})

	// parse all html files in the templates directory
	t, err := t.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatal(err)
	}

	// also parse files in subdirectories of templates
	//t, err = t.ParseGlob("templates/*/*.html")
	//if err != nil {
	//	log.Fatal(err)
	//}

	return templates{t}
}

func toString(v interface{}) string {
	return fmt.Sprint(v)
}

func link(location, name string) template.HTML {
	return escSprintf(`<a class="text-blue-600 no-underline hover:underline" href="%v">%v</a>`, location, name)
}

func escSprintf(format string, args ...interface{}) template.HTML {
	for i, arg := range args {
		args[i] = template.HTMLEscapeString(fmt.Sprint(arg))
	}
	return template.HTML(fmt.Sprintf(format, args...))
}
