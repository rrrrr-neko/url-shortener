# URL Shortener (Go)

Мини-сервис для сокращения ссылок на Go: вводишь длинный URL → получаешь короткую ссылку вида `http://localhost:8080/wow/123` → при переходе происходит редирект на исходный URL.

## Возможности

- Web UI (форма в браузере)
- Генерация короткого ключа (3 цифры: `000–999`)
- Редирект по короткому URL
- Статика подключена через `/static/*`
- Хранение ссылок **в памяти** (in-memory map)

## Демо маршруты

- `GET /` — главная страница с формой
- `POST /shorten` — принимает `url` из формы и возвращает короткую ссылку
- `GET /wow/{key}` — редирект на исходный URL по ключу
- `GET /static/style.css` — стили

## Стек

- Go (`net/http`, `html/template`)
- HTML шаблон (`templates/index.html`)
- SCSS/CSS (стили лежат в `static/`)
- Bootstrap подключен через CDN

## Структура проекта

> В коде используется `template.ParseFiles("templates/index.html")`, а статика берётся из `./static`.

## Запуск

### 1) Подготовь папки и разложи файлы

```bash
mkdir -p templates static
# перемести index.html в templates/
# перемести style.css/style.scss/reset.scss/style.css.map в static/

