package game

import (
	"brillian_voice_back/internal/domain/dto"
	"brillian_voice_back/internal/domain/entity/user"
	"brillian_voice_back/internal/domain/services/game"
	"brillian_voice_back/internal/infrustucture/conn"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/schema"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"net/http"
)

const (
	CreatePath = "/create"
	JoinPath   = "/ws/join"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

type Handler struct {
	service *game.Service
}

func NewHandler(service *game.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Register(e *gin.Engine) {
	h.registerCreateRoom(e)
	h.registerJoinRoom(e)
}

func (h *Handler) registerCreateRoom(e *gin.Engine) {
	e.POST(CreatePath, func(c *gin.Context) {
		h.createHandle(c)
	})
}

func (h *Handler) registerJoinRoom(e *gin.Engine) {
	e.GET(JoinPath, func(c *gin.Context) {
		h.joinHandle(c.Writer, c.Request)
	})
}

func (h *Handler) createHandle(c *gin.Context) {
	var body CreateInput
	if err := c.ShouldBindJSON(&body); err != nil {
		c.String(http.StatusUnprocessableEntity, "Is not valid error: %s", err.Error())
		return
	}
	r, err := h.service.CreateRoom(c.Request.Context(), &dto.InputCreateGameDto{
		Username:          body.Username,
		ID:                body.ID,
		CountPlayers:      body.CountPlayers,
		TimeDurationRound: body.TimeDurationMin,
	})
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	r.Run()
	c.JSON(http.StatusCreated, gin.H{
		"msg":       "successfully created the room", //todo mv to const msgs
		"join_code": body.ID,
	})
}

var decoder = schema.NewDecoder()

func (h *Handler) joinHandle(w http.ResponseWriter, req *http.Request) {
	var body JoinInput
	if err := decoder.Decode(&body, req.URL.Query()); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "%s", err)
		return
	}
	r, err := h.service.JoinToRoom(req.Context(), &dto.InputJoinGameDto{
		Username: body.Username,
		ID:       body.ID,
		Code:     body.Code,
	})
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "%s", err)
		return
	}

	ws, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%s", err)
		return
	}

	u := user.NewUser(body.ID, body.Username, conn.NewPlayerConn(ws))
	ws.SetCloseHandler(func(code int, text string) error {
		return u.DeleteAndClose()
	})
	log.Info().
		Str("id", body.ID).
		Str("username", body.Username).
		Str("room", r.Desc().Code).Msg("created new user and joined to room")
	if err := r.JoinToRoom(u); err != nil { //todo (может получится что owner игру создаст но присоединиться не сможет)
		log.Error().
			Err(err).
			Str("id", body.ID).
			Str("room_code", r.Desc().Code).Msg("an error occurred while joining the room")
		_ = u.Close() //todo mv to service layer
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%s", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "successfuly room created")
}
