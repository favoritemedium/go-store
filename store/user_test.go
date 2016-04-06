package store

import (
  "testing"
  "os"
  "fmt"
  "time"
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
  Catch(err)
  Catch(d.Ping())
  Init(d)
}

var allTables = []string{
  UserTable,
  UserSigninsTable,
  UserTokensTable,
  UserVerifyTable,
}

func mustTruncateTables() {
  mustConnectMysql()
  fmt.Println("Truncating all tables...")
  _, err := db.Exec("SET FOREIGN_KEY_CHECKS=0")
  Catch(err)
  for _, table := range allTables {
    _, err = db.Exec("TRUNCATE TABLE " + table)
    Catch(err)
  }
  _, err = db.Exec("SET FOREIGN_KEY_CHECKS=1")
  Catch(err)
}

var unsafeSuper = AuthUser{roles: SuperRole}
var unsafeAdmin = AuthUser{roles: AdminRole}

func isAlmostNow(when int64) bool {
  now := time.Now().Unix()
  return when >= now-2 && when <= now+2
}

func TestCreateUser(t *testing.T) {
  mustTruncateTables()

  u := User{
    Email: "email@example.com",
    FullName: "Myfull Name",
    NameToUse: "Honey",
    IsActive: true,
    Roles: AdminRole,
  }
  if err := UserCreate(unsafeSuper, &u); err != nil {
    t.Fatal(err)
  }

  // should fail the second time as it's a duplicate
  err := UserCreate(unsafeSuper, &u)
  if err == nil {
    t.Fatal("UserCreate: duplicate email failed to trigger error.")
  }
  if err != ErrDuplicateEmail {
    t.Fatal(err)
  }

  // test reading the user back
  z, err := UserRead(1)
  if err != nil {
    t.Fatal(err)
  }
  if z.Email != u.Email || z.FullName != u.FullName || z.NameToUse != u.NameToUse || !z.IsActive ||
      !isAlmostNow(z.CreatedAt) || !isAlmostNow(z.UpdatedAt) || !isAlmostNow(z.ActiveAt) {
    t.Fatalf("UserCreate: expected %#v, got %#v", u, z)
  }
  if z.Roles != 0 {
    t.Fatalf("UserCreate: non-admin was able to create admin user")
  }

  // test missing email
  u0 := u
  u0.Email = ""
  err = UserCreate(unsafeSuper, &u0)
  if err != ErrEmailRequired {
    t.Fatalf("UserCreate: expected %#v, got %#v", ErrEmailRequired, err)
  }

  // test missing fullname
  u1 := u
  u1.FullName = ""
  err = UserCreate(unsafeSuper, &u1)
  if err != ErrFullNameRequired {
    t.Fatalf("UserCreate: expected %#v, got %#v", ErrFullNameRequired, err)
  }

  // test missing nametouse
  u2 := u
  u2.NameToUse = ""
  err = UserCreate(unsafeSuper, &u2)
  if err != ErrNameToUseRequired {
    t.Fatalf("UserCreate: expected %#v, got %#v", ErrNameToUseRequired, err)
  }

  // make sure non-super can't create superuser
  u3 := u
  u3.Email = "email@example.co"
  u3.Roles = SuperRole
  err = UserCreate(unsafeAdmin, &u3)
  if err != nil {
    t.Fatal(err)
  }
  z, err = UserRead(2)
  if err != nil {
    t.Fatal(err)
  }
  if z.Roles != 0 {
    t.Fatalf("UserCreate: non-super was able to create superuser")
  }
}

