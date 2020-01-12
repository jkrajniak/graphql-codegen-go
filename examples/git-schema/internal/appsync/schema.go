package appsync

//go:generate go run ../../../../. -schema git:// -out models.go -entities Entity1

type RawEntity struct {
	ID string
}

func CopyEntity(e Entity1) RawEntity {
	return RawEntity{ID: e.Id}
}