package dbrepo

import (
	"context"
	"errors"
	"time"

	"github.com/jfk23/gobookings/internal/model"
	"golang.org/x/crypto/bcrypt"
)

func (m *postgresDBRepo) AllUsers() bool {
	return true
}

func (m *postgresDBRepo) InsertRoomRestriction(res model.RoomRestriction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `insert into room_restrictions (start_date, end_date, room_id, reservation_id,
			 restriction_id, created_at, updated_at) values
			 ($1, $2, $3, $4, $5, $6, $7)`
	_, err := m.DB.ExecContext(ctx, query, res.StartDate, res.EndDate, res.RoomID, res.ReservationID,
		res.RestrictionID, time.Now(), time.Now())

	if err != nil {
		print("inserting data into DB failed!!!")
		return err
	}

	return nil
}

func (m *postgresDBRepo) SearchAvailabilityByDateByRoomID(roomID int, startDate, endDate time.Time) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var matchedRows int

	query := `select count(id) from room_restrictions where room_id = $1 and end_date > $2 and start_date < $3`
	rows := m.DB.QueryRowContext(ctx, query, roomID, startDate, endDate)
	err := rows.Scan(&matchedRows)

	if err != nil {
		return false, err
	}

	if matchedRows == 0 {
		return true, nil
	}

	return false, nil
}

func (m *postgresDBRepo) SearchAvailabilityByDateAll(startDate, endDate time.Time) ([]model.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var rooms []model.Room

	query := `select r.id, r.room_name
	from rooms r
	where r.id not in (select rr.id 
	from room_restrictions rr
	where $1 < rr.end_date and $2 > rr.start_date);`

	rows, err := m.DB.QueryContext(ctx, query, startDate, endDate)
	if err != nil {
		return rooms, err
	}

	for rows.Next() {
		var room model.Room

		err = rows.Scan(&room.ID, &room.RoomName)

		if err != nil {
			return rooms, err
		}

		rooms = append(rooms, room)
	}

	err = rows.Err()

	if err != nil {
		return rooms, nil
	}

	return rooms, nil
}

func (m *postgresDBRepo) GetRoomByID(id int) (model.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var room model.Room

	query := `select id, room_name, created_at, updated_at from rooms where id=$1`

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&room.ID, &room.RoomName, &room.CreatedAt, &room.UpdatedAt)
	if err != nil {
		return room, err
	}

	return room, nil

}

