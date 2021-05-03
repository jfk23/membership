package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/jfk23/gobookings/cmd/web/helpers"
	"github.com/jfk23/gobookings/driver"
	"github.com/jfk23/gobookings/internal/config"
	"github.com/jfk23/gobookings/internal/forms"
	"github.com/jfk23/gobookings/internal/model"
	"github.com/jfk23/gobookings/internal/render"
	"github.com/jfk23/gobookings/internal/repository"
	"github.com/jfk23/gobookings/internal/repository/dbrepo"
)

var Repo *Repository
var existingTemplateData *model.TemplateData

type Repository struct {
	ConfigSetting *config.AppConfig
	DB            repository.DatabaseRepo
}

func CreateNewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		ConfigSetting: a,
		DB:            dbrepo.NewPostgresRepo(a, db.SQL),
	}
}

func CreateNewTestRepo(a *config.AppConfig) *Repository {
	return &Repository{
		ConfigSetting: a,
		DB:            dbrepo.NewtestingDBRepo(a),
	}
}

func SetHandler(r *Repository) {
	Repo = r
}

//Home is page for /
func (re *Repository) Home(rw http.ResponseWriter, r *http.Request) {

	render.Template(rw, r, "home.page.html", &model.TemplateData{})

}

func (re *Repository) AdminAddMember(rw http.ResponseWriter, r *http.Request) {
	render.Template(rw, r, "new_member.page.html", &model.TemplateData{
		Form: forms.New(nil),
	})

}

func (re *Repository) Reservation(rw http.ResponseWriter, r *http.Request) {
	res, ok := re.ConfigSetting.Session.Get(r.Context(), "reservation").(model.Reservation)
	if !ok {
		re.ConfigSetting.Session.Put(r.Context(), "error1", "there is no data from the session")
		http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)
		return
	}

	var stringMap = make(map[string]string)
	sd := res.StartDate.Format("2006-01-02")
	ed := res.EndDate.Format("2006-01-02")

	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	// getting room name calling function via room id passed
	room, err := re.DB.GetRoomByID(res.RoomID)

	if err != nil {
		re.ConfigSetting.Session.Put(r.Context(), "error", "there is no room found from the session")
		http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)
		return
	}

	res.Room.RoomName = room.RoomName

	re.ConfigSetting.Session.Put(r.Context(), "reservation", res)

	data := make(map[string]interface{})
	data["reservation"] = res

	render.Template(rw, r, "make-reservation.page.html", &model.TemplateData{
		Form:      forms.New(nil),
		DataMap:   data,
		StringMap: stringMap,
	})
}

