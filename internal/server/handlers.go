package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"

	"github.com/jdks/fiber-example/internal/models"
)

const pageSize = 10

func (s Server) GetUser(fiberCtx *fiber.Ctx) error {
	u, err := s.store.GetUser(fiberCtx.Context(), fiberCtx.Params("user_id"))
	if err != nil {
		log.Error().Err(err).Msg("")
		fiberCtx.JSON(models.User{})
	}

	return fiberCtx.JSON(models.User{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
	})
}

func (s Server) GetAllUsers(fiberCtx *fiber.Ctx) error {
	users, err := s.store.GetAllUsers(fiberCtx.Context())
	if err != nil {
		log.Error().Err(err).Msg("")
		fiberCtx.JSON(users)
	}
	us := make([]models.User, len(users))
	for i, u := range users {
		us[i] = models.User{
			ID:        u.ID,
			FirstName: u.FirstName,
			LastName:  u.LastName,
		}
	}

	return fiberCtx.JSON(users)
}
