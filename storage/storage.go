package storage

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/boltdb/bolt"
	"github.com/hthl85/flexapi-mock-server/model"
)

// Storage defines storage structure
type Storage struct {
	BoltDB *bolt.DB
	bucket []byte
}

// NewStorage inits new storage
func NewStorage(dbName, bucketName string) *Storage {
	db, err := bolt.Open(dbName, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}

	bucket := []byte(bucketName)

	err = db.Update(func(tx *bolt.Tx) error {
		_ = tx.DeleteBucket(bucket)

		_, err := tx.CreateBucketIfNotExists(bucket)
		if err != nil {
			log.Fatalf("Create bucket failed: %s", err)
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Create bucket failed: %s", err)
	}

	return &Storage{BoltDB: db, bucket: bucket}
}

// AddNewUser add new user to database
func (s *Storage) AddNewUser(user *model.User) error {
	// Serialize the struct to JSON
	jsonBytes, _ := json.Marshal(user)

	// Generate a key
	key := strconv.Itoa(user.ID)

	// Write the data to the AccountBucket
	err := s.BoltDB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(s.bucket)
		err := b.Put([]byte(key), jsonBytes)
		return err
	})
	if err != nil {
		return err
	}

	return nil
}

// GetUserByID gets user by userid from database
func (s *Storage) GetUserByID(userID int) (*model.User, error) {
	// Allocate an empty user instance
	var user model.User

	// Generate a key
	key := strconv.Itoa(userID)

	// Read an object from the bucket using BoltDB.View
	err := s.BoltDB.View(func(tx *bolt.Tx) error {
		// Read the bucket from the DB
		b := tx.Bucket(s.bucket)

		// Read the value identified by our userId supplied as []byte
		accountBytes := b.Get([]byte(key))
		if accountBytes == nil {
			return fmt.Errorf("No user with id %d", userID)
		}

		// Unmarshal the returned bytes into the user struct we created at
		// the top of the function
		json.Unmarshal(accountBytes, &user)

		// Return nil to indicate nothing went wrong, e.g no error
		return nil
	})

	// If there were an error, return the error
	if err != nil {
		return nil, err
	}
	// Return the user struct and nil as error.
	return &user, nil
}

// GetUsersByIDs gets user by userid from database
func (s *Storage) GetUsersByIDs(userIDs []int) ([]model.User, error) {
	// Allocate an empty user instance
	var users []model.User

	for _, userID := range userIDs {
		// Generate a key
		key := strconv.Itoa(userID)

		var user model.User

		// Read an object from the bucket using BoltDB.View
		err := s.BoltDB.View(func(tx *bolt.Tx) error {
			// Read the bucket from the DB
			b := tx.Bucket(s.bucket)

			// Read the value identified by our userId supplied as []byte
			accountBytes := b.Get([]byte(key))
			if accountBytes == nil {
				return fmt.Errorf("No user with id %d", userID)
			}

			// Unmarshal the returned bytes into the user struct we created at
			// the top of the function
			json.Unmarshal(accountBytes, &user)

			// Return nil to indicate nothing went wrong, e.g no error
			return nil
		})

		// If there were an error, return the error
		if err == nil {
			users = append(users, user)
		}
	}

	// Return the user struct and nil as error.
	return users, nil
}

// GetAllUsers gets all user from database
func (s *Storage) GetAllUsers() ([]*model.User, error) {

	var users []*model.User

	// Read an object from the bucket using BoltDB.View
	err := s.BoltDB.View(func(tx *bolt.Tx) error {
		// Read the bucket from the DB
		b := tx.Bucket(s.bucket)

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var user model.User
			json.Unmarshal(v, &user)
			users = append(users, &user)
		}

		// Return nil to indicate nothing went wrong, e.g no error
		return nil
	})

	// If there were an error, return the error
	if err != nil {
		return nil, err
	}

	// Return the user struct and nil as error.
	return users, nil
}

// UpdateUser updates user in database
func (s *Storage) UpdateUser(user *model.User) error {
	return s.ReplaceUser(user)
}

// ReplaceUser replaces user in database
func (s *Storage) ReplaceUser(user *model.User) error {
	// Serialize the struct to JSON
	jsonBytes, _ := json.Marshal(user)

	// Generate a key
	key := strconv.Itoa(user.ID)

	// Write the data to the AccountBucket
	err := s.BoltDB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(s.bucket)
		err := b.Put([]byte(key), jsonBytes)
		return err
	})
	if err != nil {
		return err
	}

	return nil
}

// DeleteUserByID deletes user in database
func (s *Storage) DeleteUserByID(userID int) error {
	// Generate a key
	key := strconv.Itoa(userID)

	return s.BoltDB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(s.bucket)
		return bucket.Delete([]byte(key))
	})
}
