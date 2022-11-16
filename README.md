# REST API using Go Fiber and GORM

<img alt="gofiber logo" src="https://gofiber.io/assets/images/embed.png" width="580"></img>

## 1. 프로젝트 개요

![Performance Benchmark of Gofiber and Others](https://taejoone.jeju.onl/assets/img/2022/11/15-benchmark-requests-fiber-crunch.png){: width="560"}

최대 6~7배 정도 빠르다는데, 정말인지는 써보면서 알아보자.

### 1) 기능 설명

GoFiber 와 GORM 을 이용해 간단한 REST API 를 구현함 (boilerplate)

#### Go-fiber 웹서버

- `.env` 로부터 DB_URL, PORT 등을 읽어 적용
- Go-Fiber 의 미들웨어 조립 : Logger, CORS, Cache, Views
- API Group 과 GET / POST / PUT / PATCH / DELETE 메소드
- 다양한 Route Parameters 형식을 등록하여 테스트
- Query Params 읽어와 DB Where 조건에 사용
- View Engine 을 마운트 하고 HTML 템플릿 페이지 출력
- Cache 미들웨어를 등록했으나, refresh 가 작동하지 않음
  - 디버깅을 위해 cacheHit 페이지 추가

#### GORM & SQLite

- 최초 샘플 데이터 입력
- SQLite DB 에 대해 CRUD 구현
- validator 로 입력 struct 에 적합한지 타입 검사
- GORM 의 sql.NullInt16 필드를 사용 (Dog.Age)
  - JSON 출력을 위해 별도의 MarshalJSON/UnmarshalJSON 함수를 구현
  - Null 업데이트를 위해 Age 에 대한 Update 문을 추가
  - 트랜잭션을 사용해 실패시 Rollback 처리

#### 그 외 (API 와 관계없지만)

- AES 암호화/복호화
- 하위 디렉토리 모듈 임포트 연습
- 여러 예제와 유틸리티들을 모두 모아서 작성

참조 : [How to Build REST API using Go Fiber and Gorm ORM](https://dev.to/franciscomendes10866/how-to-build-rest-api-using-go-fiber-and-gorm-orm-2jbe)

## 2. 프로젝트 Setup

```bash
$ mkdir fiber-example && cd fiber-example
$ go mod init example.com

$ cat <<EOF > main.go
package main
func main() {}
EOF

$ go get -u gorm.io/gorm
$ go get -u gorm.io/driver/sqlite
$ go get -u github.com/gofiber/fiber/v2
$ go get -u github.com/joho/godotenv
$ go get -u github.com/gofiber/template
$ go get -u gopkg.in/go-playground/validator.v9

$ go get -u golang.org/x/exp/maps   # maps.Keys() 함수
$ go get -u golang.org/x/exp/slices # slices.Contains() 함수

$ go mod tidy

$ go run .
9f4yohBU0rUoq6ajOcC3hA==
hello world
{1 Go}
false
2022/11/15 15:03:12 init: 3 records inserted
views: parsed template: index

Fiber v2.39.0
http://127.0.0.1:3000
# ...
2022/11/15 19:06:03 params = map[]
19:06:03 | 200 |     1ms |       127.0.0.1 | GET     | /api/dogs
```

### main.go

```go
// main.go
import (
  "example.com/db"  // DB 접속 및 CRUD 함수
  m "example.com/models"  // 모델 및 JSON 변환, 인터페이스 함수
  u "example.com/utils"  // map 처리, env 등등 유틸리티 함수들
  "example.com/web"  // 웹서버 미들웨어 및 라우터 설정
)

func main() {
  db.Connect()

  app := fiber.New()
  web.SetupFiber(app)

  var port = db.Config("PORT")
  log.Fatal(app.Listen(":" + port))
}
```

## 3. [Go-Fiber](https://docs.gofiber.io/) 웹서버

### 1) [미들웨어](https://docs.gofiber.io/api/middleware)

- Logger
- CORS
- Cache : 기본으로 메모리 캐시를 사용
  - refresh 쿼리 파라미터가 들어가면 캐시 갱신이 되어야 하는데 안됨
    - refresh 쿼리 파라미터까지 캐싱되어 통째로 무시되는듯 함

```go
  // Logger middleware
  app.Use(logger.New(logger.Config{
    Format: "[${time}] ${status} - ${latency} ${method} ${path}\n",
  }))

  // CORS middleware
  app.Use(cors.New(cors.Config{
    // AllowOrigins: "https://gofiber.io, https://gofiber.net",
    AllowOrigins: "*",
    AllowHeaders: "Origin, Content-Type, Accept",
  }))

  // 특정 API 그룹에만 캐시 적용
  cacheGroup := app.Group("/click")
  // Cache middleware
  cacheGroup.Use(cache.New(cache.Config{
    Next: func(c *fiber.Ctx) bool {
      return c.Query("refresh") == "true"
    },
    Expiration:   30 * time.Minute,
    CacheControl: true,
  }))
```

#### cache Hit 조사

- `/click`, `/click?refresh=true` 를 여러차례 요청
- `/cacheHits` 에서 캐시 Hit 비율을 출력시켰는데
  - 캐시 갱신이 먹지 않는다. (hander 함수에 진입하지 못함)

```js
{
  "cacheHits": 5,
  "cacheHitsPercentage": 83,
  "requests": 6
}
```

### 2) [HTML 템플릿](https://docs.gofiber.io/api/middleware) - server-side template engines

템플릿 view 를 `/home` 에 연결 (API 와 함께 사용)

- {PRJ_ROOT}/views
  - index.html

```html
<!-- 템플릿 파일 index.html -->
<!DOCTYPE html>
<body>
    <h1>{ {.Title} }</h1>
    <p>{ { greet "Fiber" } }</p> <!-- 사용자 함수 greet 사용 -->
</body>
</html>
```

![gofiber-template-html](https://taejoone.jeju.onl/assets/img/2022/11/15-gofiber-template-html-crunch.png){: width="420"}

### 3) 파라미터

참고 : [Stackoverflow - How to iterate over query params in Golang](https://stackoverflow.com/a/73736090/6811653)

- 설명으로는 `c.AllParams()` 로 모든 Query 파라미터를 가져온다는데
  - 안된다. 소스 코드를 봐도 딱히 안될 부분은 안보이는데.

그래서 따로 `getQueryParams` 함수를 작성했더니, 이건 된다.

```go
  params := getQueryParams(c)
  log.Printf("params = %+v", params)

  // GET /api/dogs?age=5&name=abc
  // ==> params: map[string]string{"age":"5", "name":"abc"}

////////////////////////////////////

func getQueryParams(c *fiber.Ctx) map[string]string {
  params := make(map[string]string)
  var err error
  c.Context().QueryArgs().VisitAll(func(key, val []byte) {
    if err != nil {
      return
    }
    k := utils.UnsafeString(key)
    v := utils.UnsafeString(val)
    params[k] = v
  })
  return params
}
```

### 3) Endpoints

- PUT 은 필드 전체를 업데이트하고, PATCH 는 부분 업데이트를 한다

- 라우터의 Path 파라미터에 제약사항을 설정할 수 있다.
  - 제약사항에 위배되면 `404 Not Found` 로 처리됨
  - 참고 [Route constraints](https://docs.gofiber.io/guide/routing#constraints)

```go
  // Create a new route group '/api'
  api := app.Group("/api")

  // id 는 int 만 가능
  api.Get("/dogs", db.GetDogs)
  api.Get("/dogs/:id<int>", db.GetDog)
  api.Post("/dogs", db.AddDog)
  api.Put("/dogs/:id<int>", db.UpdateDog)
  api.Patch("/dogs/:id<int>", db.UpdateDogPartial)
  api.Delete("/dogs/:id<int>", db.RemoveDog)
```

## 4. [GORM](https://gorm.io/docs/) with SQLite3

### 1) 설정

#### DB 모델을 위한 `Dog` 구조체

- sql.NullInt16 대신에 wrapper 타입 NullInt16 을 사용
  - Null 처리가 가능하면서 JSON 출력시 값만 나오게 하려고 적용
  - 참고 [How can I work with SQL NULL values and JSON?](https://stackoverflow.com/a/33072822/6811653)

```go
// Dog type with sql.NullInt16
type Dog struct {
  ID        int       `json:"id" gorm:"primaryKey"`
  Name      string    `json:"name" validate:"required,min=3,max=32"`
  Breed     string    `json:"breed" validate:"required"`
  Age       NullInt16 `json:"age" validate:"number" form:"age"`
  IsGoodBoy bool      `json:"isGoodBoy" gorm:"default:true"`
}

// NullInt16 is wrapper for sql.NullInt16
// 참고 https://stackoverflow.com/a/33072822/6811653
type NullInt16 struct {
  sql.NullInt16
}

// ToNullInt16 convert int to sql.NullInt16
func ToNullInt16(v int) NullInt16 { ... }

// MarshalJSON marshal json of NullInt16
func (v NullInt16) MarshalJSON() ([]byte, error) { ... }

// UnmarshalJSON unmarshal json of NullInt16
func (v *NullInt16) UnmarshalJSON(data []byte) error { ... }
```

#### 초기 데이터 삽입

- nullable 필드를 누락하면, null 또는 기본값이 들어간다

```go
  //You can insert multiple records too
  var dogs []m.Dog = []m.Dog{
    {Name: "Ricky", Breed: "Chihuahua", Age: m.ToNullInt16(2), IsGoodBoy: false},
    {Name: "Adam", Breed: "Pug", IsGoodBoy: true},
    {Name: "Justin", Breed: "Poodle", Age: m.ToNullInt16(3), IsGoodBoy: false},
  }
  tx := db.Create(&dogs)
```

### 2) 트랜잭션

#### [Updates multiple columns](https://gorm.io/docs/update.html#Updates-multiple-columns) - 다수의 필드 업데이트

다중 필드 업데이트는 구조체 또는 Map 인터페이스로 할 수 있다.

- 반드시 대상을 특정할 수 있는 ID 가 명시되어야 함
- 단, not-Null / non-Zero 값만 업데이트함
  - `User{Active: false}` => 무시/누락

> NOTE When updating with struct, GORM will only update non-zero fields. You might want to use map to update attributes or use Select to specify fields to update

> **주의!!** 구조체로 업데이트할 때 GORM은 0이 아닌 필드만 업데이트합니다. 지도를 사용하여 속성을 업데이트하거나 선택을 사용하여 업데이트할 필드를 지정할 수 있습니다.

#### PUT `/dogs/:id<int>` 전체 필드 업데이트

Null / Zero 값을 업데이트 하려면 Select 를 포함하여야 함

- 트랜잭션 처리 (절차식으로 나열하는 것보다 함수형이 안전하다)
  - 오류가 나면 err 를 내보내고, 맨 나중에 웹응답 처리

1. 트랜잭션 진입
2. `tx.Model(&dog)` 로 갱신 대상 테이블을 알려주고
3. `Select("*")` 로 필드 전체가 갱신 대상임을 알려주고
4. `Where("ID = ?", id)` 로 업데이트 대상을 명시하고
5. `Omit("ID")` 혹시나 중요 필드가 업데이트 되지 않도록 보호
6. Body 에서 받아온 struct 데이터로 `Updates(dog)` 적용
7. 별 문제 없으면 nil 반환 (커밋)

```go
  // Transaction return nil or error
  err := Database.Transaction(func(tx *gorm.DB) error {
    id := c.Params("id")
    // 모든 필드에 대해 업데이트 (ID 제외)
    if err := tx.Model(&dog).Select("*").Where("ID = ?", id).Omit("ID").Updates(dog).Error; err != nil {
      return err
    }
    return nil // commit
  })

  if err != nil {
    log.Fatalln(err)
    return c.Status(503).SendString(err.Error())
  }
  return c.Status(200).JSON(dog)
```

#### PATCH `/dogs/:id<int>` 부분 필드 업데이트

Body 를 통해 생성된 모델 구조체는 모든 필드를 포함하고 있다. 따라서, Select 를 이용해 갱신 대상을 제한하도록 해야 한다. (안그러면 필드 전체가 변경됨)

> Select 사용시 JSON 태그명을 구조체의 필드명으로 바꾸어 주어야함

1. 모델 구조체에서 필드명과 JSON 태그명 사전(map)을 생성
2. c.Body() 에서 사용된 JSON 태그명 슬라이스를 추출
3. 사전(map) 으로 업데이터 대상인 필드명 슬라이스를 생성
4. 트랜잭션 진입
5. `Select(fields)` 과 함께 `Updates(dog)` 적용
6. 이상 없으면 nil 반환 (커밋)

```go
  tableName, fieldNames := GetTableJSONTags(Database, dog)
  if fieldNames == nil {
    return c.Status(503).SendString("Any JSON tag is not defined")
  }

  // 업데이트 대상 json tag 추출
  var tags []string = u.ExtractFields(c.Body())
  // json tag 를 field name 로 변환
  var fields []string = u.ReplaceSliceByMap(tags, u.MapS(fieldNames).Reverse())
  log.Printf("%s: tags %v => fields %+v", tableName, tags, fields)

  err := Database.Transaction(func(tx *gorm.DB) error {
    id := c.Params("id")
    // 업데이트 대상 필드(fields)들만 업데이트
    if err := tx.Model(&dog).Select(fields).Where("ID = ?", id).Updates(dog).Error; err != nil {
      return err
    }
    return nil // commit
  })

  if err != nil {
    log.Fatalln(err)
    return c.Status(503).SendString(err.Error())
  }
  return c.Status(200).JSON(dog)
```

### 3) Delete 할 때 사전에 검사하기 위해 Hook (훅) 사용

샘플데이터 ID=[1,2,3] 에 대해 삭제하지 못하도록 검사 후 삭제

1. ID 값으로 Delete 실행
2. BeforeDelete 인터페이스 함수 (Hook) 진입
3. 검사할 수 있는 값은 구조체 값뿐이라 사전에 ID 값을 넣어두어야 함!
4. 조건을 만족하지 않으면 Error 반환 (취소됨)
5. 이상 없으면, Delete 적용

참고 [GORM - Delete Hooks](https://gorm.io/docs/delete.html#Delete-Hooks)

```go
func RemoveDog(c *fiber.Ctx) error {
  id, err := strconv.Atoi(c.Params("id"))

  var dog m.Dog = m.Dog{ID: id} // for BeforeDelete
  result := Database.Model(&dog).Delete(&dog, id)
  // ...
}

// BeforeDelete prevent delete sample data which ID < 4
// **NOTE: 같은 모듈 안에서만 정의할 수 있음
func (d *Dog) BeforeDelete(tx *gorm.DB) (err error) {
  if d.ID < 4 {
    log.Printf("cancel: ID=%d", d.ID)
    return errors.New("Sample Data (ID<4) not allowed to delete")
  }
  return
}
```

### 4) CRUD 실행

#### REST API 요청 및 결과

```js
// GET http://localhost:3000/api/dogs HTTP/1.1
[
  { "id": 1, "name": "Ricky", "breed": "Chihuahua",
    "age": 2, "isGoodBoy": false
  },
  { "id": 2, "name": "Adam", "breed": "Pug",
    "age": null,          // <-- nullable
    "isGoodBoy": true
  },
  { "id": 3, "name": "Justin", "breed": "Poodle",
    "age": 3, "isGoodBoy": false
  }
]

// POST http://localhost:3000/api/dogs HTTP/1.1
// {
//   "name": "Max Junior",
//   "breed": "Shepherd",
//   "age": 4,
//   "isGoodBoy": true
// }
{
  "id": 4,
  "name": "Max 2nd",
  "breed": "Shepherd",
  "age": 4,
  "isGoodBoy": true
}

// PUT http://localhost:3000/api/dogs/1 HTTP/1.1
// {
//   "name": "Max Junior",
//   "breed": "Shepherd (German)",
//   "age": 9
// }
{
  "id": 0,                // <-- Omit
  "name": "Max Junior",
  "breed": "Shepherd (German)",
  "age": 9,
  "isGoodBoy": false
}

// GET http://localhost:3000/api/dogs/4 HTTP/1.1
{
  "id": 4,
  "name": "Max Junior",
  "breed": "Shepherd (German)",
  "age": 9,
  "isGoodBoy": false
}

// PATCH http://localhost:3000/api/dogs/4 HTTP/1.1
// {
//     "age": null,
//     "isGoodBoy": false
// }
{
  "id": 0,
  "name": "",
  "breed": "",
  "age": null,          // <-- select
  "isGoodBoy": false    // <-- select
}

// GET http://localhost:3000/api/dogs/4 HTTP/1.1
{
  "id": 4,
  "name": "Max Junior",
  "breed": "Shepherd (German)",
  "age": null,
  "isGoodBoy": false
}

// DELETE http://localhost:3000/api/dogs/1 HTTP/1.1
[403 Forbidden]
Sample Data (ID<4) not allowed to delete

// DELETE http://localhost:3000/api/dogs/4 HTTP/1.1
OK

// GET http://localhost:3000/api/dogs/4
Not Found
```

#### Gofiber 로깅

```bash
 ┌───────────────────────────────────────────────────┐
 │                   Fiber v2.39.0                   │
 │               http://127.0.0.1:3000               │
 │       (bound on host 0.0.0.0 and port 3000)       │
 │                                                   │
 │ Handlers ............ 34  Processes ........... 1 │
 │ Prefork ....... Disabled  PID ............. 31294 │
 └───────────────────────────────────────────────────┘

14:47:10 | 200 |     1ms |       127.0.0.1 | GET     | /api/dogs
14:48:50 | 201 |     2ms |       127.0.0.1 | POST    | /api/dogs
2022/11/16 14:50:16 update: &{ID:0 Name:Max Junior Breed:Shepherd (German) Age:{NullInt16:{Int16:9 Valid:true}} IsGoodBoy:false}
14:50:16 | 200 |      0s |       127.0.0.1 | PUT     | /api/dogs/4
14:52:56 | 200 |      0s |       127.0.0.1 | GET     | /api/dogs/4
2022/11/16 14:53:06 dogs: tags [age isGoodBoy] => fields [Age IsGoodBoy]
14:53:06 | 200 |      0s |       127.0.0.1 | PATCH   | /api/dogs/4
14:53:11 | 200 |      0s |       127.0.0.1 | GET     | /api/dogs/4
2022/11/16 14:53:47 cancel: ID=1
14:53:46 | 403 |      0s |       127.0.0.1 | DELETE  | /api/dogs/1
14:54:00 | 200 |      0s |       127.0.0.1 | DELETE  | /api/dogs/4
14:54:32 | 404 |      0s |       127.0.0.1 | GET     | /api/dogs/4
```

## 5. Others

### 1) 유틸리티 함수들

#### [golang.org/x/exp/maps](https://pkg.go.dev/golang.org/x/exp/maps) - Generic 타입 실험 패키지

> `golang.org/x/exp` 는 실험적인 또는 폐기된 패키지를 포함하고 있기 때문에, 반드시 하위 디렉토리까지 지정해서 사용하도록 경고하고 있음

- map 타입에서 Key 추출하기
  - ['maps.Keys' - Go Playground](https://go.dev/play/p/fkm9PrJYTly)

```go
import (
  "fmt"

  "golang.org/x/exp/maps"
)

func main() {
  intMap := map[int]int{1: 1, 2: 2}
  intKeys := maps.Keys(intMap)
  // intKeys is []int
  fmt.Println(intKeys)

  strMap := map[string]int{"alpha": 1, "bravo": 2}
  strKeys := maps.Keys(strMap)
  // strKeys is []string
  fmt.Println(strKeys)
}
// 출력 ==>
// [2 1]
// [alpha bravo]
```

#### [slices 의 Contains 함수](https://stackoverflow.com/a/71181131/6811653)

```go
// go get golang.org/x/exp/slices
import  "golang.org/x/exp/slices"

things := []string{"foo", "bar", "baz"}
slices.Contains(things, "foo") // true
```

#### 문자열 map 의 Key 와 Value 뒤바꾸기

```go
// MapS is a map with string keys and values.
type MapS map[string]string

// Reverse returns a new map with the keys and values swapped.
func (m MapS) Reverse() map[string]string {
  n := make(map[string]string, len(m))
  for k, v := range m {
    n[v] = k
  }
  return n
}
```

#### 인터페이스 map 을 특정 Key 리스트로 필터링하기

```go
// MapT is a map with string keys and values.
type MapT map[string]interface{}

// Filter returns a new map with matched keys
func (m MapT) Filter(keys []string) map[string]interface{} {
  n := make(map[string]interface{}, len(m))
  for k, v := range m {
    if slices.Contains(keys, k) {
      n[k] = v
    }
  }
  return n
}
```

## 6. Docker build & run

놀라운 이미지 크기 `19MB` <span style="color:green;font-weight:700;font-size:40px"> 😲</span>

```bash
$ docker build -t gofiber-sqlite --no-cache .
$ docker image ls gofiber-sqlite
REPOSITORY       TAG       IMAGE ID       CREATED          SIZE
gofiber-sqlite   latest    a240a33bf435   10 minutes ago   19MB

# 환경변수 DB_URL, PORT 필요
$ docker run -lt --name gofiber-sqlite -p 3000:3000 -e DB_URL="file:memdb2?mode=memory" -e PORT=3000 gofiber-sqlite

$ docker rm -f $(docker ps -l -q)
```

## 9. Summary

[태주네이야기/Post](https://taejoone.jeju.onl/posts/2022-11-15-golang-tutorial-day5/)
