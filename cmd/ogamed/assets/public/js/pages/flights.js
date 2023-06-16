Vue.filter('formatDate', function(value) {
    if (value) {
        return moment(String(value)).format('DD.MM hh:mm:ss')
    }
});

function dezInt(num, size, prefix) {
    prefix = prefix ? prefix: "0";
    var minus = num < 0 ? "-": "",
        result = prefix === "0" ? minus: "";
    num = Math.abs(parseInt(num, 10));
    size -= ("" + num).length;
    for (var i = 1; i <= size; i++) {
        result += "" + prefix;
    }
    result += (prefix !== "0" ? minus: "") + num;
    return result;
}
function getFormatedDate(timestamp, format) {
    var currTime = new Date();
    currTime.setTime(timestamp);
    str = format;
    str = str.replace('[d]', dezInt(currTime.getDate(), 2));
    str = str.replace('[m]', dezInt(currTime.getMonth() + 1, 2));
    str = str.replace('[j]', parseInt(currTime.getDate()));
    str = str.replace('[Y]', currTime.getFullYear());
    str = str.replace('[y]', currTime.getFullYear().toString().substr(2, 4));
    str = str.replace('[G]', currTime.getHours());
    str = str.replace('[H]', dezInt(currTime.getHours(), 2));
    str = str.replace('[i]', dezInt(currTime.getMinutes(), 2));
    str = str.replace('[s]', dezInt(currTime.getSeconds(), 2));
    return str;
}

