package web

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func decodeGetUserWithIDRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	v := mux.Vars(r)
	userIDParam, ok := v["id"]
	if !ok {
		return nil, errors.New("user ID was not provided")
	}
	return userIDParam, nil
}

func decodeSearchUsersRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	filterRequest := SearchUserFilter{
		Page:     1,
		PageSize: 10,
	}

	filters := r.URL.Query()

	if v, ok := filters["city"]; ok {
		filterRequest.City = v[0]
	}

	thereAreSkills := true
	err := r.ParseForm()
	if err != nil {
		log.Println("level", "ERROR", "no skills to look for", "error", err)
		thereAreSkills = false
	}
	if thereAreSkills {
		skills := r.Form["skills"]
		skillsToSearch := make([]string, 0)
		for _, v := range skills {
			skill := v
			skillsToSearch = append(skillsToSearch, skill)
		}
		filterRequest.Skills = skills
	}

	if v, ok := filters["page"]; ok {
		page, err := strconv.Atoi(v[0])
		if err != nil {
			log.Println("level", "ERROR", "invalid page parameter, it must be an integer", "error", err)
			page = 1
		}
		filterRequest.Page = page
	}
	if v, ok := filters["pagesize"]; ok {
		pageSize, err := strconv.Atoi(v[0])
		if err != nil {
			log.Println("level", "ERROR", "invalid page size parameter, it must be an integer", "error", err)
			pageSize = 10
		}
		filterRequest.PageSize = pageSize
	}

	filter := filterRequest.toSearchUserFilter()

	return filter, nil
}

func decodeCreateUserRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	log.Println("level", "DEBUG", "msg", "decoding new user request")
	var req NewUser
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &req)
	if err != nil {
		log.Println("level", "ERROR", "new user request could not be decoded. Request: %q because of: %s", string(body), err.Error())
		return nil, err
	}

	log.Println("level", "DEBUG", "msg", "user request was decoded", "request", req)

	domainUser := req.toUser()

	return domainUser, nil
}

func decodeUpdateUserRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	log.Println("level", "DEBUG", "msg", "decoding update user request")
	var req UpdateUser
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &req)
	if err != nil {
		log.Println("level", "ERROR", "update user request could not be decoded. Request: %q because of: %s", string(body), err.Error())
		return nil, err
	}

	log.Println("level", "DEBUG", "msg", "user request was decoded", "request", req)

	domainUser := req.toUser()

	return domainUser, nil
}
