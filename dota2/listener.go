// Copyright © 2017 Matthew Kleinschafer
// 
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the “Software”), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
// 
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
// 
// THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package dota2

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

// State information packet from the DOTA 2 GSI system
type GameState struct {
	Player     Player `json:"player"`
	Hero       Hero `json:"hero"`
	Map        Map `json:"map"`
	Previously *GameState `json:"previously"`
}

// Information about the active game map
type Map struct {
	// general information
	Name    string `json:"name"`
	MatchId int `json:"match_id"`

	// time information
	GameTime             int `json:"game_time"`
	ClockTime            int `json:"clock_time"`
	WardPurchaseCooldown int `json:"ward_purchase_cooldown"`
	IsDaytime            bool `json:"daytime"`
	IsNightstalkerNight  bool `json:"nightstalker_night"`
}

// Status information for a particular player in the game
type Player struct {
	// general player info
	SteamId  int `json:"steam_id"`
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
	Id               int `json:"id"`
	Name             string `json:"name"`
	Level            int `json:"level"`
	IsAlive          bool `json:"alive"`
	SecondsToRespawn int `json:"respawn_seconds"`

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

// Starts a HTTP server listening on the given port for game status updates,
// returns a channel that is provided with states as the updates occur
func ListenForUpdates(port int) chan *GameState {
	updates := make(chan *GameState)

	// starts the http server
	start := func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			state := new(GameState)

			decoder := json.NewDecoder(r.Body)
			err := decoder.Decode(state)
			if err != nil {
				log.Print(err)
			}

			updates <- state
		})

		addr := ":" + strconv.Itoa(port)
		log.Print("Starting game state listener on port ", addr)

		err := http.ListenAndServe(addr, mux)
		if err != nil {
			log.Fatal(err)
		}
	}

	go start()

	return updates
}
