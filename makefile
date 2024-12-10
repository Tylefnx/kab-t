CC = gcc
CFLAGS = -Wall -g
LIBS = -lmicrohttpd -lsqlite3
SRC = src/main.c src/router.c src/question.c src/leaderboard.c src/db.c

kahoot-backend: $(SRC)
	$(CC) $(CFLAGS) -o kahoot-backend $(SRC) $(LIBS)

clean:
	rm -f kahoot-backend
