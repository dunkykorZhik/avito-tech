package routes

import (
	"net/http"

	"github.com/dunkykorZhik/avito-tech/internal/service"
	"github.com/go-playground/validator/v10"
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

var (
	Validate *validator.Validate
	userCtx  string = "username"
)

type handlefuncWithError func(w http.ResponseWriter, r *http.Request) error

func init() {
	Validate = validator.New(validator.WithRequiredStructEnabled())
}
func AddRoutes(mux *http.ServeMux, services *service.Service, logger *zap.SugaredLogger) {
	loggingMw := LoggingMiddleware(logger)
	authMw := AuthMiddleware(services.User)

	mux.Handle("POST /api/auth", loggingMw(Auth(services.User)))

	mux.Handle("GET /api/info", loggingMw(authMw(GetInfo(services.History))))
	mux.Handle("POST /api/sendCoin", loggingMw(authMw(SendCoin(services.Transfer))))
	mux.Handle("GET /api/buy/{item}", loggingMw(authMw(BuyItem(services.Inventory))))

}
