package network

import (
	"fmt"
	"mini_docker/internal/util"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"
)

func LS() {
	files, err := os.ReadDir(stroageRootDir)
	if err != nil {
		util.LogAndExit("failed to read network directory", err)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.Debug)
	fmt.Fprintln(w, "NAME\tSUBNET\tGATEWAY\tDRIVER")
	for _, file := range files {
		netWorkName := strings.TrimSuffix(file.Name(), ".json")
		n, err := getNetwork(netWorkName)
		if err != nil {
			os.RemoveAll(filepath.Join(stroageRootDir, file.Name()))
			continue
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", n.Name, n.Subnet, n.Gateway, n.Driver)
	}
	w.Flush()
}
