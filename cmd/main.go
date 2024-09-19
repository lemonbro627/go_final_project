package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	_ "modernc.org/sqlite"

	"github.com/lemonbro627/go_final_project/internal/config"
	"github.com/lemonbro627/go_final_project/internal/db"
	"github.com/lemonbro627/go_final_project/internal/db/repo"
	"github.com/lemonbro627/go_final_project/internal/handlers"
)

const (
	webDir = "./web/"
)

func main() {
	// создаем config, куда записываем пароль из переменной окружения и секретное слово
	config, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error")
	}

	// Получаем переменную окружения TODO_DBFILE
	//Проверяем есть ли база данных, если нет - создаем
	db.CreateDatabase(config.DbPath)

	db, err := sql.Open("sqlite", config.DbPath)
	if err != nil {
		log.Println(err)
		return
	}
	tRepository := repo.NewTasksRepository(db)

	api := handlers.NewApi(&tRepository, config)

	r := chi.NewRouter()
	r.Handle("/*", http.FileServer(http.Dir(webDir)))
	r.Get("/api/nextdate", handlers.GetNextDay)

	// api, которое ожидается в этом задании
	// 1. POST /api/task создает таск
	// 2. GET /api/task возвращает ошибку
	// 3. GET /api/tasks возвращает набор тасков без фильтрации
	// 4. GET /api/tasks?search=... возвращает набор тасков с фильрацией по параметру search
	// 5. GET /api/tasks/{id} возвращает таск по id
	r.Delete("/api/task", api.Auth(api.DeleteTaskHandler))
	r.Put("/api/task", api.Auth(api.PutTaskHandler))
	r.Post("/api/task", api.Auth(api.PostTaskHandler))
	r.Get("/api/task", api.Auth(api.GetTaskHandler))
	r.Get("/api/tasks", api.Auth(api.GetTasksHandler))            // search
	r.Get("/api/tasks/{id}", api.Auth(api.GetTaskByIdHandler))    // http://localhost:7540/api/tasks/<id>
	r.HandleFunc("/api/task/done", api.Auth(api.TaskDoneHandler)) // post И delete, здесь id - это параметр запроса

	r.Post("/api/signin", api.SigninHandler)

	r.Get("/api/status", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	log.Printf("Starting server on port: %s", config.ApiPort)
	if err := http.ListenAndServe(":"+config.ApiPort, r); err != nil {
		log.Printf("Start server error: %s", err.Error())
	}
}
