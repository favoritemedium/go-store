package store

import (
  "testing"
  "os"
  "fmt"
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
)

// connect to a mysql db or die
func mustConnectMysql() {
  if db != nil { return } // already connected
  dsn := os.Getenv("MYSQL_TEST_DSN")
  if dsn == "" {
    fmt.Println("Please set envar MYSQL_TEST_DSN.")
    os.Exit(1)
  }
  d, err := sql.Open("mysql", dsn)
  if err != nil {
    err = d.Ping()
  }
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
  Init(d)
}

var allTables = []string{
  UserTable,
//  SignedinUserTable,
//  SigninHistoryTable,
//  EmailVerifyTable,
}

func mustTruncateTables() {
  mustConnectMysql()
  fmt.Println("Truncating all tables...")
  for _, table := range allTables {
    _, err := db.Exec("TRUNCATE TABLE " + table)
    if err != nil {
      fmt.Println(err)
      os.Exit(1)
    }
  }
}

var unsafeSuper = AuthUser{isSuperUser: true}

func TestCreateUser(t *testing.T) {
  mustTruncateTables()
  u := User{
    Email: "email@example.com",
    FullName: "Myfull Name",
    NameToUse: "Honey",
    IsActive: true,
  }
  err := UserCreate(unsafeSuper, &u)
  if err != nil {
    t.Fatal(err)
  }
  fmt.Printf("%#v\n", u)

  for z := range UserRead(1) {
    fmt.Printf("%#v\n", z)
  }

  for z := range UserRead(2) {
    fmt.Printf("%#v\n", z)
  }
}

