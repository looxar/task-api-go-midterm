package main

import (
	"fmt"
	"log"
	"os"

	// "task-api/internal/auth"

	"task-api/internal/auth"
	"task-api/internal/item"
	"task-api/internal/user"

	"task-api/internal/mylog"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	fmt.Println("FOO : ", os.Getenv("FOO"))
	fmt.Println("TEST : ", os.Getenv("TEST"))

	postgres_user := os.Getenv("postgres_user")
	pass := os.Getenv("postgres_password")
	host := os.Getenv("postgres_host") // Assuming you also have a host
	port := os.Getenv("postgres_port")
	dbname := os.Getenv("postgres_dbname") // Assuming you have a database name

	// Construct the connection string
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", postgres_user, pass, host, port, dbname)
	fmt.Println(dsn)
	// Connect database
	// db, err := gorm.Open(
	// 	postgres.Open(
	// 		"postgres://postgres:password@localhost:5432/task",
	// 	),
	// )
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Panic(err)
	}

	controller := item.NewController(db)

	r := gin.Default()

	// Setup CORS middleware before starting the server
	config := cors.DefaultConfig()
	// frontend URL
	config.AllowOrigins = []string{"http://localhost:8008"}
	r.Use(cors.New(config))

	// Define your routes here
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to the API!",
		})
	})

	r.Use(mylog.Logger())
	// r.Use(mylog.Logger2())
	r.GET("/test", func(ctx *gin.Context) {
		fmt.Println("-----TEST------")

		value, exists := ctx.Get("example")
		if exists {
			log.Println("example = ", value)
		} else {
			log.Println("example does not exists")
		}

		value2, hasValue := ctx.Get("example2")
		if hasValue {
			log.Printf("example2 = %v, %T \n", value2, value2)
		} else {
			log.Println("example does not exists")
		}

		log.Println()
		// for i := 0; i < 10; i++ {
		// 	fmt.Println(i)
		// 	time.Sleep(1 * time.Second)
		// }

		// endless.DefaultHammerTime = 10 * time.Second
		ctx.JSON(200, "test response")
	})

	userController := user.NewController(db, os.Getenv("JWT_SECRET"))
	r.POST("/login", userController.Login)

	items := r.Group("/items")
	// items.Use(mylog.Logger2())

	// items.Use(auth.BasicAuth([]auth.Credential{
	// 	{"admin", "secret"},
	// 	{"admin2", "1234"},
	// }))
	items.Use(auth.Guard(os.Getenv("JWT_SECRET")))
	{
		items.POST("", controller.CreateItem)
		items.GET("", controller.FindItems)
		items.PATCH("/:id", controller.UpdateItemStatus)
		items.GET("/:id", controller.GetItemByID)
		items.PUT("/:id", controller.UpdateItem)    // Full update of an item by ID
		items.DELETE("/:id", controller.DeleteItem) // Delete an item by ID
	}
	// Start the server on port 2024
	r.Run(":2024")
	// items := r.Group("/items")
	// {
	// 	items.POST("", controller.CreateItem)
	// 	items.GET("", controller.FindItems)
	// 	items.PATCH("/:id", controller.UpdateItemStatus)
	// }
	// r.POST("/items", controller.CreateItem)
	// r.GET("/items", controller.FindItems)
	// r.PATCH("/items/:id", controller.UpdateItemStatus)

	if err := r.Run(); err != nil {
		log.Panic(err)
	}

	// srv := &http.Server{
	// 	Addr:    ":8080",
	// 	Handler: r.Handler(),
	// }
	// go func() {
	// 	// service connections
	// 	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
	// 		log.Fatalf("listen: %s\n", err)
	// 	}
	// }()
	// // Wait for interrupt signal to gracefully shutdown the server with
	// // a timeout of 5 seconds.
	// quit := make(chan os.Signal, 1)
	// // kill (no param) default send syscall.SIGTERM
	// // kill -2 is syscall.SIGINT
	// // kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	// signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// <-quit
	// log.Println("Shutdown Server ...")
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()
	// if err := srv.Shutdown(ctx); err != nil {
	// 	log.Fatal("Server Shutdown:", err)
	// }
	// // catching ctx.Done(). timeout of 5 seconds.
	// select {
	// case <-ctx.Done():
	// 	log.Println("timeout of 5 seconds.")
	// }
	// log.Println("Server exiting")

}

//เพิ่มรายการสั่งซื้อ 				POST 	/items
//เรียกดูรายการสั่งซื้อ 			GET 	/items?status=XXXXXXX
//เรียกดูรายละเอียดทีละรายการ 		GET 	/items/:id
//อัพเดทข้อมูลรายการสั้งซื้อ 		PUT 	/items/:id
//อนุมัติรายการ approve/reject	PATCH	/items/:id
//ลบรายการที่ยัง pending อยู่ไได้	DELETE /items/:id
