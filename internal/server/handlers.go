package server

import (
	"github.com/gofiber/fiber/v2"

	"github.com/jdks/fiber-example/internal/models"
	"github.com/jdks/fiber-example/internal/store"
)

const defaultPageSize = 100

func (s Server) GetUser(fiberCtx *fiber.Ctx) error {
	userID := fiberCtx.Params("user_id")
	u, err := s.store.GetUser(fiberCtx.Context(), userID)
	if err != nil {
		s.log.Error().Err(err).Msg("")
		return fiberCtx.JSON(models.User{})
	}

	return fiberCtx.JSON(models.User{
		ID:        u.ID.String(),
		FirstName: u.FirstName,
		LastName:  u.LastName,
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
	us := make([]models.User, len(users))
	for i, u := range users {
		us[i] = models.User{
			ID:        u.ID.String(),
			FirstName: u.FirstName,
			LastName:  u.LastName,
		}
	}

	return fiberCtx.JSON(users)
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

	ues := make([]models.UserEvent, len(userEvents))
	associatedUserIDs := make([]string, len(userEvents[0].AssociatedUserIDs))
	for i, ue := range userEvents {
		for j, id := range ue.AssociatedUserIDs {
			associatedUserIDs[j] = id.String()
		}
		ues[i] = models.UserEvent{
			ID:                ue.EventID.String(),
			UserID:            ue.UserID.String(),
			CreatedAt:         ue.CreatedAt,
			Action:            string(ue.Payload.Action),
			SessionID:         ue.Payload.SessionID.String(),
			AssociatedUserIDs: associatedUserIDs,
		}
	}
	return fiberCtx.JSON(ues)
}
