goose reset
goose up
valkey-cli -p 6400 FLUSHALL

BODY='{"name": "User1",  "email": "user1@mail.com",  "password": "password", "phone": "7001234567"}'
curl -X POST -i -H "Content-Type: application/json" -d "$BODY" localhost:8088/register


BODY='{ "email": "user1@mail.com",  "password": "password"}'
curl -X POST -i -H "Content-Type: application/json" -d "$BODY" localhost:8088/login
