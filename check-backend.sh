#!/usr/bin/env bash
set -Eeuo pipefail

project_dir="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
base_url="${BASE_URL:-http://localhost:8080}"
base_url="${base_url%/}"
env_file="${ENV_FILE:-$project_dir/backend/.env}"

for dependency in curl jq docker; do
  command -v "$dependency" >/dev/null || { printf 'Нужна команда: %s\n' "$dependency" >&2; exit 1; }
done
[[ -f "$env_file" ]] || { printf 'Не найден env-файл: %s\n' "$env_file" >&2; exit 1; }

compose=(docker compose --project-directory "$project_dir" -f "$project_dir/docker-compose.yaml" --env-file "$env_file")
run_id="$(cat /proc/sys/kernel/random/uuid)"
run_id="${run_id//-/}"
admin_email="check-admin-$run_id@example.com"
user_email="check-user-$run_id@example.com"
movie_name="check-movie-$run_id"
password="Check-$run_id"
work_dir="$(mktemp -d)"
checks=0
cleanup_needed=false

sql() {
  "${compose[@]}" exec -T db sh -c 'exec psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" -v ON_ERROR_STOP=1' >/dev/null
}

cleanup() {
  local result=$?
  trap - EXIT
  if [[ "$cleanup_needed" == true ]]; then
    if ! sql <<SQL
DELETE FROM movies WHERE name = '$movie_name';
DELETE FROM users WHERE email IN ('$admin_email', '$user_email');
SQL
    then
      printf 'Не удалось удалить тестовые данные: %s\n' "$run_id" >&2
      result=1
    fi
  fi
  rm -rf -- "$work_dir"
  exit "$result"
}
trap cleanup EXIT
trap 'exit 130' INT
trap 'exit 143' TERM

request() {
  local method=$1 path=$2 expected=$3 token=${4:-} body=${5:-} actual
  local args=(--silent --show-error --connect-timeout 5 --max-time 20 -X "$method" -o "$work_dir/response.json" -w '%{http_code}')
  [[ -z "$token" ]] || args+=(-H "Authorization: Bearer $token")
  [[ -z "$body" ]] || args+=(-H 'Content-Type: application/json' --data "$body")
  actual="$(curl "${args[@]}" "$base_url$path")"
  if [[ "$actual" != "$expected" ]]; then
    printf 'FAIL %s %s: HTTP %s, ожидался %s\n' "$method" "$path" "$actual" "$expected" >&2
    cat "$work_dir/response.json" >&2
    printf '\n' >&2
    exit 1
  fi
  if jq -e 'any(.. | objects; has("password_hash"))' "$work_dir/response.json" >/dev/null 2>&1; then
    printf 'FAIL: ответ содержит password_hash\n' >&2
    exit 1
  fi
  checks=$((checks + 1))
  printf 'OK %s %s → %s\n' "$method" "$path" "$actual"
}

field() {
  jq -er "$1" "$work_dir/response.json"
}

assert_json() {
  jq -e "$1" "$work_dir/response.json" >/dev/null || { printf 'FAIL: неверное содержимое ответа\n' >&2; cat "$work_dir/response.json" >&2; exit 1; }
}

request GET /healthz 200
assert_json '.status == "ok"'
sql <<< 'SELECT 1;'
cleanup_needed=true

user_body="$(jq -nc --arg username "user_$run_id" --arg email "$user_email" --arg password "$password" '{username:$username,first_name:"Test",second_name:"User",email:$email,password:$password}')"
"${compose[@]}" exec -T app /app/traci create-admin >/dev/null <<INPUT
admin_$run_id
Test
Admin
$admin_email
$password
INPUT
checks=$((checks + 1))
printf 'OK CLI create-admin\n'
request POST /auth/register 201 '' "$user_body"
user_id="$(field '.id')"
request POST /auth/register 409 '' "$user_body"
request POST /auth/register 400 '' '{}'

