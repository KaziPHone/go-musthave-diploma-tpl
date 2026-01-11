package storage

type GofferStorage struct {
	DbStorage DbStorage
}

func NewGofferStorage(dataBaseDsn string) *GofferStorage {
	return &GofferStorage{
		DbStorage: newDatabase(dataBaseDsn),
	}
}
