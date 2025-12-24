package main

import (
	"log"

	"stepikdb.ru/internal/config"
	"stepikdb.ru/internal/parser"
	"stepikdb.ru/internal/storage/postgres"
)

func main() {
	url := "https://stepik.org/course/89381"

	course, err := parser.ParseCourse(url)
	if err != nil {
		log.Fatal(err)
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	db, err := postgres.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf(
		"Course ID: %s\nTitle: %s\nDescription: %s\nPrice: %d\nURL: %s\n",
		course.CourseID,
		course.Title,
		course.Description,
		course.Price,
		course.URL,
	)
}
