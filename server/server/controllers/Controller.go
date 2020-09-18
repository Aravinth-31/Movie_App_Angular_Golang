package controllers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	PaytmChecksum "../Paytm_Go_Checksum/paytm"
	"github.com/go-yaml/yaml"
)

var database *sql.DB

// AddTicketData is ...
var AddTicketData string

// UpdateBookingData is ...
var UpdateBookingData string

// Config is ...
type Config struct {
	PaymentParams struct {
		Port           string `yaml:"port"`
		MID            string `yaml:"MID"`
		WEBSITE        string `yaml:"WEBSITE"`
		CHANNELID      string `yaml:"CHANNEL_ID"`
		INDUSTRYTYPEID string `yaml:"INDUSTRY_TYPE_ID"`
		ORDERID        string `yaml:"ORDER_ID"`
		CUSTID         string `yaml:"CUST_ID"`
		CALLBACKURL    string `yaml:"CALLBACK_URL"`
		KEY            string `yaml:"KEY"`
		TXNURL         string `yaml:"TXNURL"`
		TXNSTATUSURL   string `yaml:"TXNSTATUSURL"`
	} `yaml:"paymentParams"`
	DataBase struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Dbname   string `yaml:"dbname"`
	} `yaml:"database"`
}

// Configer is ...
func Configer() Config {
	var config Config
	file, err := os.Open("config.yml")
	check(err)
	defer file.Close()

	content, err := ioutil.ReadAll(file)
	check(err)

	err = yaml.Unmarshal(content, &config)
	check(err)
	return config
}

// Database is ...
func Database() {
	config := Configer()
	/* DataBase Connection */
	port, err := strconv.ParseInt(config.DataBase.Port, 10, 64)
	check(err)
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config.DataBase.Host, port, config.DataBase.User, config.DataBase.Password, config.DataBase.Dbname)
	db, err := sql.Open("postgres", psqlInfo)
	check(err)
	fmt.Println("Successfully connected!")
	database = db
}

