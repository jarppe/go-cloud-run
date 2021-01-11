package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jarppe/go-cloud-run/assets"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"html/template"
	"io"
	"log"
	"os"
	"time"
)

type Server struct {
	*echo.Echo
	*template.Template
	// db, queue, secrets, etc...
	assets *assets.Assets
	DB     *pgxpool.Pool
}

func requireEnv(envName string) string {
	value, found := os.LookupEnv(envName)
	if !found {
		log.Fatalf("missing environment value: %q", envName)
	}
	return value
}

func NewServer() *Server {
	assetsPath := requireEnv("ASSETS_PATH")
	if stat, err := os.Stat(assetsPath); err != nil {
		log.Fatalf("can't stat %q: %v", assetsPath, err)
	} else if !stat.IsDir() {
		log.Fatalf("assets path %q is not a directory", assetsPath)
	}
	mode := requireEnv("SERVER_MODE")
	if mode != "production" && mode != "development" {
		log.Fatalf("illegal value for mode: %q", mode)
	}

	assetsContext := assets.NewAssetsContext("assets", assetsPath)

	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s pool_min_conns=2 pool_max_conns=10",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"))
	db, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		log.Fatalf("unable to connect to database: %v", err)
	}

	var greeting string
	var now time.Time
	var data struct {
		Foo struct {
			Bar int64 `json:"bar"`
		} `json:"foo"`
	}
	err = db.QueryRow(context.Background(), `select 'heelo', now(), '{"foo": {"bar": 42}}'::jsonb`).Scan(&greeting, &now, &data)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			log.Fatalf("Query failed: %q", pgErr.Message) // => syntax error at end of input
		}
		log.Fatalf("Query error: %v\n", err)
	}

	log.Printf("DB response:\n\t%s\n\t%s\n\t%+v",
		greeting,
		now.Format(time.RFC3339),
		data)

	s := &Server{
		Echo:     echo.New(),
		assets:   assetsContext,
		Template: initTemplates(assetsContext),
		DB:       db,
	}

	s.HideBanner = true
	s.Debug = mode != "production"
	s.Renderer = s

	s.Use(middleware.Logger())
	s.Use(middleware.Recover())
	s.Use(middleware.Secure())
	s.Use(assetsContext.Middleware())

	return s
}

func (s *Server) Up() {
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
		if err := s.Start(host + ":" + port); err != nil {
			log.Printf("HTTP Server termination cause: %v (%[1]T)", err)
		}
	}()
}

func (s *Server) Down() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Printf("HTTP Server shutdown cause: %v (%[1]T)", err)
	}
	s.DB.Close()
}

func (s *Server) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return s.ExecuteTemplate(w, name, data)
}

func initTemplates(assets *assets.Assets) *template.Template {
	t := template.New("")
	t.Funcs(template.FuncMap{
		"toString": toString,
		"link":     link,
	})
	t.Funcs(template.FuncMap{
		"asset": assets.AssetRef,
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

	return t
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
