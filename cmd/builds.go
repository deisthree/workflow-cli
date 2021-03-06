package cmd

import (
	yaml "gopkg.in/yaml.v2"

	"github.com/teamhephy/controller-sdk-go/builds"
)

// BuildsList lists an app's builds.
func (d *HephyCmd) BuildsList(appID string, results int) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	if results == defaultLimit {
		results = s.Limit
	}

	builds, count, err := builds.List(s.Client, appID, results)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Printf("=== %s Builds%s", appID, limitCount(len(builds), count))

	for _, build := range builds {
		d.Println(build.UUID, build.Created)
	}
	return nil
}

// BuildsCreate creates a build for an app.
func (d *HephyCmd) BuildsCreate(appID, image, procfile string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	procfileMap := make(map[string]string)

	if procfile != "" {
		if procfileMap, err = parseProcfile([]byte(procfile)); err != nil {
			return err
		}
	}

	d.Print("Creating build... ")
	quit := progress(d.WOut)
	_, err = builds.New(s.Client, appID, image, procfileMap)
	quit <- true
	<-quit
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Println("done")

	return nil
}

func parseProcfile(procfile []byte) (map[string]string, error) {
	procfileMap := make(map[string]string)
	return procfileMap, yaml.Unmarshal(procfile, &procfileMap)
}