func check(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

// GetTickets is ...
func GetTickets(w http.ResponseWriter, r *http.Request) {
	type User struct {
		Email string `json:"email"`
	}
	var tickets []map[string]interface{}
	var (
		user        User
		id          int
		name        string
		info        string
		desc        string
		image       string
		location    string
		theatrename string
		showtime    string
		date        string
		nooftickets int
		price       int
		createdat   string
		updatedat   string
		userid      int
	)
	_ = json.NewDecoder(r.Body).Decode(&user)
	row := database.QueryRow(`SELECT id from "Users" WHERE email=$1`, user.Email)
	err := row.Scan(&id)
	check(err)
	rows, err := database.Query(`SELECT * FROM "Tickets" AS "Tickets" WHERE "Tickets"."UserId"=$1`, id)
	check(err)
	for rows.Next() {
		err := rows.Scan(&id, &name, &info, &desc, &image, &location, &theatrename, &showtime, &date, &nooftickets, &price, &createdat, &updatedat, &userid)
		check(err)
		ticket := map[string]interface{}{
			"UserId":      userid,
			"createdAt":   createdat,
			"date":        date,
			"id":          id,
			"location":    location,
			"movieDesc":   desc,
			"movieImg":    image,
			"movieInfo":   info,
			"movieName":   name,
			"noOfTickets": nooftickets,
			"price":       price,
			"showTime":    showtime,
			"theatreName": theatrename,
			"updatedAt":   updatedat,
		}
		tickets = append(tickets, ticket)
	}
	json.NewEncoder(w).Encode(tickets)
}

// AddTicket is ...
func AddTicket() {
	type User struct {
		ID          int
		NoOfTickets int
		Price       int
		Email       string
	}
	var user User
	err := json.Unmarshal([]byte(AddTicketData), &user)
	check(err)
	fmt.Println(user)
	var (
		id              int
		movieName       string
		movieInfo       string
		movieDesc       string
		movieImage      string
		theatreID       int
		date            string
		showtime        string
		TheatreName     string
		TheatreLocation string
	)
	row := database.QueryRow(`SELECT id from "Users" WHERE email=$1`, user.Email)
	err = row.Scan(&id)
	check(err)
	row = database.QueryRow(`SELECT "Movies"."name","Movies"."info","Movies"."desc","Movies"."image","Movies"."TheatreId","Movies"."Date","Movies"."showtime" FROM "Movies" AS "Movies" WHERE "Movies"."id"=$1`, user.ID)
	err = row.Scan(&movieName, &movieInfo, &movieDesc, &movieImage, &theatreID, &date, &showtime)
	check(err)
	row = database.QueryRow(`SELECT name,location from "Theatres" WHERE id=$1`, theatreID)
	err = row.Scan(&TheatreName, &TheatreLocation)
	check(err)
	_, err = database.Query(`INSERT INTO "Tickets" ("movieName","movieInfo","movieDesc","movieImg","location",
		"theatreName","showTime","date","noOfTickets","price","createdAt","updatedAt","UserId")
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13) `, movieName, movieInfo, movieDesc,
		movieImage, TheatreLocation, TheatreName, showtime, date, user.NoOfTickets, user.Price,
		time.Now().Format("2006-01-02 15:04:05.166+00:00"), time.Now().Format("2006-01-02 15:04:05.166+00:00"), id)
	check(err)
	fmt.Println("Ticket Added")
}

// UpdateBooking is ...
func UpdateBooking() {
	type User struct {
		ID     int
		Booked []string
	}
	var user User
	err := json.Unmarshal([]byte(UpdateBookingData), &user)
	check(err)
	fmt.Println(user)
	str := "{"
	n := len(str)
	for _, val := range user.Booked {
		if n == len(str) {
			str += val
		} else {
			str += "," + val
		}
	}
	str += "}"
	_, err = database.Query(`UPDATE "Movies" SET booked=$1 WHERE id=$2`, str, user.ID)
	check(err)
	fmt.Println("Updated")
}

// Theatres is ...
func Theatres(w http.ResponseWriter, r *http.Request) {
	type User struct {
		Location string `json:"location"`
	}
	var user User
	var theatre string
	_ = json.NewDecoder(r.Body).Decode(&user)
	var theatres []map[string]string
	rows, err := database.Query(`SELECT DISTINCT name from "Theatres" Where location=$1`, user.Location)
	check(err)
	for rows.Next() {
		err := rows.Scan(&theatre)
		check(err)
		temp := map[string]string{
			"name": theatre,
		}
		theatres = append(theatres, temp)
	}
	json.NewEncoder(w).Encode(theatres)
}

//Locations is ...
func Locations(w http.ResponseWriter, r *http.Request) {
	var location string
	var locations []map[string]string
	rows, err := database.Query(`SELECT DISTINCT location from "Theatres"`)
	check(err)
	for rows.Next() {
		err := rows.Scan(&location)
		check(err)
		temp := map[string]string{
			"location": location,
		}
		locations = append(locations, temp)
	}
	json.NewEncoder(w).Encode(locations)
}

//SignIn is ...
func SignIn(w http.ResponseWriter, r *http.Request) {
	// User is ...
	type User struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var user User
	_ = json.NewDecoder(r.Body).Decode(&user)
	var count int
	rows, err := database.Query(`SELECT COUNT(*) as count FROM "Users" AS "User" WHERE "User"."email" = $1 AND "User"."password" = $2`, user.Email, user.Password)
	defer rows.Close()
	check(err)
	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			check(err)
		}
	}
	if count > 0 {
		json.NewEncoder(w).Encode("success")
	} else {
		json.NewEncoder(w).Encode("failure")
	}
	// var data = struct {
	// 	Result string `json:"result"`
	// }{
	// 	Result: "success",
	// }
	// jsonBytes, err := StructToJSON(data)
	// if err != nil {
	// }
	// w.Header().Set("Content-Type", "application/json")
	// w.Write(jsonBytes)
	return
}

//SignUp is ...
func SignUp(w http.ResponseWriter, r *http.Request) {
	type User struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var user User
	_ = json.NewDecoder(r.Body).Decode(&user)
	var count int
	rows, err := database.Query(`SELECT COUNT(*) as count FROM "Users" AS "User" WHERE "User"."email" = $1 `, user.Email)
	defer rows.Close()
	check(err)
	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			check(err)
		}
	}
	if count > 0 {
		json.NewEncoder(w).Encode("failure")
	} else {
		// INSERT INTO "Users" ("id","name","password","email","createdAt","updatedAt") VALUES (DEFAULT,$1,$2,$3,$4,$5) RETURNING "id","name","password","email","createdAt","updatedAt"
		sqlStatement := `
		INSERT INTO "Users" ("name","password","email","createdAt","updatedAt")
		VALUES ($1, $2, $3,$4,$5)`
		_, err = database.Exec(sqlStatement, user.Name, user.Password, user.Email, time.Now().Format("2006-01-02 15:04:05.166+00:00"), time.Now().Format("2006-01-02 15:04:05.166+00:00"))
		if err != nil {
			check(err)
		}
		json.NewEncoder(w).Encode("success")
	}
}

