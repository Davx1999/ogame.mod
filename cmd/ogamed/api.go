package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/alaingilbert/ogame/pkg/ogame"
	"github.com/alaingilbert/ogame/pkg/utils"
	"github.com/alaingilbert/ogame/pkg/wrapper"
	"github.com/labstack/echo/v4"
)

func APIFlightTime(c echo.Context) error {
	bot := c.Get("bot").(*wrapper.OGame)

	formdata, _ := c.MultipartForm()
	//log.Println(formdata.Value)

	var ships ogame.ShipsInfos
	var origin int64
	var destination ogame.Coordinate
	var speed ogame.Speed
	var mission int64

	mission, _ = strconv.ParseInt("mission", 10, 64)

	for k, v := range formdata.Value {
		//log.Printf("%s, %s", k, v[0])
		switch k {
		case "LightFighter":
			nbr, _ := strconv.ParseInt(v[0], 10, 64)
			ships.LightFighter = nbr
			break
		case "HeavyFighter":
			nbr, _ := strconv.ParseInt(v[0], 10, 64)
			ships.HeavyFighter = nbr
			break
		case "Cruiser":
			nbr, _ := strconv.ParseInt(v[0], 10, 64)
			ships.Cruiser = nbr
			break
		case "Battleship":
			nbr, _ := strconv.ParseInt(v[0], 10, 64)
			ships.Battleship = nbr
			break
		case "Battlecruiser":
			nbr, _ := strconv.ParseInt(v[0], 10, 64)
			ships.Battlecruiser = nbr
			break
		case "Bomber":
			nbr, _ := strconv.ParseInt(v[0], 10, 64)
			ships.Bomber = nbr
			break
		case "Destroyer":
			nbr, _ := strconv.ParseInt(v[0], 10, 64)
			ships.Destroyer = nbr
			break
		case "Deathstar":
			nbr, _ := strconv.ParseInt(v[0], 10, 64)
			ships.Deathstar = nbr
			break
		case "SmallCargo":
			nbr, _ := strconv.ParseInt(v[0], 10, 64)
			ships.SmallCargo = nbr
			break
		case "LargeCargo":
			nbr, _ := strconv.ParseInt(v[0], 10, 64)
			ships.LargeCargo = nbr
			break
		case "ColonyShip":
			nbr, _ := strconv.ParseInt(v[0], 10, 64)
			ships.ColonyShip = nbr
			break
		case "Recycler":
			nbr, _ := strconv.ParseInt(v[0], 10, 64)
			ships.Recycler = nbr
			break
		case "EspionageProbe":
			nbr, _ := strconv.ParseInt(v[0], 10, 64)
			ships.EspionageProbe = nbr
			break
		case "Crawler":
			nbr, _ := strconv.ParseInt(v[0], 10, 64)
			ships.Crawler = nbr
			break
		case "Reaper":
			nbr, _ := strconv.ParseInt(v[0], 10, 64)
			ships.Reaper = nbr
			break
		case "Pathfinder":
			nbr, _ := strconv.ParseInt(v[0], 10, 64)
			ships.Pathfinder = nbr
			break
		case "origin":
			nbr, _ := strconv.ParseInt(v[0], 10, 64)
			origin = nbr
			break
		case "speed":
			nbr, _ := strconv.ParseInt(v[0], 10, 64)
			speed = ogame.Speed(nbr)
			break
		case "destGalaxy":
			nbr, _ := strconv.ParseInt(v[0], 10, 64)
			destination.Galaxy = nbr
			break
		case "destSystem":
			nbr, _ := strconv.ParseInt(v[0], 10, 64)
			destination.System = nbr
			break
		case "destPosition":
			nbr, _ := strconv.ParseInt(v[0], 10, 64)
			destination.Position = nbr
			break
		}
	}

	cel := bot.GetCachedCelestial(origin)
	secs, fuel := bot.FlightTime(cel.GetCoordinate(), destination, speed, ships, ogame.MissionID(mission))
	cargo := ships.Cargo(bot.GetCachedResearch(), false, bot.CharacterClass() == ogame.Collector, bot.IsPioneers())

	human_time := time.Duration(secs) * time.Second
	data := struct {
		Secs       int64  `json:"secs"`
		Fuel       int64  `json:"fuel"`
		Cargo      int64  `json:"cargo"`
		Human_time string `json:"human_time"`
	}{
		Secs:       secs,
		Fuel:       fuel,
		Cargo:      cargo,
		Human_time: human_time.String(),
	}
	return c.JSON(http.StatusOK, data)
}

