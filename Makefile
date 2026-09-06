.PHONY: all build dev clean run version bump-patch bump-minor bump-major commit install-hooks

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

# Ver versión actual de la aplicación
version:
	@./scripts/bump-version.sh get

# Incrementar manualmente la versión
bump-patch:
	@./scripts/bump-version.sh patch

bump-minor:
	@./scripts/bump-version.sh minor

bump-major:
	@./scripts/bump-version.sh major

# Realizar un commit convencional con auto-incremento de versión
# Uso: make commit m="feat: nueva característica"
commit:
	@./scripts/commit.sh "$(m)"

# Instalar Git Hook automático para auto-incremento con git commit estándar
install-hooks:
	git config core.hooksPath .githooks
	@echo "✓ Git Hook instalado. Ahora cada 'git commit -m ...' incrementará la versión automáticamente."

