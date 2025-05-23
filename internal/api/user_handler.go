package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"

	store "github.com/LikhithMar14/workout-tracker"
	"github.com/LikhithMar14/workout-tracker/internal/utils"
)

type registerUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Bio      string `json:"bio"`
}


type UserHandler struct {
	UserStore store.UserStore
	Logger 	*log.Logger
}


func NewUserHandler(userStore store.UserStore, logger *log.Logger) *UserHandler{
	return &UserHandler{
		UserStore: userStore,
		Logger: logger,
	}
}


func (uh *UserHandler) validateRegisterRequest(req *registerUserRequest) error {
	if req.Username == ""{
		return errors.New("username is required")
	}


	if len(req.Username) > 50 {
		return errors.New("username cannot be greater than 50 characters")
	}

	if req.Email == "" {
		return errors.New("email is required")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(req.Email) {
		return errors.New("invalid email format")
	}

	if req.Password == "" {
		return errors.New("password is required")
	}

	return nil

}

func (uh *UserHandler) HandleRegisterUser(w http.ResponseWriter, r *http.Request){
	var req registerUserRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		uh.Logger.Printf("ERROR: decoding register user: %v",err)
		utils.WriteJSON(w,http.StatusBadRequest,utils.Envelope{"error":"invalid request payload"})
		return
	}


}