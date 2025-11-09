goose reset
goose up
valkey-cli -p 6400 FLUSHALL
for i in {1..10}
do
  BODY='{"name": "User'$i'",  "email": "user'$i'@mail.com",  "password": "password", "phone": "701710544'$i'"}'
curl -X POST -i -H "Content-Type: application/json" -d "$BODY" localhost:8088/register
done

# for i in {10..100}
# do
#   BODY='{"name": "User'$i'",  "email": "user'$i'@mail.com",  "password": "password'$i'", "phone": "70000000'$i'"}'
# curl -X POST -i -H "Content-Type: application/json" -d "$BODY" localhost:8088/register
# done
# for i in {100..1000}
# do
#   BODY='{"name": "User'$i'",  "email": "user'$i'@mail.com",  "password": "password'$i'", "phone": "7000000'$i'"}'
# curl -X POST -i -H "Content-Type: application/json" -d "$BODY" localhost:8088/register
# done
