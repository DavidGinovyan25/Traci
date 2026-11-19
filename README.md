# Traci

Трекер фильмов с каталогом, личными коллекциями, оценками и отзывами.

Бэкенд на Go с PostgreSQL реализует REST API, регистрацию, JWT-вход и управление пользователями. Фронтенд на HTML, CSS и JavaScript позволяет работать с каталогом, личной коллекцией и админкой.

Результаты ДЗ №3, четыре LogsQL-запроса со скриншотами и проверка observability описаны в [Отчёт.md](Отчёт.md). Ниже приведены инструкции запуска приложения и компонентов логирования.

Для запуска через Docker скопируйте `backend/.env.example` в `backend/.env` и заполните настройки: шаблон содержит только названия переменных. `JWT_SECRET` должен быть длиной минимум 32 байта. Приложение и мигратор читают `POSTGRES_URL` из этого файла; для Compose хост БД должен быть `db`. Параметры `POSTGRES_DB`, `POSTGRES_USER`, `POSTGRES_PASSWORD` должны соответствовать URL:

```sh
make docker-up
```

Обычный `docker compose up` не читает `backend/.env`: Compose по умолчанию ищет `.env` в корне проекта. Цель `make docker-up` передаёт нужный файл явно. Остановка: `make docker-down`.

База и миграции запускаются автоматически. Интерфейс доступен на **http://localhost:3000**, API — на `http://localhost:8080`. Зарегистрируйся и войди; фильмы в каталог добавляет администратор.

Без Docker:

```sh
make build
make run
```

Перед запуском примените миграции и настройте `POSTGRES_URL` и `JWT_SECRET`.

```sh
make generate
make test
make test-unit
make test-race
```

Генерация требует sqlc v1.31.1. Сгенерированные файлы включены в проект; для обычной сборки повторная генерация не требуется.

Контракт API: `docs/api/openapi.yaml`. Схема БД: `docs/db/schema.dbml`.

Полные тесты требуют работающий Docker: testcontainers-go запускает отдельный PostgreSQL 17 и удаляет его после проверки. `make test-unit` запускает тесты без Docker.

Устройство интерфейса и инструкция по его проверке: [frontend/README.md](frontend/README.md). `make test-frontend` проверяет клиентскую логику; `make test-e2e` запускает сценарии в Chromium после установки тестовых зависимостей.

В minikube аккаунт администратора создаётся командой `make create-admin`. Она запросит данные нового аккаунта; почта должна быть свободна. Для Docker Compose используйте `docker compose --env-file backend/.env exec app /app/traci create-admin`.


## Приложение в minikube

Нужны Docker, minikube, Helm и GNU Make. Команды выполняются из корня проекта.

```bash
minikube start
minikube addons enable ingress
make minikube-load
```

`make minikube-load` собирает `traci-backend:v1` и `traci-frontend:v1` и загружает их в minikube. Перед развёртыванием заполните локальный `k8s/app/secret.yaml`: пароль PostgreSQL и JWT-секрет не короче 32 байт. Проверьте обычные настройки в `k8s/app/configmap.yaml`.

```bash
./deploy.sh
minikube kubectl -- get pods,services,ingresses,pvc -n traci
minikube ip
```

Скрипт применяет манифесты из `k8s/app/`, ожидает PostgreSQL, запускает миграции, затем backend, frontend и Ingress. Откройте `http://<IP-minikube>`, если IP узла доступен с хоста. Альтернативный доступ:

```bash
minikube kubectl -- port-forward -n traci service/traci-frontend 3001:80
```

В этом случае интерфейс доступен по <http://localhost:3001>. Порт 3001 выбран, чтобы не конфликтовать с Grafana.

## Логирование и алертинг в minikube

Приложение должно уже работать в том же minikube. Проверьте контекст Helm:

```bash
helm list --kube-context minikube -A
```

`logging.sh` использует текущий Kubernetes-контекст. Если он отличается от minikube, можно запускать скрипт как `HELM_KUBECONTEXT=minikube ./logging.sh`.

### Telegram Secret

Создайте бота через @BotFather, отправьте ему `/start` и получите ID своего чата. Перед установкой Alertmanager создайте Secret. Команды ниже рассчитаны на Bash; ввод токена скрыт.

