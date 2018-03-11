package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/mail"

	"github.com/labstack/echo"
)

type contactResponse struct {
	Response bool `json:"response"`
}

type ContactMessage struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
	Message   string `json:"message"`
}

const (
	emailAddr = "kontakt@lens-wide-shut.de"
	sender    = "contact-web@lens-wide-shut.appspotmail.com"
)

func Contact(c echo.Context) error {
	response := contactResponse{Response: true}
	ctx := appengine.NewContext(c.Request())
	firstName := c.FormValue("fname")
	lastName := c.FormValue("lname")
	email := c.FormValue("email")
	message := c.FormValue("message")

	contactMessage := ContactMessage{
		Firstname: firstName,
		Lastname:  lastName,
		Email:     email,
		Message:   message,
	}

	contactJSON, err := json.Marshal(contactMessage)

	if err != nil {
		response.Response = false
		return c.JSON(http.StatusBadRequest, response)
	}

	msg := &mail.Message{
		Sender:  sender,
		To:      []string{emailAddr},
		Subject: "Contact Website",
		Body:    fmt.Sprintf(string(contactJSON)),
	}
	if err := mail.Send(ctx, msg); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		fmt.Println(ctx, "Couldn't send email: %v", err)
	}

	return c.JSON(http.StatusOK, response)

}
