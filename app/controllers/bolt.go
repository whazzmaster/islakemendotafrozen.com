package controllers

import (
	"github.com/boltdb/bolt"
	"github.com/revel/revel"
	"islakemendotafrozen.com/app/models"
)

var (
	BoltDB *bolt.DB
)

func InitDB() {
	var (
		path  string
		found bool
		err   error
	)

	path, found = revel.Config.String("db.path")
	if !found {
		revel.ERROR.Printf("Could not find db.path configuration setting.")
	}

	revel.INFO.Printf("Using '%s' as path to data store.", path)

	BoltDB, err = bolt.Open(path, 0600, nil)
	if err != nil {
		revel.ERROR.Fatal(err)
	}

	err = BoltDB.Update(func(tx *bolt.Tx) error {
		revel.INFO.Printf("Ensuring buckets exist.")
    tx.CreateBucketIfNotExists([]byte("mendota"))
    return nil
	})
	if(err != nil) {
		revel.ERROR.Printf("Could not initialize BoltDB's buckets")
	}
}

type BoltController struct {
	*revel.Controller
	DB *bolt.DB
}

func (c *BoltController) Begin() revel.Result {
	c.DB = BoltDB
	return nil
}

func (c *BoltController) GetLake(name string) models.Lake {
	var lake = models.Lake{Name: name, Season: "fall", IsFrozen: false}

	if c.DB == nil {
		revel.ERROR.Printf("BoltDB is not available in the controller.")
	}

	revel.INFO.Printf("About to read from database")

	err := c.DB.View(func(tx *bolt.Tx) error {
		if tx == nil {
			revel.ERROR.Printf("Tx is nil in View block")
		}
		bucket := tx.Bucket([]byte(name))
		var season = string(bucket.Get([]byte("season")))
		var is_frozen = "true" == string(bucket.Get([]byte("frozen")))

		revel.INFO.Printf("Lake: %s, Season: %s, IsFrozen: %t", name, season, is_frozen)

		lake.Season = season
		lake.IsFrozen = is_frozen

		return nil
	})

	if err != nil {
		revel.ERROR.Printf("Something has gone wrong in the database read")
	}

	return lake
}
