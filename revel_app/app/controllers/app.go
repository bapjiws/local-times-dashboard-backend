package controllers

import (
	"github.com/revel/revel"
	"timezones_mc/revel_app/app"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	return c.Render()
}

func (c App) SearchCity(name string) revel.Result{
	response, err := app.ES.SuggestDocuments("city_suggest", name, "suggest", "city_id")

	// TODO: handle error better?
	if err != nil {
		return c.RenderJson(map[string]interface{}{"error": err.Error()})
	}

	return c.RenderJson(response)
}