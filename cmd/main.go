package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/spf13/viper"
	"log"
	"ridenow/notifier"
	"ridenow/notifier/models"
	"ridenow/notifier/queue"
)

type Env struct {
	db   models.Datastore
	cons *queue.QueueConsumer
}

func main() {
	// get config
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Panic(err)
	}
	dbUser := viper.GetString("DB_USER")
	dbPass := viper.GetString("DB_PASSWORD")
	dbHost := viper.GetString("DB_HOST")
	dbPort := viper.GetString("DB_PORT")
	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/notifier?sslmode=disable", dbUser, dbPass, dbHost, dbPort)

	qUser := viper.GetString("AMQP_USER")
	qPass := viper.GetString("AMQP_PASSWORD")
	qHost := viper.GetString("AMQP_HOST")
	qPort := viper.GetString("AMQP_PORT")
	qUrl := fmt.Sprintf("amqp://%s:%s@%s:%s/", qUser, qPass, qHost, qPort)

	// set up postgresql connection
	db, err := models.NewDB(dbUrl)
	if err != nil {
		log.Panic(err)
	}
	// set up rabbitmq connection
	cons, err := queue.NewQueueConsumer(qUrl)
	if err != nil {
		log.Panic(err)
	}
	env := &Env{db, cons}

	matches, err := env.cons.Subscribe("ridenow.users.match")

	wait := make(chan bool)

	go func() {
		for m := range matches {
			match := &notifier.Match{}
			err := proto.Unmarshal(m.Body, match)
			if err != nil {
				log.Panic(err)
			}
			not := models.NewNotification(*match.User.Email, *match.Location.Name, *match.WaveHeightM, *match.Time)
			fmt.Printf("%+v\n", not)
			_, err = db.StoreNotification(not)
			if err != nil {
				log.Panic(err)
			}
		}
	}()
	log.Printf(" [*] Running `notifier` service . To exit press CTRL+C")
	<-wait
}
