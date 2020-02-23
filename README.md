[![Go Report Card](https://goreportcard.com/badge/github.com/jkrajniak/graphql-codegen-go)](https://goreportcard.com/report/github.com/jkrajniak/graphql-codegen-go)

# graphql-codegen-go
Generate Go structs from your GraphQL schema.

This [code generator](https://blog.golang.org/generate) helps you derive Go structures directly from [GraphQL](https://graphql.org/) schema. The schema
can be located either locally or can be fetched from GIT repository.

Install it using `go get`:

```bash
$ go get -u github.com/jkrajniak/graphql-codegen-go
```
## Quick start

Simply, define the GQL schema
```graphql
type Person {
  name: String!
  age: Int!
  weight: Int
  likes: [String]
  donts: [String!]
}
```
and save it, e.g. in `schema.gql` file. Then run the code generator

```bash
$ graphql-codegen-go -schema schema.gql -packageName pkg -out models.go
```
As a result, you will get a `models.go` file with the following Go code

```go
// Code generated by go generate; DO NOT EDIT.
// This file was generated from GraphQL schema schema.gql

package pkg

type Person struct {
	Name   string    `json:"name"`
	Age    int64     `json:"age"`
	Weight *int64    `json:"weight"`
	Likes  []*string `json:"likes"`
	Donts  []string  `json:"donts"`
}
```

Notice that not required (`weight`) fields are converted to the pointers. The `packageName` option is optional. The generator
will try to derive the package name from the current running path of the code, or if run by `go:generate` from `$GOPACKAGE` env variable.

### GIT

The schema does not have to be located locally. The program supports also Git repositories.
Let's assume that you the `schema.gql` file is placed in `github.com/orange/repo1` repository, inside a `deployment` directory.
Then, to create the structures you can run the generator as follows

```bash
$ graphql-codegen-go -schema https://github.com/orange/repo1.git/deployment/schema.gql -packageName pkg -out models.go
```

or via ssh

```
$ graphql-codegen-go -schema git@github.com:orange/repo1.git/deployment/schema.gql -packageName pkg -out models.go
```

By default, the schema is pulled from the `HEAD`. To point a specific commit, you can place a commit hash after the file name, e.g.,

```
$ graphql-codegen-go -schema git@github.com:orange/repo1.git/deployment/schema.gql#a56351vc -packageName pkg -out models.go
```

Moreover, you can also point the specific branch or tag by using `@` sign

```
$ graphql-codegen-go -schema git@github.com:orange/repo1.git/deployment/schema.gql@branch -packageName pkg -out models.go
```

or

```
$ graphql-codegen-go -schema git@github.com:orange/repo1.git/deployment/schema.gql@tag1 -packageName pkg -out models.go
```

### Entities

By default generator will output structures for all of the entities found in the schema. To output only a subset of structures
you can use `-entities` option.

For example
```bash
$ graphql-codegen-go -schema schema.gql -packageName pkg -out models.go -entities Person
```
will create file `models.go` only with a single structure `Person`, and all related dependent structures and enums.


### YAML Config

Instead of command line parameters, the generator supports also a config file (`-config`). The example of the file can be found in `examples/config/config.yml`.

The structure of the YAML file is

```yaml

schema: A list of schema files (it will be combined into one schema before parsing)
generates: A key-value map, where key is the name of the output Go file
    <output-file>:
      config:
        packageName: A package name
        entities: A list of entities to be included into the output Go file
```

#### Example

```yaml
schema:
    - ./schema.graphql
    - ./types.graphql
    - https://github.com/jkrajniak/sc.git/schema1.gql
generates:
    internal/models.go:
      config:
        packageName: internal
        entities:
          - User
          - Person
    internal/abc/models.go:
      config:
        packageName: abc
        entities:
          - Action
```

After execution of the above config you will find two `.go` files (`internal/models.go` and `internal/abc/models.go`);
the first will contain structures `User` and `Person`, and the second `Action`.

The GQL schema will be read from two local files, from GIT repository.

### go:generate

The generator can work perfectly fine with the ```go:generate``` directive. The examples of how to include it can be found in `examples/` directory.


## Union

The [union type](https://graphql.org/learn/schema/#union-types) is supported in the following way. Let's consider a schema
```graphql
type Et {
  value: String
}

type Pt {
  label: String!
}

union Ett = Et | Pt

type Entity {
  id: String!
  etpt: Ett
}
```

The union type `Ett` is converted into a Go struct `Ett`
```go
type Et struct {
	Value *string `json:"value"`
}

type Pt struct {
	Label string `json:"label"`
}

type Ett struct {
	TypeName string `json:"__typeName"`
	Et
	Pt
}

type Entity struct {
	Id   string `json:"id"`
	Etpt *Ett   `json:"etpt"`
}
```
