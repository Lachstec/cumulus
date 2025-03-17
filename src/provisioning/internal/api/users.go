package api

import (
	"net/http"
	"strconv"

	"github.com/Lachstec/mc-hosting/internal/types"
	"github.com/gin-gonic/gin"
)

func (h *Handler) GetUsers(c *gin.Context) {
	users, err := h.UserService.ReadAllUsers()

	if err != nil {
		h.Logger.Error().Err(err).Msg("failed to fetch users from database")
		h.respondError(c, http.StatusInternalServerError, "failed to retrieve users from database", err.Error())
		return
	}

	h.respondSuccess(c, http.StatusOK, users)
}

func (h *Handler) CreateUser(c *gin.Context) {
	var user types.User
	err := BindJSONStrict(c, &user)
	if err != nil {
		h.Logger.Warn().Err(err).Msg("invalid payload for new user")
		h.respondError(c, http.StatusUnprocessableEntity, "user cannot be created from request", err.Error())
		return
	}

	userid, err := h.UserService.CreateUser(&user)
	if err != nil {
		h.Logger.Warn().Err(err).Msg("invalid payload for new user")
		h.respondError(c, http.StatusInternalServerError, "failed to create user", err.Error())
		return
	}

	h.respondSuccess(c, http.StatusCreated, userid)
}

func (h *Handler) GetUserById(c *gin.Context) {

	param := c.Param("userid")

	userid, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		h.Logger.Warn().Err(err).Msg("failed to extract user id from request")
		h.respondError(c, http.StatusBadRequest, "expected user id in url param", err.Error())
		return
	}

	users, err := h.UserService.ReadUserByUserID(userid)
	if err != nil {
		h.Logger.Warn().Err(err).Msg("failed to retrieve user from database")
		h.respondError(c, http.StatusInternalServerError, "failed to retrieve user from database", err.Error())
		return
	}
	if len(users) == 0 {
		h.respondError(c, http.StatusNotFound, "no user with given id exists", nil)
		return
	}
	user := users[0]
	h.respondSuccess(c, http.StatusOK, user)
}

func (h *Handler) UpdateUserById(c *gin.Context) {

	param := c.Param("userid")

	userid, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		h.Logger.Warn().Err(err).Msg("failed to extract user id from request")
		h.respondError(c, http.StatusBadRequest, "expected user id in url param", err.Error())
		return
	}

	users, err := h.UserService.ReadUserByUserID(userid)
	if err != nil {
		h.Logger.Warn().Err(err).Msg("failed to retrieve user from database")
		h.respondError(c, http.StatusInternalServerError, "failed to retrieve user", err.Error())
		return
	}

	if len(users) == 0 {
		h.respondError(c, http.StatusNotFound, "no user with given id.", nil)
		return
	}

	user := users[0]
	user.ID = userid
	err = BindJSONStrict(c, user)

	if err != nil {
		h.Logger.Error().Err(err).Msg("user returned from database does not match expected schema")
		h.respondError(c, http.StatusInternalServerError, "failed to retrieve user", err.Error())
		return
	}

	updated, err := h.UserService.UpdateUser(user)
	if err != nil {
		h.Logger.Warn().Err(err).Int64("userid", user.ID).Msg("failed to update user in database")
		h.respondError(c, http.StatusInternalServerError, "failed to update user", err.Error())
		return
	}

	h.respondSuccess(c, http.StatusOK, updated)
}

func (h *Handler) DeleteUserById(c *gin.Context) {

	param := c.Param("userid")

	userid, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		h.Logger.Warn().Err(err).Msg("failed to extract user id from request")
		h.respondError(c, http.StatusBadRequest, "expected user id in url param", err.Error())
		return
	}

	users, err := h.UserService.ReadUserByUserID(userid)
	if err != nil {
		h.Logger.Warn().Err(err).Msg("failed to retrieve user from database")
		h.respondError(c, http.StatusInternalServerError, "failed to retrieve user", err.Error())
		return
	}

	if len(users) == 0 {
		h.respondError(c, http.StatusNotFound, "no user with given id exists", nil)
		return
	}

	user := users[0]
	err = h.UserService.DeleteUser(user)
	if err != nil {
		h.Logger.Warn().Err(err).Msg("failed to delete user from database")
		h.respondError(c, http.StatusInternalServerError, "failed to delete user", err.Error())
		return
	}

	h.respondSuccess(c, http.StatusOK, "user deleted")
}

func (h *Handler) ServersOfUser(c *gin.Context) {

	param := c.Param("userid")

	userid, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		h.Logger.Warn().Err(err).Msg("failed to extract user id from request")
		h.respondError(c, http.StatusBadRequest, "expected user id in url param", err.Error())
		return
	}

	servers, err := h.ServerService.ReadServerByUserID(userid)
	if err != nil {
		h.Logger.Warn().Err(err).Msg("failed to fetch servers for given user")
		h.respondError(c, http.StatusInternalServerError, "failed to retrieve servers", err.Error())
		return
	}

	h.respondSuccess(c, http.StatusOK, servers)
}
