package game

import (
	"brillian_voice_back/internal/domain/dto"
	"brillian_voice_back/internal/domain/entity/game"
	gameSrv "brillian_voice_back/internal/domain/services/game"
	"brillian_voice_back/internal/infrustucture/conn"
	"errors"
	"fmt"
	"github.com/fasthttp/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"net/http"
)

const (
	CreatePath   = "/create"
	JoinPath     = "/ws/join/:code"
	WsPrefixPath = "/ws"
)

var upgrader = websocket.FastHTTPUpgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Handler struct {
	service *gameSrv.Service
}

func NewHandler(service *gameSrv.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Register(e *fiber.App) {
	h.registerCreateRoom(e)
	h.registerJoinRoom(e)
}

func (h *Handler) registerCreateRoom(e *fiber.App) {
	e.Post(CreatePath, func(c *fiber.Ctx) error {
		return h.createHandle(c)
	})
}

func (h *Handler) registerJoinRoom(e *fiber.App) {
	e.Use(WsPrefixPath, h.WebSocketRequired)
	e.Use(JoinPath, h.FindRoom)
	e.Get(JoinPath, func(c *fiber.Ctx) error {
		return h.joinHandle(c)
	})
}

func (h *Handler) createHandle(c *fiber.Ctx) error {
	var body CreateInput
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if err := body.Validate(); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	r, err := h.service.CreateRoom(c.Context(), &dto.InputCreateGameDto{
		Username:          body.Username,
		ID:                body.ID,
		CountPlayers:      body.CountPlayers,
		TimeDurationRound: body.TimeDurationSec,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	r.Run()
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":   "successfully created the room", //todo mv to const msgs
		"join_code": r.Desc().Code,
	})
}

func (h *Handler) joinHandle(c *fiber.Ctx) error {
	d, ok := c.Context().Value(JoinDtoKey).(dto.InputJoinGameDto)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%s", errors.New("error type cast"))
		return errors.New("error type cast")
	}

	ws, err := upgrader.Upgrade()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%s", err)
		return err
	}

	c := conn.NewPlayerConn(ws, d.Room.ActionChannel())
	u := game.NewUser(d.ID, d.Username, c.Adapter())
	c.SetContextInfo(d.Room.Desc(), u)
	ws.SetCloseHandler(func(code int, text string) error {
		if d.Room != nil {
			return d.Room.LeaveUser(u)
		}
		return errors.New("room is nil")
	})
	log.Info().
		Str("id", d.ID).
		Str("username", d.Username).
		Str("room", d.Room.Desc().Code).Msg("created new game and joined to room")
	if err := d.Room.JoinToRoom(u); err != nil {
		log.Error().
			Err(err).
			Str("id", d.ID).
			Str("room_code", d.Room.Desc().Code).Msg("an error occurred while joining the room")
		_ = c.Close()
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%s", err)
		return err
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "successfuly room created")
	return nil
}
