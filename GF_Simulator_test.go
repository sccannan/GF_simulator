package main

//----------------------------------------------
// Imports
//----------------------------------------------
import (
    "testing"
    "math"
    "os"
    "os/exec"
    "net"
    "time"
    "math/rand"
    "strconv"
    "reflect"
)
//----------------------------------------------

//--------------------------------------------------------------------------------------------
// List of functions
//--------------------------------------------------------------------------------------------
// read_args 			-- untested (reads in args)
// check_Error 			-- untested (simple, no need)
// debug_Print 			-- untested (simple, no need)
// hyperbolic_Distance 		-- tested
// euclidean_Distance 		-- tested
// distance_Helper 		-- tested
// dijkstras 			-- tested
// avg 				-- tested
// std 				-- tested
// stdCheck 			-- tested
// gf_Helper 			-- tested
// mgf_Helper 			-- tested
// send_Msg 			-- tested
// receive_Msg 			-- tested
// find_Node 			-- tested
// find_Neighbors 		-- tested
// index_PIT 			-- tested
// remove_PIT 			-- tested
// index_CACHE 			-- tested
// remove_CACHE 		-- tested
// determine_CACHE 		-- tested
// evict_CACHE 			-- tested
// reverse_Check 		-- tested
// empty_PIT_Check 		-- tested
// calc_Size 			-- tested
// read_Topology		-- tested
// metrics_Calc 		-- tested
// algorithm_Check 		-- tested
// check_Cache_Hit		-- tested
// cache_Data			-- tested 
// pressure_Mode_Helper		-- tested 
// interest_Collapse_Helper	-- tested 
// return_FILO			-- tested 
// return_ALL			-- tested 
// random_Selection		-- tested
// flooding_Helper		-- tested
// copy_Packet			-- tested
// update_Names			-- tested
// combine_Names		-- tested
// compare_Packet		-- tested
// combine_Packet		-- tested
// wait_Main 			-- untested (uses globals to terminate)
// setup_Node 			-- untested (server code)
// serviceConnection 		-- untested (handler for packets)
// main 			-- untested (main)
//--------------------------------------------------------------------------------------------

//--------------------------------------------------------------------------------------------
// hyperbolic_Distance(r_A float64, r_B float64, theta_A float64, theta_B float64) float64
//--------------------------------------------------------------------------------------------

