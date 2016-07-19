/*
CREATE TABLE IF NOT EXISTS `task`(
  `task_id` int(11) NOT NULL AUTO_INCREMENT,
  `mobile_number` VARCHAR(12) NOT NULL,
  `email` VARCHAR(50) NOT NULL,
  `reminder_message` VARCHAR(500) NOT NULL,
  `reminder_time`  datetime NOT NULL,
  `reminder_status` boolean NOT NULL DEFAULT 0,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Created Time',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Updated Time',
   PRIMARY KEY (`task_id`)
 ) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;
*/

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Task struct {
	TaskId          int    `json:"task_id"`          //taskid for a specific task
	Email           string `json:"email"`            //email of person
	MobileNumber    string `json:"mobile"`           // mobile number of person
	ReminderMessage string `json:"reminder_message"` // Reminder Message of person
	ReminderTime    string `json:"reminder_time"`    // Reminder Time Of Person
	ReminderStatus  int    `json:"reminder_status"`  // Reminder Status Of Person
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type ErrorMessage struct {
	Error string `json:"error"`
}

var db *sql.DB

func main() {

	//Created Database Connection with Mysql
	var err error
	db, err = sql.Open("mysql", "root:tolexo@tcp(127.0.0.1:3306)/test")
	if err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/tasks", Tasks)
	router.HandleFunc("/task/add", TaskCreate)
	fmt.Println("Service Started")
	fmt.Println("GET /tasks")
	fmt.Println("POST /task/add")
	SendNotification()
	log.Fatal(http.ListenAndServe(":8080", router))

}

//API To List All the task details
func Tasks(w http.ResponseWriter, r *http.Request) {
	var Tasks []Task
	sql := `SELECT
			  task_id,
			  mobile_number,
			  email,
			  reminder_message,
			  IFNULL(reminder_time, ''),
			  reminder_status
			FROM task`
	rows, err := db.Query(sql)
	if err != nil {
		e := ErrorMessage{err.Error()}
		fmt.Fprintln(w, json.NewEncoder(w).Encode(e))
		return
	}
	var task_id, reminder_status int
	var mobile_number, email, reminder_message, reminder_time string
	//var created_at, updated_at time.Time
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&task_id, &mobile_number, &email, &reminder_message, &reminder_time, &reminder_status)
		if err != nil {
			e := ErrorMessage{err.Error()}
			fmt.Fprintln(w, json.NewEncoder(w).Encode(e))
			return
		}
		t := Task{task_id, mobile_number, email, reminder_message, reminder_time, reminder_status, time.Now(), time.Now()}
		Tasks = append(Tasks, t)
	}
	fmt.Fprintln(w, json.NewEncoder(w).Encode(Tasks))
}

//Api to insert the Task details in Database.
func TaskCreate(w http.ResponseWriter, r *http.Request) {
	inputFormat := "2006-01-02 15:04:05"
	var t Task
	var e ErrorMessage
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err := r.Body.Close(); err != nil {
		e = ErrorMessage{err.Error()}
		fmt.Fprintln(w, json.NewEncoder(w).Encode(e))
		return
	}
	err = json.Unmarshal(body, &t)
	if err != nil {
		e = ErrorMessage{err.Error()}
		fmt.Fprintln(w, json.NewEncoder(w).Encode(e))
		return
	}
	tm, timeErr := time.Parse(inputFormat, t.ReminderTime)
	if timeErr != nil {
		e = ErrorMessage{timeErr.Error()}
		w.WriteHeader(500)
		fmt.Fprintln(w, json.NewEncoder(w).Encode(e))
		return

	}
	if compareTotime(tm, time.Now()) == -1 {
		e = ErrorMessage{"Reminder time is older then current time"}
		w.WriteHeader(500)
		fmt.Fprintln(w, json.NewEncoder(w).Encode(e))
		return
	}
	stmt, er := db.Prepare("INSERT INTO task(mobile_number,email,reminder_message,reminder_time) VALUES(?,?,?,?)")
	if er != nil {
		e = ErrorMessage{er.Error()}
		w.WriteHeader(500)
		fmt.Fprintln(w, json.NewEncoder(w).Encode(e))
		return
	}
	_, errr := stmt.Exec(t.MobileNumber, t.Email, t.ReminderMessage, t.ReminderTime)
	if errr != nil {
		e = ErrorMessage{errr.Error()}
		w.WriteHeader(500)
		fmt.Fprintln(w, json.NewEncoder(w).Encode(e))
		return
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, json.NewEncoder(w).Encode(e))
	return
}

/*
Compare time function
@input time.Time,time.Time
@return int
*/

func compareTotime(date1, date2 time.Time) int {
	if date1.After(date2) == true {
		return 1
	}
	return -1
}

/*

 */
func SendNotification() {
	ticker := time.NewTicker(60 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				var task_id int
				var mobile_number, email, reminder_message string
				sql := `SELECT task_id,
						       mobile_number,
						       email,
						       reminder_message
						FROM task
						WHERE reminder_status = 0
						  AND reminder_time <= DATE_ADD(?, INTERVAL 1 MINUTE)
						  AND reminder_time >= ? `
				rows, _ := db.Query(sql, time.Now(), time.Now())
				defer rows.Close()
				for rows.Next() {
					rows.Scan(&task_id, &mobile_number, &email, &reminder_message)
					db.Exec("update task set reminder_status = 1 where task_id = ?", task_id)
					fmt.Println("Reminder Send to ", email)
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}
