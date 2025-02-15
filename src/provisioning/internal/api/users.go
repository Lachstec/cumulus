package api

import (
	"net/http"
	"strconv"

	"github.com/Lachstec/mc-hosting/internal/types"
	"github.com/gin-gonic/gin"
)

/*

	router.PATCH("/users/:userid", func(c *gin.Context) {
		userid, err := urlParamToInt64(c.Param("userid"))
		if err != nil {
			l.Warn().Err(err).Msg("invalid payload for updating user")
			c.String(http.StatusUnprocessableEntity, "invalid user format")
			return
		}

		users, err := userService.ReadUserByUserID(userid)
		if err != nil {
			l.Warn().Err(err).Msg("failed to retrieve user from database")
			c.String(http.StatusInternalServerError, "failed to retrieve user")
			return
		}

		if len(users) == 0 {
			c.String(http.StatusNotFound, "no user with given id exists")
			return
		}

		user := users[0]
		user.ID = userid
		err = c.BindJSON(&user)
		if err != nil {
			l.Error().Err(err).Msg("user returned from database does not match expected schema")
			c.String(http.StatusInternalServerError, "failed to retrieve user")
			return
		}
		updated, err := userService.UpdateUser(user)
		if err != nil {
			l.Warn().Err(err).Int64("userid", user.ID).Msg("failed to update user in database")
			c.String(http.StatusInternalServerError, "failed to update user")
			return
		}
		c.JSON(http.StatusOK, updated)
	})

	router.DELETE("/users/:userid", func(c *gin.Context) {
		userid, err := urlParamToInt64(c.Param("userid"))
		if err != nil {
			l.Warn().Err(err).Msg("invalid payload for deleting user")
			c.String(http.StatusUnprocessableEntity, "invalid user format")
			return
		}
		users, err := userService.ReadUserByUserID(userid)
		if err != nil {
			l.Warn().Err(err).Msg("failed to retrieve user from database")
			c.String(http.StatusInternalServerError, "failed to retrieve user")
			return
		}
		if len(users) == 0 {
			c.String(http.StatusNotFound, "no user with given id exists")
			return
		}
		user := users[0]
		err = userService.DeleteUser(user)
		if err != nil {
			l.Warn().Err(err).Msg("failed to delete user from database")
			c.String(http.StatusInternalServerError, "failed to delete user")
			return
		}
		c.Status(http.StatusNoContent)
	})

	router.GET("/users/:userid/servers", func(c *gin.Context) {
		userid, err := urlParamToInt64(c.Param("userid"))
		if err != nil {
			l.Warn().Err(err).Msg("invalid payload for user")
			c.String(http.StatusUnprocessableEntity, "invalid user format")
			return
		}
		servers, err := serverService.ReadServerByUserID(userid)
		if err != nil {
			l.Warn().Err(err).Msg("failed to fetch servers for given user")
			c.String(http.StatusInternalServerError, "failed to delete user")
			return
		}
		c.JSON(http.StatusOK, servers)
	})
*/

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
	var user *types.User
	err := c.BindJSON(user)
	if err != nil {
		h.Logger.Warn().Err(err).Msg("invalid payload for new user")
		h.respondError(c, http.StatusUnprocessableEntity, "user cannot be created from request", err.Error())
		return
	}

	userid, err := h.UserService.CreateUser(user)
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
