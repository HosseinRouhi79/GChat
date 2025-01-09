package api

import (
	"chat/config"
	"chat/docs"
	"chat/src/api/middlewares/middlewares"
	"chat/src/api/routers"
	"chat/src/api/validation"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/websocket"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan string)

func handleWebSocket(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Failed to upgrade:", err)
		return
	}
	defer ws.Close()

	clients[ws] = true

	for {
		var message string

		err := ws.ReadJSON(&message)
		if err != nil {
			log.Println("Error reading message:", err)
			delete(clients, ws)
			break
		}

		message = "You: "+message 
		broadcast <- message
	}
}

func handleMessages() {
	for {
		
		message := <-broadcast

		for client := range clients {
			err := client.WriteJSON(message)
			if err != nil {
				log.Println("Error writing message:", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
func InitServer(cfg *config.Config) {
	r := gin.New()
	RegisterMainValidation()
	r.Use(gin.Logger(), gin.Recovery(), middlewares.Limitter(), middlewares.StructuredMiddleware()) // => r1 := gin.Default()
	// Load templates from the templates directory
	r.LoadHTMLGlob("../../templates/*")
	r.GET("/ws", handleWebSocket)

	
	go handleMessages()


	r.GET("/testo", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"title":   "HTMX with Gin",
			"heading": "Welcome to the Gin HTMX Example!",
		})
	})
	r.GET("/chat", func(c *gin.Context) {
		c.HTML(http.StatusOK, "chat.html", gin.H{
			"title":   "HTMX with Gin",
			"heading": "Welcome to the Gin HTMX Example!",
		})
	})
	RegisterRoute(r)
	RegisterSwagger(r, cfg)
	if err := r.Run(fmt.Sprintf(":%s", cfg.Server.InternalPort)); err != nil {
		panic(err)
	}
}

func RegisterMainValidation() {
	val, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		err := val.RegisterValidation("mobile", validation.MobileValidator, true)
		err2 := val.RegisterValidation("password", validation.PasswordValidator, true)
		if err != nil {
			log.Print(err.Error())
		}
		if err2 != nil {
			log.Print(err.Error())
		}
	}
}

func RegisterSwagger(r *gin.Engine, cfg *config.Config) {
	docs.SwaggerInfo.Title = "golang web api"
	docs.SwaggerInfo.Description = "golang web api"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Host = fmt.Sprintf("192.168.59.133:%s", cfg.Server.ExternalPort)
	docs.SwaggerInfo.Schemes = []string{"http"}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func RegisterRoute(r *gin.Engine) {
	v1 := r.Group("/")
	{
		healthGroup := v1.Group("")
		routers.Health(healthGroup)
	}

	v2 := r.Group("/api/v1/")
	{
		testGroup := v2.Group("test")
		routers.Test(testGroup)
	}

	v3 := r.Group("/api/v3/")
	{
		formGroup := v3.Group("form")
		routers.BodyBinder(formGroup)
	}

	v4 := r.Group("/api/")
	{
		formGroup := v4.Group("redis")
		routers.SetToRedisRouter(formGroup)
	}

	v5 := r.Group("/api/")
	{
		formGroup := v5.Group("redis")
		routers.GetFromRedisRouter(formGroup)
		routers.SetOtp(formGroup)
	}

	v6 := r.Group("/api/")
	{
		formGroup := v6.Group("")
		routers.GetJWT(formGroup)
		routers.Auth(formGroup)
	}
}
