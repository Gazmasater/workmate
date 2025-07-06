## Быстрый старт

- git clone  https://github.com/Gazmasater/workmate/tree/iter1
- cd workmate
- Запуск сервера
  - cd cmd/server
  - go run .
- в браузере набрать для визуализвцмм сваггера http://localhost:8080/swagger/index.html


По умолчанию сервер слушает порт 8080.

## Функционал

- REST API для задач (Task)
- CRUD-операции:

  - GET /Health - проверка сервиса
    - curl -v http://localhost:8080/health 
  - POST /tasks — создать задачу
    - curl -X POST http://localhost:8080/tasks 
  - GET /tasks — получить список всех задач
    - curl http://localhost:8080/tasks
    - Поддерживает фильтрацию и пагинацию(примеры можно посмотреть в сваггере)
  - GET /tasks/{id} — получить задачу по ID
    - curl http://localhost:8080/tasks/{id}
      - пример:
curl http://localhost:8080/tasks/123e4567-e89b-12d3-a456-426614174000
   - DELETE /tasks/{id} — удалить задачу
    - curl -X DELETE http://localhost:8080/tasks/{id} 
  - PUT /tasks/{id}/cancel - отмена задачи
    - curl -X PUT http://localhost:8080/tasks/{id}/cancel
      - пример:
curl -X PUT http://localhost:8080/tasks/123e4567-e89b-12d3-a456-426614174000/cancel
 
- Хранилище: оперативная память (in-memory)
- Структурированное логирование через кастомный логгер на основе zap
- Конфигурирование через переменные окружения
- Расширяемость: можно внедрить любую БД

## Профилирование и оптимизация

- В проекте поддерживается профилирование с помощью [pprof](https://pkg.go.dev/net/http/pprof).  
  Для анализа производительности используются стандартные инструменты Go:  
  `go test -run TestInMemoryRepo_Concurrency -cpuprofile=cpu.out -memprofile=mem.out`
- Для просмотра профилей:
  ```sh
  go tool pprof cpu.out
  go tool pprof mem.out
Архитектурные решения
In-memory репозиторий реализован с использованием шардинга (sharded map, 16-32 shard’а), что существенно уменьшает lock contention при высокой конкурентной нагрузке.

Оптимизация позволила снизить время выполнения теста в ~5 раз по сравнению с вариантом с одним mutex.

## CI/CD

- Для проекта настроен Continuous Integration 

