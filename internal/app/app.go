package app

import (
	"brillian_voice_back/internal/domain/services/game"
	"log"

	handlers "brillian_voice_back/internal/handler/rest/game"
	"brillian_voice_back/internal/infrustucture/provider"
	"github.com/gin-gonic/gin"
)

func Run() {
	e := gin.Default()

	p := provider.NewProvider()
	s := game.NewGameService(p)

	h := handlers.NewHandler(s)

	h.Register(e)
	
	if err := e.Run(":8080"); err != nil {
		log.Fatalln(err)
	}
}
