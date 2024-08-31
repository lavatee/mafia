package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"github.com/lavatee/mafia"
	"github.com/lavatee/mafia/internal/endpoint"
	"github.com/lavatee/mafia/internal/repository"
	"github.com/lavatee/mafia/internal/service"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	if err := InitConfig(); err != nil {
		log.Fatalf("config error", err.Error())
	}
	// wd, err := os.Getwd()
	// if err != nil {
	// 	log.Fatalf("wd error", err.Error())
	// }
	if err := godotenv.Load(); err != nil {
		log.Fatalf("env error", err.Error())
	}
	db, err := repository.NewPostgresDB(repository.PostgresDB{Host: viper.GetString("db.host"), Port: viper.GetString("db.port"), Username: viper.GetString("db.username"), Password: os.Getenv("DB_PASSWORD"), DBName: viper.GetString("db.dbname"), SSLMode: viper.GetString("db.sslmode")})
	if err != nil {
		log.Fatalf("db error"+os.Getenv("DB_PASSWORD"), err.Error())
	}
	mongoClient, mongoFriendsDB, err := repository.NewMongoDB()
	if err != nil {
		log.Fatalf("mongo error", err.Error())
	}
	rdb := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})
	if err != nil {
		log.Fatalf("postgres error", err.Error())
	}
	repo := repository.NewRepository(db, mongoFriendsDB, rdb)
	svc := service.NewService(repo)
	endp := endpoint.NewEndpoint(svc)
	handler := endp.InitRoutes()
	srv := new(mafia.Server)
	go func() {

		err := srv.Run(viper.GetString("port"), handler)
		if err != nil {
			log.Fatalf("srv err", err.Error())
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatalf("shutdown error", err.Error())
	}
	rdb.Close()
	mongoClient.Disconnect(context.Background())
}
func InitConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
