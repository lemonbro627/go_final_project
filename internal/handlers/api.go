package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"

	"github.com/lemonbro627/go_final_project/internal/config"
	"github.com/lemonbro627/go_final_project/internal/db/repo"
	"github.com/lemonbro627/go_final_project/internal/models"
)

type Api struct {
	repo   *repo.TasksRepository
	config *config.Config
}

// конструктор объекта API
func NewApi(repo *repo.TasksRepository, config *config.Config) *Api {
	return &Api{repo: repo, config: config}
}

func (a *Api) GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	idToSearch := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idToSearch)
	if err != nil {
		log.Println("error:", err)
		RenderApiErrorAndResponse(w, fmt.Errorf(InvalidIdError), http.StatusBadRequest)
		return
	}
	a.GetTask(w, r, id)
}

func (a *Api) PostTaskHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("We are in MethodPost")
	a.CreateTask(w, r)
}

func (a *Api) PutTaskHandler(w http.ResponseWriter, r *http.Request) {
	a.UpdateTask(w, r)
}

func (a *Api) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	idToSearch := r.URL.Query().Get("id")
	if idToSearch != "" {
		a.DeleteTask(w, r)
	} else {
		RenderApiErrorAndResponse(w, fmt.Errorf(IdMissingError), http.StatusBadRequest)
		return
	}
}

// http://localhost:7540/api/tasks/<id>
func (a *Api) GetTaskByIdHandler(w http.ResponseWriter, r *http.Request) {
	idToSearch := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idToSearch)
	if err != nil {
		log.Println("error:", err)
		RenderApiErrorAndResponse(w, fmt.Errorf(IdMissingError), http.StatusBadRequest)
		return
	}
	a.GetTask(w, r, id)
}

func (a *Api) GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("search") != "" {
		s := r.URL.Query().Get("search")
		a.SearchTasks(w, r, s)
	} else {
		a.GetAllTasks(w)
	}
}

