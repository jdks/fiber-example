package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gofrs/uuid"

	"github.com/jdks/fiber-example/internal/config"
	"github.com/jdks/fiber-example/internal/store"
)

func main() {
	// Parse command line flags
	numUsers := flag.Int("users", 10, "Number of users to seed")
	numEvents := flag.Int("events", 100, "Number of events to seed")
	flag.Parse()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create a new store
	s, err := store.New(cfg.Store)
	if err != nil {
		log.Fatalf("Failed to create store: %v", err)
	}

	// Seed users
	users := seedUsers(s, *numUsers)

	// Seed events
	_ = seedEvents(s, users, *numEvents)

	fmt.Printf("Seeded %d users and %d events\n", *numUsers, *numEvents)
}

func seedUsers(s store.Store, count int) []store.User {
	users := make([]store.User, count)
	for i := 0; i < count; i++ {
		user := store.User{
			ID:        uuid.Must(uuid.NewV4()),
			FirstName: gofakeit.FirstName(),
			LastName:  gofakeit.LastName(),
		}
		users[i] = user
		err := s.CreateUser(context.Background(), user)
		if err != nil {
			log.Printf("Failed to create user %s: %v", user.ID, err)
		}
	}
	return users
}

func seedEvents(s store.Store, users []store.User, count int) []store.UserEvent {
	events := make([]store.UserEvent, count)
	for i := 0; i < count; i++ {
		user := users[rand.Intn(len(users))]
		associatedUser := users[rand.Intn(len(users))]
		createdAt := time.Now().Add(-time.Duration(rand.Intn(300*24*60*60*1000)) * time.Millisecond)
		event := store.UserEvent{
			EventID: uuid.Must(uuid.NewV4()),
			UserID:  user.ID,
			Payload: store.EventPayload{
				SessionID:  uuid.Must(uuid.NewV4()),
				Action:     store.EventActionMap[gofakeit.RandomString(store.EventActionLabels)],
				Hour:       uint(createdAt.Hour()),
				DayInWeek:  uint(createdAt.Weekday()),
				DayInMonth: uint(createdAt.Day()),
				Month:      uint(createdAt.Month()),
				Year:       uint(createdAt.Year()),
			},
			AssociatedUserIDs: []uuid.UUID{associatedUser.ID},
			CreatedAt:         createdAt,
		}
		err := s.CreateUserEvent(event)
		if err != nil {
			log.Printf("Failed to create event %s: %v", event.EventID, err)
		}
	}
	return events
}
