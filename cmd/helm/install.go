package main

import (
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"

	"github.com/kubernetes/helm/pkg/helm"
	"github.com/kubernetes/helm/pkg/proto/hapi/release"
	"github.com/kubernetes/helm/pkg/timeconv"
)

const installDesc = `
This command installs a chart archive.

The install argument must be either a relative
path to a chart directory or the name of a
chart in the current working directory.
`

// install flags & args
var (
	// installArg is the name or relative path of the chart to install
	installArg string
	// installDryRun performs a dry-run install
	installDryRun bool
	// installValues is the filename of supplied values.
	installValues string
)

var installCmd = &cobra.Command{
	Use:   "install [CHART]",
	Short: "install a chart archive.",
	Long:  installDesc,
	RunE:  runInstall,
}

func init() {
	f := installCmd.Flags()
	f.StringVar(&tillerHost, "host", "", "address of tiller server (default \":44134\")")
	f.StringVarP(&installValues, "values", "f", "", "path to a values TOML file")
	f.BoolVar(&installDryRun, "dry-run", false, "simulate an install")

	RootCommand.AddCommand(installCmd)
}

func runInstall(cmd *cobra.Command, args []string) error {
	if err := checkArgsLength(1, len(args), "chart name"); err != nil {
		return err
	}
	installArg = args[0]

	rawVals, err := vals()
	if err != nil {
		return err
	}

	res, err := helm.InstallRelease(rawVals, installArg, installDryRun)
	if err != nil {
		return prettyError(err)
	}

	printRelease(res.GetRelease())

	return nil
}

func vals() ([]byte, error) {
	if installValues == "" {
		return []byte{}, nil
	}
	return ioutil.ReadFile(installValues)
}

func printRelease(rel *release.Release) {
	if rel == nil {
		return
	}
	if flagVerbose {
		fmt.Printf("NAME:   %s\n", rel.Name)
		fmt.Printf("INFO:   %s %s\n", timeconv.String(rel.Info.LastDeployed), rel.Info.Status)
		fmt.Printf("CHART:  %s %s\n", rel.Chart.Metadata.Name, rel.Chart.Metadata.Version)
		fmt.Printf("MANIFEST: %s\n", rel.Manifest)
	} else {
		fmt.Println(rel.Name)
	}
}