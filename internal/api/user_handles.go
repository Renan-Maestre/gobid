package api

import (
	"net/http"

	"github.com/Renan-Maestre/gobid/internal/jsonutils"
	"github.com/Renan-Maestre/gobid/internal/usecase/user"
)

func (api *Api) handleSingupUser(w http.ResponseWriter, r *http.Request) {
	_, problems, err := jsonutils.DecodeValidJson[user.CreateUserReq](r)
	if err != nil {
		_ = jsonutils.EncodeJson(w, r, http.StatusUnprocessableEntity, problems)
	}
}

func (api *Api) handleLoginUser(w http.ResponseWriter, r *http.Request) {
	panic("TODO - NOT INPLEMNTED")
}

func (api *Api) handleLogoutUser(w http.ResponseWriter, r *http.Request) {
	panic("TODO - NOT INPLEMNTED")
}