//----------------------------------------------
// Tests to make sure hyperbolic distance is calculating correctly
//----------------------------------------------
func TestHyperbolic(t *testing.T) {
	r_A := 1.0
	r_B := 2.0
	theta_A := 1.0
	theta_B := 2.0
	actual := hyperbolic_Distance(r_A, r_B, theta_A, theta_B)
	want := 1.925576229269065
	if actual != want{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
}
//----------------------------------------------

//----------------------------------------------
// Tests to make sure hyperbolic distance is calculating correctly if dtheta=0
//----------------------------------------------
func TestHyperbolic_dtheta(t *testing.T) {
	r_A := 1.0
	r_B := 3.0
	theta_A := 1.0
	theta_B := 1.0
	actual := hyperbolic_Distance(r_A, r_B, theta_A, theta_B)
	want := 2.0
	if actual != want{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
}
//----------------------------------------------

//--------------------------------------------------------------------------------------------
// euclidean_Distance(lat_A float64, lat_B float64, long_A float64, long_B float64) float64
//--------------------------------------------------------------------------------------------

//----------------------------------------------
// Tests to make sure hyperbolic distance is calculating correctly
//----------------------------------------------
func TestEuclidean(t *testing.T) {
	lat_A := 1.0
	lat_B := 3.0
	long_A := 1.0
	long_B := 2.0
	actual := euclidean_Distance(lat_A, lat_B, long_A, long_B)
	want := 2.23606797749979
	if actual != want{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
}
//----------------------------------------------

//--------------------------------------------------------------------------------------------
// distance_Helper(A *Node, B *Node, distance_Formula string) float64{
//--------------------------------------------------------------------------------------------

//----------------------------------------------
// Tests to make sure hyperbolic distance is calculating correctly via the helper
//----------------------------------------------
func TestHelperHyperbolic(t *testing.T) {
	node_A := &Node{Lat: 1, Long: 1, Theta: 1.0, R: 1.0} //create a struct for node
	node_B := &Node{Lat: 3, Long: 2, Theta: 2.0, R: 2.0} //create a struct for node
	actual := distance_Helper(node_A, node_B, "hyperbolic")
	want := 1.925576229269065
	if actual != want{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
}
//----------------------------------------------

//----------------------------------------------
// Tests to make sure euclidean distance is calculating correctly via the helper
//----------------------------------------------
func TestHelperEuclidean(t *testing.T) {
	node_A := &Node{Lat: 1, Long: 1, Theta: 1.0, R: 1.0} //create a struct for node
	node_B := &Node{Lat: 3, Long: 2, Theta: 2.0, R: 2.0} //create a struct for node
	actual := distance_Helper(node_A, node_B, "euclidean")
	want := 2.23606797749979
	if actual != want{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
}
//----------------------------------------------

//----------------------------------------------
// Tests to make sure hyperbolic distance is calculating correctly if NaN values
//----------------------------------------------
func TestHelper_NAN(t *testing.T) {
	node_A := &Node{Lat: 1, Long: 1, Theta: 1.0, R: math.NaN()} //create a struct for node
	node_B := &Node{Lat: 3, Long: 2, Theta: 2.0, R: 2.0} //create a struct for node
	actual := distance_Helper(node_A, node_B, "hyperbolic")
	want := math.MaxFloat64
	if actual != want{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
	
	node_A = &Node{Lat: 1, Long: 1, Theta: 1.0, R: 1} //create a struct for node
	node_B = &Node{Lat: 3, Long: 2, Theta: 2.0, R: math.NaN()} //create a struct for node
	actual = distance_Helper(node_A, node_B, "hyperbolic")
	want = math.MaxFloat64
	if actual != want{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
	
	node_A = &Node{Lat: 1, Long: 1, Theta: math.NaN(), R: 1} //create a struct for node
	node_B = &Node{Lat: 3, Long: 2, Theta: 2.0, R: 2.0} //create a struct for node
	actual = distance_Helper(node_A, node_B, "hyperbolic")
	want = math.MaxFloat64
	if actual != want{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
	
	node_A = &Node{Lat: 1, Long: 1, Theta: 1.0, R: 1} //create a struct for node
	node_B = &Node{Lat: 3, Long: 2, Theta: math.NaN(), R: 2.0} //create a struct for node
	actual = distance_Helper(node_A, node_B, "hyperbolic")
	want = math.MaxFloat64
	if actual != want{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
	
	node_A = &Node{Lat: math.NaN(), Long: 1, Theta: 1.0, R: 1} //create a struct for node
	node_B = &Node{Lat: 3, Long: 2, Theta: 2.0, R: 2.0} //create a struct for node
	actual = distance_Helper(node_A, node_B, "euclidean")
	want = math.MaxFloat64
	if actual != want{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
	
	node_A = &Node{Lat: 1, Long: 1, Theta: 1.0, R: 1} //create a struct for node
	node_B = &Node{Lat: math.NaN(), Long: 2, Theta: 2.0, R: 2.0} //create a struct for node
	actual = distance_Helper(node_A, node_B, "euclidean")
	want = math.MaxFloat64
	if actual != want{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
	
	node_A = &Node{Lat: 1, Long: math.NaN(), Theta: 1.0, R: 1} //create a struct for node
	node_B = &Node{Lat: 3, Long: 2, Theta: 2.0, R: 2.0} //create a struct for node
	actual = distance_Helper(node_A, node_B, "euclidean")
	want = math.MaxFloat64
	if actual != want{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
	
	node_A = &Node{Lat: 1, Long: 1, Theta: 1.0, R: 1} //create a struct for node
	node_B = &Node{Lat: 3, Long: math.NaN(), Theta: 2.0, R: 2.0} //create a struct for node
	actual = distance_Helper(node_A, node_B, "euclidean")
	want = math.MaxFloat64
	if actual != want{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
}
//----------------------------------------------

//----------------------------------------------
// Tests to make sure unknown distance function is correct
//----------------------------------------------
func TestHelperUnknown(t *testing.T) {

	//Run code via cmd
	if os.Getenv("FLAG") == "1" {
		node_A := &Node{Lat: 1, Long: 1, Theta: 1.0, R: 1.0} //create a struct for node
		node_B := &Node{Lat: 3, Long: 2, Theta: 2.0, R: 2.0} //create a struct for node
		_ = distance_Helper(node_A, node_B, "unknown")
		return
	}
	
	//Execute the code via exec to check error code
	cmd := exec.Command(os.Args[0], "-test.run=TestHelperUnknown")
	cmd.Env = append(os.Environ(), "FLAG=1")
	output, err := cmd.Output() //grab error code and output
	
	//Compare error codes
	actual := err.Error()
	want := "exit status 1"
	if actual != want{       
		t.Errorf("Actual: %s -- Want: %s", actual, want)
	}
	
	//Compare output
	actual = string(output)
	want = "Unrecognized distance formula! Exiting.\n"
	if actual != want{       
		t.Errorf("Actual: %s -- Want: %s", actual, want)
	}
}
//----------------------------------------------

//--------------------------------------------------------------------------------------------
// dijkstras(y int, z int, topology []*Node) int
//--------------------------------------------------------------------------------------------

//----------------------------------------------
// Tests to make sure dijkstras works at 0 hops difference
//----------------------------------------------
func TestDijkstras_Same(t *testing.T) {
	var topology []*Node
	topology = append(topology, &Node{Number: 0, Neighbors: []int{0}})
	actual := dijkstras(0, 0, topology)
	want := 0
	if actual != want{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
}
//----------------------------------------------

//----------------------------------------------
// Tests to make sure dijkstras works at N hops difference
//----------------------------------------------
func TestDijkstras(t *testing.T) {
	var topology []*Node
	topology = append(topology, &Node{Number: 0, Neighbors: []int{0,1,2}})
	topology = append(topology, &Node{Number: 1, Neighbors: []int{0,1,2}})
	topology = append(topology, &Node{Number: 2, Neighbors: []int{0,1,2,3}})
	topology = append(topology, &Node{Number: 3, Neighbors: []int{2,3,4,5}})
	topology = append(topology, &Node{Number: 4, Neighbors: []int{3,4,5}})
	topology = append(topology, &Node{Number: 5, Neighbors: []int{3,4,5}})
	actual := dijkstras(0, 5, topology)
	want := 3
	if actual != want{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
}
//----------------------------------------------

//----------------------------------------------
// Tests to make sure dijkstras works when there is no path
//----------------------------------------------
func TestDijkstras_Fail(t *testing.T) {
	var topology []*Node
	topology = append(topology, &Node{Number: 0, Neighbors: []int{0,1,2}})
	topology = append(topology, &Node{Number: 1, Neighbors: []int{0,1,2}})
	topology = append(topology, &Node{Number: 2, Neighbors: []int{0,1,2,3}})
	topology = append(topology, &Node{Number: 3, Neighbors: []int{2,3,4,5}})
	topology = append(topology, &Node{Number: 4, Neighbors: []int{3,4,5}})
	topology = append(topology, &Node{Number: 5, Neighbors: []int{3,4,5}})
	topology = append(topology, &Node{Number: 6, Neighbors: []int{}})
	actual := dijkstras(0, 6, topology)
	want := -1
	if actual != want{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
}
//----------------------------------------------

//--------------------------------------------------------------------------------------------
// avg(arr []float64) float64
//--------------------------------------------------------------------------------------------

//----------------------------------------------
// Tests to make sure avg works
//----------------------------------------------
func TestAverage(t *testing.T) {
	values := []float64{46,69,32,60,52,41}
	actual := avg(values)
	want := 50.0
	if actual != want{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
}
//----------------------------------------------

//----------------------------------------------
// Tests to make sure avg works if len=0
//----------------------------------------------
func TestAverage_Len_0(t *testing.T) {
	values := []float64{}
	actual := avg(values)
	want := 0.0
	if actual != want{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
}
//----------------------------------------------

//--------------------------------------------------------------------------------------------
// std(arr []float64) float64
//--------------------------------------------------------------------------------------------

//----------------------------------------------
// Tests to make sure avg works
//----------------------------------------------
func TestSTD(t *testing.T) {
	values := []float64{46,69,32,60,52,41}
	actual := std(values)
	want := 13.311649033834989
	if actual != want{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
}
//----------------------------------------------

//----------------------------------------------
// Tests to make sure avg works if len=0
//----------------------------------------------
func TestSTD_Len_0(t *testing.T) {
	values := []float64{}
	actual := avg(values)
	want := 0.0
	if actual != want{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
}
//----------------------------------------------

//--------------------------------------------------------------------------------------------
// stdCheck(arr []float64, std_range int) bool 
//--------------------------------------------------------------------------------------------

//----------------------------------------------
// Tests to make sure stdCheck works if true
//----------------------------------------------
func TestSTDCheck_True(t *testing.T) {
	values := []float64{46,69,32,60,52,41}
	actual := stdCheck(values, 2)
	want := true
	if actual != want{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
}
//----------------------------------------------

//----------------------------------------------
// Tests to make sure stdCheck works if false
//----------------------------------------------
func TestSTDCheck_False(t *testing.T) {
	values := []float64{46,69,32,60,52,10000}
	actual := stdCheck(values, 2)
	want := false
	if actual != want{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
}
//----------------------------------------------

//----------------------------------------------
// Tests to make sure stdCheck works if true
//----------------------------------------------
func TestSTDCheck_Len_0(t *testing.T) {
	values := []float64{}
	actual := stdCheck(values, 3)
	want := true
	if actual != want{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
}
//----------------------------------------------

//--------------------------------------------------------------------------------------------
// gf_Helper(curr_Node *Node, end *Node, neighbors []*Node, distance_Formula string) (int, bool){
//--------------------------------------------------------------------------------------------

//----------------------------------------------
// Tests to make sure gf_helper works with immediate neighbor
//----------------------------------------------
func TestGF_Helper_Immediate(t *testing.T) {
	var neighbors []*Node
	node_A := &Node{Number: 0, Neighbors: []int{0,1}}
	node_B := &Node{Number: 1, Neighbors: []int{0,1}}
	neighbors = append(neighbors, node_A)
	neighbors = append(neighbors, node_B)
	actual, actual_2 := gf_Helper(node_A, node_B, neighbors, "euclidian")
	want := 1
	if actual != want{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
	want_2 := false
	if actual_2 != want_2{       
		t.Errorf("Actual: %v -- Want: %v", actual_2, want_2)
	}
}
//----------------------------------------------

//----------------------------------------------
// Tests to make sure gf_helper works with non-immediate neighbor
//----------------------------------------------
func TestGF_Helper_Non_Immediate(t *testing.T) {
	var neighbors []*Node
	node_A := &Node{Number: 0, Neighbors: []int{0,1,2}, Lat: 0, Long: 1}
	node_B := &Node{Number: 1, Neighbors: []int{0,1,2}, Lat: 1, Long: 1}
	node_C := &Node{Number: 2, Neighbors: []int{0,1,2}, Lat: 2, Long: 1}
	destination := &Node{Number: 3, Neighbors: []int{0,1,2}, Lat: 3, Long: 1}
	neighbors = append(neighbors, node_A)
	neighbors = append(neighbors, node_B)
	neighbors = append(neighbors, node_C)
	actual, actual_2 := gf_Helper(node_A, destination, neighbors, "euclidean")
	want := 2
	if actual != want{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
	want_2 := false
	if actual_2 != want_2{       
		t.Errorf("Actual: %v -- Want: %v", actual_2, want_2)
	}
}
//----------------------------------------------

//----------------------------------------------
// Tests to make sure gf_helper works with no neighbors
//----------------------------------------------
func TestGF_Helper_Error_Neighbors(t *testing.T) {
	
	//Run code via cmd
	if os.Getenv("FLAG") == "1" {
		var neighbors []*Node
		node_A := &Node{Number: 0, Neighbors: []int{0,1,2}, Lat: 0, Long: 1}
		destination := &Node{Number: 3, Neighbors: []int{0,1,2}, Lat: 3, Long: 1}
		_, _ = gf_Helper(node_A, destination, neighbors, "euclidean")
		return
	}
	
	//Execute the code via exec to check error code
	cmd := exec.Command(os.Args[0], "-test.run=TestGF_Helper_Error")
	cmd.Env = append(os.Environ(), "FLAG=1")
	output, err := cmd.Output() //grab error code and output
	
	//Compare error codes
	actual := err.Error()
	want := "exit status 1"
	if actual != want{       
		t.Errorf("Actual: %s -- Want: %s", actual, want)
	}
	
	//Compare output
	actual = string(output)
	want = "Error finding node with smallest distance! Exiting.\n"
	if actual != want{       
		t.Errorf("Actual: %s -- Want: %s", actual, want)
	}
}
//----------------------------------------------

//----------------------------------------------
// Tests to make sure gf_helper works with max dist
//----------------------------------------------
func TestGF_Helper_Error_Max_Dist(t *testing.T) {
	//Run code via cmd
	if os.Getenv("FLAG") == "2" {
		var neighbors []*Node
		node_A := &Node{Number: 0, Neighbors: []int{0,1,2}, Lat: 0, Long: 1}
		node_B := &Node{Number: 1, Neighbors: []int{0,1,2}, Lat: 1, Long: 1}
		node_C := &Node{Number: 2, Neighbors: []int{0,1,2}, Lat: 2, Long: 1}
		destination := &Node{Number: 3, Neighbors: []int{0,1,2}, Lat: math.MaxFloat64, Long: math.MaxFloat64}
		neighbors = append(neighbors, node_A)
		neighbors = append(neighbors, node_B)
		neighbors = append(neighbors, node_C)
		_, _ = gf_Helper(node_A, destination, neighbors, "euclidean")
		return
	}	
	
	//Execute the code via exec to check error code
	cmd := exec.Command(os.Args[0], "-test.run=TestGF_Helper_Error")
	cmd.Env = append(os.Environ(), "FLAG=2")
	output, err := cmd.Output() //grab error code and output
	
	//Compare error codes
	actual := err.Error()
	want := "exit status 1"
	if actual != want{       
		t.Errorf("Actual: %s -- Want: %s", actual, want)
	}
	
	//Compare output
	actual = string(output)
	want = "Error finding node with smallest distance! Exiting.\n"
	if actual != want{       
		t.Errorf("Actual: %s -- Want: %s", actual, want)
	}	
}
//----------------------------------------------

//----------------------------------------------
// Tests to make sure gf_helper works with packet drop conditions
//----------------------------------------------
func TestGF_Helper_Packet_Drop(t *testing.T) {
	var neighbors []*Node
	node_A := &Node{Number: 0, Neighbors: []int{0,1,2}, Lat: 3, Long: 0}
	node_B := &Node{Number: 1, Neighbors: []int{0,1,2}, Lat: 10, Long: 1}
	node_C := &Node{Number: 2, Neighbors: []int{0,1,2}, Lat: 20, Long: 1}
	destination := &Node{Number: 3, Neighbors: []int{0,1,2}, Lat: 3, Long: 1}
	neighbors = append(neighbors, node_A)
	neighbors = append(neighbors, node_B)
	neighbors = append(neighbors, node_C)
	actual, actual_2 := gf_Helper(node_A, destination, neighbors, "euclidean")
	want := 0
	if actual != want{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
	want_2 := true
	if actual_2 != want_2{       
		t.Errorf("Actual: %v -- Want: %v", actual_2, want_2)
	}
}
//----------------------------------------------

//--------------------------------------------------------------------------------------------
// mgf_Helper(curr_Node *Node, end *Node, neighbors []*Node, previous_Node_Num int, distance_Formula string) (int, bool)
//--------------------------------------------------------------------------------------------

//----------------------------------------------
// Tests to make sure mgf_helper works normally
//----------------------------------------------
func TestMGF_Helper(t *testing.T) {
	var neighbors []*Node
	node_A := &Node{Number: 0, Neighbors: []int{0,1,2}, Lat: 3, Long: 0}
	node_B := &Node{Number: 1, Neighbors: []int{0,1,2}, Lat: 10, Long: 1}
	node_C := &Node{Number: 2, Neighbors: []int{0,1,2}, Lat: 20, Long: 1}
	destination := &Node{Number: 3, Neighbors: []int{0,1,2}, Lat: 3, Long: 1}
	neighbors = append(neighbors, node_A)
	neighbors = append(neighbors, node_B)
	neighbors = append(neighbors, node_C)
	actual, actual_2 := mgf_Helper(node_A, destination, neighbors, 0, "euclidean")
	want := 1
	if actual != want{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
	want_2 := false
	if actual_2 != want_2{       
		t.Errorf("Actual: %v -- Want: %v", actual_2, want_2)
	}
}
//----------------------------------------------

//----------------------------------------------
// Tests to make sure mgf_helper works if youre not a neighbor to your self from the start
//----------------------------------------------
func TestMGF_Helper_Not_Self(t *testing.T) {
	var neighbors []*Node
	node_A := &Node{Number: 0, Neighbors: []int{0,1,2}, Lat: 3, Long: 0}
	node_B := &Node{Number: 1, Neighbors: []int{0,1,2}, Lat: 10, Long: 1}
	node_C := &Node{Number: 2, Neighbors: []int{0,1,2}, Lat: 20, Long: 1}
	destination := &Node{Number: 3, Neighbors: []int{0,1,2}, Lat: 3, Long: 1}
	neighbors = append(neighbors, node_B)
	neighbors = append(neighbors, node_C)
	actual, actual_2 := mgf_Helper(node_A, destination, neighbors, 0, "euclidean")
	want := 1
	if actual != want{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
	want_2 := false
	if actual_2 != want_2{       
		t.Errorf("Actual: %v -- Want: %v", actual_2, want_2)
	}
}
//----------------------------------------------

//----------------------------------------------
// Tests to make sure mgf_helper works if you want to drop the packet
//----------------------------------------------
func TestMGF_Helper_Packet_Drop(t *testing.T) {
	var neighbors []*Node
	node_A := &Node{Number: 0, Neighbors: []int{0,1,2}, Lat: 3, Long: 0}
	node_B := &Node{Number: 1, Neighbors: []int{0,1,2}, Lat: 10, Long: 1}
	node_C := &Node{Number: 2, Neighbors: []int{0,1,2}, Lat: 20, Long: 1}
	destination := &Node{Number: 3, Neighbors: []int{0,1,2}, Lat: 3, Long: 1}
	neighbors = append(neighbors, node_A)
	neighbors = append(neighbors, node_B)
	neighbors = append(neighbors, node_C)
	actual, actual_2 := mgf_Helper(node_A, destination, neighbors, 1, "euclidean")
	want := 1
	if actual != want{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
	want_2 := true
	if actual_2 != want_2{       
		t.Errorf("Actual: %v -- Want: %v", actual_2, want_2)
	}
}
//----------------------------------------------

//----------------------------------------------
// Tests to make sure mgf_helper does not modify the neighbors array
//----------------------------------------------
func TestMGF_Helper_Packet_Modify_Neighbors(t *testing.T) {
	var neighbors []*Node
	node_A := &Node{Number: 0, Neighbors: []int{0,1,2}, Lat: 3, Long: 0}
	node_B := &Node{Number: 1, Neighbors: []int{0,1,2}, Lat: 10, Long: 1}
	node_C := &Node{Number: 2, Neighbors: []int{0,1,2}, Lat: 20, Long: 1}
	destination := &Node{Number: 3, Neighbors: []int{0,1,2}, Lat: 3, Long: 1}
	neighbors = append(neighbors, node_A)
	neighbors = append(neighbors, node_B)
	neighbors = append(neighbors, node_C)
	actual, actual_2 := mgf_Helper(node_A, destination, neighbors, 1, "euclidean")
	want := 1
	if actual != want{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
	want_2 := true
	if actual_2 != want_2{       
		t.Errorf("Actual: %v -- Want: %v", actual_2, want_2)
	}
	want = 3
	if len(neighbors) != want{       
		t.Errorf("Actual: %v -- Want: %v", len(neighbors), want)
	}
}
//----------------------------------------------

//--------------------------------------------------------------------------------------------
// send_Msg(packet Packet, c net.Conn) 
// receive_Msg(c net.Conn) Packet
//--------------------------------------------------------------------------------------------

//----------------------------------------------
// Tests to make sure sending and receiving packets works
//----------------------------------------------
func TestSend_Receive(t *testing.T) {
	channel := make(chan bool, 1) //channel for passing info
	go func(channel chan bool){
		l, err := net.Listen("tcp", "localhost:8080")
		check_Error(err) //check for errors
		channel <- true
		c, err := l.Accept() //accept a connection
		check_Error(err) //check for errors
		packet := receive_Msg(c) //receive and decode the message
		want := -999
		if packet.Previous != want{       
			t.Errorf("Actual: %v -- Want: %v", packet.Previous, want)
		}
		want = 999
		if packet.Destination != want{       
			t.Errorf("Actual: %v -- Want: %v", packet.Destination, want)
		}
		want = 2
		if packet.Interest_Data != want{       
			t.Errorf("Actual: %v -- Want: %v", packet.Interest_Data, want)
		}
		want_2 := "Shutdown"
		if packet.Name != want_2{       
			t.Errorf("Actual: %v -- Want: %v", packet.Name, want_2)
		}
		l.Close()
	} (channel)
	packet := Packet{Previous: -999, Destination: 999, Interest_Data: 2, Name: "Shutdown"}
	<- channel //wait until server is online
	c, err := net.Dial("tcp", "localhost:8080")
	check_Error(err) //check for errors
	send_Msg(packet, c) 
}
//----------------------------------------------

//--------------------------------------------------------------------------------------------
// find_Node(topology []*Node, number int) *Node
//--------------------------------------------------------------------------------------------

//----------------------------------------------
// Tests to make find_Node works
//----------------------------------------------
func TestFind_Node(t *testing.T) {
	var topology []*Node
	node_A := &Node{Number: 0, Neighbors: []int{0,1,2}}
	node_B := &Node{Number: 1, Neighbors: []int{0,1,2}}
	node_C := &Node{Number: 2, Neighbors: []int{0,1,2}}
	topology = append(topology, node_A)
	topology = append(topology, node_B)
	topology = append(topology, node_C)
	actual := find_Node(topology, 1)
	want := node_B
	if actual != want{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
	node_B.Number = 100
	want_2 := 100
	if actual.Number != want_2{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
}
//----------------------------------------------

//----------------------------------------------
// Tests to make find_Node works on empty topology
//----------------------------------------------
func TestFind_Node_Empty_Topology(t *testing.T) {
	var topology []*Node
	actual := find_Node(topology, 1)
	want := -1
	if actual.Number != want{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
}
//----------------------------------------------

//----------------------------------------------
// Tests to make find_Node works when you cant find
//----------------------------------------------
func TestFind_Node_Not_Exist(t *testing.T) {
	var topology []*Node
	node_A := &Node{Number: 0, Neighbors: []int{0,1,2}}
	node_B := &Node{Number: 1, Neighbors: []int{0,1,2}}
	node_C := &Node{Number: 2, Neighbors: []int{0,1,2}}
	topology = append(topology, node_A)
	topology = append(topology, node_B)
	topology = append(topology, node_C)
	actual := find_Node(topology, 4)
	want := -1
	if actual.Number != want{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
}
//----------------------------------------------

//--------------------------------------------------------------------------------------------
// find_Neighbors(topology []*Node, number int) []*Node
//--------------------------------------------------------------------------------------------

//----------------------------------------------
// Tests to make find_Node works
//----------------------------------------------
func TestFind_Neighbor(t *testing.T) {
	var topology []*Node
	var neighbors []*Node
	node_A := &Node{Number: 0, Neighbors: []int{0,1,2}}
	node_B := &Node{Number: 1, Neighbors: []int{0,1,2}}
	node_C := &Node{Number: 2, Neighbors: []int{0,1,2}}
	topology = append(topology, node_A)
	topology = append(topology, node_B)
	topology = append(topology, node_C)
	neighbors = append(neighbors, node_A)
	neighbors = append(neighbors, node_B)
	neighbors = append(neighbors, node_C)
	actual := find_Neighbors(topology, 1)
	for x := 0; x < len(actual); x++{
		if actual[x] != neighbors[x]{       
			t.Errorf("Actual: %v -- Want: %v", actual[x], neighbors[x])
		}
	}
	neighbors[1].Number = 100
	want_2 := 100
	if actual[1].Number != want_2{       
		t.Errorf("Actual: %v -- Want: %v", actual, want_2)
	}
}
//----------------------------------------------

//----------------------------------------------
// Tests to make find_Node works for empty topology
//----------------------------------------------
func TestFind_Neighbor_Empty_Topology(t *testing.T) {
	var topology []*Node
	actual := find_Neighbors(topology, 1)
	want := 0
	if len(actual) != want{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
}
//----------------------------------------------

//----------------------------------------------
// Tests to make find_Node works if the node does not exist
//----------------------------------------------
func TestFind_Neighbor_Not_Exist(t *testing.T) {
	var topology []*Node
	node_A := &Node{Number: 0, Neighbors: []int{0,1,2}}
	node_B := &Node{Number: 1, Neighbors: []int{0,1,2}}
	node_C := &Node{Number: 2, Neighbors: []int{0,1,2}}
	topology = append(topology, node_A)
	topology = append(topology, node_B)
	topology = append(topology, node_C)
	actual := find_Neighbors(topology, 4)
	want := 0
	if len(actual) != want{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
}
//----------------------------------------------

//--------------------------------------------------------------------------------------------
// index_PIT(curr_Node *Node, name string, Outgoing int, Incoming int) *PIT
//--------------------------------------------------------------------------------------------

//----------------------------------------------
// Tests to make index_PIT works with 1 entry
//----------------------------------------------
func TestIndex_PIT_1(t *testing.T) {
	node_A := &Node{Number: 0}
	var PIT_A []PIT
	PIT_A = append(PIT_A, PIT{Incoming: 0, Outgoing: 0, Name: "a"})
	PIT_A = append(PIT_A, PIT{Incoming: 1, Outgoing: 1, Name: "b"})
	PIT_A = append(PIT_A, PIT{Incoming: 2, Outgoing: 2, Name: "c"})
	node_A.Pit = PIT_A
	actual := index_PIT(node_A, "b", -2, -2)
	want := &PIT_A[1]
	if actual != want{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
	PIT_A[1].Incoming = 100
	want_2 := 100
	if actual.Incoming != want_2{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
}
//----------------------------------------------

//----------------------------------------------
// Tests to make index_PIT works with multiple entry
//----------------------------------------------
func TestIndex_PIT_Multiple(t *testing.T) {
	node_A := &Node{Number: 0}
	var PIT_A []PIT
	PIT_A = append(PIT_A, PIT{Incoming: 0, Outgoing: 0, Name: "a"})
	PIT_A = append(PIT_A, PIT{Incoming: 1, Outgoing: 1, Name: "c"})
	PIT_A = append(PIT_A, PIT{Incoming: 2, Outgoing: 2, Name: "c"})
	node_A.Pit = PIT_A
	actual := index_PIT(node_A, "c", -2, -2)
	want := &PIT_A[2]
	if actual != want{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
}
//----------------------------------------------

//----------------------------------------------
// Tests to make index_PIT works with incoming == -2
//----------------------------------------------
func TestIndex_PIT_Outgoing(t *testing.T) {
	node_A := &Node{Number: 0}
	var PIT_A []PIT
	PIT_A = append(PIT_A, PIT{Incoming: 0, Outgoing: 0, Name: "a"})
	PIT_A = append(PIT_A, PIT{Incoming: 1, Outgoing: 1, Name: "c"})
	PIT_A = append(PIT_A, PIT{Incoming: 2, Outgoing: 2, Name: "c"})
	node_A.Pit = PIT_A
	actual := index_PIT(node_A, "c", -2, 1)
	want := &PIT_A[1]
	if actual != want{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
}
//----------------------------------------------

//----------------------------------------------
// Tests to make index_PIT works with outgoing == -2
//----------------------------------------------
func TestIndex_PIT_Incoming(t *testing.T) {
	node_A := &Node{Number: 0}
	var PIT_A []PIT
	PIT_A = append(PIT_A, PIT{Incoming: 0, Outgoing: 0, Name: "a"})
	PIT_A = append(PIT_A, PIT{Incoming: 1, Outgoing: 1, Name: "c"})
	PIT_A = append(PIT_A, PIT{Incoming: 2, Outgoing: 2, Name: "c"})
	node_A.Pit = PIT_A
	actual := index_PIT(node_A, "c", 1, -2)
	want := &PIT_A[1]
	if actual != want{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
}
//----------------------------------------------

//----------------------------------------------
// Tests to make index_PIT works with incoming and outgoing != -2
//----------------------------------------------
func TestIndex_PIT_Incoming_Outgoing(t *testing.T) {
	node_A := &Node{Number: 0}
	var PIT_A []PIT
	PIT_A = append(PIT_A, PIT{Incoming: 0, Outgoing: 0, Name: "a"})
	PIT_A = append(PIT_A, PIT{Incoming: 1, Outgoing: 1, Name: "c"})
	PIT_A = append(PIT_A, PIT{Incoming: 2, Outgoing: 2, Name: "c"})
	node_A.Pit = PIT_A
	actual := index_PIT(node_A, "c", 1, 1)
	want := &PIT_A[1]
	if actual != want{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
}
//----------------------------------------------

//----------------------------------------------
// Tests to make index_PIT works if it fails to find
//----------------------------------------------
func TestIndex_PIT_Fail(t *testing.T) {
	node_A := &Node{Number: 0}
	var PIT_A []PIT
	PIT_A = append(PIT_A, PIT{Incoming: 0, Outgoing: 0, Name: "a"})
	PIT_A = append(PIT_A, PIT{Incoming: 1, Outgoing: 1, Name: "c"})
	PIT_A = append(PIT_A, PIT{Incoming: 2, Outgoing: 2, Name: "c"})
	node_A.Pit = PIT_A
	actual := index_PIT(node_A, "d", -2, -2)
	want := PIT{Incoming: -2, Outgoing: -2}
	if actual.Incoming != want.Incoming && actual.Outgoing != want.Outgoing{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
}
//----------------------------------------------

//----------------------------------------------
// Tests to make index_PIT actually returns the PIT instead of a copy
//----------------------------------------------
func TestIndex_PIT_Real(t *testing.T) {
	node_A := &Node{Number: 0}
	var PIT_A []PIT
	PIT_A = append(PIT_A, PIT{Incoming: 0, Outgoing: 0, Name: "a"})
	PIT_A = append(PIT_A, PIT{Incoming: 1, Outgoing: 1, Name: "b"})
	PIT_A = append(PIT_A, PIT{Incoming: 2, Outgoing: 2, Name: "c"})
	node_A.Pit = PIT_A
	actual := index_PIT(node_A, "b", -2, -2)
	want := &PIT_A[1]
	if actual != want{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
	actual.Name = "Hello"
	if actual.Name != want.Name{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
	want_2 := PIT_A[1]
	if actual.Name != want_2.Name{       
		t.Errorf("Actual: %v -- Want: %v", actual, want_2)
	}
}
//----------------------------------------------

//--------------------------------------------------------------------------------------------
// remove_PIT(curr_Node *Node, local_PIT *PIT) 
//--------------------------------------------------------------------------------------------

//----------------------------------------------
// Tests to make remove_PIT works for 1 entry
//----------------------------------------------
func TestRemove_PIT_1(t *testing.T) {
	node_A := &Node{Number: 0}
	var PIT_A []PIT
	PIT_A = append(PIT_A, PIT{Incoming: 0, Outgoing: 0, Name: "a"})
	PIT_A = append(PIT_A, PIT{Incoming: 1, Outgoing: 1, Name: "b"})
	PIT_A = append(PIT_A, PIT{Incoming: 2, Outgoing: 2, Name: "c"})
	node_A.Pit = PIT_A
	entry := index_PIT(node_A, "b", -2, -2)
	remove_PIT(node_A, entry)
	want := 2
	if len(node_A.Pit) != want{       
		t.Errorf("Actual: %v -- Want: %v", len(node_A.Pit), want)
	}
	want = 3
	if len(PIT_A) != want{       
		t.Errorf("Actual: %v -- Want: %v", len(PIT_A), want)
	}
	if node_A.Pit[0] != PIT_A[0] || node_A.Pit[1] != PIT_A[2]{
		t.Errorf("Actual: %v -- Want: %v", node_A.Pit, PIT_A)
	}
}
//----------------------------------------------

//----------------------------------------------
// Tests to make remove_PIT works with multiple entry
//----------------------------------------------
func TestRemove_PIT_Multiple(t *testing.T) {
	node_A := &Node{Number: 0}
	var PIT_A []PIT
	PIT_A = append(PIT_A, PIT{Incoming: 0, Outgoing: 0, Name: "a"})
	PIT_A = append(PIT_A, PIT{Incoming: 1, Outgoing: 1, Name: "b"})
	PIT_A = append(PIT_A, PIT{Incoming: 2, Outgoing: 2, Name: "c"})
	PIT_A = append(PIT_A, PIT{Incoming: 7, Outgoing: 7, Name: "b"})
	node_A.Pit = PIT_A
	entry := index_PIT(node_A, "b", -2, -2)
	remove_PIT(node_A, entry)
	want := 3
	if len(node_A.Pit) != want{       
		t.Errorf("Actual: %v -- Want: %v", len(node_A.Pit), want)
	}
	want = 4
	if len(PIT_A) != want{       
		t.Errorf("Actual: %v -- Want: %v", len(PIT_A), want)
	}
	if node_A.Pit[0] != PIT_A[0] || node_A.Pit[1] != PIT_A[1] || node_A.Pit[2] != PIT_A[2]{
		t.Errorf("Actual: %v -- Want: %v", node_A.Pit, PIT_A)
	}
}
//----------------------------------------------

//----------------------------------------------
// Tests to make remove_PIT works with if it cant find the entry
//----------------------------------------------
func TestRemove_PIT_Fail(t *testing.T) {

	//Run code via cmd
	if os.Getenv("FLAG") == "1" {
		node_A := &Node{Number: 0}
		var PIT_A []PIT
		PIT_A = append(PIT_A, PIT{Incoming: 0, Outgoing: 0, Name: "a"})
		PIT_A = append(PIT_A, PIT{Incoming: 1, Outgoing: 1, Name: "b"})
		PIT_A = append(PIT_A, PIT{Incoming: 2, Outgoing: 2, Name: "c"})
		node_A.Pit = PIT_A
		entry := index_PIT(node_A, "d", -2, -2)
		remove_PIT(node_A, entry)
	}
	
	//Execute the code via exec to check error code
	cmd := exec.Command(os.Args[0], "-test.run=TestRemove_PIT_Fail")
	cmd.Env = append(os.Environ(), "FLAG=1")
	output, err := cmd.Output() //grab error code and output
	
	//Compare error codes
	actual := err.Error()
	want := "exit status 1"
	if actual != want{       
		t.Errorf("Actual: %s -- Want: %s", actual, want)
	}
	
	//Compare output
	actual = string(output)
	want = "Error finding the PIT entry to remove! Exiting.\n"
	if actual != want{       
		t.Errorf("Actual: %s -- Want: %s", actual, want)
	}
}
//----------------------------------------------

//--------------------------------------------------------------------------------------------
// index_CACHE(curr_Node *Node, name string) *CACHE
//--------------------------------------------------------------------------------------------

//----------------------------------------------
// Tests to make index_CACHE works with 1 entry
//----------------------------------------------
func TestIndex_CACHE_1(t *testing.T) {
	node_A := &Node{Number: 0}
	var CACHE_A []CACHE
	CACHE_A = append(CACHE_A, CACHE{Name: "a", Data: "a"})
	CACHE_A = append(CACHE_A, CACHE{Name: "b", Data: "b"})
	CACHE_A = append(CACHE_A, CACHE{Name: "c", Data: "c"})
	node_A.Cache = CACHE_A
	actual := index_CACHE(node_A, "b")
	want := &CACHE_A[1]
	if actual != want{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
	actual.Name = "Hello"
	if actual.Name != want.Name{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
	want_2 := CACHE_A[1]
	if actual.Name != want_2.Name{       
		t.Errorf("Actual: %v -- Want: %v", actual, want_2)
	}
}
//----------------------------------------------

//----------------------------------------------
// Tests to make index_CACHE works with an entry that doesnt exist
//----------------------------------------------
func TestIndex_CACHE_Fail(t *testing.T) {
	node_A := &Node{Number: 0}
	var CACHE_A []CACHE
	CACHE_A = append(CACHE_A, CACHE{Name: "a", Data: "a"})
	CACHE_A = append(CACHE_A, CACHE{Name: "b", Data: "b"})
	CACHE_A = append(CACHE_A, CACHE{Name: "c", Data: "c"})
	node_A.Cache = CACHE_A
	actual := index_CACHE(node_A, "d")
	want := &CACHE{Timestamp: time.Time{}}
	if !((actual.Timestamp).Equal(want.Timestamp)){       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
}
//----------------------------------------------

//--------------------------------------------------------------------------------------------
// remove_CACHE(curr_Node *Node, local_CACHE *CACHE)
//--------------------------------------------------------------------------------------------

//----------------------------------------------
// Tests to make remove_CACHE works for 1 entry
//----------------------------------------------
func TestRemove_CACHE_1(t *testing.T) {
	node_A := &Node{Number: 0}
	var CACHE_A []CACHE
	CACHE_A = append(CACHE_A, CACHE{Name: "a", Data: "a"})
	CACHE_A = append(CACHE_A, CACHE{Name: "b", Data: "b"})
	CACHE_A = append(CACHE_A, CACHE{Name: "c", Data: "c"})
	node_A.Cache = CACHE_A
	entry := index_CACHE(node_A, "b")
	remove_CACHE(node_A, entry)
	want := 2
	if len(node_A.Cache) != want{       
		t.Errorf("Actual: %v -- Want: %v", len(node_A.Cache), want)
	}
	want = 3
	if len(CACHE_A) != want{       
		t.Errorf("Actual: %v -- Want: %v", len(CACHE_A), want)
	}
	if node_A.Cache[0] != CACHE_A[0] || node_A.Cache[1] != CACHE_A[2]{
		t.Errorf("Actual: %v -- Want: %v", node_A.Cache, CACHE_A)
	}
}
//----------------------------------------------

//----------------------------------------------
// Tests to make remove_CACHE works with if it cant find the entry
//----------------------------------------------
func TestRemove_CACHE_Fail(t *testing.T) {

	//Run code via cmd
	if os.Getenv("FLAG") == "1" {
		node_A := &Node{Number: 0}
		var CACHE_A []CACHE
		CACHE_A = append(CACHE_A, CACHE{Name: "a", Data: "a"})
		CACHE_A = append(CACHE_A, CACHE{Name: "b", Data: "b"})
		CACHE_A = append(CACHE_A, CACHE{Name: "c", Data: "c"})
		node_A.Cache = CACHE_A
		entry := index_CACHE(node_A, "d")
		remove_CACHE(node_A, entry)
	}
	
	//Execute the code via exec to check error code
	cmd := exec.Command(os.Args[0], "-test.run=TestRemove_CACHE_Fail")
	cmd.Env = append(os.Environ(), "FLAG=1")
	output, err := cmd.Output() //grab error code and output
	
	//Compare error codes
	actual := err.Error()
	want := "exit status 1"
	if actual != want{       
		t.Errorf("Actual: %s -- Want: %s", actual, want)
	}
	
	//Compare output
	actual = string(output)
	want = "Error finding the CACHE entry to remove! Exiting.\n"
	if actual != want{       
		t.Errorf("Actual: %s -- Want: %s", actual, want)
	}
}
//----------------------------------------------

//--------------------------------------------------------------------------------------------
// determine_CACHE(caching_Strat string) bool
//--------------------------------------------------------------------------------------------

//----------------------------------------------
// Tests to make sure determine_CACHE works with LCE
//----------------------------------------------
func TestDetermine_CACHE_LCE(t *testing.T) {
	actual := determine_CACHE("LCE")
	want := true
	if actual != want{
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
}
//----------------------------------------------

//----------------------------------------------
// Tests to make sure determine_CACHE works with RAND
//----------------------------------------------
func TestDetermine_CACHE_RAND(t *testing.T) {
	rand.Seed(1)
	actual := determine_CACHE("RAND")
	want := false
	if actual != want{
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
}
//----------------------------------------------

//----------------------------------------------
// Tests to make sure determine_CACHE works with NONE
//----------------------------------------------
func TestDetermine_CACHE_NONE(t *testing.T) {
	actual := determine_CACHE("NONE")
	want := false
	if actual != want{
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
}
//----------------------------------------------

//----------------------------------------------
// Tests to make sure determine_CACHE works with unknown strategy
//----------------------------------------------
func TestDetermine_CACHE_FAIL(t *testing.T) {

	//Run code via cmd
	if os.Getenv("FLAG") == "1" {
		_ = determine_CACHE("UNKNOWN")
	}
	
	//Execute the code via exec to check error code
	cmd := exec.Command(os.Args[0], "-test.run=TestDetermine_CACHE_FAIL")
	cmd.Env = append(os.Environ(), "FLAG=1")
	output, err := cmd.Output() //grab error code and output
	
	//Compare error codes
	actual := err.Error()
	want := "exit status 1"
	if actual != want{       
		t.Errorf("Actual: %s -- Want: %s", actual, want)
	}
	
	//Compare output
	actual = string(output)
	want = "Unknown caching strategy! Exiting.\n"
	if actual != want{       
		t.Errorf("Actual: %s -- Want: %s", actual, want)
	}
}
//----------------------------------------------

//--------------------------------------------------------------------------------------------
// evict_CACHE(curr_Node *Node, name string, cache_Evict string)
//--------------------------------------------------------------------------------------------

//----------------------------------------------
// Tests to make sure evict_CACHE works with when the cache is empty
//----------------------------------------------
func TestEvict_CACHE_Empty(t *testing.T) {

	//Run code via cmd
	if os.Getenv("FLAG") == "1" {
		node_A := &Node{Number: 0}
		evict_CACHE(node_A, "a", "FIFO")
	}
	
	//Execute the code via exec to check error code
	cmd := exec.Command(os.Args[0], "-test.run=TestEvict_CACHE_Empty")
	cmd.Env = append(os.Environ(), "FLAG=1")
	output, err := cmd.Output() //grab error code and output
	
	//Compare error codes
	actual := err.Error()
	want := "exit status 1"
	if actual != want{       
		t.Errorf("Actual: %s -- Want: %s", actual, want)
	}
	
	//Compare output
	actual = string(output)
	want = "Called evict_CACHE on empty cache! Exiting.\n"
	if actual != want{       
		t.Errorf("Actual: %s -- Want: %s", actual, want)
	}
}
//----------------------------------------------

//----------------------------------------------
// Tests to make sure evict_CACHE works with matching entry
//----------------------------------------------
func TestEvict_CACHE_Same(t *testing.T) {
	node_A := &Node{Number: 0}
	var CACHE_A []CACHE
	CACHE_A = append(CACHE_A, CACHE{Name: "a", Data: "a"})
	CACHE_A = append(CACHE_A, CACHE{Name: "b", Data: "b"})
	CACHE_A = append(CACHE_A, CACHE{Name: "c", Data: "c"})	
	node_A.Cache = CACHE_A
	evict_CACHE(node_A, "b", "FIFO")
	
	want := 2
	if len(node_A.Cache) != want{       
		t.Errorf("Actual: %v -- Want: %v", len(node_A.Cache), want)
	}
	want = 3
	if len(CACHE_A) != want{       
		t.Errorf("Actual: %v -- Want: %v", len(CACHE_A), want)
	}
	if node_A.Cache[0] != CACHE_A[0] || node_A.Cache[1] != CACHE_A[2]{
		t.Errorf("Actual: %v -- Want: %v", node_A.Cache, CACHE_A)
	}
}
//----------------------------------------------

//----------------------------------------------
// Tests to make sure evict_CACHE works with TTL
//----------------------------------------------
func TestEvict_CACHE_TTL(t *testing.T) {
	node_A := &Node{Number: 0}
	var CACHE_A []CACHE
	CACHE_A = append(CACHE_A, CACHE{Name: "a", Data: "a", TTL: time.Duration(time.Second * 10), Timestamp: time.Now()}) //10 seconds ttl
	CACHE_A = append(CACHE_A, CACHE{Name: "b", Data: "b", TTL: time.Duration(0), Timestamp: time.Now()}) //0 seconds ttl
	CACHE_A = append(CACHE_A, CACHE{Name: "c", Data: "c", TTL: time.Duration(time.Second * 10), Timestamp: time.Now()}) //10 seconds ttl
	node_A.Cache = CACHE_A
	evict_CACHE(node_A, "d", "FIFO")
	
	want := 2
	if len(node_A.Cache) != want{       
		t.Errorf("Actual: %v -- Want: %v", len(node_A.Cache), want)
	}
	want = 3
	if len(CACHE_A) != want{       
		t.Errorf("Actual: %v -- Want: %v", len(CACHE_A), want)
	}
	if node_A.Cache[0] != CACHE_A[0] || node_A.Cache[1] != CACHE_A[2]{
		t.Errorf("Actual: %v -- Want: %v", node_A.Cache, CACHE_A)
	}
}
//----------------------------------------------

//----------------------------------------------
// Tests to make sure evict_CACHE works with FIFO
//----------------------------------------------
func TestEvict_CACHE_FIFO(t *testing.T) {
	node_A := &Node{Number: 0}
	var CACHE_A []CACHE
	CACHE_A = append(CACHE_A, CACHE{Name: "a", Data: "a", TTL: time.Duration(-1), Timestamp: (time.Now()).Add(-time.Duration(10 * time.Second))})
	CACHE_A = append(CACHE_A, CACHE{Name: "b", Data: "b", TTL: time.Duration(-1), Timestamp: (time.Now()).Add(-time.Duration(20 * time.Second))})
	CACHE_A = append(CACHE_A, CACHE{Name: "c", Data: "c", TTL: time.Duration(-1), Timestamp: (time.Now()).Add(-time.Duration(30 * time.Second))})
	node_A.Cache = CACHE_A
	evict_CACHE(node_A, "d", "FIFO")
	
	want := 2
	if len(node_A.Cache) != want{       
		t.Errorf("Actual: %v -- Want: %v", len(node_A.Cache), want)
	}
	want = 3
	if len(CACHE_A) != want{       
		t.Errorf("Actual: %v -- Want: %v", len(CACHE_A), want)
	}
	if node_A.Cache[0] != CACHE_A[0] || node_A.Cache[1] != CACHE_A[1]{
		t.Errorf("Actual: %v -- Want: %v", node_A.Cache, CACHE_A)
	}
}
//----------------------------------------------

//----------------------------------------------
// Tests to make sure evict_CACHE works with LRU
//----------------------------------------------
func TestEvict_CACHE_LRU(t *testing.T) {
	node_A := &Node{Number: 0}
	var CACHE_A []CACHE
	CACHE_A = append(CACHE_A, CACHE{Name: "a", Data: "a", TTL: time.Duration(-1), Last_Used: (time.Now()).Add(-time.Duration(30 * time.Second))})
	CACHE_A = append(CACHE_A, CACHE{Name: "b", Data: "b", TTL: time.Duration(-1), Last_Used: (time.Now()).Add(-time.Duration(20 * time.Second))})
	CACHE_A = append(CACHE_A, CACHE{Name: "c", Data: "c", TTL: time.Duration(-1), Last_Used: (time.Now()).Add(-time.Duration(10 * time.Second))})
	node_A.Cache = CACHE_A
	evict_CACHE(node_A, "d", "LRU")
	
	want := 2
	if len(node_A.Cache) != want{       
		t.Errorf("Actual: %v -- Want: %v", len(node_A.Cache), want)
	}
	want = 3
	if len(CACHE_A) != want{       
		t.Errorf("Actual: %v -- Want: %v", len(CACHE_A), want)
	}
	if node_A.Cache[0] != CACHE_A[1] || node_A.Cache[1] != CACHE_A[2]{
		t.Errorf("Actual: %v -- Want: %v", node_A.Cache, CACHE_A)
	}
}
//----------------------------------------------

//----------------------------------------------
// Tests to make sure evict_CACHE works with LRU
//----------------------------------------------
func TestEvict_CACHE_LFU(t *testing.T) {
	node_A := &Node{Number: 0}
	var CACHE_A []CACHE
	CACHE_A = append(CACHE_A, CACHE{Name: "a", Data: "a", TTL: time.Duration(-1), Frequency: 2})
	CACHE_A = append(CACHE_A, CACHE{Name: "b", Data: "b", TTL: time.Duration(-1), Frequency: 1})
	CACHE_A = append(CACHE_A, CACHE{Name: "c", Data: "c", TTL: time.Duration(-1), Frequency: 3})
	node_A.Cache = CACHE_A
	evict_CACHE(node_A, "d", "LFU")
	
	want := 2
	if len(node_A.Cache) != want{       
		t.Errorf("Actual: %v -- Want: %v", len(node_A.Cache), want)
	}
	want = 3
	if len(CACHE_A) != want{       
		t.Errorf("Actual: %v -- Want: %v", len(CACHE_A), want)
	}
	if node_A.Cache[0] != CACHE_A[0] || node_A.Cache[1] != CACHE_A[2]{
		t.Errorf("Actual: %v -- Want: %v", node_A.Cache, CACHE_A)
	}
}
//----------------------------------------------

//----------------------------------------------
// Tests to make sure evict_CACHE works with when the cache is empty
//----------------------------------------------
func TestEvict_CACHE_Unknown(t *testing.T) {

	//Run code via cmd
	if os.Getenv("FLAG") == "1" {
		node_A := &Node{Number: 0}
		var CACHE_A []CACHE
		CACHE_A = append(CACHE_A, CACHE{Name: "a", Data: "a", TTL: time.Duration(-1), Frequency: 2})
		CACHE_A = append(CACHE_A, CACHE{Name: "b", Data: "b", TTL: time.Duration(-1), Frequency: 1})
		CACHE_A = append(CACHE_A, CACHE{Name: "c", Data: "c", TTL: time.Duration(-1), Frequency: 3})
		node_A.Cache = CACHE_A
		evict_CACHE(node_A, "d", "Unknown")
	}
	
	//Execute the code via exec to check error code
	cmd := exec.Command(os.Args[0], "-test.run=TestEvict_CACHE_Unknown")
	cmd.Env = append(os.Environ(), "FLAG=1")
	output, err := cmd.Output() //grab error code and output
	
	//Compare error codes
	actual := err.Error()
	want := "exit status 1"
	if actual != want{       
		t.Errorf("Actual: %s -- Want: %s", actual, want)
	}
	
	//Compare output
	actual = string(output)
	want = "Unknown cache eviction policy! Exiting.\n"
	if actual != want{       
		t.Errorf("Actual: %s -- Want: %s", actual, want)
	}
}
//----------------------------------------------

//--------------------------------------------------------------------------------------------
// reverse_Check(path_Traversed []int)
//--------------------------------------------------------------------------------------------

//----------------------------------------------
// Tests to make sure reverse_Check works 
//----------------------------------------------
func TestReverse_Check_Same(t *testing.T) {
	path := []int{}
	reverse_Check(path)
	path = []int{1}
	reverse_Check(path)
	path = []int{1,2,3,4,3,2,1}
	reverse_Check(path)
}
//----------------------------------------------

//----------------------------------------------
// Tests to make sure reverse_Check works when on fail
//----------------------------------------------
func TestReverse_Check_Fail(t *testing.T) {

	//Run code via cmd
	if os.Getenv("FLAG") == "1" {
		path := []int{1,2,3,4,2,3,1}
		reverse_Check(path)
	}
	
	if os.Getenv("FLAG") == "2" {
		path := []int{1,2,3,4,4,3,2,1}
		reverse_Check(path)
	}
	
	//Execute the code via exec to check error code
	cmd := exec.Command(os.Args[0], "-test.run=TestReverse_Check_Fail")
	cmd.Env = append(os.Environ(), "FLAG=1")
	output, err := cmd.Output() //grab error code and output
	
	//Compare error codes
	actual := err.Error()
	want := "exit status 1"
	if actual != want{       
		t.Errorf("Actual: %s -- Want: %s", actual, want)
	}
	
	//Compare output
	actual = string(output)
	want = "Error! Not following the reverse of the taken path!\n"
	if actual != want{       
		t.Errorf("Actual: %s -- Want: %s", actual, want)
	}
	
	//Execute the code via exec to check error code
	cmd = exec.Command(os.Args[0], "-test.run=TestReverse_Check_Fail")
	cmd.Env = append(os.Environ(), "FLAG=2")
	output, err = cmd.Output() //grab error code and output
	
	//Compare error codes
	actual = err.Error()
	want = "exit status 1"
	if actual != want{       
		t.Errorf("Actual: %s -- Want: %s", actual, want)
	}
	
	//Compare output
	actual = string(output)
	want = "Error! Not following the reverse of the taken path!\n"
	if actual != want{       
		t.Errorf("Actual: %s -- Want: %s", actual, want)
	}
}
//----------------------------------------------

//--------------------------------------------------------------------------------------------
// empty_PIT_Check(topology []*Node)
//--------------------------------------------------------------------------------------------

//----------------------------------------------
// Tests to make sure reverse_Check works 
//----------------------------------------------
func TestEmpty_PIT_Check_Empty(t *testing.T) {
	var topology []*Node
	empty_PIT_Check(topology)
	node_A := &Node{Number: 0}
	node_B := &Node{Number: 0}
	topology = append(topology, node_A)
	topology = append(topology, node_B)
	empty_PIT_Check(topology)
}
//----------------------------------------------

//----------------------------------------------
// Tests to make sure reverse_Check works when on fail
//----------------------------------------------
func TestEmpty_PIT_Check_Fail(t *testing.T) {

	//Run code via cmd
	if os.Getenv("FLAG") == "1" {
		var topology []*Node
		node_A := &Node{Number: 0}
		topology = append(topology, node_A)
		node_B := &Node{Number: 1}
		var PIT_B []PIT
		PIT_B = append(PIT_B, PIT{Incoming: 0, Outgoing: 0, Name: "a"})
		node_B.Pit = PIT_B
		topology = append(topology, node_B)
		empty_PIT_Check(topology)
	}
	
	//Execute the code via exec to check error code
	cmd := exec.Command(os.Args[0], "-test.run=Empty_PIT_Check_Fail")
	cmd.Env = append(os.Environ(), "FLAG=1")
	output, err := cmd.Output() //grab error code and output
	
	//Compare error codes
	actual := err.Error()
	want := "exit status 1"
	if actual != want{       
		t.Errorf("Actual: %s -- Want: %s", actual, want)
	}
	
	//Compare output
	actual = string(output)
	want = "Error! PIT is not empty for node 1!\n"
	if actual != want{       
		t.Errorf("Actual: %s -- Want: %s", actual, want)
	}
}
//----------------------------------------------

//--------------------------------------------------------------------------------------------
// calc_Size(packet Packet) float64
//--------------------------------------------------------------------------------------------

//----------------------------------------------
// Tests to make sure packet size works with empty packet
//----------------------------------------------
func TestCalc_Size_Empty(t *testing.T) {
	var packet Packet
	actual := calc_Size(packet)
	want := float64(4 + 4 + 4 + 4 + 0 + 0 + 8 + 8 + 4 + 0 + 0 + 8)
	//Compare output
	if actual != want{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
}
//----------------------------------------------

//----------------------------------------------
// Tests to make sure packet size works with all fields individually
//----------------------------------------------
func TestCalc_Size_Individual(t *testing.T) {
	
	//Previous
	packet := Packet{}
	packet.Previous = 100
	actual := calc_Size(packet)
	want := float64(4 + 4 + 4 + 4 + 0 + 0 + 8 + 8 + 4 + 0 + 0 + 8)
	//Compare output
	if actual != want{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
	
	//Start
	packet = Packet{}
	packet.Start = 100
	actual = calc_Size(packet)
	want = float64(4 + 4 + 4 + 4 + 0 + 0 + 8 + 8 + 4 + 0 + 0 + 8)
	//Compare output
	if actual != want{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
	
	//Destination
	packet = Packet{}
	packet.Destination = 100
	actual = calc_Size(packet)
	want = float64(4 + 4 + 4 + 4 + 0 + 0 + 8 + 8 + 4 + 0 + 0 + 8)
	//Compare output
	if actual != want{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}

	//Interest_Data
	packet = Packet{}
	packet.Interest_Data = 100
	actual = calc_Size(packet)
	want = float64(4 + 4 + 4 + 4 + 0 + 0 + 8 + 8 + 4 + 0 + 0 + 8)
	//Compare output
	if actual != want{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
	
	//Name
	packet = Packet{}
	packet.Name = "HelloWorld"
	actual = calc_Size(packet)
	want = float64(4 + 4 + 4 + 4 + 10 + 0 + 8 + 8 + 4 + 0 + 0 + 8)
	//Compare output
	if actual != want{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
	
	//Payload
	packet = Packet{}
	packet.Payload = "HelloWorld"
	actual = calc_Size(packet)
	want = float64(4 + 4 + 4 + 4 + 0 + 0 + 8 + 8 + 4 + 0 + 0 + 8)
	//Compare output
	if actual != want{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
	
	//Start_Time
	packet = Packet{}
	packet.Start_Time = time.Time{}
	actual = calc_Size(packet)
	want = float64(4 + 4 + 4 + 4 + 0 + 0 + 8 + 8 + 4 + 0 + 0 + 8)
	//Compare output
	if actual != want{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}	
	
	//TTL
	packet = Packet{}
	packet.TTL = time.Duration(time.Second * 100)
	actual = calc_Size(packet)
	want = float64(4 + 4 + 4 + 4 + 0 + 0 + 8 + 8 + 4 + 0 + 0 + 8)
	//Compare output
	if actual != want{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
	
	//GP_Mode
	packet = Packet{}
	packet.GP_Mode = 100
	actual = calc_Size(packet)
	want = float64(4 + 4 + 4 + 4 + 0 + 0 + 8 + 8 + 4 + 0 + 0 + 8)
	//Compare output
	if actual != want{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
	
	//GP_Mode -1
	packet = Packet{}
	packet.GP_Mode = -1
	actual = calc_Size(packet)
	want = float64(4 + 4 + 4 + 4 + 0 + 0 + 8 + 8 + 0 + 0 + 0 + 8)
	//Compare output
	if actual != want{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
	
	//Nodes_Visit empty
	packet = Packet{}
	packet.Nodes_Visit = []int{}
	actual = calc_Size(packet)
	want = float64(4 + 4 + 4 + 4 + 0 + 0 + 8 + 8 + 4 + 0 + 0 + 8)
	//Compare output
	if actual != want{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
	
	//Nodes_Visit 2 values
	packet = Packet{}
	packet.Nodes_Visit = []int{1,2}
	actual = calc_Size(packet)
	want = float64(4 + 4 + 4 + 4 + 0 + 0 + 8 + 8 + 4 + 8 + 0 + 8)
	//Compare output
	if actual != want{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
	
	//Num_Visit empty
	packet = Packet{}
	packet.Num_Visit = []int{}
	actual = calc_Size(packet)
	want = float64(4 + 4 + 4 + 4 + 0 + 0 + 8 + 8 + 4 + 0 + 0 + 8)
	//Compare output
	if actual != want{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
	
	//Num_Visit 2 values
	packet = Packet{}
	packet.Num_Visit = []int{1,2}
	actual = calc_Size(packet)
	want = float64(4 + 4 + 4 + 4 + 0 + 0 + 8 + 8 + 4 + 0 + 8 + 8)
	//Compare output
	if actual != want{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
	
	//GP_Distance
	packet = Packet{}
	packet.GP_Distance = 100
	actual = calc_Size(packet)
	want = float64(4 + 4 + 4 + 4 + 0 + 0 + 8 + 8 + 4 + 0 + 0 + 8)
	//Compare output
	if actual != want{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
	
	//GP_Distance -1
	packet = Packet{}
	packet.GP_Distance = -1
	actual = calc_Size(packet)
	want = float64(4 + 4 + 4 + 4 + 0 + 0 + 8 + 8 + 4 + 0 + 0 + 0)
	//Compare output
	if actual != want{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
	
	//Path_Traversed empty
	packet = Packet{}
	packet.Path_Traversed = []int{}
	actual = calc_Size(packet)
	want = float64(4 + 4 + 4 + 4 + 0 + 0 + 8 + 8 + 4 + 0 + 0 + 8)
	//Compare output
	if actual != want{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
	
	//Path_Traversed 2 values
	packet = Packet{}
	packet.Path_Traversed = []int{1,2}
	actual = calc_Size(packet)
	want = float64(4 + 4 + 4 + 4 + 0 + 0 + 8 + 8 + 4 + 0 + 0 + 8)
	//Compare output
	if actual != want{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
	
	//Pressure_Gauge
	packet = Packet{}
	packet.Pressure_Gauge = 100
	actual = calc_Size(packet)
	want = float64(4 + 4 + 4 + 4 + 0 + 0 + 8 + 8 + 4 + 0 + 0 + 8)
	//Compare output
	if actual != want{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
	
	//Index_Middle
	packet = Packet{}
	packet.Index_Middle = 100
	actual = calc_Size(packet)
	want = float64(4 + 4 + 4 + 4 + 0 + 0 + 8 + 8 + 4 + 0 + 0 + 8)
	//Compare output
	if actual != want{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
	
	//End_Time
	packet = Packet{}
	packet.End_Time = time.Time{}
	actual = calc_Size(packet)
	want = float64(4 + 4 + 4 + 4 + 0 + 0 + 8 + 8 + 4 + 0 + 0 + 8)
	//Compare output
	if actual != want{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
	
	//Collapsed
	packet = Packet{}
	packet.Collapsed = 100
	actual = calc_Size(packet)
	want = float64(4 + 4 + 4 + 4 + 0 + 0 + 8 + 8 + 4 + 0 + 0 + 8)
	//Compare output
	if actual != want{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
	
}
//----------------------------------------------

//--------------------------------------------------------------------------------------------
// read_Topology(file_Path string, ip string, port int, cache_Size float64) []*Node
//--------------------------------------------------------------------------------------------

//----------------------------------------------
// Tests to make sure reading in test topology works correctly
//----------------------------------------------
func TestRead_Topology(t *testing.T) {
	actual_topology := read_Topology("./unit_test_topologies/test_topology.csv", "localhost", 8000, 1)
	
	//Compare output
	neighbors := [][]int{{0}, {1,5}, {2,5}, {3,4}, {3,4}, {1,2,5}}
	for x := 0; x < len(actual_topology); x++{
		lat, _ := strconv.ParseFloat(strconv.Itoa(x+1), 64)
		long, _ := strconv.ParseFloat(strconv.Itoa(x+1), 64)
		theta := math.Atan(lat/long)
		r :=  math.Sqrt((lat*lat) + (long*long))
		cache_size := int(float64(len(actual_topology)) * float64(1))
		data_name := "/ucla/videos/demo.mpg/1/" + strconv.Itoa(x)
		data_value := ("Lat: " + strconv.FormatFloat(lat, 'E', -1, 64) + "Long: " + strconv.FormatFloat(long, 'E', -1, 64))
		
		if actual_topology[x].IP != "localhost"{       
			t.Errorf("Actual: %v -- Want: %v", actual_topology[x].IP, "localhost")
		}	
		if actual_topology[x].Port != 8000+x{       
			t.Errorf("Actual: %v -- Want: %v", actual_topology[x].Port, 8000+x)
		}
		if actual_topology[x].Number != x{       
			t.Errorf("Actual: %v -- Want: %v", actual_topology[x].Number, x)
		}	
		if actual_topology[x].Lat != lat{       
			t.Errorf("Actual: %v -- Want: %v", actual_topology[x].Lat, lat)
		}
		if actual_topology[x].Long != long{       
			t.Errorf("Actual: %v -- Want: %v", actual_topology[x].Long, long)
		}
		if actual_topology[x].Theta!= theta{       
			t.Errorf("Actual: %v -- Want: %v", actual_topology[x].Theta, theta)
		}	
		if actual_topology[x].R != r{       
			t.Errorf("Actual: %v -- Want: %v", actual_topology[x].R, r)
		}
		if actual_topology[x].Data_Name != data_name{       
			t.Errorf("Actual: %v -- Want: %v", actual_topology[x].Data_Name, data_name)
		}	
		if actual_topology[x].Data_Value != data_value{       
			t.Errorf("Actual: %v -- Want: %v", actual_topology[x].Data_Value, data_value)
		}			
		if !reflect.DeepEqual(actual_topology[x].Neighbors, neighbors[x]){       
			t.Errorf("Actual: %v -- Want: %v", actual_topology[x].Neighbors, neighbors[x])
		}
		if actual_topology[x].Cache_Size != cache_size{       
			t.Errorf("Actual: %v -- Want: %v", actual_topology[x].Cache_Size, cache_size)
		}
	}
}
//----------------------------------------------

//----------------------------------------------
// Tests to make sure reading in test topology fails if file is not found
//----------------------------------------------
func TestRead_Topology_Not_Found(t *testing.T) {

	//Run code via cmd
	if os.Getenv("FLAG") == "1" {
		_ = read_Topology("./unit_test_topologies/test_topology_not_real.csv", "localhost", 8000, 1)
	}
	
	//Execute the code via exec to check error code
	cmd := exec.Command(os.Args[0], "-test.run=TestRead_Topology_Not_Found")
	cmd.Env = append(os.Environ(), "FLAG=1")
	output, err := cmd.Output() //grab error code and output
	
	//Compare error codes
	actual := err.Error()
	want := "exit status 1"
	if actual != want{       
		t.Errorf("Actual: %s -- Want: %s", actual, want)
	}
	
	//Compare output
	actual = string(output)
	want = "Error! File: ./unit_test_topologies/test_topology_not_real.csv was not found!\n"
	if actual != want{       
		t.Errorf("Actual: %s -- Want: %s", actual, want)
	}
}
//----------------------------------------------

//----------------------------------------------
// Tests to make sure reading in test topology fails if there is a neighbor mismatch
//----------------------------------------------
func TestRead_Topology_Bad_Neighbor(t *testing.T) {

	//Run code via cmd
	if os.Getenv("FLAG") == "1" {
		_ = read_Topology("./unit_test_topologies/test_topology_bad_neighbors.csv", "localhost", 8000, 1)
	}
	
	//Execute the code via exec to check error code
	cmd := exec.Command(os.Args[0], "-test.run=TestRead_Topology_Bad_Neighbor")
	cmd.Env = append(os.Environ(), "FLAG=1")
	output, err := cmd.Output() //grab error code and output
	
	//Compare error codes
	actual := err.Error()
	want := "exit status 1"
	if actual != want{       
		t.Errorf("Actual: %s -- Want: %s", actual, want)
	}
	
	//Compare output
	actual = string(output)
	want = "Error! Node 2 or Node 5 has not symmetrical neighbors!\n"
	if actual != want{       
		t.Errorf("Actual: %s -- Want: %s", actual, want)
	}
}
//----------------------------------------------

//----------------------------------------------
// Tests to make sure reading in test topology fails if you are not a neighbor to yourself
//----------------------------------------------
func TestRead_Topology_Self(t *testing.T) {

	//Run code via cmd
	if os.Getenv("FLAG") == "1" {
		_ = read_Topology("./unit_test_topologies/test_topology_self.csv", "localhost", 8000, 1)
	}
	
	//Execute the code via exec to check error code
	cmd := exec.Command(os.Args[0], "-test.run=TestRead_Topology_Self")
	cmd.Env = append(os.Environ(), "FLAG=1")
	output, err := cmd.Output() //grab error code and output
	
	//Compare error codes
	actual := err.Error()
	want := "exit status 1"
	if actual != want{       
		t.Errorf("Actual: %s -- Want: %s", actual, want)
	}
	
	//Compare output
	actual = string(output)
	want = "Error! Node 2 must be a neighbor to itself!\n"
	if actual != want{       
		t.Errorf("Actual: %s -- Want: %s", actual, want)
	}
}
//----------------------------------------------

//--------------------------------------------------------------------------------------------
// metrics_Calc(packet Packet, topology []*Node) []float64
//--------------------------------------------------------------------------------------------

//----------------------------------------------
// Tests to make sure metrics are calculated correctly when not interest collapsing
//----------------------------------------------
func TestMetrics_Calc(t *testing.T) {
	time_now := time.Now()
	topology := read_Topology("./unit_test_topologies/test_topology.csv", "localhost", 8000, 1)
	packet := Packet{Start_Time: (time_now).Add(-time.Duration(30 * time.Second)), End_Time: time_now, Start: 1, Destination: 2, Path_Traversed: []int{1,5,2,5,1}, Collapsed: -1, Pressure_Gauge: 1, Index_Middle: 2, Nodes_Visit: []int{1,5,2}}
	actual := metrics_Calc(packet, topology)
	
	//Latency
	want := 30.0
	if actual[0] != want{
		t.Errorf("Latency -- Actual: %v -- Want: %v", actual[0], want)
	}
	
	//Packet Loss
	want = 1.0
	if actual[1] != want{
		t.Errorf("Packet Loss -- Actual: %v -- Want: %v", actual[1], want)
	}
	
	//Path Stretch
	want = 1.0
	if actual[2] != want{
		t.Errorf("Path Stretch -- Actual: %v -- Want: %v", actual[2], want)
	}
	
	//Packet Size
	want = 56.0
	if actual[3] != want{
		t.Errorf("Packet Size -- Actual: %v -- Want: %v", actual[3], want)
	}
	
	//Pressure Count
	want = 0.5
	if actual[4] != want{
		t.Errorf("Pressure Count -- Actual: %v -- Want: %v", actual[4], want)
	}
	
	//Pressure Percentage
	want = 1.0
	if actual[5] != want{
		t.Errorf("Pressure Percentage -- Actual: %v -- Want: %v", actual[5], want)
	}
	
	//Nodes Visited
	want = 0.5
	if actual[6] != want{
		t.Errorf("Nodes Visited -- Actual: %v -- Want: %v", actual[6], want)
	}
	
	//Hops Collapsed
	want = 0.0
	if actual[7] != want{
		t.Errorf("Hops Collapsed -- Actual: %v -- Want: %v", actual[7], want)
	}
}
//----------------------------------------------

//----------------------------------------------
// Tests to make sure metrics are calculated correctly when interest collapsing
//----------------------------------------------
func TestMetrics_Calc_Interest(t *testing.T) {
	time_now := time.Now()
	topology := read_Topology("./unit_test_topologies/test_topology.csv", "localhost", 8000, 1)
	packet := Packet{Start_Time: (time_now).Add(-time.Duration(30 * time.Second)), End_Time: time_now, Start: 1, Destination: 2, Path_Traversed: []int{0, 1, 2, 3, 4, 3, 2, 1, 0}, Collapsed: 1, Pressure_Gauge: 0, Index_Middle: 1, Nodes_Visit: []int{1,5}}
	actual := metrics_Calc(packet, topology)
	
	//Latency
	want := 30.0
	if actual[0] != want{
		t.Errorf("Latency -- Actual: %v -- Want: %v", actual[0], want)
	}
	
	//Packet Loss
	want = 1.0
	if actual[1] != want{
		t.Errorf("Packet Loss -- Actual: %v -- Want: %v", actual[1], want)
	}
	
	//Path Stretch
	want = 0.5
	if actual[2] != want{
		t.Errorf("Path Stretch -- Actual: %v -- Want: %v", actual[2], want)
	}
	
	//Packet Size
	want = 52.0
	if actual[3] != want{
		t.Errorf("Packet Size -- Actual: %v -- Want: %v", actual[3], want)
	}
	
	//Pressure Count
	want = -1.0
	if actual[4] != want{
		t.Errorf("Pressure Count -- Actual: %v -- Want: %v", actual[4], want)
	}
	
	//Pressure Percentage
	want = 0.0
	if actual[5] != want{
		t.Errorf("Pressure Percentage -- Actual: %v -- Want: %v", actual[5], want)
	}
	
	//Nodes Visited
	want = 0.3333333333333333
	if actual[6] != want{
		t.Errorf("Nodes Visited -- Actual: %v -- Want: %v", actual[6], want)
	}
	
	//Hops Collapsed
	want = 0.75
	if actual[7] != want{
		t.Errorf("Hops Collapsed -- Actual: %v -- Want: %v", actual[7], want)
	}
}
//----------------------------------------------

//--------------------------------------------------------------------------------------------
// algorithm_Check(algorithm string) (bool, bool, bool, bool, bool)
//--------------------------------------------------------------------------------------------

//----------------------------------------------
// Tests to make sure for each individual algorithm
//----------------------------------------------
func TestAlgorithm_Check(t *testing.T) {

	//FILO_GF
	actual_1, actual_2, actual_3, actual_4, actual_5 := algorithm_Check("FILO_GF")
	want_1 := true
	want_2 := false
	want_3 := false
	want_4 := false
	want_5 := false
	if actual_1 != want_1{
		t.Errorf("FILO -- Actual: %v -- Want: %v", actual_1, want_1)
	}
	if actual_2 != want_2{
		t.Errorf("Global -- Actual: %v -- Want: %v", actual_2, want_2)
	}
	if actual_3 != want_3{
		t.Errorf("Modified -- Actual: %v -- Want: %v", actual_3, want_3)
	}
	if actual_4 != want_4{
		t.Errorf("Gravity Pressure -- Actual: %v -- Want: %v", actual_4, want_4)
	}
	if actual_5 != want_5{
		t.Errorf("Multi -- Actual: %v -- Want: %v", actual_5, want_5)
	}
	
	//ALL_GF
	actual_1, actual_2, actual_3, actual_4, actual_5 = algorithm_Check("ALL_GF")
	want_1 = false
	want_2 = false
	want_3 = false
	want_4 = false
	want_5 = false
	if actual_1 != want_1{
		t.Errorf("FILO -- Actual: %v -- Want: %v", actual_1, want_1)
	}
	if actual_2 != want_2{
		t.Errorf("Global -- Actual: %v -- Want: %v", actual_2, want_2)
	}
	if actual_3 != want_3{
		t.Errorf("Modified -- Actual: %v -- Want: %v", actual_3, want_3)
	}
	if actual_4 != want_4{
		t.Errorf("Gravity Pressure -- Actual: %v -- Want: %v", actual_4, want_4)
	}
	if actual_5 != want_5{
		t.Errorf("Multi -- Actual: %v -- Want: %v", actual_5, want_5)
	}
		
	//FILO_GPGF
	actual_1, actual_2, actual_3, actual_4, actual_5 = algorithm_Check("FILO_GPGF")
	want_1 = true
	want_2 = false
	want_3 = false
	want_4 = true
	want_5 = false
	if actual_1 != want_1{
		t.Errorf("FILO -- Actual: %v -- Want: %v", actual_1, want_1)
	}
	if actual_2 != want_2{
		t.Errorf("Global -- Actual: %v -- Want: %v", actual_2, want_2)
	}
	if actual_3 != want_3{
		t.Errorf("Modified -- Actual: %v -- Want: %v", actual_3, want_3)
	}
	if actual_4 != want_4{
		t.Errorf("Gravity Pressure -- Actual: %v -- Want: %v", actual_4, want_4)
	}
	if actual_5 != want_5{
		t.Errorf("Multi -- Actual: %v -- Want: %v", actual_5, want_5)
	}
		
	//ALL_GPGF
	actual_1, actual_2, actual_3, actual_4, actual_5 = algorithm_Check("ALL_GPGF")
	want_1 = false
	want_2 = false
	want_3 = false
	want_4 = true
	want_5 = false
	if actual_1 != want_1{
		t.Errorf("FILO -- Actual: %v -- Want: %v", actual_1, want_1)
	}
	if actual_2 != want_2{
		t.Errorf("Global -- Actual: %v -- Want: %v", actual_2, want_2)
	}
	if actual_3 != want_3{
		t.Errorf("Modified -- Actual: %v -- Want: %v", actual_3, want_3)
	}
	if actual_4 != want_4{
		t.Errorf("Gravity Pressure -- Actual: %v -- Want: %v", actual_4, want_4)
	}
	if actual_5 != want_5{
		t.Errorf("Multi -- Actual: %v -- Want: %v", actual_5, want_5)
	}
		
	//G_FILO_GPGF
	actual_1, actual_2, actual_3, actual_4, actual_5 = algorithm_Check("G_FILO_GPGF")
	want_1 = true
	want_2 = true
	want_3 = false
	want_4 = true
	want_5 = false
	if actual_1 != want_1{
		t.Errorf("FILO -- Actual: %v -- Want: %v", actual_1, want_1)
	}
	if actual_2 != want_2{
		t.Errorf("Global -- Actual: %v -- Want: %v", actual_2, want_2)
	}
	if actual_3 != want_3{
		t.Errorf("Modified -- Actual: %v -- Want: %v", actual_3, want_3)
	}
	if actual_4 != want_4{
		t.Errorf("Gravity Pressure -- Actual: %v -- Want: %v", actual_4, want_4)
	}
	if actual_5 != want_5{
		t.Errorf("Multi -- Actual: %v -- Want: %v", actual_5, want_5)
	}
		
	//G_ALL_GPGF
	actual_1, actual_2, actual_3, actual_4, actual_5 = algorithm_Check("G_ALL_GPGF")
	want_1 = false
	want_2 = true
	want_3 = false
	want_4 = true
	want_5 = false
	if actual_1 != want_1{
		t.Errorf("FILO -- Actual: %v -- Want: %v", actual_1, want_1)
	}
	if actual_2 != want_2{
		t.Errorf("Global -- Actual: %v -- Want: %v", actual_2, want_2)
	}
	if actual_3 != want_3{
		t.Errorf("Modified -- Actual: %v -- Want: %v", actual_3, want_3)
	}
	if actual_4 != want_4{
		t.Errorf("Gravity Pressure -- Actual: %v -- Want: %v", actual_4, want_4)
	}	
	if actual_5 != want_5{
		t.Errorf("Multi -- Actual: %v -- Want: %v", actual_5, want_5)
	}
	
	//FILO_mGF
	actual_1, actual_2, actual_3, actual_4, actual_5 = algorithm_Check("FILO_mGF")
	want_1 = true
	want_2 = false
	want_3 = true
	want_4 = false
	want_5 = false
	if actual_1 != want_1{
		t.Errorf("FILO -- Actual: %v -- Want: %v", actual_1, want_1)
	}
	if actual_2 != want_2{
		t.Errorf("Global -- Actual: %v -- Want: %v", actual_2, want_2)
	}
	if actual_3 != want_3{
		t.Errorf("Modified -- Actual: %v -- Want: %v", actual_3, want_3)
	}
	if actual_4 != want_4{
		t.Errorf("Gravity Pressure -- Actual: %v -- Want: %v", actual_4, want_4)
	}
		if actual_5 != want_5{
		t.Errorf("Multi -- Actual: %v -- Want: %v", actual_5, want_5)
	}
	
	//ALL_mGF
	actual_1, actual_2, actual_3, actual_4, actual_5 = algorithm_Check("ALL_mGF")
	want_1 = false
	want_2 = false
	want_3 = true
	want_4 = false
	want_5 = false
	if actual_1 != want_1{
		t.Errorf("FILO -- Actual: %v -- Want: %v", actual_1, want_1)
	}
	if actual_2 != want_2{
		t.Errorf("Global -- Actual: %v -- Want: %v", actual_2, want_2)
	}
	if actual_3 != want_3{
		t.Errorf("Modified -- Actual: %v -- Want: %v", actual_3, want_3)
	}
	if actual_4 != want_4{
		t.Errorf("Gravity Pressure -- Actual: %v -- Want: %v", actual_4, want_4)
	}
	if actual_5 != want_5{
		t.Errorf("Multi -- Actual: %v -- Want: %v", actual_5, want_5)
	}
	
	//FILO_GPmGF
	actual_1, actual_2, actual_3, actual_4, actual_5 = algorithm_Check("FILO_GPmGF")
	want_1 = true
	want_2 = false
	want_3 = true
	want_4 = true
	want_5 = false
	if actual_1 != want_1{
		t.Errorf("FILO -- Actual: %v -- Want: %v", actual_1, want_1)
	}
	if actual_2 != want_2{
		t.Errorf("Global -- Actual: %v -- Want: %v", actual_2, want_2)
	}
	if actual_3 != want_3{
		t.Errorf("Modified -- Actual: %v -- Want: %v", actual_3, want_3)
	}
	if actual_4 != want_4{
		t.Errorf("Gravity Pressure -- Actual: %v -- Want: %v", actual_4, want_4)
	}
	if actual_5 != want_5{
		t.Errorf("Multi -- Actual: %v -- Want: %v", actual_5, want_5)
	}
	
	//ALL_GPmGF
	actual_1, actual_2, actual_3, actual_4, actual_5 = algorithm_Check("ALL_GPmGF")
	want_1 = false
	want_2 = false
	want_3 = true
	want_4 = true
	want_5 = false
	if actual_1 != want_1{
		t.Errorf("FILO -- Actual: %v -- Want: %v", actual_1, want_1)
	}
	if actual_2 != want_2{
		t.Errorf("Global -- Actual: %v -- Want: %v", actual_2, want_2)
	}
	if actual_3 != want_3{
		t.Errorf("Modified -- Actual: %v -- Want: %v", actual_3, want_3)
	}
	if actual_4 != want_4{
		t.Errorf("Gravity Pressure -- Actual: %v -- Want: %v", actual_4, want_4)
	}
	if actual_5 != want_5{
		t.Errorf("Multi -- Actual: %v -- Want: %v", actual_5, want_5)
	}
		
	//G_FILO_GPmGF
	actual_1, actual_2, actual_3, actual_4, actual_5 = algorithm_Check("G_FILO_GPmGF")
	want_1 = true
	want_2 = true
	want_3 = true
	want_4 = true
	want_5 = false
	if actual_1 != want_1{
		t.Errorf("FILO -- Actual: %v -- Want: %v", actual_1, want_1)
	}
	if actual_2 != want_2{
		t.Errorf("Global -- Actual: %v -- Want: %v", actual_2, want_2)
	}
	if actual_3 != want_3{
		t.Errorf("Modified -- Actual: %v -- Want: %v", actual_3, want_3)
	}
	if actual_4 != want_4{
		t.Errorf("Gravity Pressure -- Actual: %v -- Want: %v", actual_4, want_4)
	}
	if actual_5 != want_5{
		t.Errorf("Multi -- Actual: %v -- Want: %v", actual_5, want_5)
	}
		
	//G_ALL_GPmGF
	actual_1, actual_2, actual_3, actual_4, actual_5 = algorithm_Check("G_ALL_GPmGF")
	want_1 = false
	want_2 = true
	want_3 = true
	want_4 = true
	want_5 = false
	if actual_1 != want_1{
		t.Errorf("FILO -- Actual: %v -- Want: %v", actual_1, want_1)
	}
	if actual_2 != want_2{
		t.Errorf("Global -- Actual: %v -- Want: %v", actual_2, want_2)
	}
	if actual_3 != want_3{
		t.Errorf("Modified -- Actual: %v -- Want: %v", actual_3, want_3)
	}
	if actual_4 != want_4{
		t.Errorf("Gravity Pressure -- Actual: %v -- Want: %v", actual_4, want_4)
	}	
	if actual_5 != want_5{
		t.Errorf("Multi -- Actual: %v -- Want: %v", actual_5, want_5)
	}
		
	//Flooding
	actual_1, actual_2, actual_3, actual_4, actual_5 = algorithm_Check("Flooding")
	want_1 = false
	want_2 = false
	want_3 = false
	want_4 = false
	want_5 = false
	if actual_1 != want_1{
		t.Errorf("FILO -- Actual: %v -- Want: %v", actual_1, want_1)
	}
	if actual_2 != want_2{
		t.Errorf("Global -- Actual: %v -- Want: %v", actual_2, want_2)
	}
	if actual_3 != want_3{
		t.Errorf("Modified -- Actual: %v -- Want: %v", actual_3, want_3)
	}
	if actual_4 != want_4{
		t.Errorf("Gravity Pressure -- Actual: %v -- Want: %v", actual_4, want_4)
	}
	if actual_5 != want_5{
		t.Errorf("Multi -- Actual: %v -- Want: %v", actual_5, want_5)
	}
		
	//Multi_Flooding
	actual_1, actual_2, actual_3, actual_4, actual_5 = algorithm_Check("Multi_Flooding")
	want_1 = false
	want_2 = false
	want_3 = false
	want_4 = false
	want_5 = true
	if actual_1 != want_1{
		t.Errorf("FILO -- Actual: %v -- Want: %v", actual_1, want_1)
	}
	if actual_2 != want_2{
		t.Errorf("Global -- Actual: %v -- Want: %v", actual_2, want_2)
	}
	if actual_3 != want_3{
		t.Errorf("Modified -- Actual: %v -- Want: %v", actual_3, want_3)
	}
	if actual_4 != want_4{
		t.Errorf("Gravity Pressure -- Actual: %v -- Want: %v", actual_4, want_4)
	}
	if actual_5 != want_5{
		t.Errorf("Multi -- Actual: %v -- Want: %v", actual_5, want_5)
	}
	
	//Multi_FILO_GF
	actual_1, actual_2, actual_3, actual_4, actual_5 = algorithm_Check("Multi_FILO_GF")
	want_1 = true
	want_2 = false
	want_3 = false
	want_4 = false
	want_5 = true
	if actual_1 != want_1{
		t.Errorf("FILO -- Actual: %v -- Want: %v", actual_1, want_1)
	}
	if actual_2 != want_2{
		t.Errorf("Global -- Actual: %v -- Want: %v", actual_2, want_2)
	}
	if actual_3 != want_3{
		t.Errorf("Modified -- Actual: %v -- Want: %v", actual_3, want_3)
	}
	if actual_4 != want_4{
		t.Errorf("Gravity Pressure -- Actual: %v -- Want: %v", actual_4, want_4)
	}
	if actual_5 != want_5{
		t.Errorf("Multi -- Actual: %v -- Want: %v", actual_5, want_5)
	}
	
	//ALL_GF
	actual_1, actual_2, actual_3, actual_4, actual_5 = algorithm_Check("Multi_ALL_GF")
	want_1 = false
	want_2 = false
	want_3 = false
	want_4 = false
	want_5 = true
	if actual_1 != want_1{
		t.Errorf("FILO -- Actual: %v -- Want: %v", actual_1, want_1)
	}
	if actual_2 != want_2{
		t.Errorf("Global -- Actual: %v -- Want: %v", actual_2, want_2)
	}
	if actual_3 != want_3{
		t.Errorf("Modified -- Actual: %v -- Want: %v", actual_3, want_3)
	}
	if actual_4 != want_4{
		t.Errorf("Gravity Pressure -- Actual: %v -- Want: %v", actual_4, want_4)
	}
	if actual_5 != want_5{
		t.Errorf("Multi -- Actual: %v -- Want: %v", actual_5, want_5)
	}
		
	//FILO_GPGF
	actual_1, actual_2, actual_3, actual_4, actual_5 = algorithm_Check("Multi_FILO_GPGF")
	want_1 = true
	want_2 = false
	want_3 = false
	want_4 = true
	want_5 = true
	if actual_1 != want_1{
		t.Errorf("FILO -- Actual: %v -- Want: %v", actual_1, want_1)
	}
	if actual_2 != want_2{
		t.Errorf("Global -- Actual: %v -- Want: %v", actual_2, want_2)
	}
	if actual_3 != want_3{
		t.Errorf("Modified -- Actual: %v -- Want: %v", actual_3, want_3)
	}
	if actual_4 != want_4{
		t.Errorf("Gravity Pressure -- Actual: %v -- Want: %v", actual_4, want_4)
	}
	if actual_5 != want_5{
		t.Errorf("Multi -- Actual: %v -- Want: %v", actual_5, want_5)
	}
		
	//ALL_GPGF
	actual_1, actual_2, actual_3, actual_4, actual_5 = algorithm_Check("Multi_ALL_GPGF")
	want_1 = false
	want_2 = false
	want_3 = false
	want_4 = true
	want_5 = true
	if actual_1 != want_1{
		t.Errorf("FILO -- Actual: %v -- Want: %v", actual_1, want_1)
	}
	if actual_2 != want_2{
		t.Errorf("Global -- Actual: %v -- Want: %v", actual_2, want_2)
	}
	if actual_3 != want_3{
		t.Errorf("Modified -- Actual: %v -- Want: %v", actual_3, want_3)
	}
	if actual_4 != want_4{
		t.Errorf("Gravity Pressure -- Actual: %v -- Want: %v", actual_4, want_4)
	}
	if actual_5 != want_5{
		t.Errorf("Multi -- Actual: %v -- Want: %v", actual_5, want_5)
	}
		
	//G_FILO_GPGF
	actual_1, actual_2, actual_3, actual_4, actual_5 = algorithm_Check("Multi_G_FILO_GPGF")
	want_1 = true
	want_2 = true
	want_3 = false
	want_4 = true
	want_5 = true
	if actual_1 != want_1{
		t.Errorf("FILO -- Actual: %v -- Want: %v", actual_1, want_1)
	}
	if actual_2 != want_2{
		t.Errorf("Global -- Actual: %v -- Want: %v", actual_2, want_2)
	}
	if actual_3 != want_3{
		t.Errorf("Modified -- Actual: %v -- Want: %v", actual_3, want_3)
	}
	if actual_4 != want_4{
		t.Errorf("Gravity Pressure -- Actual: %v -- Want: %v", actual_4, want_4)
	}
	if actual_5 != want_5{
		t.Errorf("Multi -- Actual: %v -- Want: %v", actual_5, want_5)
	}
		
	//G_ALL_GPGF
	actual_1, actual_2, actual_3, actual_4, actual_5 = algorithm_Check("Multi_G_ALL_GPGF")
	want_1 = false
	want_2 = true
	want_3 = false
	want_4 = true
	want_5 = true
	if actual_1 != want_1{
		t.Errorf("FILO -- Actual: %v -- Want: %v", actual_1, want_1)
	}
	if actual_2 != want_2{
		t.Errorf("Global -- Actual: %v -- Want: %v", actual_2, want_2)
	}
	if actual_3 != want_3{
		t.Errorf("Modified -- Actual: %v -- Want: %v", actual_3, want_3)
	}
	if actual_4 != want_4{
		t.Errorf("Gravity Pressure -- Actual: %v -- Want: %v", actual_4, want_4)
	}	
	if actual_5 != want_5{
		t.Errorf("Multi -- Actual: %v -- Want: %v", actual_5, want_5)
	}
	
	//FILO_mGF
	actual_1, actual_2, actual_3, actual_4, actual_5 = algorithm_Check("Multi_FILO_mGF")
	want_1 = true
	want_2 = false
	want_3 = true
	want_4 = false
	want_5 = true
	if actual_1 != want_1{
		t.Errorf("FILO -- Actual: %v -- Want: %v", actual_1, want_1)
	}
	if actual_2 != want_2{
		t.Errorf("Global -- Actual: %v -- Want: %v", actual_2, want_2)
	}
	if actual_3 != want_3{
		t.Errorf("Modified -- Actual: %v -- Want: %v", actual_3, want_3)
	}
	if actual_4 != want_4{
		t.Errorf("Gravity Pressure -- Actual: %v -- Want: %v", actual_4, want_4)
	}
		if actual_5 != want_5{
		t.Errorf("Multi -- Actual: %v -- Want: %v", actual_5, want_5)
	}
	
	//ALL_mGF
	actual_1, actual_2, actual_3, actual_4, actual_5 = algorithm_Check("Multi_ALL_mGF")
	want_1 = false
	want_2 = false
	want_3 = true
	want_4 = false
	want_5 = true
	if actual_1 != want_1{
		t.Errorf("FILO -- Actual: %v -- Want: %v", actual_1, want_1)
	}
	if actual_2 != want_2{
		t.Errorf("Global -- Actual: %v -- Want: %v", actual_2, want_2)
	}
	if actual_3 != want_3{
		t.Errorf("Modified -- Actual: %v -- Want: %v", actual_3, want_3)
	}
	if actual_4 != want_4{
		t.Errorf("Gravity Pressure -- Actual: %v -- Want: %v", actual_4, want_4)
	}
	if actual_5 != want_5{
		t.Errorf("Multi -- Actual: %v -- Want: %v", actual_5, want_5)
	}
	
	//FILO_GPmGF
	actual_1, actual_2, actual_3, actual_4, actual_5 = algorithm_Check("Multi_FILO_GPmGF")
	want_1 = true
	want_2 = false
	want_3 = true
	want_4 = true
	want_5 = true
	if actual_1 != want_1{
		t.Errorf("FILO -- Actual: %v -- Want: %v", actual_1, want_1)
	}
	if actual_2 != want_2{
		t.Errorf("Global -- Actual: %v -- Want: %v", actual_2, want_2)
	}
	if actual_3 != want_3{
		t.Errorf("Modified -- Actual: %v -- Want: %v", actual_3, want_3)
	}
	if actual_4 != want_4{
		t.Errorf("Gravity Pressure -- Actual: %v -- Want: %v", actual_4, want_4)
	}
	if actual_5 != want_5{
		t.Errorf("Multi -- Actual: %v -- Want: %v", actual_5, want_5)
	}
	
	//ALL_GPmGF
	actual_1, actual_2, actual_3, actual_4, actual_5 = algorithm_Check("Multi_ALL_GPmGF")
	want_1 = false
	want_2 = false
	want_3 = true
	want_4 = true
	want_5 = true
	if actual_1 != want_1{
		t.Errorf("FILO -- Actual: %v -- Want: %v", actual_1, want_1)
	}
	if actual_2 != want_2{
		t.Errorf("Global -- Actual: %v -- Want: %v", actual_2, want_2)
	}
	if actual_3 != want_3{
		t.Errorf("Modified -- Actual: %v -- Want: %v", actual_3, want_3)
	}
	if actual_4 != want_4{
		t.Errorf("Gravity Pressure -- Actual: %v -- Want: %v", actual_4, want_4)
	}
	if actual_5 != want_5{
		t.Errorf("Multi -- Actual: %v -- Want: %v", actual_5, want_5)
	}
		
	//G_FILO_GPmGF
	actual_1, actual_2, actual_3, actual_4, actual_5 = algorithm_Check("Multi_G_FILO_GPmGF")
	want_1 = true
	want_2 = true
	want_3 = true
	want_4 = true
	want_5 = true
	if actual_1 != want_1{
		t.Errorf("FILO -- Actual: %v -- Want: %v", actual_1, want_1)
	}
	if actual_2 != want_2{
		t.Errorf("Global -- Actual: %v -- Want: %v", actual_2, want_2)
	}
	if actual_3 != want_3{
		t.Errorf("Modified -- Actual: %v -- Want: %v", actual_3, want_3)
	}
	if actual_4 != want_4{
		t.Errorf("Gravity Pressure -- Actual: %v -- Want: %v", actual_4, want_4)
	}
	if actual_5 != want_5{
		t.Errorf("Multi -- Actual: %v -- Want: %v", actual_5, want_5)
	}
		
	//G_ALL_GPmGF
	actual_1, actual_2, actual_3, actual_4, actual_5 = algorithm_Check("Multi_G_ALL_GPmGF")
	want_1 = false
	want_2 = true
	want_3 = true
	want_4 = true
	want_5 = true
	if actual_1 != want_1{
		t.Errorf("FILO -- Actual: %v -- Want: %v", actual_1, want_1)
	}
	if actual_2 != want_2{
		t.Errorf("Global -- Actual: %v -- Want: %v", actual_2, want_2)
	}
	if actual_3 != want_3{
		t.Errorf("Modified -- Actual: %v -- Want: %v", actual_3, want_3)
	}
	if actual_4 != want_4{
		t.Errorf("Gravity Pressure -- Actual: %v -- Want: %v", actual_4, want_4)
	}	
	if actual_5 != want_5{
		t.Errorf("Multi -- Actual: %v -- Want: %v", actual_5, want_5)
	}
}
//----------------------------------------------

//----------------------------------------------
// Tests to make sure unknown algorithm exits
//----------------------------------------------
func TestAlgorithm_Check_Fail(t *testing.T) {

	//Run code via cmd
	if os.Getenv("FLAG") == "1" {
		_, _, _, _,_ = algorithm_Check("UNKNOWN")
	}
	
	//Execute the code via exec to check error code
	cmd := exec.Command(os.Args[0], "-test.run=TestAlgorithm_Check_Fail")
	cmd.Env = append(os.Environ(), "FLAG=1")
	output, err := cmd.Output() //grab error code and output
	
	//Compare error codes
	actual := err.Error()
	want := "exit status 1"
	if actual != want{       
		t.Errorf("Actual: %s -- Want: %s", actual, want)
	}
	
	//Compare output
	actual = string(output)
	want = "Unrecognized algorithm! Exiting.\n"
	if actual != want{       
		t.Errorf("Actual: %s -- Want: %s", actual, want)
	}
}
//----------------------------------------------

//--------------------------------------------------------------------------------------------
// check_Cache_Hit(curr_Node *Node, name string) (*CACHE, bool)
//--------------------------------------------------------------------------------------------

//----------------------------------------------
// Tests to make sure cache hit works
//----------------------------------------------
func TestCheck_Cache_Hit(t *testing.T) {
	node_A := &Node{Number: 0}
	var cache_A []CACHE
	cache_A = append(cache_A, CACHE{Name: "Miss", Timestamp: time.Now(), TTL: time.Duration(time.Second * 1000), Last_Used: time.Time{}, Frequency: 0})
	cache_A = append(cache_A, CACHE{Name: "Hit", Timestamp: time.Now(), TTL: time.Duration(time.Second * 1000), Last_Used: time.Time{}, Frequency: 0})
	cache_A = append(cache_A, CACHE{Name: "Miss_2", Timestamp: time.Now(), TTL: time.Duration(time.Second * 1000), Last_Used: time.Time{}, Frequency: 0})
	node_A.Cache = cache_A
	actual_cache, actual_hit := check_Cache_Hit(node_A, "Hit")
	
	//Compare output
	want := cache_A[1]
	if *actual_cache != want{       
		t.Errorf("Actual: %v -- Want: %v", actual_cache, want)
	}
	
	//Compare output
	want_2 := true
	if actual_hit != want_2{       
		t.Errorf("Actual: %v -- Want: %v", actual_hit, want_2)
	}
	
	//Compare output
	want_3 := time.Time{}
	if actual_cache.Last_Used == want_3{       
		t.Errorf("Actual: %v -- Dont Want: %v", actual_cache.Last_Used, want_3)
	}
	
	//Compare output
	want_4 := 1
	if actual_cache.Frequency != want_4{       
		t.Errorf("Actual: %v -- Want: %v", actual_cache.Frequency, want_4)
	}
}
//----------------------------------------------

//----------------------------------------------
// Tests to make sure cache hit works if TTL is -1
//----------------------------------------------
func TestCheck_Cache_Hit_TTL(t *testing.T) {
	node_A := &Node{Number: 0}
	var cache_A []CACHE
	cache_A = append(cache_A, CACHE{Name: "Miss", Timestamp: time.Now(), TTL: time.Duration(time.Second * 1000), Last_Used: time.Time{}, Frequency: 0})
	cache_A = append(cache_A, CACHE{Name: "Hit", Timestamp: time.Now(), TTL: time.Duration(-1), Last_Used: time.Time{}, Frequency: 0})
	cache_A = append(cache_A, CACHE{Name: "Miss_2", Timestamp: time.Now(), TTL: time.Duration(time.Second * 1000), Last_Used: time.Time{}, Frequency: 0})
	node_A.Cache = cache_A
	actual_cache, actual_hit := check_Cache_Hit(node_A, "Hit")
	
	//Compare output
	want := &cache_A[1]
	if actual_cache != want{       
		t.Errorf("Actual: %v -- Want: %v", actual_cache, want)
	}
	
	//Compare output
	want_2 := true
	if actual_hit != want_2{       
		t.Errorf("Actual: %v -- Want: %v", actual_hit, want_2)
	}
	
	//Compare output
	want_3 := time.Time{}
	if actual_cache.Last_Used == want_3{       
		t.Errorf("Actual: %v -- Dont Want: %v", actual_cache.Last_Used, want_3)
	}
	
	//Compare output
	want_4 := 1
	if actual_cache.Frequency != want_4{       
		t.Errorf("Actual: %v -- Want: %v", actual_cache.Frequency, want_4)
	}
}
//----------------------------------------------

//----------------------------------------------
// Tests to make sure cache hit fails if no matching entry
//----------------------------------------------
func TestCheck_Cache_Hit_Fail(t *testing.T) {
	node_A := &Node{Number: 0}
	var cache_A []CACHE
	cache_A = append(cache_A, CACHE{Name: "Miss_0", Timestamp: time.Now(), TTL: time.Duration(time.Second * 1000), Last_Used: time.Time{}, Frequency: 0})
	cache_A = append(cache_A, CACHE{Name: "Miss_1", Timestamp: time.Now(), TTL: time.Duration(time.Second * 1000), Last_Used: time.Time{}, Frequency: 0})
	cache_A = append(cache_A, CACHE{Name: "Miss_2", Timestamp: time.Now(), TTL: time.Duration(time.Second * 1000), Last_Used: time.Time{}, Frequency: 0})
	node_A.Cache = cache_A
	actual_cache, actual_hit := check_Cache_Hit(node_A, "Hit")
	
	//Compare output
	want := time.Time{}
	if actual_cache.Timestamp != want{       
		t.Errorf("Actual: %v -- Want: %v", actual_cache.Timestamp, want)
	}
	
	//Compare output
	want_2 := false
	if actual_hit != want_2{       
		t.Errorf("Actual: %v -- Want: %v", actual_hit, want_2)
	}
}
//----------------------------------------------

//----------------------------------------------
// Tests to make sure cache hit fails if bad TTL
//----------------------------------------------
func TestCheck_Cache_Hit_Fail_TTL(t *testing.T) {
	node_A := &Node{Number: 0}
	time_Now := time.Now()
	var cache_A []CACHE
	cache_A = append(cache_A, CACHE{Name: "Miss_0", Timestamp: time.Now(), TTL: time.Duration(time.Second * 1000), Last_Used: time.Time{}, Frequency: 0})
	cache_A = append(cache_A, CACHE{Name: "Hit", Timestamp: time_Now, TTL: time.Duration(1), Last_Used: time.Time{}, Frequency: 0})
	cache_A = append(cache_A, CACHE{Name: "Miss_2", Timestamp: time.Now(), TTL: time.Duration(time.Second * 1000), Last_Used: time.Time{}, Frequency: 0})
	node_A.Cache = cache_A
	actual_cache, actual_hit := check_Cache_Hit(node_A, "Hit")
	
	//Compare output
	want := time_Now
	if actual_cache.Timestamp != want{       
		t.Errorf("Actual: %v -- Want: %v", actual_cache.Timestamp, want)
	}
	
	//Compare output
	want_2 := false
	if actual_hit != want_2{       
		t.Errorf("Actual: %v -- Want: %v", actual_hit, want_2)
	}
}
//----------------------------------------------

//--------------------------------------------------------------------------------------------

//--------------------------------------------------------------------------------------------
// cache_Data(curr_Node *Node, packet Packet, cache_Evict string)
//--------------------------------------------------------------------------------------------

//----------------------------------------------
// Tests to make sure cache data works if cache is empty
//----------------------------------------------
func TestCache_Data_Empty(t *testing.T) {
	node_A := &Node{Number: 0, Cache_Size: 10}
	packet_A := Packet{Name: "Miss", Payload: "Hello", TTL: time.Duration(10)}
	cache_Data(node_A, packet_A, "LFU")
	
	//Compare output
	want := 1
	if len(node_A.Cache) != want{       
		t.Errorf("Actual: %v -- Want: %v", len(node_A.Cache), want)
	}
	
	//Compare output
	want_2 := "Miss"
	if node_A.Cache[0].Name != want_2{       
		t.Errorf("Actual: %v -- Want: %v", node_A.Cache[0].Name, want_2)
	}
	
	//Compare output
	want_2 = "Hello"
	if node_A.Cache[0].Data != want_2{       
		t.Errorf("Actual: %v -- Want: %v", node_A.Cache[0].Data, want_2)
	}
	
	//Compare output
	want_3 := time.Duration(10)
	if node_A.Cache[0].TTL != want_3{      
		t.Errorf("Actual: %v -- Want: %v", node_A.Cache[0].TTL, want_3)
	}	
}
//----------------------------------------------

//----------------------------------------------
// Tests to make sure cache data works if cache is not empty
//----------------------------------------------
func TestCache_Data_Not_Empty(t *testing.T) {
	node_A := &Node{Number: 0, Cache_Size: 10}
	packet_A := Packet{Name: "Miss", Payload: "Hello", TTL: time.Duration(10)}
	var cache_A []CACHE
	cache_A = append(cache_A, CACHE{Name: "Hit_0", Timestamp: time.Now(), TTL: time.Duration(time.Second * 1000), Last_Used: time.Time{}, Frequency: 1})
	cache_A = append(cache_A, CACHE{Name: "Hit_1", Timestamp: time.Now(), TTL: time.Duration(time.Second * 1000), Last_Used: time.Time{}, Frequency: 0})
	cache_A = append(cache_A, CACHE{Name: "Hit_2", Timestamp: time.Now(), TTL: time.Duration(time.Second * 1000), Last_Used: time.Time{}, Frequency: 1})
	node_A.Cache = cache_A
	cache_Data(node_A, packet_A, "LFU")
	
	//Compare output
	want := 4
	if len(node_A.Cache) != want{       
		t.Errorf("Actual: %v -- Want: %v", len(node_A.Cache), want)
	}
	
	//Compare output
	want_2 := "Miss"
	if node_A.Cache[3].Name != want_2{       
		t.Errorf("Actual: %v -- Want: %v", node_A.Cache[3].Name, want_2)
	}
	
	//Compare output
	want_2 = "Hello"
	if node_A.Cache[3].Data != want_2{       
		t.Errorf("Actual: %v -- Want: %v", node_A.Cache[3].Data, want_2)
	}
	
	//Compare output
	want_3 := time.Duration(10)
	if node_A.Cache[3].TTL != want_3{       
		t.Errorf("Actual: %v -- Want: %v", node_A.Cache[3].TTL, want_3)
	}		
}
//----------------------------------------------

//----------------------------------------------
// Tests to make sure cache data works if cache has matching data
//----------------------------------------------
func TestCache_Data_Matching(t *testing.T) {
	node_A := &Node{Number: 0, Cache_Size: 10}
	packet_A := Packet{Name: "Hit_1", Payload: "Hello", TTL: time.Duration(10)}
	var cache_A []CACHE
	cache_A = append(cache_A, CACHE{Name: "Hit_0", Timestamp: time.Now(), TTL: time.Duration(time.Second * 1000), Last_Used: time.Time{}, Frequency: 1})
	cache_A = append(cache_A, CACHE{Name: "Hit_1", Timestamp: time.Now(), TTL: time.Duration(time.Second * 1000), Last_Used: time.Time{}, Frequency: 0})
	cache_A = append(cache_A, CACHE{Name: "Hit_2", Timestamp: time.Now(), TTL: time.Duration(time.Second * 1000), Last_Used: time.Time{}, Frequency: 1})
	node_A.Cache = cache_A
	cache_Data(node_A, packet_A, "LFU")
	
	//Compare output
	want := 3
	if len(node_A.Cache) != want{       
		t.Errorf("Actual: %v -- Want: %v", len(node_A.Cache), want)
	}
	
	//Compare output
	want_2 := "Hello"
	if node_A.Cache[2].Data != want_2{       
		t.Errorf("Actual: %v -- Want: %v", node_A.Cache[2].Data, want_2)
	}
	
	//Compare output
	want_3 := time.Duration(10)
	if node_A.Cache[2].TTL != want_3{       
		t.Errorf("Actual: %v -- Want: %v", node_A.Cache[2].TTL, want_3)
	}	
}
//----------------------------------------------

//----------------------------------------------
// Tests to make sure cache data works if cache is full
//----------------------------------------------
func TestCache_Data_Full(t *testing.T) {
	node_A := &Node{Number: 0, Cache_Size: 3}
	packet_A := Packet{Name: "Miss", Payload: "Hello", TTL: time.Duration(10)}
	var cache_A []CACHE
	cache_A = append(cache_A, CACHE{Name: "Hit_0", Timestamp: time.Now(), TTL: time.Duration(time.Second * 1000), Last_Used: time.Time{}, Frequency: 1})
	cache_A = append(cache_A, CACHE{Name: "Hit_1", Timestamp: time.Now(), TTL: time.Duration(time.Second * 1000), Last_Used: time.Time{}, Frequency: 0})
	cache_A = append(cache_A, CACHE{Name: "Hit_2", Timestamp: time.Now(), TTL: time.Duration(time.Second * 1000), Last_Used: time.Time{}, Frequency: 1})
	node_A.Cache = cache_A
	cache_Data(node_A, packet_A, "LFU")
	
	//Compare output
	want := 3
	if len(node_A.Cache) != want{       
		t.Errorf("Actual: %v -- Want: %v", len(node_A.Cache), want)
	}
	
	//Compare output
	want_2 := "Miss"
	if node_A.Cache[2].Name != want_2{       
		t.Errorf("Actual: %v -- Want: %v", node_A.Cache[2].Name, want_2)
	}
	
	//Compare output
	want_2 = "Hello"
	if node_A.Cache[2].Data != want_2{       
		t.Errorf("Actual: %v -- Want: %v", node_A.Cache[2].Data, want_2)
	}
	
	//Compare output
	want_3 := time.Duration(10)
	if node_A.Cache[2].TTL != want_3{       
		t.Errorf("Actual: %v -- Want: %v", node_A.Cache[2].TTL, want_3)
	}
}
//----------------------------------------------

//--------------------------------------------------------------------------------------------

//--------------------------------------------------------------------------------------------
// pressure_Mode_Helper(curr_Node *Node, packet Packet, topology []*Node, distance_Formula string, destination int) int
//--------------------------------------------------------------------------------------------

//----------------------------------------------
// Tests to make sure pressure helper works if a neighbor hasnt been seen
//----------------------------------------------
func TestPressure_Mode_Zero(t *testing.T) {
	var topology []*Node
	node_A := &Node{Number: 0, Neighbors: []int{0,1,2}, Lat: 0, Long: 0}
	node_B := &Node{Number: 1, Neighbors: []int{0,1,2}, Lat: 1, Long: 1}
	node_C := &Node{Number: 2, Neighbors: []int{0,1,2}, Lat: 2, Long: 2}
	node_D := &Node{Number: 3, Neighbors: []int{0,1,2}, Lat: 5, Long: 5}
	topology = append(topology, node_A)
	topology = append(topology, node_B)
	topology = append(topology, node_C)
	topology = append(topology, node_D)
	packet_A := Packet{Nodes_Visit: []int{0,2,5}, Num_Visit: []int{1,2,5}, Destination: 3}
	actual := pressure_Mode_Helper(node_A, packet_A, topology, "euclidean", packet_A.Destination)
	
	//Compare output
	want := 1
	if actual != want{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}	
}
//----------------------------------------------

//----------------------------------------------
// Tests to make sure pressure helper works if seen different times
//----------------------------------------------
func TestPressure_Mode_Diff(t *testing.T) {
	var topology []*Node
	node_A := &Node{Number: 0, Neighbors: []int{0,1,2}, Lat: 0, Long: 0}
	node_B := &Node{Number: 1, Neighbors: []int{0,1,2}, Lat: 1, Long: 1}
	node_C := &Node{Number: 2, Neighbors: []int{0,1,2}, Lat: 2, Long: 2}
	node_D := &Node{Number: 3, Neighbors: []int{0,1,2}, Lat: 5, Long: 5}
	topology = append(topology, node_A)
	topology = append(topology, node_B)
	topology = append(topology, node_C)
	topology = append(topology, node_D)
	packet_A := Packet{Nodes_Visit: []int{0,1,2,5}, Num_Visit: []int{2,3,4,1}, Destination: 3}
	actual := pressure_Mode_Helper(node_A, packet_A, topology, "euclidean", packet_A.Destination)
	
	//Compare output
	want := 0 //weve seen '5' less times, but its not a neighbor to node_A
	if actual != want{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
}
//----------------------------------------------

//----------------------------------------------
// Tests to make sure pressure helper works if seen same times
//----------------------------------------------
func TestPressure_Mode_Same(t *testing.T) {
	var topology []*Node
	node_A := &Node{Number: 0, Neighbors: []int{0,1,2}, Lat: 0, Long: 0}
	node_B := &Node{Number: 1, Neighbors: []int{0,1,2}, Lat: 1, Long: 1}
	node_C := &Node{Number: 2, Neighbors: []int{0,1,2}, Lat: 2, Long: 2}
	node_D := &Node{Number: 3, Neighbors: []int{0,1,2}, Lat: 5, Long: 5}
	topology = append(topology, node_A)
	topology = append(topology, node_B)
	topology = append(topology, node_C)
	topology = append(topology, node_D)
	packet_A := Packet{Nodes_Visit: []int{0,1,2}, Num_Visit: []int{2,2,2}, Destination: 3}
	actual := pressure_Mode_Helper(node_A, packet_A, topology, "euclidean", packet_A.Destination)
	
	//Compare output
	want := 2
	if actual != want{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
}
//----------------------------------------------

//--------------------------------------------------------------------------------------------

//--------------------------------------------------------------------------------------------
// interest_Collapse_Helper(curr_Node *Node, name string, packet Packet, filo bool, destination int) bool
//--------------------------------------------------------------------------------------------

//----------------------------------------------
// Tests to make sure you interest collapse if all conditions are met
//----------------------------------------------
func TestInterest_Collapse(t *testing.T) {
	node_A := &Node{Number: 0}
	var PIT_A []PIT
	PIT_A = append(PIT_A, PIT{Incoming: 0, Outgoing: 0, Name: "a"})
	PIT_A = append(PIT_A, PIT{Incoming: 1, Outgoing: 1, Name: "b"})
	PIT_A = append(PIT_A, PIT{Incoming: 2, Outgoing: 2, Name: "c"})
	node_A.Pit = PIT_A
	packet_A := Packet{Name: "a", GP_Mode: 0, Previous: 100}
	actual := interest_Collapse_Helper(node_A, packet_A.Name, packet_A, false, -2)
	
	//Compare output
	want := true
	if actual != want{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
	
	//Compare output
	want_2 := 4
	if len(node_A.Pit) != want_2{       
		t.Errorf("Actual: %v -- Want: %v", len(node_A.Pit), want_2)
	}
	
	//Compare output
	want_2 = 100
	if node_A.Pit[3].Incoming != want_2{       
		t.Errorf("Actual: %v -- Want: %v", node_A.Pit[3].Incoming, want_2)
	}
	
	//Compare output
	want_2 = -2
	if node_A.Pit[3].Outgoing != want_2{       
		t.Errorf("Actual: %v -- Want: %v", node_A.Pit[3].Outgoing, want_2)
	}
	
	//Compare output
	want_3 := "a"
	if node_A.Pit[3].Name != want_3{       
		t.Errorf("Actual: %v -- Want: %v", node_A.Pit[3].Name, want_3)
	}
}
//----------------------------------------------

//----------------------------------------------
// Tests to make sure you dont interest collapse for FILO algorithms
//----------------------------------------------
func TestInterest_Collapse_Fail_FILO(t *testing.T) {
	node_A := &Node{Number: 0}
	var PIT_A []PIT
	PIT_A = append(PIT_A, PIT{Incoming: 0, Outgoing: 0, Name: "a"})
	PIT_A = append(PIT_A, PIT{Incoming: 1, Outgoing: 1, Name: "b"})
	PIT_A = append(PIT_A, PIT{Incoming: 2, Outgoing: 2, Name: "c"})
	node_A.Pit = PIT_A
	packet_A := Packet{Name: "a", GP_Mode: 0, Previous: 100}
	actual := interest_Collapse_Helper(node_A, packet_A.Name, packet_A, true, -2)
	
	//Compare output
	want := false
	if actual != want{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
	
	//Compare output
	want_2 := 3
	if len(node_A.Pit) != want_2{       
		t.Errorf("Actual: %v -- Want: %v", len(node_A.Pit), want_2)
	}
}
//----------------------------------------------

//----------------------------------------------
// Tests to make sure you dont interest collapse if you are in pressure mode
//----------------------------------------------
func TestInterest_Collapse_Fail_Pressure(t *testing.T) {
	node_A := &Node{Number: 0}
	var PIT_A []PIT
	PIT_A = append(PIT_A, PIT{Incoming: 0, Outgoing: 0, Name: "a"})
	PIT_A = append(PIT_A, PIT{Incoming: 1, Outgoing: 1, Name: "b"})
	PIT_A = append(PIT_A, PIT{Incoming: 2, Outgoing: 2, Name: "c"})
	node_A.Pit = PIT_A
	packet_A := Packet{Name: "a", GP_Mode: 1, Previous: 100}
	actual := interest_Collapse_Helper(node_A, packet_A.Name, packet_A, false, -2)
	
	//Compare output
	want := false
	if actual != want{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
	
	//Compare output
	want_2 := 3
	if len(node_A.Pit) != want_2{       
		t.Errorf("Actual: %v -- Want: %v", len(node_A.Pit), want_2)
	}
}
//----------------------------------------------

//----------------------------------------------
// Tests to make sure you dont interest collapse if you dont get a match
//----------------------------------------------
func TestInterest_Collapse_Fail(t *testing.T) {
	node_A := &Node{Number: 0}
	var PIT_A []PIT
	PIT_A = append(PIT_A, PIT{Incoming: 0, Outgoing: 0, Name: "a"})
	PIT_A = append(PIT_A, PIT{Incoming: 1, Outgoing: 1, Name: "b"})
	PIT_A = append(PIT_A, PIT{Incoming: 2, Outgoing: 2, Name: "c"})
	node_A.Pit = PIT_A
	packet_A := Packet{Name: "d", GP_Mode: 0, Previous: 100}
	actual := interest_Collapse_Helper(node_A, packet_A.Name, packet_A, false, -2)
	
	//Compare output
	want := false
	if actual != want{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
	
	//Compare output
	want_2 := 3
	if len(node_A.Pit) != want_2{       
		t.Errorf("Actual: %v -- Want: %v", len(node_A.Pit), want_2)
	}
}
//----------------------------------------------

//----------------------------------------------
// Tests to make sure smart interest collapsing works
//----------------------------------------------
func TestInterest_Collapse_Smart(t *testing.T) {
	node_A := &Node{Number: 0}
	var PIT_A []PIT
	PIT_A = append(PIT_A, PIT{Incoming: 0, Outgoing: 0, Name: "a"})
	PIT_A = append(PIT_A, PIT{Incoming: 700, Outgoing: 900, Name: "b"})
	PIT_A = append(PIT_A, PIT{Incoming: 2, Outgoing: 2, Name: "c"})
	node_A.Pit = PIT_A
	packet_A := Packet{Name: "b", GP_Mode: 0, Previous: 900}
	actual := interest_Collapse_Helper(node_A, packet_A.Name, packet_A, false, -2)
	
	//Compare output
	want := false
	if actual != want{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
	
	//Compare output
	want_2 := 3
	if len(node_A.Pit) != want_2{       
		t.Errorf("Actual: %v -- Want: %v", len(node_A.Pit), want_2)
	}
}
//----------------------------------------------

//--------------------------------------------------------------------------------------------

//--------------------------------------------------------------------------------------------
// return_FILO(curr_Node *Node, packet Packet, topology []*Node) []*Node
//--------------------------------------------------------------------------------------------

//----------------------------------------------
// Tests to make sure you the next node is chosen correctly
// Also make sure all PIT entries are removed correctly
//----------------------------------------------
func TestReturn_FILO(t *testing.T) {
	var topology []*Node
	node_A := &Node{Number: 0, Neighbors: []int{0,1,2}, Lat: 0, Long: 0}
	node_B := &Node{Number: 1, Neighbors: []int{0,1,2}, Lat: 1, Long: 1}
	node_C := &Node{Number: 2, Neighbors: []int{0,1,2}, Lat: 2, Long: 2}
	node_D := &Node{Number: 3, Neighbors: []int{0,1,2}, Lat: 5, Long: 5}
	topology = append(topology, node_A)
	topology = append(topology, node_B)
	topology = append(topology, node_C)
	topology = append(topology, node_D)
	var PIT_A []PIT
	PIT_A = append(PIT_A, PIT{Incoming: 0, Outgoing: 0, Name: "a"})
	PIT_A = append(PIT_A, PIT{Incoming: 2, Outgoing: 2, Name: "b", Start: 100})
	PIT_A = append(PIT_A, PIT{Incoming: 2, Outgoing: 2, Name: "b", Start: 200})
	PIT_A = append(PIT_A, PIT{Incoming: 4, Outgoing: 4, Name: "c"})
	node_A.Pit = PIT_A
	packet_A := Packet{Name: "b", Previous: 2, Path_Traversed: []int{1,2,3}, Index_Middle: 2}
	actual := return_FILO(node_A, packet_A, topology)

	//Compare output
	want := 1
	if len(actual) != want{       
		t.Errorf("Actual: %v -- Want: %v", len(actual), want)
	}
	
	//Compare output
	want_2 := node_C
	if actual[0] != want_2{       
		t.Errorf("Actual: %v -- Want: %v", actual[0], want_2)
	}
	
	//Compare output
	want = 3
	if len(node_A.Pit) != want{       
		t.Errorf("Actual: %v -- Want: %v", len(node_A.Pit), want)
	}
	
	//Compare output
	want = 4
	if len(PIT_A) != want{       
		t.Errorf("Actual: %v -- Want: %v", len(PIT_A), want)
	}
	
	//Compare output
	want_3 := PIT_A[1]
	if node_A.Pit[1] != want_3{       
		t.Errorf("Actual: %v -- Want: %v", node_A.Pit[1], want_3)
	}	
	
	//Compare output
	want_3 = PIT_A[3]
	if node_A.Pit[2] != want_3{       
		t.Errorf("Actual: %v -- Want: %v", node_A.Pit[2], want_3)
	}	
}
//----------------------------------------------

//--------------------------------------------------------------------------------------------

//--------------------------------------------------------------------------------------------
// return_ALL(curr_Node *Node, packet Packet, topology []*Node) ([]*Node, []int, []int, []int)
//--------------------------------------------------------------------------------------------

//----------------------------------------------
// Tests to make sure you the next node is chosen correctly if only 1 choice
// Also make sure all PIT entries are removed correctly
//----------------------------------------------
func TestReturn_ALL_1(t *testing.T) {
	var topology []*Node
	node_A := &Node{Number: 0, Neighbors: []int{0,1,2}, Lat: 0, Long: 0}
	node_B := &Node{Number: 1, Neighbors: []int{0,1,2}, Lat: 1, Long: 1}
	node_C := &Node{Number: 2, Neighbors: []int{0,1,2}, Lat: 2, Long: 2}
	node_D := &Node{Number: 3, Neighbors: []int{0,1,2}, Lat: 5, Long: 5}
	topology = append(topology, node_A)
	topology = append(topology, node_B)
	topology = append(topology, node_C)
	topology = append(topology, node_D)
	var PIT_A []PIT
	PIT_A = append(PIT_A, PIT{Incoming: 0, Outgoing: 0, Name: "a"})
	PIT_A = append(PIT_A, PIT{Incoming: 1, Outgoing: 1, Name: "b_1", Start: 100})
	PIT_A = append(PIT_A, PIT{Incoming: 2, Outgoing: 2, Name: "b", Start: 200, Destination: 300, Collapsed: 400})
	PIT_A = append(PIT_A, PIT{Incoming: 4, Outgoing: 4, Name: "c"})
	node_A.Pit = PIT_A
	packet_A := Packet{Name: "b", Previous: 2, Path_Traversed: []int{1,2,3}, Index_Middle: 2}
	actual, actual_1, actual_2, actual_3 := return_ALL(node_A, packet_A, topology)

	//Compare output
	want := 200
	if actual_1[0] != want || len(actual_1) != 1{       
		t.Errorf("Actual: %v -- Want: %v", actual_1[0], want)
		t.Errorf("Actual: %v -- Want: %v", len(actual_1), 1)
	}
	
	//Compare output
	want = 300
	if actual_2[0] != want || len(actual_2) != 1{       
		t.Errorf("Actual: %v -- Want: %v", actual_2[0], want)
		t.Errorf("Actual: %v -- Want: %v", len(actual_2), 1)
	}
	
	//Compare output
	want = 400
	if actual_3[0] != want || len(actual_3) != 1{       
		t.Errorf("Actual: %v -- Want: %v", actual_3[0], want)
		t.Errorf("Actual: %v -- Want: %v", len(actual_3), 1)
	}

	//Compare output
	want = 1
	if len(actual) != want{       
		t.Errorf("Actual: %v -- Want: %v", len(actual), want)
	}
	
	//Compare output
	want_2 := node_C
	if actual[0] != want_2{       
		t.Errorf("Actual: %v -- Want: %v", actual[0], want_2)
	}
	
	//Compare output
	want = 3
	if len(node_A.Pit) != want{       
		t.Errorf("Actual: %v -- Want: %v", len(node_A.Pit), want)
	}
	
	//Compare output
	want = 4
	if len(PIT_A) != want{       
		t.Errorf("Actual: %v -- Want: %v", len(PIT_A), want)
	}
	
	//Compare output
	want_3 := PIT_A[1]
	if node_A.Pit[1] != want_3{       
		t.Errorf("Actual: %v -- Want: %v", node_A.Pit[1], want_3)
	}	
	
	//Compare output
	want_3 = PIT_A[3]
	if node_A.Pit[2] != want_3{       
		t.Errorf("Actual: %v -- Want: %v", node_A.Pit[1], want_3)
	}	
}
//----------------------------------------------

//----------------------------------------------
// Tests to make sure you the next node is chosen correctly if multiple choices
// Also make sure all PIT entries are removed correctly
//----------------------------------------------
func TestReturn_ALL_Multiple(t *testing.T) {
	var topology []*Node
	node_A := &Node{Number: 0, Neighbors: []int{0,1,2}, Lat: 0, Long: 0}
	node_B := &Node{Number: 1, Neighbors: []int{0,1,2}, Lat: 1, Long: 1}
	node_C := &Node{Number: 2, Neighbors: []int{0,1,2}, Lat: 2, Long: 2}
	node_D := &Node{Number: 3, Neighbors: []int{0,1,2}, Lat: 5, Long: 5}
	topology = append(topology, node_A)
	topology = append(topology, node_B)
	topology = append(topology, node_C)
	topology = append(topology, node_D)
	var PIT_A []PIT
	PIT_A = append(PIT_A, PIT{Incoming: 0, Outgoing: 0, Name: "a"})
	PIT_A = append(PIT_A, PIT{Incoming: 1, Outgoing: 1, Name: "b", Start: 100, Destination: 101, Collapsed: 102})
	PIT_A = append(PIT_A, PIT{Incoming: 2, Outgoing: 2, Name: "b", Start: 200, Destination: 300, Collapsed: 400})
	PIT_A = append(PIT_A, PIT{Incoming: 4, Outgoing: 4, Name: "c"})
	node_A.Pit = PIT_A
	packet_A := Packet{Name: "b", Previous: 2, Path_Traversed: []int{1,2,3}, Index_Middle: 2}
	actual, actual_1, actual_2, actual_3 := return_ALL(node_A, packet_A, topology)

	//Compare output
	if actual_1[0] != 200 || actual_1[1] != 100 || len(actual_1) != 2{       
		t.Errorf("Actual: %v -- Want: %v", actual_1[0], 200)
		t.Errorf("Actual: %v -- Want: %v", actual_1[1], 100)
		t.Errorf("Actual: %v -- Want: %v", len(actual_1), 2)
	}
	
	//Compare output
	if actual_2[0] != 300 || actual_2[1] != 101 || len(actual_1) != 2{       
		t.Errorf("Actual: %v -- Want: %v", actual_2[0], 300)
		t.Errorf("Actual: %v -- Want: %v", actual_2[1], 101)
		t.Errorf("Actual: %v -- Want: %v", len(actual_2), 2)
	}
	
	//Compare output
	if actual_3[0] != 400 || actual_3[1] != 102 || len(actual_1) != 2{       
		t.Errorf("Actual: %v -- Want: %v", actual_3[0], 400)
		t.Errorf("Actual: %v -- Want: %v", actual_3[1], 102)
		t.Errorf("Actual: %v -- Want: %v", len(actual_3), 2)
	}

	//Compare output
	want := 2
	if len(actual) != want{       
		t.Errorf("Actual: %v -- Want: %v", len(actual), want)
	}
	
	//Compare output
	want_2 := node_C
	if actual[0] != want_2{       
		t.Errorf("Actual: %v -- Want: %v", actual[0], want_2)
	}
	
	//Compare output
	want_2 = node_B
	if actual[1] != want_2{       
		t.Errorf("Actual: %v -- Want: %v", actual[1], want_2)
	}
	
	//Compare output
	want = 2
	if len(node_A.Pit) != want{       
		t.Errorf("Actual: %v -- Want: %v", len(node_A.Pit), want)
	}
	
	//Compare output
	want = 4
	if len(PIT_A) != want{       
		t.Errorf("Actual: %v -- Want: %v", len(PIT_A), want)
	}
	
	//Compare output
	want_3 := PIT_A[3]
	if node_A.Pit[1] != want_3{       
		t.Errorf("Actual: %v -- Want: %v", actual, want_3)
	}		
}
//----------------------------------------------

//--------------------------------------------------------------------------------------------
//random_Selection(topology []*Node, percentage float64) []int
//--------------------------------------------------------------------------------------------

//----------------------------------------------
// Tests to make sure random selection works with half
//----------------------------------------------
func TestRandom_Selection(t *testing.T) {
	rand.Seed(1)
	var topology []*Node
	node_A := &Node{Number: 0, Neighbors: []int{0,1,2}, Lat: 0, Long: 0}
	node_B := &Node{Number: 1, Neighbors: []int{0,1,2}, Lat: 1, Long: 1}
	node_C := &Node{Number: 2, Neighbors: []int{0,1,2}, Lat: 2, Long: 2}
	node_D := &Node{Number: 3, Neighbors: []int{0,1,2}, Lat: 5, Long: 5}
	topology = append(topology, node_A)
	topology = append(topology, node_B)
	topology = append(topology, node_C)
	topology = append(topology, node_D)
	nodes := random_Selection(topology, .5)
	
	//Compare output
	want := 2
	if len(nodes) != want{       
		t.Errorf("Actual: %v -- Want: %v", len(nodes), want)
	}
	
	//Compare output
	want = 0
	if nodes[0] != want{       
		t.Errorf("Actual: %v -- Want: %v", nodes[0], want)
	}
	
	//Compare output
	want = 1
	if nodes[1] != want{       
		t.Errorf("Actual: %v -- Want: %v", nodes[1], want)
	}	
}
//----------------------------------------------

//----------------------------------------------
// Tests to make sure random selection works with all
//----------------------------------------------
func TestRandom_Selection_All(t *testing.T) {
	rand.Seed(1)
	var topology []*Node
	node_A := &Node{Number: 0, Neighbors: []int{0,1,2}, Lat: 0, Long: 0}
	node_B := &Node{Number: 1, Neighbors: []int{0,1,2}, Lat: 1, Long: 1}
	node_C := &Node{Number: 2, Neighbors: []int{0,1,2}, Lat: 2, Long: 2}
	node_D := &Node{Number: 3, Neighbors: []int{0,1,2}, Lat: 5, Long: 5}
	topology = append(topology, node_A)
	topology = append(topology, node_B)
	topology = append(topology, node_C)
	topology = append(topology, node_D)
	nodes := random_Selection(topology, 1)
	
	//Compare output
	want := 4
	if len(nodes) != want{       
		t.Errorf("Actual: %v -- Want: %v", len(nodes), want)
	}	
}
//----------------------------------------------

//----------------------------------------------
// Tests to make sure random selection works with none
//----------------------------------------------
func TestRandom_Selection_None(t *testing.T) {
	rand.Seed(1)
	var topology []*Node
	node_A := &Node{Number: 0, Neighbors: []int{0,1,2}, Lat: 0, Long: 0}
	node_B := &Node{Number: 1, Neighbors: []int{0,1,2}, Lat: 1, Long: 1}
	node_C := &Node{Number: 2, Neighbors: []int{0,1,2}, Lat: 2, Long: 2}
	node_D := &Node{Number: 3, Neighbors: []int{0,1,2}, Lat: 5, Long: 5}
	topology = append(topology, node_A)
	topology = append(topology, node_B)
	topology = append(topology, node_C)
	topology = append(topology, node_D)
	nodes := random_Selection(topology, 0)
	
	//Compare output
	want := 0
	if len(nodes) != want{       
		t.Errorf("Actual: %v -- Want: %v", len(nodes), want)
	}	
}
//----------------------------------------------

//--------------------------------------------------------------------------------------------
//flooding_Helper(curr_Node *Node, previous_Node_Num int) []int
//--------------------------------------------------------------------------------------------

//----------------------------------------------
// Tests to make sure flooding helper works
//----------------------------------------------
func TestFlooding_Helper(t *testing.T) {
	node_A := &Node{Number: 0, Neighbors: []int{0,1,2,3,4}, Lat: 0, Long: 0}
	actual := flooding_Helper(node_A, 3)
	
	//Compare output
	want := 3
	if len(actual) != want{       
		t.Errorf("Actual: %v -- Want: %v", len(actual), want)
	}
	
	//Compare output
	want_2 := []int{1,2,4}
	for x := 0; x < len(actual); x++{
		if actual[x] != want_2[x]{       
			t.Errorf("Actual: %v -- Want: %v", actual[x], want_2[x])
		}
	}	
}
//----------------------------------------------

//----------------------------------------------
// Tests to make sure flooding helper works if you are not a neighbor to yourself
//----------------------------------------------
func TestFlooding_Helper_Self(t *testing.T) {
	node_A := &Node{Number: 0, Neighbors: []int{1,2,3,4}, Lat: 0, Long: 0}
	actual := flooding_Helper(node_A, 3)
	
	//Compare output
	want := 3
	if len(actual) != want{       
		t.Errorf("Actual: %v -- Want: %v", len(actual), want)
	}
	
	//Compare output
	want_2 := []int{1,2,4}
	for x := 0; x < len(actual); x++{
		if actual[x] != want_2[x]{       
			t.Errorf("Actual: %v -- Want: %v", actual[x], want_2[x])
		}
	}	
}
//----------------------------------------------

//--------------------------------------------------------------------------------------------
//copy_Packet(packet Packet) Packet
//--------------------------------------------------------------------------------------------

//----------------------------------------------
// Tests to make sure copy packet makes a deepcopy
//----------------------------------------------
func TestCopy_Packet(t *testing.T) {
	packet := Packet{0, 1, 2, "2", 3, "Hello", "World", time.Time{}, time.Duration(10), 4, []int{1,2,3}, []int{4,5,6}, 100, []int{1,2,3}, 1, 2, time.Time{}, 3, 0}
	new_Packet := copy_Packet(packet)
	
	time_Now := time.Now()
	packet.Previous = 100
	packet.Name = "A"
	packet.Start_Time = time_Now
	packet.TTL = time.Duration(100)
	packet.Nodes_Visit = []int{100,200,300}

	//Compare output
	want := 0
	actual := new_Packet.Previous
	if actual != want{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
	
	//Compare output
	want_1 := "Hello"
	actual_1 := new_Packet.Name
	if actual_1 != want_1{       
		t.Errorf("Actual: %v -- Want: %v", actual_1, want_1)
	}
	
	//Compare output
	want_2 := time.Time{}
	actual_2 := new_Packet.Start_Time
	if actual_2 != want_2{       
		t.Errorf("Actual: %v -- Want: %v", actual_2, want_2)
	}
	
	//Compare output
	want_3 := time.Duration(10)
	actual_3 := new_Packet.TTL
	if actual_3 != want_3{       
		t.Errorf("Actual: %v -- Want: %v", actual_3, want_3)
	}
	
	//Compare output
	want_4 := []int{1,2,3}
	actual_4 := new_Packet.Nodes_Visit
	for x := 0; x < len(actual_4); x++{
		if actual_4[x] != want_4[x]{       
			t.Errorf("Actual: %v -- Want: %v", actual_4[x], want_4[x])
		}
	}	
	
	new_Packet.Previous = 1000
	new_Packet.Name = "1"
	new_Packet.Start_Time = time_Now
	new_Packet.TTL = time.Duration(1000)
	new_Packet.Nodes_Visit = []int{1000,2000,3000}
	
	//Compare output
	want = 100
	actual = packet.Previous
	if actual != want{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
	
	//Compare output
	want_1 = "A"
	actual_1 = packet.Name
	if actual_1 != want_1{       
		t.Errorf("Actual: %v -- Want: %v", actual_1, want_1)
	}
	
	//Compare output
	want_2 = time_Now
	actual_2 = packet.Start_Time
	if actual_2 != want_2{       
		t.Errorf("Actual: %v -- Want: %v", actual_2, want_2)
	}
	
	//Compare output
	want_3 = time.Duration(100)
	actual_3 = packet.TTL
	if actual_3 != want_3{       
		t.Errorf("Actual: %v -- Want: %v", actual_3, want_3)
	}
	
	//Compare output
	want_4 = []int{100,200,300}
	actual_4 = packet.Nodes_Visit
	for x := 0; x < len(actual_4); x++{
		if actual_4[x] != want_4[x]{       
			t.Errorf("Actual: %v -- Want: %v", actual_4[x], want_4[x])
		}
	}	
	

}
//----------------------------------------------

//--------------------------------------------------------------------------------------------
//update_Names(names []string, satisfied_Names []string, destinations []string) ([]string, string, bool, []string)
//--------------------------------------------------------------------------------------------

//----------------------------------------------
// Tests to make sure names gets updated correctly
//----------------------------------------------
func TestUpdate_Name(t *testing.T) {
	names := []string{"A", "B", "C", "D"}
	destinations := []string{"0", "1", "2", "3"}
	satisfied_Names := []string{"B", "C"}
	new_Names, new_Name, drop_Packet, new_Destinations := update_Names(names, satisfied_Names, destinations)
	
	//Compare output
	want_4 := 4
	actual_4 := len(names)
	if want_4 != actual_4{       
		t.Errorf("Actual: %v -- Want: %v", actual_4, want_4)
	}
	
	//Compare output
	want := []string{"A", "B", "C", "D"}
	actual := names
	for x := 0; x < len(actual); x++{
		if want[x] != actual[x] {       
			t.Errorf("Actual: %v -- Want: %v", actual[x], want[x])
		}
	}

	//Compare output
	want_4 = 4
	actual_4 = len(names)
	if want_4 != actual_4{       
		t.Errorf("Actual: %v -- Want: %v", actual_4, want_4)
	}
	
	//Compare output
	want = []string{"0", "1", "2", "3"}
	actual = destinations
	for x := 0; x < len(actual); x++{
		if want[x] != actual[x] {       
			t.Errorf("Actual: %v -- Want: %v", actual[x], want[x])
		}
	}
	
	//Compare output
	want_4 = 2
	actual_4 = len(satisfied_Names)
	if want_4 != actual_4{       
		t.Errorf("Actual: %v -- Want: %v", actual_4, want_4)
	}
	
	//Compare output
	want = []string{"B", "C"}
	actual = satisfied_Names
	for x := 0; x < len(actual); x++{
		if want[x] != actual[x] {       
			t.Errorf("Actual: %v -- Want: %v", actual[x], want[x])
		}
	}
	
	//Compare output
	want_4 = 2
	actual_4 = len(new_Names)
	if want_4 != actual_4{       
		t.Errorf("Actual: %v -- Want: %v", actual_4, want_4)
	}
	
	//Compare output
	want = []string{"A", "D"}
	actual = new_Names
	for x := 0; x < len(actual); x++{
		if want[x] != actual[x] {       
			t.Errorf("Actual: %v -- Want: %v", actual[x], want[x])
		}
	}
	
	//Compare output
	want = []string{"0", "3"}
	actual = new_Destinations
	for x := 0; x < len(actual); x++{
		if want[x] != actual[x] {       
			t.Errorf("Actual: %v -- Want: %v", actual[x], want[x])
		}
	}
	
	//Compare output
	want_2 := "A|D"
	actual_2 := new_Name
	if want_2 != actual_2{       
		t.Errorf("Actual: %v -- Want: %v", actual_2, want_2)
	}
	
	//Compare output
	want_3 := false
	actual_3 := drop_Packet
	if want_3 != actual_3{       
		t.Errorf("Actual: %v -- Want: %v", actual_3, want_3)
	}
}
//----------------------------------------------

//----------------------------------------------
// Tests to make sure names gets updated correctly if all names are satisfied
//----------------------------------------------
func TestUpdate_Name_All_Sat(t *testing.T) {
	names := []string{"A", "B", "C", "D"}
	destinations := []string{"0", "1", "2", "3"}
	satisfied_Names := []string{"A", "B", "C", "D"}
	new_Names, new_Name, drop_Packet, new_Destinations := update_Names(names, satisfied_Names, destinations)
	
	//Compare output
	want_4 := 0
	actual_4 := len(new_Names)
	if want_4 != actual_4{       
		t.Errorf("Actual: %v -- Want: %v", actual_4, want_4)
	}
	
	//Compare output
	want := []string{}
	actual := new_Names
	for x := 0; x < len(actual); x++{
		if want[x] != actual[x] {       
			t.Errorf("Actual: %v -- Want: %v", actual[x], want[x])
		}
	}

	//Compare output
	want = []string{}
	actual = new_Destinations
	for x := 0; x < len(actual); x++{
		if want[x] != actual[x] {       
			t.Errorf("Actual: %v -- Want: %v", actual[x], want[x])
		}
	}

	//Compare output
	want_2 := ""
	actual_2 := new_Name
	if want_2 != actual_2{       
		t.Errorf("Actual: %v -- Want: %v", actual_2, want_2)
	}
	
	//Compare output
	want_3 := true
	actual_3 := drop_Packet
	if want_3 != actual_3{       
		t.Errorf("Actual: %v -- Want: %v", actual_3, want_3)
	}
}
//----------------------------------------------

//----------------------------------------------
// Tests to make sure names gets updated correctly if satisfied names is empty
//----------------------------------------------
func TestUpdate_Name_Sat_Empty(t *testing.T) {
	names := []string{"A", "B", "C", "D"}
	destinations := []string{"0", "1", "2", "3"}
	satisfied_Names := []string{}
	new_Names, new_Name, drop_Packet, new_Destinations := update_Names(names, satisfied_Names, destinations)
	
	//Compare output
	want_4 := 4
	actual_4 := len(new_Names)
	if want_4 != actual_4{       
		t.Errorf("Actual: %v -- Want: %v", actual_4, want_4)
	}
	
	//Compare output
	want := []string{"A", "B", "C", "D"}
	actual := new_Names
	for x := 0; x < len(actual); x++{
		if want[x] != actual[x] {       
			t.Errorf("Actual: %v -- Want: %v", actual[x], want[x])
		}
	}
	
	//Compare output
	want = []string{"0", "1", "2", "3"}
	actual = new_Destinations
	for x := 0; x < len(actual); x++{
		if want[x] != actual[x] {       
			t.Errorf("Actual: %v -- Want: %v", actual[x], want[x])
		}
	}
	
	//Compare output
	want_2 := "A|B|C|D"
	actual_2 := new_Name
	if want_2 != actual_2{       
		t.Errorf("Actual: %v -- Want: %v", actual_2, want_2)
	}
	
	//Compare output
	want_3 := false
	actual_3 := drop_Packet
	if want_3 != actual_3{       
		t.Errorf("Actual: %v -- Want: %v", actual_3, want_3)
	}
}
//----------------------------------------------

//----------------------------------------------
// Tests to make sure names gets updated correctly if names is empty
//----------------------------------------------
func TestUpdate_Name_Names_Empty(t *testing.T) {
	names := []string{}
	destinations := []string{}
	satisfied_Names := []string{"A", "B", "C", "D"}
	new_Names, new_Name, drop_Packet, new_Destinations := update_Names(names, satisfied_Names, destinations)
	
	//Compare output
	want_4 := 0
	actual_4 := len(new_Names)
	if want_4 != actual_4{       
		t.Errorf("Actual: %v -- Want: %v", actual_4, want_4)
	}
	
	//Compare output
	want := []string{}
	actual := new_Names
	for x := 0; x < len(actual); x++{
		if want[x] != actual[x] {       
			t.Errorf("Actual: %v -- Want: %v", actual[x], want[x])
		}
	}
	
	//Compare output
	want = []string{}
	actual = new_Destinations
	for x := 0; x < len(actual); x++{
		if want[x] != actual[x] {       
			t.Errorf("Actual: %v -- Want: %v", actual[x], want[x])
		}
	}
	
	//Compare output
	want_2 := ""
	actual_2 := new_Name
	if want_2 != actual_2{       
		t.Errorf("Actual: %v -- Want: %v", actual_2, want_2)
	}
	
	//Compare output
	want_3 := true
	actual_3 := drop_Packet
	if want_3 != actual_3{       
		t.Errorf("Actual: %v -- Want: %v", actual_3, want_3)
	}
}
//----------------------------------------------

//--------------------------------------------------------------------------------------------
//compare_Packet(packet_A Packet, packet_B Packet, exact bool) bool
//--------------------------------------------------------------------------------------------

//----------------------------------------------
// Tests to make sure all value matching returns true if exact = true
//----------------------------------------------
func TestCompare_Packet_All_True_Exact(t *testing.T) {
	time_Now := time.Now()
	packet_A := Packet{0, 1, 2, "hello", 3, "a", "b", time_Now, 4, 5, []int{1,2,3}, []int{4,5,6}, 6.0, []int{7,8,9}, 7, 8, time_Now, 9, 10}
	packet_B := Packet{0, 1, 2, "hello", 3, "a", "b", time_Now, 4, 5, []int{1,2,3}, []int{4,5,6}, 6.0, []int{7,8,9}, 7, 8, time_Now, 9, 10}
	actual := compare_Packet(packet_A, packet_B, true)
	
	//Compare output
	want := true
	if want != actual{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
}
//----------------------------------------------

//----------------------------------------------
// Tests to make sure all value matching returns true if exact = false
//----------------------------------------------
func TestCompare_Packet_All_True(t *testing.T) {
	time_Now := time.Now()
	packet_A := Packet{0, 1, 2, "hello", 3, "a", "b", time_Now, 4, 5, []int{1,2,3}, []int{4,5,6}, 6.0, []int{7,8,9}, 7, 8, time_Now, 9, 10}
	packet_B := Packet{0, 1, 100, "hello", 3, "cow", "b", time_Now, 4, 5, []int{1,2,3}, []int{4,5,6}, 6.0, []int{7,8,9}, 7, 8, time_Now, 9, 10}
	actual := compare_Packet(packet_A, packet_B, false)
	
	//Compare output
	want := true
	if want != actual{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
}
//----------------------------------------------

//----------------------------------------------
// Tests to make sure no specific value matching returns false
//----------------------------------------------
func TestCompare_Packet_All_False(t *testing.T) {

	//-----Previous-----
	time_Now := time.Now()
	packet_A := Packet{0, 1, 2, "hello", 3, "a", "b", time_Now, 4, 5, []int{1,2,3}, []int{4,5,6}, 6.0, []int{7,8,9}, 7, 8, time_Now, 9, 10}
	packet_B := Packet{100, 1, 2, "hello", 3, "a", "b", time_Now, 4, 5, []int{1,2,3}, []int{4,5,6}, 6.0, []int{7,8,9}, 7, 8, time_Now, 9, 10}
	actual := compare_Packet(packet_A, packet_B, false)
	
	//Compare output
	want := false
	if want != actual{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
	
	//-----Start-----
	time_Now = time.Now()
	packet_A = Packet{0, 1, 2, "hello", 3, "a", "b", time_Now, 4, 5, []int{1,2,3}, []int{4,5,6}, 6.0, []int{7,8,9}, 7, 8, time_Now, 9, 10}
	packet_B = Packet{0, 101, 2, "hello", 3, "a", "b", time_Now, 4, 5, []int{1,2,3}, []int{4,5,6}, 6.0, []int{7,8,9}, 7, 8, time_Now, 9, 10}
	actual = compare_Packet(packet_A, packet_B, false)
	
	//Compare output
	want = false
	if want != actual{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
	
	//-----Destination-----
	time_Now = time.Now()
	packet_A = Packet{0, 1, 2, "hello", 3, "a", "b", time_Now, 4, 5, []int{1,2,3}, []int{4,5,6}, 6.0, []int{7,8,9}, 7, 8, time_Now, 9, 10}
	packet_B = Packet{0, 1, 102, "hello", 3, "a", "b", time_Now, 4, 5, []int{1,2,3}, []int{4,5,6}, 6.0, []int{7,8,9}, 7, 8, time_Now, 9, 10}
	actual = compare_Packet(packet_A, packet_B, false)
	
	//Compare output
	want = true
	if want != actual{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
	
	//-----Destination-----
	time_Now = time.Now()
	packet_A = Packet{0, 1, 2, "hello", 3, "a", "b", time_Now, 4, 5, []int{1,2,3}, []int{4,5,6}, 6.0, []int{7,8,9}, 7, 8, time_Now, 9, 10}
	packet_B = Packet{0, 1, 102, "hello", 3, "a", "b", time_Now, 4, 5, []int{1,2,3}, []int{4,5,6}, 6.0, []int{7,8,9}, 7, 8, time_Now, 9, 10}
	actual = compare_Packet(packet_A, packet_B, true)
	
	//Compare output
	want = false
	if want != actual{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}	
	
	//-----Destination_Str-----
	time_Now = time.Now()
	packet_A = Packet{0, 1, 2, "hello", 3, "a", "b", time_Now, 4, 5, []int{1,2,3}, []int{4,5,6}, 6.0, []int{7,8,9}, 7, 8, time_Now, 9, 10}
	packet_B = Packet{0, 1, 2, "goodbye", 3, "a", "b", time_Now, 4, 5, []int{1,2,3}, []int{4,5,6}, 6.0, []int{7,8,9}, 7, 8, time_Now, 9, 10}
	actual = compare_Packet(packet_A, packet_B, false)
	
	//Compare output
	want = false
	if want != actual{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
	
	//-----Interest_Data-----
	time_Now = time.Now()
	packet_A = Packet{0, 1, 2, "hello", 3, "a", "b", time_Now, 4, 5, []int{1,2,3}, []int{4,5,6}, 6.0, []int{7,8,9}, 7, 8, time_Now, 9, 10}
	packet_B = Packet{0, 1, 2, "hello", 103, "a", "b", time_Now, 4, 5, []int{1,2,3}, []int{4,5,6}, 6.0, []int{7,8,9}, 7, 8, time_Now, 9, 10}
	actual = compare_Packet(packet_A, packet_B, false)
	
	//Compare output
	want = false
	if want != actual{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
	
	//-----Name-----
	time_Now = time.Now()
	packet_A = Packet{0, 1, 2, "hello", 3, "a", "b", time_Now, 4, 5, []int{1,2,3}, []int{4,5,6}, 6.0, []int{7,8,9}, 7, 8, time_Now, 9, 10}
	packet_B = Packet{0, 1, 2, "hello", 3, "z", "b", time_Now, 4, 5, []int{1,2,3}, []int{4,5,6}, 6.0, []int{7,8,9}, 7, 8, time_Now, 9, 10}
	actual = compare_Packet(packet_A, packet_B, false)
	
	//Compare output
	want = true
	if want != actual{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
	
	//-----Name-----
	time_Now = time.Now()
	packet_A = Packet{0, 1, 2, "hello", 3, "a", "b", time_Now, 4, 5, []int{1,2,3}, []int{4,5,6}, 6.0, []int{7,8,9}, 7, 8, time_Now, 9, 10}
	packet_B = Packet{0, 1, 2, "hello", 3, "z", "b", time_Now, 4, 5, []int{1,2,3}, []int{4,5,6}, 6.0, []int{7,8,9}, 7, 8, time_Now, 9, 10}
	actual = compare_Packet(packet_A, packet_B, true)
	
	//Compare output
	want = false
	if want != actual{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
	
	//-----Payload-----
	time_Now = time.Now()
	packet_A = Packet{0, 1, 2, "hello", 3, "a", "b", time_Now, 4, 5, []int{1,2,3}, []int{4,5,6}, 6.0, []int{7,8,9}, 7, 8, time_Now, 9, 10}
	packet_B = Packet{0, 1, 2, "hello", 3, "a", "y", time_Now, 4, 5, []int{1,2,3}, []int{4,5,6}, 6.0, []int{7,8,9}, 7, 8, time_Now, 9, 10}
	actual = compare_Packet(packet_A, packet_B, false)
	
	//Compare output
	want = false
	if want != actual{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
	
	//-----Start_Time-----
	time_Now = time.Now()
	packet_A = Packet{0, 1, 2, "hello", 3, "a", "b", time_Now, 4, 5, []int{1,2,3}, []int{4,5,6}, 6.0, []int{7,8,9}, 7, 8, time_Now, 9, 10}
	packet_B = Packet{0, 1, 2, "hello", 3, "a", "b", time.Now(), 4, 5, []int{1,2,3}, []int{4,5,6}, 6.0, []int{7,8,9}, 7, 8, time_Now, 9, 10}
	actual = compare_Packet(packet_A, packet_B, false)
	
	//Compare output
	want = false
	if want != actual{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
	
	//-----TTL-----
	time_Now = time.Now()
	packet_A = Packet{0, 1, 2, "hello", 3, "a", "b", time_Now, 4, 5, []int{1,2,3}, []int{4,5,6}, 6.0, []int{7,8,9}, 7, 8, time_Now, 9, 10}
	packet_B = Packet{0, 1, 2, "hello", 3, "a", "b", time_Now, 104, 5, []int{1,2,3}, []int{4,5,6}, 6.0, []int{7,8,9}, 7, 8, time_Now, 9, 10}
	actual = compare_Packet(packet_A, packet_B, false)
	
	//Compare output
	want = false
	if want != actual{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
	
	//-----GP_Mode-----
	time_Now = time.Now()
	packet_A = Packet{0, 1, 2, "hello", 3, "a", "b", time_Now, 4, 5, []int{1,2,3}, []int{4,5,6}, 6.0, []int{7,8,9}, 7, 8, time_Now, 9, 10}
	packet_B = Packet{0, 1, 2, "hello", 3, "a", "b", time_Now, 4, 105, []int{1,2,3}, []int{4,5,6}, 6.0, []int{7,8,9}, 7, 8, time_Now, 9, 10}
	actual = compare_Packet(packet_A, packet_B, false)
	
	//Compare output
	want = false
	if want != actual{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}

	//-----Nodes_Visited-----
	time_Now = time.Now()
	packet_A = Packet{0, 1, 2, "hello", 3, "a", "b", time_Now, 4, 5, []int{1,2,3}, []int{4,5,6}, 6.0, []int{7,8,9}, 7, 8, time_Now, 9, 10}
	packet_B = Packet{0, 1, 2, "hello", 3, "a", "b", time_Now, 4, 5, []int{101,102,103}, []int{4,5,6}, 6.0, []int{7,8,9}, 7, 8, time_Now, 9, 10}
	actual = compare_Packet(packet_A, packet_B, false)
	
	//Compare output
	want = false
	if want != actual{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}

	//-----Nodes_Visited-----
	time_Now = time.Now()
	packet_A = Packet{0, 1, 2, "hello", 3, "a", "b", time_Now, 4, 5, []int{1,2,3}, []int{4,5,6}, 6.0, []int{7,8,9}, 7, 8, time_Now, 9, 10}
	packet_B = Packet{0, 1, 2, "hello", 3, "a", "b", time_Now, 4, 5, []int{103}, []int{4,5,6}, 6.0, []int{7,8,9}, 7, 8, time_Now, 9, 10}
	actual = compare_Packet(packet_A, packet_B, false)
	
	//Compare output
	want = false
	if want != actual{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
	
	//-----Num_Visited-----
	time_Now = time.Now()
	packet_A = Packet{0, 1, 2, "hello", 3, "a", "b", time_Now, 4, 5, []int{1,2,3}, []int{4,5,6}, 6.0, []int{7,8,9}, 7, 8, time_Now, 9, 10}
	packet_B = Packet{0, 1, 2, "hello", 3, "a", "b", time_Now, 4, 5, []int{1,2,3}, []int{104,105,106}, 6.0, []int{7,8,9}, 7, 8, time_Now, 9, 10}
	actual = compare_Packet(packet_A, packet_B, false)
	
	//Compare output
	want = false
	if want != actual{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}

	//-----Num_Visited-----
	time_Now = time.Now()
	packet_A = Packet{0, 1, 2, "hello", 3, "a", "b", time_Now, 4, 5, []int{1,2,3}, []int{4,5,6}, 6.0, []int{7,8,9}, 7, 8, time_Now, 9, 10}
	packet_B = Packet{0, 1, 2, "hello", 3, "a", "b", time_Now, 4, 5, []int{1,2,3}, []int{106}, 6.0, []int{7,8,9}, 7, 8, time_Now, 9, 10}
	actual = compare_Packet(packet_A, packet_B, false)
	
	//Compare output
	want = false
	if want != actual{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
	
	//-----GP_Distance-----
	time_Now = time.Now()
	packet_A = Packet{0, 1, 2, "hello", 3, "a", "b", time_Now, 4, 5, []int{1,2,3}, []int{4,5,6}, 6.0, []int{7,8,9}, 7, 8, time_Now, 9, 10}
	packet_B = Packet{0, 1, 2, "hello", 3, "a", "b", time_Now, 4, 5, []int{1,2,3}, []int{4,5,6}, 106.0, []int{7,8,9}, 7, 8, time_Now, 9, 10}
	actual = compare_Packet(packet_A, packet_B, false)
	
	//Compare output
	want = false
	if want != actual{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
	
	//-----Path_Traversed-----
	time_Now = time.Now()
	packet_A = Packet{0, 1, 2, "hello", 3, "a", "b", time_Now, 4, 5, []int{1,2,3}, []int{4,5,6}, 6.0, []int{7,8,9}, 7, 8, time_Now, 9, 10}
	packet_B = Packet{0, 1, 2, "hello", 3, "a", "b", time_Now, 4, 5, []int{1,2,3}, []int{4,5,6}, 6.0, []int{107,108,109}, 7, 8, time_Now, 9, 10}
	actual = compare_Packet(packet_A, packet_B, false)
	
	//Compare output
	want = false
	if want != actual{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
	
	//-----Path_Traversed-----
	time_Now = time.Now()
	packet_A = Packet{0, 1, 2, "hello", 3, "a", "b", time_Now, 4, 5, []int{1,2,3}, []int{4,5,6}, 6.0, []int{7,8,9}, 7, 8, time_Now, 9, 10}
	packet_B = Packet{0, 1, 2, "hello", 3, "a", "b", time_Now, 4, 5, []int{1,2,3}, []int{4,5,6}, 6.0, []int{109}, 7, 8, time_Now, 9, 10}
	actual = compare_Packet(packet_A, packet_B, false)
	
	//Compare output
	want = false
	if want != actual{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}	
	
	//-----Pressure_Gauge-----
	time_Now = time.Now()
	packet_A = Packet{0, 1, 2, "hello", 3, "a", "b", time_Now, 4, 5, []int{1,2,3}, []int{4,5,6}, 6.0, []int{7,8,9}, 7, 8, time_Now, 9, 10}
	packet_B = Packet{0, 1, 2, "hello", 3, "a", "b", time_Now, 4, 5, []int{1,2,3}, []int{4,5,6}, 6.0, []int{7,8,9}, 107, 8, time_Now, 9, 10}
	actual = compare_Packet(packet_A, packet_B, false)
	
	//Compare output
	want = false
	if want != actual{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}

	//-----Index_Middle-----
	time_Now = time.Now()
	packet_A = Packet{0, 1, 2, "hello", 3, "a", "b", time_Now, 4, 5, []int{1,2,3}, []int{4,5,6}, 6.0, []int{7,8,9}, 7, 8, time_Now, 9, 10}
	packet_B = Packet{0, 1, 2, "hello", 3, "a", "b", time_Now, 4, 5, []int{1,2,3}, []int{4,5,6}, 6.0, []int{7,8,9}, 7, 108, time_Now, 9, 10}
	actual = compare_Packet(packet_A, packet_B, false)
	
	//Compare output
	want = false
	if want != actual{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
	
	//-----End_Time-----
	time_Now = time.Now()
	packet_A = Packet{0, 1, 2, "hello", 3, "a", "b", time_Now, 4, 5, []int{1,2,3}, []int{4,5,6}, 6.0, []int{7,8,9}, 7, 8, time_Now, 9, 10}
	packet_B = Packet{0, 1, 2, "hello", 3, "a", "b", time_Now, 4, 5, []int{1,2,3}, []int{4,5,6}, 6.0, []int{7,8,9}, 7, 8, time.Now(), 9, 10}
	actual = compare_Packet(packet_A, packet_B, false)
	
	//Compare output
	want = false
	if want != actual{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
	
	//-----Collapsed-----
	time_Now = time.Now()
	packet_A = Packet{0, 1, 2, "hello", 3, "a", "b", time_Now, 4, 5, []int{1,2,3}, []int{4,5,6}, 6.0, []int{7,8,9}, 7, 8, time_Now, 9, 10}
	packet_B = Packet{0, 1, 2, "hello", 3, "a", "b", time_Now, 4, 5, []int{1,2,3}, []int{4,5,6}, 6.0, []int{7,8,9}, 7, 8, time_Now, 109, 10}
	actual = compare_Packet(packet_A, packet_B, false)
	
	//Compare output
	want = false
	if want != actual{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
	
	//-----Cache_Hit-----
	time_Now = time.Now()
	packet_A = Packet{0, 1, 2, "hello", 3, "a", "b", time_Now, 4, 5, []int{1,2,3}, []int{4,5,6}, 6.0, []int{7,8,9}, 7, 8, time_Now, 9, 10}
	packet_B = Packet{0, 1, 2, "hello", 3, "a", "b", time_Now, 4, 5, []int{1,2,3}, []int{4,5,6}, 6.0, []int{7,8,9}, 7, 8, time_Now, 9, 110}
	actual = compare_Packet(packet_A, packet_B, false)
	
	//Compare output
	want = false
	if want != actual{       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
}
//----------------------------------------------

//--------------------------------------------------------------------------------------------
//combine_Packet(temp_Interest_Node []*Node, temp_Interest_Packet []Packet) ([]*Node, []Packet)
//--------------------------------------------------------------------------------------------

//----------------------------------------------
// Tests to make sure names gets combined correctly
//----------------------------------------------
func TestCombined_Name(t *testing.T) {

	//A B A
	node_A := &Node{Number: 0} //create a struct for node
	node_B := &Node{Number: 1} //create a struct for node
	var temp_Interest_Node []*Node
	temp_Interest_Node = append(temp_Interest_Node, node_A)
	temp_Interest_Node = append(temp_Interest_Node, node_B)
	temp_Interest_Node = append(temp_Interest_Node, node_A)
	temp_Interest_Node = append(temp_Interest_Node, node_A)
	
	//4 unique packets
	packet_A := Packet{GP_Mode: 0, Name: "a", Destination: 0}
	packet_B := Packet{GP_Mode: 0, Name: "b", Destination: 1}
	packet_C := Packet{GP_Mode: 1, Name: "c", Destination: 2}
	packet_D := Packet{GP_Mode: 0, Name: "d", Destination: 3}
	var temp_Interest_Packet []Packet
	temp_Interest_Packet = append(temp_Interest_Packet, packet_A)
	temp_Interest_Packet = append(temp_Interest_Packet, packet_B)
	temp_Interest_Packet = append(temp_Interest_Packet, packet_C)
	temp_Interest_Packet = append(temp_Interest_Packet, packet_D)
	
	next_Nodes, next_Packets := combine_Packet(temp_Interest_Node, temp_Interest_Packet)
	
	//Compare output
	var want []*Node
	want = append(want, node_A)
	want = append(want, node_B)
	want = append(want, node_A)
	actual := next_Nodes
	for x := 0; x < len(actual); x++{
		if want[x].Number != actual[x].Number {       
			t.Errorf("Actual: %v -- Want: %v", actual[x], want[x])
		}
	}
	
	//Compare output
	var want_2 []Packet
	packet_A_New := Packet{Name: "a|d", Destination_Str: "0|3", GP_Mode: 0}
	packet_B_New := Packet{Name: "b", Destination_Str: "1", GP_Mode: 0}
	packet_C_New := Packet{Name: "c", Destination_Str: "2", GP_Mode: 1}
	want_2 = append(want_2, packet_A_New)
	want_2 = append(want_2, packet_B_New)
	want_2 = append(want_2, packet_C_New)
	actual_2 := next_Packets
	for x := 0; x < len(actual_2); x++{
		if want_2[x].Name != actual_2[x].Name {       
			t.Errorf("Actual: %v -- Want: %v", actual_2[x].Name, want_2[x].Name)
		}
		if want_2[x].Destination_Str != actual_2[x].Destination_Str {       
			t.Errorf("Actual: %v -- Want: %v", actual_2[x].Destination_Str, want_2[x].Destination_Str)
		}
		if want_2[x].GP_Mode != actual_2[x].GP_Mode {       
			t.Errorf("Actual: %v -- Want: %v", actual_2[x].GP_Mode, want_2[x].GP_Mode)
		}
	}
}

//----------------------------------------------
// Tests to make sure names gets combined correctly if slices empty
//----------------------------------------------
func TestCombined_Name_Empty(t *testing.T) {

	var temp_Interest_Node []*Node
	var temp_Interest_Packet []Packet
	
	next_Nodes, next_Packets := combine_Packet(temp_Interest_Node, temp_Interest_Packet)
	
	//Compare output
	want := 0
	actual := len(next_Nodes)
	if want != actual {       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
		
	//Compare output
	actual = len(next_Packets)
	if want != actual {       
		t.Errorf("Actual: %v -- Want: %v", actual, want)
	}
}
//----------------------------------------------

//----------------------------------------------
// Tests to make sure mismatching lengths exits
//----------------------------------------------
func TestCombined_Name_Mismatch(t *testing.T) {

	//Run code via cmd
	if os.Getenv("FLAG") == "1" {
		//A B A
		node_A := &Node{Number: 0} //create a struct for node
		node_B := &Node{Number: 1} //create a struct for node
		var temp_Interest_Node []*Node
		temp_Interest_Node = append(temp_Interest_Node, node_A)
		temp_Interest_Node = append(temp_Interest_Node, node_B)
		temp_Interest_Node = append(temp_Interest_Node, node_A)
		
		//4 unique packets
		packet_A := Packet{GP_Mode: 0, Name: "a", Destination: 0}
		packet_B := Packet{GP_Mode: 0, Name: "b", Destination: 1}
		packet_C := Packet{GP_Mode: 1, Name: "c", Destination: 2}
		packet_D := Packet{GP_Mode: 0, Name: "d", Destination: 3}
		var temp_Interest_Packet []Packet
		temp_Interest_Packet = append(temp_Interest_Packet, packet_A)
		temp_Interest_Packet = append(temp_Interest_Packet, packet_B)
		temp_Interest_Packet = append(temp_Interest_Packet, packet_C)
		temp_Interest_Packet = append(temp_Interest_Packet, packet_D)
		_, _ = combine_Packet(temp_Interest_Node, temp_Interest_Packet)
	}
	
	//Execute the code via exec to check error code
	cmd := exec.Command(os.Args[0], "-test.run=TestCombined_Name_Mismatch")
	cmd.Env = append(os.Environ(), "FLAG=1")
	output, err := cmd.Output() //grab error code and output
	
	//Compare error codes
	actual := err.Error()
	want := "exit status 1"
	if actual != want{       
		t.Errorf("Actual: %s -- Want: %s", actual, want)
	}
	
	//Compare output
	actual = string(output)
	want = "Error! Mismatching lengths for packets and next nodes! Exiting.\n"
	if actual != want{       
		t.Errorf("Actual: %s -- Want: %s", actual, want)
	}
}
//----------------------------------------------
