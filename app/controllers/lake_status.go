package controllers

import (
	"encoding/json"
	"github.com/boltdb/bolt"
	"github.com/revel/revel"
	"islakemendotafrozen.com/app/models"
	"net/http"
)

type LakeStatus struct {
	BoltController
}

func (c LakeStatus) Show(name string) revel.Result {
	lake := c.GetLake(name)
	c.Response.Status = http.StatusOK
	return c.RenderJson(lake)
}

func (c LakeStatus) CanUpdate() bool {
	auth_required := revel.Config.BoolDefault("auth_required", true)
	if auth_required {
		// TODO check headers for params
		return false
	} else {
		return true
	}
}

func (c LakeStatus) Create(name string) revel.Result {
	lake := new(models.Lake)
	decode_err := json.NewDecoder(c.Request.Body).Decode(&lake)
	if decode_err != nil {
		revel.ERROR.Printf("Could not decide JSON request body")
	} else {
		lake.Name = name
	}

	revel.INFO.Printf("[POST] [Deserialized] Lake: %s, Season: %s, IsFrozen: %t", lake.Name, lake.Season, lake.IsFrozen)

	if c.CanUpdate() {
		err := c.DB.Update(func(tx *bolt.Tx) error {
			bucket := tx.Bucket([]byte(name))
			err := bucket.Put([]byte("season"), []byte(lake.Season))

			is_frozen := "false"
			if lake.IsFrozen {
				is_frozen = "true"
			}
			err = bucket.Put([]byte("frozen"), []byte(is_frozen))
			return err
		})
		if err != nil {
			revel.ERROR.Fatal("Could not update lake status.", err)
		}

		c.Response.Status = http.StatusOK
		return c.RenderJson(lake)
	}

	c.Response.Status = http.StatusUnauthorized
	return c.RenderText("")
}
