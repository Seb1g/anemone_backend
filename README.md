# Anemone Notes Backend

    XXXX   XXX   XX   XXXXXXXX   XXX   XXX    XXXXX     XXX    XX   XXXXXXXX
   XX XX   XXXX  XX   XX         XXXX XXXX   XX    XX   XXXX   XX   XX      
XXX   XX   XX XX XX   XXXXXXXX   XX XXX XX   XX    XX   XX XX  XX   XXXXXXXX
 XXXXXXXX  XX  XXXX   XXXXXXXX   XX  X  XX   XX    XX   XX  XX XX   XXXXXXXX
XX    XX   XX   XXX   XX         XX     XX   XX    XX   XX   XXXX   XX      
XX    XX   XX    XX   XXXXXXXX   XX     XX    XXXXX     XX    XXX   XXXXXXXX
                                                                            
                                                                            
                                                                            
   XX    XX   XXXXXXX   XX        XXXXXX    XXXXXXX   XX   XX   XXXXXXX     
   XX    XX   XXXXXXX   XX        XXXXXXX   XXXXXXX   XXX  XX   XXXXXXX     
   XXXXXXXX   XX   XX   XX        XX   XX     XXX     XXXX XX   XX          
   XXXXXXXX   XX   XX   XX        XX  XXX     XXX     XX XXXX   XX   XX     
   XX    XX   XXXXXXX   XXXXXXX   XXXXXXX   XXXXXXX   XX  XXX   XXXXXXX     
   XX    XX   XXXXXXX   XXXXXXX   XXXXX     XXXXXXX   XX   XX   XXXXXXX     

**Anemone Notes Backend** — это монолитный сервис для управления заметками (временно только ими).  
Проект написан на **Go** и использует **PostgreSQL** в качестве основного хранилища данных.  

---

## Технологии

| Компонент      | Библиотека                                         | Назначение                                    |
|----------------|----------------------------------------------------|-----------------------------------------------|
| Язык           | Go (1.25.1)                                        | Основной язык разработки                      |
| Роутинг        | [gorilla/mux](https://github.com/gorilla/mux)      | HTTP-роутер и диспетчер маршрутов             |
| Работа с БД    | [sqlx](https://github.com/jmoiron/sqlx)            | Упрощённая работа с SQL-запросами             |
| Драйвер БД     | [lib/pq](https://github.com/lib/pq)                | PostgreSQL-драйвер                            |
| Конфигурация   | [godotenv](https://github.com/joho/godotenv)       | Загрузка переменных окружения из `.env`       |
| JWT            | [golang-jwt](https://github.com/golang-jwt/jwt)    | Генерация и проверка токенов                  |
| Хеширование    | [x/crypto](https://pkg.go.dev/golang.org/x/crypto) | Безопасное хеширование паролей                |

---

## Архитектура

Проект реализован по принципам **Layered Architecture**:

- **Handler/API** → Обработка HTTP-запросов (JSON, валидация, ответы).  
- **Service/Usecase** → Бизнес-логика и координация данных.  
- **Repository/Repo** → Взаимодействие с базой данных (SQL-запросы).  

## Дорожная карта

В будущем сервис будет расширяться и разделяться на отдельные модули для поддержки приложений:

- **Anemone Kanban**
- **Anemone Mail**
- **Anemone Quiz**