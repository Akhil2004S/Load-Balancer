/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"loadBalancer/internal/algorithms"
	loadbalancer "loadBalancer/internal/loadBalancer"
	lbServer "loadBalancer/internal/server"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type ServerConfig struct {
	Address string `mapstructure:"address"`
}

var serverConfig []ServerConfig

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts the server",
	RunE: func(cmd *cobra.Command, args []string) error {
		// We get the configuration value from Viper, not from the flag directly.
		port := viper.GetInt("port")
		algorithm := viper.GetString("algorithm")
		alpha := viper.GetFloat64("alpha")
		fmt.Printf("Starting load balancer on port: %d with algorithm:%s and the alpha value is:%f\n", port, algorithm, alpha)
		// In a real app, you would start a server here.
		data, err := initializeData(alpha, algorithm)
		if err != nil {
			log.Fatal(err)
		}
		if err := loadbalancer.StartServer(data, port); err != nil {
			log.Fatal(err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.Flags().Int("port", 3000, "Port to run the load balancer")
	startCmd.Flags().Float64("alpha", 0.2, "Alpha value for EMA calculation")
	startCmd.Flags().String("algorithm", "Least connections", "Algorithm to start the load balancer with.")

}

func initializeData(alpha float64, algorithm string) (*loadbalancer.Data, error) {
	data := &loadbalancer.Data{}
	viper.UnmarshalKey("servers", &serverConfig)

	switch algorithm {
	case "least_connections":
		data.Algorithm = algorithms.LeastConnections{}
	case "least_response_time":
		data.Algorithm = algorithms.LeastResponseTime{}
	case "ip_hash":
		data.Algorithm = algorithms.IPHash{}
	default:
		return nil, errors.New("ERROR. Choose one of the following algorithms. \n1.least_connections\n2.least_response_time\n3.ip_hash")
	}

	for id, server := range serverConfig {
		server := lbServer.CreateServer(id+1, server.Address)
		data.Servers = append(data.Servers, server)
	}
	data.Alpha = alpha
	return data, nil
}
