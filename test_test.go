package main

import (
	"testing"
	"time"
)

func setupTestStore(t *testing.T) *SqliteStore {
	store, err := NewSqliteStoreTesting()
	if err != nil {
		t.Fatalf("Failed to create a new sqlite store: %v", err)
	}
	if err := store.Init(); err != nil {
		t.Fatalf("Failed to init the sqlite store: %v", err)
	}
	return store
}
func TestCreateUser(t *testing.T) {
	store := setupTestStore(t)
	defer store.db.Close()
	user, err := NewUser(123, "testuser", "password")
	if err != nil {
		t.Fatalf("Failed to create user object: %v", err)
	}
	if err := store.CreateUser(*user); err != nil {
		t.Fatalf("Failed to insert user to db: %v", err)
	}
	// test if user is existed after creating
	var count int
	err = store.db.QueryRow("SELECT COUNT(*) FROM users WHERE id = ?", user.ID).Scan(&count)
	if err != nil {
		t.Fatalf("Failed to query user: %v", err)
	}
	if count != 1 {
		t.Fatalf("Expected 1 user, got %d", count)
	}
}
func TestGetUser(t *testing.T) {
	store := setupTestStore(t)
	defer store.db.Close()
	user, err := NewUser(123, "testuser", "password")
	if err != nil {
		t.Fatalf("Failed to create user object: %v", err)
	}
	if err := store.CreateUser(*user); err != nil {
		t.Fatalf("Failed to insert user to db: %v", err)
	}
	// test if user is existed after creating
	fetchedUser, err := store.GetUser(user.ID)
	if err != nil {
		t.Fatalf("Failed to get user: %v", err)
	}
	if fetchedUser.Username != user.Username {
		t.Fatalf("Expected username %s, got %s", user.Username, fetchedUser.Username)
	}
}
func TestGetUsers(t *testing.T) {
	store := setupTestStore(t)
	defer store.db.Close()
	var users []User
	for i := 0; i < 10; i++ {
		user, err := NewUser(int64(i), "testuser", "password")
		if err != nil {
			t.Fatalf("Failed to create user object %v", err)
		} else {
			users = append(users, *user)
		}
	}

	for _, user := range users {
		if err := store.CreateUser(user); err != nil {
			t.Fatalf("Failed to insert user to db: %v", err)
		}
	}

	fetchedUsers, err := store.GetUsers()
	if err != nil {
		t.Fatalf("Failed to get users: %v", err)
	}

	if len(fetchedUsers) != len(users) {
		t.Fatalf("Expected %d users, got %d", len(users), len(fetchedUsers))
	}
}
func TestGetUserByUserID(t *testing.T) {
	store := setupTestStore(t)
	defer store.db.Close()
	user, err := NewUser(123, "testuser", "password")
	if err != nil {
		t.Fatalf("Failed to create user object: %v", err)
	}
	if err := store.CreateUser(*user); err != nil {
		t.Fatalf("Failed to insert user to db: %v", err)
	}
	fetchedUser, err := store.GetUserByUserID(123)
	if err != nil {
		t.Fatalf("Failed to query user with userID: %v", err)
	}
	if fetchedUser.UserID != user.UserID {
		t.Fatalf("Expected userID %v, got %v", user.UserID, fetchedUser.UserID)
	}
}
func TestUpdateUser(t *testing.T) {
	store := setupTestStore(t)
	user, err := NewUser(123, "testuser", "password")
	if err != nil {
		t.Fatalf("Failed to create user object: %v", err)
	}
	if err := store.CreateUser(*user); err != nil {
		t.Fatalf("Failed to insert user to db: %v", err)
	}
	// create test for chenging username
	updatedUsername := "updateduser"
	user.Username = updatedUsername
	if err = store.UpdateUser(user.ID, *user); err != nil {
		t.Fatalf("Failed to update user: %v", err)
	}
	fetchedUser, err := store.GetUser(user.ID)
	if err != nil {
		t.Fatalf("Failed to get user: %v", err)
	}
	if fetchedUser.Username != updatedUsername {
		t.Fatalf("Expected username %s, got %s", updatedUsername, fetchedUser.Username)
	}
	// TODO create other test for updating feilds
}
func TestDeleteUser(t *testing.T) {
	store := setupTestStore(t)
	defer store.db.Close()
	user, err := NewUser(123, "testuser", "password")
	if err != nil {
		t.Fatalf("Failed to create user object: %v", err)
	}
	if err := store.CreateUser(*user); err != nil {
		t.Fatalf("Failed to insert user to db: %v", err)
	}

	if err := store.DeleteUser(user.ID); err != nil {
		t.Fatalf("Failed to delete user: %v", err)
	}
	var count int
	if err := store.db.QueryRow("SELECT COUNT(*) FROM users WHERE id = ?", user.ID).Scan(&count); err != nil {
		t.Fatalf("Failed to query user: %v", err)
	}

	if count != 0 {
		t.Fatalf("Expected 0 users, got %d", count)
	}
}
func TestCreateAlert(t *testing.T) {
	store := setupTestStore(t)
	defer store.db.Close()
	alert := NewAlert(123, "testalert", time.Time{})
	if err := store.CreateAlert(*alert); err != nil {
		t.Fatalf("Failed to insert alert to db: %v", err)
	}
	var count int
	if err := store.db.QueryRow("SELECT COUNT(*) FROM alerts WHERE id = ?", alert.ID).Scan(&count); err != nil {
		t.Fatalf("Failed to query alert: %v", err)
	}
	if count != 1 {
		t.Fatalf("Expected 1 user, got %d", count)
	}
}

