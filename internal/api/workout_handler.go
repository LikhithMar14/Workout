package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/LikhithMar14/workout-tracker/internal/utils"
)


type WorkoutHandler struct{
	Logger *log.Logger
}

func NewWorkoutHandler(logger *log.Logger)*WorkoutHandler{
	return &WorkoutHandler{
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
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"workout": fmt.Sprintf("Fetched Successfully %v",workoutID)})
	
}

func (wh *WorkoutHandler) HandleCreateWorkout(w http.ResponseWriter, r *http.Request){
	wh.Logger.Printf("HIT HandleCreateWorkout")
}

func (wh *WorkoutHandler) HandleUpdateWorkoutByID(w http.ResponseWriter, r *http.Request){
	wh.Logger.Printf("HIT HandleUpdateWorkoutByID")
}

func (wh *WorkoutHandler) HandleDeleteWorkoutByID(w http.ResponseWriter, r *http.Request) {	
	wh.Logger.Printf("HIT HandleDeleteWorkoutByID")
}