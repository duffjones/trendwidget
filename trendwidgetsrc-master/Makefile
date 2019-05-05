CFLAGS = -v
INCS =
SOURCES = $(MAIN) $(TEST)
TARGET	=	main
TTARGET = maintest
MAIN = $(TARGET).go
TEST = $(TTARGET).go
SRCFOLDER = source/
LIBS =
CC = go

all: format test run clean

format:
	$(CC) fmt ./...

test:
	$(CC) build $(CFLAGS) $(SRCFOLDER)$(TEST)
	$(CC) run $(CFLAGS) $(SRCFOLDER)$(TEST)
run:
	$(CC) build $(SRCFOLDER)$(MAIN)
	$(CC) run $(SRCFOLDER)$(MAIN)

clean:
	rm -f $(TARGET)
	rm -f $(TTARGET)
