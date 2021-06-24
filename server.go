package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

func registerUser(w http.ResponseWriter, r *http.Request) {
	var newUser User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&newUser); err != nil {
		params := mux.Vars(r)
		errorString :=
			fmt.Sprintf("Invalid request payload, error: %s, json: %v", err.Error(), params)
		respondWithError(w, http.StatusBadRequest, errorString)
		return
	}

	hashedPassword, err := Hash(newUser.Password)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	newUser.Password = string(hashedPassword)

	if clonedDb := db.Create(&newUser); clonedDb.Error != nil {
		respondWithError(w, http.StatusBadRequest, clonedDb.Error.Error())
		return
	}
	checkPendingTaskAssignsFor(newUser)
	respondWithMessage(w, http.StatusOK, "User registered")
}

func loginUser(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, http.StatusUnprocessableEntity, "Could not read json body")
		return
	}
	passedUser := User{}
	json.Unmarshal(body, &passedUser)

	fetchedUser := User{}

	if clonedDb := db.First(&fetchedUser, "email = ?", passedUser.Email); clonedDb.Error != nil {
		respondWithError(w, http.StatusBadRequest, "O pa g, ay user db which nai higa!" /*clonedDb.Error.Error()*/)
		return
	}
	isValid := fetchedUser.VerifyPassword(passedUser.Password)

	if !isValid {
		respondWithError(w, http.StatusUnprocessableEntity, "Password to sahi daas dey yar")
		return
	} else {
		token, err := createToken(fetchedUser.ID)
		if err != nil {
			respondWithError(w, http.StatusUnprocessableEntity, "Error creating token")
			return
		}
		respondWithMessageAndToken(w, http.StatusOK, "Login successful", token)
		return
	}
}

func addTask(w http.ResponseWriter, r *http.Request) {

	var newTask Task
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&newTask); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if clonedDb := db.Create(&newTask); clonedDb.Error != nil {
		respondWithError(w, http.StatusBadRequest, clonedDb.Error.Error())
		return
	}
	respondWithMessage(w, http.StatusOK, "Task added successfully")
}

func editTask(w http.ResponseWriter, r *http.Request) {
	userId, err := extractTokenID(r)
	if err != nil {
		respondWithError(w, http.StatusUnprocessableEntity, "Unauthorized")
		return
	}

	params := mux.Vars(r)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, http.StatusUnprocessableEntity, "Could not read json body")
		return
	}
	passedTask := Task{}
	err = json.Unmarshal(body, &passedTask)
	if err != nil {
		respondWithError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	fetchedTask := Task{}

	if clonedDb := db.First(&fetchedTask, params["id"]); clonedDb.Error != nil {
		respondWithError(w, http.StatusBadRequest, clonedDb.Error.Error())
		return
	}

	if userId != fetchedTask.UserId {
		respondWithError(w, http.StatusUnprocessableEntity, "Tsk, Tsk, Tsk, why u try to edit task not assigned to you huh ? NO! You can only edit tasks assigned to you!")
		return
	}

	err = json.Unmarshal(body, &fetchedTask)
	if err != nil {
		respondWithError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	if clonedDb := db.Save(&fetchedTask); clonedDb.Error != nil {
		respondWithError(w, http.StatusBadRequest, clonedDb.Error.Error())
		return
	}
	respondWithMessage(w, http.StatusOK, "Task edited successfully")
}

func deleteTask(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	var task Task
	if clonedDb := db.Delete(&task, params["id"]); clonedDb.Error != nil {
		respondWithError(w, http.StatusBadRequest, clonedDb.Error.Error())
		return
	}

	//------ Adding below code coz cascade delete is not yet working for me
	deletePendingAssignTaskWithId(params["id"])
	//------

	respondWithMessage(w, http.StatusOK, "Task deleted successfully")
}

func getAllUserTasks(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var tasks []Task
	if clonedDb := db.Where("user_id = ?", params["userId"]).Find(&tasks); clonedDb.Error != nil {
		respondWithError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, tasks)
}

func assignTask(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, http.StatusUnprocessableEntity, "Could not read json body")
		return
	}
	passedAssignedTask := TaskAssign{}

	err = json.Unmarshal(body, &passedAssignedTask)
	if err != nil {
		respondWithError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	fetchedTask := Task{}

	if clonedDb := db.First(&fetchedTask, passedAssignedTask.TaskId); clonedDb.Error != nil { // check if task with the passed id exists in the tasks table
		respondWithError(w, http.StatusBadRequest, "Task with passed Id does not exist")
		return
	}

	fetchedUser := User{}
	assigneEmail := passedAssignedTask.AssigneEmail

	if clonedDb := db.First(&fetchedUser, "email = ?", assigneEmail); clonedDb.Error != nil { // check task assign table for the email passed with this request

		if clonedDb := db.Create(&passedAssignedTask); clonedDb.Error != nil { // create a pending assign task in task assign table
			respondWithError(w, http.StatusBadRequest, "Passed task already assigned to future signup against the passed email")
			return
		}
		// sendEmail("Yu hu!", "Wahahaha, a task was assigned to you, sucker!", passedAssignedTask.AssigneEmail)
		respondWithMessage(w, http.StatusOK, "Task assigned to user for future signup")
		return
	}

	fetchedTask.UserId = fetchedUser.ID

	if clonedDb := db.Save(&fetchedTask); clonedDb.Error != nil { // Update a task with the user id to which the task was assigned
		respondWithError(w, http.StatusBadRequest, clonedDb.Error.Error())
		return
	}
	respondWithMessage(w, http.StatusOK, "Task assigned to user")
}

func deletePendingAssignTaskWithId(id string) {
	var pendingAssignTask TaskAssign
	db.Where("task_id = ?", id).Delete(&pendingAssignTask)
}

func checkPendingTaskAssignsFor(user User) {
	pendingTask := TaskAssign{}
	if clonedDb := db.Where("assigne_email = ?", user.Email).Find(&pendingTask); clonedDb.Error != nil {
		return
	}

	fetchedTask := Task{}

	if clonedDb := db.First(&fetchedTask, pendingTask.TaskId); clonedDb.Error != nil {
		return
	}
	fetchedTask.UserId = user.ID

	db.Save(&fetchedTask)
	u32Str := fmt.Sprint(fetchedTask.ID)
	deletePendingAssignTaskWithId(u32Str)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	updatedMessage := "Task Failed successfully: " + message
	respondWithJSON(w, code, map[string]string{"error": updatedMessage})
}

func respondWithMessage(w http.ResponseWriter, code int, message string) {

	defaultResponse := Response{}
	defaultResponse.Message = message
	respondWithJSON(w, code, defaultResponse)
}

func respondWithMessageAndToken(w http.ResponseWriter, code int, message string, token string) {

	defaultResponse := Response{}
	defaultResponse.Message = message
	defaultResponse.Token = token
	respondWithJSON(w, code, defaultResponse)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
