package main

import (
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/paudarco/todo-app"
	"github.com/paudarco/todo-app/pkg/handler"
	"github.com/paudarco/todo-app/pkg/repository"
	"github.com/spf13/viper"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("error reading config file: %v", err)
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		User:     viper.GetString("db.user"),
		Password: viper.GetString("db.password"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})

	if err != nil {
		log.Fatalf("failed to initialize db: %v", err.Error())
	}

	// db.MustExec(repository.Schema)

	fmt.Println("okok")

	repos := repository.NewRepository(db)
	handlers := handler.NewHandler(repos)

	srv := new(todo.Server)
	if err := srv.Run(viper.GetString("port"), handlers.InitRouters()); err != nil {
		log.Fatalf("error while running: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
