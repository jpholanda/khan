// khan
// https://github.com/topfreegames/khan
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2016 Top Free Games <backend@tfgco.com>

package api

import (
	"github.com/kataras/iris"
	"github.com/topfreegames/khan/models"
)

type playerDataChangePayload struct {
	PublicID string
	Name     string
	Metadata string
}

// CreatePlayerHandler is the handler responsible for creating new players
func CreatePlayerHandler(app *App) func(c *iris.Context) {
	return func(c *iris.Context) {
		gameID := c.Param("gameID")

		var payload playerDataChangePayload
		if err := c.ReadJSON(&payload); err != nil {
			FailWith(400, err.Error(), c)
			return
		}

		db := GetCtxDB(c)

		player, err := models.CreatePlayer(
			db,
			gameID,
			payload.PublicID,
			payload.Name,
			payload.Metadata,
		)

		if err != nil {
			FailWith(500, err.Error(), c)
			return
		}

		SucceedWith(map[string]interface{}{
			"publicID": player.PublicID,
		}, c)
	}
}

// UpdatePlayerHandler is the handler responsible for updating existing
func UpdatePlayerHandler(app *App) func(c *iris.Context) {
	return func(c *iris.Context) {
		gameID := c.Param("gameID")
		playerPublicID := c.Param("playerPublicID")

		var payload playerDataChangePayload
		if err := c.ReadJSON(&payload); err != nil {
			FailWith(400, err.Error(), c)
			return
		}

		db := GetCtxDB(c)

		_, err := models.UpdatePlayer(
			db,
			gameID,
			playerPublicID,
			payload.Name,
			payload.Metadata,
		)

		if err != nil {
			FailWith(500, err.Error(), c)
			return
		}

		SucceedWith(map[string]interface{}{}, c)
	}
}
