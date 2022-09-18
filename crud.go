package ampqgotest

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

const (
	CREATE = "req.create.todo"
)

type TODO struct {
	ID      string    `json:"id,omitempty"`
	Name    string    `json:"name,omitempty"`
	Created time.Time `json:"created,omitempty"`
}

type ReqTODO struct {
	Name string `json:"name,omitempty"`
}

func (c *Config) SetupRoute() *Config {
	c.Route.HandleFunc("/todo", c.Create).Methods("POST")
	return c
}

func (c *Config) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	var p ReqTODO
	err := dec.Decode(&p)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	var data TODO
	uuid := uuid.New()
	data.ID = uuid.String()
	data.Name = p.Name
	data.Created = time.Now()

	m := make(map[string]interface{})
	m["id"] = data.ID
	m["name"] = data.Name
	m["created"] = data.Created

	if err := c.EsCreate(m); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	if err := c.Publisher(ctx, CREATE, m); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	jsonResp, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonResp)

}
