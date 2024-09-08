package server

import (
	"encoding/json"
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/jdks/fiber-example/internal/models"
	"github.com/jdks/fiber-example/internal/store"
)

const defaultPageSize = 100

func (s Server) GetUser(fiberCtx *fiber.Ctx) error {
	userID := fiberCtx.Params("user_id")
	s.log.Info().Str("user_id", userID).Msg("Getting user")
	user, err := s.store.GetUser(fiberCtx.Context(), userID)
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to get user")
		return fiberCtx.JSON(models.User{})
	}

	return fiberCtx.JSON(models.User{
		ID:        user.ID.String(),
		FirstName: user.FirstName,
		LastName:  user.LastName,
	})
}

func (s Server) GetAllUsers(fiberCtx *fiber.Ctx) error {
	pageNumber := fiberCtx.QueryInt("page_number", 1)
	pageSize := fiberCtx.QueryInt("page_size", defaultPageSize)
	users, err := s.store.GetAllUsers(fiberCtx.Context(), pageSize, pageNumber)
	if err != nil {
		s.log.Error().Err(err).Msg("")
		return fiberCtx.JSON(users)
	}
	resp := make([]models.User, len(users))
	for i, user := range users {
		resp[i] = models.User{
			ID:        user.ID.String(),
			FirstName: user.FirstName,
			LastName:  user.LastName,
		}
	}

	return fiberCtx.JSON(resp)
}

func (s Server) QueryEvents(fiberCtx *fiber.Ctx) error {
	action := new(store.EventAction)
	err := action.UnmarshalJSON([]byte(fiberCtx.Query("action")))
	if err != nil {
		s.log.Error().Err(err).Msg("")
		return fiberCtx.JSON([]models.UserEvent{})
	}

	pageNumber := fiberCtx.QueryInt("page_number", 1)
	pageSize := fiberCtx.QueryInt("page_size", defaultPageSize)
	params := store.EventQueryParams{
		EventID:    fiberCtx.Query("event_id"),
		SessionID:  fiberCtx.Query("session_id"),
		UserID:     fiberCtx.Query("user_id"),
		DayInWeek:  fiberCtx.QueryInt("day_in_week"),
		DayInMonth: fiberCtx.QueryInt("day_in_month"),
		Month:      fiberCtx.QueryInt("month"),
		Year:       fiberCtx.QueryInt("year"),
		Hour:       fiberCtx.QueryInt("hour"),
	}
	if action != nil {
		params.Action = *action
	}
	userEvents, err := s.store.QueryEvents(fiberCtx.Context(), params, pageSize, pageNumber)
	if err != nil {
		s.log.Error().Err(err).Msg("")
	}

	resp := make([]models.UserEvent, len(userEvents))
	for i, userEvent := range userEvents {
		associatedUserIDs := make([]string, len(userEvent.AssociatedUserIDs))
		for j, id := range userEvent.AssociatedUserIDs {
			associatedUserIDs[j] = id.String()
		}
		action, _ := json.Marshal(userEvent.Payload.Action)
		resp[i] = models.UserEvent{
			ID:                userEvent.EventID.String(),
			UserID:            userEvent.UserID.String(),
			CreatedAt:         userEvent.CreatedAt,
			Action:            strings.Trim(string(action), `"`),
			SessionID:         userEvent.Payload.SessionID.String(),
			AssociatedUserIDs: associatedUserIDs,
		}
	}
	return fiberCtx.JSON(resp)
}
