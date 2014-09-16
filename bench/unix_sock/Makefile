all:client server

client:client.c
	cc -o client client.c

server:server.c
	cc -o server server.c

bench:all
	@pkill server
	./server 1000&
	@sleep 1
	./client

clean:
	rm -f client server demo.sock
