package mobiapi

import (
	"errors"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Grade struct {
	Description string
	Category    string
	Value       string
	Subject     string
	Semester    int
}

func (api *MobiAPI) GetGrades(semester int) (map[string][]Grade, error) {
	if !(semester == 1 || semester == 2) {
		return nil, errors.New("WrongSemester")
	}
	_, doc, err := api.request("GET", "oceny/?semestr="+strconv.Itoa(semester), "")
	if err != nil {
		return nil, err
	}

	subjects := map[string][]Grade{}
	var subject string
	var category string
	var description string
	var value string

	doc.Find("tr").Each(func(i int, s *goquery.Selection) {
		subject = s.Children().First().Text()
		s.Children().Last().Children().Each(func(ii int, ss *goquery.Selection) {
			if ss.Is("span") {
				category = strings.TrimSpace(strings.Trim(ss.Text(), ":"))
			} else if ss.Is("a") {
				value = ss.Children().First().Text()
				if ss.Children().Length() == 4 {
					ss.Children().First().Remove()
					ss.Children().First().Remove()
					description = ss.Children().First().Text()
					description = strings.TrimSpace(description[1 : len(description)-1])
				} else {
					description = "No description"
				}
				subjects[subject] = append(subjects[subject], Grade{
					Description: description,
					Category:    category,
					Value:       value,
					Subject:     subject,
					Semester:    semester,
				})
			}
		})
	})
	return subjects, nil
}