func APIResources(c echo.Context) error {
	bot := c.Get("bot").(*wrapper.OGame)
	planetID, err := strconv.ParseInt(c.Param("celestialID"), 10, 64)
	if err != nil || planetID < 0 {
		return c.JSON(http.StatusBadRequest, wrapper.ErrorResp(400, "invalid planet id"))
	}
	//res, _ := bot.GetResources(ogame.CelestialID(planetID))
	res, _ := bot.GetResourcesDetails(ogame.CelestialID(planetID))
	data := struct {
		Resources ogame.Resources `json:"resources"`
	}{
		Resources: res.Available(),
	}

	return c.JSON(http.StatusOK, data)
}

func APIShips(c echo.Context) error {
	bot := c.Get("bot").(*wrapper.OGame)
	planetID, err := strconv.ParseInt(c.Param("celestialID"), 10, 64)
	if err != nil || planetID < 0 {
		return c.JSON(http.StatusBadRequest, wrapper.ErrorResp(400, "invalid planet id"))
	}
	//res, _ := bot.GetShips(ogame.CelestialID(planetID))
	res, _ := bot.GetShips(ogame.CelestialID(planetID))

	data := struct {
		Ships ogame.ShipsInfos `json:"ships"`
	}{
		Ships: res,
	}
	return c.JSON(http.StatusOK, data)
}

func APIFleets(c echo.Context) error {
	bot := c.Get("bot").(*wrapper.OGame)
	fleets, slots := bot.GetFleets()
	for k, f := range fleets {
		if !f.ReturnFlight {
			arrivalIn := f.ArrivalTime.Sub(time.Now())
			fleets[k].ArriveIn = int64(arrivalIn.Seconds())
		}

		if f.ReturnFlight {
			backIn := f.BackTime.Sub(time.Now())
			fleets[k].BackIn = int64(backIn.Seconds())
		}
	}
	data := struct {
		Fleets []ogame.Fleet `json:"fleets"`
		Slots  ogame.Slots   `json:"slots"`
	}{
		Fleets: fleets,
		Slots:  slots,
	}

	return c.JSON(http.StatusOK, data)
}

