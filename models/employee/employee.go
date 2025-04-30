package employee

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Employee struct {
	id       int64  `json:"id"`
	email    string `json:"email"`
	password string `json:"password"`
	ipAddr   string `json:"ip"`
}

func CreateEmployee(r *gin.Engine) {
	r.GET("/employee", func(c *gin.Context) {
		employee := Employee{
			id:       1,
			email:    "flakfl@mail.ru",
			password: "adsasd",
			ipAddr:   "127.0.0.1",
		}
		c.JSON(http.StatusOK, gin.H{
			"employee": gin.H{
				"id":       employee.id,
				"email":    employee.email,
				"password": employee.password,
				"ipAddr":   employee.ipAddr,
			},
		})
		//if err := c.ShouldBindJSON(&employee); err != nil {
		//	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		//}
		//c.JSON(http.StatusOK, gin.H{"employee": employee})
	})
}
