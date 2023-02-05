package game

import (
	"brillian_voice_back/internal/domain/dto"
	"brillian_voice_back/internal/domain/entity/user"
	"brillian_voice_back/internal/domain/services/game"
	"brillian_voice_back/internal/infrustucture/conn"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/schema"
	"github.com/gorilla/websocket"
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
		func(w http.ResponseWriter, r *http.Request) {
			h.createHandle(c)
		}(c.Writer, c.Request)
	})
}

func (h *Handler) registerJoinRoom(e *gin.Engine) {
	e.GET(JoinPath, func(c *gin.Context) {
		func(w http.ResponseWriter, r *http.Request) {
			h.joinHandle(c)
		}(c.Writer, c.Request)
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
		"join_code": r.Desc().Code,
	})
}

var decoder = schema.NewDecoder()

func (h *Handler) joinHandle(c *gin.Context) {
	var body JoinInput
	if err := decoder.Decode(&body, c.Request.URL.Query()); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	r, err := h.service.JoinToRoom(c.Request.Context(), &dto.InputJoinGameDto{
		Username: body.Username,
		ID:       body.ID,
		Code:     body.Code,
	})
	if err != nil {
		c.String(http.StatusNotFound, err.Error())
		return
	}

	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	u := user.NewUser(body.ID, body.Username, conn.NewPlayerConn(ws, r))
	if err := r.JoinToRoom(u); err != nil { //todo (может получится что owner игру создаст но присоединиться не сможет)
		_ = u.Close() //todo mv to service layer
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "successfully joined to the room",
	})
}
