package rawData

import (
	"context"
	_ "context"
	"encoding/csv"
	"encoding/json"
	"eyeident/internal/db"
	_ "eyeident/internal/db"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	_ "github.com/jackc/pgx/v5/pgtype"
)

type Params struct {
	Id   []string `json:"id"`
	Type []string `json:"type"`
}

type Dataset struct {
	Id        string  `json:"id"`
	Timestamp int64   `json:"timestamp"`
	Type      string  `json:"type"`
	AccelX    float32 `json:"accel_x"`
	AccelY    float32 `json:"accel_y"`
	AccelZ    float32 `json:"accel_z"`
	GyroX     float32 `json:"gyro_x"`
	GyroY     float32 `json:"gyro_y"`
	GyroZ     float32 `json:"gyro_z"`
	QX        float32 `json:"qx"`
	QY        float32 `json:"qy"`
	QZ        float32 `json:"qz"`
	QW        float32 `json:"qw"`
	Yaw       float32 `json:"yaw"`
	Pitch     float32 `json:"pitch"`
	Roll      float32 `json:"roll"`
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
	Type      string     `json:"type"`
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

func GetUserEnable(p SensorPacket) (bool, error) {
	sqlGetUserEnable, _ := db.LoadQuery("get_user_enable.sql")

	rows, _ := db.DB.Query(
		context.Background(),
		sqlGetUserEnable,
		p.Id)
	defer rows.Close()

	fmt.Println(rows)

	var enabled bool
	for rows.Next() {
		if err := rows.Scan(&enabled); err != nil {
			log.Fatalln("Error parsing users", err)
			return false, err
		}
	}
	if !enabled {
		log.Println("Collecting for", p.Id, "is disabled!")
	}

	return enabled, nil
}

func GetUsers() ([]User, error) {
	log.Println("Getting users...")

	sqlGetUser, _ := db.LoadQuery("get_user.sql")
	sqlUpdateUser, _ := db.LoadQuery("update_user.sql")

	now := (time.Now().UTC().Add(time.Hour * 3)).Format("2006-01-02 15:04:05")
	_, err := db.DB.Exec(
		context.Background(),
		sqlUpdateUser,
		now, "User diconnected!")
	if err != nil {
		return nil, fmt.Errorf("failed to disconnect user: %v", err)
	}

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
			p.Id, s.Timestamp, s.Type,
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

func GetDataset(ids []string, types []string, time1 int64, time2 int64, limit int64, outputPath string) (string, error) {
	if len(ids) == 0 || len(types) == 0 {
		return "", fmt.Errorf("ids and types must not be empty")
	}

	ctx := context.Background()

	log.Println("Start getting dataset...")

	conn, err := db.DB.Acquire(ctx)
	if err != nil {
		log.Println("Error acquiring connection", err)
		return "", err
	}
	defer conn.Release()

	file, err := os.Create(outputPath)
	if err != nil {
		log.Println("Error creating file", err)
		return "", err
	}
	defer file.Close()

	log.Println("Temp file created!")
	log.Println("Forming SQL query!")

	idList := make([]string, len(ids))
	for i, v := range ids {
		idList[i] = fmt.Sprintf("'%s'", v)
	}

	typeList := make([]string, len(types))
	for i, v := range types {
		typeList[i] = fmt.Sprintf("'%s'", v)
	}

	sqlGetDataset, _ := db.LoadQuery("generate_dataset.sql")

	sql := fmt.Sprintf(sqlGetDataset,
		strings.Join(idList, ","),
		strings.Join(typeList, ","),
		time1,
		time2,
		limit,
	)

	log.Println("SQL formed successfully!")

	countCT, err := conn.Conn().PgConn().CopyTo(ctx, file, sql)
	if err != nil {
		log.Println("Error copying to file", err)
		return "", err
	}
	//var count, _ = strconv.ParseInt(countCT.String(), 10, 64)

	log.Println("Dataset saved successfully! Amount of rows =", countCT.String())

	return countCT.String(), nil
}

func ReadCSVPreview(path string, limit int) ([]Dataset, error) {
	file, err := os.Open(path)
	if err != nil {
		log.Println("Error opening file", err)
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// читаем header
	_, err = reader.Read()
	if err != nil {
		log.Println("Error reading file", err)
		return nil, fmt.Errorf("cannot read csv header: %w", err)
	}

	var result []Dataset

	for i := 0; i < limit; i++ {
		record, err := reader.Read()
		if err != nil {
			break
		}

		row, err := parseDatasetRow(record)
		if err != nil {
			log.Println("Error parsing dataset row", err)
			return nil, err
		}

		result = append(result, row)
	}

	return result, nil
}

func parseDatasetRow(r []string) (Dataset, error) {
	if len(r) < 16 {
		return Dataset{}, fmt.Errorf("invalid csv row length")
	}

	ts, _ := strconv.ParseInt(r[1], 10, 64)

	parse := func(s string) float32 {
		v, _ := strconv.ParseFloat(s, 32)
		return float32(v)
	}

	return Dataset{
		Id:        r[0],
		Timestamp: ts,
		Type:      r[2],

		AccelX: parse(r[3]),
		AccelY: parse(r[4]),
		AccelZ: parse(r[5]),

		GyroX: parse(r[6]),
		GyroY: parse(r[7]),
		GyroZ: parse(r[8]),

		QX: parse(r[9]),
		QY: parse(r[10]),
		QZ: parse(r[11]),
		QW: parse(r[12]),

		Yaw:   parse(r[13]),
		Pitch: parse(r[14]),
		Roll:  parse(r[15]),
	}, nil
}

func GetParams() (Params, error) {
	log.Println("Getting params...")

	sqlGetParamsId, _ := db.LoadQuery("get_ids.sql")
	sqlGetParamsType, _ := db.LoadQuery("get_types.sql")

	rows1, err := db.DB.Query(
		context.Background(),
		sqlGetParamsId)
	if err != nil {
		log.Println("Error getting params", err)
		return Params{}, fmt.Errorf("failed to get params: %v", err)
	}
	defer rows1.Close()

	rows2, err := db.DB.Query(
		context.Background(),
		sqlGetParamsType)
	if err != nil {
		log.Println("Error getting params", err)
		return Params{}, fmt.Errorf("failed to get params: %v", err)
	}
	defer rows2.Close()

	var ids []string
	var types []string

	for rows1.Next() {
		var id string

		if err := rows1.Scan(&id); err != nil {
			log.Println("Error parsing ids", err)
			return Params{}, err
		}

		ids = append(ids, id)
	}

	for rows2.Next() {
		var typee string

		if err := rows2.Scan(&typee); err != nil {
			log.Println("Error parsing types", err)
			return Params{}, err
		}

		types = append(types, typee)
	}

	log.Println("Got params successfully!")

	return Params{
		Id:   ids,
		Type: types,
	}, nil
}

func ChangeAble(id string, mode string) error {
	log.Println("Change able...")

	var sqlAbleUser string
	if mode == "enable" {
		sqlAbleUser, _ = db.LoadQuery("enable_user.sql")
	} else {
		sqlAbleUser, _ = db.LoadQuery("disable_user.sql")
	}

	_, err := db.DB.Exec(
		context.Background(),
		sqlAbleUser,
		id)

	if err != nil {
		return fmt.Errorf("failed to change user able: %v", err)
	}

	log.Println("User disable/enable successful")

	return nil
}
