goose reset
goose up
valkey-cli -p 6400 FLUSHALL
BODY='{"name": "Jane Doe",  "email": "jaapm@ail.com",  "password": "adlkfalsdfjldj", "phone": "7017105448"}'
curl -X POST -i -H "Content-Type: application/json" -d "$BODY" localhost:8088/register
valkey -p 6400 LLEN 

