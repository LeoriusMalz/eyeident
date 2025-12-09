package rawData

import (
	_ "context"
	"encoding/json"
	"fmt"
	"log"
	_ "sensorsProject/internal/db"
)

type RawData struct {
	DeviceID string      `json:"deviceId"`
	Model    string      `json:"model"`
	Data     interface{} `json:"data"`
}

func AddRaw(d RawData) error {
	fmt.Println(d)
	log.Println(d)

	rawJson, _ := json.Marshal(d)

	fmt.Println(rawJson)
	log.Println(string(rawJson))

	//
	//_, err := db.DB.Exec(
	//	context.Background(),
	//	"INSERT INTO raw_data (device_id, model, data) VALUES ($1, $2, $3)",
	//	d.DeviceID, d.Model, rawJson,
	//)

	return nil
}