func (re *Repository) PostReservation(rw http.ResponseWriter, r *http.Request) {
	var startDate time.Time
	var endDate time.Time
	var roomID int
	// log.Println("post reservation function called!!!")

	// res, ok := re.ConfigSetting.Session.Get(r.Context(), "reservation").(model.Reservation)
	// if !ok {
	// 	re.ConfigSetting.Session.Put(r.Context(), "error", "cann't get reservation from the session")
	// 	http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)
	// 	return
	// }

	// var res model.Reservation

	e := r.ParseForm()
	// err = errors.New("this is artificial error for testing")
	if e != nil {
		re.ConfigSetting.Session.Put(r.Context(), "error", "cannot parse the form")
		http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)
		return
	}

	var sd string
	var ed string

	if existingTemplateData == nil {

		//log.Println("existingTempdata called!!!")
		//set up values to store into DB
		//2020-01-20 --> Mon Jan 2 15:04:05 MST 2006
		sd = r.Form.Get("start_date")
		ed = r.Form.Get("end_date")

		var err error

		layout := "2006-01-02"
		startDate, err = time.Parse(layout, sd)
		if err != nil {
			re.ConfigSetting.Session.Put(r.Context(), "error", "cannot convert start time")
			http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)
			return
		}
		endDate, err = time.Parse(layout, ed)
		if err != nil {
			re.ConfigSetting.Session.Put(r.Context(), "error", "cannot convert end time")
			http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)
			return
		}

		roomID, err = strconv.Atoi(r.Form.Get("room_id"))
		if err != nil {
			re.ConfigSetting.Session.Put(r.Context(), "error", "cannot convert room id string into integer")
			http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)
			return
		}
		//print("room id being printed")
		//fmt.Printf("this is room id: %d\n", roomID)
	} else {
		//fmt.Println("this is temp data: ", existingTemplateData)
		//fmt.Println("skipping existingTempdata??")
		res, ok := re.ConfigSetting.Session.Get(r.Context(), "reservation").(model.Reservation)
		if !ok {
			//fmt.Println("no res data available, so for testing purpose use below dummy data")

			startDate, _ = time.Parse("2006-01-02", "2022-02-20")
			endDate, _ = time.Parse("2006-01-02", "2022-02-22")
			roomID = 1

			//print("game over here?????")
			//http.Error(rw, "invalid form data", http.StatusSeeOther)

		} else {
			startDate = res.StartDate
			endDate = res.EndDate
			roomID = res.RoomID
		}

	}

	// var res model.Reservation
	// res.StartDate = startDate
	// res.EndDate = endDate
	// res.RoomID = roomID

	// res.FirstName = r.Form.Get("first_name")
	// res.LastName = r.Form.Get("last_name")
	// res.Phone = r.Form.Get("phone")
	// res.Email = r.Form.Get("email")

	res := model.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Phone:     r.Form.Get("phone"),
		Email:     r.Form.Get("email"),
		RoomID:    roomID,
		StartDate: startDate,
		EndDate:   endDate,
	}

	startdate := startDate.Format("2006-01-02")
	enddate := endDate.Format("2006-01-02")

	stringMap := map[string]string{"start_date": startdate, "end_date": enddate}

	//fmt.Println("this is res: ", res)

	res.Room, _ = re.DB.GetRoomByID(res.RoomID)

	formdata := forms.New(r.PostForm)

	// formdata.Has("first_name", r)
	formdata.Required("first_name", "last_name", "email", "phone")
	formdata.MinLength("first_name", 4)
	formdata.EmailValidate("email")

	if !formdata.Valid() {
		//log.Println("invalid form function called!!!")
		data := make(map[string]interface{})
		data["reservation"] = res
		//http.Error(rw, "invalid form data", http.StatusSeeOthe1)

		existingTemplateData = &model.TemplateData{
			Form:      formdata,
			DataMap:   data,
			StringMap: stringMap,
		}
		re.ConfigSetting.Session.Put(r.Context(), "reservation", res)

		render.Template(rw, r, "make-reservation.page.html", existingTemplateData)
		return
	}
	//log.Println("just before insertreservation call")

	reservationID, e := re.DB.InsertReservation(res)
	//log.Println("just after insertreservation call")
	if e != nil {
		//log.Println("inserting data into DB failed!!! because room id =2")
		re.ConfigSetting.Session.Put(r.Context(), "error", "cannot insert reservation into DB")
		http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)
		return
	}

	restriction := model.RoomRestriction{
		RoomID:        res.RoomID,
		ReservationID: reservationID,
		StartDate:     res.StartDate,
		EndDate:       res.EndDate,
		RestrictionID: 1,
	}

	e = re.DB.InsertRoomRestriction(restriction)
	//log.Println("just before insert restriction call")
	if e != nil {
		//log.Println("inserting data into DB failed!!!")
		re.ConfigSetting.Session.Put(r.Context(), "error", "cannot insert room restriction into DB")
		http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)
		return
	}
	//log.Println("just after insert restriction call")

	htmlContent := fmt.Sprintf(
		`<strong> Reservation summary</strong><br>
		Name: %s <br>
		Arrival date: %s <br>
		Departure date: %s <br>
	`, res.FirstName, sd, ed)

	msg := model.MailData{
		To:      res.Email,
		From:    "good@friend.com",
		Subject: "Confirmation of your reservation",
		Content: htmlContent,
	}

	re.ConfigSetting.MailChan <- msg

	htmlContent = fmt.Sprintf(
		`<strong> Reservation summary from guest</strong><br>
		Client Name: %s <br>
		Room: %s <br>
		Arrival date: %s <br>
		Departure date: %s <br>
	`, res.FirstName, res.Room.RoomName, sd, ed)

	msg = model.MailData{
		To:       "owner@friend.com",
		From:     "good@friend.com",
		Subject:  "Confirmation of client reservation",
		Content:  htmlContent,
		Template: "drip.html",
	}

	re.ConfigSetting.MailChan <- msg

	re.ConfigSetting.Session.Put(r.Context(), "reservation", res)
	//log.Println("got to the statuseeother")
	http.Redirect(rw, r, "/reservation-summary", http.StatusSeeOther)

}