// Movies is ...
func Movies(w http.ResponseWriter, r *http.Request) {
	type User struct {
		Name     string `json:"name"`
		Location string `json:"location"`
		Date     string `json:"date"`
		Showtime string `json:"showtime"`
	}
	var user User
	_ = json.NewDecoder(r.Body).Decode(&user)
	type Theatre struct {
		ID        int
		Name      string
		Row       int
		Col       int
		Booked    []string
		Location  string
		CreatedAt string
		UpdatedAt string
		Movies    struct {
			Mname    string
			Mid      int
			Minfo    string
			Mdesc    string
			Mimg     string
			Mtheid   int
			Mdate    string
			Mtime    string
			Mrow     int
			Mcol     int
			Mbooked  []string
			Mcreated string
			Mupdated string
		}
	}
	var theatre Theatre
	var movies []map[string]interface{}
	rows, err := database.Query(`SELECT "Theatre"."id", "Theatre"."name", "Theatre"."row", "Theatre"."col","Theatre"."booked", "Theatre"."location", "Theatre"."createdAt", "Theatre"."updatedAt", "Movies"."id" AS "Movies.id", "Movies"."name" AS "Movies.name", "Movies"."info" AS "Movies.info","Movies"."desc" AS "Movies.desc", "Movies"."image" AS "Movies.image", "Movies"."TheatreId" AS "Movies.TheatreId", "Movies"."Date" AS "Movies.Date", "Movies"."showtime" AS "Movies.showtime","Movies"."row" AS "Movies.row", "Movies"."col" AS "Movies.col", "Movies"."booked" AS "Movies.booked","Movies"."createdAt" AS "Movies.createdAt", "Movies"."updatedAt" AS "Movies.updatedAt" FROM "Theatres" AS "Theatre" INNER JOIN "Movies" AS "Movies" ON "Theatre"."id"= "Movies"."TheatreId"`)
	defer rows.Close()
	check(err)
	var (
		temp1 []uint8
		temp2 []uint8
	)
	for rows.Next() {
		err := rows.Scan(&theatre.ID, &theatre.Name, &theatre.Row, &theatre.Col, &temp1, &theatre.Location, &theatre.CreatedAt, &theatre.UpdatedAt, &theatre.Movies.Mid, &theatre.Movies.Mname, &theatre.Movies.Minfo, &theatre.Movies.Mdesc, &theatre.Movies.Mimg, &theatre.Movies.Mtheid, &theatre.Movies.Mdate, &theatre.Movies.Mtime, &theatre.Movies.Mrow, &theatre.Movies.Mcol, &temp2, &theatre.Movies.Mcreated, &theatre.Movies.Mupdated)
		if err != nil {
			check(err)
		}
		theatre.Booked = strings.Split(string(temp1)[1:len(string(temp1))-1], ",")
		theatre.Movies.Mbooked = strings.Split(string(temp2)[1:len(string(temp2))-1], ",")
		if len(user.Name) > 0 && user.Name != theatre.Name {
			continue
		}
		if len(user.Location) > 0 && user.Location != theatre.Location {
			continue
		}
		if len(user.Showtime) > 0 && user.Showtime != theatre.Movies.Mtime {
			continue
		}
		if len(user.Date) > 0 && user.Date != theatre.Movies.Mdate {
			continue
		}
		f := 1
		for _, value := range movies {
			if value["name"] == theatre.Movies.Mname {
				f = 0
			}
		}
		if f == 1 {
			movie := map[string]interface{}{
				"Date":      theatre.Movies.Mdate,
				"TheatreId": theatre.ID,
				"booked":    theatre.Movies.Mbooked,
				"col":       theatre.Movies.Mcol,
				"desc":      theatre.Movies.Mdesc,
				"id":        theatre.Movies.Mid,
				"image":     theatre.Movies.Mimg,
				"info":      theatre.Movies.Minfo,
				"name":      theatre.Movies.Mname,
				"row":       theatre.Movies.Mrow,
				"showtime":  theatre.Movies.Mname,
			}
			movies = append(movies, movie)
		}
	}
	json.NewEncoder(w).Encode(movies)
}

