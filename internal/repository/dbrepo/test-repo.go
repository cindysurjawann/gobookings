package dbrepo

import (
	"errors"
	"time"

	"github.com/cindysurjawann/gobookings/internal/models"
)

func (m *testDBRepo) AllUsers() bool {
	return true
}

// inserts a reservation into the database
func (m *testDBRepo) InsertReservation(res models.Reservation) (int, error) {
	//if the room id is 2, then fail; otherwise, pass
	if res.RoomID == 2 {
		return 0, errors.New("error in testDBRepo.InsertReservation")
	}
	return 1, nil
}

// insert a room restriction into the database
func (m *testDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	if r.RoomID == 1000 {
		return errors.New("error in testDBRepo.InsertRoomRestriction")
	}
	return nil
}

// returns true if availability exists for roomid, and false if no availaibility
func (m *testDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error) {
	return false, nil
}

// return a slice of available rooms, if any for given date range
func (m *testDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	var rooms []models.Room
	return rooms, nil
}

// get a room by ID
func (m *testDBRepo) GetRoomByID(id int) (models.Room, error) {
	var room models.Room
	if id > 2 {
		return room, errors.New("some error")
	}

	return room, nil
}
