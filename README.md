# Файлы для итогового задания

В директории `tests` находятся тесты для проверки API, которое должно быть реализовано в веб-сервере.

Директория `web` содержит файлы фронтенда.

Список выполенных заданий со звёздочкой - не выполнялись

Описание проекта - планировщик задач

Основные функции:

1. Регистрация маршрутов API:
    - /api/nextdate - Получение следующей даты выполнения задачи.
    - /api/task - Добавление новой задачи (POST) и обновление существующей задачи (PUT).
    - /api/tasks - Получение списка всех задач.
    - /api/task - Получение задачи по идентификатору (GET) и удаление задачи (DELETE).
    - /api/task/done - Отметка задачи как выполненной.

2. Обработка API-запросов:
    - HandleNextDate: Обрабатывает запрос для получения следующей даты выполнения задачи, учитывая текущую дату, начальную дату и период повторения.
    - HandleAddTask: Обрабатывает добавление новой задачи. Вызывает метод AddTask из пакета task.
    - getTaskHandler: Получает задачу по идентификатору и возвращает ее данные в формате JSON.
    - updateTaskHandler: Обновляет данные существующей задачи. Проверяет корректность входных данных, и если они валидны, вызывает метод UpdateTask из пакета db.
    - handleTaskDelete: Удаляет задачу из базы данных.
    - handleTaskDone: Отмечает задачу как выполненную и планирует следующую дату выполнения, если задача имеет период повторения.

Пример структуры задачи:

- ID: Уникальный идентификатор задачи.
- Date: Дата выполнения задачи.
- Title: Заголовок задачи.
- Comment: Дополнительный комментарий к задаче.
- Repeat: Период повторения задачи (например, ежедневно, ежегодно и т.д.).

Ключевые пакеты:

- transport: Содержит обработчики HTTP-запросов.
- db: Реализует функции для взаимодействия с базой данных (добавление, обновление, удаление и получение задач).
- task: Реализует логику по работе с задачами, а именно функции добавления и предоставления списка задач.
- utils: Содержит утилиты, такие как функция NextDate, для расчета следующей даты выполнения задачи.

Основные обработчики функций:

- RegisterAPIRoutes: Регистрация всех маршрутов API.
- HandleNextDate: Обработка запроса на получение следующей даты задачи.
- HandleAddTask: Обработка запроса на добавление новой задачи.
- getTaskHandler: Обработка запроса на получение задачи по идентификатору.
- updateTaskHandler: Обработка запроса на обновление задачи.
- handleTaskDelete: Обработка запроса на удаление задачи.
- handleTaskDone: Обработка запроса на отметку задачи как выполненной.

Инструкция по запуску проекта локально

Чтобы запустить проект на вашем компьютере, выполните следующие шаги:

1. Запуск кода:
   Откройте терминал и перейдите в директорию с вашим проектом. Затем введите команду:

   go run main.go


2. Доступ к проекту:
   После успешного запуска, откройте веб-браузер и перейдите по следующему адресу:

   http://localhost:7540/


3. Запуск тестов:
   Для выполнения тестов используйте команду:

   go test ./tests

