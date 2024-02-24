package main

import (
	"github.com/gin-gonic/gin"
)

const defaultFlashCookieName = "flash"

func cookieName(name []string) string {
	if len(name) > 0 {
		return name[0]
	}
	return defaultFlashCookieName
}

func Flash(c *gin.Context, msg string, name ...string) {
	c.SetCookie(cookieName(name), msg, 0, "", "", false, false)
}

func FlashMessage(c *gin.Context, name ...string) string {
	flashCookieName := cookieName(name)
	msg, err := c.Cookie(flashCookieName)
	defer c.SetCookie(flashCookieName, "", -1, "", "", false, false)
	if err != nil {
		return ""
	}

	return msg
}
