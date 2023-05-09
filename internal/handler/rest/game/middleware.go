package game

import (
	"brillian_voice_back/internal/domain/dto"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

const (
	JoinDtoKey string = "join-dto"
)

func (h *Handler) WebSocketRequired(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}

func (h *Handler) FindRoom(c *fiber.Ctx) error {
	body := JoinInput{}
	if err := c.ParamsParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	if err := c.QueryParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	r, err := h.service.FindRoomByCode(c.Context(), body.Code)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	joinDto := dto.InputJoinGameDto{
		Username: body.Username,
		ID:       body.ID,
		Code:     body.Code,
		Room:     r,
	}

	c.Locals(JoinDtoKey, joinDto)

	return c.Next()
}
