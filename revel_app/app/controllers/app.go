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
	response, _ := app.ES.SuggestDocuments("city_suggest", name, "suggest", "city_id")
	// TODO: handle error

	//response := make(map[string]interface{})
	//response["the_city_you_searched_for"] = name

	return c.RenderJson(response)
}