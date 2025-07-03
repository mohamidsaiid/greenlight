package main

import (
	"context"
	"database/sql"
	"flag"
	"os"
	"sync"
	"time"

	_ "github.com/lib/pq"
	"github.com/mohamidsaiid/greenlight/internal/data"
	"github.com/mohamidsaiid/greenlight/internal/jsonlog"
	"github.com/mohamidsaiid/greenlight/internal/mailer"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxIdleTime  string
		maxOpenConns int
		maxIdleConns int
	}
	limiter struct {
		rps     float64
		brust   int
		enabled bool
	}
	smtp struct {
		host     string
		port     int
		username string
		password string
		sender   string
	}
}

type application struct {
	config config
	logger *jsonlog.Logger
	models data.Models
	mailer mailer.Mailer
	wg     sync.WaitGroup
}

func main() {

	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")

	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("GREENLIGHT_DB_DSN"), "PostgreSQL DSN")

	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")

	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.limiter.brust, "limiter-brust", 4, "Rate limiter maximum brust")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")

	flag.StringVar(&cfg.smtp.host, "smtp-host", "sandbox.smtp.mailtrap.io", "SMTP host")
	flag.IntVar(&cfg.smtp.port, "smtp-port", 587, "SMTP port")
	flag.StringVar(&cfg.smtp.username, "smtp-username", "cc40015f5c7095", "SMTP username")
	flag.StringVar(&cfg.smtp.password, "smpt-password", "e68b317fb0aa6f", "SMTP password")
	flag.StringVar(&cfg.smtp.sender, "smtp-sender", "Mohamed Saeed <no-reply@greenlight.com>", "SMTP sender")
	flag.Parse()

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	db, err := openDSN(cfg)
	if err != nil {
		logger.PrintFatal(err, nil)
	}
	defer db.Close()

	logger.PrintInfo("database connection pool established", nil)

	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
		mailer: mailer.New(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password, cfg.smtp.sender),
	}

	err = app.serve()
	if err != nil {
		logger.PrintFatal(err, nil)
	}
}

func openDSN(cfg config) (*sql.DB, error) {

	// start conntection with the db with the received dsn (a.k.a data source name)
	db, err := sql.Open("postgres", cfg.db.dsn)

	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(cfg.db.maxIdleConns)

	db.SetMaxOpenConns(cfg.db.maxOpenConns)

	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(duration)
	// create context with 5 seconds time out
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
