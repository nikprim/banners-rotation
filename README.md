# Banners-rotation
[![Go Report Card](https://goreportcard.com/badge/nikprim/banners-rotation)](https://goreportcard.com/report/nikprim/banners-rotation)
[![CI](https://github.com/nikprim/banners-rotation/actions/workflows/tests.yml/badge.svg)](https://github.com/nikprim/banners-rotation/actions/workflows/tests.yml)
## Начало
### Скопировать проект
```git clone git@github.com:nikprim/banners-rotation.git```
### Настроить env
`.env #скопировать с .env.dist` <br/>
`.env.test #скопировать с .env.test.dist`
### Команды Makefile
```
make run - запустить приложение
make build - сборка приложения
make up - поднять контейнеры
make down - отключить контейнеры
make stop - остоновить контейнеры
make migrate - запустить миграции
make test - запустить тесты
make lint - запустить линтер
make generate - сгенерировать протобаф
make integration-tests - запустить интеграционные тесты
```
