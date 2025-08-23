package routes

import (
	"net/http"

	"github.com/dunkykorZhik/avito-tech/internal/entity"
	errs "github.com/dunkykorZhik/avito-tech/internal/errors"
	"github.com/dunkykorZhik/avito-tech/internal/service"
)

// @Summary	Получить информацию о монетах, инвентаре и истории транзакций.
// @Security	BearerAuth
// @Success	200	{object}	service.InfoResponse
// @Failure	400	{object}	ErrorResponse
// @Failure	401	{object}	ErrorResponse
// @Failure	500	{object}	ErrorResponse
// @Router		/info [get]
func GetInfo(service service.History) handlefuncWithError {
	return func(w http.ResponseWriter, r *http.Request) error {
		//get user from context -> call service -> send answer
		user, ok := r.Context().Value(userCtx).(userInfo)
		if !ok || user.id == 0 {
			return errs.ErrUnAuth
		}
		infoResponse, err := service.GetHistory(r.Context(), user.id)
		if err != nil {
			return err
		}
		return writeJSON(w, http.StatusOK, &infoResponse)
	}
}

type SendCoinRequest struct {
	ToUser string `json:"toUser" validate:"required,max=100"`
	Amount int64  `json:"amount" validate:"required,gt=0"`
}

// @Summary	Отправить монеты другому пользователю.
// @Security	BearerAuth
// @Accept		json
// @Produce	json
// @Param		request	body	SendCoinRequest	true	"SendCoinRequest"
// @Success	200
// @Failure	400	{object}	ErrorResponse
// @Failure	401	{object}	ErrorResponse
// @Failure	500	{object}	ErrorResponse
// @Router		/sendCoin [post]
func SendCoin(service service.Transfer) handlefuncWithError {

	var sendCoinReq SendCoinRequest
	return func(w http.ResponseWriter, r *http.Request) error {
		//get user from context -> call service -> send answer
		user, ok := r.Context().Value(userCtx).(userInfo)
		if !ok || user.id == 0 {
			return errs.ErrUnAuth
		}
		if err := readJSON(w, r, &sendCoinReq); err != nil {
			return errs.WrapError(err, http.StatusBadRequest)
		}
		if err := Validate.Struct(sendCoinReq); err != nil {
			return errs.WrapError(err, http.StatusBadRequest)
		}

		transfer := entity.Transfer{
			SenderID:     user.id,
			ReceiverName: sendCoinReq.ToUser,
			Amount:       sendCoinReq.Amount,
		}
		if err := service.CreateTransfer(r.Context(), transfer); err != nil {
			return err
		}
		w.WriteHeader(http.StatusOK)
		return nil

	}
}

// @Summary	Купить предмет за монеты.
// @Security	BearerAuth
// @Param		item	path	string	true	"Item Name"
// @Success	200
// @Failure	400	{object}	ErrorResponse
// @Failure	401	{object}	ErrorResponse
// @Failure	500	{object}	ErrorResponse
// @Router		/buy/{item} [get]
func BuyItem(service service.Inventory) handlefuncWithError {
	return func(w http.ResponseWriter, r *http.Request) error {
		//get user from context -> call service -> send answer
		user, ok := r.Context().Value(userCtx).(userInfo)
		if !ok || user.id == 0 {
			return errs.ErrUnAuth
		}
		item_name := r.PathValue("item")
		if err := service.BuyItem(r.Context(), user.id, item_name); err != nil {
			return err
		}
		w.WriteHeader(http.StatusOK)
		return nil

	}

}

type AuthRequest struct {
	Username string `json:"username" validate:"required,min=1,max=100"`
	Password string `json:"password" validate:"required,min=6,max=100"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

// @Summary	Аутентификация и получение JWT-токена.
// @Accept		json
// @Produce	json
// @Param		request	body		AuthRequest	true	"AuthRequest"
// @Success	200		{object}	AuthResponse
// @Failure	400		{object}	ErrorResponse
// @Failure	401		{object}	ErrorResponse
// @Failure	500		{object}	ErrorResponse
// @Router		/auth [post]
func Auth(service service.User) handlefuncWithError {

	var authRequest AuthRequest
	return func(w http.ResponseWriter, r *http.Request) error {

		if err := readJSON(w, r, &authRequest); err != nil {
			return errs.WrapError(err, http.StatusBadRequest)
		}
		if err := Validate.Struct(authRequest); err != nil {
			return errs.WrapError(err, http.StatusBadRequest)
		}
		token, err := service.GenerateToken(r.Context(), authRequest.Username, authRequest.Password)
		if err != nil {
			return err
		}
		return writeJSON(w, http.StatusOK, &AuthResponse{Token: token})

	}

}
