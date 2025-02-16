package api

import (
	"net/http"
	"strconv"

	"github.com/Lachstec/mc-hosting/internal/types"
	"github.com/gin-gonic/gin"
)

func (h *Handler) GetServers(c *gin.Context) {
	servers, err := h.ServerService.ReadAllServers()

	if err != nil {
		h.Logger.Warn().Err(err).Msg("failed to fetch servers from database")
		h.respondError(c, http.StatusInternalServerError, "failed to retrieve servers", err.Error())
		return
	}

	h.respondSuccess(c, http.StatusOK, servers)
}

func (h *Handler) CreateServer(c *gin.Context) {
	var server types.Server
	err := BindJSONStrict(c, &server)
	if err != nil {
		h.Logger.Warn().Err(err).Msg("invalid server payload")
		h.respondError(c, http.StatusBadRequest, "invalid server payload", err.Error())
		return
	}

	user := types.User{
		ID:    1,
		Sub:   "Samplesub",
		Name:  "Sampleuser",
		Class: types.Admin.Value(),
	}

	srv, err := h.Provisioner.NewGameServer(c, &server, &user)
	if err != nil {
		h.Logger.Error().Err(err).Int64("user_id", user.ID).Msg("failed to create game server")
		h.respondError(c, http.StatusInternalServerError, "failed to create new game server", err.Error())
	}

	h.respondSuccess(c, http.StatusCreated, srv.ID)
}

func (h *Handler) GetServerById(c *gin.Context) {

	param := c.Param("serverid")

	serverid, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		h.Logger.Warn().Err(err).Msg("failed to extract server id from request")
		h.respondError(c, http.StatusBadRequest, "expected server id in url param", err.Error())
		return
	}

	servers, err := h.ServerService.ReadServerByServerID(serverid)
	if err != nil {
		h.Logger.Warn().Err(err).Msg("failed to retrieve server from database")
		h.respondError(c, http.StatusInternalServerError, "failed to retrieve server", err.Error())
		return
	}

	if len(servers) == 0 {
		h.respondError(c, http.StatusNotFound, "no server with given id", nil)
		return
	}

	server := servers[0]
	h.respondSuccess(c, http.StatusOK, server)
}

func (h *Handler) StartServerById(c *gin.Context) {

	param := c.Param("serverid")

	serverid, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		h.Logger.Warn().Err(err).Msg("failed to extract server id from request")
		h.respondError(c, http.StatusBadRequest, "expected server id in url param", err.Error())
		return
	}

	servers, err := h.ServerService.ReadServerByServerID(serverid)
	if err != nil {
		h.Logger.Warn().Err(err).Msg("failed to retrieve server from database")
		h.respondError(c, http.StatusInternalServerError, "failed to retrieve server", err.Error())
		return
	}

	if len(servers) == 0 {
		h.respondError(c, http.StatusNotFound, "no server with given id", nil)
		return
	}

	server := servers[0]
	if server.Status != types.Stopped {
		h.respondError(c, http.StatusBadRequest, "server already running/starting", nil)
		return
	}

	h.respondSuccess(c, http.StatusOK, server)
}

func (h *Handler) Put(c *gin.Context) {
	h.respondError(c, http.StatusBadRequest, "not implemented", nil)
}

func (h *Handler) UpdateServerById(c *gin.Context) {
	param := c.Param("serverid")

	serverid, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		h.Logger.Warn().Err(err).Msg("failed to extract server id from request")
		h.respondError(c, http.StatusBadRequest, "expected server id in url param", err.Error())
		return
	}

	servers, err := h.ServerService.ReadServerByServerID(serverid)
	if err != nil {
		h.Logger.Warn().Err(err).Msg("failed to retrieve server from database")
		h.respondError(c, http.StatusInternalServerError, "failed to retrieve server", err.Error())
		return
	}

	if len(servers) == 0 {
		h.respondError(c, http.StatusNotFound, "no server with given id", nil)
		return
	}

	server := servers[0]

	err = BindJSONStrict(c, server)
	if err != nil {
		h.Logger.Warn().Err(err).Msg("invalid server payload")
		h.respondError(c, http.StatusUnprocessableEntity, "server payload not valid", err.Error())
	}

	srv, err := h.ServerService.UpdateServer(server)
	if err != nil {
		h.Logger.Warn().Err(err).Msg("failed to update server from database")
		h.respondError(c, http.StatusInternalServerError, "failed to update server", nil)
		return
	}

	h.respondSuccess(c, http.StatusOK, srv)
}

func (h *Handler) DeleteServerById(c *gin.Context) {
	param := c.Param("serverid")

	serverid, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		h.Logger.Warn().Err(err).Msg("failed to extract server id from request")
		h.respondError(c, http.StatusBadRequest, "expected server id in url param", err.Error())
		return
	}

	servers, err := h.ServerService.ReadServerByServerID(serverid)
	if err != nil {
		h.Logger.Warn().Err(err).Msg("failed to retrieve server from database")
		h.respondError(c, http.StatusInternalServerError, "failed to retrieve server", err.Error())
		return
	}

	if len(servers) == 0 {
		h.respondError(c, http.StatusNotFound, "no server with given id", nil)
		return
	}

	server := servers[0]
	err = h.ServerService.DeleteServer(server)

	if err != nil {
		h.Logger.Warn().Err(err).Msg("failed to delete server from database")
		h.respondError(c, http.StatusInternalServerError, "failed to delete server", err.Error())
		return
	}

	h.respondSuccess(c, http.StatusNoContent, "server was successfully deleted")
}

func (h *Handler) GetIpByServerId(c *gin.Context) {
	param := c.Param("serverid")

	serverid, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		h.Logger.Warn().Err(err).Msg("failed to extract server id from request")
		h.respondError(c, http.StatusBadRequest, "expected server id in url param", err.Error())
		return
	}

	ip, err := h.FloatingIPService.ReadIpByServerID(serverid)

	if err != nil {
		h.Logger.Warn().Err(err).Msg("failed to retrieve floating ip for server from database")
		h.respondError(c, http.StatusInternalServerError, "failed to retrieve server ip address", err.Error())
		return
	}

	h.respondSuccess(c, http.StatusOK, ip)
}
