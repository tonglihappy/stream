CC = g++
AR = ar

CFLAGS = -O0 -DDEBUG -Wall -c -g

SRCS = $(wildcard ./src/*.cpp)
OBJS = $(patsubst %.cpp, %.o, $(SRCS))
CSRCS = $(wildcard ./src/*.c)

HEADER_PATH = -I./include/
LIB_PATH = -L./lib/

LIBS = -lavformat -lavcodec -lswscale -lavutil -lavfilter -lswresample -lavdevice -lpostproc -lx264 -ldl -lz -lm -lrt -lpthread

$(OBJS):
	$(CC) $(CFLAGS) $(SRCS) $(HEADER_PATH)
	$(CC) $(CFLAGS) $(CSRCS) $(HEADER_PATH)
	$(AR) cr libffmpeg.a *.o
	mv libffmpeg.a ./lib
	rm -rf *.o
#$(TARGET) : $(OBJS)
#	$(AR) cr libffmpeg.a *.o
#	rm -rf *.o

test:
	g++ -O0 -DDEBUG -Wall -g -o test test.cpp $(HEADER_PATH) $(LIB_PATH) -lffmpeg -ldl -lm -lrt -lpthread -lstdc++ -lavformat -lavcodec -lswscale -lavutil -lavfilter -lswresample -lavdevice -lpostproc -lx264 

clean:
	rm -rf ./*.a
	echo "hello"
