package repo

import (
	"database/sql"
	"log"
	"strings"
	"time"

	"github.com/lemonbro627/go_final_project/internal/dateutil"
	"github.com/lemonbro627/go_final_project/internal/models"
)

const (
	limit = 20
)

// чтобы оперировать Tasks (TaskCreationRequest), нужна всегда ссылка на БД
type TasksRepository struct {
	db *sql.DB
}

func NewTasksRepository(db *sql.DB) TasksRepository {
	return TasksRepository{db: db}
}

func (tr TasksRepository) AddTask(t models.Task) (int, error) {
	task, err := tr.db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)",
		sql.Named("date", t.Date),
		sql.Named("title", t.Title),
		sql.Named("comment", t.Comment),
		sql.Named("repeat", t.Repeat))

	if err != nil {
		return 0, err
	}

	id, err := task.LastInsertId()
	if err != nil {
		return 0, err
	}

	// возвращаем ID последней добавленной записи
	return int(id), nil
}

// PostTaskDone moves task according the repeat rule
func (tr TasksRepository) PostTaskDone(id int) (*models.Task, error) {
	t, err := tr.GetTask(id)
	if err != nil {
		return nil, err
	}

	dt, err := time.Parse(dateutil.DateFormat, t.Date)
	if err != nil {
		return nil, err
	}

	if t.Repeat == "" {
		log.Println("Repeat is null")
		err = tr.DeleteTask(id)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}

	now := time.Now()
	nextDate, err := dateutil.NextDate(now, dt, t.Repeat)
	if err != nil {
		return nil, err
	}
	err = tr.UpdateTaskDate(t, nextDate)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

// UpdateTask - put Method, обновляет задание БД.
func (tr TasksRepository) UpdateTaskIn(t models.Task) error {
	_, err := tr.db.Exec("UPDATE scheduler SET date = :date, title = :title, comment = :comment,"+
		"repeat = :repeat WHERE id = :id",
		sql.Named("date", t.Date),
		sql.Named("title", t.Title),
		sql.Named("comment", t.Comment),
		sql.Named("repeat", t.Repeat),
		sql.Named("id", t.ID))

	if err != nil {
		return err
	}

	return nil
}

// Чтение строки по заданному id.
// Из таблицы должна вернуться одна строка.
func (tr TasksRepository) GetTask(id int) (models.Task, error) {
	s := models.Task{}
	row := tr.db.QueryRow("SELECT id, date, title, comment, repeat from scheduler WHERE id = :id",
		sql.Named("id", id))

	// заполняем TaskCreationRequest данными из таблицы
	err := row.Scan(&s.ID, &s.Date, &s.Title, &s.Comment, &s.Repeat)
	if err != nil {
		return models.Task{}, err
	}
	return s, nil
}

// Из таблицы должны вернуться сроки с ближайшими датами.
func (tr TasksRepository) GetAllTasks() ([]models.Task, error) {
	today := time.Now().Format(dateutil.DateFormat)

	rows, err := tr.db.Query("SELECT id, date, title, comment, repeat FROM scheduler WHERE date >= :today "+
		"ORDER BY date LIMIT :limit",
		sql.Named("today", today),
		sql.Named("limit", limit))

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	result := []models.Task{}
	// заполняем Task данными из таблицы
	for rows.Next() { // идём по записям
		s := models.Task{} // создаем новый Task и заполняем его данными из текущего row
		err := rows.Scan(&s.ID, &s.Date, &s.Title, &s.Comment, &s.Repeat)
		if err != nil {
			return nil, err
		}
		result = append(result, s)
	}
	//Проверяем успешное завершение цикла
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

// Из таблицы должна вернуться строка в соответсвии с критерием поиска search.
func (tr TasksRepository) SearchTasks(searchData SearchQueryData) ([]models.Task, error) {
	var rows *sql.Rows

	queryData := searchData.GetQueryData()

	querySql := strings.Join([]string{
		"SELECT id, date, title, comment, repeat FROM scheduler",
		queryData.Condition,
		"ORDER BY date LIMIT :limit",
	}, " ")

	rows, err := tr.db.Query(querySql,
		sql.Named("search", queryData.Param),
		sql.Named("limit", limit))

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	result := []models.Task{}
	// заполняем объект Task данными из таблицы
	for rows.Next() { // идём по записят
		s := models.Task{} // создаем новый Task и заполняем его данными из текущего row

		if err := rows.Scan(&s.ID, &s.Date, &s.Title, &s.Comment, &s.Repeat); err != nil {
			return nil, err
		}
		result = append(result, s)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

// Удаление строки по заданному id.
func (tr TasksRepository) DeleteTask(id int) error {
	_, err := tr.db.Exec("DELETE FROM scheduler WHERE id = :id",
		sql.Named("id", id))
	if err != nil {
		return err
	}

	return nil
}

// UpdateTask обновляет задание в БД соблюдая дату из правила повторний
func (tr TasksRepository) UpdateTaskDate(t models.Task, newDate string) error {
	_, err := tr.db.Exec("UPDATE scheduler SET date = :date WHERE id = :id",
		sql.Named("date", newDate),
		sql.Named("id", t.ID))

	if err != nil {
		return err
	}

	return nil
}
