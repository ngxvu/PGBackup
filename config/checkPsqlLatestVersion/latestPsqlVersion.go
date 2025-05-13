package checkPsqlLatestVersion

import (
	"backup/model"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
)

// GetAndParseServerVersion gets and parses the PostgreSQL server version
func GetAndParseServerVersion(db *sql.DB) (*model.PostgresqlVersion, error) {
	serverPsqlVersion, err := checkDatabaseVersion(db)
	if err != nil {
		return nil, fmt.Errorf("error checking database version: %v", err)
	}

	parsedVersion, err := parsePostgresqlVersion(*serverPsqlVersion)
	if err != nil {
		return nil, fmt.Errorf("error parsing version: %v", err)
	}

	return parsedVersion, nil
}

func checkDatabaseVersion(db *sql.DB) (*string, error) {
	var version string
	err := db.QueryRow("SHOW server_version").Scan(&version)
	if err != nil {
		return nil, fmt.Errorf("error querying database version: %v", err)
	}
	return &version, nil
}

func parsePostgresqlVersion(versionString string) (*model.PostgresqlVersion, error) {
	re := regexp.MustCompile(`(\d+\.\d+)`)
	matches := re.FindStringSubmatch(versionString)
	if len(matches) < 2 {
		return nil, fmt.Errorf("invalid version string: %s", versionString)
	}

	latestVersionWithMinor := matches[1]
	versionParts := strings.Split(latestVersionWithMinor, ".")
	versionMinor := versionParts[0]
	patchVersion := versionParts[1]

	return &model.PostgresqlVersion{
		LatestVersionWithMinor: &latestVersionWithMinor,
		VersionMinor:           &versionMinor,
		PatchVersion:           &patchVersion,
	}, nil
}

func CheckCurrentPostgresqlLatestVersion() (*model.PostgresqlVersion, error) {

	url := "https://www.postgresql.org/versions.json"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		err = fmt.Errorf("error creating request: %v", err)
		log.Println(err)
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		err = fmt.Errorf("error sending request: %v", err)
		log.Println(err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		err = fmt.Errorf("error reading response body: %v", err)
		log.Println(err)
	}

	var versions []model.CheckPostgresqlLatestVersionModel
	err = json.Unmarshal(body, &versions)

	if err != nil {
		err = fmt.Errorf("error unmarshalling response: %v", err)
		log.Println(err)
		return nil, err
	}

	var latestVersionWithMinor string
	var versionMinor string
	var psqlUrl *string

	for i := len(versions) - 1; i >= 0; i-- {
		if versions[i].Current == true {
			latestVersionWithMinor = versions[i].Major + "." + versions[i].LatestMinor
			versionMinor = versions[i].Major
			psqlUrl, err = getHrefLatestWindowsVersion(latestVersionWithMinor)
			if err != nil {
				err = fmt.Errorf("error getting href latest windows version: %v", err)
				log.Println(err)
				return nil, err
			}
			break
		}
	}

	postgresqlVersion := model.PostgresqlVersion{
		LatestVersionWithMinor: &latestVersionWithMinor,
		VersionMinor:           &versionMinor,
		PsqlUrl:                psqlUrl,
	}

	return &postgresqlVersion, nil
}

func getHrefLatestWindowsVersion(pgLatestVersion string) (*string, error) {

	url := model.PG_LATEST_VERSION_DOWNLOADS_URL

	// Make HTTP request
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Failed to fetch URL: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Fatalf("Non-200 HTTP status: %d", resp.StatusCode)
	}

	// Parse the HTML
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatalf("Failed to parse HTML: %v", err)
	}

	var link string
	// Find the target tbody with the desired class
	doc.Find("tbody.border-y.border-opacity-100.border-white").Each(func(i int, tbody *goquery.Selection) {
		// Iterate over each row
		tbody.Find("tr.border-y.border-white").Each(func(j int, tr *goquery.Selection) {
			// Find the version cell
			version := tr.Find("td.py-2.text-center.font-family-table-body").Text()
			if version == pgLatestVersion {
				// Find the 4th `text-center py-4` cell and extract the href
				link = tr.Find("td.text-center.py-4").Eq(3).Find("a").AttrOr("href", "")

			}
		})
	})
	return &link, nil
}
