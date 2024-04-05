package models

import (
	"time"
)

type RouteLogWithOneTime struct {
	UserID    int
	RouteID   int
	TimeStamp time.Time
}

type RouteLog struct {
	UserID    int
	RouteId   int
	StartTime time.Time
	EndTime   time.Time
}
