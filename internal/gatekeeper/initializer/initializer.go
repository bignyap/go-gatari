package initializer

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/bignyap/go-admin/internal/caching"
	"github.com/bignyap/go-admin/internal/common"
	"github.com/bignyap/go-admin/internal/database/sqlcgen"
	cachemanagement "github.com/bignyap/go-admin/internal/gatekeeper/service/CacheManagement"
	gatekeeping "github.com/bignyap/go-admin/internal/gatekeeper/service/GateKeeping"
	pubsublistener "github.com/bignyap/go-admin/internal/gatekeeper/service/PubSubListener"
	"github.com/bignyap/go-admin/internal/initialize"
	"github.com/bignyap/go-admin/internal/router"
	"github.com/bignyap/go-utilities/counter"
	"github.com/bignyap/go-utilities/logger/api"
	"github.com/bignyap/go-utilities/logger/config"
	"github.com/bignyap/go-utilities/logger/factory"
	"github.com/bignyap/go-utilities/pubsub"
	"github.com/bignyap/go-utilities/server"
	"github.com/go-playground/validator"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type GateKeeperService struct {
	Logger         api.Logger
	ResponseWriter *server.ResponseWriter
	DB             *sqlcgen.Queries
	Conn           *pgxpool.Pool
	Validator      *validator.Validate
	CacheContoller *caching.CacheController
	CacheManager   *cachemanagement.CacheManagementService
	// CounterWorker  *counter.CounterWorker
	Matcher      *gatekeeping.Matcher
	PubSubClient pubsub.PubSubClient
	Mode         string
	Target       string
	stopFlush    chan struct{}
}

func NewGateKeeperService(
	logger api.Logger,
	conn *pgxpool.Pool,
	validator *validator.Validate,
	pubSubClient pubsub.PubSubClient,
	cacheController *caching.CacheController,
	redisClient redis.UniversalClient,
	counterWorker *counter.CounterWorker,
	mode string,
	target string,
) *GateKeeperService {

	db := sqlcgen.New(conn)

	cacheManager := cachemanagement.NewCacheManagementService(
		db,
		conn,
		logger,
		validator,
		cacheController,
		redisClient,
		counterWorker,
	)

	return &GateKeeperService{
		Logger:         logger,
		Validator:      validator,
		DB:             db,
		Conn:           conn,
		CacheContoller: cacheController,
		CacheManager:   cacheManager,
		PubSubClient:   pubSubClient,
		Mode:           mode,
		Target:         target,
		stopFlush:      make(chan struct{}),
	}
}

func (s *GateKeeperService) Setup(server server.Server) error {
	setupLogger := s.Logger.WithComponent("server.Setup")
	setupLogger.Info("Starting")

	s.ResponseWriter = server.GetResponseWriter()

	router.RegisterGateKeeperHandlers(
		server.Router(),
		s.Logger,
		s.ResponseWriter,
		s.DB,
		s.Conn,
		s.Validator,
		s.Matcher,
		s.CacheContoller,
		s.Mode,
		s.Target,
	)

	// Start periodic DB flush (Redis -> DB only)
	cachemanagement.StartPeriodicFlush(s.CacheManager, 30*time.Second, s.stopFlush)

	setupLogger.Info("Completed")
	return nil
}

func (s *GateKeeperService) Shutdown() error {
	shtLogger := s.Logger.WithComponent("server.Shutdown")
	shtLogger.Info("Starting")

	// Stop periodic DB flush
	close(s.stopFlush)

	ctx := context.Background()

	// Flush local counters to Redis
	if err := s.CacheManager.CounterWorker.FlushNow(string(common.Usageprefix), ctx); err != nil {
		shtLogger.Error("Flush to Redis failed", err)
	}

	// Flush Redis -> DB
	s.CacheManager.SyncAggregatedToDB(ctx, string(common.Usageprefix), func(key string, val map[string]float64) error {
		return s.CacheManager.IncrementUsageFromCacheKey(ctx, key, val)
	})
	shtLogger.Info("Cache flushed")

	// Close Redis
	if err := s.CacheContoller.Close(); err != nil {
		shtLogger.Error("Error closing Redis", err)
	} else {
		shtLogger.Info("Redis connection closed")
	}

	if s.Conn != nil {
		s.Conn.Close()
		shtLogger.Info("Database connection pool closed")
	}

	shtLogger.Info("Completed")
	return nil
}