```bash
minikube kubectl -- create namespace observability --dry-run=client -o yaml | minikube kubectl -- apply -f -
read -rsp 'Telegram bot token: ' TELEGRAM_BOT_TOKEN; echo
read -rp 'Telegram chat ID: ' TELEGRAM_CHAT_ID
printf '%s' "$TELEGRAM_BOT_TOKEN" | minikube kubectl -- create secret generic alertmanager-telegram \
  -n observability --from-file=bot-token=/dev/stdin \
  --from-literal=chat-id="$TELEGRAM_CHAT_ID" --dry-run=client -o yaml | minikube kubectl -- apply -f -
unset TELEGRAM_BOT_TOKEN TELEGRAM_CHAT_ID
```

Токен и ID чата не записываются в YAML проекта. При смене токена повторите создание Secret и перезапустите Alertmanager:

```bash
minikube kubectl -- rollout restart statefulset/alertmanager -n observability
```

### Установка компонентов

```bash
./logging.sh
```

В полном репозитории также доступна цель `make minikube-logging`. Для отдельного архива ДЗ №3 достаточно `logging.sh`, Helm values и заранее развёрнутого приложения.

Скрипт устанавливает или обновляет VictoriaLogs, Vector, Grafana, Alertmanager и vmalert в namespace `observability`, передаёт values и ожидает готовности. Конфиги находятся в `k8s/observability/`:

| Файл | Компонент |
|---|---|
| `victoria-logs.yaml` | Хранилище логов |
| `vector-values.yaml` | Сбор и доставка логов |
| `grafana.yaml` | Grafana, плагин и источник данных |
| `alertmanager.yaml` | Получатель Telegram и монтирование Secret |
| `vmalert.yaml` | Источник данных, адрес Alertmanager и правила |

Проверка состояния и логов:

```bash
minikube kubectl -- get deployments,daemonsets,statefulsets,pods,services,pvc -n observability
minikube kubectl -- logs -n observability -l app.kubernetes.io/name=vector --tail=30
minikube kubectl -- logs -n observability -l app.kubernetes.io/instance=vma --tail=30
minikube kubectl -- logs -n observability alertmanager-0 -c alertmanager --tail=30
```

### Открытие интерфейсов

Каждую нужную команду запустите в отдельном терминале и оставьте работающей. `Ctrl+C` закрывает соответствующий проброс.

```bash
minikube kubectl -- port-forward -n observability service/vls-victoria-logs-single-server 9428:9428
minikube kubectl -- port-forward -n observability service/grafana 3000:80
minikube kubectl -- port-forward -n observability service/vma-victoria-metrics-alert-server 8880:8880
minikube kubectl -- port-forward -n observability service/alertmanager 9093:9093
```

| Интерфейс | Адрес |
|---|---|
| vmui | <http://localhost:9428/vmui/> |
| Grafana | <http://localhost:3000> |
| vmalert | <http://localhost:8880/vmalert/> |
| Alertmanager | <http://localhost:9093> |

Логин Grafana — `admin`. Получить пароль локально:

```bash
minikube kubectl -- get secret grafana -n observability -o jsonpath='{.data.admin-password}' | base64 --decode; echo
```

В Grafana откройте Explore и выберите VictoriaLogs. Для событий backend используйте `kubernetes.container_name:="traci-backend" | unpack_json`. Четыре запроса для vmui и результаты приведены в отчёте.

Grafana использует временное хранилище: созданные вручную дашборды могут потеряться после пересоздания Pod. Плагин и источник данных восстанавливаются из values. Данные VictoriaLogs находятся на PVC; изменение расположения конфигов не удаляет ресурсы кластера.

### Проверка алерта

`TooManyLogs` проверяется раз в минуту и срабатывает при более чем 100 событиях backend за минуту, если условие сохраняется одну минуту. Обычные действия в UI могут не создать столько событий. Состояние правила видно в vmalert, полученный алерт — в Alertmanager, уведомление — в Telegram. После прекращения событий правило разрешается при следующих вычислениях; отправка уведомления о разрешении включена.

Для временной остановки всего кластера используйте `minikube stop`, для запуска — `minikube start`. Данные PVC сохраняются до удаления соответствующего хранилища.
