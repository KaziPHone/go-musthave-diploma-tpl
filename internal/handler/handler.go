package handler

import "github.com/KaziPHone/go-musthave-diploma-tpl/pkg/storage"

type Handler struct {
	Storage storage.GofferStorage
	JwtKey  []byte
}

func NewHandler(dsn string, jwtKey []byte) *Handler {
	return &Handler{
		Storage: *storage.NewGofferStorage(dsn),
		JwtKey:  jwtKey,
	}
}
