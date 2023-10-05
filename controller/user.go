package controller

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
)

type User struct {
	Id        int    `json:"_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

func (u *User) GetUserProfile(c *gin.Context) {

	endpoint := c.Request.URL

	cachedKey := endpoint.String()

	resp := map[string]string{
		"first_name": u.FirstName,
		"last_name":  u.LastName,
		"email":      u.Email,
	}

	m, _ := json.Marshal(&resp)

	cacheErr := redis

}
