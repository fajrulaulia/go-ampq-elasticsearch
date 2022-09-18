package ampqgotest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/elastic/go-elasticsearch/esapi"
	amqp "github.com/rabbitmq/amqp091-go"
)

func (c *Config) Publisher(ctx context.Context, name string, payload map[string]interface{}) (err error) {
	ch, err := c.Ampq.Channel()
	if err != nil {
		err = fmt.Errorf("failed to open a channel: %s", err)
		return
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(name, false, false, false, false, nil)
	if err != nil {
		err = fmt.Errorf("failed to open a channel: %s", err)

	}

	var b bytes.Buffer
	encoder := json.NewEncoder(&b)
	if err := encoder.Encode(payload); err != nil {
		return err
	}

	if err = ch.PublishWithContext(ctx, "", q.Name, false, false, amqp.Publishing{
		ContentType: "text/json",
		Body:        b.Bytes(),
	}); err != nil {
		err = fmt.Errorf("failed to open a publish with context: %s", err)
	}
	return err
}

func (c *Config) EsCreate(payload map[string]interface{}) (err error) {

	data, err := json.Marshal(payload)
	if err != nil {
		log.Fatalf("Error marshaling document: %s", err)
	}

	req := esapi.IndexRequest{
		Index:      "testdb",
		DocumentID: payload["id"].(string),
		Body:       bytes.NewReader(data),
		Refresh:    "true",
	}
	log.Println(req)
	res, err := req.Do(context.Background(), c.Es)
	if err != nil {
		return fmt.Errorf("error getting response: %s", err.Error())
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("[%s] error indexing document ID=%s", res.Status(), payload["id"])
	}

	return err

}

// func serialize(msg Message) ([]byte, error) {
// 	var b bytes.Buffer
// 	encoder := json.NewEncoder(&b)
// 	err := encoder.Encode(msg)
// 	return b.Bytes(), err
// }

// func deserialize(b []byte) (Message, error) {
// 	var msg Message
// 	buf := bytes.NewBuffer(b)
// 	decoder := json.NewDecoder(buf)
// 	err := decoder.Decode(&msg)
// 	return msg, err
// }
