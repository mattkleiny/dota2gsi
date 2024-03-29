// Copyright 2017, the project authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package dota2

import (
	"encoding/json"
	"log"
	"net/http"
)

// State information packet from the DOTA 2 GSI system
type GameState struct {
	Player     Player     `json:"player"`
	Hero       Hero       `json:"hero"`
	Map        Map        `json:"map"`
	Previously *GameState `json:"previously"`
}

// Information about the active game map
type Map struct {
	// general information
	Name    string `json:"name"`
	MatchId int    `json:"match_id"`

	// time information
	GameTime             int  `json:"game_time"`
	ClockTime            int  `json:"clock_time"`
	WardPurchaseCooldown int  `json:"ward_purchase_cooldown"`
	IsDaytime            bool `json:"daytime"`
	IsNightstalkerNight  bool `json:"nightstalker_night"`
}

// Status information for a particular player in the game
type Player struct {
	// general player info
	SteamId  int    `json:"steam_id"`
	Name     string `json:"name"`
	TeamName string `json:"team_name"`
	Activity string `json:"activity"`

	// k/d/a ratio
	Kills      int `json:"kills"`
	KillStreak int `json:"kill_streak"`
	Deaths     int `json:"deaths"`
	Assists    int `json:"assists"`

	// creep stats
	LastHits int `json:"last_hits"`
	Denies   int `json:"denies"`

	// gold stats
	Gold           int `json:"gold"`
	GoldReliable   int `json:"gold_reliable"`
	GoldUnreliable int `json:"gold_unreliable"`
	GoldPerMinute  int `json:"gpm"`

	// experience stats
	ExperiencePerMinute int `json:"xpm"`
}

// Status information for a particular hero in te game
type Hero struct {
	// general hero info
	Id               int    `json:"id"`
	Name             string `json:"name"`
	Level            int    `json:"level"`
	IsAlive          bool   `json:"alive"`
	SecondsToRespawn int    `json:"respawn_seconds"`

	// buyback information
	BuybackCost     int `json:"buyback_cost"`
	BuybackCooldown int `json:"buyback_cooldown"`

	// health stats
	Health        int `json:"health"`
	MaxHealth     int `json:"max_health"`
	HealthPercent int `json:"health_percent"`

	// mana stats
	Mana        int `json:"mana"`
	MaxMana     int `json:"max_mana"`
	ManaPercent int `json:"mana_percent"`

	// effect statuses
	IsStunned     bool `json:"stunned"`
	IsDisarmed    bool `json:"disarmed"`
	IsMagicImmune bool `json:"magicimmune"`
	IsHexed       bool `json:"hexed"`
	IsMuted       bool `json:"muted"`
	IsBroken      bool `json:"break"`
	HasDebuff     bool `json:"has_debuff"`
}

// Starts an HTTP server listening on the given address for game state updates.
// Returns a channel that is provided with states as the updates occur.
func ListenForUpdates(address string) chan *GameState {
	updates := make(chan *GameState)

	start := func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/", GameStateHandler(updates))

		log.Print("Starting game state listener on ", address)
		err := http.ListenAndServe(address, mux)
		if err != nil {
			log.Fatal(err)
		}
	}

	go start()

	return updates
}

// Creates a new GSI state handler HTTP func that writes to the given channel.
func GameStateHandler(updates chan *GameState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		state := new(GameState)

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(state)
		if err != nil {
			log.Print(err)
		}

		updates <- state
	}
}
