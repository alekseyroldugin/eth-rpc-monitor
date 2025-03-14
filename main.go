package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/yaml.v3"
)

// List of Ethereum RPC Endpoints
// var rpcEndpoints = []string{
// 	"https://rpc.ankr.com/eth",
// 	"https://eth.llamarpc.com",
// 	"https://mainnet.gateway.tenderly.co",
// 	"https://eth.rpc.blxrbdn.com",
// }

// Struct to store RPC endpoints from config
type RPCConfig struct {
	RPCEndpoints []struct {
		Name string `yaml:"name"`
		URL  string `yaml:"url"`
	} `yaml:"rpc_endpoints"`
}

var rpcConfig RPCConfig

// Track the latest block number
var latestBlockNumber uint64
var mu sync.Mutex // Mutex to prevent race conditions

// Prometheus Histogram with labels for measuring request latency percentiles
var rpcLatencyHistogram = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "ethereum_rpc_latency_seconds",
		Help:    "Latency of Ethereum RPC requests",
		Buckets: prometheus.ExponentialBuckets(0.01, 2, 10), // Buckets: 10ms, ~2s, ~10s
	},
	[]string{"rpc_name", "rpc_url"},
)

func init() {
	// Register the Prometheus histogram
	prometheus.MustRegister(rpcLatencyHistogram)
}

// Load RPC endpoints from config.yaml
func loadConfig() {
	file, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("❌ Failed to read config.yaml: %v", err)
	}
	err = yaml.Unmarshal(file, &rpcConfig)
	if err != nil {
		log.Fatalf("❌ Failed to parse config.yaml: %v", err)
	}
	fmt.Println("✅ Config loaded successfully!")
}

// Function to check block number & latency for a given RPC
func checkRPC(rpcName, rpcURL string) {
	// startTime := time.Now()

	// Connect to Ethereum node
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Printf("❌ Failed to connect to %s (%s): %v\n", rpcName, rpcURL, err)
		return
	}
	// connectionLatency := time.Since(startTime)

	// Measure block number request latency
	startTime := time.Now()
	blockNumber, err := client.BlockNumber(context.Background())
	if err != nil {
		log.Printf("❌ Failed to fetch block number from %s (%s): %v\n", rpcName, rpcURL, err)
		return
	}
	requestLatency := time.Since(startTime)

	// Record latency in Prometheus histogram
	rpcLatencyHistogram.WithLabelValues(rpcName, rpcURL).Observe(requestLatency.Seconds())

	// Print results
	fmt.Printf("\n🌍 RPC: %s (%s)\n", rpcName, rpcURL)
	// fmt.Printf("⏳ Connection Latency: %v ms\n", connectionLatency.Milliseconds())
	fmt.Printf("🔍 Latest Block: %d\n", blockNumber)
	fmt.Printf("⏱️ Request Latency: %v ms\n", requestLatency.Milliseconds())
}

// Function to monitor the latest Ethereum block and trigger checkRPC only when it changes
func monitorBlockChanges() {
	client, err := ethclient.Dial(rpcConfig.RPCEndpoints[0].URL) // Use the first RPC for monitoring
	if err != nil {
		log.Fatalf("❌ Failed to connect to Ethereum RPC: %v", err)
	}

	for {
		blockNumber, err := client.BlockNumber(context.Background())
		if err != nil {
			log.Printf("❌ Failed to fetch latest block: %v", err)
			time.Sleep(3 * time.Second) // Wait before retrying
			continue
		}

		mu.Lock()
		if blockNumber > latestBlockNumber { // Only update if a new block is found
			fmt.Printf("\n🚀 New Block Mined: %d\n", blockNumber)
			latestBlockNumber = blockNumber

			// Run checks for all RPCs
			for _, rpc := range rpcConfig.RPCEndpoints {
				go checkRPC(rpc.Name, rpc.URL) // Run in a separate goroutine
			}
		}
		mu.Unlock()

		time.Sleep(1 * time.Second) // Check again in 1 second
	}
}

// Start an HTTP server for Prometheus metrics
func startMetricsServer() {
	http.Handle("/metrics", promhttp.Handler())
	log.Println("📡 Prometheus metrics available at http://localhost:9090/metrics")
	log.Fatal(http.ListenAndServe(":9090", nil))
}

func main() {
	fmt.Println("🚀 Ethereum RPC Monitor Started...")

	// Load configuration
	loadConfig()

	// Start Prometheus metrics server in a separate goroutine
	go startMetricsServer()

	// Start monitoring for block changes
	monitorBlockChanges()
}