admin_login="$(jq -nc --arg email "$admin_email" --arg password "$password" '{email:$email,password:$password}')"
user_login="$(jq -nc --arg email "$user_email" --arg password "$password" '{email:$email,password:$password}')"
request POST /auth/login 200 '' "$admin_login"
admin_token="$(field '.access_token')"
admin_id="$(jq -nr --arg token "$admin_token" '$token | split(".")[1] | gsub("-"; "+") | gsub("_"; "/") | @base64d | fromjson | .sub')"
request POST /auth/login 200 '' "$user_login"
user_token="$(field '.access_token')"
request POST /auth/login 401 '' "$(jq -nc --arg email "$user_email" '{email:$email,password:"wrong-password"}')"
request GET /catalog 401
request GET /catalog 401 invalid-token
request GET /users 403 "$user_token"

movie_body="$(jq -nc --arg name "$movie_name" '{name:$name,duration_min:120,genre:"DRAMA",release_year:2020,description:"Test movie"}')"
request POST /catalog 403 "$user_token" "$movie_body"
request POST /catalog 400 "$admin_token" '{"name":"Bad","duration_min":0,"genre":"DRAMA","release_year":2020}'
request POST /catalog 201 "$admin_token" "$movie_body"
movie_id="$(field '.id')"
request GET '/catalog?limit=10&offset=0' 200 "$user_token"
assert_json '.limit == 10 and .offset == 0 and (.items | type == "array")'
request GET "/catalog?genre=DRAMA&search=$movie_name&limit=10&offset=0" 200 "$user_token"
assert_json '.total == 1 and (.items | length == 1)'
request GET '/catalog?limit=0' 400 "$user_token"
request GET /catalog/not-a-uuid 400 "$user_token"
request GET "/catalog/$movie_id" 200 "$user_token"
request PATCH "/catalog/$movie_id" 200 "$admin_token" '{"duration_min":130,"description":null}'
assert_json '.duration_min == 130 and .description == null'

collection_body="$(jq -nc --arg id "$movie_id" '{movie_id:$id,status:"PLANNED"}')"
request POST /collection 201 "$user_token" "$collection_body"
request POST /collection 409 "$user_token" "$collection_body"
request GET '/collection?status=PLANNED&limit=10&offset=0' 200 "$user_token"
assert_json '.total == 1 and (.items | length == 1)'
request GET "/collection/$movie_id" 200 "$user_token"
request GET "/collection/$movie_id" 404 "$admin_token"
request PATCH "/collection/$movie_id" 400 "$user_token" '{"personal_rating":11}'
request PATCH "/collection/$movie_id" 200 "$user_token" '{"status":"WATCHED","personal_rating":9,"review":"Хороший фильм"}'
assert_json '.status == "WATCHED" and .personal_rating == 9'
request PATCH "/collection/$movie_id" 200 "$user_token" '{"personal_rating":null,"review":null}'
assert_json '.personal_rating == null and .review == null'
request DELETE "/collection/$movie_id" 204 "$user_token"
request GET "/collection/$movie_id" 404 "$user_token"
request POST /collection 201 "$user_token" "$collection_body"

request GET '/users?limit=10&offset=0' 200 "$admin_token"
request GET "/users/$user_id" 200 "$admin_token"
assert_json '.role == "user" and .is_blocked == false'
request PATCH "/users/$admin_id" 409 "$admin_token" '{"is_blocked":true}'
request DELETE "/users/$admin_id" 409 "$admin_token"
request PATCH "/users/$user_id" 200 "$admin_token" '{"is_blocked":true}'
request GET /collection 403 "$user_token"
request POST /auth/login 403 '' "$user_login"
request PATCH "/users/$user_id" 200 "$admin_token" '{"is_blocked":false}'
request GET /collection 200 "$user_token"

request DELETE "/catalog/$movie_id" 204 "$admin_token"
request GET "/catalog/$movie_id" 404 "$user_token"
request GET "/collection/$movie_id" 404 "$user_token"
request DELETE "/users/$user_id" 204 "$admin_token"
request GET "/users/$user_id" 404 "$admin_token"
request GET /catalog 401 "$user_token"
printf '\nВсе проверки прошли: %d.\n' "$checks"