func (m *postgresDBRepo) GetUserByID(id int) (model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id, first_name, last_name, email, password, access_level, created_at, updated_at from users where id=$1`

	row := m.DB.QueryRowContext(ctx, query, id)

	var r model.User

	err := row.Scan(&r.ID, &r.FirstName, &r.LastName, &r.Email, &r.Password, &r.AccessLevel, &r.CreatedAt, &r.UpdatedAt)

	if err != nil {
		return r, err
	}

	return r, nil

}

func (m *postgresDBRepo) UpdateUser(u model.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `update users set first_name=$1, last_name=$2, email=$3, access_level=$4, updated_at=$5`

	_, err := m.DB.ExecContext(ctx, query, u.FirstName, u.LastName, u.Email, u.AccessLevel, time.Now())

	if err != nil {
		return err
	}

	return nil

}

func (m *postgresDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := m.DB.QueryRowContext(ctx, "select id, password from users where email = $1", email)

	var id int
	var hashedPassword string

	err := row.Scan(&id, &hashedPassword)

	if err != nil {
		return 0, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(testPassword))

	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", errors.New("invalid password")
	} else if err != nil {
		return 0, "", err
	}

	return id, hashedPassword, nil
}

func (m *postgresDBRepo) AllReservations() ([]model.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservations []model.Reservation

	query := `select r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date, r.end_date,
			r.room_id, r.created_at, r.updated_at, r.processed, rm.id, rm.room_name
			from reservations r
			left join rooms rm on (r.room_id = rm.id)
			order by r.start_date asc`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return reservations, err
	}

	defer rows.Close()

	for rows.Next() {
		var i model.Reservation
		err = rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.Phone,
			&i.StartDate,
			&i.EndDate,
			&i.RoomID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Processed,
			&i.Room.ID,
			&i.Room.RoomName,
		)
		if err != nil {
			return reservations, err
		}

		reservations = append(reservations, i)
	}

	if err = rows.Err(); err != nil {
		return reservations, err
	}

	return reservations, nil

}

func (m *postgresDBRepo) AllNewReservations() ([]model.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservations []model.Reservation

	query := `select r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date, r.end_date,
			r.room_id, r.created_at, r.updated_at, rm.id, rm.room_name
			from reservations r
			left join rooms rm on (r.room_id = rm.id)
			where r.processed = 0
			order by r.start_date asc`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return reservations, err
	}

	defer rows.Close()

	for rows.Next() {
		var i model.Reservation
		err = rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.Phone,
			&i.StartDate,
			&i.EndDate,
			&i.RoomID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Room.ID,
			&i.Room.RoomName,
		)
		if err != nil {
			return reservations, err
		}

		reservations = append(reservations, i)
	}

	if err = rows.Err(); err != nil {
		return reservations, err
	}

	return reservations, nil

}

func (m *postgresDBRepo) GetReservationByID(id int) (model.Reservation, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var res model.Reservation

	query := `select r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date, r.end_date,
			r.room_id, r.created_at, r.updated_at, r.processed, rm.id, rm.room_name
			from reservations r
			left join rooms rm on (r.room_id = rm.id)
			where r.id = $1`

	row := m.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&res.ID,
		&res.FirstName,
		&res.LastName,
		&res.Email,
		&res.Phone,
		&res.StartDate,
		&res.EndDate,
		&res.RoomID,
		&res.CreatedAt,
		&res.UpdatedAt,
		&res.Processed,
		&res.Room.ID,
		&res.Room.RoomName,
	)
	if err != nil {
		return res, err
	}

	return res, nil

}

func (m *postgresDBRepo) UpdateReservations(u model.Reservation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `update reservations set first_name=$1, last_name=$2, email=$3, phone=$4, updated_at=$5
				where id=$6`

	_, err := m.DB.ExecContext(ctx, query, u.FirstName, u.LastName, u.Email, u.Phone, time.Now(), u.ID)

	if err != nil {
		print(err)
		return err
	}

	return nil
}

func (m *postgresDBRepo) DeleteReservationByID(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `delete from reservations where id = $1`

	_, err := m.DB.ExecContext(ctx, query, id)

	if err != nil {
		return err
	}

	return nil
}

func (m *postgresDBRepo) UpdateProcessedForReservation(id, processed int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `update reservations set processed = $1 where id = $2`

	_, err := m.DB.ExecContext(ctx, query, processed, id)

	if err != nil {
		return err
	}

	return nil
}

func (m *postgresDBRepo) AllRooms() ([]model.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var rooms []model.Room

	query := `select * from rooms
			order by room_name asc`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return rooms, err
	}

	defer rows.Close()

	for rows.Next() {
		var room model.Room

		err := rows.Scan(
			&room.ID,
			&room.RoomName,
			&room.CreatedAt,
			&room.UpdatedAt,
		)

		if err != nil {
			return rooms, err
		}

		rooms = append(rooms, room)
	}

	if err = rows.Err(); err != nil {
		return rooms, err
	}

	return rooms, nil
}

func (m *postgresDBRepo) GetRoomRestrictionByDate(roomID int, start, end time.Time) ([]model.RoomRestriction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var restrictions []model.RoomRestriction

	query := `select id, start_date, end_date, room_id, coalesce(reservation_id, 0), restriction_id from room_restrictions
			where ($1 < end_date and $2 >= start_date)
			and room_id = $3`

	rows, err := m.DB.QueryContext(ctx, query, start, end, roomID)
	if err != nil {
		return restrictions, err
	}

	defer rows.Close()

	for rows.Next() {
		var r model.RoomRestriction

		err := rows.Scan(
			&r.ID,
			&r.StartDate,
			&r.EndDate,
			&r.RoomID,
			&r.ReservationID,
			&r.RestrictionID,
		)

		if err != nil {
			return restrictions, err
		}

		restrictions = append(restrictions, r)
	}

	if err = rows.Err(); err != nil {
		return restrictions, err
	}

	return restrictions, nil
}

func (m *postgresDBRepo) InsertBlockForRoom(RoomID int, start time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `insert into room_restrictions (start_date, end_date, room_id, restriction_id, created_at, updated_at) 
			values($1, $2, $3, $4, $5, $6)`

	_, err := m.DB.ExecContext(ctx, query, start, start.AddDate(0, 0, 1), RoomID, 2, time.Now(), time.Now())
	if err != nil {
		return err
	}
	return nil
}

func (m *postgresDBRepo) InsertNewMember(RoomID int, start time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `insert into room_restrictions (start_date, end_date, room_id, restriction_id, created_at, updated_at) 
			values($1, $2, $3, $4, $5, $6)`

	_, err := m.DB.ExecContext(ctx, query, start, start.AddDate(0, 0, 1), RoomID, 2, time.Now(), time.Now())
	if err != nil {
		return err
	}
	return nil
}

func (m *postgresDBRepo) DeleteBlockByID(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `delete from room_restrictions where id = $1`

	_, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

func (m *postgresDBRepo) InsertReservation(res model.Reservation) (int, error) {
	print("insert db called")
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var ReservationID int

	query := `insert into reservations (first_name, last_name, email, phone,
			 start_date, end_date, room_id, created_at, updated_at) values
			 ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`
	err := m.DB.QueryRowContext(ctx, query, res.FirstName, res.LastName, res.Email, res.Phone,
		res.StartDate, res.EndDate, res.RoomID, time.Now(), time.Now()).Scan(&ReservationID)

	if err != nil {
		print("inserting data into DB failed!!!")
		return 0, err
	}

	return ReservationID, nil
}

// KCPC membership functions below

func (m *postgresDBRepo) InsertMember(member model.Member) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var MemberID int

	query := `insert into members (member_class, community_group_id, small_group_id, korean_name,
			 english_name, address, email, phone, family_members, created_at, updated_at) values
			 ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) returning id`
	err := m.DB.QueryRowContext(ctx, query, member.MemberClass, member.CommunityGroupID,
		member.SmallGroupID, member.KORName, member.ENGName, member.Address, member.Email,
		member.Phone, member.FamilyMembers, time.Now(), time.Now()).Scan(&MemberID)

	if err != nil {
		print("inserting memeber into DB failed!!!")
		return 0, err
	}

	return MemberID, nil
}

func (m *postgresDBRepo) AllMembers() ([]model.Member, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var members []model.Member

	query := `select *
			from members
			order by english_name asc`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return members, err
	}

	defer rows.Close()

	for rows.Next() {
		var i model.Member
		err = rows.Scan(
			&i.ID,
			&i.MemberClass,
			&i.CommunityGroupID,
			&i.SmallGroupID,
			&i.KORName,
			&i.ENGName,
			&i.Address,
			&i.Email,
			&i.Phone,
			&i.FamilyMembers,
			&i.CreatedAt,
			&i.UpdatedAt,
		)
		if err != nil {
			return members, err
		}

		members = append(members, i)
	}

	if err = rows.Err(); err != nil {
		return members, err
	}

	return members, nil

}
