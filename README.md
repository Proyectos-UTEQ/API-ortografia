# Proyecto Ortografia xD

Este proyecto es la parte del back de la aplicación para mejorar la ortografia.

## Requisitos

- Docker
- Visual Studio Code

## Configuración del Entorno de Desarrollo

### 1. Clona el Repositorio

```bash
$ git clone git@github.com:Proyectos-UTEQ/api-ortografia.git
$ cd api-ortografia
```

### 2. Abre el Proyecto en Visual Studio Code

Si tienes Visual Studio Code instalado, puedes abrir el proyecto con los contenedores de desarrollo proporcionados.

- Instala la extensión "Dev Containers" en VS Code.
- Abre el proyecto en VS Code.
- Selecciona "Reopen in Container" cuando se te pregunte sobre la apertura en un contenedor.

### 3. Ejecuta el Proyecto

Usa el comando `make run` para ejecutar el proyecto.
```bash
make run
```

Este comando realizará las siguientes acciones:

- Descargará las dependencias del proyecto.
- Compilará el código.
- Ejecutará la aplicación.

## Comandos Útiles
- `$ make run` : Ejecuta la aplicación.
- `$ make build` : Compila la aplicación.
- `$ make seed` : Pobla la base de datos.