func APISendFleets(c echo.Context) error {
	//log.Println(c.Request().Method)
	bot := c.Get("bot").(*wrapper.OGame)

	params, _ := c.FormParams()
	var ships ogame.ShipsInfos

	duration, _ := strconv.ParseInt(params.Get("duration"), 10, 64)

	ships.LightFighter, _ = strconv.ParseInt(params.Get("LightFighter"), 10, 64)
	ships.HeavyFighter, _ = strconv.ParseInt(params.Get("HeavyFighter"), 10, 64)
	ships.Cruiser, _ = strconv.ParseInt(params.Get("Cruiser"), 10, 64)
	ships.Battleship, _ = strconv.ParseInt(params.Get("Battleship"), 10, 64)
	ships.Battlecruiser, _ = strconv.ParseInt(params.Get("Battlecruiser"), 10, 64)
	ships.Bomber, _ = strconv.ParseInt(params.Get("Bomber"), 10, 64)
	ships.Destroyer, _ = strconv.ParseInt(params.Get("Destroyer"), 10, 64)
	ships.Deathstar, _ = strconv.ParseInt(params.Get("Deathstar"), 10, 64)
	ships.SmallCargo, _ = strconv.ParseInt(params.Get("SmallCargo"), 10, 64)
	ships.LargeCargo, _ = strconv.ParseInt(params.Get("LargeCargo"), 10, 64)
	ships.ColonyShip, _ = strconv.ParseInt(params.Get("ColonyShip"), 10, 64)
	ships.Recycler, _ = strconv.ParseInt(params.Get("Recycler"), 10, 64)
	ships.EspionageProbe, _ = strconv.ParseInt(params.Get("EspionageProbe"), 10, 64)
	ships.Crawler, _ = strconv.ParseInt(params.Get("Crawler"), 10, 64)
	ships.Reaper, _ = strconv.ParseInt(params.Get("Reaper"), 10, 64)
	ships.Pathfinder, _ = strconv.ParseInt(params.Get("Pathfinder"), 10, 64)

	o, _ := strconv.ParseInt(params.Get("origin"), 10, 64)

	var dest ogame.Coordinate
	dest.Galaxy, _ = strconv.ParseInt(params.Get("destGalaxy"), 10, 64)
	dest.System, _ = strconv.ParseInt(params.Get("destSystem"), 10, 64)
	dest.Position, _ = strconv.ParseInt(params.Get("destPosition"), 10, 64)
	pT := params.Get("planetType")
	pType, _ := strconv.ParseInt(pT, 10, 64)
	dest.Type = ogame.CelestialType(pType)

	speed, _ := strconv.ParseInt(params.Get("speed"), 10, 64)

	mission, _ := strconv.ParseInt(params.Get("mission"), 10, 64)

	metal, _ := strconv.ParseInt(params.Get("metal"), 10, 64)
	crystal, _ := strconv.ParseInt(params.Get("crystal"), 10, 64)
	deuterium, _ := strconv.ParseInt(params.Get("deuterium"), 10, 64)

	res := ogame.Resources{Metal: metal, Crystal: crystal, Deuterium: deuterium}

	//log.Printf("Origin: %d %d %d %s %s", o, speed, mission, dest.String(), ships.String())
	//log.Println("%d", pType)

	allShips, _ := strconv.ParseBool(params.Get("allShips"))
	allResources, _ := strconv.ParseBool(params.Get("allResources"))

	_, scheduleBtn := params["scheduleBtn"]
	if scheduleBtn {
		var sdb ScheduleFleet
		log.Println("Schedule")
		hour := params.Get("schedule_hour")
		min := params.Get("schedule_min")
		sec := params.Get("schedule_sec")
		tNext := utils.GetNextExecutionTime(time.Now(), hour+":"+min+":"+sec)
		sdb.AllShips = allShips
		sdb.AllResources = allResources
		sdb.SendAt = tNext
		sdb.ShipsInfos = ships
		sdb.Resources = res
		sdb.Universe = bot.Universe
		sdb.Language = bot.GetServer().Language
		sdb.PlayerID = bot.Player.PlayerID
		db.Save(&sdb)
	}

	_, sendBtn := params["sendBtn"]
	if sendBtn {
		fleetsender := wrapper.NewFleetBuilder(bot)
		fleetsender.SetDuration(duration)
		fleetsender.SetOrigin(ogame.CelestialID(o))
		fleetsender.SetDestination(dest)
		fleetsender.SetMission(ogame.MissionID(mission))
		fleetsender.SetSpeed(ogame.Speed(speed))
		if allResources {
			fleetsender.SetAllResources()
		} else {
			fleetsender.SetResources(res)
		}
		if allShips {
			fleetsender.SetAllShips()
		} else {
			fleetsender.SetShips(ships)
		}

		_, err := fleetsender.SendNow()
		if err != nil {
			log.Println(err.Error())
		}
		bot.GetShips(ogame.CelestialID(o))
	}
	return c.Redirect(http.StatusFound, "/bot/flights")
}

// CancelFleetHandler ...
func APIFlightsCancel(c echo.Context) error {
	bot := c.Get("bot").(*wrapper.OGame)
	fleetID, err := utils.ParseI64(c.Param("fleetID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, wrapper.ErrorResp(400, err.Error()))
	}
	bot.CancelFleet(ogame.FleetID(fleetID))
	return c.Redirect(http.StatusFound, "/bot/flights")
}
