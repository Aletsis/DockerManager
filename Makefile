.PHONY: all build dev clean run

WAILS := $(HOME)/go/bin/wails

all: build

# Compilar la aplicación completa para producción
build:
	$(WAILS) build -tags webkit2_41

# Ejecutar en modo desarrollo con recarga en vivo (hot-reload)
dev:
	$(WAILS) dev -tags webkit2_41

# Ejecutar el binario ya compilado
run:
	./build/bin/dockermanager

# Limpiar binarios generados
clean:
	rm -rf build/bin
