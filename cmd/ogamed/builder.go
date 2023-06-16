package main

import (
	"log"
	"math"
	"time"

	"github.com/alaingilbert/ogame/pkg/ogame"
	"github.com/alaingilbert/ogame/pkg/utils"
	"gorm.io/gorm"
)

type BuildQueue struct {
	gorm.Model
	ServerName     string
	ServerLanguage string
	PlayerID       int64
	CelestialID    int64 `json:"planetid"`
	OGameID        int64 `json:"id"`
	Nbr            int64 `json:"nbr"`
	Idx            int64 `json:"idx"`
}

func addQueue(serverName, serverLanguage string, playerID, celestialID, oGameID, nbr int64) {
	b := BuildQueue{
		ServerName:     serverName,
		ServerLanguage: serverLanguage,
		PlayerID:       playerID,
		CelestialID:    celestialID,
		OGameID:        oGameID,
		Nbr:            nbr,
	}
	log.Println(&b) //b db.Create(&b)

	if db == nil {
		panic("db is nil")
	}
	db.Create(&b)
}

func getQueue(serverName, serverLanguage string, playerID int64) []BuildQueue {
	var b []BuildQueue
	db.Where(&BuildQueue{ServerName: serverName, ServerLanguage: serverLanguage, PlayerID: playerID}).Find(&b)
	return b
}

type missingResources struct {
	ResourcesNeeded ogame.Resources
	Countdown       int64
	ResourceName    string
}

func ResourceCountdown(price ogame.Resources, res ogame.ResourcesDetails) missingResources {
	var result missingResources
	if res.Available().CanAfford(price) {
		return result
	}

	metallNeeded := res.Available().Metal - price.Metal
	crystalNeeded := res.Available().Crystal - price.Crystal
	deuteriumNeeded := res.Available().Deuterium - price.Deuterium

	if metallNeeded < 0 {
		metallNeeded = utils.AbsInt64(metallNeeded)
	}

	if crystalNeeded < 0 {
		crystalNeeded = utils.AbsInt64(crystalNeeded)
	}

	if deuteriumNeeded < 0 {
		deuteriumNeeded = utils.AbsInt64(deuteriumNeeded)
	}

	var metalCountdown time.Duration
	var crystalCountdown time.Duration
	var deuteriumCountdown time.Duration

	if res.Metal.CurrentProduction != 0 {
		metalCountdown = time.Duration(int64(math.Ceil(float64(metallNeeded/res.Metal.CurrentProduction)))) * time.Hour
	}

	if res.Crystal.CurrentProduction != 0 {
		crystalCountdown = time.Duration(int64(math.Ceil(float64(crystalNeeded/res.Crystal.CurrentProduction)))) * time.Hour
	}

	if res.Deuterium.CurrentProduction != 0 {
		deuteriumCountdown = time.Duration(int64(math.Ceil(float64(deuteriumNeeded/res.Deuterium.CurrentProduction)))) * time.Hour
	}

	var maxCountdown int64
	var maxResName string

	if metalCountdown > crystalCountdown {
		maxCountdown = int64(metalCountdown.Seconds())
		maxResName = "Metal"
	} else {
		maxCountdown = int64(crystalCountdown.Seconds())
		maxResName = "Crystal"
	}

	if deuteriumCountdown > time.Duration(maxCountdown)*time.Second {
		maxCountdown = int64(deuteriumCountdown.Seconds())
		maxResName = "Deuterium"
	}

	result.ResourcesNeeded = price.Sub(res.Available())
	result.Countdown = maxCountdown
	result.ResourceName = maxResName

	return result
}

// func ResourceCountdown(price ogame.Resources, res ogame.ResourcesDetails) int64 {
// 	if res.Available().CanAfford(price) {
// 		return 0
// 	}

// 	metallNeeded := res.Available().Metal - price.Metal
// 	crystalNeeded := res.Available().Crystal - price.Crystal
// 	deuteriumNeeded := res.Available().Deuterium - price.Deuterium

