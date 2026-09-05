# 🐳 DockerManager

Una herramienta de escritorio moderna, minimalista y de alto rendimiento para la gestión y monitoreo en tiempo real de contenedores Docker en Linux.

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

## 🚀 Cómo Ejecutar

### 1. Ejecución Inmediata
```bash
./build/bin/dockermanager
```
O usando el Makefile:
```bash
make run
```

### 2. Modo Desarrollo (Hot-Reload)
```bash
make dev
```

### 3. Recompilar para Producción
```bash
make build
```
