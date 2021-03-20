package repository

import (
	"time"
	"github.com/jfk23/gobookings/internal/model"
)

type DatabaseRepo interface {
	AllUsers() bool
	InsertReservation(res model.Reservation) (int, error)
	InsertRoomRestriction(res model.RoomRestriction) error
	SearchAvailabilityByDateByRoomID(roomID int, startDate, endDate time.Time) (bool, error)
	SearchAvailabilityByDateAll (startDate, endDate time.Time) ([]model.Room, error)
	GetRoomByID (id int) (model.Room, error)

	GetUserByID (id int) (model.User, error)
	UpdateUser (u model.User) error
	Authenticate (email, testPassword string) (int, string, error)
	AllReservations() ([]model.Reservation, error)
	AllNewReservations() ([]model.Reservation, error)
	GetReservationByID(id int) (model.Reservation, error)
	UpdateReservations (u model.Reservation) error
	DeleteReservationByID (id int) error
	UpdateProcessedForReservation (id, processed int) error
	AllRooms() ([]model.Room, error)
	GetRoomRestrictionByDate(roomID int, start, end time.Time) ([]model.RoomRestriction, error)
	InsertBlockForRoom(RoomID int, start time.Time) error
	DeleteBlockByID(id int) error
}