package storage

type Stats struct {
	TotalInserted int // сколько всего добавили
	UniqueCount   int // уникальные строки
}

type Storage interface {
	// Add добавляет строку (с учётом логики хранилища)
	Add(s string)

	// GetAll возвращает все добавленные строки (как они лежат)
	GetAll() []string

	// Stats возвращает информацию о хранилище
	Stats() Stats

	// Reset очищает хранилище
	Reset()
}
