package routes


import (
"github.com/gin-gonic/gin"
"gorm.io/gorm"
"workout-tracker/controllers"
"workout-tracker/middlewares"
)


func Register(r *gin.Engine, db *gorm.DB) {
ac := &controllers.AuthController{DB: db}
ec := &controllers.ExerciseController{DB: db}
wc := &controllers.WorkoutController{DB: db}
rc := &controllers.ReportController{DB: db}


api := r.Group("/api")
{
auth := api.Group("/auth")
auth.POST("/signup", ac.Signup)
auth.POST("/login", ac.Login)


api.Use(middlewares.AuthRequired())
api.GET("/exercises", ec.List)


api.POST("/workouts", wc.Create)
api.GET("/workouts", wc.List)
api.GET("/workouts/:id", wc.Get)
api.PUT("/workouts/:id", wc.Update)
api.DELETE("/workouts/:id", wc.Delete)


api.GET("/reports", rc.Summary)
}
}