package main

import (
	"log"

	_ "github.com/skndash96/lastnight-backend/docs"
	api "github.com/skndash96/lastnight-backend/internal"
)

// @title Lastnight API
// @version 1.0
// @description This is the backend server for Lastnight
// @termsOfService http://swagger.io/terms/
// @contact.name skndash96
// @contact.email dashskndash@gmail.com
//
// @securityDefenitions.cookie ApiCookieAuth
// @in cookie
// @name lastnight_token
//
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
func main() {
	if err := api.Server(); err != nil {
		log.Fatal(err)
	}
}