var app = new Vue({
    delimiters: ['${', '}'],
    el: '#app1',
    data: {
        schedule: '',
        planets: serverVars.planets,
        planetLoading: false,
        loading: true,
        fleets: [],
        slots: {InUse: 0, Total: 0, ExpInUse: 0, ExpTotal: 0},
        ships: {
            LightFighter: 0,
            HeavyFighter: 0,
            Cruiser: 0,
            Battleship: 0,
            Battlecruiser: 0,
            Bomber: 0,
            Destroyer: 0,
            Deathstar: 0,
            SmallCargo: 0,
            LargeCargo: 0,
            ColonyShip: 0,
            Recycler: 0,
            EspionageProbe: 0,
            Crawler: 0,
            Reaper: 0,
            Pathfinder: 0,
            SolarSatellite: 0,
        },
        resources: {
            Metal: 0,
            Crystal: 0,
            Deuterium: 0,
        },
        metal: 0,
        crystal: 0,
        deuterium: 0,
        LightFighter: 0,
        HeavyFighter: 0,
        Cruiser: 0,
        Battleship: 0,
        Battlecruiser: 0,
        Bomber: 0,
        Destroyer: 0,
        Deathstar: 0,
        SmallCargo: 0,
        LargeCargo: 0,
        ColonyShip: 0,
        Recycler: 0,
        EspionageProbe: 0,
        Crawler: 0,
        Reaper: 0,
        Pathfinder: 0,
        SolarSatellite: 0,
        mission: 1,
        allShips: false,
        allResources: false,
        planetType: 1,
        speed: 10,
        origin: '',
        destGalaxy: null,
        destSystem: null,
        destPosition: null,
        humanTime: '-',
        fuel: 0,
        fuelCapacity: 0,
        cargo: 0,
        arrivalFmt: '-',
        returnFmt: '-',
        arrivalSecs: 0,
        missionExpeditionActive:   false,
        missionColonisationActive: false,
        missionDeploymentActive:   false,
        missionRecycleActive:      false,
        missionTransportActive:    false,
        missionEspionageActive:    false,
        missionACSDefendActive:    false,
        missionAttackActive:       false,
        missionACSAttackActive:    false,
        missionDestructionActive:  false,
    },
    watch: {
        LightFighter:   function() { this.calcFlightTime(); },
        HeavyFighter:   function() { this.calcFlightTime(); },
        Cruiser:        function() { this.calcFlightTime(); },
        Battleship:     function() { this.calcFlightTime(); },
        Battlecruiser:  function() { this.calcFlightTime(); },
        Bomber:         function() { this.calcFlightTime(); },
        Destroyer:      function() { this.calcFlightTime(); },
        Deathstar:      function() { this.calcFlightTime(); },
        SmallCargo:     function() { this.calcFlightTime(); },
        LargeCargo:     function() { this.calcFlightTime(); },
        ColonyShip:     function() { this.calcFlightTime(); },
        Recycler:       function() { this.calcFlightTime(); },
        EspionageProbe: function() { this.calcFlightTime(); },
        Crawler:        function() { this.calcFlightTime(); },
        Reaper:         function() { this.calcFlightTime(); },
        Pathfinder:     function() { this.calcFlightTime(); },
        speed:          function() { this.calcFlightTime(); },
        origin:         function() { this.calcFlightTime(); },
        planetType:     function() { this.updateMissions(); this.calcFlightTime(); },
        destGalaxy:     function() { this.updateMissions(); this.calcFlightTime(); },
        destSystem:     function() { this.updateMissions(); this.calcFlightTime(); },
        destPosition:   function() { this.updateMissions(); this.calcFlightTime(); },
        mission:        function() { this.missionChanged(); },
        allShips: function() { if (this.allShips) { this.arrivalSecs = 0; } else { this.calcFlightTime(); } },
    },
    computed: {
        lcNeeded: function() {
            var total = this.cargoNeeded()
            return Math.max(this.LargeCargo + Math.ceil(total / serverVars.lcCapacity ), 0);
        },
        scNeeded: function() {
            var total = this.cargoNeeded()
            return Math.max(this.SmallCargo + Math.ceil(total / serverVars.scCapacity ), 0);
        },
        dsNeeded: function() {
            var total = this.cargoNeeded()
            return Math.max(this.Deathstar + Math.ceil(total / serverVars.dsCapacity ), 0);
        },
        pfNeeded: function() {
            var total = this.cargoNeeded()
            return Math.max(this.Pathfinder + Math.ceil(total / serverVars.pfCapacity ), 0);
        },
        rNeeded: function() {
            var total = this.cargoNeeded()
            return Math.max(this.Recycler + Math.ceil(total / serverVars.rCapacity ), 0);
        },
        hasShipSelected: function() {
            return this.LightFighter > 0 ||
                this.HeavyFighter > 0 ||
                this.Cruiser > 0 ||
                this.Battleship > 0 ||
                this.Battlecruiser > 0 ||
                this.Bomber > 0 ||
                this.Destroyer > 0 ||
                this.Deathstar > 0 ||
                this.SmallCargo > 0 ||
                this.LargeCargo > 0 ||
                this.ColonyShip > 0 ||
                this.Recycler > 0 ||
                this.EspionageProbe > 0 ||
                this.Crawler > 0 ||
                this.Reaper > 0 ||
                this.Pathfinder > 0;
        },
    },
    methods: {
        toggleMetal:     function() { this.metal     = (this.metal     === this.resources.Metal    ) ? 0 : this.resources.Metal;     },
        toggleCrystal:   function() { this.crystal   = (this.crystal   === this.resources.Crystal  ) ? 0 : this.resources.Crystal;   },
        toggleDeuterium: function() { this.deuterium = (this.deuterium === this.resources.Deuterium) ? 0 : this.resources.Deuterium; },

        cargoNeeded: function() {
            var metal = this.metal || 0;
            var crystal = this.crystal || 0;
            var deuterium = this.deuterium || 0;
            var total = metal+crystal+deuterium;
            total -= this.LargeCargo * serverVars.lcCapacity;
            total -= this.SmallCargo * serverVars.scCapacity;
            total -= this.Deathstar * serverVars.dsCapacity;
            total -= this.Pathfinder * serverVars.pfCapacity;
            total -= this.Recycler * serverVars.rCapacity;
            return total;
        },

        setMission: function(newMission) {
            this.mission = newMission;
        },
        canUseJumpGate: function() {
            return this.destIsOurs() &&
                this.planetType === 3 &&
                this.originIsMoon();
        },
        originIsMoon: function(planetID) {
            // for (var i = 0; i < this.planets.length; i++) {
            //     if (this.origin === this.planets[i].PlanetID) {
            //         if (this.planets[i].PlanetType === 3) {
            //             return true;
            //         }
            //         return false;
            //     }
            // }
            for (var i = 0; i < this.planets.length; i++) {
                if (this.origin === this.planets[i].ID) {
                    if (this.planets[i].Coordinate.Type === 3) {
                        return true;
                    }
                    return false;
                }
            }
            return false;
        },
        destIsOurs: function() {
            // for (var i = 0; i < this.planets.length; i++) {
            //     if (this.destGalaxy === this.planets[i].Galaxy &&
            //         this.destSystem === this.planets[i].System &&
            //         this.destPosition === this.planets[i].Position &&
            //         this.planetType === this.planets[i].PlanetType) {
            //         return true;
            //     }
            // }
            for (var i = 0; i < this.planets.length; i++) {
                if (this.destGalaxy === this.planets[i].Coordinate.Galaxy &&
                    this.destSystem === this.planets[i].Coordinate.System &&
                    this.destPosition === this.planets[i].Coordinate.Position &&
                    this.planetType === this.planets[i].Coordinate.Type) {
                    return true;
                }
            }
            return false;
        },

        deactivateMissions: function() {
            this.missionExpeditionActive   = false;
            this.missionColonisationActive = false;
            this.missionDeploymentActive   = false;
            this.missionRecycleActive      = false;
            this.missionTransportActive    = false;
            this.missionEspionageActive    = false;
            this.missionACSDefendActive    = false;
            this.missionAttackActive       = false;
            this.missionACSAttackActive    = false;
            this.missionDestructionActive  = false;
        },

        updateMissions: function() {
            this.deactivateMissions();
            if (this.planetType === 2) { // Debris field
                this.missionRecycleActive      = true;
            } else if (this.destPosition === 16) {
                this.missionExpeditionActive   = true;
            } else if (this.destIsOurs()) {
                this.missionDeploymentActive   = true;
                this.missionTransportActive    = true;
                this.missionACSDefendActive    = true;
            } else {
                this.missionColonisationActive = true;
                this.missionTransportActive    = true;
                this.missionEspionageActive    = true;
                this.missionACSDefendActive    = true;
                this.missionAttackActive       = true;
                this.missionACSAttackActive    = true;
                if (this.planetType === 3) { // Moon
                    this.missionDestructionActive  = true;
                }
            }

        },

        calcFlightTime: function() {
            if (this.allShips) {
                return;
            }
            this.debouncedFlightTime();
        },

        debouncedFlightTime: _.debounce(function() {
            var self = this;
            var formData = new FormData();
            formData.append('csrf', serverVars.csrf);
            formData.append('mission',         this.mission);
            formData.append('origin',         this.origin);
            formData.append('destGalaxy',     this.destGalaxy);
            formData.append('destSystem',     this.destSystem);
            formData.append('destPosition',   this.destPosition);
            formData.append('speed',          this.speed);
            formData.append('LightFighter',   this.LightFighter);
            formData.append('HeavyFighter',   this.HeavyFighter);
            formData.append('Cruiser',        this.Cruiser);
            formData.append('Battleship',     this.Battleship);
            formData.append('Battlecruiser',  this.Battlecruiser);
            formData.append('Bomber',         this.Bomber);
            formData.append('Destroyer',      this.Destroyer);
            formData.append('Deathstar',      this.Deathstar);
            formData.append('SmallCargo',     this.SmallCargo);
            formData.append('LargeCargo',     this.LargeCargo);
            formData.append('ColonyShip',     this.ColonyShip);
            formData.append('Recycler',       this.Recycler);
            formData.append('EspionageProbe', this.EspionageProbe);
            formData.append('Crawler',        this.Crawler);
            formData.append('Reaper',         this.Reaper);
            formData.append('Pathfinder',     this.Pathfinder);
            $.ajax({
                url: "/bot/flighttime",
                data: formData,
                processData: false,
                contentType: false,
                type: 'POST',
                success: function(res) {
                    self.arrivalSecs = res.secs;
                    self.humanTime = res.human_time;
                    if (self.humanTime === '') { self.humanTime = '-'; }
                    self.fuel = res.fuel;
                    self.cargo = res.cargo;
                    self.calcArrival();
                },
                error: function(err) { console.log(err); },
            });
        }, 500),

        calcArrival: function() {
            var now = new Date().getTime();
            var oneWay = now + this.arrivalSecs * 1000;
            var twoWay = now + this.arrivalSecs * 1000 * 2;
            if (this.arrivalSecs > 0) {
                this.arrivalFmt = getFormatedDate(oneWay, '[d].[m].[y] [G]:[i]:[s]');
                this.returnFmt = getFormatedDate(twoWay, '[d].[m].[y] [G]:[i]:[s]');
            } else {
                this.arrivalFmt = '-';
                this.returnFmt = '-';
            }
        },

        shortcutClicked: function(galaxy, system, position, planetType) {
            this.destGalaxy = galaxy;
            this.destSystem = system;
            this.destPosition = position;
            this.planetType = planetType;
            if (this.mission === 15) {
                this.planetType = 1;
                this.destPosition = 16;
            } else if (this.mission === 8) {
                this.planetType = 2;
            }
        },
        selectAllShips: function() {
            this.LightFighter = this.ships.LightFighter;
            this.HeavyFighter = this.ships.HeavyFighter;
            this.Cruiser = this.ships.Cruiser;
            this.Battleship = this.ships.Battleship;
            this.Battlecruiser = this.ships.Battlecruiser;
            this.Bomber = this.ships.Bomber;
            this.Destroyer = this.ships.Destroyer;
            this.Deathstar = this.ships.Deathstar;
            this.SmallCargo = this.ships.SmallCargo;
            this.LargeCargo = this.ships.LargeCargo;
            this.ColonyShip = this.ships.ColonyShip;
            this.Recycler = this.ships.Recycler;
            this.EspionageProbe = this.ships.EspionageProbe;
            this.Crawler = this.ships.Crawler;
            this.Reaper = this.ships.Reaper;
            this.Pathfinder = this.ships.Pathfinder;
            this.SolarSatellite = this.ships.SolarSatellite;
        },
        noShips: function() {
            this.LightFighter = 0;
            this.HeavyFighter = 0;
            this.Cruiser = 0;
            this.Battleship = 0;
            this.Battlecruiser = 0;
            this.Bomber = 0;
            this.Destroyer = 0;
            this.Deathstar = 0;
            this.SmallCargo = 0;
            this.LargeCargo = 0;
            this.ColonyShip = 0;
            this.Recycler = 0;
            this.EspionageProbe = 0;
            this.Crawler = 0;
            this.Reaper = 0;
            this.Pathfinder = 0;
            this.SolarSatellite = 0;
        },
        selectAllResources: function() {
            this.metal = this.resources.Metal;
            this.crystal = this.resources.Crystal;
            this.deuterium = this.resources.Deuterium;
        },
        noResources: function() {
            this.metal = 0;
            this.crystal = 0;
            this.deuterium = 0;
        },
        planetChanged: function(e) {
            let celestialID = e.target.value;
            if (celestialID === "") return;
            let self = this;
            this.planetLoading = true;
            let p1 = fetch('/bot/ships/'+celestialID).then(res => res.json());
            let p2 = fetch('/bot/resources/'+celestialID).then(res => res.json());
            Promise.all([p1, p2]).then(function(values) {
                self.ships = values[0].ships;
                self.resources = values[1].resources;
            }).finally(function() {
                self.planetLoading = false;
            });
        },
        missionChanged: function(e) {
            if (this.mission === 8) { // Recycle
                this.planetType = 2;
            } else if (this.mission === 15) { // Expedition
                this.destPosition = 16;
                this.planetType = 1;
            } else if (this.mission === 9) { // Destroy
                this.planetType = 3; // Moon
            } else {
                if (this.planetType === 2) { // DebrisType
                    this.planetType = 1;
                }
            }
        },
    },
    created () {
        var self = this;
        $.ajax({url: '/bot/get-fleets',
            success: function(json) {
                self.slots = json.slots;
                self.fleets = json.fleets;
                self.loading = false;
            },
            error: function(err) { console.log(err); },
        });


        var self = this;
        setInterval(function () { self.calcArrival(); }, 1000);
    }
});