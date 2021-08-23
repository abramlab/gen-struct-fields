# Gen-struct-fields
Gen-struct-fields is a tool to generate the names of the structure fields in various formats in Go code.

The generated code can be used, for example, to simplify work with the database and etc.

## Install
Gen-struct-fields can be installed as any other go command:

```
go get github.com/abramlab/gen-struct-fields
```
After that, the `gen-struct-fields` executable will be in "$GOPATH/bin" folder, and you can use it with `go generate`

## How to use

```
gen-struct-fields -struct=User -custom_name=user
```
or with go generate:
```
//go:generate gen-struct-fields -struct=User -custom_name=user
```

By default, for each structure will be generated a separate file, with name: <struct_name>_fields.go.

You can use ```output``` tag to specified custom output file name, then all data will be generated in it.

```
gen-struct-fields -struct=User -output=custom_generated_file.go
```


For multiple generation, separate tag struct value with coma, for example:
```
gen-struct-fields -struct=User,Car,Plane -tag=custom_tag -custom_name=user,,plane
```
will generate fields for ```User,Car,Plane``` structs with custom names for ```User``` and ```Plane```.

## Simple example

Initial data:
```
//go:generate gen-struct-fields -struct=User -tag=custom_tag -custom_name=user

type User struct {
	Name     string `custom_tag:"name"`
	Login    string `custom_tag:"username"`
	Password string `custom_tag:"-"`
	AuthType int    `custom_tag:"auth_type"`
}

```
Output(generated) data:
```
const UserName = "user"

var UserFields = struct {
	Name     string
	Login    string
	AuthType string
}{
	Name:     "name",
	Login:    "username",
	AuthType: "auth_type",
}

var UserFieldsArray = []string{
	UserFields.Name,
	UserFields.Login,
	UserFields.AuthType,
}
```
