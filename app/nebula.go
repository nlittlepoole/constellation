package main

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/nlittlepoole/observatory/rover"
	"path/filepath"
	"time"
)

var db *gorm.DB
var DATE_FORMATS = map[time.Duration]string{
	time.Hour:   "%Y-%m-%d %H",
	time.Minute: "%Y-%m-%d %H:%M",
}

const DB_FILENAME string = "nebula.db"

func init() {
	var err error
	dbPath := filepath.Join(getCachePath(), DB_FILENAME)
	if db, err = gorm.Open("sqlite3", dbPath); err != nil {
		log.Error(err)
		panic(err)
	}
	db.AutoMigrate(&Event{})
}

func getDateFormatString(granularity time.Duration) string {
	if format, ok := DATE_FORMATS[granularity]; ok {
		return format
	} else {
		return "%Y-%m-%d"
	}
}

type Timeseries struct {
	X []string `json:"x"`
	Y []int64  `json:"y"`
}

type Event struct {
	Timestamp     time.Time `gorm:primary_key`
	DeviceId      string    `gorm:primary_key`
	RadioStrength int64
	Location      string
}

func anonymize(address string) string {
	h := sha1.New()
	h.Write([]byte(address))
	return hex.EncodeToString(h.Sum(nil))
}

func logEvent(probe rover.Probe) (err error) {
	event := Event{
		Timestamp:     probe.Timestamp,
		DeviceId:      anonymize(probe.Address),
		RadioStrength: probe.Strength,
		Location:      probe.Location,
	}
	if ok := db.NewRecord(event); ok {
		db.Create(&event)
	} else {
		err = errors.New("Failed to write event: event already exists")
	}
	return
}

func GetUniques(start time.Time, end time.Time, granularity time.Duration) (Timeseries, error) {
	x := make([]string, 0)
	y := make([]int64, 0)
	rows, err := db.Raw(
		`
     SELECT
	cast(strftime(?, timestamp) as text) as Timestamp,
	count(distinct device_id) as Uniques
     FROM events
     WHERE
	timestamp >= ? AND
	timestamp <= ?
     GROUP BY 1 ORDER BY 1 ASC
     `,
		DATE_FORMATS[granularity],
		start,
		end,
	).Rows()
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var key string
			var val int64
			err = rows.Scan(&key, &val)
			if err == nil {
				x = append(x, key)
				y = append(y, val)
			}
		}
	}
	return Timeseries{x, y}, err
}

func GetAllUniques(granularity time.Duration) (Timeseries, error) {
	return GetUniques(time.Unix(0, 0), time.Now(), granularity)
}

func GetCurrentUniques(window time.Duration) (result int64, err error) {
	series, err := GetUniques(time.Now().Add(-1*window), time.Now(), window)
	total := int64(0)
	for _, point := range series.Y {
		total += point
	}
	return total, err
}

type Returning struct {
	Old int64 `json:"old"`
	New int64 `json:"new"`
}

func GetReturningUniques(start time.Time, end time.Time) (Returning, error) {
	diff := end.Sub(start)
	previous := start.Add(-1 * diff)
	row := db.Raw(
		`
     SELECT
	SUM(CASE WHEN y.device_id is not null THEN 1 ELSE 0 END) as old,
	SUM(CASE WHEN y.device_id is null THEN 1 ELSE 0 END) as new	
     FROM
	(SELECT
	    distinct device_id
        FROM events
        WHERE
	   timestamp >= ? AND
	   timestamp <= ?
	) x LEFT JOIN
	(SELECT
	    distinct device_id
        FROM events
        WHERE
	   timestamp >= ? AND
	   timestamp <= ?
	) y ON x.device_id = y.device_id
     `,
		start,
		end,
		previous,
		start,
	).Row()
	result := Returning{}
	err := row.Scan(&result.Old, &result.New)
	return result, err
}

func GetStrengthHistogram(start time.Time, end time.Time) (Timeseries, error) {
	x := make([]string, 0)
	y := make([]int64, 0)
	rows, err := db.Raw(
		`
     SELECT
	cast(10 * (radio_strength / 10) as text)  as dbm,
	count(distinct device_id) as uniques
     FROM events
     WHERE
	timestamp >= ? AND
	timestamp <= ?
     GROUP BY 1 ORDER BY 1 ASC
     `,
		start,
		end,
	).Rows()
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var key string
			var val int64
			err = rows.Scan(&key, &val)
			if err == nil {
				x = append(x, key)
				y = append(y, val)
			}
		}
	}
	return Timeseries{x, y}, err
}
