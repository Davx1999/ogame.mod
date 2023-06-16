package main

import (
	"time"

	"github.com/alaingilbert/ogame/pkg/ogame"
	"gorm.io/gorm"
)

type BotPlanet struct {
	gorm.Model
	UniverseName             string
	Language                 string
	PlayerID                 int64
	CelestialID              int64
	Name                     string
	Diameter                 int64
	Galaxy                   int64
	System                   int64
	Position                 int64
	CelestialType            int64
	FieldsBuilt              int64
	FieldsTotal              int64
	TemperatureMin           int64
	TemperatureMax           int64
	HasMoon                  bool
	ResourcesDetails         []byte
	ogame.Resources          `gorm:"embedded"`
	ogame.ResourcesBuildings `gorm:"embedded"`
	ogame.Facilities         `gorm:"embedded"`
	ogame.ShipsInfos         `gorm:"embedded"`
	ogame.DefensesInfos      `gorm:"embedded"`
	ogame.LfBuildings        `gorm:"embedded"`
	ogame.LfResearches       `gorm:"embedded"`
	ConstructionFinishedAt   *time.Time
	ConstructionBuildingID   *int64
	ProductionQueue          []byte
	ProductionFinishedAt     *time.Time
	LfBuildingFinishedAt     *time.Time
	LfConstructionBuildingID *int64
	LfResearchFinishedAt     *time.Time
	LfResearchID             *int64
}
