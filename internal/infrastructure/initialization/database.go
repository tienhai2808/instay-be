package initialization

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/InstaySystem/is_v2-be/internal/infrastructure/config"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/persistence/orm/translator"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type SQL struct {
	gorm *gorm.DB
	sql  *sql.DB
}

type NoSQL struct {
	client *mongo.Client
	db     *mongo.Database
}

func InitSQL(cfg config.SQLConfig) (*SQL, error) {
	dsn := fmt.Sprintf(
		"host=%s dbname=%s user=%s password=%s sslmode=%s",
		cfg.Host,
		cfg.DBName,
		cfg.User,
		cfg.Password,
		cfg.SSLMode,
	)

	gormLogger := logger.New(
		log.New(os.Stdout, "[DB] ", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
			ParameterizedQueries:      true,
		},
	)

	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:                 gormLogger,
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})
	if err != nil {
		return nil, err
	}
	if err = translator.RegisterGORMPostgresTranslator(gormDB); err != nil {
		return nil, err
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, err
	}

	db := &SQL{
		gorm: gormDB,
		sql:  sqlDB,
	}

	return db, nil
}

func (s *SQL) ORM() *gorm.DB {
	return s.gorm
}

func (s *SQL) Close() {
	_ = s.sql.Close()
}

func InitNoSQL(cfg config.NoSQLConfig) (*NoSQL, error) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOpts := options.Client().ApplyURI(cfg.URI).
		SetServerAPIOptions(serverAPI).
		SetConnectTimeout(5 * time.Second).
		SetServerSelectionTimeout(5 * time.Second).
		SetMaxPoolSize(cfg.MaxPoolSize).
		SetMinPoolSize(cfg.MinPoolSize).
		SetMaxConnIdleTime(cfg.MaxConnIdleTime).
		SetAppName("is_v2-be").
		SetRetryWrites(true).
		SetRetryReads(true)

	client, err := mongo.Connect(clientOpts)
	if err != nil {
		return nil, err
	}

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		disConnCtx, disconnectCancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer disconnectCancel()

		_ = client.Disconnect(disConnCtx)
		return nil, err
	}

	noSQL := &NoSQL{
		client: client,
		db:     client.Database(cfg.DBName),
	}

	return noSQL, nil
}

func (n *NoSQL) DB() *mongo.Database {
	return n.db
}

func (n *NoSQL) Collection(name string) *mongo.Collection {
	return n.db.Collection(name)
}

func (n *NoSQL) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_ = n.client.Disconnect(ctx)
}
