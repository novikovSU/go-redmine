package redmine

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	//"fmt"
)

const (
	maxLimit = 100
	minLimit = 10
)

type timeEntriesResult struct {
	TimeEntries []TimeEntry `json:"time_entries"`
	TotalCount int `json:"total_count"`
}

type timeEntryResult struct {
	TimeEntry TimeEntry `json:"time_entry"`
}

type timeEntryRequest struct {
	TimeEntry TimeEntry `json:"time_entry"`
}

// TimeEntry -- main exported type
type TimeEntry struct {
	ID          int            `json:"id"`
	Project      IdName         `json:"project"`
	Issue        Id             `json:"issue"`
	User         IdName         `json:"user"`
	Activity     IdName         `json:"activity"`
	Hours        float32        `json:"hours"`
	Comments     string         `json:"comments"`
	SpentOn      string         `json:"spent_on"`
	CreatedOn    string         `json:"created_on"`
	UpdatedOn    string         `json:"updated_on"`
	CustomFields []*CustomField `json:"custom_fields,omitempty"`
}

// TimeEntriesWithFilter send query and return parsed result
func (c *Client) TimeEntriesWithFilter(filter Filter) ([]TimeEntry, error) {
	var result []TimeEntry
	
	var limit int
	limit, err := strconv.Atoi(filter.filters["limit"])
	if err != nil {
		limit = minLimit
	}

	var offset int
	offset, err = strconv.Atoi(filter.filters["offset"])
	if err != nil {
		offset = 0
	}

	for i := 0; i < limit/maxLimit; i++ {
		filter.filters["offset"] = strconv.Itoa(offset+maxLimit*i)
		filter.filters["limit"] = strconv.Itoa(maxLimit)

		uri, err := c.URLWithFilter("/time_entries.json", filter)
		if err != nil {
			return nil, err
		}
		req, err := http.NewRequest("GET", uri, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Add("X-Redmine-API-Key", c.apikey)
		res, err := c.Do(req)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		decoder := json.NewDecoder(res.Body)
		var r timeEntriesResult
		if res.StatusCode == 404 {
			return nil, errors.New("Not Found")
		}
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
		result = append(result, r.TimeEntries...)
	}

	return result, nil
}

// TimeEntries -- get time entries by project ID
func (c *Client) TimeEntries(projectID int) ([]TimeEntry, error) {
	res, err := c.Get(c.endpoint + "/projects/" + strconv.Itoa(projectID) + "/time_entries.json?key=" + c.apikey + c.getPaginationClause())
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	var r timeEntriesResult
	if res.StatusCode == 404 {
		return nil, errors.New("Not Found")
	}
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
	return r.TimeEntries, nil
}

// TimeEntry -- get single time entry by its ID
func (c *Client) TimeEntry(id int) (*TimeEntry, error) {
	res, err := c.Get(c.endpoint + "/time_entries/" + strconv.Itoa(id) + ".json?key=" + c.apikey)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	var r timeEntryResult
	if res.StatusCode == 404 {
		return nil, errors.New("Not Found")
	}
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
	return &r.TimeEntry, nil
}

// CreateTimeEntry -- no comments
func (c *Client) CreateTimeEntry(timeEntry TimeEntry) (*TimeEntry, error) {
	var ir timeEntryRequest
	ir.TimeEntry = timeEntry
	s, err := json.Marshal(ir)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", c.endpoint+"/time_entries.json?key="+c.apikey, strings.NewReader(string(s)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	var r timeEntryResult
	if res.StatusCode != 201 {
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
	return &r.TimeEntry, nil
}

// UpdateTimeEntry -- no comments
func (c *Client) UpdateTimeEntry(timeEntry TimeEntry) error {
	var ir timeEntryRequest
	ir.TimeEntry = timeEntry
	s, err := json.Marshal(ir)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("PUT", c.endpoint+"/time_entries/"+strconv.Itoa(timeEntry.ID)+".json?key="+c.apikey, strings.NewReader(string(s)))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := c.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		return errors.New("Not Found")
	}
	if res.StatusCode != 200 {
		decoder := json.NewDecoder(res.Body)
		var er errorsResult
		err = decoder.Decode(&er)
		if err == nil {
			err = errors.New(strings.Join(er.Errors, "\n"))
		}
	}
	if err != nil {
		return err
	}
	return err
}

// DeleteTimeEntry -- no comments
func (c *Client) DeleteTimeEntry(id int) error {
	req, err := http.NewRequest("DELETE", c.endpoint+"/time_entries/"+strconv.Itoa(id)+".json?key="+c.apikey, strings.NewReader(""))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := c.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		return errors.New("Not Found")
	}

	decoder := json.NewDecoder(res.Body)
	if res.StatusCode != 200 {
		var er errorsResult
		err = decoder.Decode(&er)
		if err == nil {
			err = errors.New(strings.Join(er.Errors, "\n"))
		}
	}
	return err
}
