/*
Copyright © 2022 balajisa09

*/
package cmd

import (
	//"context"
	"fmt"
	"os"
	"flag"
	"path/filepath"
	"context"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	calico_cs "github.com/projectcalico/api/pkg/client/clientset_generated/clientset"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// rootCmd represents the base command when called without any subcommands

var rootCmd = &cobra.Command{
	Use:   "k8s-netpol",
	Short: "k8s-netpol is a CLI application for k8s network security analysis",
	Long: `k8s-netpol is a CLI application for k8s network security analysis , currenty compactible with calico CNI`,
	Run: func(cmd *cobra.Command, args []string) {
		var kubeconfig *string

		//set kube config path
		if home := homedir.HomeDir(); home != "" {
			kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		} else {
			kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
		}
		flag.Parse()

		//build config
		config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
		
		if err != nil {
			panic(err)
		}
		fmt.Println("The host path for the current cluster is: ",config.Host)

		//create clientset
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			panic(err)
		}

		//get list of standard network policies
		stdnetpols ,err := clientset.NetworkingV1().NetworkPolicies("calico-new").List(context.Background(), metav1.ListOptions{})
		if(err != nil){
			panic(err)
		}
		foundstdNetpols := false
		for _,netpol := range stdnetpols.Items{
			foundstdNetpols = true
			fmt.Println(netpol.Name)
		}

		if(!foundstdNetpols){
			fmt.Println("No standard networkpolicies found")
		}

		//create calico clientset
		calico_clientset, err := calico_cs.NewForConfig(config)
		if err != nil {
			panic(err)
		}
		//get list of calico networkpolicies
		calico_netpols, err := calico_clientset.ProjectcalicoV3().NetworkPolicies("calico-new").List(context.Background(), metav1.ListOptions{})

		if(err != nil){
			panic(err)
		}
		foundcalicoNetpols := false
		for _,netpol := range calico_netpols.Items{
			foundcalicoNetpols = true
			fmt.Println(netpol.Name)
		}

		if(!foundcalicoNetpols){
			fmt.Println("No Calico networkpolicies found")
		}

	
	 },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.k8s-netpol.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", true, "Help message for toggle")
}


