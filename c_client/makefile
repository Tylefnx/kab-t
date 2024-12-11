# Compiler
CC = gcc

# Directories
SRCDIR = src
OBJDIR = obj
BINDIR = bin

# Executable
EXEC = $(BINDIR)/server

# Source and Object files
SOURCES = $(wildcard $(SRCDIR)/*.c)
OBJECTS = $(patsubst $(SRCDIR)/%.c, $(OBJDIR)/%.o, $(SOURCES))

# Include directories
INCLUDES = -Iinclude

# Libraries
LIBS = -lncurses -lwebsockets -lcurl -lparson

# Compiler flags
CFLAGS = -Wall -Wextra $(INCLUDES)
LDFLAGS = $(LIBS)

# Default target
all: $(EXEC)

# Link
$(EXEC): $(OBJECTS)
	@mkdir -p $(BINDIR)
	$(CC) $^ -o $@ $(LDFLAGS)

# Compile
$(OBJDIR)/%.o: $(SRCDIR)/%.c
	@mkdir -p $(OBJDIR)
	$(CC) $(CFLAGS) -c $< -o $@

# Clean
clean:
	rm -rf $(OBJDIR) $(BINDIR)

# Run
run: $(EXEC)
	$(EXEC)

.PHONY: all clean run
