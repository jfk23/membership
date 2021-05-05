package model

import (
	"time"
)

type User struct {
	ID          int
	FirstName   string
	LastName    string
	Email       string
	Password    string
	AccessLevel int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Room struct {
	ID        int
	RoomName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Restriction struct {
	ID              int
	RestrictionName string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type Reservation struct {
	ID        int
	FirstName string
	LastName  string
	Email     string
	Phone     string
	RoomID    int
	StartDate time.Time
	EndDate   time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	Room      Room
	Processed int
}

type RoomRestriction struct {
	ID            int
	RoomID        int
	ReservationID int
	RestrictionID int
	StartDate     time.Time
	EndDate       time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Room          Room
	Restriction   Restriction
	Reservation   Reservation
}

type MailData struct {
	To       string
	From     string
	Subject  string
	Content  string
	Template string
}

// KCPC membership models

type Member struct {
	ID               int
	MemberClass      int
	CommunityGroupID int
	SmallGroupID     int
	KORName          string
	ENGName          string
	Address          string
	Email            string
	Phone            string
	FamilyMembers    string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
