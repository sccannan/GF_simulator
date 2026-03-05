package main

//----------------------------------------------
// Imports
//----------------------------------------------
import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"net"
	"io"
	"io/ioutil"
	"math"
	"bytes"
	"encoding/gob"
	"time"
	"reflect"
	"sync"
	"math/rand"
	"flag"
	"strings"
	"errors"
	"sort"
)
//----------------------------------------------

//----------------------------------------------
// Globals
//----------------------------------------------
//Metrics
var num_Collapsed 	int 		//the number of times we interest collapsed
var num_Collapsed_Lock 	sync.Mutex 	//lock for interest collapsed
var open_Processes 	int 		//the number of open packets (not dropped or successed)
var open_Processes_Lock sync.Mutex 	//lock for processes
var open_Processes_Max 	int 		//the highest concurrent number of open packets
var cache_Hit 		int 		//used to track the number of cache hits we got 
var cache_Hit_Lock 	sync.Mutex 	//lock for cache hits
var hops_Made		int		//number of packets sent
var hops_Made_Lock	sync.Mutex 	//lock for number of hops made

//Toggles
var debug		bool		//true for extra debug info
//----------------------------------------------

//----------------------------------------------
// Node Struct Declaration
//----------------------------------------------
type Node struct {
	IP		string		//the ip
	Port		int		//the port
	Number		int		//the node number
	Lat		float64 	//latitude
	Long		float64 	//longitude
	Theta		float64 	//hyperbolic theta
	R		float64 	//hyperbolic R
	Data_Name	string  	//a node will satisfy a request with this data name
	Data_Value 	string  	//what the payload will be
	Neighbors	[]int		//a list of all the neighbors nodes as numbers 
	Cache		[]CACHE 	//a cache for each node
	Cache_Size	int		//how many objects we store in a cache
	Cache_Lock	sync.Mutex 	//lock for CACHE
	Pit		[]PIT		//PIT needs to record incoming interface, outgoing interface, and name of the data
	Pit_Entries	int		//maximum number of PIT entries
	Pit_Lock	sync.Mutex 	//lock for PIT
	Hops		int		//number of hops this node has sent
	Hops_Lock	sync.Mutex 	//lock for hops
}
//----------------------------------------------

//----------------------------------------------
// PIT Struct Declaration
//----------------------------------------------
type PIT struct {
	Incoming	int		//the node number of the incoming interface
	Outgoing	int		//the node number of the outgoing interface
	Name		string		//the name for the interest
	Start 		int 		//only used for metrics -- tracks the start
	Destination 	int 		//only used for metrics -- tracks the destination
	Collapsed	int		//only used for metrics -- if this value is not -1, means we interest collapsed on that hop
}
//----------------------------------------------

//----------------------------------------------
// CACHE Struct Declaration
//----------------------------------------------
type CACHE struct {
	Name		string		//the name for the interest
	Data		string  	//the payload of a packet
	Number		int		//the number of the producer
	Timestamp	time.Time 	//timestamp of when this was placed in cache
	TTL		time.Duration 	//how long FROM THE TIMESTAMP in seconds until this entry is no longer valid
	Last_Used	time.Time	//used for LRU
	Frequency	int		//used for LFU
}
//----------------------------------------------

//----------------------------------------------
// Packet Struct Declaration
//----------------------------------------------
type Packet struct {

	//For all packets
	Previous	int		//the node number of the previous interface
	Start		int		//the starting index of where we want to put the packet
	Destination	int		//the node number of the destination interface
	Destination_Str string		//the node number(s) of the destination interface, seperated by '\'
	Interest_Data	int 		//0 if interest, 1 if data packet, 2 for shutdown signal
	Name		string		//the name for the interest
	Payload		string		//the payload if data
	Start_Time 	time.Time 	//time we sent this packet
	TTL		time.Duration 	//how long FROM THE TIMESTAMP in seconds until this packet is no longer valid
	
	//For gravity/pressure mode
	GP_Mode		int		//0 for gravity, 1 for pressure
	Nodes_Visit	[]int		//a list of the node numbers we have visited just in pressure node
	Num_Visit	[]int		//a list of the number of times we have visited a node just in pressure node
	GP_Distance 	float64 	//the distance at local minimum
	
	//For metrics, do not count when adding
	Path_Traversed 	[]int 		//records the path this packet has taken
	Pressure_Gauge 	int 		//counter for how many times this packet went it pressure mode
	Index_Middle 	int		//index of the destination 
	End_Time	time.Time 	//time we received this packet
	Collapsed 	int		//what index we interest collapsed on (-1 if N/A)
	Cache_Hit	int		//1 if this packet was due to a cache hit, 0 else
}
//----------------------------------------------

//----------------------------------------------
// Reads in command line arguments and prints them
//----------------------------------------------
func read_args() (string, string, int, string, int, []string, bool, bool, string, string, float64, float64, float64, float64, string) {

	//Reads in all the args
	data_Input := flag.String("topology_dir", "./topologies/", "Where we are reading topologies from (specify directory where each item is a topology)")
	ip := flag.String("ip", "localhost", "What IP we want to work on")
	port := flag.Int("port", 6000, "What Port we want to work on (increments by 1)")
	temp_debug := flag.Bool("debug", false, "True means extra print statements")
	distance_Formula := flag.String("distance", "hyperbolic", "The distance formula we want to use ('euclidean' or 'hyperbolic')")
	average_Num := flag.Int("repeat", 1, "How many times we want to run each experiment")
	algorithms_Input := flag.String("algorithms", "ALL_GPmGF", "The algorithm(s) we want to test, comma seperated - 'FILO_GF', 'ALL_GF', 'FILO_GPGF', 'ALL_GPGF', 'G_FILO_GPGF', 'G_ALL_GPGF', 'FILO_mGF', 'ALL_mGF', 'FILO_GPmGF', 'ALL_GPmGF', 'G_FILO_GPmGF', 'G_ALL_GPmGF', 'Flooding', 'Multi_FILO_GF', 'Multi_ALL_GF', 'Multi_FILO_GPGF', 'Multi_ALL_GPGF', 'Multi_G_FILO_GPGF', 'Multi_G_ALL_GPGF', 'Multi_FILO_mGF', 'Multi_ALL_mGF', 'Multi_FILO_GPmGF', 'Multi_ALL_GPmGF', 'Multi_G_FILO_GPmGF', 'Multi_G_ALL_GPmGF', and/or 'Multi_Flooding' - additionally, the following groupings can be specified: 'all' for all algorithms, 'all_normal' for all non-multicast algorithms, or 'all_multi' for all multicast algorithms")  
	one_at_a_time := flag.Bool("one_at_a_time", false, "Do we want to send packets one at a time or all at once")
	interest_Collapse := flag.Bool("collapse", true, "Do we want to do interest collapsing")
	caching_Strat := flag.String("cache_strat", "LCE", "Caching strategy - 'LCE' (Leave Copy Everywhere), 'RAND' (currently random at 50/50), or 'NONE' (no caching)")
	cache_Evict := flag.String("cache_evict", "FIFO", "Caching eviction policy - 'FIFO', 'LRU', 'LFU'")
	TTL := flag.Float64("ttl", -1, "How long (in seconds) we want an interest and associated data to be valid (packet TTL and cache TTL) (-1 is infinite TTL)")
	cache_Size := flag.Float64("cache_size", .5, "How many objects we want out cache to hold (0 for none, .5 for 50% of topology size, 1 for 100% of topology size (infinite size))")
	start_Nodes := flag.Float64("start", 1, "The percentage of nodes (chosen at random) to start sending from -- 1 means all nodes in the topology, .5 means 50% of all nodes in the topology")
	end_Nodes := flag.Float64("end", 1, "The percentage of nodes (chosen at random) to receive -- 1 means all nodes in the topology, .5 means 50% of all nodes in the topology")
	output := flag.String("output", "out.csv", "What file the metrics will be output to")
	
	flag.Parse()
	
	//Prints out all the args
	fmt.Println("Data Input:", *data_Input)
	fmt.Println("IP:", *ip)
	fmt.Println("Port:", *port)
	fmt.Println("Debug:", *temp_debug)
	debug = *temp_debug
	fmt.Println("Distance Formula:", *distance_Formula)
	fmt.Println("Average Number:", *average_Num)
	algorithms := strings.Split(*algorithms_Input, ",")
	if algorithms[0] == "all"{
		algorithms = []string{"FILO_GF", "ALL_GF", "FILO_GPGF", "ALL_GPGF", "G_FILO_GPGF", "G_ALL_GPGF", "FILO_mGF", "ALL_mGF", "FILO_GPmGF", "ALL_GPmGF", "G_FILO_GPmGF", "G_ALL_GPmGF", "Flooding", "Multi_FILO_GF", "Multi_ALL_GF", "Multi_FILO_GPGF", "Multi_ALL_GPGF", "Multi_G_FILO_GPGF", "Multi_G_ALL_GPGF", "Multi_FILO_mGF", "Multi_ALL_mGF", "Multi_FILO_GPmGF", "Multi_ALL_GPmGF", "Multi_G_FILO_GPmGF", "Multi_G_ALL_GPmGF", "Multi_Flooding"}	
	}
	if algorithms[0] == "all_normal"{
		algorithms = []string{"FILO_GF", "ALL_GF", "FILO_GPGF", "ALL_GPGF", "G_FILO_GPGF", "G_ALL_GPGF", "FILO_mGF", "ALL_mGF", "FILO_GPmGF", "ALL_GPmGF", "G_FILO_GPmGF", "G_ALL_GPmGF", "Flooding"}
	}
	if algorithms[0] == "all_multi"{
		algorithms = []string{"Multi_FILO_GF", "Multi_ALL_GF", "Multi_FILO_GPGF", "Multi_ALL_GPGF", "Multi_G_FILO_GPGF", "Multi_G_ALL_GPGF", "Multi_FILO_mGF", "Multi_ALL_mGF", "Multi_FILO_GPmGF", "Multi_ALL_GPmGF", "Multi_G_FILO_GPmGF", "Multi_G_ALL_GPmGF", "Multi_Flooding"}
	}
	fmt.Print("Algorithms: ")
	for x := 0; x < len(algorithms); x++{
		if x != len(algorithms)-1{
			fmt.Print(algorithms[x] + ", ")
		} else{
			fmt.Print(algorithms[x] + "\n")
		}
	}
	fmt.Println("One at a time?:", *one_at_a_time)
	fmt.Println("Interest Collapse?:", *interest_Collapse)
	fmt.Println("Caching Strat:", *caching_Strat)
	fmt.Println("Cache Eviction Policy:", *cache_Evict)
	fmt.Println("Cache/Packet TTL:", *TTL)
	fmt.Println("Cache Size:", *cache_Size)
	fmt.Println("Start Percentage:", *start_Nodes)
	fmt.Println("End Percentage:", *end_Nodes)
	fmt.Println("Output File:", *output)
	
	return *data_Input, *ip, *port, *distance_Formula, *average_Num, algorithms, *one_at_a_time, *interest_Collapse, *caching_Strat, *cache_Evict, *TTL, *cache_Size, *start_Nodes, *end_Nodes, *output
}
//----------------------------------------------

