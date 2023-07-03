package types

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/ms-henglu/armstrong/coverage"
)

type PassReport struct {
	Resources []Resource
}

type Resource struct {
	ApiPath string
	Type    string
	Address string
}

type CoverageReport struct {
	Coverages map[Resource]*coverage.Model
}

type DiffReport struct {
	Diffs []Diff
	Logs  []RequestTrace
}

type Diff struct {
	Id      string
	Type    string
	Address string
	Change  Change
}

type Change struct {
	Before string
	After  string
}

type ErrorReport struct {
	Errors []Error
	Logs   []RequestTrace
}

type Error struct {
	Id      string
	Type    string
	Label   string
	Message string
}

func (c *CoverageReport) AddCoverageFromState(resourceId, resourceType, address string, jsonBody map[string]interface{}) error {
	var apiPath, modelName, modelSwaggerPath *string
	var err error

	apiVersion := strings.Split(resourceType, "@")[1]
	if !regexp.MustCompile(`^[0-9]{4}-[0-9]{2}-[0-9]{2}$`).MatchString(apiVersion) {
		return fmt.Errorf("could not parse apiVersion from resourceType: %s", resourceType)
	}

	apiPath, modelName, modelSwaggerPath, err = coverage.GetModelInfoFromIndex(resourceId, apiVersion)
	if err != nil {
		return fmt.Errorf("error find the path for %s from index:%s", resourceId, err)

	}

	log.Printf("[INFO] matched API path:%s modelSwawggerPath:%s\n", *apiPath, *modelSwaggerPath)

	resource := Resource{
		ApiPath: *apiPath,
		Type:    resourceType,
		Address: address,
	}

	if _, ok := c.Coverages[resource]; !ok {
		expanded, err := coverage.Expand(*modelName, *modelSwaggerPath)
		if err != nil {
			return fmt.Errorf("error expand model %s property:%s", *modelName, err)
		}

		c.Coverages[resource] = expanded
	}
	c.Coverages[resource].MarkCovered(jsonBody)
	c.Coverages[resource].CountCoverage()

	return nil
}
