# Миграции

Пошаговая инструкция, как в этом шаблоне создавать, заполнять и прогонять
миграции схемы БД. Механизм — версионированные миграции через
[gormigrate](https://github.com/go-gormigrate/gormigrate), прогоняются
`PostgresKernel` при старте приложения, до того как остальные kernels
(HTTP, Rabbit) начнут работать.

## Содержание

1. [Как это устроено](#1-как-это-устроено)
2. [Шаг 1 — сгенерировать миграцию](#шаг-1--сгенерировать-миграцию)
3. [Шаг 2 — заполнить тело миграции](#шаг-2--заполнить-тело-миграции)
4. [Шаг 3 — подключить в app.go](#шаг-3--подключить-в-appgo-один-раз-на-проект)
5. [Пример: таблица + сид данных](#пример-таблица--сид-данных)
6. [Откат (rollback)](#откат-rollback)
7. [Как проверить, что миграция уже применена](#как-проверить-что-миграция-уже-применена)
8. [Несколько инстансов одновременно](#несколько-инстансов-одновременно)
9. [Локальная проверка](#локальная-проверка)
10. [Частые ошибки](#частые-ошибки)

---

## 1. Как это устроено

- `internal/migrations/` — пакет со всеми миграциями проекта.
  - `registry.go` — функция `All()`, список миграций в порядке применения.
    Генерируется и дополняется командой `codegen migration add`, руками не
    редактируется.
  - `<timestamp>_<domain>_<module>_<description>.go` — один файл на одну
    миграцию, с двумя функциями: `Migrate` (применить) и `Rollback` (откатить).
- При старте `PostgresKernel.Init()`:
  1. открывает соединение с БД;
  2. берёт Postgres advisory lock (защита от гонки при одновременном старте
     нескольких инстансов — см. [раздел ниже](#несколько-инстансов-одновременно));
  3. прогоняет ещё не применённые миграции по порядку — какие уже применены,
     хранится в служебной таблице `migrations` в самой БД (см.
     [раздел ниже](#как-проверить-что-миграция-уже-применена));
  4. снимает лок и только после этого регистрирует соединение в DI —
     остальные kernels стартуют, когда миграции гарантированно готовы.

---

## Шаг 1 — сгенерировать миграцию

```bash
go run github.com/exgamer/go-sdk-generator/cmd/codegen@latest \
    migration add <domain>/<module> <description>
```

Например, для домена `handbook/city`:

```bash
go run github.com/exgamer/go-sdk-generator/cmd/codegen@latest \
    migration add handbook/city create_city_table
```

`<domain>/<module>` и `<description>` тут — не путь и не ссылка на реальный
код, а просто название файла: генератор ничего не проверяет по файловой
системе (даже если `internal/domains/<domain>/<module>` не существует —
команда всё равно отработает). Единственная валидация — формат: строчные
латинские буквы/цифры/`_`, начинается с буквы, не зарезервированное слово Go
(`create_city_table`, `seed_initial_cities`, но не `create city table`).
По конвенции стоит указывать `domain/module` существующего домена — так
понятнее, к чему относится миграция, — но это только договорённость, а не
требование инструмента.

Команда создаёт:
- `internal/migrations/<timestamp>_<domain>_<module>_<description>.go` —
  заготовку с `Migrate`/`Rollback` (`// TODO: implement`);
- `internal/migrations/registry.go` — создаёт или дополняет `All()` новой
  записью, в конец списка.

Повторные вызовы только добавляют новые файлы и дописывают `All()` —
существующие файлы не трогают.

---

## Шаг 2 — заполнить тело миграции

Открыть сгенерированный файл и заполнить `Migrate`/`Rollback`. `tx` — обычный
`*gorm.DB` (в рамках транзакции этой миграции): можно писать сырой SQL
(`tx.Exec(...)`) или использовать GORM-стиль (`tx.AutoMigrate(&model{})`).

```go
func migration20260824050406() *migration.Migration {
	return &migration.Migration{
		ID: "20260824050406_handbook_city_create_city_table",
		Migrate: func(tx *gorm.DB) error {
			return tx.Exec(`CREATE TABLE city (
				id SERIAL PRIMARY KEY,
				name VARCHAR(100) NOT NULL,
				status INT NOT NULL DEFAULT 0
			)`).Error
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Exec(`DROP TABLE city`).Error
		},
	}
}
```

> `ID` — не трогать, это первичный ключ в служебной таблице `migrations`,
> генератор проставляет его сам из timestamp + domain + module + description.

---

## Шаг 3 — подключить в `app.go` (один раз на проект)

В `internal/app/app.go`:

```go
import "github.com/exgamer/go-sdk-rest-template/internal/migrations"

appInstance.RegisterAndInitKernels(
	(&postgres.PostgresKernel{}).WithMigrations(migrations.All()...),
	&http.HttpKernel{},
	...
)
```

Дальше это делать не нужно: каждая следующая миграция — это просто шаг 1
(`codegen migration add ...`) + шаг 2 (заполнить тело). `app.go` трогать
больше не надо, `migrations.All()` уже подключён и подтягивает новые записи
сама.

---

## Пример: таблица + сид данных

В этом шаблоне уже есть рабочий пример на домене `handbook/city` —
`internal/migrations/`:

1. `20260824050406_handbook_city_create_city_table.go` — создаёт таблицу
   `city` (`CREATE TABLE` / `DROP TABLE`).
2. `20260824050450_handbook_city_seed_initial_cities.go` — сидит в неё
   стартовые данные:

```go
Migrate: func(tx *gorm.DB) error {
	return tx.Exec(`INSERT INTO city (name, status) VALUES
		('Astana', 1),
		('Almaty', 1),
		('Shymkent', 1)`).Error
},
Rollback: func(tx *gorm.DB) error {
	return tx.Exec(`DELETE FROM city WHERE name IN ('Astana', 'Almaty', 'Shymkent')`).Error
},
```

Схема и данные — раздельными миграциями, в порядке применения (сначала
`CREATE TABLE`, потом `INSERT`). Так же стоит поступать и для новых доменов:
одна миграция — одно логическое изменение (создать таблицу, добавить
колонку, засеять данные, и т.д.), не смешивать в одном файле.

---

## Откат (rollback)

`Rollback` вызывается не автоматически — только через `gormigrate` API,
если он используется отдельно от `migration.Run` (в этом шаблоне `Run`
всегда идёт только вперёд, `RollbackLast`/`RollbackTo` не подключены). На
практике `Rollback` в первую очередь документирует, как обратить миграцию,
и пригождается при ручном вмешательстве через psql/консоль. Пишите его
всегда, даже если не будете гонять автоматически.

---

## Как проверить, что миграция уже применена

`gormigrate` хранит состояние в отдельной таблице `migrations` в той же БД
(создаёт сама при первом запуске):

- перед прогоном миграции с `ID = X` выполняется
  `SELECT count(*) FROM migrations WHERE id = 'X'` — если запись есть,
  `Migrate` не вызывается вообще;
- после успешного `Migrate` в той же транзакции делается
  `INSERT INTO migrations (id) VALUES ('X')` — если `Migrate` упал, записи
  не будет, и при следующем запуске миграция попробует применится снова.

Проверить руками:

```sql
SELECT * FROM migrations;
```

---

## Несколько инстансов одновременно

При rolling deploy может подняться несколько инстансов сервиса сразу.
`PostgresKernel` берёт Postgres advisory lock (`pg_advisory_lock`) вокруг
прогона миграций:

- инстанс, который первым взял лок, применяет ещё не применённые миграции;
- остальные инстансы ждут на этом же `SELECT pg_advisory_lock(...)`, и после
  разблокировки видят уже применённое состояние (по таблице `migrations`) —
  просто пропускают все миграции без ошибок.

Никакой отдельной настройки для этого не требуется — включено всегда, как
только передан непустой список в `WithMigrations(...)`.

---

## Локальная проверка

Поднять локальный Postgres (см. `.env` / `.env.example` для параметров
подключения) и прогнать миграции без поднятия всего приложения (в т.ч. без
Rabbit-kernel) — небольшой отдельный `main.go`:

```go
package main

import (
	"log"

	"github.com/exgamer/go-sdk-rest-template/internal/migrations"
	"github.com/exgamer/gosdk-core/pkg/app"
	postgres "github.com/exgamer/gosdk-postgres-core/pkg/app"
)

func main() {
	a := app.NewApp()
	kernel := (&postgres.PostgresKernel{}).WithMigrations(migrations.All()...)

	if err := a.RegisterAndInitKernels(kernel); err != nil {
		log.Fatal(err)
	}

	log.Println("migrations applied")
}
```

Запустить дважды подряд — второй раз ничего не должно примениться (только
`SELECT count(*) FROM migrations WHERE id = ...`, без `CREATE`/`INSERT`).
Запустить несколько раз параллельно на пустой БД — только один инстанс
реально выполнит DDL/DML, остальные подождут на локе и выйдут без ошибок.

---

## Частые ошибки

| Ошибка | Причина | Решение |
|---|---|---|
| `codegen: invalid description "..."` | В `<description>` пробелы/заглавные буквы | Использовать `snake_case`: `create_city_table` |
| `relation "city" already exists` при чистом старте | Таблица создана вручную/раньше, а таблицы `migrations` с записью о ней — нет | Либо удалить таблицу и дать миграции создать её самой, либо создать миграцию, которая просто регистрирует уже существующую схему (`Migrate` — no-op) |
| Миграция подвисла при параллельном старте нескольких инстансов | Нормальное поведение — инстанс ждёт `pg_advisory_lock`, пока первый не закончит | Если висит подозрительно долго — проверить, не упал ли инстанс, державший лок, не сняв его (соединение оборвётся — Postgres освободит advisory lock сам) |
| Правки в `registry.go` руками потерялись | Файл перегенерируется командой `codegen migration add` | Не редактировать `registry.go` руками — добавлять миграции только через генератор |