//----------------------------------------------
// Exits code if error
//----------------------------------------------
func check_Error(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
//----------------------------------------------

//----------------------------------------------
// Print statements in debug mode
//----------------------------------------------
func debug_Print(str string) {
	if debug == true{
		fmt.Printf(str)
	}
}
//----------------------------------------------

//----------------------------------------------
// Returns the hyperbolic distance from A to B
//----------------------------------------------
func hyperbolic_Distance(r_A float64, r_B float64, theta_A float64, theta_B float64) float64{
	dtheta := math.Pi - math.Abs(math.Pi - math.Abs(theta_A - theta_B))
	
	//this is here to solve a rounding problem
	//without this here, if dtheta=0, r_A=1 and r_B=3, 1.999999... gets returned instead of 2
	//if this isnt here, you get functionally the same answer
	if dtheta == 0{
		return math.Abs(r_A - r_B)
	}
	return math.Acosh(math.Cosh(r_A) * math.Cosh(r_B) - math.Sinh(r_A) * math.Sinh(r_B) * math.Cos(dtheta))
}
//----------------------------------------------

//----------------------------------------------
// Returns the euclidean distance from A to B
//----------------------------------------------
func euclidean_Distance(lat_A float64, lat_B float64, long_A float64, long_B float64) float64{
	return math.Sqrt(math.Pow((lat_A-lat_B),2) + math.Pow((long_A-long_B),2))
}
//----------------------------------------------

//----------------------------------------------
// Distance Helper Function
//----------------------------------------------
func distance_Helper(A *Node, B *Node, distance_Formula string) float64{
	var distance float64
	if distance_Formula == "euclidean"{ //euclidean distance
		distance = euclidean_Distance(A.Lat, B.Lat, A.Long, B.Long)
	} else if distance_Formula == "hyperbolic"{ //hyperbolic distance
		distance = hyperbolic_Distance(A.R, B.R, A.Theta, B.Theta)
	} else{
		fmt.Println("Unrecognized distance formula! Exiting.")
		os.Exit(1)
	}
	
	//this just exists as error checking
	//this will trigger is Acosh(x<1)
	//should never trigger realistically, but if it does it is the same as a link not existing, so max out the distance
	//realistically occurs if input value is NaN
	if math.IsNaN(distance) == true{
		return math.MaxFloat64
	}
	return distance
}
//----------------------------------------------

//----------------------------------------------
// Dijkstra's, only counting number of hops
// Returns the optimal number of hops
// Only counts hops, NOT weight
//----------------------------------------------
func dijkstras(y int, z int, topology []*Node) int{
	local_hops := 0
	if y == z{ //if we are already done
		return local_hops
	}
	nodes := []int{y} //start with the current node
	seen_Before := []int{y} //list of nodes weve seen before
	for x := 0; x < len(topology); x++{ //maximum we go for each node in topology
	
		//for each node, grab all its neighbors as ints
		var new_Nodes []int
		for y := 0; y < len(nodes); y++{ 
			new_Nodes = append(new_Nodes, find_Node(topology, nodes[y]).Neighbors...)
		}

		//Only look for unique nodes, as anything repeating will always be more hops 
		//append all unique nodes to temp and to see before
		var temp []int
		for y := 0; y < len(new_Nodes); y++{
			add_Node := true
			for z := 0; z < len(seen_Before); z++{
				if new_Nodes[y] == seen_Before[z]{
					add_Node = false
					break
				}
			}
			if add_Node == true{
				temp = append(temp, new_Nodes[y])
			}
		}
		local_hops++

		//if one of the values is at the destination, return
		for y := 0; y < len(temp); y++{
			if temp[y] == z{
				return local_hops
			}
		}
		
		//set nodes to unique nodes from previous run, so you dont recompute nodes already computed
		nodes = temp
	}
	return -1 //return -1 if fail to find
}
//----------------------------------------------

//----------------------------------------------
// Returns the average a slice
//----------------------------------------------
func avg(arr []float64) float64{
	if len(arr) == 0{
		return 0.0
	}
	sum := 0.0
	for x := 0; x < len(arr); x++{
		sum = sum + arr[x]
	}
	return sum/float64(len(arr))
}
//----------------------------------------------

//----------------------------------------------
// Returns the std of a slice
//----------------------------------------------
func std(arr []float64) float64 {
	if len(arr) == 0{
		return 0.0
	}
	mean := avg(arr)
	std := 0.0
	for x := 0; x < len(arr); x++ {
		std = std + math.Pow(arr[x]-mean, 2)
	}
	std = math.Sqrt(std/float64(len(arr)-1)) //assume the value is a sample, not the entire population
	return std
}
//----------------------------------------------

//----------------------------------------------
// Checks if the data is within std_range standard deviations
// std_range=1 is 68% of data, std_range=2 is 95% of data, std_range=3 is 99.7% of data
//----------------------------------------------
func stdCheck(arr []float64, std_range int) bool {
	if len(arr) == 0{ //if there is no data, its all within any standard deviation
		return true
	}
	standard_Deviation := std(arr)
	mean := avg(arr)
	check := true
	for x := 0; x < len(arr); x++ {
		if arr[x] > mean + (float64(std_range) * standard_Deviation) || arr[x] < mean - (standard_Deviation * float64(std_range)){ //if the data is within std_range standard deviations
			check = false	
			debug_Print("Value: " + strconv.FormatFloat(arr[x], 'f', -1, 64)  + "\n")
		}
	}
	return check
}
//----------------------------------------------

//----------------------------------------------
// Returns the number of the neighbor with the closest distance to the end
// AND whether or not the drop condition is met
// NOTE: Use neighbors instead of curr_Node.neighbors for mgf_helper and GPGF pressure mode (to exclude certain neighbors)
//----------------------------------------------
func gf_Helper(curr_Node *Node, end *Node, neighbors []*Node, distance_Formula string) (int, bool){
	min_Distance := math.MaxFloat64
	
	//first check to see if a neighbor is the destination
	for x := 0; x < len(neighbors); x++{
		if neighbors[x].Number == end.Number{
			return neighbors[x].Number, false
		}
	}
	
	//calculate each distance, tracking the minimum
	index := -1
	for x := 0; x < len(neighbors); x++{
		distance := distance_Helper(end, neighbors[x], distance_Formula)
		debug_Print("\tNeighbor " + strconv.Itoa(neighbors[x].Number) + ", distance " + fmt.Sprintf("%f", distance) + " \n")
		if distance < min_Distance{
			min_Distance = distance
			index = x
		}
	}

	//If all distances are math.MaxFloat64 or neighbors = []
	if index == -1{
		fmt.Println("Error finding node with smallest distance! Exiting.")
		os.Exit(1)
	}
	
	//Packet drop condition
	drop_Packet := false
	if curr_Node.Number == neighbors[index].Number{ //drop packet if we send the packet to the current node
		drop_Packet = true
	}
	return neighbors[index].Number, drop_Packet
}
//----------------------------------------------

//----------------------------------------------
// Returns the number of the neighbor with the closest distance to the end
// AND whether or not the drop condition is met
//----------------------------------------------
func mgf_Helper(curr_Node *Node, end *Node, neighbors []*Node, previous_Node_Num int, distance_Formula string) (int, bool){

	//Remove yourself as a neighbor if you are a neighbors
	index := -1
	for x := 0; x < len(neighbors); x++{
		if neighbors[x].Number == curr_Node.Number{
			index = x
			break
		}
	}
	
	//Actual removal
	usable_Neighbors := neighbors
	if index != -1{
		temp_Neighbors := make([]*Node, 0)
		temp_Neighbors = append(temp_Neighbors, neighbors[:index]...)
		temp_Neighbors = append(temp_Neighbors, neighbors[index+1:]...)
		usable_Neighbors = temp_Neighbors
	}
	
	next_Node, _ := gf_Helper(curr_Node, end, usable_Neighbors, distance_Formula)
	
	//Packet drop condition
	drop_Packet := false
	if previous_Node_Num == next_Node{ //drop packet if we send the packet to its previous node
		drop_Packet = true
	}
	return next_Node, drop_Packet
}
//----------------------------------------------

//----------------------------------------------
// Encodes and sends a message
//----------------------------------------------
func send_Msg(packet Packet, c net.Conn) {
	var buf bytes.Buffer //where the encoded message is stored
	enc := gob.NewEncoder(&buf) //create an encoder
	err := enc.Encode(packet) //encode the packet
	check_Error(err) //check for errors
	_, err = c.Write(buf.Bytes()) //send to destination over the connection
	check_Error(err) //check for errors
	err = c.Close() //close the connection
	check_Error(err) //check for errors
}
//----------------------------------------------

//----------------------------------------------
// Receives and decodes a message
// Returns the transmitted packet
//----------------------------------------------
func receive_Msg(c net.Conn) Packet{
	var packet Packet //where the decoded message is stored
	raw := make([]byte, 99999) //default buffer size of 
	_, err := c.Read(raw) //read the data into raw
	check_Error(err) //check for errors
	dec := gob.NewDecoder(bytes.NewBuffer(raw)) //creating a decoder
	err = dec.Decode(&packet) //decoding the value and storing it in args
	check_Error(err) //check for errors
	err = c.Close() //close the connection
	check_Error(err) //check for errors
	return packet
}
//----------------------------------------------

//----------------------------------------------
// Returns a node with the given value
// Returns a dummy node with -1 if that node doesnt exist or topology is empty
// NOTE: While this is functionally the same as returning topology[x],
// this accounts for the '-1' node
//----------------------------------------------
func find_Node(topology []*Node, number int) *Node{
	for x := 0; x < len(topology); x++{
		if topology[x].Number == number{
			return topology[x]
		}
	}
	return &Node{Number: -1} //return a dummy node if the node is not in the topology
}
//----------------------------------------------

//----------------------------------------------
//Returns the []*Node for each neighbor a node has
//----------------------------------------------
func find_Neighbors(topology []*Node, number int) []*Node{
	curr_node := find_Node(topology, number) //grab the current node
	var node_neighbors []*Node
	for y := 0; y < len(curr_node.Neighbors); y++{ //for each neighbor of the current node
		node_neighbors = append(node_neighbors, find_Node(topology, curr_node.Neighbors[y])) //based on the number, append the node
	}
	return node_neighbors
}
//----------------------------------------------

//----------------------------------------------
// Returns the LAST PIT entry corresponding to the named data we are looking for
// Must also match the outgoing interface (this function is used for return trips) if Outgoing isnt -2
// Must also match the incoming interface (this function is used for return trips) if Incoming isnt -2
// Returns dummy PIT with -2 if the PIT entry is not found
//----------------------------------------------
func index_PIT(curr_Node *Node, name string, Outgoing int, Incoming int) *PIT{
	for x := len(curr_Node.Pit)-1; x >=0 ; x--{
		if curr_Node.Pit[x].Name == name{
			if Outgoing == -2 && Incoming == -2{
				return &curr_Node.Pit[x]
			} else if Outgoing != -2 && Incoming != -2 && curr_Node.Pit[x].Outgoing == Outgoing && curr_Node.Pit[x].Incoming == Incoming{
				return &curr_Node.Pit[x]
			} else if Outgoing != -2 && Incoming == -2 && curr_Node.Pit[x].Outgoing == Outgoing{
				return &curr_Node.Pit[x]
			} else if Outgoing == -2 && Incoming != -2 && curr_Node.Pit[x].Incoming == Incoming{
				return &curr_Node.Pit[x]	
			} 
		}
	}
	return &PIT{Incoming: -2, Outgoing: -2}
}
//----------------------------------------------

//----------------------------------------------
// Removes the last copy of local_PIT from the current nodes PIT
// Exits if not found
//----------------------------------------------
func remove_PIT(curr_Node *Node, local_PIT *PIT) {
	
	//Find index of removal
	index := -1
	for x := len(curr_Node.Pit)-1; x >= 0; x--{ //find the index of the local_PIT inside all pit entries
		if curr_Node.Pit[x] == *local_PIT{
			index = x
			break
		}
	}
	
	//Actual removal
	if index != -1{
		temp_PIT := make([]PIT, 0)
		temp_PIT = append(temp_PIT, curr_Node.Pit[:index]...)
		temp_PIT = append(temp_PIT, curr_Node.Pit[index+1:]...)
		curr_Node.Pit = temp_PIT
	} else{
		fmt.Println("Error finding the PIT entry to remove! Exiting.")
		os.Exit(1)	
	}
}
//----------------------------------------------

//----------------------------------------------
// Returns the cache entry if we find the name we like
// Will return the first entry, but caching is set up here that there should 
// only be one entry per name
// Returns a cache entry with '0 time' as the timestamp if nothing is found
//----------------------------------------------
func index_CACHE(curr_Node *Node, name string) *CACHE{
	for x := 0; x < len(curr_Node.Cache); x++{
		if curr_Node.Cache[x].Name == name{
			return &curr_Node.Cache[x]
		}	
	}
	return &CACHE{Timestamp: time.Time{}}
}
//----------------------------------------------

//----------------------------------------------
// Removes the cache local_CACHE from the current nodes cache
// Exits if not found
//----------------------------------------------
func remove_CACHE(curr_Node *Node, local_CACHE *CACHE)  {
	
	//Remove yourself as a neighbor if you are a neighbors
	index := -1
	for x := 0; x < len(curr_Node.Cache); x++{ //find the index of the local_CACHE inside all cache entries
		if curr_Node.Cache[x] == *local_CACHE{
			index = x
			break
		}
	}
	
	//Actual removal
	if index != -1{
		temp_CACHE := make([]CACHE, 0)
		temp_CACHE = append(temp_CACHE, curr_Node.Cache[:index]...)
		temp_CACHE = append(temp_CACHE, curr_Node.Cache[index+1:]...)
		curr_Node.Cache = temp_CACHE
	} else{
		fmt.Println("Error finding the CACHE entry to remove! Exiting.")
		os.Exit(1)	
	}
}
//----------------------------------------------

//----------------------------------------------
// Determines if we cache based off of the caching strategy 
// Returns true if we cache, false else
//----------------------------------------------
func determine_CACHE(caching_Strat string) bool{
	if caching_Strat == "LCE"{
		return true
	} else if caching_Strat == "RAND"{
		if rand.Intn(2) == 0{
			return true	
		}else{
			return false
		}
	} else if caching_Strat == "NONE"{
		return false
	} else{
		fmt.Println("Unknown caching strategy! Exiting.")
		os.Exit(1)	
	}
	return false
}
//----------------------------------------------

//----------------------------------------------
// Evicts a cache based on TTL and cache eviction policy ((1)Sameness -> (2)TTL -> (3)Policy takes priority)
// Only works if there is at least 1 entry in the cache
//----------------------------------------------
func evict_CACHE(curr_Node *Node, name string, cache_Evict string) {

	//Error if empty cache
	if len(curr_Node.Cache) == 0{
		fmt.Println("Called evict_CACHE on empty cache! Exiting.")
		os.Exit(1)	
	}

	cache_To_Evict := curr_Node.Cache[0]
	same_Eviction := false
	ttl_Eviction := false

	//First, evict a packet if its an older version of the current packet
	for x := 0; x < len(curr_Node.Cache); x++{
		if name == curr_Node.Cache[x].Name{
			cache_To_Evict = curr_Node.Cache[x]
			same_Eviction = true
			break
		}
	}
		
	//Second, evict a packet if it no longer alive (based on TTL)
	if same_Eviction == false{
		for x := 0; x < len(curr_Node.Cache); x++{
			if curr_Node.Cache[x].TTL >= time.Duration(0){ //dont evict based on TTL if TTL is negative
				time_Added := (curr_Node.Cache[x].Timestamp).Add(curr_Node.Cache[x].TTL)
				if (time.Now()).After(time_Added){
					cache_To_Evict = curr_Node.Cache[x]
					ttl_Eviction = true
					break
				}	
			}
		}
	}
	
	//If neither of those options worked, evict per policy
	if ttl_Eviction == false && same_Eviction == false{ //If we cant evict based on TTL 
		if cache_Evict == "FIFO"{//Remove the oldest cache entry
			oldest_Value := curr_Node.Cache[0].Timestamp
			for x := 1; x < len(curr_Node.Cache); x++{
				if (curr_Node.Cache[x].Timestamp).Before(oldest_Value){
					oldest_Value = curr_Node.Cache[x].Timestamp
					cache_To_Evict = curr_Node.Cache[x]
				}	
			}
		}else if cache_Evict == "LRU"{ //Remove the least recently used cache entry
			oldest_Value := curr_Node.Cache[0].Last_Used
			for x := 1; x < len(curr_Node.Cache); x++{
				if (curr_Node.Cache[x].Last_Used).Before(oldest_Value){
					oldest_Value = curr_Node.Cache[x].Last_Used
					cache_To_Evict = curr_Node.Cache[x]
				}	
			}
		}else if cache_Evict == "LFU"{  //Remove the least frequently used cache entry
			oldest_Value := curr_Node.Cache[0].Frequency
			for x := 1; x < len(curr_Node.Cache); x++{
				if curr_Node.Cache[x].Frequency < oldest_Value{
					oldest_Value = curr_Node.Cache[x].Frequency
					cache_To_Evict = curr_Node.Cache[x]
				}	
			}
		} else{
			fmt.Println("Unknown cache eviction policy! Exiting.")
				os.Exit(1)	
		}
	}
	
	if same_Eviction == true{
		debug_Print("Same Evicted!\n")
	} else if ttl_Eviction == true{
		debug_Print("TTL Evicted!\n")
	} else if ttl_Eviction == false && same_Eviction == false{
		debug_Print("Policy Evicted!\n")
	} else {
		debug_Print("Error!\n")
		os.Exit(1)
	}
	remove_CACHE(curr_Node, &cache_To_Evict) 
}
//----------------------------------------------

//----------------------------------------------
// Checks if we are following the reverse of the path 
// we took for successful packets for
// If reverse isn't taken, error and exit
//----------------------------------------------
func reverse_Check(path_Traversed []int) {
	check := true

	//Path cant be even in length (1,2,3,2,1 = good, 1,2,3,3,2,1 = bad)
	//Path can be 0 or 1 in length (path of 1 means the first hop satisfies the request because we do not include the driver node in the path)
	if (len(path_Traversed)%2 == 0 && len(path_Traversed) != 0){
		check = false
	}
	
	//Individually check path
	for x := 0; x < len(path_Traversed); x++{
		if path_Traversed[x] != path_Traversed[len(path_Traversed)-1-x]{
			check = false
			break
		}
	}
	
	if check == false{
		fmt.Printf("Error! Not following the reverse of the taken path!\n")
		for x := 0; x < len(path_Traversed); x++{
			if x == len(path_Traversed)-1{
				debug_Print(strconv.Itoa(path_Traversed[x]) + "\n")
			} else{
				debug_Print(strconv.Itoa(path_Traversed[x]) + ", ")
			}
		}
		os.Exit(1)
	}
}
//----------------------------------------------

//----------------------------------------------
// Checks to make sure all PITs are empty on a packet success 
// If PIT isnt clear, error and exit
//----------------------------------------------
func empty_PIT_Check(topology []*Node) {
	for x := 0; x < len(topology); x++{
		if len(topology[x].Pit) != 0{
			fmt.Printf("Error! PIT is not empty for node " + strconv.Itoa(topology[x].Number) + "!\n")
			for y := 0; y < len(topology[x].Pit); y++{
				debug_Print("IN: " + strconv.Itoa(topology[x].Pit[y].Incoming) + ", ")
				debug_Print("OUT: " + strconv.Itoa(topology[x].Pit[y].Outgoing) + "\n")
			}
			os.Exit(1)
		}
	}
}
//----------------------------------------------

//----------------------------------------------
// Calculates the size of a packet in bytes
//----------------------------------------------
func calc_Size(packet Packet) float64{
	size := 0.0
	v := reflect.ValueOf(packet)
	for x := 0; x < v.NumField(); x++{
		field_Name := v.Type().Field(x).Name
		if field_Name == "Previous" || field_Name == "Start" || field_Name == "Destination" || field_Name == "Interest_Data"{ //int
			size = size + 4
		} else if field_Name == "Name" || field_Name == "Destination_Str"{ //string
			size = size + float64(len(v.Field(x).Interface().(string))) //1 byte per character
		} else if field_Name == "Start_Time"{ //time.Time *can* be represented in 8 bytes
			size = size + 8 
		} else if field_Name == "TTL"{ //time.Duration
			size = size + 8
		} else if field_Name == "GP_Mode"{
			if v.Field(x).Interface().(int) != -1{ //int != -1
				size = size + 4
			}
		} else if field_Name == "Nodes_Visit" || field_Name == "Num_Visit"{ //[]int
			size = size + float64(len(v.Field(x).Interface().([]int))*4)
		} else if field_Name == "GP_Distance"{
			if v.Field(x).Interface().(float64) != -1.0{ //float != -1
				size = size + 8
			}
		} else if field_Name == "Payload" || field_Name == "Path_Traversed" || field_Name == "Pressure_Gauge" || field_Name == "Index_Middle" || field_Name == "End_Time" || field_Name == "Collapsed" || field_Name == "Cache_Hit"{ //skip
			continue
		} else{
			fmt.Printf("Error! Undefined field: " + field_Name + "\n")
			os.Exit(1)
		}
	}
	return size
}
//----------------------------------------------

//----------------------------------------------
// Reads in and returns a topology based on a file
//----------------------------------------------
func read_Topology(file_Path string, ip string, port int, cache_Size float64) []*Node{
	var topology []*Node //list of nodes in the topology
	_, err := os.Stat(file_Path)
	if errors.Is(err, os.ErrNotExist){ //see if file is real
		fmt.Printf("Error! File: " + file_Path + " was not found!\n")
		os.Exit(1)	
	}
	file, err := os.Open(file_Path)// open file
	check_Error(err) //check for errors
	csvReader := csv.NewReader(file) //create a csvReader object

	// Read in the values line by line
	counter := -1 //to see how many items we have (-1 b/c increment on header)
	for {
		rec, err := csvReader.Read() //read the next line
		if err == io.EOF { //break out if its the end of the file
			break
		}
		if counter == -1{ //if the line we read in is the header or empty, skip it
			counter++
		} else{
			check_Error(err) //check for errors
			var neighbors []int //list of neighbors
			for x := 3; x < len(rec); x++{
				neighbor, err := strconv.Atoi(rec[x]) //cast to int
				check_Error(err) //check for errors
				if neighbor == 1{
					neighbors = append(neighbors, x-3) //append to the list
				}
			}
			node_Number, err := strconv.Atoi(rec[0]) //cast string to int
			check_Error(err) //check for errors
			node_Lat, err := strconv.ParseFloat(rec[1], 64) //cast string to float
			check_Error(err) //check for errors
			node_Long, err := strconv.ParseFloat(rec[2], 64) //cast string to float
			check_Error(err) //check for errors
			temp_Node := &Node{IP: ip, Port: (port + counter), Number: node_Number, Lat: node_Lat, Long: node_Long, Theta: -1.0, R: -1.0, Data_Name: "/ucla/videos/demo.mpg/1/" + strconv.Itoa(node_Number), Data_Value: ("Lat: " + strconv.FormatFloat(node_Lat, 'E', -1, 64) + "Long: " + strconv.FormatFloat(node_Long, 'E', -1, 64)), Neighbors: neighbors, Hops: 0} //create a struct for node
			topology = append(topology, temp_Node) //append to the list
			counter++
		}
	}
	file.Close() //close the file
	
	//Update the hyperbolic coordinates and cache size
	for x := 0; x < len(topology); x++{
		topology[x].Theta = math.Atan(topology[x].Lat/topology[x].Long)
		topology[x].R = math.Sqrt((topology[x].Lat*topology[x].Lat) + (topology[x].Long*topology[x].Long))
		topology[x].Cache_Size = int(float64(len(topology)) * float64(cache_Size))
	}
	
	//Assertion - You must be a neighbor to yourself
	for x := 0; x < len(topology); x++{
		status := false
		for y := 0; y < len(topology[x].Neighbors); y++{
			if topology[x].Neighbors[y] == x{
				status = true
				break
			}
		}
		if status == false{
			fmt.Printf("Error! Node " + strconv.Itoa(x) + " must be a neighbor to itself!\n")
			os.Exit(1)
		}
	}	
	
	//Assertion - You must have symmetrical neighbors (if A -> B exists, B -> A must exist)
	for x := 0; x < len(topology); x++{
		for y := 0; y < len(topology[x].Neighbors); y++{
			status := false
			for z := 0; z < len(topology[topology[x].Neighbors[y]].Neighbors); z++{
				if topology[topology[x].Neighbors[y]].Neighbors[z] == topology[x].Number{
					status = true
					break
				}
			}
			if status == false{
				fmt.Printf("Error! Node " + strconv.Itoa(x) + " or Node " + strconv.Itoa(topology[x].Neighbors[y]) + " has not symmetrical neighbors!\n")
				os.Exit(1)
			}
		}
	}
	return topology
}
//----------------------------------------------

//----------------------------------------------
// Calculates metrics
// Only called on a successful packet
//----------------------------------------------
func metrics_Calc(packet Packet, topology []*Node) []float64{

	//Latency
	latency := (packet.End_Time).Sub(packet.Start_Time).Seconds() //measures RTT
	
	//Packet Success
	packet_success := 1.0 //this function is only called for a successful packet
	
	//Stretch
	d_hops := dijkstras(packet.Start, packet.Destination, topology) //calculate optimal hops
	RTT_d_hops := float64(d_hops*2) //multiplied by 2 because its to and from
	total_Hops := float64(len(packet.Path_Traversed)-1)
	packet_Hops := float64(len(packet.Path_Traversed)-1)
	if packet.Collapsed != -1{ //we did interest collapse
		//if our packet says it took the path [0, 1, 2, 3, 4, 3, 2, 1, 0] and we interest collapsed at 1, we only made 2 hops instead of 8
		packet_Hops = float64(packet.Collapsed * 2)
	}
	path_stretch := packet_Hops/RTT_d_hops
	
	//Packet Size
	packet_size := calc_Size(packet)
	
	//Pressure Count and Pressure Mode Used
	pressure_Count := -1.0
	pressure_Mode_Used := 0.0
	if packet.Pressure_Gauge > 0{ //we only want to look at algorithms that use pressure mode at least once
		pressure_Count = float64(packet.Pressure_Gauge)/float64(packet.Index_Middle) //number of hops we made in pressure / number of hops to the destination (aka the middle)
		pressure_Mode_Used = 1.0
	}
	
	//Nodes Visited
	nodes_Visited := float64(len(packet.Nodes_Visit))/float64(len(topology))
	
	//Hops collapsed
	hops_Collapsed := 0.0 //by default, assume we didnt interest collapse, therefor moaking the ratio 0
	if packet.Collapsed != -1{ //we did interest collapse
		//if our packet says it took the path [0, 1, 2, 3, 4, 3, 2, 1, 0] and we interest collapsed at 1, we only made 2 hops instead of 8, making the ratio .75 (we saved 75% of hops reported)
		hops_Collapsed = (total_Hops - packet_Hops)/total_Hops
	}
	
	//Cache Hit
	cache_Hit := float64(packet.Cache_Hit)
	
	fmt.Println(packet.Name)
	fmt.Println(packet.Destination_Str)
	
	return []float64{latency, packet_success, path_stretch, packet_size, pressure_Count, pressure_Mode_Used, nodes_Visited, hops_Collapsed, cache_Hit}
}
//----------------------------------------------

//----------------------------------------------
// Sets globals
//----------------------------------------------
func algorithm_Check(algorithm string) (bool, bool, bool, bool, bool) {	
	filo := false
	global := false
	modified := false
	gravity_pressure := false
	multi := false
	
	if !(algorithm == "FILO_GF" || algorithm == "ALL_GF" ||algorithm == "FILO_mGF" ||algorithm == "ALL_mGF" ||algorithm == "FILO_GPGF" || algorithm == "ALL_GPGF" || algorithm == "G_FILO_GPGF" || algorithm == "G_ALL_GPGF" || algorithm == "FILO_GPmGF" || algorithm == "ALL_GPmGF" || algorithm == "G_FILO_GPmGF" || algorithm == "G_ALL_GPmGF" || algorithm == "Flooding" || algorithm == "Multi_Flooding" || algorithm == "Multi_FILO_GF" || algorithm == "Multi_ALL_GF" ||algorithm == "Multi_FILO_mGF" ||algorithm == "Multi_ALL_mGF" ||algorithm == "Multi_FILO_GPGF" || algorithm == "Multi_ALL_GPGF" || algorithm == "Multi_G_FILO_GPGF" || algorithm == "Multi_G_ALL_GPGF" || algorithm == "Multi_FILO_GPmGF" || algorithm == "Multi_ALL_GPmGF" || algorithm == "Multi_G_FILO_GPmGF" || algorithm == "Multi_G_ALL_GPmGF"){
		fmt.Println("Unrecognized algorithm! Exiting.")
		os.Exit(1)
	}
	if algorithm == "FILO_GF" || algorithm == "FILO_mGF" || algorithm == "FILO_GPGF" || algorithm == "G_FILO_GPGF" || algorithm == "FILO_GPmGF" || algorithm == "G_FILO_GPmGF" || algorithm == "Multi_FILO_GF" || algorithm == "Multi_FILO_mGF" || algorithm == "Multi_FILO_GPGF" || algorithm == "Multi_G_FILO_GPGF" || algorithm == "Multi_FILO_GPmGF" || algorithm == "Multi_G_FILO_GPmGF"{
		filo = true
	}
	if algorithm == "G_FILO_GPGF" || algorithm == "G_ALL_GPGF" || algorithm == "G_FILO_GPmGF" || algorithm == "G_ALL_GPmGF" || algorithm == "Multi_G_FILO_GPGF" || algorithm == "Multi_G_ALL_GPGF" || algorithm == "Multi_G_FILO_GPmGF" || algorithm == "Multi_G_ALL_GPmGF"{
		global = true
	}
	if algorithm == "FILO_mGF" || algorithm == "ALL_mGF" || algorithm == "FILO_GPmGF" || algorithm == "ALL_GPmGF" || algorithm == "G_FILO_GPmGF" || algorithm == "G_ALL_GPmGF" || algorithm == "Multi_FILO_mGF" || algorithm == "Multi_ALL_mGF" || algorithm == "Multi_FILO_GPmGF" || algorithm == "Multi_ALL_GPmGF" || algorithm == "Multi_G_FILO_GPmGF" || algorithm == "Multi_G_ALL_GPmGF"{
		modified = true
	}
	if algorithm == "FILO_GPGF" || algorithm == "ALL_GPGF" || algorithm == "G_FILO_GPGF" || algorithm == "G_ALL_GPGF" || algorithm == "FILO_GPmGF" || algorithm == "ALL_GPmGF" || algorithm == "G_FILO_GPmGF" || algorithm == "G_ALL_GPmGF" || algorithm == "Multi_FILO_GPGF" || algorithm == "Multi_ALL_GPGF" || algorithm == "Multi_G_FILO_GPGF" || algorithm == "Multi_G_ALL_GPGF" || algorithm == "Multi_FILO_GPmGF" || algorithm == "Multi_ALL_GPmGF" || algorithm == "Multi_G_FILO_GPmGF" || algorithm == "Multi_G_ALL_GPmGF"{
		gravity_pressure = true
	}
	if algorithm == "Multi_Flooding" || algorithm == "Multi_FILO_GF" || algorithm == "Multi_ALL_GF" ||algorithm == "Multi_FILO_mGF" ||algorithm == "Multi_ALL_mGF" ||algorithm == "Multi_FILO_GPGF" || algorithm == "Multi_ALL_GPGF" || algorithm == "Multi_G_FILO_GPGF" || algorithm == "Multi_G_ALL_GPGF" || algorithm == "Multi_FILO_GPmGF" || algorithm == "Multi_ALL_GPmGF" || algorithm == "Multi_G_FILO_GPmGF" || algorithm == "Multi_G_ALL_GPmGF"{
		multi = true
	}
	return filo, global, modified, gravity_pressure, multi
}
//----------------------------------------------

//----------------------------------------------
// Returns the cache entry and whether a cache hit occured
// If a cache hit occured, increment the global
//----------------------------------------------
func check_Cache_Hit(curr_Node *Node, name string) (*CACHE, bool){
	curr_Node.Cache_Lock.Lock()
	cache_Hit_Bool := false
	temp_Cache := index_CACHE(curr_Node, name) //grab the index if a cache hit
	if !((temp_Cache.Timestamp).Equal(time.Time{})){ //index_cache returns TTL at the beginning of time if fail to find
		time_Added := (temp_Cache.Timestamp).Add(temp_Cache.TTL)
		if (time.Now()).Before(time_Added) || temp_Cache.TTL < time.Duration(0){ // if we are within TTL
			cache_Hit_Lock.Lock()
			cache_Hit = cache_Hit + 1 //update cache hit
			cache_Hit_Lock.Unlock()
			temp_Cache.Last_Used = time.Now() //update last used for LRU
			temp_Cache.Frequency = temp_Cache.Frequency + 1 //update frequency for LFU
			debug_Print("Cache Hit\n")
			cache_Hit_Bool = true
		}
	}
	curr_Node.Cache_Lock.Unlock()
	return temp_Cache, cache_Hit_Bool
}
//----------------------------------------------

//----------------------------------------------
// Caches the data
// Will always cache at the end of the cache, even if overwriting
//----------------------------------------------
func cache_Data(curr_Node *Node, packet Packet, cache_Evict string){
	curr_Node.Cache_Lock.Lock()
			
	//Check to see if the data is already in the cache
	same_Eviction := false
	for x := 0; x < len(curr_Node.Cache); x++{
		if packet.Name == curr_Node.Cache[x].Name{
			same_Eviction = true
			break
		}
	}
	
	//If the cache is full or we can update an old copy, evict an entry
	if (len(curr_Node.Cache) == curr_Node.Cache_Size) || same_Eviction == true{ 
		evict_CACHE(curr_Node, packet.Name, cache_Evict)
	}
	
	//Add the new cache entry
	curr_Node.Cache = append(curr_Node.Cache, CACHE{packet.Name, packet.Payload, packet.Destination, time.Now(), packet.TTL, time.Time{}, 0})//cache the data
	curr_Node.Cache_Lock.Unlock()
}
//----------------------------------------------

//----------------------------------------------
// Calculates the next packet it we are in pressure mode
// Returns the node number
//----------------------------------------------
func pressure_Mode_Helper(curr_Node *Node, packet Packet, topology []*Node, distance_Formula string, destination int) int{

	//Calculate how many times we have visited each neighbor
	var times_Visited []int
	for x := 0; x < len(curr_Node.Neighbors); x++{
		visit_Before := false
		for y := 0; y < len(packet.Nodes_Visit); y++{
			if curr_Node.Neighbors[x] == packet.Nodes_Visit[y]{
				times_Visited = append(times_Visited, packet.Num_Visit[y])
				visit_Before = true
				break
			}
		}
		if visit_Before == false{
			times_Visited = append(times_Visited, 0)
		}
	}

	//Calculate the nodes we have visited the least
	var min []int
	num_Visited := math.MaxInt64
	for x := 0; x < len(times_Visited); x++{ //first pass to get the minimum
		if times_Visited[x] <= num_Visited{
			num_Visited = times_Visited[x]
		}
	}
	for x := 0; x < len(times_Visited); x++{ //second pass to get the indexs
		if times_Visited[x] == num_Visited{
			min = append(min, curr_Node.Neighbors[x])
		}
	}

	//Keep only nodes from the list of minimum visits
	var new_Neighbors []*Node
	for x := 0; x < len(min); x++{
		new_Neighbors = append(new_Neighbors, find_Node(topology, min[x]))
	}
	end := find_Node(topology, destination)
	next_Node_Num, _ := gf_Helper(curr_Node, end, new_Neighbors, distance_Formula) //call the helper
	return next_Node_Num
}
//----------------------------------------------

//----------------------------------------------
// Returns whether or not the packet was interest collapsed
// If interest collapsing occured, increment the global
// Name is seperate for stitched names
//----------------------------------------------
func interest_Collapse_Helper(curr_Node *Node, name string, packet Packet, filo bool, destination int) bool{
	status := false
	
	//Never interest collapse for FILO algorithms, as they wont be following the reverse path
	//Never interest collapse while a packet is in pressure mode or was just in pressure mode
	if filo != false || packet.GP_Mode == 1{
		return status
	}
	
	curr_Node.Pit_Lock.Lock()
	
	//Smart Interest Collapsing
	//Check if the node we just came from is already an outgoing interface in the PIT
	//If it is, do NOT interest collapse
	smart := false
	temp_Local_PIT := index_PIT(curr_Node, name, packet.Previous, -2)
	if (*temp_Local_PIT != PIT{Incoming: -2, Outgoing: -2}){ 
		debug_Print("Smart Interest Collapsing - Not Collapsed\n")
		smart = true
	}
	
	//Normal interest collapsing
	temp_Local_PIT = index_PIT(curr_Node, name, -2, -2) //see if the value is already in the PIT
	if (*temp_Local_PIT != PIT{Incoming: -2, Outgoing: -2}) && smart == false{ //we can only interest collapse if the value is already in the PIT
		//If we are here: the name is already in the PIT, we are NOT a FILO algorithm, and if we are a GPGF algorithm we are NOT in pressure mode
		debug_Print("Collapsed\n")
		num_Collapsed_Lock.Lock()
		num_Collapsed++
		num_Collapsed_Lock.Unlock()
		curr_Node.Pit = append(curr_Node.Pit, PIT{packet.Previous, -2, name, packet.Start, destination, len(packet.Path_Traversed)-1}) //add it to the PIT			
		status = true
	}
	curr_Node.Pit_Lock.Unlock()
	return status
}
//----------------------------------------------

//----------------------------------------------
// Calculates the next node when an interest is returning for FILO
// Also clears the PITs
//----------------------------------------------
func return_FILO(curr_Node *Node, packet Packet, topology []*Node) []*Node{
	var next_Nodes []*Node
	
	//Find the index of where we currently are
	current_Index := len(packet.Path_Traversed) - 1 
	
	//Calculate the next index
	next_Index := packet.Index_Middle - (current_Index - packet.Index_Middle) - 1
	
	//Presumed Incoming Interface is on the path or -1
	presumed_Incoming := -1		
	if next_Index != -1{ //we are not about to be done
		presumed_Incoming = packet.Path_Traversed[next_Index]
	}
	
	curr_Node.Pit_Lock.Lock()
	temp_Local_PIT := index_PIT(curr_Node, packet.Name, packet.Previous, presumed_Incoming) //we are going backwards, old incoming = new outgoing and old outgoing = new incoming
	next_Nodes = append(next_Nodes, find_Node(topology, temp_Local_PIT.Incoming))
	remove_PIT(curr_Node, temp_Local_PIT)
	curr_Node.Pit_Lock.Unlock()
	return next_Nodes
}
//----------------------------------------------

//----------------------------------------------
// Calculates the next node when an interest is returning for ALL
// Also clears the PITs
//----------------------------------------------
func return_ALL(curr_Node *Node, packet Packet, topology []*Node) ([]*Node, []int, []int, []int){
	curr_Node.Pit_Lock.Lock()
	var next_Nodes []*Node
	var start []int
	var destination []int
	var collapsed []int
	for{
		temp_Local_PIT := index_PIT(curr_Node, packet.Name, -2, -2) //grab every PIT entry with a matching name
		if (*temp_Local_PIT == PIT{Incoming: -2, Outgoing: -2}){ //returns -2 if no PIT entries
			break
		}
		next_Nodes = append(next_Nodes, find_Node(topology, temp_Local_PIT.Incoming))
		start = append(start, temp_Local_PIT.Start)
		destination = append(destination, temp_Local_PIT.Destination) 
		collapsed = append(collapsed, temp_Local_PIT.Collapsed)
		remove_PIT(curr_Node, temp_Local_PIT)
	}
	curr_Node.Pit_Lock.Unlock()
	return next_Nodes, start, destination, collapsed
}
//----------------------------------------------

//----------------------------------------------
// Randomly returns a percentage of nodes from the topology
//----------------------------------------------
func random_Selection(topology []*Node, percentage float64) []int{
	
	//Make a list of numbers
	var numbers []int
	for x := 0; x < len(topology); x++{
		numbers = append(numbers, x)
	}
	
	//Calculate how many values we need
	num_random := int(math.Ceil(percentage * float64(len(numbers))))
	
	//Select nodes
	var ret_arr_num []int
	for x := 0; x < num_random; x++{
		node_Num := rand.Intn(len(numbers))
		ret_arr_num = append(ret_arr_num, numbers[node_Num])
		numbers = append(numbers[:node_Num], numbers[node_Num+1:]...)
	}
	
	//Sort the numbers
	sort.Ints(ret_arr_num)
	return ret_arr_num
}
//----------------------------------------------

//----------------------------------------------
// Returns all node numbers that at neighbors to the current node
// Except where the packet just came from and the current node
//----------------------------------------------
func flooding_Helper(curr_Node *Node, previous_Node_Num int) []int{

	//Create a list of numbers of neighbors
	original_Neighbors := make([]int, len(curr_Node.Neighbors))
	copy(original_Neighbors, curr_Node.Neighbors)
	
	//Remove yourself as a neighbor if you are a neighbors
	index := -1
	for x := 0; x < len(original_Neighbors); x++{
		if original_Neighbors[x] == curr_Node.Number{
			index = x
			break
		}
	}
	
	//Actual removal
	original_Neighbors_Not_Self := original_Neighbors
	if index != -1{
		temp_Neighbors := make([]int, 0)
		temp_Neighbors = append(temp_Neighbors, original_Neighbors[:index]...)
		temp_Neighbors = append(temp_Neighbors, original_Neighbors[index+1:]...)
		original_Neighbors_Not_Self = temp_Neighbors
	}
	
	//Remove the previous node as a neighbor
	index = -1
	for x := 0; x < len(original_Neighbors_Not_Self); x++{
		if original_Neighbors_Not_Self[x] == previous_Node_Num{
			index = x
			break
		}
	}
	
	//Actual removal
	original_Neighbors_Not_Self_Prev := original_Neighbors_Not_Self
	if index != -1{
		temp_Neighbors := make([]int, 0)
		temp_Neighbors = append(temp_Neighbors, original_Neighbors_Not_Self[:index]...)
		temp_Neighbors = append(temp_Neighbors, original_Neighbors_Not_Self[index+1:]...)
		original_Neighbors_Not_Self_Prev = temp_Neighbors
	}	
	
	return original_Neighbors_Not_Self_Prev
}
//----------------------------------------------

//----------------------------------------------
// Returns a copy of the packet
//----------------------------------------------
func copy_Packet(packet Packet) Packet{
	return packet
}
//----------------------------------------------

//----------------------------------------------
// Updates list of which names we still need
//----------------------------------------------
func update_Names(names []string, satisfied_Names []string, destinations []string) ([]string, string, bool, []string){
	var removal_Index []int
	drop_Packet := false
	new_Concat_Name := ""
	
	//If we have already satisfied a certain name, remove it going forward
	for x := 0; x < len(names); x++{
		for y := 0; y < len(satisfied_Names); y++{	
			if names[x] == satisfied_Names[y]{
				removal_Index = append(removal_Index, x)
				break
			}
		}
	}
	for x := len(removal_Index)-1; x >= 0; x--{
		temp_Name := make([]string, 0)
		temp_Name = append(temp_Name, names[:removal_Index[x]]...)
		temp_Name = append(temp_Name, names[removal_Index[x]+1:]...)
		names = temp_Name
		temp_Destination := make([]string, 0)
		temp_Destination = append(temp_Destination, destinations[:removal_Index[x]]...)
		temp_Destination = append(temp_Destination, destinations[removal_Index[x]+1:]...)
		destinations = temp_Destination
	}
	
	//Check drop packet condition
	if len(names) == 0{
		drop_Packet = true
	} else{
		drop_Packet = false
		for x := 0; x < len(names); x++{
			new_Concat_Name = new_Concat_Name + names[x]
			if x != len(names)-1{
				new_Concat_Name = new_Concat_Name + "|"
			}
		}
	}
	return names, new_Concat_Name, drop_Packet, destinations
}
//----------------------------------------------

//----------------------------------------------
// Packet Comparison
// If exact == false, used for multicasting (ignore Destination and Name)
//----------------------------------------------
func compare_Packet(packet_A Packet, packet_B Packet, exact bool) bool{
	if packet_A.Previous != packet_B.Previous{
		return false
	}
	if packet_A.Start != packet_B.Start{
		return false
	}
	if packet_A.Destination != packet_B.Destination{
		if exact == true{
			return false
		}
	}
	if packet_A.Destination_Str != packet_B.Destination_Str{
		return false
	}	
	if packet_A.Interest_Data != packet_B.Interest_Data{
		return false
	}
	if packet_A.Name != packet_B.Name{
		if exact == true{
			return false
		}
	}		
	if packet_A.Payload != packet_B.Payload{
		return false
	}		
	if !((packet_A.Start_Time).Equal(packet_B.Start_Time)){
		return false
	}
	if packet_A.TTL != packet_B.TTL{
		return false
	}	
	if packet_A.GP_Mode != packet_B.GP_Mode{
		return false
	}
	if len(packet_A.Nodes_Visit) != len(packet_B.Nodes_Visit){
		return false
	}
	for x := 0; x < len(packet_A.Nodes_Visit); x++{
		if packet_A.Nodes_Visit[x] != packet_B.Nodes_Visit[x]{
			return false
		}
	}		
	if len(packet_A.Num_Visit) != len(packet_B.Num_Visit){
		return false
	}
	for x := 0; x < len(packet_A.Num_Visit); x++{
		if packet_A.Num_Visit[x] != packet_B.Num_Visit[x]{
			return false
		}
	}
	if packet_A.GP_Distance != packet_B.GP_Distance{
		return false
	}			
	if len(packet_A.Path_Traversed) != len(packet_B.Path_Traversed){
		return false
	}
	for x := 0; x < len(packet_A.Path_Traversed); x++{
		if packet_A.Path_Traversed[x] != packet_B.Path_Traversed[x]{
			return false
		}
	}	
	if packet_A.Pressure_Gauge != packet_B.Pressure_Gauge{
		return false
	}
	if packet_A.Index_Middle != packet_B.Index_Middle{
		return false
	}
	if !((packet_A.End_Time).Equal(packet_B.End_Time)){
		return false
	}
	if packet_A.Collapsed != packet_B.Collapsed{
		return false
	}
	if packet_A.Cache_Hit != packet_B.Cache_Hit{
		return false
	}
	return true
}
//----------------------------------------------

//----------------------------------------------
// Stitches names of similar packets together
//----------------------------------------------
func combine_Packet(temp_Interest_Node []*Node, temp_Interest_Packet []Packet) ([]*Node, []Packet){

	var next_Nodes []*Node
	var next_Packets []Packet
	var used []int
	
	if len(temp_Interest_Node) != len(temp_Interest_Packet){
		fmt.Println("Error! Mismatching lengths for packets and next nodes! Exiting.")
		os.Exit(1)
	}
	
	//For each packet
	for x := 0; x < len(temp_Interest_Packet); x++{
	
		//If we have not seen this packet before
		contin := true
		for z := 0; z < len(used); z++{
			if used[z] == x{
				contin = false
				break
			}
		}
		if contin == true{
		
			//We have now seen this packet
			used = append(used, x)
			
			//Init the packet for comparison
			new_Out := copy_Packet(temp_Interest_Packet[x])
			temp_Names := new_Out.Name
			temp_Destinations := strconv.Itoa(new_Out.Destination)
			
			//Compare it against all other packets that havent been see
			for y := x+1; y < len(temp_Interest_Packet); y++{
				same := compare_Packet(temp_Interest_Packet[x], temp_Interest_Packet[y], false)
				
				//If the packets are the same and the next node is the ssame, and the second packet hasnt been seen before
				if same == true && temp_Interest_Node[x].Number == temp_Interest_Node[y].Number{
					contin = true
					for z := 0; z < len(used); z++{
						if used[z] == y{
							contin = false
							break
						}
					}
					if contin == true{
						used = append(used, y)
						temp_Names = temp_Names + "|" + temp_Interest_Packet[y].Name
						temp_Destinations = temp_Destinations + "|" + strconv.Itoa(temp_Interest_Packet[y].Destination)
					}
				}
			}
			new_Out.Name = temp_Names
			new_Out.Destination_Str = temp_Destinations
			new_Out.Destination = -2
			next_Nodes = append(next_Nodes, temp_Interest_Node[x])
			next_Packets = append(next_Packets, new_Out)
		}
	}
	return next_Nodes, next_Packets
}
//----------------------------------------------

//----------------------------------------------
// Waits until the app node pairings are received
// Returns all the packets weve gotten
//----------------------------------------------
func wait_Main(packet_Contents chan Packet, packets [][]Packet) [][]Packet{
	for { //while processes are still active
		select{
			case return_Packet := <-packet_Contents: //grab the successful packet
				start := return_Packet.Start	
				end := return_Packet.Destination
				if end != -2{
					if packets[start][end].Name == return_Packet.Name && packets[start][end].Payload == "" {
						if ((time.Now()).Before((return_Packet.Start_Time).Add(return_Packet.TTL)) || return_Packet.TTL < time.Duration(0)){ //if the TTL is valid
							return_Packet.Start_Time = packets[start][end].Start_Time
							return_Packet.End_Time = time.Now()
							packets[start][end] = return_Packet //in case we have multiple successful packets, keep only the first one
						}	
					}
				}
		}
		open_Processes_Lock.Lock()
		open_Processes = open_Processes - 1
		if open_Processes == 0{
			open_Processes_Lock.Unlock()
			break	
		}
		open_Processes_Lock.Unlock()
	}
	return packets
}
//----------------------------------------------

//----------------------------------------------
// Sets up a node to listen on an IP and port
// Calls serviceConnection on a connection
//----------------------------------------------
func setup_Node(curr_Node *Node, shutdown_ACK chan int, packet_Contents chan Packet, topology []*Node, distance_Formula string, algorithm string, caching_Strat string, cache_Evict string, interest_Collapse bool) {
	l, err := net.Listen("tcp", curr_Node.IP+ ":" + strconv.Itoa(curr_Node.Port)) //listen on ip:port
	check_Error(err) //check for errors
	debug_Print("Node " + strconv.Itoa(curr_Node.Number) + " up on " + curr_Node.IP + ":" + strconv.Itoa(curr_Node.Port) + " \n")
	shutdown_ACK <- 1 //reuse this channel, send value when properly setup
	loop:
	for { //forever until blreak
		c, err := l.Accept() //accept a connection
		check_Error(err) //check for errors
		packet := receive_Msg(c) //receive and decode the message
		if packet.Previous == -999 && packet.Destination == -999 && packet.Interest_Data == 2 && packet.Name == "Shutdown"{
			break loop
		}
		go serviceConnection(curr_Node, packet_Contents, topology, packet, distance_Formula, algorithm, caching_Strat, cache_Evict, interest_Collapse)
	}
	debug_Print("Shutting down Node " + strconv.Itoa(curr_Node.Number) + " \n")
	l.Close()
	shutdown_ACK <- 1
}
//----------------------------------------------

//----------------------------------------------
// Thread to handle connections
//----------------------------------------------
func serviceConnection(curr_Node *Node, packet_Contents chan Packet, topology []*Node, packet Packet, distance_Formula string, algorithm string, caching_Strat string, cache_Evict string, interest_Collapse bool){
	filo, global, modified, gravity_pressure, multi := algorithm_Check(algorithm)
	debug_Print("Node " + strconv.Itoa(curr_Node.Number) + " received a message from Node " + strconv.Itoa(packet.Previous) + " \n")
	var next_Nodes []*Node
	var next_Packets []Packet
	var names []string
	var destinations []string
	var satisfied_Names []string
	continue_Forward := true	
		
	//Create names
	temp_Name := strings.Clone(packet.Name)
	for strings.Index(temp_Name, "|") != -1{
		index := strings.Index(temp_Name, "|")
		names = append(names, temp_Name[0:index])
		temp_Name = temp_Name[index+1:]
	}
	names = append(names, temp_Name)
	
	//Create destinations
	temp_Destination := strings.Clone(packet.Destination_Str)
	for strings.Index(temp_Destination, "|") != -1{
		index := strings.Index(temp_Destination, "|")
		destinations = append(destinations, temp_Destination[0:index])
		temp_Destination = temp_Destination[index+1:]
	}
	destinations = append(destinations, temp_Destination)
		
	//Check if we have a cache hit within TTL (only relevant for interest packets)
	//Can have multiple cache hits for Multi_Flooding
	var temp_Cache_Store []*CACHE
	var name_Hit_Bool bool
	if packet.Interest_Data == 0{
		for x := 0; x < len(names); x++{
			temp_Cache, cache_Hit_Bool := check_Cache_Hit(curr_Node, names[x])
			if cache_Hit_Bool == true{ //cache hit
				temp_Cache_Store = append(temp_Cache_Store, temp_Cache)
				satisfied_Names = append(satisfied_Names, names[x])
			}
			if curr_Node.Data_Name == names[x]{ //at producer
				name_Hit_Bool = true
				satisfied_Names = append(satisfied_Names, names[x])
			}
		}		
	}
	
	//(if the names match or we have a cache hit), we are going to start going backwards
	if name_Hit_Bool == true || len(temp_Cache_Store) > 0{
		
		//For each cache hit, grab the data
		if len(temp_Cache_Store) > 0{ 
			for x := 0; x < len(temp_Cache_Store); x++{
				debug_Print("Reached Destination (cache hit) for name: " + temp_Cache_Store[x].Name + "\n")
				if curr_Node.Number == temp_Cache_Store[x].Number{ //if the producer cached its own data in its own cache, skip (should never happen)
					continue
				}
				temp_Packet := copy_Packet(packet)
				temp_Packet.Interest_Data = 1
				if temp_Packet.GP_Mode == 1{
					temp_Packet.GP_Mode = 0
				}
				temp_Packet.Index_Middle = len(temp_Packet.Path_Traversed)-1
				temp_Packet.Payload = temp_Cache_Store[x].Data 
				temp_Packet.Name = temp_Cache_Store[x].Name
				temp_Packet.Destination = temp_Cache_Store[x].Number
				//temp_Packet.Destination_Str = strconv.Itoa(temp_Cache_Store[x].Number)
				temp_Packet.Cache_Hit = 1
				next_Nodes = append(next_Nodes, find_Node(topology, temp_Packet.Previous))
				next_Packets = append(next_Packets, temp_Packet)
			}
		}
		
		//If you reach the producer, grab the data
		if name_Hit_Bool == true{ 
			debug_Print("Reached Destination (producer) for name: " + curr_Node.Data_Name + "\n")
			temp_Packet := copy_Packet(packet)
			temp_Packet.Interest_Data = 1
			if temp_Packet.GP_Mode == 1{
				temp_Packet.GP_Mode = 0
			}
			temp_Packet.Index_Middle = len(temp_Packet.Path_Traversed)-1
			temp_Packet.Payload = curr_Node.Data_Value 
			temp_Packet.Name = curr_Node.Data_Name
			temp_Packet.Destination = curr_Node.Number
			//temp_Packet.Destination_Str = strconv.Itoa(curr_Node.Number)
			next_Nodes = append(next_Nodes, find_Node(topology, temp_Packet.Previous))
			next_Packets = append(next_Packets, temp_Packet)
		}
		if multi == false{
			continue_Forward = false
		}
		
	} else if packet.Interest_Data == 1{ //if we are going back
		
		//Determine if we are going to cache the data here
		if determine_CACHE(caching_Strat) == true && curr_Node.Cache_Size != 0{
			cache_Data(curr_Node, packet, cache_Evict)	
		}
		if filo == true{ //these algorithms treat the PIT like a stack
			next_Nodes = return_FILO(curr_Node, packet, topology)
			next_Packets = append(next_Packets, copy_Packet(packet))
		} else{ //these algorithms forward to ALL entries in the PIT
			temp_next_Nodes, start, destination, collapsed := return_ALL(curr_Node, packet, topology)
			next_Nodes = temp_next_Nodes
			for x := 0; x < len(next_Nodes); x++{
				temp_Packet := copy_Packet(packet)
				temp_Packet.Start = start[x]
				if multi == false{
					temp_Packet.Destination = destination[x] //for Multi algorithms, we dont know the destination when we set the PIT entries, only when we get the packet
				}
				temp_Packet.Collapsed = collapsed[x]
				next_Packets = append(next_Packets, temp_Packet)
			}
		}
		continue_Forward = false
	} 
	
	if continue_Forward == true {//else, we need to find where to send the packet to, as we are not on the return trip (we can easily check this if there is a payload)
	
		//If we want to interest collapse
		drop_Packet := false
		if interest_Collapse == true{
			for x := 0; x < len(names); x++{
				temp_Dest, err := strconv.Atoi(destinations[x])
				check_Error(err) //check for errors
				drop_Packet = interest_Collapse_Helper(curr_Node, names[x], packet, filo, temp_Dest)
				if drop_Packet == true{
					satisfied_Names = append(satisfied_Names, names[x])
				}
			}
		}
		
		new_Concat_Name := ""
		if multi == true{
			names, new_Concat_Name, drop_Packet, destinations = update_Names(names, satisfied_Names, destinations)
		}
	
		//If we dont interest collapse / have unsatisfied names, find the neighbor
		if drop_Packet == false{
		
			if algorithm == "Flooding" || algorithm == "Multi_Flooding"{
				next_Node_Num := flooding_Helper(curr_Node, packet.Previous) //find all neighbors that arent the previous node or the current node
				
				//Update the PIT
				//Only log the incoming interface since there may be multiple outgoing interfaces
				curr_Node.Pit_Lock.Lock()
				
				for x := 0; x < len(names); x++{
					temp_Dest, err := strconv.Atoi(destinations[x])
					check_Error(err) //check for errors
					curr_Node.Pit = append(curr_Node.Pit, PIT{packet.Previous, -2, names[x], packet.Start, temp_Dest, -1}) //add it to the PIT
					curr_Node.Pit_Entries++
				}
				curr_Node.Pit_Lock.Unlock()

				for x := 0; x < len(next_Node_Num); x++{ //for each neighbor
					next_Nodes = append(next_Nodes, find_Node(topology, next_Node_Num[x])) //get the node
					if algorithm == "Multi_Flooding"{
						temp_Copy_Packet := copy_Packet(packet)
						temp_Copy_Packet.Name = new_Concat_Name
						next_Packets = append(next_Packets, temp_Copy_Packet) //only the unsatisfied names
					} else{
						next_Packets = append(next_Packets, copy_Packet(packet)) //only the unsatisfied names
					}
				}
			} else{
		
				var temp_Interest_Node []*Node
				var temp_Interest_Packet []Packet	
		
				//for each unsatisfied name
					//see what the next node is (next_i) according to the specific algorithm
						//if drop, consider it 'satisfied'
						//else, add to PIT
				for x := 0; x < len(names); x++{
					temp_Dest, err := strconv.Atoi(destinations[x])
					check_Error(err) //check for errors
					end := find_Node(topology, temp_Dest)
					neighbors := find_Neighbors(topology, curr_Node.Number)
					temp_Packet := copy_Packet(packet)
					temp_Packet.Name = names[x]
					temp_Packet.Destination = temp_Dest
					
					//Based on what algorithm we are doing, execute the helper function	
					next_Node_Num := -1
					just_In_Pressure := false
					if modified == false && temp_Packet.GP_Mode != 1{
						next_Node_Num, drop_Packet = gf_Helper(curr_Node, end, neighbors, distance_Formula) //Call the helper
					} else if modified == true && temp_Packet.GP_Mode != 1{
						next_Node_Num, drop_Packet = mgf_Helper(curr_Node, end, neighbors, temp_Packet.Previous, distance_Formula) //Call the helper
					} 
					next_Node := find_Node(topology, next_Node_Num)
					if gravity_pressure == true{
						if temp_Packet.GP_Mode == 0 && drop_Packet == true{
							temp_Distance := distance_Helper(end, curr_Node, distance_Formula)
							if temp_Distance < temp_Packet.GP_Distance || temp_Packet.GP_Distance == 0{
								temp_Packet.GP_Distance = temp_Distance
							}	
							temp_Packet.GP_Mode = 1 //enter pressure mode
							debug_Print("Entering Pressure Mode -- Distance: " + strconv.FormatFloat(temp_Packet.GP_Distance, 'g', -1, 64) + "\n")
						}
						drop_Packet = false
						if temp_Packet.GP_Mode == 1{ //if we are in pressure mode

							//Calculate the next node in pressure mode
							next_Node_Num = pressure_Mode_Helper(curr_Node, temp_Packet, topology, distance_Formula, temp_Dest)

							//Calculate the new distance 
							next_Node = find_Node(topology, next_Node_Num)
							distance := distance_Helper(end, next_Node, distance_Formula)

							//If the new distance is less, go back to gravity mode
							if distance < temp_Packet.GP_Distance{
								debug_Print("Leaving Pressure Mode -- Distance: " + strconv.FormatFloat(distance, 'g', -1, 64) + "\n")
								temp_Packet.GP_Mode = 0
							}
							temp_Packet.Pressure_Gauge = temp_Packet.Pressure_Gauge + 1 //increment the pressure gauge
							just_In_Pressure = true
						}
						if global == true || just_In_Pressure == true{
							//Track in pressure mode if just/currently in pressure mode
							//Globally track all nodes it has visited if G algorithm
							//Globally track how many times it visited each node if G algorithm
							visit_Before := false
							for x := 0; x < len(temp_Packet.Nodes_Visit); x++{
								if temp_Packet.Nodes_Visit[x] == next_Node_Num{
									temp_Packet.Num_Visit[x] = temp_Packet.Num_Visit[x] + 1
									visit_Before = true
									break
								}
							}
							if visit_Before == false{
								temp_Packet.Nodes_Visit = append(temp_Packet.Nodes_Visit, next_Node_Num)
								temp_Packet.Num_Visit = append(temp_Packet.Num_Visit, 1)
							}
						}
					}
					if next_Node_Num == -1 && drop_Packet == false{
						fmt.Println("Algorithm Failed! Exiting.")
						os.Exit(1)
					}
					if drop_Packet == false{
						curr_Node.Pit_Lock.Lock()
						curr_Node.Pit = append(curr_Node.Pit, PIT{temp_Packet.Previous, next_Node_Num, names[x], temp_Packet.Start, temp_Dest, -1}) //update the Pit if we did the algorithm
						temp_Interest_Node = append(temp_Interest_Node, next_Node)
						temp_Interest_Packet = append(temp_Interest_Packet, temp_Packet)
						curr_Node.Pit_Entries++
						curr_Node.Pit_Lock.Unlock()
					}
									
				}
				next_Nodes_Temp, next_Packets_Temp := combine_Packet(temp_Interest_Node, temp_Interest_Packet)
				next_Nodes = append(next_Nodes, next_Nodes_Temp...)
				next_Packets = append(next_Packets, next_Packets_Temp...)	
			}
		}
	}
	if len(next_Nodes) != len(next_Packets){
		fmt.Println("Error! Different number of nodes and packets!")
		os.Exit(1)
	}
	if len(next_Nodes) == 0{ //no next nodes
		debug_Print("Dropped Packet\n")
		packet.Name = ""
		packet_Contents <- packet
	} else{
		open_Processes_Lock.Lock()
		open_Processes = open_Processes + len(next_Nodes) - 1
		if open_Processes > open_Processes_Max{
			open_Processes_Max = open_Processes
		}
		open_Processes_Lock.Unlock()
	}
	
	for x := 0; x < len(next_Nodes); x++{
		if !((time.Now()).Before((packet.Start_Time).Add(packet.TTL)) || packet.TTL < time.Duration(0)){ //if the TTL has expired, throw out the packet
			debug_Print("Packet TTL Expired\n")
			packet.Name = ""
			packet_Contents <- packet
			continue
		}
		if next_Nodes[x].Number == -1{
			debug_Print("Successful Packet\n")
			packet_Contents <- next_Packets[x]
			continue
		}
		temp_Path := make([]int, len(next_Packets[x].Path_Traversed))
		copy(temp_Path, next_Packets[x].Path_Traversed) //deep copy so we dont affect path
		temp_Path = append(temp_Path, next_Nodes[x].Number)
		c, err := net.Dial("tcp", next_Nodes[x].IP + ":" + strconv.Itoa(next_Nodes[x].Port)) //send the packet to the next node
		check_Error(err) //check for errors
		debug_Print("Node " + strconv.Itoa(curr_Node.Number) + " sent a message to Node " + strconv.Itoa(next_Nodes[x].Number) + " \n")
		next_Packets[x].Path_Traversed = temp_Path
		next_Packets[x].Previous = curr_Node.Number
		send_Msg(next_Packets[x], c)
		
		//Could be a global, but do it per node to prevent bottle necks
		curr_Node.Hops_Lock.Lock()
		curr_Node.Hops++
		curr_Node.Hops_Lock.Unlock()
	}
}
//----------------------------------------------

func main() {

	//----------------------------------------------
	//Step 0: Read in the command line args, define metrics
	//----------------------------------------------
	data_Input, ip, port, distance_Formula, average_Num, algorithms, one_at_a_time, interest_Collapse, caching_Strat, cache_Evict, TTL, cache_Size, start_Nodes, end_Nodes, output := read_args()
	metric_names := []string{"Latency", "Packet Loss", "Path Stretch", "Packet Size", "Pressure count", "Pressure percentage", "Nodes visited", "Hops collapsed", "Cache Hit Ratio", "Number Cache hits", "Num Collapsed", "Max Concurrent", "Hops Made", "Producers Visited", "Ratio interest/producer", "PIT Entries", "Hanging PIT Entries"}
	offset := 8 //number of metrics metrics_calc returns
	var final_results [][]float64
	
	//----------------------------------------------
	//Step 1: For each run
	//----------------------------------------------
	for a := 0; a < average_Num; a++{
		fmt.Printf("\nLoop: %d\n", a)
		
		//----------------------------------------------
		//Step 2: For each algorithm
		//----------------------------------------------
		var algorithm_storage []float64
		for b := 0; b < len(algorithms); b++{

			//----------------------------------------------
			//Step 3: For each topology
			//----------------------------------------------
			filo, _, _, gravity_pressure, multi := algorithm_Check(algorithms[b])
			files, err := ioutil.ReadDir(data_Input)
			check_Error(err) //check for errors
			
			var topology_metrics [][]float64
			for x := 0; x < len(metric_names); x++{
				topology_metrics = append(topology_metrics, []float64{})
			}	
			for _, file_From_Files := range files{

				//----------------------------------------------
				//Step 4: Have the “driver” read in a topology file that will list each node, with each node's latitude, longitude, and neighbors
				//----------------------------------------------
				topology := read_Topology(data_Input+file_From_Files.Name(), ip, port, cache_Size)
				debug_Print("Topology established for: " + file_From_Files.Name() + "\n")
				//----------------------------------------------

				//----------------------------------------------
				//Step 5: The driver will then set up a number of workers equal to the number of nodes in the topology
				// 	  Each worker will listen on its own port, and be assigned its hyperbolic coordinates
				//----------------------------------------------
				shutdown_ACK := make(chan int, len(topology)) //channel for passing info
				packet_Contents := make(chan Packet, len(topology)*len(topology)*len(topology)) //channel for retrieving the return value -- topology^3 to ensure no blocking	
				for x := 0; x < len(topology); x++{
					go setup_Node(topology[x], shutdown_ACK, packet_Contents, topology, distance_Formula, algorithms[b], caching_Strat, cache_Evict, interest_Collapse) //setup the nodes
				}
				//wait for all nodes to be setup
				for x := 0; x < len(topology); x++{
					<-shutdown_ACK
				}
				debug_Print("All nodes setup for: " + file_From_Files.Name() + "\n")
				//----------------------------------------------

				//----------------------------------------------
				//Step 6: The driver will then iteratively select a node combination for each node pairing in the topology
				//----------------------------------------------
				cache_Hit = 0 //0 out cache_hit because its a global
				num_Collapsed = 0 //0 out num_collapsed because its a global
				
				//Setting up the packets array -- can index packets[start_node][end_node]
				var packets [][]Packet
				for y := 0; y < len(topology); y++{
					var temp_Packets []Packet
					for z := 0; z < len(topology); z++{
						//Each node will have its own named item being requested of it (hence why we append the number of the node to the desintation)
						message_To_Send := Packet{-1, y, z, strconv.Itoa(z), 0, "/ucla/videos/demo.mpg/1/" + strconv.Itoa(z), "" , time.Now(), time.Duration(time.Second * time.Duration(TTL)), 0, []int{}, []int{}, -1.0, []int{y}, 0, -1, time.Now(), -1, 0}
						if gravity_pressure == false{
							message_To_Send.GP_Mode = -1
							message_To_Send.Pressure_Gauge = -1
						} 
						temp_Packets = append(temp_Packets, message_To_Send)	
					}
					packets = append(packets, temp_Packets)
				}
				
				//Setup all node pairing 
				local_open_Processes := 1.0
				
				//Select nodes
				start_nodes := random_Selection(topology, start_Nodes)
				end_nodes := random_Selection(topology, end_Nodes)
				num_Packets := float64(len(start_nodes)*len(end_nodes))
				
				//Prints selected nodes
				if debug == true{
					print_String := ""
					for y := 0; y < len(start_nodes); y++{
						print_String = print_String + strconv.Itoa(start_nodes[y])
						if y != len(start_nodes)-1{
							print_String = print_String + ","
						} else{
							print_String = print_String + "\n"
						}
					}
					debug_Print("\nStarting Nodes: " + print_String)
					print_String = ""
					for y := 0; y < len(end_nodes); y++{
						print_String = print_String + strconv.Itoa(end_nodes[y])
						if y != len(end_nodes)-1{
							print_String = print_String + ","
						} else{
							print_String = print_String + "\n"
						}
					}
					debug_Print("Destination Nodes: " + print_String + "\n")
				}
				
				//Configure the destination string
				dest_String := ""
				for y := 0; y < len(end_nodes); y++{						
					dest_String = dest_String + strconv.Itoa(end_nodes[y])
					if y != len(end_nodes)-1{
						dest_String = dest_String + "|"
					}
				}
				
				//Multi Specific
				var Multi_Packets []Packet
				if multi == true{
					num_Packets = float64(len(start_nodes))
					
					//Stitch together the name
					name := ""
					for z := 0; z < len(end_nodes); z++{
						name = name + "/ucla/videos/demo.mpg/1/" + strconv.Itoa(end_nodes[z])
						if z != len(end_nodes)-1{
							name = name + "|"
						}
					}
					
					//Create new packets to send
					for y := 0; y < len(start_nodes); y++{
						message_To_Send := Packet{-1, start_nodes[y], -2, dest_String, 0, name, "" , time.Now(), time.Duration(time.Second * time.Duration(TTL)), 0, []int{}, []int{}, -1.0, []int{start_nodes[y]}, 0, -1, time.Now(), -1, 0}	
						if gravity_pressure == false{
							message_To_Send.GP_Mode = -1
							message_To_Send.Pressure_Gauge = -1
						}
						Multi_Packets = append(Multi_Packets, message_To_Send) 
					}
				}
				
				if one_at_a_time == false{
					local_open_Processes = num_Packets
					open_Processes = int(num_Packets)
					open_Processes_Max = open_Processes
				}
				for y := 0; y < len(start_nodes); y++{
					for z := 0; z < len(end_nodes); z++{
					
						if multi == true && z > 0{ //only send 1 packet per start node 
							break
						}
					
						if one_at_a_time == true{
							open_Processes = 1
							open_Processes_Max = open_Processes
						}
						start := find_Node(topology, start_nodes[y]) //start node
						//----------------------------------------------
						//Step 7: For each node pairing, each forwarding algorithm will be used to send a packet across the topology based on those hyperbolic coordinates
						//----------------------------------------------
						c, err := net.Dial("tcp", start.IP + ":" + strconv.Itoa(start.Port)) //Start by sending a packet to the start node
						check_Error(err) //check for errors
						
						if multi == true{
							debug_Print("Sending a message to Node " + strconv.Itoa(start_nodes[y]) + " telling it to go to multicast \n")
							time_Now := time.Now()
							for a := 0; a < len(end_nodes); a++{
								packets[start_nodes[y]][end_nodes[a]].Start_Time = time_Now
							}
							send_Msg(Multi_Packets[y], c)
						}else{
							debug_Print("Sending a message to Node " + strconv.Itoa(start_nodes[y]) + " telling it to go to Node " + strconv.Itoa(end_nodes[z]) +" \n")
							packets[start_nodes[y]][end_nodes[z]].Start_Time = time.Now()
							send_Msg(packets[start_nodes[y]][end_nodes[z]], c)
						}
						
						
						if one_at_a_time == true{
							packets = wait_Main(packet_Contents, packets)
							debug_Print("Test for Node " + strconv.Itoa(start_nodes[y]) + " and Node " + strconv.Itoa(end_nodes[z]) +" done\n\n")
							
							//Set the default Pit values (should only be needed for dropped packets -- GF or mGF) 
							for a := 0; a < len(topology); a++{
								topology[a].Pit = []PIT{} 
							}
							//Clears the cache
							for a := 0; a < len(topology); a++{
								topology[a].Cache = []CACHE{} 
							}
						}
					}
				}
				if one_at_a_time == false{
					packets = wait_Main(packet_Contents, packets)
					debug_Print("Test for all nodes done\n\n")
				}
				//----------------------------------------------
				//Step 8: Check assertions
				//----------------------------------------------
				ttl_expired := false
				for y := 0; y < len(start_nodes); y++{
					for z := 0; z < len(end_nodes); z++{
						packet := packets[start_nodes[y]][end_nodes[z]]
						
						//Assertion we didnt drop a packet for GPGF algorithms unless it was TTL based
						if !((time.Now()).Before((packet.Start_Time).Add(packet.TTL)) || packet.TTL < time.Duration(0)){
							ttl_expired = true
						} else{
							if gravity_pressure == true && packet.Payload == ""{ 
								fmt.Println("Dropped Packet for GPGF algorithm!")
								os.Exit(1)
							}
						}
						if filo == true && packet.Payload != ""{ //Assertion we are following the reverse of the taken path for FILO algorithms
							reverse_Check(packet.Path_Traversed)
						}
					}
				}
				
				//Assertion the PIT is empty (only if a packet is not dropped and a TTL has not expired)
				if gravity_pressure == true && ttl_expired == false{
					empty_PIT_Check(topology)
				}
				
				//----------------------------------------------
				//Step 9: An analysis will be calculated and output, showing things like latency, packet loss, path stretch per topology
				//GOlang is doing TRUNCATION, not ROUNDING - error of .001 at most
				//----------------------------------------------
				var node_pairing_metrics[][]float64
				for x := 0; x < len(metric_names); x++{
					node_pairing_metrics = append(node_pairing_metrics, []float64{})
				}

				for y := 0; y < len(start_nodes); y++{
					for z := 0; z < len(end_nodes); z++{
						if start_nodes[y] == end_nodes[z]{ //A->A is a meaningless metric
							continue
						}

						//Calculate Metrics
						if packets[start_nodes[y]][end_nodes[z]].Payload != ""{ //the packet was not dropped
							temp_metrics := metrics_Calc(packets[start_nodes[y]][end_nodes[z]], topology)
							for x := 0; x < len(temp_metrics); x++{
								if (x == 4 && temp_metrics[4] == -1.0) || (x == 7 && temp_metrics[7] == 0.0){
									continue
								}
								node_pairing_metrics[x] = append(node_pairing_metrics[x], temp_metrics[x])
							}
						} else{ //the packet was dropped
							node_pairing_metrics[1] = append(node_pairing_metrics[1], 0)
						}
					}
				}
				
				for x := 0; x < len(metric_names); x++{
					if x == offset+1{ //raw cache hits
						topology_metrics[x] = append(topology_metrics[x], float64(cache_Hit))
					} else if x == offset+2{ //num collapsed
						topology_metrics[x] = append(topology_metrics[x], float64(num_Collapsed)/num_Packets)
					} else if x == offset+3{ //concurrent packets
						topology_metrics[x] = append(topology_metrics[x], float64(open_Processes_Max)/local_open_Processes)
					} else if x == offset+4{ //hops
						hops_Made := 0
						for x := 0; x < len(topology); x++{
							hops_Made = hops_Made + topology[x].Hops	
						}
						topology_metrics[x] = append(topology_metrics[x], float64(hops_Made)/num_Packets)
					} else if x == offset+5{ //satisfied
						if multi == true{
							topology_metrics[x] = append(topology_metrics[x], float64(len(end_nodes)))
						} else{
							topology_metrics[x] = append(topology_metrics[x], float64(1))
						}
					} else if x == offset+6{ //Hops per satisfied
						hops_Made := 0
						for x := 0; x < len(topology); x++{
							hops_Made = hops_Made + topology[x].Hops	
						}
						hops_Per_Interest := float64(float64(hops_Made)/num_Packets)
						producers_Per_Interest := float64(1)
						if multi == true{
							producers_Per_Interest = float64(len(end_nodes))
						}
						topology_metrics[x] = append(topology_metrics[x], hops_Per_Interest/producers_Per_Interest)
					} else if x == offset+7{ //PIT entries
						PIT_Entries := 0
						for x := 0; x < len(topology); x++{
							PIT_Entries = PIT_Entries + topology[x].Pit_Entries	
						}
						topology_metrics[x] = append(topology_metrics[x], float64(PIT_Entries)/float64(len(topology)))
					} else if x == offset+8{ //Hanging PIT
						hanging_PIT := 0
						for x := 0; x < len(topology); x++{
							hanging_PIT = hanging_PIT + len(topology[x].Pit)
						}
						topology_metrics[x] = append(topology_metrics[x], float64(hanging_PIT))
					} else{
						topology_metrics[x] = append(topology_metrics[x], avg(node_pairing_metrics[x]))
						if stdCheck(node_pairing_metrics[x], 3) == false{
							debug_Print("Warning: " + metric_names[x] + " has outlier values!\n")
						}
					}
				}
				fmt.Println("All tests complete for: " + file_From_Files.Name() + "--" + algorithms[b])
				
				//----------------------------------------------
				//Step 10: Cleanup an experiment
				//----------------------------------------------
				//send the shutdown signal
				debug_Print("Sending shutdown signals\n")
				for x := 0; x < len(topology); x++{
					c, err := net.Dial("tcp", topology[x].IP + ":" + strconv.Itoa(topology[x].Port)) //Start by sending a packet to the start node
					check_Error(err) //check for errors
					send_Msg(Packet{Previous: -999, Destination: -999, Interest_Data: 2, Name: "Shutdown"}, c)
				}
				//wait for all nodes to shutdown
				for x := 0; x < len(topology); x++{
					<-shutdown_ACK
				}
				debug_Print("All nodes shutdown for: " + file_From_Files.Name() + "\n")
			}
			//----------------------------------------------
			//Step 11: Store results for each algorithm
			//----------------------------------------------
			for x := 0; x < len(metric_names); x++{
				if x == 0 || x == 2 || x == 3 || x == offset+1 || x == offset+3 || x == offset+4 || x == offset+5 || x == offset+6 || x == offset+7 || x == offset+8{
					algorithm_storage = append(algorithm_storage, avg(topology_metrics[x]))
				} else{
					algorithm_storage = append(algorithm_storage, avg(topology_metrics[x])*100)
				}
			}
		}
		
		//----------------------------------------------
		//Step 12: Store results for each run
		//----------------------------------------------	
		final_results = append(final_results, algorithm_storage)
	}
	
	//----------------------------------------------
	//Step 13: An analysis will be calculated and output, showing things like latency, packet loss, path stretch averaged for all topologies
	//GOlang is doing TRUNCATION, not ROUNDING - error of .001 at most
	//----------------------------------------------
	var final_Values []float64
	for x := 0; x < len(final_results[0]); x=x+1{	
		var temp_Final_Values []float64
		for y := 0; y < len(final_results); y++{
			temp_Final_Values = append(temp_Final_Values, final_results[y][x])
		}
		final_Values = append(final_Values, avg(temp_Final_Values))
	}
	
	output_String := ""
	if one_at_a_time == true{
		output_String = output_String + "One at a time - "
	} else{
		output_String = output_String + "All at once - "
	}
	if distance_Formula == "hyperbolic"{
		output_String = output_String + "Hyperbolic - "
	} else{
		output_String = output_String + "Euclidean - "
	}
	if interest_Collapse == true{
		output_String = output_String + "interest collapsing - "
	} else{
		output_String = output_String + "no interest collapsing- "
	}
	if caching_Strat == "NONE" || TTL == 0 || cache_Size == 0{
		output_String = output_String + "no caching - "
	} else{
		output_String = output_String + "caching - "
	}
	output_String = output_String + "Start = " + strconv.FormatFloat(float64(start_Nodes), 'f', -1, 64) + " - "
	output_String = output_String + "End = " + strconv.FormatFloat(float64(end_Nodes), 'f', -1, 64) + " - "
	output_String = output_String + "averaged over " + strconv.Itoa(average_Num) + " result(s)"
	
	output_Metrics := []string{"Average latency in seconds", "Average percentage of packets successfully delivered", "Average path stretch", "Average packet size in bytes", "Average percentage of time spent in pressure mode on initial trip if entered pressure at all", "Average packets that used pressure mode as a percentage of total packets", "Average percentage of nodes visited", "Average number of hops we saved by interest collapsing as a percentage of total hops if we collapsed", "Average cache hit ratio for successful packets based on -- Caching Strat: " + caching_Strat + " -- Cache Eviction: " + cache_Evict + " -- TTL: " + strconv.FormatFloat(float64(TTL), 'f', -1, 64) + " -- Cache Size: " + strconv.FormatFloat(float64(cache_Size), 'f', -1, 64), "Raw Number of Cache Hits", "Average number of times we interest collapsed as a percentage of total packets", "Ratio of the average maximum number of concurrent packets over number of starting maximum packets", "Average number of hops per interest", "Average number of satisfied requests per interest", "Ratio of hops per interest / satisfied requests per interest", "Average number of PIT entries per node", "Number of Hanging PIT Entries"}
	
	//Creating a file
	file, err := os.Create(output)
	check_Error(err)
	defer file.Close()
	
	fmt.Printf("\n\n" + output_String)
	_, err = io.WriteString(file, output_String)
	check_Error(err)
	for x := 0; x < len(algorithms); x++{
		fmt.Printf("\t" + algorithms[x])
		_, err = io.WriteString(file, "," + algorithms[x])
		check_Error(err)
	}
	fmt.Printf("\n")
	_, err = io.WriteString(file, "\n")
	check_Error(err)
	
	for x := 0; x < len(metric_names); x++{
		fmt.Printf(output_Metrics[x] + "\t")
		_, err = io.WriteString(file, output_Metrics[x] + ",")
		check_Error(err)
		for y := x; y < len(final_Values); y=y+len(metric_names){	
			fmt.Printf("%.6v", final_Values[y])
			ratio := math.Pow(10, float64(6))
			_, err = io.WriteString(file, strconv.FormatFloat(float64(math.Round(final_Values[y]*ratio)/ratio), 'f', -1, 64))
			check_Error(err)
			if y+len(metric_names) < len(final_Values){
				fmt.Printf("\t")
				_, err = io.WriteString(file, ",")
				check_Error(err)
			}
		}
		fmt.Printf("\n")
		_, err = io.WriteString(file, "\n")
		check_Error(err)
	}
	fmt.Printf("\n")
	check_Error(err)
}
