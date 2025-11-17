
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJzY29wZSI6ImF1dGhlbnRpY2F0aW9uIiwiZXhwIjoxNzYzMzcwNTUxLCJpYXQiOjE3NjMzNjk2NTF9.G6zNRNs6a5slzlK9YrfV4GXSjcGkrywTCu3HP6qJTrM"

curl -X POST "http://localhost:8088/products" \
-H "Authorization: Bearer $TOKEN" \
-F "name=Stylish Cotton T-Shirt" \
-F "description=A comfortable and high-quality cotton t-shirt." \
-F "condition=new" \
-F "category=apparel" \
-F "price=2500" \
-F "stock=100" \
-F "images=@img1.png" \
-F "images=@img2.png" &


curl -X GET 'http://localhost:8088/wallet/balance' \
-H "Authorization: Bearer $TOKEN"

