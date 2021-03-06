package cmd

import (
	"os"
	"text/tabwriter"

	"github.com/k14s/starlark-go/starlark"
	"github.com/sap/kubernetes-deployment-orchestrator/pkg/k8s"
	"github.com/sap/kubernetes-deployment-orchestrator/pkg/kdo"

	"github.com/spf13/cobra"
)

var listOptions = &kdo.RepoListOptions{}
var listK8sArgs = &k8s.Configs{}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list kdo charts",
	Long:  ``,
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		k8s, err := newK8s(listK8sArgs.Merge())
		if err != nil {
			exit(err)
		}
		exit(list(k8s, listOptions))
	},
}

func list(k k8s.K8s, listOptions *kdo.RepoListOptions) error {
	repo, err := repo()
	if err != nil {
		return err
	}
	thread := &starlark.Thread{Name: "main", Load: rootExecuteOptions.load}
	charts, err := repo.List(thread, k, listOptions)
	if err != nil {
		return err
	}
	writer := tabwriter.NewWriter(os.Stdout, 3, 4, 1, ' ', 0)
	defer writer.Flush()
	writer.Write([]byte("GENUS\tNAMESPACE\tVERSION\n"))
	for _, c := range charts {
		writer.Write([]byte(c.GetGenus() + "\t" + c.GetNamespace() + "\t" + c.GetVersion().String() + "\n"))
	}
	return nil
}

func init() {
	listOptions.AddFlags(listCmd.Flags())
	listK8sArgs.AddFlags(listCmd.Flags())
}
