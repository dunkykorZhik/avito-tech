package routes

import (
	"net/http"

	"github.com/dunkykorZhik/avito-tech/internal/entity"
	errs "github.com/dunkykorZhik/avito-tech/internal/errors"
	"github.com/dunkykorZhik/avito-tech/internal/service"
)

func GetInfo(service service.History) handlefuncWithError {
	return func(w http.ResponseWriter, r *http.Request) error {
		//get user from context -> call service -> send answer
		username, _ := r.Context().Value(userCtx).(string)
		if username == "" {
			return errs.ErrUnAuth
		}
		output, err := service.GetHistory(r.Context(), username)
		if err != nil {
			return err
		}
		return writeJSON(w, &output)
	}
}

func SendCoin(service service.Transfer) handlefuncWithError {
	type sendCoinRequest struct {
		ToUser string `json:"toUser" validate:"required,max=100"`
		Amount int64  `json:"amount" validate:"required,gt=0"`
	}
	var sendCoinPayload sendCoinRequest
	return func(w http.ResponseWriter, r *http.Request) error {
		//get user from context -> call service -> send answer
		username := r.Context().Value(userCtx).(string)
		if username == "" {
			return errs.ErrUnAuth
		}
		if err := readJSON(w, r, &sendCoinPayload); err != nil {
			return errs.WrapError(err, http.StatusBadRequest)
		}
		if err := Validate.Struct(sendCoinPayload); err != nil {
			return errs.WrapError(err, http.StatusBadRequest)
		}
		transfer := entity.Transfer{
			Sender:   username,
			Receiver: sendCoinPayload.ToUser,
			Amount:   sendCoinPayload.Amount,
		}
		return service.CreateTransfer(r.Context(), transfer)

	}
}

func BuyItem(service service.Inventory) handlefuncWithError {
	return func(w http.ResponseWriter, r *http.Request) error {
		//get user from context -> call service -> send answer
		username := r.Context().Value(userCtx).(string)
		if username == "" {
			return errs.ErrUnAuth
		}
		item_name := r.PathValue("item")
		return service.BuyItem(r.Context(), username, item_name)

	}

}

func Auth(service service.User) handlefuncWithError {
	type AuthRequest struct {
		Username string `json:"username" validate:"required,min=1,max=100"`
		Password string `json:"password" validate:"required,min=6,max=100"`
	}

	type envelope struct {
		Token string `json:"token"`
	}
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
		return writeJSON(w, &envelope{Token: token})

	}

}
