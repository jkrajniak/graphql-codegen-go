package appsync

//go:generate go run ../../../../. -schema ../../schema.graphql -out models.go

type RawEntity struct {
	ID string
}

func CopyEntity(e Entity1) RawEntity {
	return RawEntity{ID: e.Id}
}