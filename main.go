package main

import (
	"gpt/cmd/handler"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	vip := r.Group("/vip")
	{
		vip.POST("/open", handler.OpenVIP)
	}
	conversation := r.Group("/conversation")
	{
		conversation.POST("/new", handler.NewConversation)
		conversation.POST("/start", handler.StartConversation)
	}
	if err := r.Run("0.0.0.0:8099"); err != nil {
		log.Fatal("server init failed")
	}
}
