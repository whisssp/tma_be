package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"net/http"
	_ "onboarding_test/docs"
	h "onboarding_test/internal/delivery/http"
	"onboarding_test/internal/entity"
	"onboarding_test/internal/repository"
	"onboarding_test/internal/service"
	"onboarding_test/internal/usecase"
	"onboarding_test/pkg/db"
	"onboarding_test/pkg/redis"
	"onboarding_test/pkg/supabase"
)

const PORT = 8080

// @title 	Tasks Service API
// @version	1.0
// @description A Tasks service API in Go using Gin framework

// @host 	localhost:8080
// @BasePath
func main() {
	// connect local database
	db, err := db.NewDBConnection()
	if err != nil {
		log.Fatal(err)
	}
	if err := db.AutoMigrate(&entity.Task{}); err != nil {
		fmt.Println("Failed to migrate entity")
		log.Fatal(err)
	}

	// connect to redis
	redisClient := redis.NewRedisClient(nil)
	service.NewRedisService(redisClient)

	// connect to supabase
	supaStorageClient := supabase.NewSupaStorageClient()
	service.NewSupaStorageService(supaStorageClient)

	// config router
	router := gin.Default()
	router.UseRawPath = true
	router.UnescapePathValues = false
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// initialize task - repository
	taskRepoGorm := repository.NewTaskRepoGorm(db)
	taskRepoRedis := repository.NewTaskRepoRedis()

	// initialize image - repository
	imgRepoSupa := repository.NewImgRepoSupabase()

	// initialize usecase
	taskUsecase := usecase.NewTaskUsecase(taskRepoGorm, taskRepoRedis)
	fileUsecase := usecase.NewFileUsecase(imgRepoSupa)

	// register task handler
	taskHandler := h.NewTaskHandler(taskUsecase)
	taskHandler.RegisterRoutes(&router.RouterGroup)

	// register file handler
	fileHandler := h.NewFileHandler(fileUsecase)
	fileHandler.RegisterRoutes(&router.RouterGroup)

	// cors config
	router.Use(cors.Default())

	c := cron.New()
	err = c.AddFunc("*/5 * * * *", func() {
		LoadTasks(taskRepoRedis, taskRepoGorm)
	})

	if err != nil {
		fmt.Println("Error adding cron job:", err)
		return
	}

	go func() {
		c.Run()
	}()

	fmt.Printf("/nListening on port: %v", PORT)
	err = http.ListenAndServe(fmt.Sprintf(":%v", PORT), router)
	if err != nil {
		log.Fatal(err)
		return
	}
}

func LoadTasks(repoRedis *repository.TaskRepoRedis, repoGorm *repository.TaskRepoGorm) {
	tasks, err := repoGorm.GetAllTasks()
	if err != nil {
		fmt.Printf(fmt.Sprintf("error getting task from db"))
	}
	if len(tasks) != 0 {
		repoRedis.LoadTasksToRedis(tasks)
	}
}