# Image Previewer
## Сборка и запуск
```bash
make run
```
## URL
```bash
http://localhost:8000/600/400/{URL}
```
Можно развернуть docker-контейнеры через ```docker-compose up```. Тогда развернется контейнер с приложением и nginx-контейнер с примером картинки. Тогда запустить приложение можно через
```bash
http://localhost:8000/600/400/http://image-previewer-nginx/picture.jpg
```