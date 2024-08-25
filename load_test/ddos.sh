PROXY_SERVER_IP="http://localhost:7001"
ACTUAL_SERVER_IP="http://localhost:7000"

# Loop 20 times
for i in {1..20}
do
    echo "Call #$i"
    curl "{$PROXY_SERVER_IP}/test1"
    echo ""
done