package mysql_test

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	root "xiaoyun/pkg"
	"xiaoyun/pkg/mysql"

	"github.com/DavidHuie/gomigrate"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

// Migrate 数据库迁移对象
type Migrate struct {
	log         *root.Log
	migratiions string
	mockSession *MockSession
}

// 无升级文件不返回错误
func TestMigrate_NoFile(t *testing.T) {

	migrate, err := NewMigrate()
	if err != nil {
		t.Error(err)
	}
	mockSession := migrate.mockSession

	defer mockSession.Close()
	defer os.RemoveAll(migrate.migratiions)

	rows := sqlmock.NewRows([]string{"table_name"}).AddRow("gomigrate")
	mock := mockSession.mock
	mock.ExpectQuery("SELECT (.+) FROM information_schema.tables").WillReturnRows(rows)

	m, err := mysql.NewMigrate(migrate.mockSession.mockDB, migrate.migratiions, migrate.log.Logger)
	if err != nil {
		t.Error(err)
	}

	if err := m.Up(); err != nil {
		t.Error(err)
	}

}

// 文件名不符合{id}_{name}_{up_or_down}.sql需通过
func TestMigrate_FileInvalid(t *testing.T) {

	migrate, err := NewMigrate()
	mockSession := migrate.mockSession
	defer mockSession.Close()
	defer os.RemoveAll(migrate.migratiions)

	f, err := ioutil.TempFile(migrate.migratiions, "")
	if err != nil {
		t.Error(err)
	}
	defer f.Close()

	rows := sqlmock.NewRows([]string{"table_name"}).AddRow("gomigrate")
	mock := mockSession.mock
	mock.ExpectQuery("SELECT (.+) FROM information_schema.tables").WillReturnRows(rows)

	m, err := mysql.NewMigrate(mockSession.mockDB, migrate.migratiions, migrate.log.Logger)
	if err != nil {
		t.Error(err)
		return
	}

	if err := m.Up(); err != nil {
		t.Error(err)
	}

}

// 存在不完整的up和down需报错
func TestMigrate_PairInvalid(t *testing.T) {

	migrate, err := NewMigrate()
	if err != nil {
		t.Error(err)
	}
	mockSession := migrate.mockSession
	defer mockSession.Close()
	defer os.RemoveAll(migrate.migratiions)

	f_up, err := ioutil.TempFile(migrate.migratiions, "3_test_up.sql")
	defer f_up.Close()

	mock := mockSession.mock
	rows := sqlmock.NewRows([]string{"table_name"}).AddRow("gomigrate")
	mock.ExpectQuery("SELECT (.+) FROM information_schema.tables").WillReturnRows(rows)

	_, err = mysql.NewMigrate(mockSession.mockDB, migrate.migratiions, migrate.log.Logger)
	if err != gomigrate.InvalidMigrationPair {
		t.Error(err)
		return
	}
}

// 文件内容为空需要通过
func TestMigrate_NilContent(t *testing.T) {

	migrate, err := NewMigrate()
	if err != nil {
		t.Error(err)
	}
	mockSession := migrate.mockSession
	defer mockSession.Close()
	defer os.RemoveAll(migrate.migratiions)

	f_up, err := ioutil.TempFile(migrate.migratiions, "3_test_up.sql")
	defer f_up.Close()
	defer os.Remove(f_up.Name())

	f_down, err := ioutil.TempFile(migrate.migratiions, "3_test_down.sql")
	defer f_down.Close()

	mock := mockSession.mock
	rows_table := sqlmock.NewRows([]string{"table_name"}).AddRow("gomigrate")
	rows_migration := sqlmock.NewRows([]string{"migrate_id"}).AddRow("1")

	mock.ExpectQuery("SELECT (.+) FROM information_schema.tables").WillReturnRows(rows_table)
	mock.ExpectQuery("SELECT (.+) FROM  gomigrate WHERE").WillReturnRows(rows_migration)

	m, err := mysql.NewMigrate(mockSession.mockDB, migrate.migratiions, migrate.log.Logger)
	if err != nil {
		t.Error(err)
	}

	if err = m.Up(); err != nil {
		t.Error(err)
	}
}

