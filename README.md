# GF_Simulator
Multicasting GF Simulator

# Pseudo-code
1) Have the "driver" read in a topology file that will list each node, with each node’s latitude, longitude, and neighbors  
2) The driver will then calculate the hyperbolic coordinates of each node, based off of the latitude and longitude  
3) The driver will then set up a number of workers equal to the number of nodes in the topology  
4) Each worker will listen on its own port, and be assigned its hyperbolic coordinates  
5) The driver will then iteratively select some amount of node combination (a node combination can either be in the format (A->B), meaning A wants information from B, or (A->B,C), meaning A wants information from B and C)  
6) For each node combination, each forwarding algorithm will be used to send a packet across the topology based on those hyperbolic coordinates  
7) An analysis will be calculated and output, showing things like latency (measuring RTT), packet loss (seeing if packets get dropped at a local minimum), path stretch (comparing the forwarding algorithm path with dijkstra's path), etc. for all forwarding algorithms specified  

# Running the Code
```./run_simulator.sh``` will run the code with all the default arguments  
```go run GF_Simulator.go``` will allow you to specify your own arguments  
```./run_unit_tests.sh``` will run the unit tests for the code to make sure everything is working  
**Note: It is recommended to set the number of open files to at least 64000 via the command 'ulimit -n 64000', as some algorithms send a large number of packets**  

### Command Line Arguments
```-topology_dir``` Default = "./topologies/" -- Where we are reading topologies from (specify directory where each item is a topology)  
```-ip``` Default = "localhost" -- What IP we want to work on  
```-port``` Default = 6000 -- What Port we want to work on (increments by 1)  
```-debug``` Default = false -- True means extra print statements   
```-distance``` Default = "hyperbolic" -- The distance formula we want to use ('euclidean' or 'hyperbolic')  
```-repeat``` Default = 1 -- How many times we want to run each experiment  
```-algorithms``` Default = "ALL_GPmGF" -- The algorithm(s) we want to test, comma seperated - 'FILO_GF', 'ALL_GF', 'FILO_GPGF', 'ALL_GPGF', 'G_FILO_GPGF', 'G_ALL_GPGF', 'FILO_mGF', 'ALL_mGF', 'FILO_GPmGF', 'ALL_GPmGF', 'G_FILO_GPmGF', 'G_ALL_GPmGF', 'Flooding', 'Multi_FILO_GF', 'Multi_ALL_GF', 'Multi_FILO_GPGF', 'Multi_ALL_GPGF', 'Multi_G_FILO_GPGF', 'Multi_G_ALL_GPGF', 'Multi_FILO_mGF', 'Multi_ALL_mGF', 'Multi_FILO_GPmGF', 'Multi_ALL_GPmGF', 'Multi_G_FILO_GPmGF', 'Multi_G_ALL_GPmGF', and/or 'Multi_Flooding' - additionally, the following groupings can be specified: 'all' for all algorithms, 'all_normal' for all non-multicast algorithms, or 'all_multi' for all multicast algorithms  
```-one_at_a_time``` Default = false -- Do we want to send packets one at a time or all at once  
```-collapse``` Default = true -- Do we want to do interest collapsing  
```-cache_strat``` Default = "LCE" -- Caching strategy - 'LCE' (Leave Copy Everywhere), 'RAND' (currently random at 50/50), or 'NONE' (no caching)  
```-cache_evict``` Default = "FIFO" -- Caching eviction policy - 'FIFO', 'LRU', 'LFU'  
```-ttl``` Default = -1 -- How long (in seconds) we want an item to remain fresh in a cache (-1 is infinite TTL)  
```-cache_size``` Default = .5 -- How many objects we want out cache to hold (0 for none, .5 for 50% of topology size, 1 for 100% of topology size (infinite size))  
```-start``` Default = 1, -- The percentage of nodes (chosen at random) to start sending from -- 1 means all nodes in the topology, .5 means 50% of all nodes in the topology"  
```-end```, Default = 1, -- The percentage of nodes (chosen at random) to receive -- 1 means all nodes in the topology, .5 means 50% of all nodes in the topology"  
```-output```, Default = "out.csv" -- What file the metrics will be output to  

# Metrics
1) Latency: measuring RTT  
2) Packet Loss: seeing if packets get dropped at a local minimum  
3) Path Stretch: comparing the forwarding algorithm path with dijkstra's path  
4) Packet Size: the size of the packet  
5) Pressure Count: used to count number of hops made in pressure mode  
6) Pressure Percentage: percentage of packets that entered pressure mode  
7) Nodes Visited: used to see the number of unique nodes visited  
8) Hops Collapsed: used to see the number of hops we saved by interest collapsing as a percentage of total hops if we collapsed (0 means we saved 0% of hops)  
9) Cache Hit ratio: used to see the number of cache hits we have -- number of cache hits for successful packets / successful packets 
10) Raw Number of Cache hits: used to see the number of times a cache hit occurred
11) Num Collapsed: used to see how many times we interest collapsed  
12) Consec Packets: used to see the ratio of the maximum number of consecutive packets / starting number of packets  
13) Hops: used to see the average number of hops an interest packet makes  
14) Satisfied Interests: used to see the number of satisfied requests per interest  
15) Ratio Hops per Satisfied: used to see the ratio of hops per interest / satisfied requests per interest  
14) PIT Entries: used to see the average number of PIT entries per node
15) Hanging PIT Entries: used to see how many PIT entries remain once the experiment is over

### Example Output
```
All at once, interest collapsing, caching, averaged over 1 result(s)	ALL_GPmGF
Average latency in seconds	0.202785	
Average percentage of packets successfully delivered	100	
Average path stretch	0.982857	
Average packet size in bytes	79.4747	
Average percentage of time spent in pressure mode on initial trip if entered pressure at all	79.9415	
Average packets that used pressure mode as a percentage of total packets	26.4227	
Average percentage of nodes visited	6.01565	
Average number of hops we saved by interest collapsing as a percentage of total hops if we collapsed	100	
Average cache hit ratio for successful packets based on -- Caching Strat: LCE -- Cache Eviction: FIFO -- TTL: 5 -- Cache Size: 0.5	75.2421	
Raw Number of Cache Hits	1398	
Average number of times we interest collapsed as a percentage of total packets	6.64686	
Ratio of the average maximum number of consecutive packets over number of starting maximum packets	2.52702	
Average number of hops per packet	6.74104	
Average number of satisfied requests per interest	1
Ratio of hops per interest / satisfied requests per interest	6.74104
Average number of PIT entries per node	140.876	
Number of Hanging PIT Entries	0	
```
