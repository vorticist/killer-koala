module github.com/vorticist/killer-koala

go 1.13

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/go-stack/stack v1.8.0 // indirect
	github.com/golang/snappy v0.0.1 // indirect
	github.com/labstack/echo v3.3.10+incompatible
	github.com/labstack/gommon v0.3.0 // indirect
	github.com/vorticist/logger v0.0.0-00010101000000-000000000000
	github.com/xdg/scram v0.0.0-20180814205039-7eeb5667e42c // indirect
	github.com/xdg/stringprep v1.0.0 // indirect
	go.mongodb.org/mongo-driver v1.1.2
	golang.org/x/crypto v0.0.0-20191002192127-34f69633bfdc // indirect
	golang.org/x/sync v0.0.0-20190911185100-cd5d95a43a6e // indirect
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22
)

replace github.com/vorticist/logger => ../logger
