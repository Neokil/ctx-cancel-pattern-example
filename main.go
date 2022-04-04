package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

const (
	TIMEOUT_GENERAL  = 1 * time.Minute
	TIMEOUT_API_CALL = 10 * time.Second
	TIMEOUT_DB_QUERY = 10 * time.Second

	// Those constats will be used by the mocks to determine how long their corresponding functions are taking to return
	// You can play around with those to test scenarios like long-running API-Calls or slow DB-Queries and see how they
	// behave with the timeout-context
	API_MOCK_TAKES                = 1 * time.Second
	DB_MOCK_TAKES_FOR_GET_USERID  = 1 * time.Second
	DB_MOCK_TAKES_FOR_UPDATE_USER = 1 * time.Second
	DB_MOCK_TAKES_FOR_CREATE_USER = 1 * time.Second
)

func main() {
	fmt.Println("Main started")
	defer fmt.Println("Main ended")

	// Context with our general Timeout
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT_GENERAL)
	defer cancel()

	users, err := getUsersFromAPI(ctx)
	if err != nil {
		panic(err)
	}

	err = updateDatabase(ctx, users)
	if err != nil {
		panic(err)
	}
}

func updateDatabase(ctx context.Context, users []ApiUser) error {
	fmt.Println("updateDatabase started")
	defer fmt.Println("updateDatabase ended")

	for _, u := range users {
		err := updateOrCreateUser(ctx, u)
		if err != nil {
			return fmt.Errorf("Failed to update or create user: %w", err)
		}
	}
	return nil
}

func updateOrCreateUser(ctx context.Context, user ApiUser) error {
	fmt.Println("updateOrCreateUser started")
	defer fmt.Println("updateOrCreateUser ended")

	// derrive a new context with a new timeout and make sure it gets canceled
	// when we get out of scope so it can be cleaned up
	ctx, cancel := context.WithTimeout(ctx, TIMEOUT_DB_QUERY)
	defer cancel()

	// Create a general result channel and defer the close
	rc := make(chan error, 1)
	defer close(rc)

	go func() {
		db := &DB{}
		id, err := db.getUserIdByEmail(user.Email)
		if errors.Is(err, UserNotFoundError) {
			err = db.createUser(user.Firstname, user.Lastname, user.Email)
			if err != nil {
				rc <- fmt.Errorf("Failed to create User: %w", err)
				return
			}
			rc <- nil
			return
		}
		if err != nil {
			rc <- fmt.Errorf("Failed to get User-ID: %w", err)
			return
		}
		err = db.updateUser(id, user.Firstname, user.Lastname, user.Email)
		if err != nil {
			rc <- fmt.Errorf("Failed to update User: %w", err)
			return
		}
		rc <- nil
		return
	}()

	select {
	case err := <-rc:
		if err != nil {
			return fmt.Errorf("Failed to save or update Userin DB: %w", err)
		}
		return nil
	case <-ctx.Done():
		return fmt.Errorf("Context was terminated in updateOrCreateUser: %w", ctx.Err())
	}
}

func getUsersFromAPI(ctx context.Context) ([]ApiUser, error) {
	fmt.Println("getUsersFromAPI started")
	defer fmt.Println("getUsersFromAPI ended")

	// derrive a new context with a new timeout
	ctx, cancel := context.WithTimeout(ctx, TIMEOUT_API_CALL)
	defer cancel()

	// Create a result channel and defer the close
	rc := make(chan []ApiUser, 1)
	defer close(rc)
	// Create an error channel and defer the close
	ec := make(chan error, 1)
	defer close(ec)

	go func() {
		api := &API{}
		u, err := api.getList()
		if err != nil {
			ec <- err
			return
		}
		rc <- u
	}()

	select {
	case r := <-rc:
		return r, nil
	case err := <-ec:
		return nil, fmt.Errorf("Failed to retrieve Users from API: %w", err)
	case <-ctx.Done():
		return nil, fmt.Errorf("Context was terminated in getUsersFromAPI: %w", ctx.Err())
	}
}
