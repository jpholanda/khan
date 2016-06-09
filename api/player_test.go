// khan
// https://github.com/topfreegames/khan
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2016 Top Free Games <backend@tfgco.com>

package api

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/Pallinder/go-randomdata"
	. "github.com/franela/goblin"
	. "github.com/onsi/gomega"
	"github.com/topfreegames/khan/models"
)

func TestPlayerHandler(t *testing.T) {
	g := Goblin(t)

	//special hook for gomega
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("Create Player Handler", func() {
		g.It("Should create player", func() {
			gameID := "api-cr-1"
			publicID := randomdata.FullName(randomdata.RandomGender)
			playerName := randomdata.FullName(randomdata.RandomGender)
			metadata := "{\"x\": 1}"

			a := GetDefaultTestApp()
			payload := map[string]string{
				"gameID":   gameID,
				"publicID": publicID,
				"name":     playerName,
				"metadata": metadata,
			}
			res := PostJSON(a, "/players", t, payload)

			res.Status(http.StatusOK)
			var result map[string]interface{}
			json.Unmarshal([]byte(res.Body().Raw()), &result)
			g.Assert(result["success"]).IsTrue()

			dbPlayer, err := models.GetPlayerByPublicID(gameID, publicID)
			g.Assert(err == nil).IsTrue()
			g.Assert(dbPlayer.GameID).Equal(gameID)
			g.Assert(dbPlayer.PublicID).Equal(publicID)
			g.Assert(dbPlayer.Name).Equal(playerName)
			g.Assert(dbPlayer.Metadata).Equal(metadata)
		})

		g.It("Should not create player if invalid payload", func() {
			a := GetDefaultTestApp()
			res := PostBody(a, "/players", t, "invalid")

			res.Status(http.StatusBadRequest)
			var result map[string]interface{}
			json.Unmarshal([]byte(res.Body().Raw()), &result)
			g.Assert(result["success"]).IsFalse()
			g.Assert(result["reason"]).Equal(
				"\n[IRIS]  Error: While trying to read [JSON invalid character 'i' looking for beginning of value] from the request body. Trace %!!(MISSING)s(MISSING)",
			)
		})

		g.It("Should not create player if invalid data", func() {
			gameID := "game-id-is-too-large-for-this-field-should-be-less-than-10-chars"
			playerID := randomdata.FullName(randomdata.RandomGender)
			playerName := randomdata.FullName(randomdata.RandomGender)
			metadata := "{\"x\": 1}"

			a := GetDefaultTestApp()
			payload := map[string]string{
				"gameID":   gameID,
				"playerID": playerID,
				"name":     playerName,
				"metadata": metadata,
			}
			res := PostJSON(a, "/players", t, payload)

			res.Status(http.StatusInternalServerError)
			var result map[string]interface{}
			json.Unmarshal([]byte(res.Body().Raw()), &result)
			g.Assert(result["success"]).IsFalse()
			g.Assert(result["reason"]).Equal("pq: value too long for type character varying(10)")
		})
	})
}