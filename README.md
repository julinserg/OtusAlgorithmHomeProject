![GitHub CI](https://github.com/julinserg/go_home_project/actions/workflows/tests.yml/badge.svg)
# Cервис "Превьювер изображений"

## Общее описание
Сервис предназначен для изготовления preview (создания изображения
с новыми размерами на основе имеющегося изображения).

#### Пример превьюшек в папке [examples](./examples/image-previewer)

## Архитектура
Сервис представляет собой web-сервер (прокси), загружающий изображения,
масштабирующий/обрезающий их до нужного формата и возвращающий пользователю.

## Основной обработчик
http://cut-service.com/fill/300/200/raw.githubusercontent.com/OtusGolang/final_project/master/examples/image-previewer/_gopher_original_1024x504.jpg

<---- микросервис ----><- размеры превью -><--------- URL исходного изображения --------------------------------->

В URL выше мы видим:
- http://cut-service.com/fill/300/200/ - endpoint нашего сервиса,
в котором 300x200 - это размеры финального изображения.
- https://raw.githubusercontent.com/OtusGolang/final_project/master/examples/image-previewer/_gopher_original_1024x504.jpg - 
адрес исходного изображения; сервис скачивает его, производит resize, кэширует и отдает клиенту.

Сервис получает URL исходного изображения, скачивает его, изменяет до необходимых размеров и возвращает как HTTP-ответ.

- Работает только с HTTP.
- Ошибки удалённого сервиса проксирует как есть.
- Поддерживает только JPEG.

Сервис сохраняет (кэширует) полученное preview на локальном диске и при повторном запросе
отдает изображение с диска, без запроса к удаленному HTTP-серверу.

Поскольку размер места для кэширования ограничен, то для удаления редко используемых изображений
используется алгоритм **"Least Recent Used"**.

## Конфигурация
Основной параметр конфигурации сервиса - разрешенный размер LRU-кэша.

Он измеряется количеством закэшированных изображений.

## Развертывание
Развертывание микросервиса должно осуществляться командой `make up` (внутри `docker compose up`)
в директории с проектом.

## Тестирование
Реализация алгоритма LRU покрыта unit-тестами.

Для интеграционного тестирования используется контейнер с Nginx в качестве удаленного HTTP-сервера,
раздающего заданный набор изображений.

Проверена работа сервера в разных сценариях:
* картинка найдена в кэше;
* удаленный сервер не существует;
* удаленный сервер существует, но изображение не найдено (404 Not Found);
* удаленный сервер существует, но изображение не изображение, а скажем, exe-файл;
* удаленный сервер вернул ошибку;
* удаленный сервер вернул изображение;
* изображение меньше, чем нужный размер;
