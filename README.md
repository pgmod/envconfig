# envconfig

Библиотека для работы с переменными окружения в Go. Предоставляет удобные функции для загрузки конфигурации из переменных окружения и `.env` файлов.

## Установка

```bash
go get github.com/pgmod/envconfig
```

## Основные возможности

- Загрузка переменных окружения из `.env` файлов
- Автоматическая загрузка конфигурации в структуры с использованием тегов
- Поддержка значений по умолчанию
- Поддержка типов: `string`, `bool`, `int`, `int64`, `[]int`, `[]int64`, массивы `int`
- Простые функции для получения значений с дефолтами
- Функции для получения массивов чисел: `GetIntSlice()`, `GetInt64Slice()`

## Быстрый старт

### Загрузка .env файла

```go
package main

import (
    "log"
    "github.com/pgmod/envconfig"
)

func main() {
    if err := envconfig.Load(); err != nil {
        log.Printf("Ошибка загрузки .env файла: %v", err)
    }
    // Теперь переменные из .env доступны через os.Getenv()
}
```

### Загрузка конфигурации в структуру

```go
package main

import (
    "fmt"
    "log"
    "github.com/pgmod/envconfig"
)

type Config struct {
    Host     string `env:"HOST" default:"localhost"`
    Port     int    `env:"PORT" default:"8080"`
    Debug    bool   `env:"DEBUG" default:"false"`
    Ports    []int  `env:"PORTS" default:"8080,8081,8082"`
    Fixed    [3]int `env:"FIXED" default:"1,2,3"`
}

func main() {
    var cfg Config
    
    if err := envconfig.LoadStruct(&cfg); err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Host: %s\n", cfg.Host)
    fmt.Printf("Port: %d\n", cfg.Port)
    fmt.Printf("Debug: %v\n", cfg.Debug)
    fmt.Printf("Ports: %v\n", cfg.Ports)
    fmt.Printf("Fixed: %v\n", cfg.Fixed)
}
```

## API Reference

### Load()

Загружает переменные окружения из `.env` файла. По умолчанию ищет файл `.env` в текущей директории. Можно указать другой файл через переменную окружения `ENV_FILE`.

```go
err := envconfig.Load()
```

**Пример:**

```bash
# Установить путь к файлу через переменную окружения
export ENV_FILE=config.env
```

```go
envconfig.Load() // Загрузит config.env
```

### LoadStruct(cfg any) error

Загружает конфигурацию из переменных окружения в структуру. Использует теги `env` для указания имени переменной окружения и `default` для значения по умолчанию.

**Поддерживаемые типы:**
- `string`
- `bool`
- `int`, `int64`
- `[]int`, `[]int64` (слайсы)
- `[N]int`, `[N]int64` (массивы фиксированного размера)

**Теги:**
- `env:"VAR_NAME"` - имя переменной окружения
- `default:"value"` - значение по умолчанию (используется, если переменная не установлена)

**Пример:**

```go
type Config struct {
    // Строковое значение
    Host string `env:"HOST" default:"localhost"`
    
    // Числовое значение
    Port int `env:"PORT" default:"8080"`
    
    // Булево значение
    Enabled bool `env:"ENABLED" default:"true"`
    
    // Слайс int (значения разделяются запятыми)
    Ports []int `env:"PORTS" default:"8080,8081"`
    
    // Массив фиксированного размера
    Fixed [3]int `env:"FIXED" default:"1,2,3"`
}

var cfg Config
if err := envconfig.LoadStruct(&cfg); err != nil {
    log.Fatal(err)
}
```

**Переменные окружения:**

```bash
export HOST=example.com
export PORT=9000
export ENABLED=true
export PORTS=8080,8081,8082
export FIXED=10,20,30
```

**Примечания:**
- Поля без тега `env` игнорируются
- Если переменная окружения не установлена, используется значение из `default`
- Для массивов количество значений должно совпадать с размером массива
- Пробелы вокруг значений в массивах автоматически удаляются

### Get(key, defaultValue string) string

Получает строковое значение переменной окружения с дефолтным значением.

```go
host := envconfig.Get("HOST", "localhost")
```

### GetBool(key string, defaultValue bool) bool

Получает булево значение переменной окружения с дефолтным значением.

```go
debug := envconfig.GetBool("DEBUG", false)
```

**Поддерживаемые значения:** `true`, `false`, `1`, `0`, `t`, `f`, `T`, `F`, `TRUE`, `FALSE`, `True`, `False`

### GetInt(key string, defaultValue int) int

Получает целочисленное значение переменной окружения с дефолтным значением.

```go
port := envconfig.GetInt("PORT", 8080)
```

### GetInt64(key string, defaultValue int64) int64

Получает 64-битное целочисленное значение переменной окружения с дефолтным значением.

```go
maxSize := envconfig.GetInt64("MAX_SIZE", 1024)
```

### GetIntSlice(key string, defaultValue []int) []int

Получает слайс целых чисел из переменной окружения с дефолтным значением. Значения должны быть разделены запятыми. Пробелы вокруг значений автоматически удаляются. Возвращает значение по умолчанию, если переменная окружения не установлена, пуста или содержит невалидные значения.

