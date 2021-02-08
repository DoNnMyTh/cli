package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var kubernetesShowCmd = &cobra.Command{
	Use:     "show",
	Aliases: []string{"get", "inspect"},
	Example: `civo kubernetes show ID/HOSTNAME -o custom -f "ID: Code (DefaultUsername)"`,
	Args:    cobra.MinimumNArgs(1),
	Short:   "Show Kubernetes cluster",
	Long: `Show a specified Kubernetes cluster.
If you wish to use a custom format, the available fields are:

	* ID
	* Code
	* Name
	* Nodes
	* Size
	* CPUCores
	* RAMMegabytes
	* DiskGigabytes
	* Status
	* KubernetesVersion
	* APIEndPoint
	* MasterIP
	* DNSEntry`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if regionSet != "" {
			client.Region = regionSet
		}
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		kubernetesCluster, err := client.FindKubernetesCluster(args[0])
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriter()
		ow.StartLine()

		ow.AppendData("ID", kubernetesCluster.ID)
		ow.AppendData("Name", kubernetesCluster.Name)
		ow.AppendData("Nodes", strconv.Itoa(kubernetesCluster.NumTargetNode))
		ow.AppendData("Size", kubernetesCluster.TargetNodeSize)
		ow.AppendData("Status", kubernetesCluster.Status)

		if kubernetesCluster.UpgradeAvailableTo != "" {
			ow.AppendDataWithLabel("KubernetesVersion", utility.Red(kubernetesCluster.KubernetesVersion+" *"), "Version")
		} else {
			ow.AppendDataWithLabel("KubernetesVersion", kubernetesCluster.KubernetesVersion, "Version")
		}

		ow.AppendDataWithLabel("APIEndPoint", kubernetesCluster.APIEndPoint, "API Endpoint")
		ow.AppendDataWithLabel("MasterIP", kubernetesCluster.MasterIP, "Master IP")
		ow.AppendDataWithLabel("DNSEntry", kubernetesCluster.DNSEntry, "DNS A record")

		if outputFormat == "json" || outputFormat == "custom" {
			//ow.AppendData("CloudConfig", template.CloudConfig)

			if kubernetesCluster.UpgradeAvailableTo != "" {
				ow.AppendDataWithLabel("KubernetesVersion", kubernetesCluster.KubernetesVersion, "Version")
			} else {
				ow.AppendDataWithLabel("KubernetesVersion", kubernetesCluster.KubernetesVersion, "Version")
			}

			ow.AppendDataWithLabel("UpgradeAvailableTo", kubernetesCluster.UpgradeAvailableTo, "Upgrade Available")

			if outputFormat == "json" {
				ow.WriteSingleObjectJSON()
			} else {
				ow.WriteCustomOutput(outputFields)
			}
		} else {
			ow.WriteKeyValues()

			if kubernetesCluster.UpgradeAvailableTo != "" {
				fmt.Println()
				fmt.Printf(utility.Red("* An upgrade to v%s is available. To upgrade use: civo k3s upgrade %s --version %s"), kubernetesCluster.UpgradeAvailableTo, kubernetesCluster.Name, kubernetesCluster.UpgradeAvailableTo)
				fmt.Println()
			}

			if len(kubernetesCluster.Instances) > 0 {
				fmt.Println()
				ow.WriteHeader("Nodes")
				owNode := utility.NewOutputWriter()

				for _, instance := range kubernetesCluster.Instances {
					owNode.StartLine()

					owNode.AppendData("Name", instance.Hostname)
					owNode.AppendData("IP", instance.PublicIP)
					owNode.AppendData("Status", instance.Status)
					owNode.AppendData("Size", instance.Size)
					owNode.AppendDataWithLabel("CPUCores", strconv.Itoa(instance.CPUCores), "Cpu Cores")
					owNode.AppendDataWithLabel("RAMMegabytes", strconv.Itoa(instance.RAMMegabytes), "Ram")
					owNode.AppendDataWithLabel("DiskGigabytes", strconv.Itoa(instance.DiskGigabytes), "SSD disk")
				}
				owNode.WriteTable()
			}

			if len(kubernetesCluster.InstalledApplications) > 0 {
				fmt.Println()
				ow.WriteHeader("Applications")
				owApp := utility.NewOutputWriter()

				for _, app := range kubernetesCluster.InstalledApplications {
					owApp.StartLine()

					owApp.AppendData("Name", app.Application)
					owApp.AppendData("Version", app.Version)
					owApp.AppendData("Installed", strconv.FormatBool(app.Installed))
					owApp.AppendData("Category", app.Category)
				}
				owApp.WriteTable()
			}

		}

	},
}
