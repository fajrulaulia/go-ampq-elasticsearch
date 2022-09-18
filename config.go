package ampqgotest

import (
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/gorilla/mux"
	amqp "github.com/rabbitmq/amqp091-go"
)

// Register confgiration and integration
type Config struct {
	Route *mux.Router
	Ampq  *amqp.Connection
	Es    *elasticsearch.Client
}

func (c *Config) Apply() *Config {
	c.SetupRoute()
	return c
}

func InitalizeRabbitMQ() (conn *amqp.Connection, err error) {
	conn, err = amqp.Dial("amqp://user:password@localhost:5672/")
	if err != nil {
		return conn, err
	}
	return conn, err
}

func IntializeElasticSearch() (es *elasticsearch.Client, err error) {
	es, err = elasticsearch.NewDefaultClient()
	if err != nil {
		return es, err
	}
	return es, err
}