```go
ports := envconfig.GetIntSlice("PORTS", []int{3000, 3001})
```

**Пример:**

```bash
export PORTS=8080,8081,8082
```

```go
ports := envconfig.GetIntSlice("PORTS", []int{3000, 3001})
// ports = []int{8080, 8081, 8082}
```

**Примечания:**
- Значения разделяются запятыми
- Пробелы вокруг значений автоматически удаляются
- Пустые значения заменяются на 0
- Поддерживаются отрицательные числа

### GetInt64Slice(key string, defaultValue []int64) []int64

Получает слайс 64-битных целых чисел из переменной окружения с дефолтным значением. Значения должны быть разделены запятыми. Пробелы вокруг значений автоматически удаляются. Возвращает значение по умолчанию, если переменная окружения не установлена, пуста или содержит невалидные значения.

```go
maxSizes := envconfig.GetInt64Slice("MAX_SIZES", []int64{512, 1024})
```

**Пример:**

```bash
export MAX_SIZES=1024,2048,4096
```

```go
maxSizes := envconfig.GetInt64Slice("MAX_SIZES", []int64{512, 1024})
// maxSizes = []int64{1024, 2048, 4096}
```

**Примечания:**
- Значения разделяются запятыми
- Пробелы вокруг значений автоматически удаляются
- Пустые значения заменяются на 0
- Поддерживаются отрицательные числа

### ToList(value string, separator string) ([]string, error)

Разделяет строку на список строк по указанному разделителю.

```go
values, err := envconfig.ToList("a,b,c", ",")
// values = []string{"a", "b", "c"}
```

## Примеры использования

### Пример 1: Простая конфигурация сервера

```go
package main

import (
    "fmt"
    "log"
    "github.com/pgmod/envconfig"
)

type ServerConfig struct {
    Host     string `env:"SERVER_HOST" default:"0.0.0.0"`
    Port     int    `env:"SERVER_PORT" default:"8080"`
    Timeout  int    `env:"TIMEOUT" default:"30"`
    Debug    bool   `env:"DEBUG" default:"false"`
}

func main() {
    var cfg ServerConfig
    
    if err := envconfig.LoadStruct(&cfg); err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Запуск сервера на %s:%d\n", cfg.Host, cfg.Port)
}
```

### Пример 2: Конфигурация с массивами

```go
type AppConfig struct {
    DatabaseURL string   `env:"DATABASE_URL" default:"postgres://localhost/db"`
    WorkerIDs   []int    `env:"WORKER_IDS" default:"1,2,3"`
    Ports       [3]int   `env:"PORTS" default:"8080,8081,8082"`
    Workers     int      `env:"WORKERS" default:"4"`
}
```

**Переменные окружения:**

```bash
export DATABASE_URL=postgres://prod.example.com/mydb
export WORKER_IDS=1,2,3,4,5
export PORTS=9000,9001,9002
export WORKERS=8
```

### Пример 3: Комбинирование Load() и LoadStruct()

```go
package main

import (
    "log"
    "github.com/pgmod/envconfig"
)

func main() {
    // Сначала загружаем .env файл
    if err := envconfig.Load(); err != nil {
        log.Printf("Предупреждение: не удалось загрузить .env файл: %v", err)
    }
    
    // Затем загружаем конфигурацию в структуру
    var cfg Config
    if err := envconfig.LoadStruct(&cfg); err != nil {
        log.Fatal(err)
    }
}
```

### Пример 4: Использование простых функций

```go
package main

import (
    "fmt"
    "github.com/pgmod/envconfig"
)

func main() {
    host := envconfig.Get("HOST", "localhost")
    port := envconfig.GetInt("PORT", 8080)
    debug := envconfig.GetBool("DEBUG", false)
    ports := envconfig.GetIntSlice("PORTS", []int{3000, 3001})
    maxSizes := envconfig.GetInt64Slice("MAX_SIZES", []int64{512, 1024})
    
    fmt.Printf("Подключение к %s:%d (debug: %v)\n", host, port, debug)
    fmt.Printf("Порты: %v\n", ports)
    fmt.Printf("Максимальные размеры: %v\n", maxSizes)
}
```

## Формат .env файла

```env
# Комментарии начинаются с #
HOST=localhost
PORT=8080
DEBUG=true
DATABASE_URL=postgres://localhost/mydb
PORTS=8080,8081,8082
FIXED=1,2,3
```

## Обработка ошибок

Все функции возвращают ошибки, которые следует обрабатывать:

```go
var cfg Config
if err := envconfig.LoadStruct(&cfg); err != nil {
    // Ошибка может возникнуть при:
    // - передаче не указателя на структуру
    // - невалидных значениях переменных окружения
    // - несоответствии размера массива количеству значений
    log.Fatalf("Ошибка загрузки конфигурации: %v", err)
}
```

## Ограничения

- В `LoadStruct()` поддерживаются только типы: `string`, `bool`, `int`, `int64`, `[]int`, `[]int64`, массивы `int`
- Для массивов и слайсов поддерживаются только типы `int` и `int64`
- Значения массивов должны быть разделены запятыми
- Массивы требуют точного соответствия количества значений размеру массива

## Лицензия

MIT
