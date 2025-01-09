package handlers

import (
	"chat/config"
	"chat/data/cache"
	"chat/services"
	"chat/src/api/pkg/logging"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Health struct{}
type Test struct {
	UserID  string `json:"userID"`
	Browser string `json:"browser"`
}
type PersonData struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Mobile    string `json:"mobile" binding:"mobile"`
	Password  string `json:"password" binding:"password"`
}

type Redis struct {
	Key   string
	Value string
}

type RedisKey struct {
	Key string
}

type JWT struct {
	service services.TokenService
}

type User struct {
	service services.UserService
}

func NewHealth() *Health {
	return &Health{}
}

func NewTest() *Test {
	return &Test{}
}

// HealthCheck godoc
// @Summary Health Check
// @Description Health Check
// @Tags health
// @Accept  json
// @Produce  json
// @Success 200 {object} helper.HTTPResponse "Success"
// @Failure 400 {object} helper.HTTPResponse "Failed"
// @Router /v1/health/ [get]
func (h Health) Health(c *gin.Context) {
	c.JSON(http.StatusOK, "health GET")
}

func (h Health) HealthPost(c *gin.Context) {
	message := c.DefaultPostForm("message", "")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(200, gin.H{
			"message": err.Error(),
		})
	}
	defer conn.Close()
	for {
		if message != "" {
			c.Header("HX-Trigger", "newMessage") // Optional: You can use this for triggering other actions with htmx
			// c.String(http.StatusOK, `
			// 	<div class="chat-message">
			// 		<div class="sender">You</div>
			// 		<div class="message-text">`+message+`</div>
			// 	</div>
			// `)
		} else {
			c.Status(http.StatusBadRequest)
		}
		conn.WriteMessage(websocket.TextMessage, []byte("Hello, WebSocket!"))
		time.Sleep(time.Second)
	}
	
}

func (h Health) HealthPostByID(c *gin.Context) {
	// c.JSON(http.StatusOK, "health GET ID")
	id := c.Params.ByName("id")
	c.JSON(http.StatusOK, fmt.Sprintf("health GET: %v", id))
}

func (t Test) HeaderBind(c *gin.Context) {
	if err := c.BindHeader(&t); err != nil {
		panic("Error binding")
	}
	c.JSON(http.StatusOK, gin.H{
		"Browser": t.Browser,
		"ID":      t.UserID,
	})
}

func (q Test) QueryBind(c *gin.Context) {
	id := c.Query("id") //using QueryArray instead of Query in case we have many ids
	name := c.Query("name")
	c.JSON(http.StatusOK, gin.H{
		"ID":   id,
		"name": name,
	})
}

func (u Test) UriBind(c *gin.Context) {
	id := c.Param("id")
	name := c.Param("name")
	c.JSON(http.StatusOK, gin.H{
		"ID":   id,
		"name": name,
	})
}

func (p PersonData) BodyBind(c *gin.Context) {

	err := c.ShouldBind(&p)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false})
		return
	}
	c.JSON(200, gin.H{"status": true})

}

func (inputs Redis) SetToRedis(c *gin.Context) {
	logger := logging.NewLogger(config.GetConfig())
	redisClient := cache.GetRedis()
	err := c.ShouldBind(&inputs)
	if err != nil {
		logger.Infof("failed to marshal:%v", err)
	}
	key := inputs.Key
	value := inputs.Value

	cache.Set[string](redisClient, key, value, 3600*time.Second)

}

func (inputs RedisKey) GetFromRedis(c *gin.Context) {
	logger := logging.NewLogger(config.GetConfig())
	redisClient := cache.GetRedis()
	keyString := c.Param("key")
	inputs.Key = keyString

	dest, err := cache.Get[string](redisClient, inputs.Key)
	fmt.Println(dest)
	if err != nil {
		logger.Infof("failed to get data from redis:%v", err)
		return
	}
	c.JSON(200, gin.H{
		"dest": dest,
	})

}

func (j JWT) Generate(c *gin.Context) {
	cfg := config.GetConfig()
	dto := &services.TokenDto{
		UserID:       1,
		FirstName:    "test",
		LastName:     "test",
		UserName:     "test",
		Email:        "test_email",
		MobileNumber: "test_mobile",
		Roles:        []string{"test_role"},
	}
	j.service = *services.NewTokenService(cfg)
	tk, _ := j.service.GenerateToken(dto)
	fmt.Println(tk.AccessToken)
}

func (j JWT) Validate(c *gin.Context) {
	cfg := config.GetConfig()
	tk := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InRlc3RfZW1haWwiLCJleHAiOjE3MTc5MjU2NTAsImZpcnN0TmFtZSI6InRlc3QiLCJpZCI6MSwibGFzdE5hbWUiOiJ0ZXN0IiwibW9iaWxlTnVtYmVyIjoidGVzdF9tb2JpbGUiLCJyb2xlcyI6WyJ0ZXN0X3JvbGUiXSwidXNlck5hbWUiOiJ0ZXN0In0.8fUlh5Brv_90R_enmfVFzUF7rX_s07on5RD_FFBhv7o"
	j.service = *services.NewTokenService(cfg)
	tk2, err := j.service.ValidateToken(tk)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(tk2)
}

func (user User) TestRepo(c *gin.Context) {
	cfg := config.GetConfig()
	user.service = *services.NewUserService(cfg)
	_, v := user.service.ExistBytMobile("09199810079")
	if v {
		fmt.Println("exist")
	} else {
		fmt.Println("not exist")
	}
}
