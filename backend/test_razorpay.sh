#!/bin/bash
set -e

if [ -z "$TOKEN" ]; then
  echo "ERROR: TOKEN environment variable not set"
  echo "Usage: TOKEN=your_jwt ./test_razorpay.sh"
  exit 1
fi

BASE_URL="http://localhost:8088"
AMOUNT=5000

if [ -f .env ]; then
  export $(grep -v '^#' .env | grep RAZORPAY_WEBHOOK_SECRET | xargs)
fi

if [ -z "$RAZORPAY_WEBHOOK_SECRET" ]; then
  echo "❌ ERROR: RAZORPAY_WEBHOOK_SECRET not found in .env"
  exit 1
fi

WEBHOOK_SECRET="$RAZORPAY_WEBHOOK_SECRET"
echo "Using webhook secret: $WEBHOOK_SECRET"

ORDER_RESPONSE=$(curl -s -X POST "$BASE_URL/wallet/create-topup-order" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"amount\": $AMOUNT}")

echo "Response: $ORDER_RESPONSE"

ORDER_ID=$(echo "$ORDER_RESPONSE" | jq -r ".order_id")
if [[ -z "$ORDER_ID" || "$ORDER_ID" == "null" ]]; then
  echo "could not extract order_id"
  exit 1
fi

echo "created order: $ORDER_ID"

echo "generating payload and signature..."

PAYLOAD="{\"event\":\"payment.captured\",\"payload\":{\"payment\":{\"entity\":{\"id\":\"pay_test_123\",\"order_id\":\"$ORDER_ID\",\"status\":\"captured\"}}}}"

echo "Payload: $PAYLOAD"
echo ""

SIGNATURE=$(printf '%s' "$PAYLOAD" | openssl dgst -sha256 -hmac "$WEBHOOK_SECRET" -binary | xxd -p -c 256 | tr -d '\n')

echo "Generated signature: $SIGNATURE"
echo ""

echo "Payload hex (for debugging):"
printf '%s' "$PAYLOAD" | xxd -p | tr -d '\n'
echo ""
echo ""

###########################################
# SEND WEBHOOK
###########################################
echo "🔵 3. Sending webhook to backend..."
echo "-----------------------------------"

WEBHOOK_RESPONSE=$(curl -s -X POST "$BASE_URL/wallet/webhook" \
  -H "Content-Type: application/json" \
  -H "X-Razorpay-Signature: $SIGNATURE" \
  -d "$PAYLOAD")

echo "Webhook response: $WEBHOOK_RESPONSE"

###########################################
# CHECK WALLET BALANCE
###########################################
echo ""
echo "🔵 4. Checking wallet balance..."
echo "-----------------------------------"

BALANCE_RESPONSE=$(curl -s -X GET "$BASE_URL/wallet/balance" \
  -H "Authorization: Bearer $TOKEN")

echo "💰 Wallet Balance:"
echo "$BALANCE_RESPONSE"
echo ""
echo "🎉 Test completed!"
