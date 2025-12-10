package rawData

import (
	"context"
	_ "context"
	"encoding/json"
	"eyeident/internal/db"
	_ "eyeident/internal/db"
	"fmt"
	"log"
	"time"
)

type RawData struct {
	DeviceID string      `json:"deviceId"`
	Model    string      `json:"model"`
	Data     interface{} `json:"data"`
}

func AddRaw(d RawData) error {
	rawJson, _ := json.Marshal(d)

	log.Println("Saving RawData:", string(rawJson))

	sqlAddUser, _ := db.LoadQuery("add_user.sql")

	// TODO: убрать потом
	now := time.Now().Format("2006-01-02 15:04:05")

	_, err := db.DB.Exec(
		context.Background(),
		sqlAddUser,
		d.Data, now, "Success!")

	if err != nil {
		return fmt.Errorf("failed to insert raw_data: %v", err)
	}

	log.Println("Saved successful:", string(rawJson))

	//
	//_, err := db.DB.Exec(
	//	context.Background(),
	//	"INSERT INTO raw_data (device_id, model, data) VALUES ($1, $2, $3)",
	//	d.DeviceID, d.Model, rawJson,
	//)

	return nil
}
