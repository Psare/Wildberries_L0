package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"syscall"
	"net/http"
	"os/signal"
	"wb_l0/service"
	"wb_l0/subscriber"

	"github.com/gin-gonic/gin"
)

var o_service *service.OrderService

func getByUID(ctx *gin.Context) {
	uid := ctx.Query("order_uid")
	order, err := o_service.GetOrderByUID(&uid)
	if err != nil {
		_, err = ctx.Writer.Write([]byte(fmt.Sprintf("Can't find order with this UID: %v", err)))
		if err != nil {
			fmt.Printf("ERROR: can't find order with this UID: %v\n", err)
		}
		return
	}
	bytes, err := json.Marshal(order)
	if err != nil {
		fmt.Printf("ERROR: can't marshal order: %v\n", err)
		return
	}
	ctx.Data(http.StatusOK, "application/json", bytes)
}

func main() {
	o_service = service.NewOrderService()

	router := gin.Default()
	router.LoadHTMLFiles("web/html/index.html")
	router.Static("/css", "web/css")
	router.GET("", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", gin.H{
			"title": "main page",
		})
	})
	router.GET("/wb", getByUID)

	go subscriber.Subscribe(o_service)

	router.Run(":8000")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	log.Println("app Shutting down")
	os.Exit(0)
}