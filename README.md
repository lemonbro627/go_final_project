# ToDo List

#### Простое приложение - планировщик задач

В этом приложении реализован следующий функционал:
- Добавление задач
- Получение списка задач
- Удаление задач
- Возможность планирования задач с интервалами:
  Каждый год
  Каждый месяц
  Каждую неделю
  Каждые N дней
- Поиск задач по ID
- Реализовано API для взаимодействия с задачами
- Возможность запуска в Docker контейнере
- Поиск задач по дате и времени


## Запуск проекта

Для запуска этого приложения:
```
go run ./cmd/main.go - не рекоменжованный способ, при подобном запуске база данных каждый раз будет создаваться новая во временной директории

go build -o server ./cmd/main.go - рекоммендованный способ, создаст базу рядом с собой при запуске бинарного файла
```
Возможные переменные:

- `TODO_PASSWORD` - Дефолтно 123456, если хотите изменить - нужно задать переменную окружения .
- `TODO_DBFILE` - Расположение базы данных SQLite, обязательно если мы запускаем тесты, во всех остальных случаях определяется рядом с бинарником
- `TODO_PORT` - Порт на котором работает приложение, дефолт 7540.

#### Запуск в докере 
``` bash
docker build -t go_final_project:v1.0.0 .
docker run -p 7540:7540 -e TODO_PASSWORD=1234567 golang-go_final_project:v1.0.0
```

## Запуск тестов

# Starting tests requires running application.
Необходимо получить JWT токен из API авторизации по пути:
``` bash
http://appurl:7540/api/signin
```
Необходимо отправить POST запрос, в котором в body будет передан JSON с содержанием пароля

Далее после получения JWT токена его необходимо задать в переменную token по пути tests/settings.go

Очень важно экспортировать переменную пути DB
``` bash
export TODO_DBFILE="current dbpath location"
go test -v ./tests
```