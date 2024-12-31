package services

import (
	"chat/config"
	"chat/data/cache"
	"chat/src/api/pkg/logging"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v7"
)

type OtpService struct {
	cfg         *config.Config
	logger      logging.Logger
	redisClient *redis.Client
}

type OtpCred struct {
	Value string
	Used  bool
}

func NewOtpService(cfg *config.Config) *OtpService {

	logger := logging.NewLogger(cfg)
	redis := cache.GetRedis()
	return &OtpService{cfg: cfg, logger: logger, redisClient: redis}
}

func (otpService *OtpService) SetOtp(mobile string, otp string) error {

	var prefix string = "redis"
	key := fmt.Sprintf("%s:%s", prefix, mobile)
	val := OtpCred{Value: otp, Used: false}

	dest, err := cache.Get[OtpCred](otpService.redisClient, key)

	if err == nil && dest.Used {
		errMessage := errors.New("the otp is used")
		return errMessage
	} else if err == nil && !dest.Used {
		errMessage := errors.New("the otp exists")
		return errMessage
	}

	err = cache.Set[OtpCred](otpService.redisClient, key, val, time.Second*3600)

	if err != nil {
		return err
	}
	return nil
}

func (otpService *OtpService) ValidateOtp(mobile string, otp string) error {
	var prefix string = "redis"
	key := fmt.Sprintf("%s:%s", prefix, mobile)

	dest, err := cache.Get[OtpCred](otpService.redisClient, key)

	if err != nil {
		return err
	} else if dest.Used {
		errMessage := errors.New("the otp is used")
		return errMessage
	} else if !dest.Used && dest.Value != otp {
		errMessage := errors.New("the otp is not valid")
		return errMessage
	} else if !dest.Used && dest.Value == otp {
		dest.Used = true
		err = cache.Set[OtpCred](otpService.redisClient, key, dest, time.Second*3600)

		if err != nil {
			return err
		}
	}

	return nil
}
