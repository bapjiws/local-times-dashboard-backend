package controllers

import "github.com/revel/revel"

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	return c.Render()
}

func (c App) SearchCity(name string) revel.Result{
	response := make(map[string]interface{})
	response["the_city_you_searched_for"] = name
	return c.RenderJson(response)
}