package main

import (
	"Architorture-Backend/DataLayer"
	"Architorture-Backend/GameLogic"
	"encoding/json"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
)

func main() {
	db := DataLayer.Connect()
	defer db.Close()

	hub := GameLogic.InitHub(db)
	go hub.Run()

	router := gin.New()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 12 * time.Hour,
	}))

	router.GET("/ws", func(context *gin.Context) {
		params := context.Request.URL.Query()
		roomId := params.Get("roomId")
		userName := params.Get("userName")
		expansion := params.Get("expansion")
		expansionInt, err := strconv.Atoi(expansion)
		if err != nil {
			return
		}

		log.Println("User Connecting\n  Username:", userName, "\n  RoomId:", roomId, "\n  Expansion:", expansionInt)
		GameLogic.ServeWebSocket(context.Writer, context.Request, roomId, userName, expansionInt, hub)
	})

	router.GET("CheckRoomId/:roomId", func(context *gin.Context) {
		roomId := context.Param("roomId")
		log.Println("Checking Room Id:", roomId)

		result := hub.CheckRoomId(roomId)
		resultBytes, err := json.Marshal(result)
		if err != nil {
			log.Println("Error:", err)
			return
		}

		context.Data(http.StatusOK, gin.MIMEJSON, resultBytes)
	})

	router.GET("cards/:expansion", func(context *gin.Context) {
		log.Println("Getting Available Cards")
		expansion := context.Param("expansion")
		expansionInt, err := strconv.Atoi(expansion)
		result := hub.GetCardsByExpansion(expansionInt)
		resultBytes, err := json.Marshal(result)
		if err != nil {
			log.Println("Error:", err)
			return
		}

		context.Data(http.StatusOK, gin.MIMEJSON, resultBytes)
	})

	// ROUTES FOR DEBUGGING
	// router.GET("Rooms", func(context *gin.Context) {
	// 	log.Println("Checking All Rooms")

	// 	result := hub.GetAllRooms()
	// 	resultBytes, err := json.Marshal(result)
	// 	if err != nil {
	// 		log.Println("Error:", err)
	// 		return
	// 	}

	// 	context.Data(http.StatusOK, gin.MIMEJSON, resultBytes)
	// })

	// router.GET("Rooms/:roomId/Card", func(context *gin.Context) {
	// 	log.Println("Getting Available Cards")
	// 	roomId := context.Param("roomId")
	// 	result := hub.GetCards(roomId)
	// 	resultBytes, err := json.Marshal(result)
	// 	if err != nil {
	// 		log.Println("Error:", err)
	// 		return
	// 	}

	// 	context.Data(http.StatusOK, gin.MIMEJSON, resultBytes)
	// })

	// router.GET("GameState/:id", func(context *gin.Context) {
	// 	roomId := context.Param("id")
	// 	hub.PrintGame(roomId)
	// 	context.Data(http.StatusOK, gin.MIMEJSON, []byte{})
	// })

	// router.GET("room/:roomId", func(context *gin.Context) {
	// 	roomId := context.Param("roomId")
	// 	result := hub.GetDrawPile(roomId)
	// 	println(len(result))
	// 	resultBytes, err := json.Marshal(result)
	// 	if err != nil {
	// 		log.Println("Error:", err)
	// 		return
	// 	}

	// 	context.Data(http.StatusOK, gin.MIMEJSON, resultBytes)
	// })

	router.Run("0.0.0.0:8080")
}
