# 🚀 REST & GraphQL Task Server

[![Go Version](https://img.shields.io/badge/Go-1.24.0-00ADD8?style=flat&logo=go)](https://golang.org/)
[![GraphQL](https://img.shields.io/badge/GraphQL-E10098?style=flat&logo=graphql&logoColor=white)](https://graphql.org/)
[![Swagger](https://img.shields.io/badge/Swagger-85EA2D?style=flat&logo=swagger&logoColor=black)](https://swagger.io/)

> API REST y GraphQL para la gestión de tareas con soporte HTTPS, documentación Swagger y GraphQL Playground.

---

## 📋 Tabla de Contenidos

- [Características](#-características)
- [Requisitos](#-requisitos)
- [Instalación](#-instalación)
- [Configuración](#-configuración)
- [Uso](#-uso)
- [Estructura del Proyecto](#-estructura-del-proyecto)
- [API REST](#-api-rest)
- [GraphQL](#-graphql)
- [Comandos Útiles](#-comandos-útiles)
- [Ejemplos](#-ejemplos)

---

## ✨ Características

- 🌐 **Dual API**: REST y GraphQL en el mismo servidor
- 🔒 **HTTPS**: Servidor seguro con certificados autofirmados
- 📚 **Documentación Swagger**: Interfaz interactiva para probar endpoints
- 🎮 **GraphQL Playground**: Editor interactivo para queries y mutations
- 🔐 **Basic Auth**: Autenticación opcional para endpoints privados
- 📝 **Gestión completa de tareas**: CRUD con tags, fechas y adjuntos
- 🔄 **Store compartido**: REST y GraphQL usan la misma base de datos en memoria
- ⚡ **Thread-safe**: Store con sync.RWMutex para concurrencia

---

## 📦 Requisitos

- **Go 1.24.0** o superior
- **mkcert** (para generar certificados HTTPS)
- **swag** (para documentación Swagger)
- **gqlgen** (para código GraphQL)

---

## 🔧 Instalación

### 1. Clonar el repositorio

```bash
git clone https://github.com/tu-usuario/restServer.git
cd restServer
```

### 2. Instalar dependencias

```bash
go mod download
```

### 3. Instalar herramientas de desarrollo

```powershell
# Instalar mkcert (con Chocolatey)
choco install mkcert

# O con Scoop
scoop install mkcert

# Instalar swag para Swagger
go install github.com/swaggo/swag/cmd/swag@latest

# Instalar gqlgen para GraphQL
go install github.com/99designs/gqlgen@latest
```

---

## ⚙️ Configuración

### Generar Certificados HTTPS

```bash
# Instalar CA local de mkcert
mkcert -install

# Generar certificados para localhost
mkcert localhost
```

Esto genera:
- `localhost.pem` (certificado público)
- `localhost-key.pem` (clave privada)

### Generar Documentación Swagger

```bash
swag init
```

Esto genera el directorio `docs/` con:
- `docs.go`
- `swagger.json`
- `swagger.yaml`

### Regenerar Código GraphQL (si modificas el schema)

```bash
go run github.com/99designs/gqlgen generate
```

---

## 🚀 Uso

### Iniciar el Servidor

```powershell
go run main.go
```

**Salida esperada:**
```
TaskServer store creado: 0x...
GraphQL Resolver store configurado: 0x...
Listening in https://localhost:8443
```

### Acceder a las Interfaces

| Interfaz | URL | Descripción |
|----------|-----|-------------|
| **Swagger UI** | https://localhost:8443/docs/ | Documentación REST interactiva |
| **GraphQL Playground** | https://localhost:8443/playground | Editor GraphQL interactivo |
| **Endpoint GraphQL** | https://localhost:8443/graphql | API GraphQL |
| **API REST** | https://localhost:8443/task/ | Endpoints REST |

⚠️ **Importante**: Al acceder por primera vez, acepta el certificado autofirmado en tu navegador:
1. Ve a `https://localhost:8443/graphql`
2. Click en "Avanzado" → "Continuar a localhost"
3. Luego accede a las demás URLs

---

## 📁 Estructura del Proyecto

```
restServer/
│
├── 📄 main.go                    # Servidor principal (EJECUTAR ESTE)
├── 📄 server.go                  # Servidor GraphQL standalone (generado por gqlgen)
├── 📄 gen_cert.go               # Generador de certificados
├── 📄 go.mod                    # Dependencias
├── 📄 gqlgen.yml                # Configuración de gqlgen
├── 📄 task.json                 # Ejemplo de tarea
├── 🔐 localhost.pem             # Certificado HTTPS
├── 🔐 localhost-key.pem         # Clave privada
│
├── 📂 docs/                     # Documentación Swagger (auto-generada)
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
│
├── 📂 graph/                    # GraphQL
│   ├── generated.go            # Código generado por gqlgen
│   ├── resolver.go             # Inyección de dependencias
│   ├── schema.graphqls         # Schema GraphQL
│   ├── schema.resolvers.go     # Implementación de resolvers
│   └── model/
│       └── models_gen.go       # Modelos generados
│
├── 📂 internal/                 # Código interno
│   └── middleware.go           # Middlewares (Auth, Logging)
│
├── 📂 server/                   # Lógica REST
│   └── serverRest.go           # Handlers REST
│
└── 📂 taskstore/                # Almacenamiento
    └── taskstore.go            # Store en memoria (CRUD)
```

---

## 🌐 API REST

### Endpoints Disponibles

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| `POST` | `/task/` | Crear una tarea |
| `GET` | `/task/` | Obtener todas las tareas |
| `GET` | `/task/{id}/` | Obtener tarea por ID |
| `GET` | `/tag/{tag}/` | Obtener tareas por tag |
| `GET` | `/due/{year}/{month}/{day}/` | Obtener tareas por fecha |
| `DELETE` | `/task/{id}/` | Eliminar tarea por ID |
| `DELETE` | `/task/` | Eliminar todas las tareas |

### Ejemplos con curl (PowerShell)

#### Crear una tarea

```powershell
curl.exe -k -X POST https://localhost:8443/task/ `
  -H "Content-Type: application/json" `
  -d '{
    "text": "Comprar pan",
    "tags": ["personal", "urgente"],
    "due": "2025-12-25T23:59:59Z",
    "attachments": []
  }'
```

**Respuesta:**
```json
{"id":"1"}
```

#### Obtener todas las tareas

```powershell
curl.exe -k https://localhost:8443/task/
```

**Respuesta:**
```json
[
  {
    "Id": "1",
    "Text": "Comprar pan",
    "Tags": ["personal", "urgente"],
    "Due": "2025-12-25T23:59:59Z",
    "Attachments": []
  }
]
```

#### Obtener tareas por tag

```powershell
curl.exe -k https://localhost:8443/tag/personal/
```

#### Obtener tareas por fecha

```powershell
curl.exe -k https://localhost:8443/due/2025/12/25/
```

#### Eliminar una tarea

```powershell
curl.exe -k -X DELETE https://localhost:8443/task/1/
```

#### Eliminar todas las tareas

```powershell
curl.exe -k -X DELETE https://localhost:8443/task/
```

### Modelo de Datos (JSON)

#### Task Request
```json
{
  "text": "string (requerido)",
  "tags": ["string", "string"],
  "due": "2025-12-25T23:59:59Z (formato RFC3339, requerido)",
  "attachments": [
    {
      "Name": "string",
      "Date": "2025-12-23T10:00:00Z",
      "Contents": "string"
    }
  ]
}
```

#### Task Response
```json
{
  "Id": "string (UUID)",
  "Text": "string",
  "Tags": ["string"],
  "Due": "2025-12-25T23:59:59Z",
  "Attachments": [
    {
      "Name": "string",
      "Date": "2025-12-23T10:00:00Z",
      "Contents": "string"
    }
  ]
}
```

---

## 🎮 GraphQL

### Schema

```graphql
type Query {
    getAllTasks: [Task]
    getTask(id: ID!): Task
    getTasksByTag(tag: String!): [Task]
    getTasksByDue(due: Time!): [Task]
}

type Mutation {
    createTask(input: NewTask!): Task!
    deleteTask(id: ID!): Boolean
    deleteAllTasks: Boolean
}

type Task {
    Id: ID!
    Text: String!
    Tags: [String!]
    Due: Time!
    Attachments: [Attachment!]
}

input NewTask {
    Text: String!
    Tags: [String!]
    Due: Time!
    Attachments: [NewAttachment!]
}
```

### Ejemplos de Queries

#### Obtener todas las tareas

```graphql
query {
  getAllTasks {
    Id
    Text
    Tags
    Due
    Attachments {
      Name
      Date
      Contents
    }
  }
}
```

#### Obtener tarea por ID

```graphql
query {
  getTask(id: "1") {
    Id
    Text
    Tags
    Due
  }
}
```

#### Obtener tareas por tag

```graphql
query {
  getTasksByTag(tag: "personal") {
    Id
    Text
    Tags
  }
}
```

#### Obtener tareas por fecha

```graphql
query {
  getTasksByDue(due: "2025-12-25T23:59:59Z") {
    Id
    Text
    Due
  }
}
```

### Ejemplos de Mutations

#### Crear una tarea

```graphql
mutation {
  createTask(input: {
    Text: "Estudiar GraphQL"
    Tags: ["estudio", "importante"]
    Due: "2025-12-30T20:00:00Z"
    Attachments: []
  }) {
    Id
    Text
    Tags
    Due
  }
}
```

#### Crear tarea con adjuntos

```graphql
mutation {
  createTask(input: {
    Text: "Revisar documentos"
    Tags: ["trabajo", "documentos"]
    Due: "2025-12-28T18:00:00Z"
    Attachments: [
      {
        Name: "reporte.pdf"
        Date: "2025-12-23T10:00:00Z"
        Contents: "base64encodedcontent..."
      }
    ]
  }) {
    Id
    Text
    Attachments {
      Name
      Date
    }
  }
}
```

#### Eliminar una tarea

```graphql
mutation {
  deleteTask(id: "1")
}
```

#### Eliminar todas las tareas

```graphql
mutation {
  deleteAllTasks
}
```

### Usar GraphQL con curl

```powershell
curl.exe -k -X POST https://localhost:8443/graphql `
  -H "Content-Type: application/json" `
  -d '{
    "query": "query { getAllTasks { Id Text Tags Due } }"
  }'
```

Con variables:

```powershell
curl.exe -k -X POST https://localhost:8443/graphql `
  -H "Content-Type: application/json" `
  -d '{
    "query": "mutation($input: NewTask!) { createTask(input: $input) { Id Text } }",
    "variables": {
      "input": {
        "Text": "Tarea desde curl",
        "Tags": ["curl"],
        "Due": "2025-12-25T23:59:59Z",
        "Attachments": []
      }
    }
  }'
```

---

## 🛠️ Comandos Útiles

### Desarrollo

```bash
# Ejecutar el servidor
go run main.go

# Compilar binario
go build -o taskserver.exe main.go

# Ejecutar binario
./taskserver.exe

# Ejecutar tests
go test ./...

# Test con cobertura
go test -cover ./...

# Ver dependencias
go mod graph
```

### Certificados HTTPS

```bash
# Instalar mkcert
choco install mkcert

# Instalar CA local
mkcert -install

# Generar certificados
mkcert localhost

# Para múltiples dominios
mkcert localhost 127.0.0.1 ::1

# Ver información del certificado
openssl x509 -in localhost.pem -text -noout
```

### Swagger (Documentación REST)

```bash
# Generar documentación
swag init

# Formatear anotaciones
swag fmt

# Validar anotaciones
swag validate
```

### GraphQL

```bash
# Generar código GraphQL
go run github.com/99designs/gqlgen generate

# Validar schema
go run github.com/99designs/gqlgen validate

# Inicializar proyecto GraphQL (ya hecho)
go run github.com/99designs/gqlgen init
```

### Gestión de Dependencias

```bash
# Descargar dependencias
go mod download

# Limpiar dependencias no usadas
go mod tidy

# Verificar dependencias
go mod verify

# Actualizar dependencia específica
go get -u github.com/99designs/gqlgen

# Actualizar todas las dependencias
go get -u ./...
```

---

## 💡 Ejemplos Completos

### Ejemplo con PowerShell

```powershell
# Crear múltiples tareas
$tareas = @(
    @{text="Comprar pan"; tags=@("personal"); due="2025-12-25T10:00:00Z"; attachments=@()},
    @{text="Reunión equipo"; tags=@("trabajo"); due="2025-12-26T15:00:00Z"; attachments=@()},
    @{text="Estudiar Go"; tags=@("estudio"); due="2025-12-27T20:00:00Z"; attachments=@()}
)

foreach ($tarea in $tareas) {
    $body = $tarea | ConvertTo-Json
    Invoke-RestMethod -Uri "https://localhost:8443/task/" `
      -Method Post `
      -Body $body `
      -ContentType "application/json" `
      -SkipCertificateCheck
    
    Write-Host "✅ Tarea creada: $($tarea.text)"
}

# Obtener todas las tareas
$allTasks = Invoke-RestMethod -Uri "https://localhost:8443/task/" `
  -Method Get `
  -SkipCertificateCheck

Write-Host "`nTotal de tareas: $($allTasks.Count)"
$allTasks | Format-Table Id, Text, Tags
```

### Ejemplo con Python

```python
import requests
import json
from urllib3.exceptions import InsecureRequestWarning

# Deshabilitar warnings SSL (solo desarrollo)
requests.packages.urllib3.disable_warnings(InsecureRequestWarning)

# Crear tarea
task = {
    "text": "Tarea desde Python",
    "tags": ["python", "automatización"],
    "due": "2025-12-25T23:59:59Z",
    "attachments": []
}

response = requests.post(
    'https://localhost:8443/task/',
    json=task,
    verify=False
)
print(f"✅ Tarea creada: {response.json()}")

# Obtener todas las tareas
response = requests.get('https://localhost:8443/task/', verify=False)
tasks = response.json()

print(f"\nTotal de tareas: {len(tasks)}")
for task in tasks:
    print(f"- [{task['Id']}] {task['Text']}")
```

### Ejemplo con JavaScript/Node.js

```javascript
const https = require('https');

// Ignorar certificados autofirmados (solo desarrollo)
const agent = new https.Agent({
  rejectUnauthorized: false
});

// Función auxiliar para hacer peticiones
function makeRequest(options, data = null) {
  return new Promise((resolve, reject) => {
    const req = https.request({ ...options, agent }, (res) => {
      let body = '';
      res.on('data', chunk => body += chunk);
      res.on('end', () => resolve(JSON.parse(body)));
    });
    
    req.on('error', reject);
    if (data) req.write(JSON.stringify(data));
    req.end();
  });
}

// Crear tarea
const newTask = {
  text: "Tarea desde Node.js",
  tags: ["javascript", "node"],
  due: "2025-12-25T23:59:59Z",
  attachments: []
};

makeRequest({
  hostname: 'localhost',
  port: 8443,
  path: '/task/',
  method: 'POST',
  headers: { 'Content-Type': 'application/json' }
}, newTask)
.then(result => {
  console.log('✅ Tarea creada:', result);
  return makeRequest({
    hostname: 'localhost',
    port: 8443,
    path: '/task/',
    method: 'GET'
  });
})
.then(tasks => {
  console.log(`\nTotal de tareas: ${tasks.length}`);
  tasks.forEach(task => {
    console.log(`- [${task.Id}] ${task.Text}`);
  });
})
.catch(err => console.error('❌ Error:', err));
```

### Ejemplo con Go

```go
package main

import (
    "bytes"
    "crypto/tls"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
)

type Task struct {
    Text        string   `json:"text"`
    Tags        []string `json:"tags"`
    Due         string   `json:"due"`
    Attachments []any    `json:"attachments"`
}

type TaskResponse struct {
    ID          string   `json:"Id"`
    Text        string   `json:"Text"`
    Tags        []string `json:"Tags"`
    Due         string   `json:"Due"`
    Attachments []any    `json:"Attachments"`
}

func main() {
    // Cliente que ignora certificados (solo desarrollo)
    tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }
    client := &http.Client{Transport: tr}

    // Crear tarea
    task := Task{
        Text:        "Tarea desde Go",
        Tags:        []string{"golang", "código"},
        Due:         "2025-12-25T23:59:59Z",
        Attachments: []any{},
    }

    body, _ := json.Marshal(task)
    resp, err := client.Post(
        "https://localhost:8443/task/",
        "application/json",
        bytes.NewBuffer(body),
    )
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    var createResp map[string]string
    json.NewDecoder(resp.Body).Decode(&createResp)
    fmt.Printf("✅ Tarea creada: %s\n", createResp["id"])

    // Obtener todas las tareas
    resp, _ = client.Get("https://localhost:8443/task/")
    defer resp.Body.Close()

    data, _ := io.ReadAll(resp.Body)
    var tasks []TaskResponse
    json.Unmarshal(data, &tasks)

    fmt.Printf("\nTotal de tareas: %d\n", len(tasks))
    for _, t := range tasks {
        fmt.Printf("- [%s] %s\n", t.ID, t.Text)
    }
}
```

---

## 🔐 Autenticación Basic Auth (Opcional)

Para habilitar autenticación en endpoints específicos, descomenta las líneas en `main.go`:

```go
// Endpoints privados (con autenticación básica)
mux.Handle("GET /task/", internal.BasicAuth("admin", "1234",
    http.HandlerFunc(taskServer.GetAllTasksHandler)))
mux.Handle("DELETE /task/", internal.BasicAuth("admin", "1234",
    http.HandlerFunc(taskServer.DeleteAllTasksHandler)))
mux.Handle("DELETE /task/{id}/", internal.BasicAuth("admin", "1234",
    http.HandlerFunc(taskServer.DeleteTaskHandler)))
```

### Usar con Basic Auth

**Con curl:**
```powershell
curl.exe -k -u admin:1234 https://localhost:8443/task/
```

**Con PowerShell:**
```powershell
$credential = Get-Credential
Invoke-RestMethod -Uri "https://localhost:8443/task/" `
  -Authentication Basic `
  -Credential $credential `
  -SkipCertificateCheck
```

**Con GraphQL Playground:**

Agrega el header en el panel inferior:
```json
{
  "Authorization": "Basic YWRtaW46MTIzNA=="
}
```

(Base64 de `admin:1234`)

---

## 🐛 Troubleshooting

### Error: "transport not supported"

**Solución:** Verifica que en `main.go` estén configurados los transportes:

```go
graphqlServer.AddTransport(transport.Options{})
graphqlServer.AddTransport(transport.GET{})
graphqlServer.AddTransport(transport.POST{})
```

### Error: "Server cannot be reached" en Playground

**Solución:** Acepta el certificado autofirmado:
1. Ve a `https://localhost:8443/graphql` primero
2. Acepta la advertencia de seguridad
3. Luego ve al playground

### Error: "parsing time ... cannot parse 'T'"

**Solución:** Usa formato RFC3339 con "T":
```
✅ Correcto: "2025-12-25T23:59:59Z"
❌ Incorrecto: "2025-12-25 23:59:59"
```

### Puerto 8443 en uso

**Solución:**
```powershell
# Ver qué proceso usa el puerto
netstat -ano | findstr :8443

# Detener el proceso
taskkill /PID <PID> /F
```

---

## 📚 Recursos

- [Documentación Go](https://golang.org/doc/)
- [GraphQL Spec](https://spec.graphql.org/)
- [gqlgen Docs](https://gqlgen.com/)
- [Swagger/OpenAPI](https://swagger.io/specification/)
- [mkcert GitHub](https://github.com/FiloSottile/mkcert)

---

## 📝 Notas Importantes

⚠️ **Este servidor usa almacenamiento en memoria. Los datos se pierden al reiniciar.**

Para producción considera:
- Base de datos persistente (PostgreSQL, MongoDB)
- Certificados SSL válidos (Let's Encrypt)
- Autenticación JWT
- Rate limiting
- CORS configurado
- Variables de entorno para configuración

---

<div align="center">

### ⭐ Si este proyecto te fue útil, considera darle una estrella!

**Hecho con ❤️ y Go**

</div>
