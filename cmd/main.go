package main

import (
	"flag"
	"log"
	"os"
	"r01api/internal"
)

func main() {
	// Добавляем флаги для выбора режима работы
	auth := flag.Bool("auth", false, "Добавить TXT-запись")
	cleanup := flag.Bool("cleanup", false, "Удалить TXT-запись")
	flag.Parse()

	if !*auth && !*cleanup {
		log.Fatal("Укажи флаг --auth или --cleanup")
	}

	if *auth && *cleanup {
		log.Fatal("Нужно указать ровно один флаг: либо --auth, либо --cleanup")
	}

	// Читаем данные из переменной среды
	accessToken := os.Getenv("R01_ACCESS_TOKEN")
	domain := os.Getenv("CERTBOT_DOMAIN")
	validation := os.Getenv("CERTBOT_VALIDATION")

	if domain == "" || accessToken == "" {
		log.Fatal("CERTBOT_DOMAIN или R01_ACCESS_TOKEN не заданы")
	}

	// Собственно, выбираем режим работы
	if *auth {
		if validation == "" {
			log.Fatal("CERTBOT_VALIDATION не задан")
		}
		if err := internal.RunAuth(accessToken, domain, validation); err != nil {
			log.Fatal(err)
		}
	} else if *cleanup {
		if err := internal.RunCleanup(accessToken, domain); err != nil {
			log.Fatal(err)
		}
	}
}
