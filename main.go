package main

import "GeekCoding/router"

// @title           GeekCoding API
// @version         1.0
// @description     This is a GeekCoding Online Judge API server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /

func main() {
	r := router.Router()

	r.Run()
}