func (re *Repository) ReservationSummary(rw http.ResponseWriter, r *http.Request) {
	reservation, ok := re.ConfigSetting.Session.Get(r.Context(), "reservation").(model.Reservation)

	sd := reservation.StartDate.Format("2006-01-02")
	ed := reservation.EndDate.Format("2006-01-02")

	stringMap := map[string]string{"start_date": sd, "end_date": ed}

	if !ok {
		//re.ConfigSetting.ErrorLog.Println("error reading from session data")
		re.ConfigSetting.Session.Put(r.Context(), "error", "we can't get data from session")
		http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)
		return
	}
	re.ConfigSetting.Session.Remove(r.Context(), "reservation")

	data := make(map[string]interface{})
	data["reservation"] = reservation

	render.Template(rw, r, "reservation-summary.page.html", &model.TemplateData{
		DataMap:   data,
		StringMap: stringMap,
	})
}

func (re *Repository) Generals(rw http.ResponseWriter, r *http.Request) {
	render.Template(rw, r, "generals.page.html", &model.TemplateData{})
}

func (re *Repository) Majors(rw http.ResponseWriter, r *http.Request) {
	render.Template(rw, r, "majors.page.html", &model.TemplateData{})
}

func (re *Repository) Contact(rw http.ResponseWriter, r *http.Request) {
	render.Template(rw, r, "contact.page.html", &model.TemplateData{})
}

func (re *Repository) Search(rw http.ResponseWriter, r *http.Request) {
	render.Template(rw, r, "search.page.html", &model.TemplateData{})
}

func (re *Repository) PostSearch(rw http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")
	fmt.Println("this is raw start and end date: ", start, end)

	layout := "01/02/2006"
	startDate, err := time.Parse(layout, start)
	if err != nil {
		helpers.SeverError(rw, err)
	}
	endDate, err := time.Parse(layout, end)
	if err != nil {
		helpers.SeverError(rw, err)
	}

	rooms, err := re.DB.SearchAvailabilityByDateAll(startDate, endDate)
	if err != nil {
		helpers.SeverError(rw, err)
	}

	for _, v := range rooms {
		re.ConfigSetting.InfoLog.Println("Rooms available: ", v.ID, v.RoomName)
	}

	if len(rooms) == 0 {
		re.ConfigSetting.InfoLog.Println("No avail")
		re.ConfigSetting.Session.Put(r.Context(), "error", "There is no room available")
		http.Redirect(rw, r, "/search", http.StatusSeeOther)
	}

	reservation := model.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}

	re.ConfigSetting.Session.Put(r.Context(), "reservation", reservation)

	data := make(map[string]interface{})
	data["room"] = rooms

	render.Template(rw, r, "choose-room.page.html", &model.TemplateData{
		DataMap: data,
	})

	// rw.Write([]byte(fmt.Sprintf("start dat is %s, and end date is %s", start, end)))

}

