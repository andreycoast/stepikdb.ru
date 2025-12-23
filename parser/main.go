package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

type Course struct {
	Title		string
	Description string
	Price		int
}

func main() {
	if err := runCourseParser(); err != nil {
		log.Fatal(err)
	}
}

func runCourseParser() error {
	url := "https://stepik.org/course/89381/promo"

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP error: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
	}

	course := parseCourse(doc)

	fmt.Println("Title:", course.Title)
	fmt.Println("Description:", course.Description)
	fmt.Println("Price:", course.Price)

	return nil
}

func parseCourse(doc *goquery.Document) Course {
	return Course {
		Title:       parseTitle(doc),
		Description: parseDescription(doc),
		Price:       parsePrice(doc),
	}
}

func parseTitle(doc *goquery.Document) string {
	if title, ok := doc.Find(`meta[property="ya:ovs:title"]`).Attr("content"); ok {
		return title
	}

	return "unknown title"
}

func parseDescription(doc *goquery.Document) string {
	if description, ok := doc.Find(`meta[name="description"]`).Attr("content"); ok {
		return description
	}

	return "unknown description"
}

func parsePrice(doc *goquery.Document) int {
	priceContainer := doc.Find(".course-promo-enrollment__price-container").First()
	if priceContainer.Length() == 0 {
		return 0
	}

	formatPrice := priceContainer.Find(".format-price")
	if formatPrice.Length() == 0 {
		return 0
	}

	var priceParts []string
	formatPrice.Find(`span[data-type="integer"]`).Each(func(i int, s *goquery.Selection) {
		if val, ok := s.Attr("data-value"); ok {
			priceParts = append(priceParts, val)
		}
	})

	priceStr := strings.Join(priceParts, "")
	price, err := strconv.Atoi(priceStr)
	if err != nil {
		return 0
	}

	return price
}