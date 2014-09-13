all:client server

client:client.c
	cc -o client client.c

server:server.c
	cc -o server server.c

run:all
	./server 100&
	./client

clean:
	rm -f client server demo.sock
