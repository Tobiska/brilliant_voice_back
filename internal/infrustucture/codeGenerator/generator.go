package codeGenerator

import (
	"brillian_voice_back/pkg/randString"
	"context"
)

type Generator struct {
	lenCodes int
}

func NewGenerator(len int) *Generator {
	return &Generator{
		lenCodes: len,
	}
}

func (g *Generator) Generate(_ context.Context) string {
	return randString.RandStringBytes(g.lenCodes)
}
