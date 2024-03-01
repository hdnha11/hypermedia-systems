package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var contactRepo InMemContactRepository

func main() {
	r := gin.Default()
	r.Static("/static", "./static")
	r.HTMLRender = loadTemplates("./templates")
	r.GET("/", handleIndex)
	r.GET("/contacts", handleContacts)
	r.GET("/contacts/new", handleContactNewGet)
	r.POST("/contacts/new", handleContactNew)
	r.GET("/contacts/:id", handleContactView)
	r.GET("/contacts/:id/edit", handleContactEditGet)
	r.POST("/contacts/:id/edit", handleContactEdit)
	r.POST("/contacts/:id/delete", handleContactDelete)

	r.Run(":8080")
}

func renderError(c *gin.Context, err error) {
	c.HTML(http.StatusOK, "error.html", templateData(c, gin.H{"error": err}))
}

func handleIndex(c *gin.Context) {
	c.Redirect(http.StatusFound, "/contacts")
}

func handleContacts(c *gin.Context) {
	search := c.Query("q")
	var (
		contacts []*Contact
		err      error
	)
	if search != "" {
		contacts, err = contactRepo.Search(c.Request.Context(), search)
	} else {
		contacts, err = contactRepo.GetAll(c.Request.Context())
	}

	if err != nil {
		log.Println(err)
		renderError(c, err)
		return
	}

	c.HTML(http.StatusOK, "index.html", templateData(c, contacts))
}

func handleContactNewGet(c *gin.Context) {
	c.HTML(http.StatusOK, "new.html", templateData(c, gin.H{
		"contact": Contact{},
		"errors":  FieldErrors(nil),
	}))
}

func handleContactNew(c *gin.Context) {
	var contact Contact
	if err := c.Bind(&contact); err != nil {
		log.Println(err)
		renderError(c, err)
		return
	}

	if err := contactRepo.Save(c.Request.Context(), &contact); err != nil {
		log.Println(err)
		c.HTML(http.StatusOK, "new.html", templateData(c, gin.H{
			"contact": contact,
			"errors":  FieldErrors(err),
		}))
	}

	Flash(c, "Created New Contact!")
	c.Redirect(http.StatusSeeOther, "/contacts")
}

func handleContactView(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Println(err)
		renderError(c, err)
		return
	}

	contact, err := contactRepo.Find(c.Request.Context(), int(id))
	if err != nil {
		log.Println(err)
		renderError(c, err)
		return
	}

	c.HTML(http.StatusOK, "show.html", templateData(c, contact))
}

func handleContactEditGet(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Println(err)
		renderError(c, err)
		return
	}

	contact, err := contactRepo.Find(c.Request.Context(), int(id))
	if err != nil {
		log.Println(err)
		renderError(c, err)
		return
	}

	c.HTML(http.StatusOK, "edit.html", templateData(c, gin.H{
		"contact": contact,
		"errors":  map[string]string{},
	}))
}

func handleContactEdit(c *gin.Context) {
	var contactUpd Contact
	if err := c.Bind(&contactUpd); err != nil {
		log.Println(err)
		renderError(c, err)
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Println(err)
		renderError(c, err)
		return
	}

	contactUpd.ID = int(id)

	if err := contactRepo.Save(c.Request.Context(), &contactUpd); err != nil {
		log.Println(err)
		c.HTML(http.StatusOK, "edit.html", templateData(c, gin.H{
			"contact": contactUpd,
			"errors":  FieldErrors(err),
		}))
	}

	Flash(c, "Updated Contact!")
	c.Redirect(http.StatusSeeOther, fmt.Sprintf("/contacts/%d", id))
}

func handleContactDelete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Println(err)
		renderError(c, err)
		return
	}

	if err := contactRepo.Delete(c.Request.Context(), int(id)); err != nil {
		log.Println(err)
		renderError(c, err)
		return
	}

	Flash(c, "Deleted Contact!")
	c.Redirect(http.StatusSeeOther, "/contacts")
}
