// khan
// https://github.com/topfreegames/khan
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2016 Top Free Games <backend@tfgco.com>

package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	. "github.com/franela/goblin"
	"github.com/topfreegames/khan/models"
	"github.com/topfreegames/khan/util"
)

func TestHookHandler(t *testing.T) {
	g := Goblin(t)

	testDb, err := models.GetTestDB()

	g.Assert(err == nil).IsTrue()

	g.Describe("Create Hook Handler", func() {
		g.It("Should create hook", func() {
			a := GetDefaultTestApp()
			game := models.GameFactory.MustCreate().(*models.Game)
			err := a.Db.Insert(game)
			AssertNotError(g, err)

			payload := util.JSON{
				"type":    models.GameUpdatedHook,
				"hookURL": "http://test/create",
			}
			res := PostJSON(a, GetGameRoute(game.PublicID, "/hooks"), t, payload)

			g.Assert(res.Raw().StatusCode).Equal(http.StatusOK)
			var result util.JSON
			json.Unmarshal([]byte(res.Body().Raw()), &result)
			g.Assert(result["success"]).IsTrue()
			g.Assert(result["publicID"] != "").IsTrue()

			dbHook, err := models.GetHookByPublicID(
				a.Db, game.PublicID, result["publicID"].(string),
			)
			AssertNotError(g, err)
			g.Assert(dbHook.GameID).Equal(game.PublicID)
			g.Assert(dbHook.PublicID).Equal(result["publicID"])
			g.Assert(dbHook.EventType).Equal(payload["type"])
			g.Assert(dbHook.URL).Equal(payload["hookURL"])
		})

		g.It("Should not create hook if missing parameters", func() {
			a := GetDefaultTestApp()
			route := GetGameRoute("game-id", "/hooks")
			res := PostJSON(a, route, t, util.JSON{})

			g.Assert(res.Raw().StatusCode).Equal(http.StatusBadRequest)
			var result util.JSON
			json.Unmarshal([]byte(res.Body().Raw()), &result)
			g.Assert(result["success"]).IsFalse()
			g.Assert(result["reason"]).Equal("type is required, hookURL is required")
		})

		g.It("Should not create hook if invalid payload", func() {
			a := GetDefaultTestApp()
			route := GetGameRoute("game-id", "/hooks")
			res := PostBody(a, route, t, "invalid")

			g.Assert(res.Raw().StatusCode).Equal(http.StatusBadRequest)
			var result util.JSON
			json.Unmarshal([]byte(res.Body().Raw()), &result)
			g.Assert(result["success"]).IsFalse()
			g.Assert(strings.Contains(result["reason"].(string), "While trying to read JSON")).IsTrue()
		})
	})

	g.Describe("Delete Hook Handler", func() {
		g.It("Should delete hook", func() {
			a := GetDefaultTestApp()

			hook, err := models.CreateHookFactory(testDb, "", models.GameUpdatedHook, "http://test/update")
			g.Assert(err == nil).IsTrue()

			res := Delete(a, GetGameRoute(hook.GameID, fmt.Sprintf("/hooks/%s", hook.PublicID)), t)

			g.Assert(res.Raw().StatusCode).Equal(http.StatusOK)

			var result util.JSON
			json.Unmarshal([]byte(res.Body().Raw()), &result)
			g.Assert(result["success"]).IsTrue()

			number, err := testDb.SelectInt("select count(*) from hooks where id=$1", hook.ID)
			g.Assert(err == nil).IsTrue()
			g.Assert(number == 0).IsTrue()
		})
	})
}
