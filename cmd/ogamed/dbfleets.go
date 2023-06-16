package main

import (
	"time"

	"github.com/alaingilbert/ogame/pkg/ogame"
	"gorm.io/gorm"
)

type ScheduleFleet struct {
	gorm.Model
	Universe            string
	Language            string
	PlayerID            int64
	Origin              int64
	DestinationGalaxy   int64
	DestinationSystem   int64
	DestinationPosition int64
	DestinationType     int64
	Speed               int64
	Mission             int64
	ogame.Resources     `gorm:"embedded"`
	ogame.ShipsInfos    `gorm:"embedded"`
	AllShips            bool
	AllResources        bool
	SendAt              time.Time
}
