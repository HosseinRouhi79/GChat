package migrations

import (
	"chat/config"
	"chat/data/db"
	"chat/data/models"
	"chat/src/api/pkg/logging"
	"reflect"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Up_1() {
	var cfg = config.GetConfig()
	var zap = logging.NewLogger(cfg)
	db := db.GetDB()
	tables := []interface{}{}
	var countries = models.Country{}
	var cities = models.City{}
	var users = models.User{}
	var roles = models.Role{}
	var messages = models.Message{}
	var userRoles = models.UserRole{}

	tables = addNewTable(db, countries, tables)
	tables = addNewTable(db, cities, tables)
	tables = addNewTable(db, users, tables)
	tables = addNewTable(db, roles, tables)
	tables = addNewTable(db, userRoles, tables)
	tables = addNewTable(db, messages, tables)

	err := db.Migrator().CreateTable(tables...)
	if err != nil {
		zap.Infof("%v", "Error creating table", err)
	}
	{
		ci := []models.City{{Name: "Test City", CountryID: 1}}
		c := models.Country{Name: "Test Country"}
		c1 := models.Country{Name: "Brazil"}
		c2 := models.Country{Name: "Italy"}
		c3 := models.Country{Name: "Germany"}
		c4 := models.Country{Name: "Iran"}
		c5 := models.Country{Name: "China"}
		c6 := models.Country{Name: "Bolivia"}

		c.Cities = &ci

		countryList := []models.Country{
			c, c1, c2, c3, c4, c5, c6,
		}

		db.Create(&countryList)
	}
	{
		pass := "12345678"
		hashedPass, _ := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
		roles := []*models.Role{
			{Name: "Admin"},
			{Name: "User"},
		}
		for _, role := range roles {
			db.Create(role)
		}
		// Create sample users
		users := []*models.User{
			{
				Username:  "user1",
				FirstName: "test1",
				LastName:  "test1",
				Email:     "test1@example.com",
				Password:  string(hashedPass),
				Enabled:   true,
				UserRoles: &[]models.UserRole{
					{RoleId: roles[0].Id}, // Assigning Admin role
					{RoleId: roles[1].Id}, // Assigning User role
				},
			},
			{
				Username:  "user2",
				FirstName: "test2",
				LastName:  "test2",
				Email:     "test2@example.com",
				Password:  string(hashedPass),
				Enabled:   true,
				UserRoles: &[]models.UserRole{
					{RoleId: roles[1].Id}, // Assigning User role
				},
			},
		}
		for _, user := range users {
			db.Create(user)
		}
	}
	zap.Infof("%v", reflect.TypeOf(tables[0]), reflect.TypeOf(tables[1]))
	zap.Infof("%v", db.Name())

}

func Down_1() {

}

func addNewTable(db *gorm.DB, table interface{}, tables []interface{}) []interface{} {
	if !db.Migrator().HasTable(table) {
		tables = append(tables, table)
	}
	return tables
}
