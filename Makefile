PREFIX =				/usr/local/
DESTDIR =

.PHONY: all clean

all: saturn

saturn:
				go build

clean:
				rm -f saturn

install: saturn
				mkdir -p ${DESTDIR}${PREFIX}/bin/
				install -m 0555 saturn ${DESTDIR}${PREFIX}/bin