// 	if metallNeeded < 0 {
// 		metallNeeded = utils.AbsInt64(metallNeeded)
// 	}

// 	if crystalNeeded < 0 {
// 		crystalNeeded = utils.AbsInt64(crystalNeeded)
// 	}

// 	if deuteriumNeeded < 0 {
// 		deuteriumNeeded = utils.AbsInt64(deuteriumNeeded)
// 	}

// 	var metalCountdown time.Duration
// 	var crystalCountdown time.Duration
// 	var deuteriumCountdown time.Duration

// 	if res.Metal.CurrentProduction != 0 {
// 		metalCountdown = time.Duration(int64(math.Ceil(float64(metallNeeded/res.Metal.CurrentProduction)))) * time.Hour
// 	}

// 	if res.Crystal.CurrentProduction != 0 {
// 		crystalCountdown = time.Duration(int64(math.Ceil(float64(crystalNeeded/res.Crystal.CurrentProduction)))) * time.Hour
// 	}

// 	if res.Deuterium.CurrentProduction != 0 {
// 		deuteriumCountdown = time.Duration(int64(math.Ceil(float64(deuteriumNeeded/res.Deuterium.CurrentProduction)))) * time.Hour
// 	}

// 	var maxCountdown int64
// 	//var maxResName string

// 	if metalCountdown > crystalCountdown {
// 		maxCountdown = int64(metalCountdown.Seconds())
// 		//maxResName = "Metal"
// 	} else {
// 		maxCountdown = int64(crystalCountdown.Seconds())
// 		//maxResName = "Crystal"
// 	}

// 	if deuteriumCountdown > time.Duration(maxCountdown)*time.Second {
// 		maxCountdown = int64(deuteriumCountdown.Seconds())
// 		//maxResName = "Deuterium"
// 	}

// 	return maxCountdown
// }

// func ResourceNeeded(price ogame.Resources, res ogame.ResourcesDetails) string {
// 	if res.Available().CanAfford(price) {
// 		return ""
// 	}

// 	metallNeeded := res.Available().Metal - price.Metal
// 	crystalNeeded := res.Available().Crystal - price.Crystal
// 	deuteriumNeeded := res.Available().Deuterium - price.Deuterium

// 	if metallNeeded < 0 {
// 		metallNeeded = utils.AbsInt64(metallNeeded)
// 	}

// 	if crystalNeeded < 0 {
// 		crystalNeeded = utils.AbsInt64(crystalNeeded)
// 	}

// 	if deuteriumNeeded < 0 {
// 		deuteriumNeeded = utils.AbsInt64(deuteriumNeeded)
// 	}

// 	var metalCountdown time.Duration
// 	var crystalCountdown time.Duration
// 	var deuteriumCountdown time.Duration

// 	if res.Metal.CurrentProduction != 0 {
// 		metalCountdown = time.Duration(int64(math.Ceil(float64(metallNeeded/res.Metal.CurrentProduction)))) * time.Hour
// 	}

// 	if res.Crystal.CurrentProduction != 0 {
// 		crystalCountdown = time.Duration(int64(math.Ceil(float64(crystalNeeded/res.Crystal.CurrentProduction)))) * time.Hour
// 	}

// 	if res.Deuterium.CurrentProduction != 0 {
// 		deuteriumCountdown = time.Duration(int64(math.Ceil(float64(deuteriumNeeded/res.Deuterium.CurrentProduction)))) * time.Hour
// 	}

// 	var maxCountdown int64
// 	var maxResName string

// 	if metalCountdown > crystalCountdown {
// 		maxCountdown = int64(metalCountdown.Seconds())
// 		maxResName = "Metal"
// 	} else {
// 		maxCountdown = int64(crystalCountdown.Seconds())
// 		maxResName = "Crystal"
// 	}

// 	if deuteriumCountdown > time.Duration(maxCountdown)*time.Second {
// 		maxCountdown = int64(deuteriumCountdown.Seconds())
// 		maxResName = "Deuterium"
// 	}

// 	return maxResName
// }