func TestMigrate_SqlInvalid(t *testing.T) {

	migrate, err := NewMigrate()
	if err != nil {
		t.Error(err)
	}
	mockSession := migrate.mockSession
	defer mockSession.Close()
	defer os.RemoveAll(migrate.migratiions)

	f_up, err := ioutil.TempFile(migrate.migratiions, "3_test_up.sql")
	if err != nil {
		t.Error(err)
	}
	f_up.WriteString("CREATE TABLE tests;")
	f_up.Close()

	f_down, err := ioutil.TempFile(migrate.migratiions, "3_test_down.sql")
	if err != nil {
		t.Error(err)
	}
	f_down.WriteString("DROP TABLE tests;")
	f_down.Close()

	mock := mockSession.mock
	rows_table := sqlmock.NewRows([]string{"table_name"}).AddRow("gomigrate")

	mock.ExpectQuery("SELECT (.+) FROM information_schema.tables").WillReturnRows(rows_table)
	mock.ExpectQuery("SELECT (.+) FROM gomigrate WHERE").WillReturnError(sql.ErrNoRows)
	mock.ExpectBegin()
	mock.ExpectExec("CREATE TABLE tests").WillReturnError(errors.New("migrate_exec_error"))
	mock.ExpectRollback()

	m, err := mysql.NewMigrate(mockSession.mockDB, migrate.migratiions, migrate.log.Logger)
	if err != nil {
		t.Error(err)
	}

	if err := m.Up(); err.Error() != "migrate_exec_error" {
		t.Error(err)
	}
}

func TestMigrate_Up(t *testing.T) {

	migrate, err := NewMigrate()
	if err != nil {
		t.Error(err)
	}
	mockSession := migrate.mockSession
	defer mockSession.Close()
	defer os.RemoveAll(migrate.migratiions)

	for i := 1; i < 10; i++ {
		upName := fmt.Sprintf("%d_test_migrate_up.sql", i)
		downName := fmt.Sprintf("%d_test_migrate_down.sql", i)
		fUp, err := ioutil.TempFile(migrate.migratiions, upName)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		fUp.WriteString(fmt.Sprintf("CREATE TABLE tests_%d;", i))
		fUp.Close()

		fDown, err := ioutil.TempFile(migrate.migratiions, downName)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		fDown.WriteString(fmt.Sprintf("DROP TABLE tests_%d;", i))
		fDown.Close()

	}

	mock := mockSession.mock
	rows_table := sqlmock.NewRows([]string{"table_name"}).AddRow("gomigrate")
	result := sqlmock.NewResult(0, 0)

	mock.ExpectQuery("SELECT (.+) FROM information_schema.tables").WillReturnRows(rows_table)

	for i := 1; i < 10; i++ {
		mock.ExpectQuery("SELECT (.+) FROM gomigrate WHERE").WillReturnError(sql.ErrNoRows)
	}
	for i := 1; i < 10; i++ {
		mock.ExpectBegin()
		mock.ExpectExec(fmt.Sprintf("CREATE TABLE tests_%d", i)).WillReturnResult(result)
		mock.ExpectExec("INSERT INTO gomigrate").WillReturnResult(result)
		mock.ExpectCommit()
	}

	m, err := mysql.NewMigrate(mockSession.mockDB, migrate.migratiions, migrate.log.Logger)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if err := m.Up(); err != nil {
		t.Error(err)
		t.FailNow()
	}
}

func TestMigrate_Down(t *testing.T) {

	migrate, err := NewMigrate()
	if err != nil {
		t.Error(err)
	}
	mockSession := migrate.mockSession
	defer mockSession.Close()
	defer os.RemoveAll(migrate.migratiions)

	for i := 1; i < 10; i++ {
		upName := fmt.Sprintf("%d_test_migrate_up.sql", i)
		downName := fmt.Sprintf("%d_test_migrate_down.sql", i)
		fUp, err := ioutil.TempFile(migrate.migratiions, upName)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		fUp.WriteString(fmt.Sprintf("CREATE TABLE tests_%d;", i))
		fUp.Close()

		fDown, err := ioutil.TempFile(migrate.migratiions, downName)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		fDown.WriteString(fmt.Sprintf("DROP TABLE tests_%d;", i))
		fDown.Close()

	}

	mock := mockSession.mock
	rows_table := sqlmock.NewRows([]string{"table_name"}).AddRow("gomigrate")
	result := sqlmock.NewResult(0, 0)

	mock.ExpectQuery("SELECT (.+) FROM information_schema.tables").WillReturnRows(rows_table)

	for i := 1; i < 10; i++ {
		rows := sqlmock.NewRows([]string{"migration_id"}).AddRow(i)
		mock.ExpectQuery("SELECT (.+) FROM gomigrate WHERE").WillReturnRows(rows)
	}

	mock.ExpectBegin()
	mock.ExpectExec(fmt.Sprintf("DROP TABLE tests_%d", 9)).WillReturnResult(result)
	mock.ExpectExec("DELETE FROM gomigrate").WillReturnResult(result)
	mock.ExpectCommit()

	m, err := mysql.NewMigrate(mockSession.mockDB, migrate.migratiions, migrate.log.Logger)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if err := m.Down(); err != nil {
		t.Error(err)
		t.FailNow()
	}

}

func NewMigrate() (*Migrate, error) {

	mockSession, err := NewMockSession()
	if err != nil {
		return nil, err
	}

	tempDir, err := ioutil.TempDir("", "migrate")

	log := root.NewLogStdOut()

	m := &Migrate{
		mockSession: &mockSession,
		log:         log,
		migratiions: tempDir,
	}

	return m, nil

}
