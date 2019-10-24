package main

import (
	"os"
)

func main() {
	config := AppConfig{
		MongoDBUrl:  os.Getenv("MONGO_URL"),
		MongoDBName: os.Getenv("MONGO_DB_NAME"),
		PrivateKey:  os.Getenv("PRIVATE_KEY"),
		PublicKey:   os.Getenv("PUBLIC_KEY"),
		Port:        os.Getenv("PORT"),
	}
	app := NewAppWithConfig(&config)

	app.Serve()
}
