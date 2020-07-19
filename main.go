package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"gitub.com/BradErz/grafana-dashboard-importer/grafana"

	"github.com/sirupsen/logrus"
)

var (
	token    = flag.String("token", "", "grafana api token")
	url      = flag.String("url", "", "base url of the grafana server")
	dir      = flag.String("dir", ".", "directory containing only grafana json dashboards")
	folderID = flag.Int("folder-id", 0, "the folder id that you want to import to")
)

func loadFiles(dir string) (map[string][]byte, error) {
	objs, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("error reading directory")
	}

	res := map[string][]byte{}
	for _, obj := range objs {
		// completely ignores dirs
		if obj.IsDir() {
			continue
		}
		fileName := filepath.Join(dir, obj.Name())
		content, err := ioutil.ReadFile(fileName)
		if err != nil {
			return nil, fmt.Errorf("failed to read %s: %w", fileName, err)
		}
		res[fileName] = content
	}
	return res, nil
}

func main() {
	flag.Parse()
	logrus.Infof("importing dashboards from %s to %s folder %d", *dir, *url, *folderID)

	reqs, err := loadFiles(*dir)
	if err != nil {
		logrus.WithError(err).Fatal("error loading files")
	}

	grafanaCl := grafana.New(*token, *url, *folderID)
	if err := grafanaCl.CreateDashboards(reqs); err != nil {
		logrus.WithError(err).Fatal("failed to create dashboards")
	}
	logrus.Infof("successfully created %d dashboards", len(reqs))
}
