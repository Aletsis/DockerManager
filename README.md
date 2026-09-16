# 🐳 DockerManager

Una herramienta de escritorio moderna, minimalista y de alto rendimiento para la gestión y monitoreo en tiempo real de contenedores Docker en **Linux**, **macOS** y **Windows**.

Construida con **Go (Golang)** en el núcleo, **Wails v2** para el enlace con la ventana nativa y **Svelte 5 + TypeScript + Tailwind CSS** para una interfaz visual ultraligera con selector de **Modo Claro / Modo Oscuro**.

---

## ✨ Características Principales

* **Consumo Mínimo de Recursos (Svelte 5 + Go)**:
  * Sin sobrecarga de Virtual DOM: el frontend se actualiza de manera reactiva y quirúrgica.
  * Tamaño del bundle JavaScript comprimido de apenas **~29 KB**.
  * Binario único autocontenido de apenas **~13 MB** y uso de memoria RAM inferior a **~35-40 MB**.
* **Control Completo del Ciclo de Vida de Contenedores**:
  * Iniciar (`Start`), Detener (`Stop`), Reiniciar (`Restart`), Pausar (`Pause`), Reanudar (`Unpause`) y Eliminar (`Remove`) con diálogo de confirmación.
* **Monitoreo en Vivo**:
  * Métricas en tiempo real de uso de **CPU (%)** y **Memoria RAM** por contenedor.
  * Medidores de red (**Download Rx / Upload Tx**), **I/O de disco** y cantidad de procesos (**PIDs**).
* **Visor de Logs en Tiempo Real**:
  * Transmisión en vivo de registros (`stdout` y `stderr`) del contenedor seleccionado.
  * Selector de líneas (`tail: 50, 100, 200, 500`), filtro de búsqueda y auto-scroll inteligente.
* **Diseño Minimalista Claro / Oscuro**:
  * Selector de tema sol/luna (☀️ / 🌙) con persistencia local.
  * Lista de contenedores con scroll suave y barra minimalista.
  * Selector de frecuencia de refresco (1s, 2s, 5s, 10s o Pausa) adaptado a cada tema.
  * Filtro de búsqueda instantánea por nombre, imagen o ID.

---

## 📋 Prerrequisitos del Sistema

### 1. Comunes a todas las plataformas

* **Docker**: Debe estar instalado y corriendo.
  * Verificar en terminal con: `docker ps`
* **Go (Golang)**: Versión `1.21` o superior (se recomienda `1.23+`).
  * Verificar con: `go version`
* **Node.js**: Versión `18` o superior con `npm`.
  * Verificar con: `node -v` y `npm -v`
* **Wails CLI (v2)**:
  * Instalar ejecutando:
    ```bash
    go install github.com/wailsapp/wails/v2/cmd/wails@latest
    ```
  * Asegúrate de que el directorio de binarios de Go (`$GOPATH/bin` o `~/go/bin`) esté incluido en tu variable de entorno `PATH`.
  * Puedes verificar el estado de dependencias del sistema en cualquier momento con:
    ```bash
    wails doctor
    ```

---

### 2. Específicos por Sistema Operativo

#### 🍎 macOS
* **Xcode Command Line Tools**:
  ```bash
  xcode-select --install
  ```
* **Docker**: Docker Desktop u OrbStack activo.

#### 🐧 Linux
* **Librerías de desarrollo C / GTK / WebKit**:
  * **Ubuntu / Debian**:
    ```bash
    sudo apt update
    sudo apt install build-essential libgtk-3-dev libwebkit2gtk-4.1-dev
    ```
    *(Si tu distribución utiliza WebKit 4.0, instala `libwebkit2gtk-4.0-dev`).*
  * **Fedora**:
    ```bash
    sudo dnf install gcc-c++ gtk3-devel webkit2gtk4.1-devel
    ```
  * **Arch Linux**:
    ```bash
    sudo pacman -S base-devel gtk3 webkit2gtk-4.1
    ```
* **Permisos de Docker**:
  * Asegúrate de que tu usuario pertenezca al grupo `docker` para no requerir `sudo`:
    ```bash
    sudo usermod -aG docker $USER
    newgrp docker
    ```

#### 🪟 Windows
* **Compilador C/C++ (CGO)**:
  * Wails requiere un compilador C en Windows. Se recomienda **MinGW-w64**.
  * Vía Chocolatey:
    ```powershell
    choco install mingw
    ```
  * Vía Scoop:
    ```powershell
    scoop install mingw
    ```
  * O vía MSYS2.
* **WebView2 Runtime**: Viene preinstalado en Windows 10 y 11. Si no lo tienes, descárgalo del sitio oficial de Microsoft.
* **Docker**: Docker Desktop en ejecución (con WSL 2 habilitado).

---

## 🚀 Modo Desarrollo (Hot-Reload)

El modo desarrollo inicia el servidor de Vite con recarga en vivo en el frontend y compila el backend en Go automáticamente al detectar cambios. Además, instala las dependencias de `npm` automáticamente la primera vez.

### 🍎 macOS
```bash
wails dev
```

### 🐧 Linux
Usando el Makefile (configurado con las etiquetas de WebKit correspondientes):
```bash
make dev
```
O directamente con el CLI:
```bash
wails dev -tags webkit2_41
```

### 🪟 Windows
Desde PowerShell o Command Prompt:
```powershell
wails dev
```

> 💡 **Tip**: Una vez levantado en modo desarrollo, también puedes abrir `http://localhost:34115` en tu navegador para depurar con las herramientas de inspección (DevTools).

---

## 📦 Compilar y Ejecutar en Producción

### 🍎 macOS

1. **Compilar la aplicación:**
   ```bash
   wails build
   ```
   *(Para compilar para una arquitectura específica, puedes usar `-platform darwin/arm64` para Apple Silicon o `-platform darwin/amd64` para Intel).*

2. **Ejecutar:**
   ```bash
   open build/bin/dockermanager.app
   ```
   O directamente el ejecutable interior:
   ```bash
   ./build/bin/dockermanager.app/Contents/MacOS/dockermanager
   ```

---

### 🐧 Linux

1. **Compilar la aplicación:**
   ```bash
   make build
   # O directamente: wails build -tags webkit2_41
   ```

2. **Ejecutar:**
   ```bash
   make run
   # O directamente: ./build/bin/dockermanager
   ```

---

### 🪟 Windows

1. **Compilar la aplicación:**
   ```powershell
   wails build
   ```

2. **Ejecutar:**
   ```powershell
   .\build\bin\dockermanager.exe
   ```

---

## 🩺 Diagnóstico y Verificación

Si encuentras algún problema para compilar o ejecutar, corre la herramienta de autodiagnóstico de Wails:

```bash
wails doctor
```

Este comando analizará tu entorno y te indicará con precisión si falta alguna librería, compilador o configuración en tu sistema operativo.
