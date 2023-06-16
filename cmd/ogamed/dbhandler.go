package main

import (
	"encoding/json"
	"log"
	"net/url"
	"time"

	v9 "github.com/alaingilbert/ogame/pkg/extractor/v9"
	"github.com/alaingilbert/ogame/pkg/wrapper"
)

func getHTMLInterceptor(bot *wrapper.OGame) func(method, url string, params, payload url.Values, pageHTML []byte) {
	myExtractor := bot.GetExtractor()

	universeName := bot.Universe
	language := bot.GetLanguage()
	playerID := bot.Player.PlayerID

	fn := func(method, url string, params, payload url.Values, pageHTML []byte) {
		var page string
		if params.Get("page") == "ingame" {
			page = params.Get("component")
		} else {
			page = params.Get("page")
		}

		var currentDBPlanet BotPlanet

		if wrapper.IsKnowFullPage(params) {
			currentCelestialID, err := myExtractor.ExtractPlanetID(pageHTML)
			if err != nil {
				return
			}

			db.Where(&BotPlanet{UniverseName: universeName, Language: language, PlayerID: playerID, CelestialID: int64(currentCelestialID)}).Find(&currentDBPlanet)
			currentDBPlanet.UniverseName = universeName
			currentDBPlanet.Language = language
			currentDBPlanet.PlayerID = playerID
			currentDBPlanet.CelestialID = int64(currentCelestialID)
			result := myExtractor.ExtractResourcesDetailsFromFullPage(pageHTML)

			currentDBPlanet.CelestialID = int64(currentCelestialID)
			currentDBPlanet.ResourcesDetails, err = json.Marshal(result)
			if err != nil {
				return
			}
			currentDBPlanet.Resources = result.Available()
			currentCoordinates, err := myExtractor.ExtractPlanetCoordinate(pageHTML)
			if err != nil {
				return
			}
			currentDBPlanet.Galaxy = currentCoordinates.Galaxy
			currentDBPlanet.System = currentCoordinates.System
			currentDBPlanet.Position = currentCoordinates.Position
			currentDBPlanet.CelestialType = int64(currentCoordinates.Type)

			currentDBPlanet.UniverseName = universeName
			currentDBPlanet.Language = language
			currentDBPlanet.PlayerID = playerID

			currentTime := time.Unix(myExtractor.ExtractOgameTimestamp(pageHTML), 0)

			switch page {
			case wrapper.OverviewPageName:
				// Extract all Constructions
				buildingID, buildingCountDown, _, _, lfBuildingID, lfBuildingCountDown, lfResearchID, lfResearchCountDown := myExtractor.ExtractConstructions(pageHTML)

				if buildingID != 0 {
					tmp := int64(buildingID)
					currentDBPlanet.ConstructionBuildingID = &tmp
					tmp2 := currentTime.Add(time.Duration(buildingCountDown) * time.Second)
					currentDBPlanet.ConstructionFinishedAt = &tmp2
				} else {
					currentDBPlanet.ConstructionBuildingID = nil
					currentDBPlanet.ConstructionFinishedAt = nil
				}

				if lfBuildingID != 0 {
					tmp := int64(lfBuildingID)
					currentDBPlanet.LfConstructionBuildingID = &tmp
					tmp2 := currentTime.Add(time.Duration(lfBuildingCountDown) * time.Second)
					currentDBPlanet.LfBuildingFinishedAt = &tmp2
				} else {
					currentDBPlanet.LfConstructionBuildingID = nil
					currentDBPlanet.LfBuildingFinishedAt = nil
				}

				if lfResearchID != 0 {
					tmp := int64(lfResearchID)
					currentDBPlanet.LfResearchID = &tmp
					tmp2 := currentTime.Add(time.Duration(lfResearchCountDown) * time.Second)
					currentDBPlanet.LfResearchFinishedAt = &tmp2
				} else {
					currentDBPlanet.LfResearchID = nil
					currentDBPlanet.LfResearchFinishedAt = nil
				}
				// pageX, _ := parser.ParsePage[parser.OverviewPage](myExtractor, pageHTML)
				// productionQf, productionCountDown, err := pageX.ExtractOverviewProduction()

				ext := v9.NewExtractor()
				ext.SetLanguage(bot.GetServer().Language)
				ext.SetLifeformEnabled(bot.GetExtractor().GetLifeformEnabled())
				//productionQf, productionCountDown, err := v9.NewExtractor().ExtractOverviewProduction(pageHTML)
				productionQf, productionCountDown, err := ext.ExtractOverviewProduction(pageHTML)
				//productionQf, productionCountDown, err := myExtractor.ExtractOverviewProduction(pageHTML)
				log.Println(productionQf)

				if err != nil {
					return
				}
				if productionCountDown != 0 {
					tmp2 := currentTime.Add(time.Duration(productionCountDown) * time.Second)
					currentDBPlanet.ProductionFinishedAt = &tmp2
					log.Println(productionQf)
					tmp1, err := json.Marshal(productionQf)
					if err != nil {
						log.Fatal(err)
					}
					currentDBPlanet.ProductionQueue = tmp1

				} else {
					currentDBPlanet.ProductionFinishedAt = nil
					currentDBPlanet.ProductionQueue = nil
				}

			case wrapper.SuppliesPageName:
				result, err := myExtractor.ExtractResourcesBuildings(pageHTML)
				if err != nil {
					return
				}
				currentDBPlanet.ResourcesBuildings = result

				buildingID, buildingCountDown, _, _, _, _, _, _ := myExtractor.ExtractConstructions(pageHTML)
				if buildingID != 0 {
					tmp := int64(buildingID)
					currentDBPlanet.ConstructionBuildingID = &tmp
					tmp2 := currentTime.Add(time.Duration(buildingCountDown) * time.Second)
					currentDBPlanet.ConstructionFinishedAt = &tmp2
				} else {
					currentDBPlanet.ConstructionBuildingID = nil
					currentDBPlanet.ConstructionFinishedAt = nil
				}
			case wrapper.FacilitiesPageName:
				result, err := myExtractor.ExtractFacilities(pageHTML)
				if err != nil {
					return
				}
				currentDBPlanet.Facilities = result

				buildingID, buildingCountDown, _, _, _, _, _, _ := myExtractor.ExtractConstructions(pageHTML)
				if buildingID != 0 {
					tmp := int64(buildingID)
					currentDBPlanet.ConstructionBuildingID = &tmp
					tmp2 := currentTime.Add(time.Duration(buildingCountDown) * time.Second)
					currentDBPlanet.ConstructionFinishedAt = &tmp2
				} else {
					currentDBPlanet.ConstructionBuildingID = nil
					currentDBPlanet.ConstructionFinishedAt = nil
				}
			case wrapper.ShipyardPageName:
				// pageX, _ := parser.ParsePage[parser.OverviewPage](myExtractor, pageHTML)
				// productionQf, productionCountDown, err := pageX.ExtractOverviewProduction()

				ext := v9.NewExtractor()
				ext.SetLanguage(bot.GetServer().Language)
				ext.SetLifeformEnabled(bot.GetExtractor().GetLifeformEnabled())
				//productionQf, productionCountDown, err := v9.NewExtractor().ExtractOverviewProduction(pageHTML)
				productionQf, productionCountDown, err := ext.ExtractOverviewProduction(pageHTML)
				//productionQf, productionCountDown, err := myExtractor.ExtractOverviewProduction(pageHTML)
				log.Println(productionQf)

				if err != nil {
					return
				}
				if productionCountDown != 0 {
					tmp2 := currentTime.Add(time.Duration(productionCountDown) * time.Second)
					currentDBPlanet.ProductionFinishedAt = &tmp2
					log.Println(productionQf)
					tmp1, err := json.Marshal(productionQf)
					if err != nil {
						log.Fatal(err)
					}
					currentDBPlanet.ProductionQueue = tmp1

				} else {
					currentDBPlanet.ProductionFinishedAt = nil
					currentDBPlanet.ProductionQueue = nil
				}
			}
			db.Save(&currentDBPlanet)
		}

		//case wrapper.OverviewPageName, wrapper.FetchResourcesPageName, wrapper.FacilitiesPageName, wrapper.ResearchPageName, wrapper.TraderOverviewPageName, wrapper.ShipyardPageName, wrapper.DefensesPageName, wrapper.FleetdispatchPageName, wrapper.MovementPageName, wrapper.GalaxyPageName, wrapper.AlliancePageName, wrapper.PremiumPageName, wrapper.ShopPageName, wrapper.MessagesPageName, wrapper.ChatPageName, wrapper.ResourceSettingsPageName, "characterclassselection", "lfsettings", wrapper.BuddiesPageName, wrapper.PreferencesPageName, wrapper.HighscorePageName:

	}

	return fn
}
