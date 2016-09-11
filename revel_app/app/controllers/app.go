package controllers

import (
	"github.com/revel/revel"
	"github.com/bapjiws/timezones_mc/revel_app/app"
)

type App struct {
	*revel.Controller
}

func (a App) Index() revel.Result {
	return a.Render()
}

func (a App) SuggestCities(name string) revel.Result{
	response, err := app.ES.SuggestDocuments("city_suggest", name, "suggest", "city_id")

	if err != nil {
		return a.RenderJson(map[string]interface{}{"error": err.Error()})
	}

	return a.RenderJson(response)
}

func (a App) FindCityById(id string) revel.Result {
	response, err := app.ES.FindDocumentById(id)

	if err != nil {
		return a.RenderJson(map[string]interface{}{"error": err.Error()})
	}

	return a.RenderJson(response)
}