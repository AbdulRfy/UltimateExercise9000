package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	models "ultimate.com/exercise/Models"
	response "ultimate.com/exercise/apiresponse"
	jwtToken "ultimate.com/exercise/jwtToken"
)

func registerUser(w http.ResponseWriter, r *http.Request) {
	var newUser models.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&newUser); err != nil {
		params := mux.Vars(r)
		errorString :=
			fmt.Sprintf("Invalid request payload, error: %s, json: %v", err.Error(), params)
		response.RespondWithError(w, http.StatusBadRequest, errorString)
		return
	}

	hashedPassword, err := models.Hash(newUser.Password)
	if err != nil {
		response.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	newUser.Password = string(hashedPassword)

	if clonedDb := db.Create(&newUser); clonedDb.Error != nil {
		response.RespondWithError(w, http.StatusBadRequest, clonedDb.Error.Error())
		return
	}
	checkPendingTaskAssignsFor(newUser)
	response.RespondWithMessage(w, http.StatusOK, "User registered")
}

func loginUser(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response.RespondWithError(w, http.StatusUnprocessableEntity, "Could not read json body")
		return
	}
	passedUser := models.User{}
	json.Unmarshal(body, &passedUser)

	fetchedUser := models.User{}

	if clonedDb := db.First(&fetchedUser, "email = ?", passedUser.Email); clonedDb.Error != nil {
		response.RespondWithError(w, http.StatusBadRequest, "O pa g, ay user db which nai higa!" /*clonedDb.Error.Error()*/)
		return
	}
	isValid := fetchedUser.VerifyPassword(passedUser.Password)

	if !isValid {
		response.RespondWithError(w, http.StatusUnprocessableEntity, "Password to sahi daas dey yar")
		return
	} else {
		token, err := jwtToken.CreateToken(fetchedUser.ID)
		if err != nil {
			response.RespondWithError(w, http.StatusUnprocessableEntity, "Error creating token")
			return
		}
		response.RespondWithMessageAndToken(w, http.StatusOK, "Login successful", token)
		return
	}
}

func addTask(w http.ResponseWriter, r *http.Request) {

	userId, err := jwtToken.ExtractTokenID(r)
	if err != nil {
		response.RespondWithError(w, http.StatusUnprocessableEntity, "Unauthorized")
		return
	}

	var newTask models.Task
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&newTask); err != nil {
		response.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	newTask.OwnerId = userId
	if clonedDb := db.Create(&newTask); clonedDb.Error != nil {
		response.RespondWithError(w, http.StatusBadRequest, clonedDb.Error.Error())
		return
	}
	response.RespondWithMessage(w, http.StatusOK, "Task added successfully")
}

func editTask(w http.ResponseWriter, r *http.Request) {
	userId, err := jwtToken.ExtractTokenID(r)
	if err != nil {
		response.RespondWithError(w, http.StatusUnprocessableEntity, "Unauthorized")
		return
	}

	params := mux.Vars(r)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response.RespondWithError(w, http.StatusUnprocessableEntity, "Could not read json body")
		return
	}
	passedTask := models.Task{}
	err = json.Unmarshal(body, &passedTask)
	if err != nil {
		response.RespondWithError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	fetchedTask := models.Task{}

	if clonedDb := db.First(&fetchedTask, params["id"]); clonedDb.Error != nil {
		response.RespondWithError(w, http.StatusBadRequest, clonedDb.Error.Error())
		return
	}

	if userId != fetchedTask.UserId {
		response.RespondWithError(w, http.StatusUnprocessableEntity, "Tsk, Tsk, Tsk, why u try to edit task not assigned to you huh ? NO! You can only edit tasks assigned to you!")
		return
	}

	err = json.Unmarshal(body, &fetchedTask)
	if err != nil {
		response.RespondWithError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	if clonedDb := db.Save(&fetchedTask); clonedDb.Error != nil {
		response.RespondWithError(w, http.StatusBadRequest, clonedDb.Error.Error())
		return
	}
	response.RespondWithMessage(w, http.StatusOK, "Task edited successfully")
}

func deleteTask(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	var task models.Task
	if clonedDb := db.Delete(&task, params["id"]); clonedDb.Error != nil {
		response.RespondWithError(w, http.StatusBadRequest, clonedDb.Error.Error())
		return
	}

	//------ Adding below code coz cascade delete is not yet working for me
	deletePendingAssignTaskWithId(params["id"])
	//------

	response.RespondWithMessage(w, http.StatusOK, "Task deleted successfully")
}

func getAllUserTasks(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var tasks []models.Task
	if clonedDb := db.Where("user_id = ? OR owner_id = ?", params["userId"], params["userId"]).Find(&tasks); clonedDb.Error != nil {
		response.RespondWithError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	response.RespondWithJSON(w, http.StatusOK, tasks)
}

func assignTask(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response.RespondWithError(w, http.StatusUnprocessableEntity, "Could not read json body")
		return
	}
	passedAssignedTask := models.TaskAssign{}

	err = json.Unmarshal(body, &passedAssignedTask)
	if err != nil {
		response.RespondWithError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	fetchedTask := models.Task{}

	if clonedDb := db.First(&fetchedTask, passedAssignedTask.TaskId); clonedDb.Error != nil { // check if task with the passed id exists in the tasks table
		response.RespondWithError(w, http.StatusBadRequest, "Task with passed Id does not exist")
		return
	}

	fetchedUser := models.User{}
	assigneEmail := passedAssignedTask.AssigneEmail

	if clonedDb := db.First(&fetchedUser, "email = ?", assigneEmail); clonedDb.Error != nil { // check task assign table for the email passed with this request

		if clonedDb := db.Create(&passedAssignedTask); clonedDb.Error != nil { // create a pending assign task in task assign table
			response.RespondWithError(w, http.StatusBadRequest, "Passed task already assigned to future signup against the passed email")
			return
		}
		// sendEmail("Yu hu!", "Wahahaha, a task was assigned to you, sucker!", passedAssignedTask.AssigneEmail)
		response.RespondWithMessage(w, http.StatusOK, "Task assigned to user for future signup")
		return
	}

	fetchedTask.UserId = fetchedUser.ID

	if clonedDb := db.Save(&fetchedTask); clonedDb.Error != nil { // Update a task with the user id to which the task was assigned
		response.RespondWithError(w, http.StatusBadRequest, clonedDb.Error.Error())
		return
	}
	response.RespondWithMessage(w, http.StatusOK, "Task assigned to user")
}

func deletePendingAssignTaskWithId(id string) {
	var pendingAssignTask models.TaskAssign
	db.Where("task_id = ?", id).Delete(&pendingAssignTask)
}

func checkPendingTaskAssignsFor(user models.User) {
	var pendingTasks []models.TaskAssign
	if clonedDb := db.Where("assigne_email = ?", user.Email).Find(&pendingTasks); clonedDb.Error != nil {
		return
	}

	for _, pTask := range pendingTasks {
		fetchedTask := models.Task{}
		if clonedDb := db.First(&fetchedTask, pTask.TaskId); clonedDb.Error != nil {
			return
		}
		fetchedTask.UserId = user.ID

		db.Save(&fetchedTask)
		u32Str := fmt.Sprint(fetchedTask.ID)
		deletePendingAssignTaskWithId(u32Str)
	}
}
