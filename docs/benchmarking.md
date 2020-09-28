# Tests

Common test setup:
 - Load test was performed from a VM placed inside of the same subnet as the fiber
 - nginx is used as the ingress into a fiber  
 - Payload is a standard driver rank request with 90 drivers (20kB)
 - Number of fiber replicas – *1*
 - Load test was conducted with `vegeta`

## Test 1

#### Setup:
 - Execution graph: eager router with 3 routes defined as a proxy call to the echo service
 - Echo service is deployed to the same cluster and subnet as the router
 - Echo service is deployed with 3 replicas and nginx as an ingress

#### Methodology:
 To estimate the latency added by the fiber itself, it's required to measure the latency of echo service, defined as its routes.
Since fiber's eager router sends the incoming request to all of its routes simultaneous, it would be fair to measure the latency of the
echo service with as `<number of routes> * <rps to fiber>` RPS.

Example:
For estimating the latency of eager router at 25 RPS, we need to 
 1. run `vegeta` test against echo service with 75 RPS (3 routes * 25)
 2. run `vegeta` test against fiber service with 25 RPS
 3. subtract echo service latencies (1) from fiber service latencies (2)

#### Results:

25 RPS
```text
Requests      [total, rate, throughput]  1500, 25.02, 24.99
Duration      [total, attack, wait]      1m0.021971994s, 59.960013204s, 61.95879ms
Latencies     [mean, 50, 95, 99, max]    9.383368ms, 8.228431ms, 12.872032ms, 15.933171ms, 24.977658ms
Bytes In      [total, mean]              29620500, 19747.00
Bytes Out     [total, mean]              29620500, 19747.00
Success       [ratio]                    100.00%
Status Codes  [code:count]               200:1500
Error Set:
```

50 RPS
```text
Requests      [total, rate, throughput]  3000, 50.02, 49.96
Duration      [total, attack, wait]      1m0.044889623s, 59.980140052s, 64.749571ms
Latencies     [mean, 50, 95, 99, max]    9.710702ms, 8.711095ms, 14.421771ms, 16.204055ms, 6.250714ms
Bytes In      [total, mean]              59241000, 19747.00
Bytes Out     [total, mean]              59241000, 19747.00
Success       [ratio]                    100.00%
Status Codes  [code:count]               200:3000
Error Set:
```

75 RPS
```text
Requests      [total, rate, throughput]  4500, 75.02, 74.92
Duration      [total, attack, wait]      1m0.060802736s, 59.986666282s, 74.136454ms
Latencies     [mean, 50, 95, 99, max]    10.289462ms, 9.661406ms, 13.832404ms, 22.80185ms, 48.900073ms
Bytes In      [total, mean]              88861500, 19747.00
Bytes Out     [total, mean]              88861500, 19747.00
Success       [ratio]                    100.00%
Status Codes  [code:count]               200:4500
Error Set:
```

100 RPS
```text
Requests      [total, rate, throughput]  6000, 100.02, 99.90
Duration      [total, attack, wait]      1m0.051105094s, 59.989938604s, 61.16649ms
Latencies     [mean, 50, 95, 99, max]    8.881556ms, 8.302955ms, 13.37626ms, 21.725708ms, -21.29657ms
Bytes In      [total, mean]              118462416, 19743.74
Bytes Out     [total, mean]              118482000, 19747.00
Success       [ratio]                    99.98%
Status Codes  [code:count]               200:5999  502:1
Error Set:
502 Bad Gateway
```

150+ RPS
```text
TBU
```

## Test 2

#### Setup:
 - Execution graph: eager router with 3 routes defined as a mock echo component, that waits for `50ms` 
   and returns request payload (no requests to the external services)
 
#### Results:

25 RPS
```text
Requests      [total, rate, throughput]  1500, 25.02, 24.99
Duration      [total, attack, wait]      1m0.01693917s, 59.959925167s, 57.014003ms
Latencies     [mean, 50, 95, 99, max]    57.565112ms, 57.450232ms, 58.79543ms, 60.822437ms, 75.69441ms
Bytes In      [total, mean]              30000, 20.00
Bytes Out     [total, mean]              30000, 20.00
Success       [ratio]                    100.00%
Status Codes  [code:count]               200:1500
Error Set:
```

50 RPS
```text
Requests      [total, rate, throughput]  3000, 50.02, 49.97
Duration      [total, attack, wait]      1m0.03728246s, 59.979956528s, 57.325932ms
Latencies     [mean, 50, 95, 99, max]    57.487151ms, 57.315496ms, 58.81139ms, 61.515856ms, 76.683483ms
Bytes In      [total, mean]              60000, 20.00
Bytes Out     [total, mean]              60000, 20.00
Success       [ratio]                    100.00%
Status Codes  [code:count]               200:3000
Error Set:
```

75 RPS
```text
Requests      [total, rate, throughput]  4500, 75.02, 74.94
Duration      [total, attack, wait]      1m0.045237452s, 59.98668001s, 58.557442ms
Latencies     [mean, 50, 95, 99, max]    59.339303ms, 59.120167ms, 61.861432ms, 65.080973ms, 87.96062ms
Bytes In      [total, mean]              90000, 20.00
Bytes Out     [total, mean]              90000, 20.00
Success       [ratio]                    100.00%
Status Codes  [code:count]               200:4500
Error Set:
```

100 RPS
```text
Requests      [total, rate, throughput]  6000, 100.02, 99.92
Duration      [total, attack, wait]      1m0.046049204s, 59.989899543s, 56.149661ms
Latencies     [mean, 50, 95, 99, max]    57.878004ms, 57.5575ms, 60.137835ms, 63.625902ms, 124.472525ms
Bytes In      [total, mean]              120000, 20.00
Bytes Out     [total, mean]              120000, 20.00
Success       [ratio]                    100.00%
Status Codes  [code:count]               200:6000
Error Set:
```

150 RPS
```text
Requests      [total, rate, throughput]  9000, 150.02, 149.87
Duration      [total, attack, wait]      1m0.050703606s, 59.993350511s, 57.353095ms
Latencies     [mean, 50, 95, 99, max]    57.912805ms, 57.445643ms, 60.540051ms, 67.29764ms, 112.959971ms
Bytes In      [total, mean]              180000, 20.00
Bytes Out     [total, mean]              180000, 20.00
Success       [ratio]                    100.00%
Status Codes  [code:count]               200:9000
Error Set:
```