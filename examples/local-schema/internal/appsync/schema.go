package appsync

//go:generate go run github.com/jkrajniak/graphql-codegen-go -schema ../../schema.graphql -out models.go -entities Entity1

type RawEntity struct {
	ID string
}

func CopyEntity(e Entity1) RawEntity {
	return RawEntity{ID: e.Id}
}