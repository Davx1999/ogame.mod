package main

import (
	"html/template"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/alaingilbert/ogame/pkg/ogame"
	"github.com/alaingilbert/ogame/pkg/utils"
	"github.com/alaingilbert/ogame/pkg/wrapper"
	"github.com/labstack/echo/v4"
)

type MyTemplate struct {
	templates *template.Template
}

func (t *MyTemplate) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func GetEmpirePlanetHandler(c echo.Context) error {
	bot := c.Get("bot").(*wrapper.OGame)
	res := struct {
		Bot                 *wrapper.OGame
		EmpireCelestial     ogame.EmpireCelestial
		LifeformString      string
		LfBuildings         ogame.LfBuildings
		ResourcesDetails    ogame.ResourcesDetails
		Objs                map[ogame.ID]ogame.BaseOgameObj
		BuildingID          ogame.ID
		BuildingCountdown   int64
		ResearchID          ogame.ID
		ResearchCountdown   int64
		LfBuildingID        ogame.ID
		LfBuildingCountdown int64
		LfResearchID        ogame.ID
		LfResearchCountdown int64
	}{}

	res.Bot = bot
	res.Objs = ogame.Objs.GetAllObjs()
	planetID, err := utils.ParseI64(c.Param("planetID"))
	if err != nil {
		planetID = int64(bot.GetCachedCelestials()[0].GetID())
		//return c.JSON(http.StatusBadRequest, wrapper.ErrorResp(400, "invalid planet id"))
	}

	var botPlanet BotPlanet
	db.Where(&BotPlanet{UniverseName: bot.Universe, Language: bot.GetServer().Language, PlayerID: bot.Player.PlayerID, CelestialID: planetID}).Find(&botPlanet)
	lfbtmp := botPlanet.LfBuildings
	res.LfBuildings = lfbtmp
	res.LifeformString = lfbtmp.LifeformType.String()

	res.BuildingID, res.BuildingCountdown, res.ResearchID, res.ResearchCountdown, res.LfBuildingID, res.LfBuildingCountdown, res.LfResearchID, res.LfResearchCountdown = bot.ConstructionsBeingBuilt(ogame.CelestialID(planetID))

	planets, err := bot.GetEmpire(ogame.PlanetType)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, wrapper.ErrorResp(500, err.Error()))
	}
	moons, _ := bot.GetEmpire(ogame.MoonType)
	for _, p := range planets {
		if p.ID == ogame.CelestialID(planetID) {
			res.EmpireCelestial = p
			res.ResourcesDetails, _ = bot.GetResourcesDetails(p.ID)
			//return c.JSON(http.StatusOK, wrapper.SuccessResp(res))
		}
	}
	for _, p := range moons {
		if p.ID == ogame.CelestialID(planetID) {
			res.EmpireCelestial = p
			res.ResourcesDetails, _ = bot.GetResourcesDetails(p.ID)
			//return c.JSON(http.StatusOK, wrapper.SuccessResp(res))
		}
	}
	return c.Render(http.StatusOK, "empirePlanet", res)
	//return c.JSON(http.StatusBadRequest, wrapper.ErrorResp(400, "invalid planet id"))
}

func GetQueue(c echo.Context) error {
	bot := c.Get("bot").(*wrapper.OGame)
	queue := getQueue(bot.Universe, bot.GetServer().Language, bot.Player.PlayerID)
	return c.JSON(http.StatusOK, wrapper.SuccessResp(queue))
}