type responseJson struct {
	OK        bool   `json:"ok"`
	Message   string `json:"message"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

func (re *Repository) PostSearchJson(rw http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()

	if err != nil {
		var resp = responseJson{
			OK:      false,
			Message: "we got error when parsing the form",
		}
		out, _ := json.MarshalIndent(resp, " ", "    ")
		rw.Header().Add("Content-Type", "application/json")
		rw.Write(out)
		return
	}

	sd := r.Form.Get("popup-start-date")
	ed := r.Form.Get("popup-end-date")

	layout := "01/02/2006"

	start, _ := time.Parse(layout, sd)
	end, _ := time.Parse(layout, ed)
	roomID, _ := strconv.Atoi(r.Form.Get("room_id"))

	available, er := re.DB.SearchAvailabilityByDateByRoomID(roomID, start, end)

	if er != nil {
		var resp = responseJson{
			OK:      false,
			Message: "we got error when searching room by ID",
		}
		out, _ := json.MarshalIndent(resp, " ", "    ")
		rw.Header().Add("Content-Type", "application/json")
		rw.Write(out)
		return
	}

	var resp = responseJson{
		OK:        available,
		Message:   "",
		StartDate: sd,
		EndDate:   ed,
	}

	out, err := json.MarshalIndent(resp, " ", "    ")

	// this should work!
	if err != nil {
		helpers.SeverError(rw, err)
		fmt.Println(err)
	}

	// fmt.Println(string(out))

	rw.Header().Add("Content-Type", "application/json")
	rw.Write(out)

}

// About is page for /
func (re *Repository) About(rw http.ResponseWriter, r *http.Request) {

	render.Template(rw, r, "about.page.html", &model.TemplateData{})

}

func (re *Repository) ChooseRoom(rw http.ResponseWriter, r *http.Request) {
	roomID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.SeverError(rw, err)
		return
	}

	// res := re.ConfigSetting.Session.Get(r.Context(), "reservation")

	res, ok := re.ConfigSetting.Session.Get(r.Context(), "reservation").(model.Reservation)
	if !ok {
		helpers.SeverError(rw, errors.New("got error when getting reservation model from session"))
		return
	}

	res.RoomID = roomID

	re.ConfigSetting.Session.Put(r.Context(), "reservation", res)

	http.Redirect(rw, r, "/make-reservation", http.StatusSeeOther)

}

func (re *Repository) BookRoom(rw http.ResponseWriter, r *http.Request) {
	roomID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		helpers.SeverError(rw, err)
		return
	}

	sd := r.URL.Query().Get("s")
	ed := r.URL.Query().Get("e")

	layout := "01/02/2006"

	start, _ := time.Parse(layout, sd)
	end, _ := time.Parse(layout, ed)

	room, err := re.DB.GetRoomByID(roomID)

	if err != nil {
		helpers.SeverError(rw, err)
		return
	}

	var res model.Reservation

	res.StartDate = start
	res.EndDate = end
	res.RoomID = roomID
	res.Room = room

	re.ConfigSetting.Session.Put(r.Context(), "reservation", res)
	http.Redirect(rw, r, "/make-reservation", http.StatusSeeOther)

}

func (re *Repository) ShowLogin(rw http.ResponseWriter, r *http.Request) {
	render.Template(rw, r, "login.page.html", &model.TemplateData{
		Form: forms.New(nil),
	})
}

func (re *Repository) PostShowLogin(rw http.ResponseWriter, r *http.Request) {
	//print("logging in ")
	_ = re.ConfigSetting.Session.RenewToken(r.Context())

	err := r.ParseForm()

	if err != nil {
		print("error with parseform()")
		log.Println(err)
	}

	forms := forms.New(r.PostForm)

	forms.Required("email", "password")
	forms.EmailValidate("email")

	if !forms.Valid() {
		//print("invalid form block called!")

		render.Template(rw, r, "login.page.html", &model.TemplateData{
			Form: forms,
		})
		return

	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	id, _, err := re.DB.Authenticate(email, password)

	if err != nil {
		//log.Println("here comes", err)
		re.ConfigSetting.Session.Put(r.Context(), "error", "invalid crendentials")
		// re.ConfigSetting.Session.Put(r.Context(), "flash", "invalid crendentials")

		// fmt.Println("session data", re.ConfigSetting.Session.Get(r.Context(), "error1"))
		// fmt.Println("ErrorMsg: ", re.ConfigSetting.Session.PopString(r.Context(), "error"))

		http.Redirect(rw, r, "/user/login", http.StatusSeeOther)
		return
	}

	re.ConfigSetting.Session.Put(r.Context(), "user_id", id)
	re.ConfigSetting.Session.Put(r.Context(), "flash", "hi logged in successfully")
	http.Redirect(rw, r, "/admin/dashboard", http.StatusSeeOther)

}

func (re *Repository) ShowLogout(rw http.ResponseWriter, r *http.Request) {
	_ = re.ConfigSetting.Session.Destroy(r.Context())
	_ = re.ConfigSetting.Session.RenewToken(r.Context())

	http.Redirect(rw, r, "/user/login", http.StatusSeeOther)
}

func (re *Repository) AdminDashBoard(rw http.ResponseWriter, r *http.Request) {
	render.Template(rw, r, "admin-dashboard.page.html", &model.TemplateData{})

}

func (re *Repository) AdminNewReservations(rw http.ResponseWriter, r *http.Request) {
	reservations, err := re.DB.AllNewReservations()
	if err != nil {
		helpers.SeverError(rw, err)
	}
	data := make(map[string]interface{})
	data["reservation"] = reservations

	render.Template(rw, r, "admin-new-reservations.page.html", &model.TemplateData{
		DataMap: data,
	})

}

func (re *Repository) AdminAllReservations(rw http.ResponseWriter, r *http.Request) {
	reservations, err := re.DB.AllReservations()
	if err != nil {
		helpers.SeverError(rw, err)
	}
	data := make(map[string]interface{})
	data["reservation"] = reservations

	render.Template(rw, r, "admin-all-reservations.page.html", &model.TemplateData{
		DataMap: data,
	})

}

func (re *Repository) AdminReservationsCalendar(rw http.ResponseWriter, r *http.Request) {
	now := time.Now()

	if r.URL.Query().Get("y") != "" {
		year, _ := strconv.Atoi(r.URL.Query().Get("y"))
		month, _ := strconv.Atoi(r.URL.Query().Get("m"))
		now = time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	}

	next := now.AddDate(0, 1, 0)
	last := now.AddDate(0, -1, 0)

	nextMonth := next.Format("01")
	nextMonthYear := next.Format("2006")

	lastMonth := last.Format("01")
	lastMonthYear := last.Format("2006")

	thisMonth := now.Format("01")
	thisMonthYear := now.Format("2006")

	var stringMap = map[string]string{
		"next_month":      nextMonth,
		"next_month_year": nextMonthYear,
		"last_month":      lastMonth,
		"last_month_year": lastMonthYear,
		"this_month":      thisMonth,
		"this_month_year": thisMonthYear,
	}

	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()
	firstDayMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastDayMonth := firstDayMonth.AddDate(0, 1, -1)

	DaysInMonth := lastDayMonth.Day()

	intMap := map[string]int{"days_in_month": DaysInMonth}

	rooms, err := re.DB.AllRooms()
	if err != nil {
		helpers.SeverError(rw, err)
	}

	data := map[string]interface{}{
		"data":  now,
		"rooms": rooms,
	}

	for _, rm := range rooms {
		reservationMap := make(map[string]int)
		blockMap := make(map[string]int)

		for d := firstDayMonth; !d.After(lastDayMonth); d = d.AddDate(0, 0, 1) {
			reservationMap[d.Format("2006-01-2")] = 0
			blockMap[d.Format("2006-01-2")] = 0
		}

		restrictions, err := re.DB.GetRoomRestrictionByDate(rm.ID, firstDayMonth, lastDayMonth)
		if err != nil {
			helpers.SeverError(rw, err)
			return
		}

		for _, restriction := range restrictions {
			if restriction.ReservationID > 0 {
				for d := restriction.StartDate; !d.After(restriction.EndDate); d = d.AddDate(0, 0, 1) {
					reservationMap[d.Format("2006-01-2")] = restriction.ReservationID
				}
			} else {
				blockMap[restriction.StartDate.Format("2006-01-2")] = restriction.ID
			}
		}
		data[fmt.Sprintf("reservation_map_%d", rm.ID)] = reservationMap
		data[fmt.Sprintf("block_map_%d", rm.ID)] = blockMap

		re.ConfigSetting.Session.Put(r.Context(), fmt.Sprintf("block_map_%d", rm.ID), blockMap)
	}

	render.Template(rw, r, "admin-reservations-calendar.page.html", &model.TemplateData{
		StringMap: stringMap,
		DataMap:   data,
		IntMap:    intMap,
	})

}

func (re *Repository) AdminShowReservation(rw http.ResponseWriter, r *http.Request) {
	exploded := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(exploded[4])
	if err != nil {
		helpers.SeverError(rw, err)
		return
	}

	reservation, err := re.DB.GetReservationByID(id)
	if err != nil {
		helpers.SeverError(rw, err)
		return
	}
	data := map[string]interface{}{"reservation": reservation}

	src := exploded[3]
	month := r.URL.Query().Get("m")
	year := r.URL.Query().Get("y")
	var stringMap = map[string]string{"src": src, "month": month, "year": year}

	render.Template(rw, r, "admin-reservations-show.page.html", &model.TemplateData{
		StringMap: stringMap,
		DataMap:   data,
		Form:      forms.New(nil),
	})
}

func (re *Repository) AdminPostShowReservation(rw http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		helpers.SeverError(rw, err)
		return
	}

	exploded := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(exploded[4])
	if err != nil {
		helpers.SeverError(rw, err)
		return
	}

	src := exploded[3]

	m := r.Form.Get("month")
	y := r.Form.Get("year")

	reservation, err := re.DB.GetReservationByID(id)
	if err != nil {
		helpers.SeverError(rw, err)
		return
	}

	reservation.FirstName = r.Form.Get("first_name")
	reservation.LastName = r.Form.Get("last_name")
	reservation.Email = r.Form.Get("email")
	reservation.Phone = r.Form.Get("phone")

	err = re.DB.UpdateReservations(reservation)
	if err != nil {
		helpers.SeverError(rw, err)
		return
	}

	re.ConfigSetting.Session.Put(r.Context(), "flash", "reservation successfully updated")

	if y == "" {
		http.Redirect(rw, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)
	} else {
		http.Redirect(rw, r, fmt.Sprintf("/admin/reservations-calendar?y=%s&m=%s", y, m), http.StatusSeeOther)
	}

}

func (re *Repository) AdminProcessReservation(rw http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	src := chi.URLParam(r, "src")

	year := r.URL.Query().Get("y")
	month := r.URL.Query().Get("m")

	_ = re.DB.UpdateProcessedForReservation(id, 1)
	re.ConfigSetting.Session.Put(r.Context(), "flash", "reservation is now processed")
	if year == "" {
		http.Redirect(rw, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)
	} else {
		http.Redirect(rw, r, fmt.Sprintf("/admin/reservations-calendar?y=%s&m=%s", year, month), http.StatusSeeOther)
	}
}

func (re *Repository) AdminDeleteReservation(rw http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	src := chi.URLParam(r, "src")

	year := r.URL.Query().Get("y")
	month := r.URL.Query().Get("m")

	_ = re.DB.DeleteReservationByID(id)
	re.ConfigSetting.Session.Put(r.Context(), "flash", "reservation deleted")
	if year == "" {
		http.Redirect(rw, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)
	} else {
		http.Redirect(rw, r, fmt.Sprintf("/admin/reservations-calendar?y=%s&m=%s", year, month), http.StatusSeeOther)
	}
}

func (re *Repository) AdminPostReservationsCalendar(rw http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.SeverError(rw, err)
		return
	}
	month, _ := strconv.Atoi(r.Form.Get("m"))
	year, _ := strconv.Atoi(r.Form.Get("y"))

	//process blocks

	rooms, err := re.DB.AllRooms()
	if err != nil {
		helpers.SeverError(rw, err)
		return
	}

	form := forms.New(r.PostForm)
	for _, rm := range rooms {
		currentMap := re.ConfigSetting.Session.Get(r.Context(), fmt.Sprintf("block_map_%d", rm.ID)).(map[string]int)

		for name, value := range currentMap {
			if val, ok := currentMap[name]; ok {
				if val > 0 {
					if !form.Has(fmt.Sprintf("remove_block_%d_%s", rm.ID, name)) {
						log.Println("block to be removed: ", value)
						err := re.DB.DeleteBlockByID(value)
						if err != nil {
							helpers.SeverError(rw, err)
							return
						}
					}
				}
			}
		}
	}

	// add blocks
	for name := range r.PostForm {
		//log.Println(name, value)
		if strings.HasPrefix(name, "add_block") {
			exploded := strings.Split(name, "_")
			id, _ := strconv.Atoi(exploded[2])
			t, _ := time.Parse("2006-01-2", exploded[3])
			log.Printf("would add block to room %d at %s", id, exploded[3])
			err := re.DB.InsertBlockForRoom(id, t)
			if err != nil {
				helpers.SeverError(rw, err)
				return
			}

		}
	}

	re.ConfigSetting.Session.Put(r.Context(), "flash", "successfully updated")
	http.Redirect(rw, r, fmt.Sprintf("/admin/reservations-calendar?y=%d&m=%d", year, month), http.StatusSeeOther)

}
