package dbrepo

import (
	"errors"
	"time"

	"github.com/jfk23/gobookings/internal/model"
)

func (m *testDBRepo) AllUsers() bool {
	return true
}

func (m *testDBRepo) testDBRepo(res model.Reservation) (int, error) {

	return 1, nil
}

func (m *testDBRepo) InsertReservation(res model.Reservation) (int, error) {
	//print("test insert db called ^^")
	if res.RoomID == 2 {
		return 0, errors.New("some error")
	}

	return 1, nil
}

func (m *testDBRepo) InsertRoomRestriction(res model.RoomRestriction) error {
	if res.RoomID == 1000 {
		return errors.New("some error")
	}

	return nil
}

func (m *testDBRepo) SearchAvailabilityByDateByRoomID(roomID int, startDate, endDate time.Time) (bool, error) {

	return false, nil
}

func (m *testDBRepo) SearchAvailabilityByDateAll(startDate, endDate time.Time) ([]model.Room, error) {
	var rooms []model.Room

	return rooms, nil
}

func (m *testDBRepo) GetRoomByID(id int) (model.Room, error) {
	var room model.Room

	if id > 2 {
		return room, errors.New("there is no room where id is greater than 2")
	}

	return room, nil

}

func (m *testDBRepo) GetUserByID(id int) (model.User, error) {
	var u model.User

	return u, nil
}

func (m *testDBRepo) UpdateUser(u model.User) error {
	return nil
}

func (m *testDBRepo) Authenticate(email, testPassword string) (int, string, error) {

	if email == "admin@a.a" {
		return 0, "", nil
	}
	return 0, "", errors.New("invaild credentials")
}

func (m *testDBRepo) AllReservations() ([]model.Reservation, error) {

	var reservations []model.Reservation

	return reservations, nil

}

func (m *testDBRepo) AllNewReservations() ([]model.Reservation, error) {

	var reservations []model.Reservation
	return reservations, nil
}

func (m *testDBRepo) GetReservationByID(id int) (model.Reservation, error) {
	var res model.Reservation
	return res, nil
}

func (m *testDBRepo) UpdateReservations(u model.Reservation) error {

	return nil
}

func (m *testDBRepo) DeleteReservationByID(id int) error {

	return nil
}

func (m *testDBRepo) UpdateProcessedForReservation(id, processed int) error {

	return nil
}

func (m *testDBRepo) AllRooms() ([]model.Room, error) {

	var rooms []model.Room
	return rooms, nil
}

func (m *testDBRepo) GetRoomRestrictionByDate(roomID int, start, end time.Time) ([]model.RoomRestriction, error) {
	var restrictions []model.RoomRestriction
	return restrictions, nil
}

func (m *testDBRepo) InsertBlockForRoom(RoomID int, start time.Time) error {
	return nil
}

func (m *testDBRepo) DeleteBlockByID(id int) error {
	return nil
}

func (m *testDBRepo) InsertMember(member model.Member) (int, error) {
	return 0, nil
}

func (m *testDBRepo) AllMembers() ([]model.Member, error) {
	var members []model.Member
	return members, nil
}