func (a *Api) GetAllTasks(w http.ResponseWriter) {
	foundTasks, err := a.repo.GetAllTasks()
	if err != nil {
		log.Println("err:", err)
		RenderApiErrorAndResponse(w, fmt.Errorf(InternalServerError), http.StatusInternalServerError) // 500
		return
	}

	result := make(map[string][]models.Task)
	result["tasks"] = foundTasks

	resp, err := json.Marshal(result)
	if err != nil {
		log.Println("err:", err)
		RenderApiErrorAndResponse(w, fmt.Errorf(MarshallingError), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)
	if err != nil {
		log.Println("error:", err)
		return
	}
}

func (a *Api) SearchTasks(w http.ResponseWriter, r *http.Request, search string) {
	foundTasks, err := a.repo.SearchTasks(repo.QueryDataFromString(search))
	if err != nil {
		log.Println("error:", err)
		RenderApiErrorAndResponse(w, fmt.Errorf(InternalServerError), http.StatusInternalServerError) // 500
		return
	}

	result := make(map[string][]models.Task)
	result["tasks"] = foundTasks

	resp, err := json.Marshal(result)
	if err != nil {
		log.Println("error:", err)
		RenderApiErrorAndResponse(w, fmt.Errorf(MarshallingError), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(resp)
	if err != nil {
		log.Println("error:", err)
		return
	}
}

// CreateTask создаёт задание в БД
func (a *Api) CreateTask(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		log.Println("err:", err)
		RenderApiErrorAndResponse(w, fmt.Errorf(ReadingError), http.StatusBadRequest) // 400
		return
	}
	log.Println("received:", buf.String())

	parseBody := models.Task{}
	err = json.Unmarshal(buf.Bytes(), &parseBody)
	if err != nil {
		log.Println("err:", err)
		RenderApiErrorAndResponse(w, fmt.Errorf(UnMarshallingError), http.StatusBadRequest)
		return
	}

	err = parseBody.ValidateAndNormalizeDate()
	if err != nil {
		log.Println("err:", err)
		RenderApiErrorAndResponse(w, fmt.Errorf(ValidatingDateError), http.StatusBadRequest)
		return
	}

	id, err := a.repo.AddTask(parseBody)
	if err != nil {
		log.Println("err:", err)
		RenderApiErrorAndResponse(w, fmt.Errorf(InternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	WriteResponse(w, []byte(fmt.Sprintf("{\"id\":%d}", id))) //
}

// UpdateTask обновляет задание в БД
func (a *Api) UpdateTask(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		log.Println("err:", err)
		RenderApiErrorAndResponse(w, fmt.Errorf(ReadingError), http.StatusBadRequest) // 400
		return
	}

	parseBody := models.Task{}
	err = json.Unmarshal(buf.Bytes(), &parseBody)
	if err != nil {
		log.Println("err:", err)
		RenderApiErrorAndResponse(w, fmt.Errorf(UnMarshallingError), http.StatusBadRequest)
		return
	}

	err = parseBody.ValidateAndNormalizeDate()
	if err != nil {
		log.Println("error:", err)
		RenderApiErrorAndResponse(w, fmt.Errorf(ValidatingDateError), http.StatusBadRequest)
		return
	}
	idToSearch, err := strconv.Atoi(parseBody.ID)
	if err != nil {
		log.Println("error:", err)
		RenderApiErrorAndResponse(w, fmt.Errorf(InvalidIdError), http.StatusBadRequest)
		return
	}

	_, err = a.repo.GetTask(idToSearch)
	if err != nil {
		log.Println("error:", err)
		RenderApiErrorAndResponse(w, fmt.Errorf(InvalidIdError), http.StatusBadRequest)
		return
	}

	err = a.repo.UpdateTaskIn(parseBody)
	if err != nil {
		log.Println("error:", err)
		RenderApiErrorAndResponse(w, fmt.Errorf(InternalServerError), http.StatusInternalServerError)
		return
	}

	jsonItem, err := json.Marshal(parseBody)
	if err != nil {
		log.Println("error:", err)
		RenderApiErrorAndResponse(w, fmt.Errorf(MarshallingError), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	WriteResponse(w, jsonItem)
}

// DeleteTask удаяет задание из БД
func (a *Api) DeleteTask(w http.ResponseWriter, r *http.Request) {
	log.Println("We are in DeleteTask")
	idToSearch := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idToSearch)
	if err != nil {
		log.Println("error:", err)
		RenderApiErrorAndResponse(w, fmt.Errorf(InvalidIdError), http.StatusBadRequest)
		return
	}

	err = a.repo.DeleteTask(id)
	if err != nil {
		log.Println("error:", err)
		RenderApiErrorAndResponse(w, fmt.Errorf(InvalidIdError), http.StatusInternalServerError) // 500
		return
	}
	w.WriteHeader(http.StatusOK)
	WriteResponse(w, []byte("{}")) // пустой JSON
}

func (a *Api) TaskDoneHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("We are in TaskDoneHandler")
	idToSearch := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idToSearch)
	if err != nil {
		log.Println("error:", err)
		RenderApiErrorAndResponse(w, fmt.Errorf(InvalidIdError), http.StatusBadRequest)
		return
	}

	newTask, err := a.repo.PostTaskDone(id)
	if newTask == nil {
		w.WriteHeader(http.StatusOK)
		WriteResponse(w, []byte("{}")) // строка с пустым json
		return
	}
	if err != nil {
		log.Println("error:", err)
		RenderApiErrorAndResponse(w, fmt.Errorf(InternalServerError), http.StatusInternalServerError) // 500
		return
	}
	w.WriteHeader(http.StatusOK)
	WriteResponse(w, []byte("{}")) // w.Write(resp)
}

func (a *Api) GetTask(w http.ResponseWriter, r *http.Request, id int) {
	foundTask, err := a.repo.GetTask(id)
	log.Println("we are in GetTask", "foundTask:", foundTask)
	if err != nil {
		log.Println("error:", err)
		RenderApiErrorAndResponse(w, fmt.Errorf(InternalServerError), http.StatusInternalServerError) // 500
		return
	}

	resp, err := json.Marshal(foundTask)
	if err != nil {
		log.Println("error:", err)
		RenderApiErrorAndResponse(w, fmt.Errorf(MarshallingError), http.StatusInternalServerError) // 500
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)
	if err != nil {
		log.Println("error:", err)
		return
	}
}

// SigninHandler проверяет пароль и генерирует jwt token, если пароль верный
func (a *Api) SigninHandler(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		RenderApiErrorAndResponse(w, fmt.Errorf(ReadingError), http.StatusBadRequest)
		return
	}

	// берем пароль из request Body и записываем его в структуру signinRequest{} в поле Password
	reqBody := signinRequest{}
	err = json.Unmarshal(buf.Bytes(), &reqBody)
	if err != nil {
		RenderApiErrorAndResponse(w, fmt.Errorf(UnMarshallingError), http.StatusBadRequest)
		return
	}

	secret := []byte(a.config.EncryptionSecretKey)
	hashedUserPassword := HashPassword([]byte(reqBody.Password), secret)
	hashedEnvPassword := HashPassword([]byte(a.config.AppPassword), secret)

	if hashedUserPassword != hashedEnvPassword {
		RenderApiErrorAndResponse(w, fmt.Errorf("wrong password"), http.StatusUnauthorized)
		return
	}

	// получаем подписанный токен
	tokenValue, err := createToken(reqBody.Password, a.config.EncryptionSecretKey)
	if err != nil {
		RenderApiErrorAndResponse(w, fmt.Errorf(InternalServerError), http.StatusInternalServerError)
	}

	// записываем в response Body токен
	response := signinResponse{Token: tokenValue}
	respBody, err := json.Marshal(response)
	if err != nil {
		RenderApiErrorAndResponse(w, fmt.Errorf(MarshallingError), http.StatusInternalServerError)
		return
	}

	WriteResponse(w, respBody)
	w.WriteHeader(http.StatusOK)
}

// добавим проверку аутентификации для API-запросов
func (a *Api) Auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		pass := a.config.AppPassword
		if len(pass) > 0 {
			var jwtFromRequest string // JWT-токен из куки
			// получаем куки
			cookie, err := r.Cookie("token")
			if err != nil {
				RenderApiErrorAndResponse(w, fmt.Errorf("empty token"), http.StatusUnauthorized)
				return
			}
			jwtFromRequest = cookie.Value

			secret := []byte(a.config.EncryptionSecretKey)

			// валидация и проверка JWT-токена
			// парсим токен
			jwtToken, err := jwt.Parse(jwtFromRequest, func(t *jwt.Token) (interface{}, error) {
				return secret, nil
			})
			if err != nil {
				RenderApiErrorAndResponse(w, fmt.Errorf("invalid token"), http.StatusUnauthorized)
				return
			}

			// приводим поле Claims к типу jwt.MapClaims
			res, ok := jwtToken.Claims.(jwt.MapClaims)
			if !ok {
				RenderApiErrorAndResponse(w, fmt.Errorf("failed to typecast to jwt.MapCalims"), http.StatusUnauthorized)
				return
			}

			pass := res["password"]
			_, ok = pass.(string)
			if !ok {
				RenderApiErrorAndResponse(w, fmt.Errorf("failed to typecast to string"), http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	})
}