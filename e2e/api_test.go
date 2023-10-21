package e2e

import (
	"context"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"testing"

	"github.com/cucumber/godog"
)

const appURL = "http://localhost:9080"

type apiFeature struct {
	resp *http.Response
	jar  *cookiejar.Jar
}

func (a *apiFeature) resetResponse(*godog.Scenario) {
	a.resp = nil
	a.jar = nil
}

func (a *apiFeature) iMakeAGETRequestTo(url string) error {
	resp, err := http.Get(fmt.Sprintf("%s%s", appURL, url))
	a.resp = resp
	return err
}

func (a *apiFeature) iMakeAPOSTRequestToWithTheFollowingFormData(endpoint string, table *godog.Table) error {
	if len(table.Rows) <= 1 {
		return fmt.Errorf("expected at least 2 row in table, but actual is: %d", len(table.Rows))
	}

	tableData := tableToDataMap(table)

	form := url.Values{}
	for key, value := range tableData {
		form.Add(key, value)
	}

	// Make the POST request
	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", appURL, endpoint), strings.NewReader(form.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	jar, err := cookiejar.New(nil)
	if err != nil {
		return fmt.Errorf("failed to create cookie jar: %w", err)
	}
	client := &http.Client{
		Jar: jar,
	}
	resp, err := client.Do(req)
	a.resp = resp
	a.jar = jar
	return err
}

func (a *apiFeature) theResponseCodeShouldBe(code int) error {
	if a.resp.StatusCode != code {
		return fmt.Errorf("expected response code to be: %d, but actual is: %d", code, a.resp.StatusCode)
	}

	return nil
}

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t, // Testing instance that will run subtests.
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	api := &apiFeature{}

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		api.resetResponse(sc)
		return ctx, nil
	})

	ctx.Step(`^I make a GET request to "([^"]*)"$`, api.iMakeAGETRequestTo)
	ctx.Step(`^I make a POST request to "([^"]*)" with the following form data:$`, api.iMakeAPOSTRequestToWithTheFollowingFormData)
	ctx.Step(`^the response code should be (\d+)$`, api.theResponseCodeShouldBe)
}

func tableToDataMap(table *godog.Table) map[string]string {
	dataMap := make(map[string]string, len(table.Rows)-1)
	for _, row := range table.Rows[1:] {
		for i, cell := range row.Cells {
			dataMap[table.Rows[0].Cells[i].Value] = cell.Value
		}
	}
	return dataMap
}
