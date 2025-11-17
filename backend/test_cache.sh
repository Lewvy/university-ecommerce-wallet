#!/bin/bash

TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJzY29wZSI6ImF1dGhlbnRpY2F0aW9uIiwiZXhwIjoxNzYzMzc1MjY2LCJpYXQiOjE3NjMzNzQzNjZ9.mMcnF0gVW76nhR6BrLILyUDOuqjAkkfhJeayGNDMB_c"
BASE_URL="http://localhost:8088"

echo "--- 1. Adding 2 of Product 1 to cart ---"
curl -X POST "$BASE_URL/cart/add" \
     -H "Authorization: Bearer $TOKEN" \
     -H "Content-Type: application/json" \
     -d '{"product_id": 1, "quantity": 2}'
echo -e "\n"


echo "--- 2. Getting cart details ---"
curl -X GET "$BASE_URL/cart" \
     -H "Authorization: Bearer $TOKEN"
echo -e "\n"


echo "--- 3. Adding 1 of Product 2 to cart ---"
curl -X POST "$BASE_URL/cart/add" \
     -H "Authorization: Bearer $TOKEN" \
     -H "Content-Type: application/json" \
     -d '{"product_id": 2, "quantity": 1}'
echo -e "\n"


echo "--- 4. Getting cart details again ---"
curl -X GET "$BASE_URL/cart" \
     -H "Authorization: Bearer $TOKEN"
echo -e "\n"


echo "--- 5. Updating Product 1 quantity to 5 ---"
curl -X PUT "$BASE_URL/cart/update" \
     -H "Authorization: Bearer $TOKEN" \
     -H "Content-Type: application/json" \
     -d '{"product_id": 1, "quantity": 5}'
echo -e "\n"


echo "--- 6. Deleting Product 2 from cart ---"
curl -X DELETE "$BASE_URL/cart/item/2" \
     -H "Authorization: Bearer $TOKEN"
echo -e "\n"


echo "--- 7. Getting final cart details ---"
curl -X GET "$BASE_URL/cart" \
     -H "Authorization: Bearer $TOKEN"
echo -e "\n"


echo "--- 8. Clearing the cart ---"
curl -X DELETE "$BASE_URL/cart/clear" \
     -H "Authorization: Bearer $TOKEN"
echo -e "\n"


echo "--- 9. Verifying cart is empty ---"
curl -X GET "$BASE_URL/cart" \
     -H "Authorization: Bearer $TOKEN"
echo -e "\n"
