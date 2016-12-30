package controller

import (
	"gopkg.in/kataras/iris.v5"
)

type EmployeeController struct {
}

func (ctr EmployeeController) EmployeeList(c *iris.Context) {
	var msg struct {
		Name    string
		Message string
		Number  int
	}
	msg.Name = "Lena"
	msg.Message = "hey"
	msg.Number = 123
	// Note that msg.Name becomes "user" in the JSON
	// Will output  :   {"user": "Lena", "Message": "hey", "Number": 123}
	c.JSON(iris.StatusOK, msg)
}
