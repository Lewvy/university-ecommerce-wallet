#!/bin/bash

TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJzY29wZSI6ImF1dGhlbnRpY2F0aW9uIiwiZXhwIjoxNzYzMzc2MzYzLCJpYXQiOjE3NjMzNzU0NjN9.sjXEMPNsuuO5uyfZ6y2kj-EWxidIkp-XXDoBDb59Wms"
BASE_URL="http://localhost:8088"

echo "--- Attempting to create order from cart ---"
curl -X POST "$BASE_URL/orders" \
     -H "Authorization: Bearer $TOKEN" \
     -H "Content-Type: application/json"

echo -e "\n"
