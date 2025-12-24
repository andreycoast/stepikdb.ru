package parser

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"stepikdb.ru/internal/models"
)

func ParseCourse(url string) (models.Course, error) {
	resp, err := http.Get(url)
	if err != nil {
		return models.Course{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return models.Course{}, fmt.Errorf("HTTP error: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return models.Course{}, err
	}

	course := models.Course{
		CourseID:    parseCourseID(url),
		Title:       parseTitle(doc),
		Description: parseDescription(doc),
		Price:       parsePrice(doc),
		URL:         url,
	}

	return course, nil
}

func parseCourseID(url string) int {
	courseIDRegex, _ := regexp.Compile(`https://stepik.org/course/(\\d+)`)

	matches := courseIDRegex.FindStringSubmatch(url)
	if len(matches) < 2 {
		return 0
	}

	id, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0
	}

	return id
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
