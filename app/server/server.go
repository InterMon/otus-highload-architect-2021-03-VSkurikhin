package server

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	sa "github.com/savsgio/atreugo/v11"
	"github.com/savsgio/go-logger/v2"
	"github.com/vskurikhin/otus-highload-architect-2021-03-VSkurikhin/app/config"
	"github.com/vskurikhin/otus-highload-architect-2021-03-VSkurikhin/app/domain"
	"github.com/vskurikhin/otus-highload-architect-2021-03-VSkurikhin/app/security"
	"log"
	"os"
)

// Server определяет параметры для запуска HTTP-сервера.
type Server struct {
	DAO    *domain.DAO
	JWT    *security.JWT
	Server *sa.Atreugo
}

func gracefulClose(db *sql.DB) {
	// Настраиваем канал для отправки сигнальных уведомлений.
	// Нужно использовать буферизованный канал или есть риск пропустить сигнал
	// если не готовы принять сигнал при отправке.
	c := make(chan os.Signal, 1)

	// Блокировать до получения сигнала.
	s := <-c
	fmt.Println("Got signal:", s)
	err := db.Close()
	if err != nil {
		log.Println(err)
	}
}

// New инициализирует сервер для ответа на сетевые запросы HTTP.
func New(cfg *config.Config) *Server {

	c := sa.Config{
		Addr:             cfg.Server.Host + `:` + cfg.Server.Port,
		Compress:         true,
		Name:             "httpd",
		GracefulShutdown: true,
	}
	db := openDB(cfg)
	go gracefulClose(db)
	versionDB(db)

	return &Server{DAO: domain.New(db), JWT: security.New(cfg), Server: sa.New(c)}
}

func (s *Server) UseBefore(fns sa.Middleware) *sa.Router {
	return s.Server.UseBefore(fns)
}

func (s *Server) StaticCustom() *sa.Path {

	pathRewriteCalled := false

	return s.Server.StaticCustom("/", &sa.StaticFS{
		Root:               "web/public",
		GenerateIndexPages: true,
		AcceptByteRange:    true,
		PathRewrite: func(ctx *sa.RequestCtx) []byte {
			pathRewriteCalled = true

			return ctx.Path()
		},
		PathNotFound: func(ctx *sa.RequestCtx) error {
			return ctx.TextResponse("File not found", 404)
		},
	})
}

// GET устанавливает обработчик для GET запросов
func (s *Server) GET(url string, viewFn sa.View) *sa.Path {
	return s.Server.GET(url, viewFn)
}

// POST устанавливает обработчик для POST запросов
func (s *Server) POST(url string, viewFn sa.View) *sa.Path {
	return s.Server.POST(url, viewFn)
}

// PUT устанавливает обработчик для PUT запросов
func (s *Server) PUT(url string, viewFn sa.View) *sa.Path {
	return s.Server.PUT(url, viewFn)
}

func (s *Server) DELETE(url string, viewFn sa.View) *sa.Path {
	return s.Server.DELETE(url, viewFn)
}

// ListenAndServe запускает сервер для ответа на сетевые запросы HTTP.
func (s *Server) ListenAndServe() error {
	return s.Server.ListenAndServe()
}

func openDB(cfg *config.Config) *sql.DB {

	dbCFG := cfg.DataBase
	dsn := fmt.Sprintf(`%s:%s@tcp(%s:%d)/%s`, dbCFG.Username, dbCFG.Password, dbCFG.Host, dbCFG.Port, dbCFG.DBName)
	fmt.Println(dsn)
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		panic(err.Error())
	}
	return db
}

func versionDB(db *sql.DB) {

	if logger.DebugEnabled() {
		var version string
		err := db.QueryRow("SELECT VERSION()").Scan(&version)

		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		fmt.Println(version)
	}
}
