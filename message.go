package mobiapi

import (
	"errors"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type MessageInfo struct {
	Kind   string
	Title  string
	Author string
	ID     int
	Read   bool
}

type MessageContent struct {
	Info       MessageInfo
	Content    string
	RawContent string
	Downloads  map[string]string
}

// Scrapes and returns message IDs and titles from first or every subsequent page in the form of MessageInfo. Use GetMessageContent() with MessageInfo to read it.
func (api *MobiAPI) GetReadMessages(firstpage bool) ([]MessageInfo, error) {
	pages := 1
	messages := []MessageInfo{}
	for i := 1; i <= pages; i++ {
		resp, doc, err := api.request("GET", "wiadomosci/?sortuj_wg=otrzymano&sortuj_typ=desc&odebrane="+strconv.Itoa(i), "")
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		if !firstpage && i == 1 {
			pages = doc.Find(".stronnicowanie").Children().Length()
		}
		doc.Find(".podswietl").Each(func(mi int, s *goquery.Selection) {
			message := MessageInfo{Read: s.Find("td span").HasClass("wiadomosc_przeczytana")}
			sid, _ := s.Attr("rel")
			iid, _ := strconv.Atoi(sid)
			message.ID = iid
			html, _ := s.Children().Html()
			message.Title = strings.TrimSpace(html)
			s.Children().First().Remove()
			s.Children().First().Remove()
			s.Children().First().Remove()
			html, _ = s.Children().Html()
			message.Author = strings.ReplaceAll(strings.ReplaceAll(html, "<small>", ""), "</small>", "")
			messages = append(messages, message)
		})
	}
	if len(messages) > 0 {
		return messages, nil
	}
	return nil, errors.New("Unprocessed")
}

// Searches messages using MobiDziennik's built-in search feature.
func (api *MobiAPI) SearchMessages(phrase string) ([]MessageInfo, error) {
	resp, doc, err := api.request("GET", "dziennik/wyszukiwarkawiadomosci?q="+phrase, "")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	messages := []MessageInfo{}
	doc.Find(".podswietl").Each(func(i int, s *goquery.Selection) {
		if attr, exists := s.Attr("rel"); exists && strings.Contains(attr, "wiadodebrana") {
			id, _ := strconv.Atoi(strings.ReplaceAll(attr, "wiadodebrana?id=", ""))
			messages = append(messages, MessageInfo{
				Title:  s.Find("td div.ellipsis").Text(),
				Author: s.Find("td div.autoTooltip").Text(),
				ID:     id,
				Read:   false,
			})
		}
	})
	if len(messages) > 0 {
		return messages, nil
	}

	return nil, errors.New("Unprocessed")
}

// Read Received Message from MessageInfo into MessageContent.
func (api *MobiAPI) GetMessageContent(message MessageInfo) (MessageContent, error) {
	messagecontent := MessageContent{Info: message}
	resp, doc, err := api.request("GET", "dziennik/wiadodebrana/?id="+strconv.Itoa(message.ID), "")
	if err != nil {
		return messagecontent, err
	}
	if resp.StatusCode != 200 {
		return messagecontent, errors.New("NotFound")
	}

	contents := doc.Find(".wiadomosc_tresc").Children().First()
	messagecontent.RawContent, err = contents.Html()
	if err != nil {
		return messagecontent, err
	}
	messagecontent.Content = strings.ReplaceAll(strings.ReplaceAll(contents.Text(), "<p>", ""), "</p>", "\n")

	if doc.Find("#zalaczniki").Length() == 1 {
		messagecontent.Downloads = map[string]string{}
		doc.Find("#zalaczniki li a").Each(func(ii int, s *goquery.Selection) {
			s.Children().Remove()
			if url, exists := s.Attr("href"); exists {
				messagecontent.Downloads[strings.TrimSpace(s.Text())] = url
			}
		})
	}

	return messagecontent, nil
}
