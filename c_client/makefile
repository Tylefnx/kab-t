CC = gcc
CFLAGS = -Wall -g -fsanitize=address -I.
LIBS = -lncurses -lwebsockets -lparson

TARGET = program
SRCS = main.c
OBJS = $(SRCS:.c=.o)

all: $(TARGET)

$(TARGET): $(OBJS)
	$(CC) $(CFLAGS) -o $@ $^ $(LIBS)

%.o: %.c
	$(CC) $(CFLAGS) -c $< -o $@

clean:
	rm -f $(OBJS) $(TARGET)
