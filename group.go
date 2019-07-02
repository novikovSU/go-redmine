package redmine

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
)

type groupResult struct {
	Group Group `json:"group"`
}

type groupsResult struct {
	Groups []Group `json:"groups"`
}

// Group AAA
type Group IdName

// Groups AAA
func (c *Client) Groups() ([]Group, error) {
	res, err := c.Get(c.endpoint + "/groups.json?key=" + c.apikey + c.getPaginationClause())
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	var r groupsResult
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
	return r.Groups, nil
}

// Group AAA
func (c *Client) Group(id int) (*Group, error) {
	res, err := c.Get(c.endpoint + "/groups/" + strconv.Itoa(id) + ".json?key=" + c.apikey)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	var r groupResult
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
	return &r.Group, nil
}
