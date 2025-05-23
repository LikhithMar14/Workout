package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/LikhithMar14/workout-tracker/internal/store"
	"github.com/LikhithMar14/workout-tracker/internal/utils"
)


type WorkoutHandler struct{
	WorkoutStore store.WorkoutStore
	Logger *log.Logger
}

func NewWorkoutHandler(workoutStore store.WorkoutStore ,logger *log.Logger)*WorkoutHandler{
	return &WorkoutHandler{
		WorkoutStore: workoutStore,
		Logger: logger,
	}
}

func (wh *WorkoutHandler) HandleGetWorkoutByID(w http.ResponseWriter, r *http.Request){
	workoutID,err := utils.ReadIDParam(r)
	if err != nil{
		wh.Logger.Printf("Error: readIDParam: %v",err)
		utils.WriteJSON(w,http.StatusBadRequest,utils.Envelope{"error":"invalid workout id"})
		return
	}
	workout,err := wh.WorkoutStore.GetWorkoutByID(workoutID)
	if err != nil {
		wh.Logger.Printf("ERROR: getWorkoutByID: %v",err)
		utils.WriteJSON(w,http.StatusInternalServerError,utils.Envelope{"error":"internal server error"})
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"workout": workout})
	
}

func (wh *WorkoutHandler) HandleCreateWorkout(w http.ResponseWriter, r *http.Request){
	var workout store.Workout
	err := json.NewDecoder(r.Body).Decode(&workout)

	if err != nil{
		fmt.Println(err)
		http.Error(w,"Failed to create workout",http.StatusInternalServerError)
		return
	}

	createdWorkout,err := wh.WorkoutStore.CreateWorkout(&workout)
	if err != nil{
		wh.Logger.Printf("ERROR: getWorkoutByID: %v",err)
		utils.WriteJSON(w,http.StatusInternalServerError,utils.Envelope{"error":"internal server error"})
		return
	}

	utils.WriteJSON(w,http.StatusOK,utils.Envelope{"workout": createdWorkout})


}

func (wh *WorkoutHandler) HandleUpdateWorkoutByID(w http.ResponseWriter, r *http.Request){
	workoutID,err := utils.ReadIDParam(r)
	if err != nil{
		wh.Logger.Printf("ERROR: readIDParam: %v",err)
		utils.WriteJSON(w,http.StatusInternalServerError,utils.Envelope{"error":"invalid workout id"})
		return
	}
	existingWorkout,err := wh.WorkoutStore.GetWorkoutByID(workoutID)
	if err != nil {
		wh.Logger.Printf("ERROR: getWorkoutByID: %v",err)
		utils.WriteJSON(w,http.StatusInternalServerError,utils.Envelope{"error":"internal server error"})
		return
	}

	if existingWorkout == nil {
		http.NotFound(w,r)
		return;
	}

	var updateWorkoutRequest struct{
		Title           *string         `json:"title,omitempty"`
		Description     *string         `json:"description,omitempty"`
		DurationMinutes *int            `json:"duration_minutes"`
		CaloriesBurned  *int           `json:"calories_burned,omitempty"`
		Entries         []store.WorkoutEntry `json:"entries"`
	}

	err = json.NewDecoder(r.Body).Decode(&updateWorkoutRequest)
	if err != nil {
		wh.Logger.Printf("ERRORP: decodingUpdateRequest: %v",err)
		utils.WriteJSON(w,http.StatusBadRequest,utils.Envelope{"error":"invalid request payload"})
		return
	}

	if updateWorkoutRequest.Title != nil {
		existingWorkout.Title = *updateWorkoutRequest.Title
	}
	if updateWorkoutRequest.Description != nil {
		existingWorkout.Description = *updateWorkoutRequest.Description
	}
	if updateWorkoutRequest.DurationMinutes != nil {
		existingWorkout.DurationMinutes = *updateWorkoutRequest.DurationMinutes
	}
	if updateWorkoutRequest.CaloriesBurned != nil {
		existingWorkout.CaloriesBurned = updateWorkoutRequest.CaloriesBurned
	}

	if updateWorkoutRequest.Entries != nil {
		existingWorkout.Entries = updateWorkoutRequest.Entries
	}

	err = wh.WorkoutStore.UpdateWorkout(existingWorkout)
	if err != nil {
		wh.Logger.Printf("ERROR: updatingWorkout: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"workout": existingWorkout})

}

func (wh *WorkoutHandler) HandleDeleteWorkoutByID(w http.ResponseWriter, r *http.Request) {	
	paramsWorkoutID,err := utils.ReadIDParam(r)
	if err != nil {
		wh.Logger.Printf("ERROR: readIDParam: %v",err)
		utils.WriteJSON(w,http.StatusBadRequest,utils.Envelope{"error":"invalid workout id"})
		return
	}
	err = wh.WorkoutStore.DeleteWorkoutByID(paramsWorkoutID)
	if err == sql.ErrNoRows {
		http.Error(w, "workout not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, "error deleting workout", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}