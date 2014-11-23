package controllers

import (
	"github.com/revel/revel"
)

type App struct {
	BoltController
}

func (c App) Index() revel.Result {
	mendota := c.GetLake("mendota")
	return c.Render(mendota)
}
