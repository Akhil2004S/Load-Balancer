/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"loadBalancer/internal/metrics/proto"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// metricsCmd represents the metrics command
var metricsCmd = &cobra.Command{
	Use:   "metrics",
	Short: "Prints the current metrics of the load balancer",
	RunE: func(cmd *cobra.Command, args []string) error {
		metrics, err := creteClient()
		if err != nil {
			return err
		}
		fmt.Printf("The metrics are ready\nNumber of servers available:%d\nTotal requests handled:%d\nHealth of server 1:%s\nHealth of server 2:%s", metrics.NumServers, metrics.TotalRequests, metrics.ServerHealth[0], metrics.ServerHealth[1])
		return nil
	},
}

func init() {
	rootCmd.AddCommand(metricsCmd)
}

func creteClient() (*proto.Metrics, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.NewClient("127.0.0.1:9090", opts...)
	if err != nil {
		return nil, err
	}
	client := proto.NewAppMetricsClient(conn)
	metrics, err := client.GetMetrics(context.Background(), &proto.Empty{})
	if err != nil {
		return nil, err
	}
	return metrics, nil
}
