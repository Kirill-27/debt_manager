package main

import (
	"github.com/go-redis/redis"
	"github.com/kirill-27/debt_manager"
	"github.com/kirill-27/debt_manager/pkg/handler"
	"github.com/kirill-27/debt_manager/pkg/repository"
	"github.com/kirill-27/debt_manager/pkg/service"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"log"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := initConfig(); err != nil {
		logrus.Fatalf("error init config: %s", err.Error())
	}
	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: viper.GetString("db.password"),
	})
	if err != nil {
		logrus.Fatalf("failed to init db: %s", err.Error())
	}

	client := redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.host") + ":" + viper.GetString("redis.port"),
		Password: "",
		DB:       0,
	})

	repo := repository.NewRepository(db, client)
	services := service.NewService(repo)
	handlers := handler.NewHandler(services)
	srv := new(debt_manager.Server)
	if err := srv.Run(viper.GetString("port"), handlers.InitRouters()); err != nil {
		log.Fatalf("error when run http server: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
