package storage

type GofferStorage struct {
	DbStorage DBStorage
}

func NewGofferStorage(dataBaseDsn string) *GofferStorage {
	return &GofferStorage{
		DbStorage: newDatabase(dataBaseDsn),
	}
}