func TestGetAlert(t *testing.T) {
	store := setupTestStore(t)
	defer store.db.Close()
	alert := NewAlert(123, "testalert", time.Time{})
	if err := store.CreateAlert(*alert); err != nil {
		t.Fatalf("Failed to insert alert to db: %v", err)
	}
	fetchedAlert, err := store.GetAlert(alert.ID)
	if err != nil {
		t.Fatalf("Faild to query alert: %v", err)
	}
	if fetchedAlert.UserID != alert.UserID {
		t.Fatalf("Expected userID %v, got %v", alert.UserID, fetchedAlert.UserID)
	}
	if fetchedAlert.Lable != alert.Lable {
		t.Fatalf("Expected label %s, got %s", alert.Lable, fetchedAlert.Lable)
	}
}

func TestGetAlertByUserID(t *testing.T) {
	store := setupTestStore(t)
	defer store.db.Close()
	testingAlerts := []Alert{
		*NewAlert(123, "testalert", time.Time{}),
		*NewAlert(123, "testalert", time.Time{}),
		*NewAlert(12345, "testalert", time.Time{}),
		*NewAlert(12345, "testalert", time.Time{}),
		*NewAlert(12345, "testalert", time.Time{}),
	}
	for _, alert := range testingAlerts {
		if err := store.CreateAlert(alert); err != nil {
			t.Fatalf("Failed to insert alert to db: %v", err)
		}
	}

	fetchedAlerts, err := store.GetAlertsByUserID(123)
	if err != nil {
		t.Fatalf("Faild to query alert: %v", err)
	}
	if len(fetchedAlerts) != 2 {
		t.Fatalf("Expected alerts length is %d, got %d", 2, len(fetchedAlerts))
	}
}

func TestGetAlerts(t *testing.T) {
	store := setupTestStore(t)
	defer store.db.Close()
	testingAlerts := []Alert{
		*NewAlert(123, "testalert", time.Time{}),
		*NewAlert(123, "testalert", time.Time{}),
		*NewAlert(12345, "testalert", time.Time{}),
		*NewAlert(12345, "testalert", time.Time{}),
		*NewAlert(12345, "testalert", time.Time{}),
	}
	for _, alert := range testingAlerts {
		if err := store.CreateAlert(alert); err != nil {
			t.Fatalf("Failed to insert alert to db: %v", err)
		}
	}
	fetchedAlerts, err := store.GetAlerts()
	if err != nil {
		t.Fatalf("Failed to query alerts: %v", err)
	}
	if len(fetchedAlerts) != 5 {
		t.Fatalf("Expected alerts length is %d, got %d", 2, len(fetchedAlerts))
	}
}

func TestUpdateAlert(t *testing.T) {

}

func TestDeleteAlert(t *testing.T) {
	store := setupTestStore(t)
	defer store.db.Close()
	alert := NewAlert(123, "testalert", time.Time{})
	if err := store.CreateAlert(*alert); err != nil {
		t.Fatalf("Failed to insert alert to db: %v", err)
	}
	if err := store.DeleteAlert(alert.ID); err != nil {
		t.Fatalf("Failed to delete alert: %v", err)
	}
	var count int
	if err := store.db.QueryRow("SELECT COUNT(*) FROM alerts WHERE id = ?", alert.ID).Scan(&count); err != nil {
		t.Fatalf("Failed to query alert: %v", err)
	}

	if count != 0 {
		t.Fatalf("Expected 0 alerts, got %d", count)
	}

}
