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

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	_ "github.com/jackc/pgx/v5/pgtype"
)

type Dataset struct {
	Id        string
	Timestamp int64
	Type      string
	AccelX    float32
	AccelY    float32
	AccelZ    float32
	GyroX     float32
	GyroY     float32
	GyroZ     float32
	QX        float32
	QY        float32
	QZ        float32
	QW        float32
	Yaw       float32
	Pitch     float32
	Roll      float32
}

type UserData struct {
	Id string `json:"id"`
}

type User struct {
	Id         string
	IsActive   bool
	LastCommit pgtype.Timestamp
	LastStatus string
	IsEnabled  bool
}

type Vector3 struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
	Z float32 `json:"z"`
}

type Quaternion struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
	Z float32 `json:"z"`
	W float32 `json:"w"`
}

type SensorSample struct {
	Timestamp int64      `json:"timestamp"`
	Accel     Vector3    `json:"acc"`
	Gyro      Vector3    `json:"gyro"`
	Quat      Quaternion `json:"quat"`
	Yaw       float32    `json:"yaw"`
	Pitch     float32    `json:"pitch"`
	Roll      float32    `json:"roll"`
}

type SensorPacket struct {
	Id      string
	Samples []SensorSample
}

func AddUser(d UserData) error {
	rawJson, _ := json.Marshal(d)

	log.Println("Connecting user:", string(rawJson))

	sqlAddUser, _ := db.LoadQuery("add_user.sql")

	// TODO: убрать потом
	now := (time.Now().UTC().Add(time.Hour * 3)).Format("2006-01-02 15:04:05")
	log.Println(now)

	_, err := db.DB.Exec(
		context.Background(),
		sqlAddUser,
		d.Id, now, "User connected!")

	if err != nil {
		return fmt.Errorf("failed to connect user: %v", err)
	}

	log.Println("User connected successful:", string(rawJson))

	return nil
}

func RemoveUser(d UserData) error {
	rawJson, _ := json.Marshal(d)

	log.Println("Disconnecting user:", string(rawJson))

	sqlRemoveUser, _ := db.LoadQuery("remove_user.sql")

	// TODO: убрать потом
	now := (time.Now().UTC().Add(time.Hour * 3)).Format("2006-01-02 15:04:05")
	log.Println(now)

	time.Sleep(time.Second)
	_, err := db.DB.Exec(
		context.Background(),
		sqlRemoveUser,
		d.Id, now, "User diconnected!")

	if err != nil {
		return fmt.Errorf("failed to disconnect user: %v", err)
	}

	log.Println("User disconnected successful:", string(rawJson))

	return nil
}

func GetUsers() ([]User, error) {
	log.Println("Getting users...")

	sqlGetUser, _ := db.LoadQuery("get_user.sql")

	rows, err := db.DB.Query(
		context.Background(),
		sqlGetUser)

	if err != nil {
		log.Fatalln("Error getting users", err)
		return nil, fmt.Errorf("failed to get users list: %v", err)
	}
	defer rows.Close()

	var users []User

	for rows.Next() {
		var u User

		if err := rows.Scan(&u.Id, &u.IsActive, &u.LastCommit, &u.LastStatus, &u.IsEnabled); err != nil {
			log.Fatalln("Error parsing users", err)
			return nil, err
		}

		users = append(users, u)
	}
	log.Println("Got users successfully! Amount =", len(users))

	return users, rows.Err()
}

func Add2RawSet(p SensorPacket) error {
	batch := &pgx.Batch{}

	log.Println("Adding", len(p.Samples), "samples in a raw dataset")

	sqlAdd2Raw, _ := db.LoadQuery("add_2_raw.sql")

	for _, s := range p.Samples {
		batch.Queue(sqlAdd2Raw,
			p.Id, s.Timestamp, "test",
			s.Accel.X, s.Accel.Y, s.Accel.Z,
			s.Gyro.X, s.Gyro.Y, s.Gyro.Z,
			s.Quat.X, s.Quat.Y, s.Quat.Z, s.Quat.W,
			s.Yaw, s.Pitch, s.Roll)
	}

	now := (time.Now().UTC().Add(time.Hour * 3)).Format("2006-01-02 15:04:05")
	sqlAddUser, _ := db.LoadQuery("add_user.sql")

	br := db.DB.SendBatch(context.Background(), batch)
	defer func(br pgx.BatchResults) {
		err := br.Close()
		if err != nil {
			log.Println("Error closing batch", err)
			_, err = db.DB.Exec(
				context.Background(),
				sqlAddUser,
				p.Id, now, "Error:"+err.Error())
		}
	}(br)

	log.Println("Added", len(p.Samples), "samples successfully!")

	_, err := db.DB.Exec(
		context.Background(),
		sqlAddUser,
		p.Id, now, "Samples added successfully!")

	if err != nil {
		return fmt.Errorf("failed to connect user: %v", err)
	}

	_, err = br.Exec()
	return err
}

func GetDataset(id []string, types []string, time1 int64, time2 int64) ([]Dataset, error) {
	log.Println("Generating dataset...")

	sqlGetDataset, _ := db.LoadQuery("generate_dataset.sql")

	rows, err := db.DB.Query(
		context.Background(),
		sqlGetDataset,
		id, types, time1, time2)

	if err != nil {
		log.Fatalln("Error generating dataset", err)
		return nil, fmt.Errorf("failed to generate dataset: %v", err)
	}
	defer rows.Close()

	var dataset []Dataset

	for rows.Next() {
		var d Dataset

		if err := rows.Scan(&d.Id, &d.Timestamp, &d.Type,
			&d.AccelX, &d.AccelY, &d.AccelZ,
			&d.GyroX, &d.GyroY, &d.GyroZ,
			&d.QX, &d.QY, &d.QZ, &d.QW,
			&d.Yaw, &d.Pitch, &d.Roll); err != nil {
			log.Fatalln("Error parsing dataset", err)
			return nil, err
		}

		dataset = append(dataset, d)
	}
	log.Println("Generated dataset successfully! Rows amount =", len(dataset))
	dataset = dataset[:100]

	return dataset, rows.Err()
}