// Result is ...
type Result struct {
	TXNID       string
	BANKTXNID   string
	ORDERID     string
	TXNAMOUNT   string
	STATUS      string
	TXNTYPE     string
	GATEWAYNAME string
	RESPCODE    string
	RESPMSG     string
	BANKNAME    string
	MID         string
	PAYMENTMODE string
	REFUNDAMT   string
	TXNDATE     string
}

//IndexHandler is
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	config := Configer()
	Amount := r.FormValue("amt")
	PhnNumber := r.FormValue("number")
	Email := r.FormValue("email")
	AddTicketData = r.FormValue("addTicketData")
	UpdateBookingData = r.FormValue("updatebookingData")
	paytmParams := make(map[string]string)
	paytmParams = map[string]string{
		"MID":              config.PaymentParams.MID,
		"WEBSITE":          config.PaymentParams.WEBSITE,
		"CHANNEL_ID":       config.PaymentParams.CHANNELID,
		"INDUSTRY_TYPE_ID": config.PaymentParams.INDUSTRYTYPEID,
		"ORDER_ID":         config.PaymentParams.ORDERID + time.Now().Format("20060102150405"),
		"CUST_ID":          config.PaymentParams.CUSTID,
		"TXN_AMOUNT":       Amount,
		"CALLBACK_URL":     config.PaymentParams.CALLBACKURL,
		"EMAIL":            Email,
		"MOBILE_NO":        PhnNumber,
	}
	fmt.Println(paytmParams)
	paytmChecksum := PaytmChecksum.GenerateSignature(paytmParams, config.PaymentParams.KEY)
	verifyChecksum := PaytmChecksum.VerifySignature(paytmParams, config.PaymentParams.KEY, paytmChecksum)
	fmt.Println(verifyChecksum)

	formFields := ""
	for key, value := range paytmParams {
		formFields += `<input type="hidden" name="` + key + `" value="` + value + `">`
	}
	formFields += `<input type="hidden" name="CHECKSUMHASH" value="` + paytmChecksum + `">`

	fmt.Fprintf(w, `<html><head><title>Merchant Checkout Page</title></head>
					<body><center><h1>Please do not refresh this page...</h1></center>
					<form method="post" action="`+config.PaymentParams.TXNURL+`" name="f1">`+formFields+`</form>
					<script type="text/javascript">document.f1.submit();</script>
					</body></html>`)
}

// CallBackHandler is
func CallBackHandler(w http.ResponseWriter, r *http.Request) {
	config := Configer()
	r.ParseForm()
	postData := make(map[string]string)
	for key, value := range r.Form {
		postData[key] = value[0]
	}

	// Send Server-to-Server request to verify Order Status
	Body := map[string]string{
		"MID":     postData["MID"],
		"ORDERID": postData["ORDERID"],
	}
	checksum := PaytmChecksum.GenerateSignature(Body, config.PaymentParams.KEY)
	var jsonStr = []byte(`JsonData={"MID":"` + config.PaymentParams.MID + `","ORDERID":"` + postData["ORDERID"] + `","CHECKSUMHASH":"` + checksum + `"}`)
	resp, err := http.Post(config.PaymentParams.TXNSTATUSURL, "application/json", bytes.NewBuffer(jsonStr))
	check(err)
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var ans Result
	json.Unmarshal(body, &ans)
	html := "<body><style>h1,h3{margin-top:100px;}body,#butt{display:flex;justify-content:center;}"
	html += `button{padding:10px 20px;cursor:pointer}</style><div><div id="butt">`
	if ans.RESPCODE == "01" {
		html += "<h1>Payment Success<h1>"
		UpdateBooking()
		AddTicket()
	} else {
		html += "<h3>" + ans.RESPMSG + "</h3>"
	}
	html += `</div><div id="butt"><button onclick="movies()">Return To Movie Booking</button></div></div>`
	html += `<script type="text/javascript">function movies(){window.location.replace("http://localhost:4200/movies");}</script></body>`

	fmt.Fprintf(w, html)
}
