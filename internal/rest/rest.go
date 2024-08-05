package rest

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"github.com/saur4ig/file-storage/internal/config"
	"github.com/saur4ig/file-storage/internal/database"
	"github.com/saur4ig/file-storage/internal/rest/api"
	"github.com/saur4ig/file-storage/internal/rest/middleware"
	"github.com/saur4ig/file-storage/internal/services"
	si "github.com/saur4ig/file-storage/internal/services/interface"
)

func CreateServer(conf config.Config) {
	// initialize Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", conf.Cache.Host, conf.Cache.Port),
		DB:   0,
	})
	log.Info().Msg("Redis client initialized")

	// initialize database connection
	dbClient := newConnection(conf.DB)
	log.Info().Msg("Database connected")

	// initialize services and cache
	rc := database.NewRedisCache(redisClient)
	folderS, fileS, transactionS, s3S := initDBServices(dbClient)
	log.Info().Msg("Services initialized")

	// create API handler
	handler := api.New(folderS, fileS, transactionS, s3S, rc)

	// setup routes
	router := http.NewServeMux()
	withRoutes := routes(router, handler)
	log.Info().Msg("Routes set")

	// setup middleware
	withMiddleware := middleware.Logging(middleware.Auth(withRoutes))
	log.Info().Msg("Middleware initialized")

	// create and start server
	server := http.Server{
		Addr:    ":8080",
		Handler: withMiddleware,
	}

	log.Info().Msg("Server listening on port 8080")
	log.Fatal().Err(server.ListenAndServe())
}

// newConnection initializes the database connection pool with retries
func newConnection(cfg config.DbConfig) *sql.DB {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name)

	var db *sql.DB
	var err error

	// retry up to 10 times
	for attempts := 1; attempts <= 10; attempts++ {
		db, err = sql.Open("postgres", dsn)
		if err != nil {
			log.Info().Msgf("Attempt %d: Error opening database: %s", attempts, err.Error())
			time.Sleep(1 * time.Second)
			continue
		}

		// set connection pool settings
		db.SetMaxOpenConns(10)                 // maximum number of open connections to the database
		db.SetMaxIdleConns(10)                 // maximum number of idle connections in the pool
		db.SetConnMaxLifetime(5 * time.Minute) // maximum amount of time a connection may be reused

		if err = db.Ping(); err == nil {
			log.Info().Msgf("Successfully connected to the database on attempt %d", attempts)
			return db // connection successful
		}

		// if not successful -> close the db and wait for 1 second before retrying
		db.Close()
		log.Info().Msgf("Attempt %d: Error connecting to database: %s", attempts, err.Error())
		time.Sleep(1 * time.Second)
	}

	// if all attempts fail -> log and panic
	log.Panic().Msg("Unable to connect to the database after 10 attempts")
	return nil
}

// initializes all services
func initDBServices(db *sql.DB) (si.FolderService, si.FileService, si.TransactionService, si.FileStorage) {
	folderRepo := database.NewFolderRepository(db)
	fileRepo := database.NewFileRepository(db)
	transactionRepo := database.NewTransactionRepository(db)

	folderService := services.NewFolderService(folderRepo, fileRepo, db)
	fileService := services.NewFileService(folderRepo, fileRepo, db)
	transactionService := services.NewTransactionService(transactionRepo)
	s3Service := services.NewS3Service()

	return folderService, fileService, transactionService, s3Service
}
