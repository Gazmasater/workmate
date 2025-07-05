// @title           Tasks API
// @version         1.0
// @description     Сервис управления задачами
// @host            localhost:8080
// @BasePath        /
package main

import (
	//_ "net/http/pprof"

	_ "github.com/gaz358/myprog/workmate/cmd/server/docs"
	"github.com/gaz358/myprog/workmate/internal/app"
)

func main() {
	app.Run()
}