func GetShips(c echo.Context) error {
	bot := c.Get("bot").(*wrapper.OGame)
	celestialID, err := utils.ParseI64(c.Param("celestialID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, wrapper.ErrorResp(400, "invalid planet id"))
	}
	vals := url.Values{}
	vals.Add("page", "ingame")
	vals.Add("component", "fleetdispatch")
	vals.Add("cp", strconv.FormatInt(celestialID, 10))
	pageHTML, _ := bot.GetPageContent(vals)
	ext := bot.GetExtractor()
	ships := ext.ExtractFleet1Ships(pageHTML)

	resources := ext.ExtractResourcesDetailsFromFullPage(pageHTML)

	return c.JSON(http.StatusOK, struct {
		Ships     ogame.ShipsInfos `json:"ships"`
		Resources ogame.Resources  `json:"resources"`
	}{
		ships,
		resources.Available(),
	})
}

func GetFlights(c echo.Context) error {
	bot := c.Get("bot").(*wrapper.OGame)
	//empire, _ := bot.GetEmpire()
	data := struct {
		Bot                 *wrapper.OGame
		EmpireCelestial     []ogame.EmpireCelestial
		ResourcesDetails    ogame.ResourcesDetails
		Objs                map[ogame.ID]ogame.BaseOgameObj
		BuildingID          ogame.ID
		BuildingCountdown   int64
		ResearchID          ogame.ID
		ResearchCountdown   int64
		LfBuildingID        ogame.ID
		LfBuildingCountdown int64
		LfResearchID        ogame.ID
		LfResearchCountdown int64
		LcCapacity          int64
		ScCapacity          int64
		DsCapacity          int64
		PfCapacity          int64
		RCapacity           int64
		Planets             []ogame.Celestial
	}{
		Bot:  bot,
		Objs: ogame.Objs.GetAllObjs(),
		//EmpireCelestial: empire,
		LcCapacity: ogame.LargeCargo.GetCargoCapacity(bot.GetCachedResearch(), false, bot.CharacterClass().IsCollector(), bot.IsPioneers()),
		ScCapacity: ogame.SmallCargo.GetCargoCapacity(bot.GetCachedResearch(), false, bot.CharacterClass().IsCollector(), bot.IsPioneers()),
		DsCapacity: ogame.Deathstar.GetCargoCapacity(bot.GetCachedResearch(), false, bot.CharacterClass().IsCollector(), bot.IsPioneers()),
		PfCapacity: ogame.Pathfinder.GetCargoCapacity(bot.GetCachedResearch(), false, bot.CharacterClass().IsCollector(), bot.IsPioneers()),
		RCapacity:  ogame.Recycler.GetCargoCapacity(bot.GetCachedResearch(), false, bot.CharacterClass().IsCollector(), bot.IsPioneers()),
	}
	var planets []ogame.Celestial
	for _, v := range bot.GetCachedCelestials() {
		planets = append(planets, v.(ogame.Celestial))
	}
	data.Planets = planets

	return c.Render(http.StatusOK, "flights", data)
}

func GetQueueHandler(c echo.Context) error {
	bot := c.Get("bot").(*wrapper.OGame)
	bot.GetCachedCelestials()[0].GetCoordinate()
	planetID, err := utils.ParseI64(c.Param("planetID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, wrapper.ErrorResp(400, "invalid planet id"))
	}
	bq := getQueue(bot.Universe, bot.GetServer().Language, bot.Player.PlayerID)

	var planetBQ []BuildQueue
	for _, v := range bq {
		if planetID == v.CelestialID {
			planetBQ = append(planetBQ, v)
		}
	}

	return c.JSON(http.StatusOK, planetBQ)
}

func AddQueue(c echo.Context) error {
	bot := c.Get("bot").(*wrapper.OGame)
	planetID, err := utils.ParseI64(c.Param("planetID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, wrapper.ErrorResp(400, "invalid planet id"))
	}
	ogameID, err := utils.ParseI64(c.Param("ogameID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, wrapper.ErrorResp(400, "invalid ogame id"))
	}
	nbr, err := utils.ParseI64(c.Param("nbr"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, wrapper.ErrorResp(400, "invalid nbr"))
	}
	addQueue(bot.Universe, bot.GetServer().Language, bot.Player.PlayerID, planetID, ogameID, nbr)
	queue := getQueueForCelestial(bot.Universe, bot.GetServer().Language, bot.Player.PlayerID, ogame.CelestialID(planetID))
	if len(queue) == 1 {
		localBrain.ReRun(ogame.CelestialID(planetID))
	}
	return c.JSON(http.StatusOK, wrapper.SuccessResp(nil))
}

func GetTechDetailsHandler(c echo.Context) error {
	bot := c.Get("bot").(*wrapper.OGame)
	celestialID, err := utils.ParseI64(c.Param("celestialID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, wrapper.ErrorResp(400, "invalid celestial id"))
	}
	ogameid, err := utils.ParseI64(c.Param("ogameID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, wrapper.ErrorResp(400, "invalid ogame id"))
	}
	res, err := bot.TechnologyDetails(ogame.CelestialID(celestialID), ogame.ID(ogameid))
	if err != nil {
		return c.JSON(http.StatusBadRequest, wrapper.ErrorResp(400, err.Error()))

	}
	return c.JSON(http.StatusOK, wrapper.SuccessResp(res))
}

var ObjIDs = []ogame.ID{ogame.AllianceDepotID,
	ogame.CrystalMineID,
	ogame.CrystalStorageID,
	ogame.DeuteriumSynthesizerID,
	ogame.DeuteriumTankID,
	ogame.FusionReactorID,
	ogame.MetalMineID,
	ogame.MetalStorageID,
	ogame.MissileSiloID,
	ogame.NaniteFactoryID,
	ogame.ResearchLabID,
	ogame.RoboticsFactoryID,
	ogame.SeabedDeuteriumDenID,
	ogame.ShieldedMetalDenID,
	ogame.ShipyardID,
	ogame.SolarPlantID,
	ogame.SpaceDockID,
	ogame.LunarBaseID,
	ogame.SensorPhalanxID,
	ogame.JumpGateID,
	ogame.TerraformerID,
	ogame.UndergroundCrystalDenID,
	ogame.SolarSatelliteID,
	ogame.AntiBallisticMissilesID,
	ogame.GaussCannonID,
	ogame.HeavyLaserID,
	ogame.InterplanetaryMissilesID,
	ogame.IonCannonID,
	ogame.LargeShieldDomeID,
	ogame.LightLaserID,
	ogame.PlasmaTurretID,
	ogame.RocketLauncherID,
	ogame.SmallShieldDomeID,
	ogame.BattlecruiserID,
	ogame.BattleshipID,
	ogame.BomberID,
	ogame.ColonyShipID,
	ogame.CruiserID,
	ogame.DeathstarID,
	ogame.DestroyerID,
	ogame.EspionageProbeID,
	ogame.HeavyFighterID,
	ogame.LargeCargoID,
	ogame.LightFighterID,
	ogame.RecyclerID,
	ogame.SmallCargoID,
	ogame.CrawlerID,
	ogame.ReaperID,
	ogame.PathfinderID,
	ogame.ArmourTechnologyID,
	ogame.AstrophysicsID,
	ogame.CombustionDriveID,
	ogame.ComputerTechnologyID,
	ogame.EnergyTechnologyID,
	ogame.EspionageTechnologyID,
	ogame.GravitonTechnologyID,
	ogame.HyperspaceDriveID,
	ogame.HyperspaceTechnologyID,
	ogame.ImpulseDriveID,
	ogame.IntergalacticResearchNetworkID,
	ogame.IonTechnologyID,
	ogame.LaserTechnologyID,
	ogame.PlasmaTechnologyID,
	ogame.ShieldingTechnologyID,
	ogame.WeaponsTechnologyID,
	ogame.ResidentialSectorID,
	ogame.BiosphereFarmID,
	ogame.ResearchCentreID,
	ogame.AcademyOfSciencesID,
	ogame.NeuroCalibrationCentreID,
	ogame.HighEnergySmeltingID,
	ogame.FoodSiloID,
	ogame.FusionPoweredProductionID,
	ogame.SkyscraperID,
	ogame.BiotechLabID,
	ogame.MetropolisID,
	ogame.PlanetaryShieldID,
	ogame.MeditationEnclaveID,
	ogame.CrystalFarmID,
	ogame.RuneTechnologiumID,
	ogame.RuneForgeID,
	ogame.OriktoriumID,
	ogame.MagmaForgeID,
	ogame.DisruptionChamberID,
	ogame.MegalithID,
	ogame.CrystalRefineryID,
	ogame.DeuteriumSynthesiserID,
	ogame.MineralResearchCentreID,
	ogame.MetalRecyclingPlantID,
	ogame.AssemblyLineID,
	ogame.FusionCellFactoryID,
	ogame.RoboticsResearchCentreID,
	ogame.UpdateNetworkID,
	ogame.QuantumComputerCentreID,
	ogame.AutomatisedAssemblyCentreID,
	ogame.HighPerformanceTransformerID,
	ogame.MicrochipAssemblyLineID,
	ogame.ProductionAssemblyHallID,
	ogame.HighPerformanceSynthesiserID,
	ogame.ChipMassProductionID,
	ogame.NanoRepairBotsID,
	ogame.SanctuaryID,
	ogame.AntimatterCondenserID,
	ogame.VortexChamberID,
	ogame.HallsOfRealisationID,
	ogame.ForumOfTranscendenceID,
	ogame.AntimatterConvectorID,
	ogame.CloningLaboratoryID,
	ogame.ChrysalisAcceleratorID,
	ogame.BioModifierID,
	ogame.PsionicModulatorID,
	ogame.ShipManufacturingHallID,
	ogame.SupraRefractorID,
	ogame.IntergalacticEnvoysID,
	ogame.HighPerformanceExtractorsID,
	ogame.FusionDrivesID,
	ogame.StealthFieldGeneratorID,
	ogame.OrbitalDenID,
	ogame.ResearchAIID,
	ogame.HighPerformanceTerraformerID,
	ogame.EnhancedProductionTechnologiesID,
	ogame.LightFighterMkIIID,
	ogame.CruiserMkIIID,
	ogame.ImprovedLabTechnologyID,
	ogame.PlasmaTerraformerID,
	ogame.LowTemperatureDrivesID,
	ogame.BomberMkIIID,
	ogame.DestroyerMkIIID,
	ogame.BattlecruiserMkIIID,
	ogame.RobotAssistantsID,
	ogame.SupercomputerID,
	ogame.VolcanicBatteriesID,
	ogame.AcousticScanningID,
	ogame.HighEnergyPumpSystemsID,
	ogame.CargoHoldExpansionCivilianShipsID,
	ogame.MagmaPoweredProductionID,
	ogame.GeothermalPowerPlantsID,
	ogame.DepthSoundingID,
	ogame.IonCrystalEnhancementHeavyFighterID,
	ogame.ImprovedStellaratorID,
	ogame.HardenedDiamondDrillHeadsID,
	ogame.SeismicMiningTechnologyID,
	ogame.MagmaPoweredPumpSystemsID,
	ogame.IonCrystalModulesID,
	ogame.OptimisedSiloConstructionMethodID,
	ogame.DiamondEnergyTransmitterID,
	ogame.ObsidianShieldReinforcementID,
	ogame.RuneShieldsID,
	ogame.RocktalCollectorEnhancementID,
	ogame.CatalyserTechnologyID,
	ogame.PlasmaDriveID,
	ogame.EfficiencyModuleID,
	ogame.DepotAIID,
	ogame.GeneralOverhaulLightFighterID,
	ogame.AutomatedTransportLinesID,
	ogame.ImprovedDroneAIID,
	ogame.ExperimentalRecyclingTechnologyID,
	ogame.GeneralOverhaulCruiserID,
	ogame.SlingshotAutopilotID,
	ogame.HighTemperatureSuperconductorsID,
	ogame.GeneralOverhaulBattleshipID,
	ogame.ArtificialSwarmIntelligenceID,
	ogame.GeneralOverhaulBattlecruiserID,
	ogame.GeneralOverhaulBomberID,
	ogame.GeneralOverhaulDestroyerID,
	ogame.ExperimentalWeaponsTechnologyID,
	ogame.MechanGeneralEnhancementID,
	ogame.HeatRecoveryID,
	ogame.SulphideProcessID,
	ogame.PsionicNetworkID,
	ogame.TelekineticTractorBeamID,
	ogame.EnhancedSensorTechnologyID,
	ogame.NeuromodalCompressorID,
	ogame.NeuroInterfaceID,
	ogame.InterplanetaryAnalysisNetworkID,
	ogame.OverclockingHeavyFighterID,
	ogame.TelekineticDriveID,
	ogame.SixthSenseID,
	ogame.PsychoharmoniserID,
	ogame.EfficientSwarmIntelligenceID,
	ogame.OverclockingLargeCargoID,
	ogame.GravitationSensorsID,
	ogame.OverclockingBattleshipID,
	ogame.PsionicShieldMatrixID,
	ogame.KaeleshDiscovererEnhancementID,
}
