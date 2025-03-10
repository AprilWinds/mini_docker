package network

import (
	"encoding/json"
	"fmt"
	"mini_docker/internal/util"
	"os"
	"path/filepath"
	"text/tabwriter"
)

func LS() {
	files, err := os.ReadDir(stroageRootDir)
	if err != nil {
		util.LogAndExit("failed to read network directory", err)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.Debug)
	fmt.Fprintln(w, "ID\tNAME\tSUBNET\tGATEWAY\tDRIVER")
	for _, file := range files {
		var n Network
		f, err := os.Open(filepath.Join(stroageRootDir, file.Name()))
		if err != nil {
			os.RemoveAll(filepath.Join(stroageRootDir, file.Name()))
			continue
		}
		defer f.Close()
		json.NewDecoder(f).Decode(&n)
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", file.Name(), n.Name, n.Subnet, n.Gateway, n.Driver)
	}
	w.Flush()
}
