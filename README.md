# subscription-aggregation-service
[![forthebadge](https://forthebadge.com/images/badges/made-with-go.svg)](https://forthebadge.com)
[![forthebadge](https://forthebadge.com/images/badges/built-with-love.svg)](https://forthebadge.com)

## 📑 Содержание

- [Описание](#-описание)
- [Команды](#-команды)
- [Тестирование API](#-тестирование-api)

## 📖 Описание

Этот проект представляет собой реализацию REST-сервис для агрегации данных об онлайн-подписках пользователей.

**Функциональность:**
- Создание новой подписки
- Получение информации о конкретной подписке
- Редактирование существующей подписки
- Удаление подписки
- Получение списка всех подписок с возможностью фильтрации по ID пользователя, названию сервиса и промежутку действия подписки
- Расчёт суммарной стоимости подписок с возможностью фильтрации по пользователю и названию сервиса (подсчёт учитывает пересечение периода действия подписки с указанным интервалом)

## ⚙️ Команды
### Запуск
```bash
git clone git@github.com:kasparovgs/subscription-aggregation-service.git
cd subscription-aggregation-service
make launch_services
```
### Остановка
```bash
make stop_services
```
## 🔗 Тестирование API

| 📜 **Swagger UI** | [http://localhost:8080/swagger/docs/index.html](http://localhost:8080/swagger/docs/index.html) | Просмотр и тестирование API |