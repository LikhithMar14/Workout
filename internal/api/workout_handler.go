package api

import (
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

}

func (wh *WorkoutHandler) HandleDeleteWorkoutByID(w http.ResponseWriter, r *http.Request) {	
	wh.Logger.Printf("HIT HandleDeleteWorkoutByID")
}