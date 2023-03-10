package app

import (
	"brillian_voice_back/internal/domain/services/game"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"os"

	handlers "brillian_voice_back/internal/handler/rest/game"
	"brillian_voice_back/internal/infrustucture/provider"
	"github.com/gin-gonic/gin"
)

func Run() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	e := gin.Default()

	p := provider.NewProvider()
	s := game.NewGameService(p)

	h := handlers.NewHandler(s)

	h.Register(e)

	log.Info().Str("host", ":8080").Msg("Running server...")
	if err := e.Run(":8080"); err != nil {
		log.Fatal().Err(err)
	}
}
