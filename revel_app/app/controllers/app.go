package controllers

import (
	"github.com/revel/revel"
	"timezones_mc/revel_app/app"
	"timezones_mc/revel_app/app/models"
	"encoding/json"
	//"timezones_mc/utils"
)

type App struct {
	*revel.Controller
}

func (a App) Index() revel.Result {
	return a.Render()
}

func (a App) SuggestCities(name string) revel.Result{
	response, err := app.ES.SuggestDocuments("city_suggest", name, "suggest", "city_id")

	// TODO: handle error better? Or create a utils wrappper to avoid code duplication
	if err != nil {
		return a.RenderJson(map[string]interface{}{"error": err.Error()})
	}

	return a.RenderJson(response)
}

func (a App) FindCityById(id string) revel.Result {
	response, err := app.ES.FindDocumentById(id)

	// TODO: handle error better?
	if err != nil {
		return a.RenderJson(map[string]interface{}{"error": err.Error()})
	}

	city := &models.City{}

	//data, err := utils.GetBytes(response)
	//if err != nil {
	//	return a.RenderJson(map[string]interface{}{"error": err.Error()})
	//}

	if err = json.Unmarshal(response.([]byte), city); err != nil {
		return a.RenderJson(map[string]interface{}{"error": err.Error()})
	}

	return a.RenderJson(city)
}