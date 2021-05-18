PREFIX =				/usr/local/
DESTDIR =

.PHONY: all clean

all: saturn

saturn:
				go build -o saturn

clean:
				rm -f saturn

install: saturn
				mkdir -p ${DESTDIR}${PREFIX}/bin/
				install -m 0555 saturn ${DESTDIR}${PREFIX}/bin