func (s *GateKeeperService) InitializeEPMatcher() {
	limit := int32(10000)

	listEndpoints, err := common.FetchAll(
		func(offset, limit int32) ([]sqlcgen.ListApiEndpointRow, error) {
			endpoints, err := s.DB.ListApiEndpoint(context.Background(), sqlcgen.ListApiEndpointParams{
				Limit:  limit,
				Offset: offset,
			})
			if err != nil {
				s.Logger.Fatal("couldn't retrieve endpoints", err)
			}
			return endpoints, nil
		}, limit,
	)
	if err != nil {
		s.Logger.Fatal("couldn't retrieve endpoints", err)
	}

	var endpoints []gatekeeping.Endpoint
	for _, endpoint := range listEndpoints {
		endpoints = append(endpoints, gatekeeping.Endpoint{
			Path:   endpoint.PathTemplate,
			Method: endpoint.HttpMethod,
			Code:   endpoint.EndpointName,
		})
	}

	s.Matcher = gatekeeping.NewMatcher()
	s.Matcher.Load(endpoints)
}

func (s *GateKeeperService) InitializePubSubListener() {
	pubSubListener := pubsublistener.NewPubSubListener(
		s.Logger, s.CacheContoller, s.Matcher, s.PubSubClient,
	)
	if err := pubSubListener.UpdateEPMatcher(); err != nil {
		s.Logger.Fatal("Failed to load pubsub listener", err)
	}
}

func InitializeGateKeeperServer() {
	if err := initialize.GetEnvVals(); err != nil {
		log.Fatalf("Failed to load environment variables: %v", err)
	}

	environment := os.Getenv("ENVIRONMENT")
	if environment == "" {
		environment = "dev"
	}
	mode := os.Getenv("GATEKEEPER_MODE")
	target := os.Getenv("PROXY_TARGET")

	if mode == "proxy" && target == "" {
		log.Fatal("PROXY_TARGET must be set in proxy mode")
	}

	// Logger
	var logConfig config.LogConfig
	if environment == "prod" {
		logConfig = config.ProductionConfig()
	} else {
		logConfig = config.DevelopmentConfig()
	}
	logger, _ := factory.NewLogger(logConfig)

	logWithComponent := func(component string, fn func() error) {
		logger.WithComponent(component).Info("Started")
		if err := fn(); err != nil {
			log.Fatalf("Failed in %s: %v", component, err)
		}
		logger.WithComponent(component).Info("Completed")
	}

	var conn *pgxpool.Pool
	logWithComponent("LoadDBConn", func() error {
		var err error
		conn, err = initialize.LoadDBConn()
		return err
	})
	defer conn.Close()

	var pubSubClient pubsub.PubSubClient
	logWithComponent("LoadPubSub", func() error {
		var err error
		pubSubClient, err = initialize.LoadPubSub()
		return err
	})
	defer pubSubClient.Close()

	validator := validator.New()

	var cacheController *caching.CacheController
	logWithComponent("LoadRedisController", func() error {
		var err error
		cacheController, err = initialize.LoadRedisController()
		return err
	})

	// Redis client wrapper
	redisClient := cacheController.Redis()

	// Create counter worker
	counterWorker := counter.NewCounterWorker(redisClient, 5*time.Second, 100, 10000)
	go counterWorker.Start(context.Background())

	gkService := NewGateKeeperService(
		logger, conn, validator, pubSubClient,
		cacheController, redisClient, counterWorker,
		mode, target,
	)

	logWithComponent("InitializeEPMatcher", func() error {
		gkService.InitializeEPMatcher()
		return nil
	})

	logWithComponent("InitializePubSubListener", func() error {
		gkService.InitializePubSubListener()
		return nil
	})

	logWithComponent("InitializeWebServer", func() error {
		return initialize.InitializeWebServer(logger, gkService)
	})
}
