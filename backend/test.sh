
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJzY29wZSI6ImF1dGhlbnRpY2F0aW9uIiwiZXhwIjoxNzYzMjA5MjQ4LCJpYXQiOjE3NjMyMDgzNDh9.IW8G_c_T2X16MCTW5SOqGzalz2Jj8_vZyuqVfCAmVy0"

curl -X POST "http://localhost:8088/products" \
-H "Authorization: Bearer $TOKEN" \
-F "name=Stylish Cotton T-Shirt" \
-F "description=A comfortable and high-quality cotton t-shirt." \
-F "price=2500" \
-F "stock=100" \
-F "images=@img1.png" \
-F "images=@img2.png"


curl -X GET 'http://localhost:8088/wallet/balance' \
-H "Authorization: Bearer $TOKEN"

