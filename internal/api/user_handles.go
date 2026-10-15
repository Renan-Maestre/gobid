package api

import (
	"errors"
	"net/http"

	"github.com/Renan-Maestre/gobid/internal/jsonutils"
	"github.com/Renan-Maestre/gobid/internal/services"
	"github.com/Renan-Maestre/gobid/internal/usecase/user"
)

func (api *Api) handleSingupUser(w http.ResponseWriter, r *http.Request) {
	data, problems, err := jsonutils.DecodeValidJson[user.CreateUserReq](r)
	if err != nil {
		_ = jsonutils.EncodeJson(w, r, http.StatusUnprocessableEntity, problems)
		return
	}

	id, err := api.UserService.CreateUser(r.Context(), data.UserName, data.Email, data.Password, data.Bio)
	if err != nil {
		if errors.Is(err, services.ErrDuplicateEmailOrUsername) {
			_ = jsonutils.EncodeJson(w, r, http.StatusUnprocessableEntity, map[string]any{
				"error": "email our username already exists",
			})
			return
		}
		return
	}
	_ = jsonutils.EncodeJson(w, r, http.StatusCreated, map[string]any{
		"user_id": id,
	})
}

func (api *Api) handleLoginUser(w http.ResponseWriter, r *http.Request) {
	data, problems, err := jsonutils.DecodeValidJson[user.LoginUerReq](r)
	if err != nil {
		jsonutils.EncodeJson(w, r, http.StatusUnprocessableEntity, problems)
	}
	id, err := api.UserService.AuthenticateUser(r.Context(), data.Email, data.Password)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			jsonutils.EncodeJson(w, r, http.StatusBadRequest, map[string]any{"error": "invalid email or password"})
			return
		}
		jsonutils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{"error": "unexpected internal server error"})
		return
	}

	if err = api.Sessions.RenewToken(r.Context()); err != nil {
		jsonutils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{"error": "unexpected internal server error"})
		return
	}

	api.Sessions.Put(r.Context(), "AuthenticateUserId", id)
	jsonutils.EncodeJson(w, r, http.StatusOK, map[string]any{"message": "logged in successfully"})
}

func (api *Api) handleLogoutUser(w http.ResponseWriter, r *http.Request) {
	err := api.Sessions.RenewToken(r.Context())
	if err != nil {
		jsonutils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{"error": "unexpected internal server error"})
		return
	}
	api.Sessions.Remove(r.Context(), "AuthenticateUserId")
	jsonutils.EncodeJson(w, r, http.StatusOK, map[string]any{"message": "logged out successfully"})
}
