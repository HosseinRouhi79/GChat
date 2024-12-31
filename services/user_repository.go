package services

import (
	"chat/config"
	"chat/data/models"
	"chat/src/api/pkg/logging"
	"chat/src/constants"
	"errors"

	"gorm.io/gorm"
)

var cfg = config.GetConfig()
var logger = logging.NewLogger(cfg)

func (userService *UserService) ExistBytMobile(mobile string) (error, bool) {
	model := models.User{}
	err := userService.Db.Table("users").Where("mobile = ?", mobile).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Info(logging.Postgres, logging.Api, "Record(mobile) not found", nil)
			return nil, false
		}
		logger.Info(logging.Postgres, logging.Api, err.Error(), nil)
		return err, false

	}

	return nil, true
}

func (userService *UserService) ExistByEmail(email string) (error, bool) {
	model := models.User{}
	err := userService.Db.Table("users").Where("email = ?", email).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Info(logging.Postgres, logging.Api, "Record(email) not found", nil)
			return nil, false
		}
		logger.Info(logging.Postgres, logging.Api, err.Error(), nil)
		return err, false
	}
	return nil, true
}

func (userService *UserService) ExistByUsername(username string) (error, bool) {
	model := models.User{}
	err := userService.Db.Table("users").Where("username =?", username).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Info(logging.Postgres, logging.Api, "Record(username) not found", nil)
			return nil, false
		}
		logger.Info(logging.Postgres, logging.Api, err.Error(), nil)
		return err, false
	}
	return nil, true
}

func (userService *UserService) GetDefaultRole(username string) (roleID int, err error) {

	role := models.Role{}
	err = userService.Db.Table("roles").Where("name =?", constants.DefaultRole).First(&role).Error
	if err != nil {
		userService.Logger.Info(logging.Postgres, logging.DefaultRoleNotFound, "Role not found", nil)
		return 0, err
	}
	return role.Id, nil

}
