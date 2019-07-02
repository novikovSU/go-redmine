package redmine

import (
	"encoding/xml"
	"errors"
	"strings"
	"time"
)

// Author -- AAA
type Author struct {
	XMLName xml.Name `xml:"author"`
	Name    string   `xml:"name"`
	Email   string   `xml:"email,omitempty"`
}

// Content -- AAA
type Content struct {
	XMLName xml.Name `xml:"content"`
	Type    string   `xml:"type,attr"`
}

// Link -- AAA
type Link struct {
	XMLName xml.Name `xml:"link"`
	Rel     string   `xml:"rel,attr"`
	Href    string   `xml:"href,attr"`
}

// Activity -- AAA
type Activity struct {
	XMLName xml.Name  `xml:"entry"`
	Title   string    `xml:"title"`
	Link    Link      `xml:"link"`
	ID      string    `xml:"id"`
	Updated time.Time `xml:"updated"`
	Author  Author    `xml:"author"`
	Content Content   `xml:"content"`
}

// Generator -- AAA
type Generator struct {
	XMLName xml.Name `xml:"generator"`
	URI     string   `xml:"uri,attr"`
}

// Feed -- AAA
type Feed struct {
	XMLName    xml.Name   `xml:"feed"`
	Xmlns      string     `xml:"xmlns,attr"`
	Title      string     `xml:"title"`
	Links      []Link     `xml:"link"`
	ID         string     `xml:"id"`
	Icon       string     `xml:"icon"`
	Updated    time.Time  `xml:"updated"`
	Author     Author     `xml:"author"`
	Generator  Generator  `xml:"generator"`
	Activities []Activity `xml:"entry"`
}

type ActivityFilter struct {
	ProjectID int
	FromDate  time.Time
	ToDate    time.Time
}

func uniqueActivities(activities []Activity) []Activity {
	keys := make(map[string]bool)
	list := []Activity{}
	for _, entry := range activities {
		if _, value := keys[entry.ID]; !value {
			keys[entry.ID] = true
			list = append(list, entry)
		}
	}
	return list
}

// ActivityOf AAA
func (c *Client) ActivityOf(pID int) ([]Activity, error) {
	project, err := c.Project(pID)
	if err != nil {
		return nil, err
	}

	res, err := c.Get(c.endpoint + "/projects/" + project.Identifier + "/activity.atom?key=" + c.Atomkey)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	decoder := xml.NewDecoder(res.Body)

	var r Feed
	if res.StatusCode != 200 {
		var er errorsResult
		err = decoder.Decode(&er)
		if err == nil {
			err = errors.New(strings.Join(er.Errors, "\n"))
		}
	} else {
		err = decoder.Decode(&r)
	}
	if err != nil {
		return nil, err
	}

	return r.Activities, nil
}

// ActivityByFilter AAA
func (c *Client) ActivityByFilter(f *ActivityFilter) ([]Activity, error) {
	result := []Activity{}
	toDate := f.ToDate.AddDate(0, 0, 1)
	for i := f.FromDate; i.Before(toDate); i = i.AddDate(0, 0, 1) {
		activity, err := c.activityFromDate(f.ProjectID, i)
		if err != nil {
			return nil, err
		}
		result = append(result, activity...)
	}
	result = uniqueActivities(result)
	return result, nil
}

func (c *Client) activityFromDate(pID int, from time.Time) ([]Activity, error) {
	project, err := c.Project(pID)
	if err != nil {
		return nil, err
	}

	res, err := c.Get(c.endpoint + "/projects/" + project.Identifier + "/activity.atom?key=" + c.Atomkey + "&from=" + from.Format("2006-01-02"))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	decoder := xml.NewDecoder(res.Body)

	var r Feed
	if res.StatusCode != 200 {
		var er errorsResult
		err = decoder.Decode(&er)
		if err == nil {
			err = errors.New(strings.Join(er.Errors, "\n"))
		}
	} else {
		err = decoder.Decode(&r)
	}
	if err != nil {
		return nil, err
	}

	return r.Activities, nil
}
