package main

import (
	"context"
	"encoding/json"
	"log"
	"math"
	"time"

	"github.com/alaingilbert/ogame/pkg/ogame"
	"github.com/alaingilbert/ogame/pkg/wrapper"
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

func getQueueForCelestial(serverName, serverLanguage string, playerID int64, celestialID ogame.CelestialID) []BuildQueue {
	var b []BuildQueue
	db.Where(&BuildQueue{ServerName: serverName, ServerLanguage: serverLanguage, PlayerID: playerID, CelestialID: int64(celestialID)}).Find(&b)
	return b
}

func getCelestialFromDB(serverName, serverLanguage string, playerID int64, celestialID ogame.CelestialID) BotPlanet {
	var b BotPlanet
	result := db.Where(&BotPlanet{UniverseName: serverName, Language: serverLanguage, PlayerID: playerID, CelestialID: int64(celestialID)}).Find(&b)
	if result.Error != nil {
		log.Printf("Database Error %s", result.Error)
	}
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
		//metallNeeded = utils.AbsInt64(metallNeeded)
		metallNeeded = 0
	}

	if crystalNeeded < 0 {
		//crystalNeeded = utils.AbsInt64(crystalNeeded)
		crystalNeeded = 0
	}

	if deuteriumNeeded < 0 {
		//deuteriumNeeded = utils.AbsInt64(deuteriumNeeded)
		deuteriumNeeded = 0
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

type BotBrain struct {
	*wrapper.OGame
	ctx                   context.Context
	cancelFunc            context.CancelFunc
	celestialsCtx         map[ogame.CelestialID]context.Context
	celestialsCancelFuncs map[ogame.CelestialID]context.CancelFunc
	reRunCh               chan ogame.CelestialID
}

func NewBotBrain(bot *wrapper.OGame) (b *BotBrain) {
	b = &BotBrain{}
	b.OGame = bot
	b.celestialsCtx = map[ogame.CelestialID]context.Context{}
	b.celestialsCancelFuncs = map[ogame.CelestialID]context.CancelFunc{}
	b.ctx, b.cancelFunc = context.WithCancel(context.Background())
	b.reRunCh = make(chan ogame.CelestialID)
	return
}

func (brain *BotBrain) registerCelestial(celestialID ogame.CelestialID) {
	log.Printf("Register Celestial %d", celestialID)
	brain.celestialsCtx[celestialID], brain.celestialsCancelFuncs[celestialID] = context.WithCancel(brain.ctx)
}

func (brain *BotBrain) removeCelestial(celestialID ogame.CelestialID) {
	log.Printf("Remove Celestial %d", celestialID)
	delete(brain.celestialsCtx, celestialID)
	delete(brain.celestialsCancelFuncs, celestialID)
}

func (brain *BotBrain) ReRun(celestialID ogame.CelestialID) {
	log.Printf("re run brain logic for %d\n", celestialID)
	cancel, ex := brain.celestialsCancelFuncs[celestialID]
	if ex {
		cancel()
		brain.reRunCh <- celestialID
	}
}

var localBrain *BotBrain

func StartBrain(bot *wrapper.OGame) {
	log.Println("--- START BRAIN  ---")
	localBrain = NewBotBrain(bot)
	for {
		if bot.IsLoggedIn() && bot.Player.PlayerID != 0 {
			go localBrain.BuilderStart()
			break
		} else {
			time.Sleep(3 * time.Second)
		}
	}

}

func (brain *BotBrain) BuilderStart() {
	for {
		celestialWorkerFunc := func(ctx context.Context, celestialID ogame.CelestialID) {
			for {
				time.Sleep(1 * time.Second)

				queue := getQueueForCelestial(brain.Universe, brain.GetServer().Language, brain.Player.PlayerID, celestialID)

				var waitTime time.Duration
				var waitConstuctionBuildings time.Duration = 4 * time.Hour
				var waitResearches time.Duration = 4 * time.Hour
				var waitLfBuildings time.Duration = 4 * time.Hour

				if len(queue) == 0 {
					log.Println("Item not found!!!")
					waitTime = time.Duration(30 * time.Minute)
				} else {
					log.Println("Item found!!!")
					item := queue[0]
					var nextItem *BuildQueue
					if len(queue) > 1 {
						nextItem = &queue[1]
					}
					//log.Printf("U: %s, L: %s P: %d C:%d\n", brain.Universe, brain.GetServer().Language, brain.Player.PlayerID, celestialID)
					botPlanet := getCelestialFromDB(brain.Universe, brain.GetServer().Language, brain.Player.PlayerID, celestialID)

					build := func() {
						tech, err := brain.TechnologyDetails(celestialID, ogame.ID(item.OGameID))
						if err != nil {
							log.Println(err)
							return
						}

						if botPlanet.Resources.CanAfford(tech.Price) && tech.UpgradeEnabled {
							log.Printf("Start Upgrade: %s\n", ogame.ID(item.OGameID).String())
							brain.Build(celestialID, ogame.ID(item.OGameID), item.Nbr)
							db.Delete(&item)
							time.Sleep(3 * time.Second)
							if ogame.ID(item.OGameID).IsBuilding() {
								_, sec, _, _, _, _, _, _ := brain.ConstructionsBeingBuilt(celestialID)
								if nextItem != nil {
									if ogame.ID(nextItem.OGameID).IsBuilding() {
										waitTime = time.Duration(sec * int64(time.Second))
									}
								}
							}
							if ogame.ID(item.OGameID).IsTech() {
								_, _, _, sec, _, _, _, _ := brain.ConstructionsBeingBuilt(celestialID)
								if nextItem != nil {
									if ogame.ID(nextItem.OGameID).IsTech() {
										waitTime = time.Duration(sec * int64(time.Second))
									}
								}
							}
							if ogame.ID(item.OGameID).IsLfBuilding() {
								_, _, _, _, _, sec, _, _ := brain.ConstructionsBeingBuilt(celestialID)
								if nextItem != nil {
									if ogame.ID(nextItem.OGameID).IsLfBuilding() {
										waitTime = time.Duration(sec * int64(time.Second))
									}
								}
							}
						} else {
							var resDetails ogame.ResourcesDetails
							json.Unmarshal(botPlanet.ResourcesDetails, &resDetails)
							result := ResourceCountdown(tech.Price, resDetails)
							log.Printf("Price: %s | Resources: %s %s", tech.Price, resDetails.Available(), result)
							waitTime = time.Duration(result.Countdown) * time.Second
							if result.Countdown == 0 {
								waitTime = 15 * time.Minute
							}
						}

					}

					if ogame.ID(item.OGameID).IsBuilding() && botPlanet.ConstructionBuildingID == nil {
						build()
					} else if botPlanet.ConstructionFinishedAt != nil {
						waitTime = time.Until(*botPlanet.ConstructionFinishedAt)
						waitConstuctionBuildings = time.Until(*botPlanet.ConstructionFinishedAt)
					}

					if ogame.ID(item.OGameID).IsTech() && botPlanet.ResearchFinishedAt == nil {
						build()
					} else if botPlanet.ResearchFinishedAt != nil {
						waitTime = time.Until(*botPlanet.ResearchFinishedAt)
						waitResearches = time.Until(*botPlanet.ResearchFinishedAt)
					}

					if ogame.ID(item.OGameID).IsLfBuilding() && botPlanet.LfBuildingFinishedAt == nil {
						build()
					} else if botPlanet.LfBuildingFinishedAt != nil {
						waitTime = time.Until(*botPlanet.LfBuildingFinishedAt)
						waitLfBuildings = time.Until(*botPlanet.LfBuildingFinishedAt)
					}

				}

				log.Printf("wait for %s", waitTime)
				select {
				case <-ctx.Done():
					return

				case <-time.After(waitTime):
					log.Printf("Waited %s ", waitTime)
					time.Sleep(3 * time.Second)
					brain.ConstructionsBeingBuilt(celestialID)
				case <-time.After(waitConstuctionBuildings):
					log.Printf("Finished Construction")
					time.Sleep(3 * time.Second)
					brain.ConstructionsBeingBuilt(celestialID)
				case <-time.After(waitResearches):
					log.Printf("Finished Research")
					time.Sleep(3 * time.Second)
					brain.ConstructionsBeingBuilt(celestialID)
				case <-time.After(waitLfBuildings):
					log.Printf("Finished LfBuilding")
					time.Sleep(3 * time.Second)
					brain.ConstructionsBeingBuilt(celestialID)

				}
			}
		}

		cachedCelestials, err := brain.GetCelestials()
		if err != nil {
			time.Sleep(3 * time.Second)
		}
		for _, c := range cachedCelestials {
			if _, exists := brain.celestialsCancelFuncs[c.GetID()]; !exists {
				log.Println("Register Celestial Worker!")
				brain.registerCelestial(c.GetID())
				go celestialWorkerFunc(brain.celestialsCtx[c.GetID()], c.GetID())
			}
		}

		select {
		case celestialID := <-brain.reRunCh:
			time.Sleep(3 * time.Second)
			cancel := brain.celestialsCancelFuncs[celestialID]
			cancel()
			brain.removeCelestial(celestialID)
		case <-brain.ctx.Done():
			return
		}
	}
}
