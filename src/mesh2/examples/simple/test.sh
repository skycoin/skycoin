go run test_node.go -config node_b.json &> ./node_b.log &
go run test_node.go -config node_c.json &> ./node_c.log &

go run test_node.go -config node_a.json 