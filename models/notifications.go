package models

import (
	"fmt"
	"time"
)

type Status string

const (
	Pending  Status = "Pending"
	Sent     Status = "Sent"
	Rejected Status = "Rejected"
)

type Notification struct {
	Id           int64
	UserEmail    string
	LocationName string
	WaveHeightM  float64
	Status       Status
	ForecastTime time.Time
	Created      time.Time
	Scheduled    time.Time
	Sent         time.Time
}

func NewNotification(email, locName string, wh float64, ts int64) *Notification {
	fcTime := time.Unix(0, ts)
	sched := fcTime.AddDate(0, 0, -1)
	if sched.Before(time.Now()) {
		sched = time.Now().Add(time.Hour)
	}
	not := &Notification{
		UserEmail:    email,
		LocationName: locName,
		Scheduled:    sched,
		ForecastTime: fcTime,
		WaveHeightM:  wh}
	return not
}

func (db *DB) StoreNotification(not *Notification) (bool, error) {
	query := `INSERT INTO notification as n (user_email, location_name, scheduled, forecast_time)
			  SELECT CAST($1 AS VARCHAR), CAST($2 AS VARCHAR), CAST($3 AS TIMESTAMP), $4
			  WHERE NOT EXISTS(
			      SELECT * FROM notification
			          WHERE user_email = $1
			          AND status = 'Pending'
			          AND scheduled::date = $3::date
			  );`
	stmt, err := db.Prepare(query)
	if err != nil {
		return false, err
	}
	res, err := stmt.Exec(not.UserEmail, not.LocationName, not.Scheduled, not.ForecastTime)
	if err != nil {
		return false, err
	}
	change, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	return change > 0, err
}
