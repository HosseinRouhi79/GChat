package main

import (
	"chat/config"
	"chat/data/cache"
	"chat/data/db"
	"chat/src/api"
	"log"
	"net/http"
)

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1
// @securityDefinitions.apikey AuthBearer
// @in header
// @name Authorization
func main() {

	go func() {
		log.Println("Starting pprof on :6060")
		if err := http.ListenAndServe(":6060", nil); err != nil {
			log.Fatalf("pprof failed: %v", err)
		}
	}()

	cfg := config.GetConfig()
	cache.InitRedis(cfg)
	err := db.InitDB(cfg)
	if err != nil {
		panic("cannot connect to database(main error)")
	}
	// migrations.Up_1()
	api.InitServer(cfg)
	defer cache.CloseRedis()
	defer db.CloseDB()
}
