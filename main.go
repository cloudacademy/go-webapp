package main

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/labstack/echo/v4"
)

func getSHA(c echo.Context) error {
	data := make([]byte, 10)
	for i := range data {
		data[i] = byte(rand.Intn(256))
	}

	c.Response().Header().Add("Content-Type", "text/plain")
	return c.String(http.StatusOK, fmt.Sprintf("%x", sha256.Sum256(data)))
}

func ok(c echo.Context) error {
	c.Response().Header().Add("Content-Type", "text/plain")
	return c.String(http.StatusOK, "OK!")
}

func main() {
	e := echo.New()

	e.GET("/ok", ok)
	e.GET("/sha256", getSHA)

	e.Logger.Fatal(e.Start(":8080"))
}
