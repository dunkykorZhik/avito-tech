package routes

import (
	"net/http"

	"github.com/dunkykorZhik/avito-tech/docs"
	"github.com/dunkykorZhik/avito-tech/internal/service"
	"github.com/go-playground/validator/v10"
	httpSwagger "github.com/swaggo/http-swagger" // http-swagger middleware
	"go.uber.org/zap"
)

/*
/info  historyService
/send  transferService
/buy{item}  inventoryService

middleware
/auth


//error types:
-not enough balance
-no such user
-no such item
*/

type userInfo struct {
	id       int64
	username string
}

type handlefuncWithError func(w http.ResponseWriter, r *http.Request) error

var (
	Validate *validator.Validate
	userCtx  = "user"
)

func init() {
	Validate = validator.New(validator.WithRequiredStructEnabled())
}

// @title						API Avito Shop
// @version					1.0.0
// @description				API for managing shop transactions.
// @host						localhost:8080
// @BasePath					/api
// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
func AddRoutes(mux *http.ServeMux, services *service.Service, logger *zap.SugaredLogger) {
	loggingMw := LoggingMiddleware(logger)
	authMw := AuthMiddleware(services.User)
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/api"

	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)
	mux.Handle("POST /api/auth", loggingMw(Auth(services.User)))

	mux.Handle("GET /api/info", loggingMw(authMw(GetInfo(services.History))))
	mux.Handle("POST /api/sendCoin", loggingMw(authMw(SendCoin(services.Transfer))))
	mux.Handle("GET /api/buy/{item}", loggingMw(authMw(BuyItem(services.Inventory))))

}
