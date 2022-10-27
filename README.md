# GO-Service-Example

## Running the Application

1. Create database for mock
```sh
# Pembuatan database
psql -U <username>
create database <database name>;
exit

psql -U <username> -f sql/user.sql
```
2. Configur **env** file inside `config/*` (you can choose to use `config.yaml` or `.env`)

```go
// file main.go
func main() {
    // using config.yaml
	c := config.InitConfig()

    // using .env
    c := config.ReadEnv()
}
```

3. Running the application
```go
go run .
```

***

# Project Structure 

```
├───config
├───datasource
├───job
├───service
├───test
└───utils
    └───logger
```

***

# Overview
Dalam projek ini digunakan library external seperti :
> [**viper**](https://github.com/spf13/viper) untuk membaca file env\
> [**cron**](https://github.com/robfig/cron) untuk menerapkan cron job\
> [**pgx**](https://github.com/jackc/pgx) sebagai database driver untuk database postgres\
> [**zap**](https://github.com/uber-go/zap) sebagai logger

Projek ini berisi 2 dummy service yang bernama `Service A` dan `Service B`.

> ### Service A
> Service A berisi 2 function yaitu `DoSomething()` yang mengembalikan *hello world* beserta nama service dan fungsi `CopyToCSV()` yang mengambil data dari database dan menyimpan hasil tersebut kedalam file `.csv`

> ### Service B
> Service B berisi 1 function yaitu `DoSomething()` yang mengembalikan *hello world* berserta nama service.

Kemudian kedua service tersebut akan dijalankan oleh `cron job` yang terdapat dalam package `job`. Berikut penerapan cron job untuk kedua service diatas :
```go
// file job/job.go
func InitJob(...) {
    // ...

    // Job every Mon-Sat at 23:50
	c.AddFunc("50 23 * * 1-6", func() {
		s.ServiceA().DoSomething()
	})

	// Job every 5s
	c.AddFunc("@every 10s", func() {
		s.ServiceB().DoSomething()
	})

	// Job every 1m
	c.AddFunc("@every 1m", func() {
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, 40*time.Second)
		defer cancel()

		s.ServiceA().CopyToCSV(ctx)
	})

    c.Start()
}
```

***
# How to
Berikut adalah langkah - langkah untuk membuat service seperti repository ini. Dalam dokumentasi ini terdapat 3 langkah yang terdiri dari `Single Service`, `Multiple Service`, dan `Cron Job`.

## Single Service
Untuk membuat sebuah service, pertama kali dilakukan adalah `planning` struktur sebuah service. Pada tahap `planning` ini, kita akan menentukan fungsi yang terdapat pada sebuah service. Pada `Service A` akan dibuat 2 fungsi yang terdiri dari `DoSomething()` yang melakukan print hello world dan `CopyToCsv()` yang akan melakukan operasi database.

### Service A
Pertama kali akan dibuat folder berupa `service` yang memiliki isi folder `service-a`. Setelah itu akan dibuat sebuah file bernama `main.go` dengan kedua fungsi diatas.

> **Fungsi** `DoSomething()`
> ```go
> // file service/service-a/main.go
> package servicea
>
> func DoSomething() {
>	fmt.Println("Hello world from service-a")
> }
> ```

> **Fungsi** `CopyToCsv()`
> ```go
> // file service/service-a/main.go
> package servicea
>
> func CopyToCsv() {
>	query := `
>		COPY (SELECT * 
>			FROM users 
>		) TO STDOUT DELIMITER ',' CSV HEADER
>	`
>   // Create database connection
>   // ...
> }
> ```

Dalam fungsi `CopyToCsv()`, dibutuhkan koneksi database. Oleh karena itu akan digunakan sebuah `Database Driver` atau `ORM`. Untuk membuat sebuah `Database Driver` atau `ORM` dibuatlah sebuah package bernama `datasource` yang menyimpan objek koneksi database.

### Datasource
Dalam projek ini digunakan library `pgx` yang merupakan database driver untuk `Postgres`.

```go
// file datasource/postgres.go
package datasource

func NewPgConn(dsn string) *pgx.Conn {
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		panic(fmt.Errorf("failed to create database connection\n%v", err))
	}

	return conn
}
```

### Config
Dari potongan kode sebelumnya, fungsi tersebut tidak dapat langsung digunakan karena dibutuhkan `dsn` yang berupa `Data Source Name`. Oleh karena itu akan dibuat sebuah package yang menghandle variabel dsn dari sebuah file. Package tersebut adalah `config`. Untuk membaca variabel dari sebuah file dapat digunakan library `viper`.

```go
// file config/config.go

// variabel yang digunakan untuk menyimpan nilai
// dari file env
type config struct {
	Dsn  string `mapstructure:"dsn"`
}

func InitConfig() *config {
	var c config

    // Mengatur file konfigurasi
	viper.AddConfigPath(filepath.Join(rootDir(), "config"))

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

    // Melakukan pengecekan file env
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("config file not found\n%v", err))
	}

    // Mengambil nilai dari file env
	if err := viper.Unmarshal(&c); err != nil {
		panic(fmt.Errorf("failed to read configuration\n%v", err))
	}

	return &c
}

func rootDir() string {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	return filepath.Dir(d)
}
```

Berikut contoh isi dari file `.env` atau `config.yml`.
```env
dsn=postgres://username:password@host:port/dbName
```
```yaml
dsn: postgres://username:password@host:port/dbName
```

### Implementasi
Setelah membuat package `config`, kita akan mengirimkan data tersebut kedalam `datasource`. Proses ini akan dilakukan dalam file `main.go`. Berikut kodingan pada `main.go`.

```go
// file main.go
package main

func main() {
	// Read env configs
	c := config.InitConfig()

	// Init datasources
	pgConn := datasource.NewPgConn(c.Dsn)
}
```

Apabila kita ingin menambahkan koneksi database lain, kita dapat menerapkan `struct` **DB** dalam package `datasource`. Berikut contoh penerapannya :

```go
// file datasource/postgres.go
package datasource

type DB struct {
	p *pgx.Conn
    // Another db connection
    m *anotherConn
}

func NewDB(dsn string) *DB {
	return &DB{
		newPgConn(dsn),
        anotherConn(dsn),
	}
}

func newPgConn(dsn string) *pgx.Conn {
	// ...
	return conn
}

func anotherConn(dsn string) *anotherConn {
    // ...
    return conn
}
```

```go
// file main.go

func main() {
	c := config.InitConfig()
	db := datasource.NewDB(c.Dsn)
}
```
Kita telah mendapatkan database connection, kemudian langkah selanjutnya adalah menerapkan service. Sebelumnya kita dapat langsung menerapkan service seperti ini :

```go
// file main.go
package main

func main() {
	c := config.InitConfig()
	db := datasource.NewDB(c.Dsn)

	servicea.DoSomething()
	servicea.CopyToCsv(db)
}

// file service/service-a/main.go
package servicea

func CopyToCsv(pg *pgx.Conn) {
	query := `
		COPY (SELECT * 
			FROM users 
		) TO STDOUT DELIMITER ',' CSV HEADER
	`
   	res, err := pg.PgConn().CopyTo(context.Background(), file, query)

	// ...
}
```

Potongan kode diatas memang langsung bisa memanggil fungsi, hanya saja terdapat kelemahan apabila kita memiliki banyak fungsi dalam sebuah `service`. Hal tersebut membuat kita memasukkan satu per satu `dependency` fungsi tersebut.

```go
func main() {
	c := config.InitConfig()
	db := datasource.NewDB(c.Dsn)
	something := something.Something()

	// Contoh dari pemanggilan berulang yang
	// dapat memenuhi main.go
	servicea.DoSomething()
	servicea.CopyToCsv(db)
	servicea.DoSomethingSomething(something)
	servicea.CallSomething(something)
}
```

### Service Struct
Untuk menghindari hal diatas, kita dapat membuat sebuah `struct` pada **Service A**. Struct tersebut akan diisi oleh `dependency` yang dibutuhkan Service A. Contoh dari dependency tersebut adalah `db.Conn`, `logger`, dan sebagainya. Hal ini menghindari kita untuk melakukan **passing dependency** berulang - ulang ke dalam function 

```go
// file service/servicea/main.go
package servicea

// Struct untuk menyimpan dependency
type ServiceA struct {
	db *pgx.Conn
}

// Mengembalikan struct untuk mengakses function-nya
func NewServiceA(db *pgx.Conn) *ServiceA {
	return &ServiceA{
		db: db,
	}
}

func (s *ServiceA) DoSomething() {
	fmt.Println("Hello world from service-a")
}

func (s *ServiceA) CopyToCSV() {
	query := ...

	// Kita hanya perlu memanggil s.db milik `struct ServiceA`
   	res, err := s.db.PgConn()...
}

// file main.go
package main

func main() {
	c := config.InitConfig()
	db := datasource.NewDB(c.Dsn)

	s := servicea.NewServiceA(db)
	s.DoSomething()
	s.CopyToCsv()

	// ...
}
```

## Multiple Service

Setelah membuat sebuah service bernama `Service A`, developer mendapatkan permintaan untuk menambahkan service lain bernama `Service B` dengan isi sebagai berikut :

```go
// file service/service-b/main.go
package serviceb

type ServiceB struct {
	db *pgx.Conn
}

func NewServiceB() *ServiceB {
	return &ServiceB{}
}

func (s *ServiceB) DoSomething() {
	fmt.Println("Hello World from service-b")
}

func (s *ServiceB) ConnectToDB() {
	res, err := db.Conn...
	// ...
}

```

Tentu untuk menerapkan service tersebut kita dapat menggunakkan cara penerapan `Service A`, hanya saja penerapan tersebut kurang efektif apabila kedua service (`Service A` dan `Service B`) memiliki dependency yang sama. Oleh karena itu akan dibuat sebuah `service registry` untuk mendaftarkan berbagai macam service.

### Service Registry
Untuk menerapkan service registry, pertama kali yang dilakukan adalah membuat file `service.go`. File tersebut berisi `struct` yang berperan untuk menyimpan `global dependency` bagi seluruh service.

```go
// file service/service.go
package service

type Service struct {
	db *pgx.Conn
}

func NewService(db *pgx.Conn) *Service {
	return &Service{
		db: db,
	}
}
```

Setelah membuat service tersebut, kita hanya perlu memanggil sekali pada file `main.go`. Berikut contoh penerapannya.

```go
// file main.go
func main() {
	c := config.InitConfig()
	db := datasource.NewDB(c.Dsn)

	s := service.NewService(db)
}
```

Kemudian agar mendapatkan akses pada Service A dan Service B, kita akan mendaftarkan kedua service tersebut kedalam `service registry` yang merupakan file `service.go`. Penerapan service tersebut dilakukan dengan cara penerapan service yaitu menggunakan `struct function`.

```go
// file service/service.go
package service 

type Service struct {
	db *datasource.DB
}

func NewService(db *datasource.DB) *Service {
	return &Service{
		db: db,
	}
}

func (s *Service) ServiceA() *servicea.ServiceA {
	return servicea.NewServiceA(s.db)
}

func (s *Service) ServiceB() *serviceb.ServiceB {
	return serviceb.NewServiceB(s.db)
}

// file main.go
package main

func main() {
	c := config.InitConfig()
	db := datasource.NewDB(c.Dsn)

	s := service.NewService(db)

	s.ServiceA().DoSomething()
	s.ServiceB().DoSomething()
}
```

### *Penerapan Optional*

Untuk menerapkan `multiple service`, dapat digunakan cara sebagai berikut.
```go
type Services struct {
	ServiceA *servicea.ServiceA
	ServiceB *serviceb.ServiceB
}

type Service struct {
	l  *logger.Logger
	db *datasource.DB
}

func New(l *logger.Logger, db *datasource.DB) *Service {
	return &Service{
		l:  l,
		db: db,
	}
}

func (s *Service) NewServices() *Services {
	return &Services{
		ServiceA: servicea.NewServiceA(s.l, s.db),
		ServiceB: serviceb.NewServiceB(),
	}
}
```

Cara yang digunakan adalah dengan membuat sebuah struct `Service` yang berfungsi untuk menyimpan global dependency dan struct `Services` yang menyimpan seluruh services yang ada. Kemudian untuk mengakses keseluruhan service, developer dapat menggunakkan cara berikut :

```go
func main() {
	// ...
	s := service.New(l, db)
	ss := s.NewServices()

	ss.ServiceA.DoSomething()

	// ...
}

```

## Cron Job
Setelah menerapkan kedua service (Service A dan Service B), developer diminta untuk membuat sebuah `cron job` yang menjalankan fungsi dalam kedua service tersebut. Untuk menerapkan cron job, kita akan menggunakkan sebuah library yaitu [robfig/cron](https://github.com/robfig/cron).

Pertama kali kita akan membuat sebuah package yang mendaftarkan seluruh cron job. Package tersebut dinamakan `job`. Dalam package tersebut dibuat sebuah file bernama `job.go`.

```go
// file job/job.go
package job

import "github.com/robfig/cron/v3"

func InitJob() {
	c := cron.New()

	// init job here

	c.Start()
}
```

Setelah membuat file `job/job.go` kita akan mendaftarkan fungsi service yang akan dijalankan oleh cron job. Oleh karena itu kita membutuhkan **service registry** berupa struct `Service` dalam file `service/service.go`.

```go
// file job/job.go
package job

import (
	"github.com/MinorvaFalk/go-service-example/service"
	"github.com/robfig/cron/v3"
)

func InitJob(s *service.Service) {
	c := cron.New()

	// init job here

	c.Start()
}
```

Kemudian akan dilakukan passing service ke job dalam file `main.go`. Berikut penerapannya dalam file `main.go` :

```go
func main() {
	c := config.InitConfig()
	db := datasource.NewDB(c.Dsn)

	s := service.NewService(db)
	job.InitJob(s)

	// Blocking operation for cron job
	// operasi ini dapat diganti dengan infinite loop hingga mutex
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
```

Setelah mendapatkan service dalam job, kita akan menjalankan fungsi dari masing masing service yang ada. Untuk cara menjalankan `cron job` sebuah service sesuai kebutuhan, kita dapat membaca dokumentasi [godoc robfig/cron](https://pkg.go.dev/github.com/robfig/cron?utm_source=godoc). Berikut contoh penerapan cron job.

```go
package job

func InitJob(s *service.Service) {
	c := cron.New(cronOption...)
	// refer https://pkg.go.dev/github.com/robfig/cron for more about cronjob functions

	// Job every Mon-Sat at 23:50
	c.AddFunc("50 23 * * 1-6", func() {
		s.ServiceA().DoSomething()
	})

	// Job every 10 second
	c.AddFunc("@every 10s", func() {
		s.ServiceB().DoSomething()
	})

	// Job every 1 minute
	c.AddFunc("@every 1m", func() {
		s.ServiceA().CopyToCSV()
	})

	c.Start()
}
```