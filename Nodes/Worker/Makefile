run: main.go heartbeat.go nodeManager.go file_interface.go
	go build -o node_app main.go heartbeat.go nodeManager.go file_interface.go
	./node_app --node-id=node-01 & \
	./node_app --node-id=node-02 & \
	./node_app --node-id=node-03 & \
	./node_app --node-id=node-04 & \
	./node_app --node-id=node-05 & \
	wait