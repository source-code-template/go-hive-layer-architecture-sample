package app

import (
	"context"
	"reflect"

	. "github.com/beltran/gohive"
	v "github.com/core-go/core/v10"

	"github.com/core-go/hive"
	"github.com/core-go/log"
	"github.com/core-go/search/hive"
	"go-service/internal/handler"
	"go-service/internal/model"
	"go-service/internal/service"
)

type ApplicationContext struct {
	User handler.UserPort
}

func NewApp(ctx context.Context, conf Config) (*ApplicationContext, error) {
	configuration := NewConnectConfiguration()
	configuration.Database = "masterdata"
	connection, errConn := Connect(conf.Hive.Host, conf.Hive.Port, conf.Hive.Auth, configuration)
	if errConn != nil {
		return nil, errConn
	}

	logError := log.LogError
	validator := v.NewValidator()

	userType := reflect.TypeOf(model.User{})
	userQuery := query.NewBuilder("users", userType)
	userSearchBuilder, err := hive.NewSearchBuilder(connection, userType, userQuery.BuildQuery)
	if err != nil {
		return nil, err
	}
	userRepository, err := hive.NewWriter(connection, "users", userType)
	if err != nil {
		return nil, err
	}
	userService := service.NewUserService(userRepository)
	userHandler := handler.NewUserHandler(userSearchBuilder.Search, userService, validator.Validate, logError)

	return &ApplicationContext{
		User: userHandler,
	}, nil
}
