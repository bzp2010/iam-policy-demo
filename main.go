package main

import (
	"context"
	"log"
	"os"
	"policydemo/pkg/database"
	ladonhelper "policydemo/pkg/ladon_helper"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/ory/ladon"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := os.Getenv("DB_URL")
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic("failed to connect database")
	}

	runServer()
}

type hasAccess struct {
	Resource           string            `json:"resource"`
	Action             string            `json:"action"`
	UseToken           bool              `json:"use_token"`
	SimulationBelongTo string            `json:"sim_belong_to"`
	SimulationLabels   map[string]string `json:"sim_labels"`
}
type hasAccessMap map[string]hasAccess

func runServer() {
	manager, err := ladonhelper.NewCachedDBManager(db)
	if err != nil {
		panic("Failed to initialize IAM system")
	}
	warden := &ladon.Ladon{
		Manager: manager,
	}

	r := gin.Default()
	r.GET("/initdb", func(c *gin.Context) {
		database.InitDatabase(db)
		manager.RebuildCache(db)
		c.AbortWithStatus(200)
	})
	r.POST("/hasAccess", func(c *gin.Context) {
		username := c.Query("username")
		var asserts hasAccessMap
		c.BindJSON(&asserts)

		var user database.User
		err := db.
			Model(&database.User{}).
			Where("username = ?", username).
			First(&user).Error
		if err != nil {
			return
		}

		ctx := context.Background()
		result := map[string]bool{}
		for key, ha := range asserts {
			err = warden.IsAllowed(ctx, &ladon.Request{
				Subject:  username,
				Action:   ha.Action,
				Resource: ha.Resource,
				Context: ladon.Context{
					"needToken":  ha.UseToken,
					"belongTo":   ha.SimulationBelongTo,
					"api:labels": ha.SimulationLabels,
				},
			})

			if err != nil {
				//fmt.Printf("scanln: %+v", err)
				result[key] = false
				continue
			}

			result[key] = true
		}

		c.AbortWithStatusJSON(200, result)
	})
	r.GET("/caniuse_service", func(c *gin.Context) {
		// Assuming the system has the following services, it should query from the database.
		services := []string{"api:gateway:service:a", "api:gateway:service:b", "api:gateway:service:c", "api:gateway:service:1", "api:gateway:service:2", "api:gateway:service:3"}
		username := c.Query("username")

		ctx := context.Background()
		result := map[string]bool{}
		for _, service := range services {
			err = warden.IsAllowed(ctx, &ladon.Request{
				Subject:  username,
				Action:   "gateway:GetService",
				Resource: service,
				Context: ladon.Context{
					"needToken": false,
				},
			})
			if err != nil {
				//fmt.Printf("scanln: %+v", err)
				result[service] = false
				continue
			}
			result[service] = true
		}

		c.AbortWithStatusJSON(200, result)

	})
	r.Run() // listen and serve on 0.0.0.0:8080
}
