package infra

//go:generate go run github.com/vektra/mockery/v2@v2.53.2 --name=Database --output=../../../mocks --filename=mock_database.go --with-expecter

type Database interface {
	// TODO: Add CRUD methods
}