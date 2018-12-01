package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type parameters struct {
	alg string
	lowerCYL int
	upperCYL int
	initCYL int
}

type request struct {
	val int
	read bool
}

type err struct {
	val int
	err bool
}


const maxRequestSize int = 20

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func checkAbortConditions(params parameters) bool {

	if params.upperCYL < params.lowerCYL {
		fmt.Printf("ABORT(13):upper (%d) < lower (%d)\n", params.upperCYL, params.lowerCYL)
		return true
	}

	if params.initCYL < params.lowerCYL {
		fmt.Printf("ABORT(12):initial (%d) < lower (%d)\n", params.initCYL, params.lowerCYL)
		return true
	}

	if params.initCYL > params.upperCYL {
		fmt.Printf("ABORT(11):initial (%d) > upper (%d)\n", params.initCYL, params.upperCYL)
		return true
	}


	return false

}

func parseFile(filePath string) (parameters, []request, []int){

	file, _ := os.Open(filePath)

	var params parameters

	//once we return from parseFile() we can close the file we opened
	defer file.Close()

	scanner := bufio.NewScanner(file)


	//Scan the first line of text
	scanner.Scan()

	//Split the line into strings separated by white space.
	//this will be the algorithm we need to use
	params.alg = strings.Fields(scanner.Text())[1]

	//Scan the second line
	scanner.Scan()

	//convert the string we just pulled and convert it to an integer.
	//This will be the lower position
	params.lowerCYL,_ = strconv.Atoi(strings.Fields(scanner.Text())[1])

	//Scan the third line
	scanner.Scan()

	//convert the string we just pulled and convert it to an integer.
	//This will be the upper position
	params.upperCYL,_ = strconv.Atoi(strings.Fields(scanner.Text())[1])


	//Scan the forth line
	scanner.Scan()

	//convert the string we just pulled and convert it to an integer.
	//This will be the initial position
	params.initCYL,_ = strconv.Atoi(strings.Fields(scanner.Text())[1])



	//Create a new array to store all of the processes
	requests := make([]request, 0)
	requestsPrim := make([]int,0)

	//since we already know how many processes there are
	//we can loop through each line and extract the process
	//name, arrival, and burst time.
	for i:=0; i<maxRequestSize; i++ {
		scanner.Scan()
		s := strings.Fields(scanner.Text())

		if strings.Compare(s[0],"end") == 0 {
			break
		}

		val,_ := strconv.Atoi(s[1])

		requests = append(requests,request{val,false})
		requestsPrim = append(requestsPrim,val)
	}


	return params,requests,requestsPrim
}


//First come first serve
func fcfs(params parameters, requests []int) ([]err, int) {

	var traversals = 0
	var currentPos = params.initCYL
	errors := make([]err,0)

	//Here we processes requests in order unless we hit an error request
	for i:=0; i<len(requests); i++ {

		if requests[i] > params.upperCYL || requests[i] < params.lowerCYL {
			errors = append(errors,err{requests[i],true})
			continue
		}

		errors = append(errors,err{requests[i],false})
		traversals += Abs(currentPos - requests[i])

		currentPos = requests[i]
	}


	return errors, traversals
}


//Find the index of the position with the shortest seek time from the
//current position. This also avoids positions that have already
//been read
func getShortestIndex(requests []request, position int, max int) int {

	if len(requests) < 1 {
		return -1
	}

	var diff = max+1
	var index = -1

	for i:=0; i<len(requests); i++ {
		if requests[i].read == true {
			continue
		}

		if Abs(requests[i].val-position) < diff {
			diff = Abs(requests[i].val-position)
			index = i
		}
	}

	return index
}





//Shortest seek time first
func sstf(params parameters, requests []request) ([]err, int) {

	var traversals = 0
	var currentPos = params.initCYL
	errors := make([]err,0)

	var nextIndex = getShortestIndex(requests,currentPos,params.upperCYL)


	for i:=0; i<len(requests); i++ {

		if requests[nextIndex].val > params.upperCYL || requests[nextIndex].val < params.lowerCYL {
			errors = append(errors,err{requests[nextIndex].val,true})
			requests[nextIndex].read = true
			continue
		}


		traversals += Abs(currentPos-requests[nextIndex].val)

		currentPos = requests[nextIndex].val
		requests[nextIndex].read = true

		errors = append(errors,err{requests[nextIndex].val,false})

		nextIndex = getShortestIndex(requests,currentPos,params.upperCYL)

	}

	return errors, traversals
}






func scan(requests []int, params parameters) ([]err, int) {

	var traversals = 0
	var currentPos = params.initCYL
	var startIndex = len(requests)-1
	errors := make([]err,0)

	//Sort all the requests smallest to largest
	sort.Ints(requests)

	//Find the index of the nearest request to
	//the starting position
	for i:=0; i<len(requests); i++ {
		if requests[i] > currentPos {
			startIndex = i
			break
		}
	}


	//service all the requests in an ascending order
	for i:=startIndex; i<len(requests); i++ {

		if requests[i] > params.upperCYL || requests[i] < params.lowerCYL {
			errors = append(errors,err{requests[i],true})
			continue
		}


		errors = append(errors,err{requests[i],false})
		traversals += Abs(currentPos - requests[i])

		currentPos = requests[i]
	}

	//If out starting position was the first request we are done
	//and can return
	if startIndex < 1 {
		return errors, traversals
	}

	//scan to the upper cylinder limit
	traversals += Abs(params.upperCYL-currentPos)

	//set our head position to the index right before our
	//starting index to simulate scanning down past previous
	//requests
	currentPos = requests[startIndex-1]

	//compute traversal time for moving from upper cylinder to
	//the new index
	traversals += Abs(params.upperCYL-currentPos)


	//Now we service the remaining requests in descending order
	for i:=startIndex-1; i>=0; i-- {

		if requests[i] > params.upperCYL || requests[i] < params.lowerCYL {
			errors = append(errors,err{requests[i],true})
			continue
		}


		errors = append(errors,err{requests[i],false})
		traversals += Abs(currentPos - requests[i])

		currentPos = requests[i]
	}

	return errors, traversals
}





func cscan(requests []int, params parameters) ([]err, int) {
	var traversals = 0
	var currentPos = params.initCYL
	var startIndex = len(requests)-1
	errors := make([]err,0)

	//Sort all the requests smallest to largest
	sort.Ints(requests)

	//Find the index of the nearest request to
	//the starting position
	for i:=0; i<len(requests); i++ {
		if requests[i] > currentPos {
			startIndex = i
			break
		}
	}


	//we are at the top or bottom already
	if startIndex == len(requests)-1  || startIndex == 0{

		//If we are at the top we need to traverse to bottom
		//so this accounts for that seek time
		if startIndex == len(requests)-1 {
			traversals += Abs(params.upperCYL-params.lowerCYL)
		}

		//service all the requests in an ascending order from
		//the bottom
		for i:=0; i<len(requests); i++ {

			if requests[i] > params.upperCYL || requests[i] < params.lowerCYL {
				errors = append(errors,err{requests[i],true})
				continue
			}

			errors = append(errors,err{requests[i],false})
			traversals += Abs(currentPos - requests[i])

			currentPos = requests[i]
		}


		return errors, traversals
	}



	//service all the requests in an ascending order from
	//out starting index
	for i:=startIndex; i<len(requests); i++ {

		if requests[i] > params.upperCYL || requests[i] < params.lowerCYL {
			errors = append(errors,err{requests[i],true})
			continue
		}


		errors = append(errors,err{requests[i],false})
		traversals += Abs(currentPos - requests[i])

		currentPos = requests[i]
	}


	//scan to the upper cylinder limit
	traversals += Abs(params.upperCYL-currentPos)

	//set our head position to the first index
	currentPos = requests[0]

	//compute traversal time for moving from upper cylinder to
	//the bottom
	traversals += Abs(params.upperCYL)

	//compute traversal time for moving from the bottom to
	//the new index
	traversals += Abs(currentPos)


	//Now we service the remaining requests in ascending order
	for i:=0; i<startIndex; i++ {

		if requests[i] > params.upperCYL || requests[i] < params.lowerCYL {
			errors = append(errors,err{requests[i],true})
			continue
		}


		errors = append(errors,err{requests[i],false})
		traversals += Abs(currentPos - requests[i])

		currentPos = requests[i]
	}

	return errors, traversals

}

func look(requests []int, params parameters) ([]err, int) {
	var traversals = 0
	var currentPos = params.initCYL
	var startIndex = len(requests)-1
	errors := make([]err,0)

	//Sort all the requests smallest to largest
	sort.Ints(requests)

	//Find the index of the nearest request to
	//the starting position
	for i:=0; i<len(requests); i++ {
		if requests[i] > currentPos {
			startIndex = i
			break
		}
	}


	//service all the requests in an ascending order
	for i:=startIndex; i<len(requests); i++ {

		if requests[i] > params.upperCYL || requests[i] < params.lowerCYL {
			errors = append(errors,err{requests[i],true})
			continue
		}


		errors = append(errors,err{requests[i],false})
		traversals += Abs(currentPos - requests[i])

		currentPos = requests[i]
	}

	//If out starting position was the first request we are done
	//and can return
	if startIndex < 1 {
		return errors, traversals
	}


	//compute traversal time for moving from current postiion to
	//the new index
	traversals += Abs(currentPos-requests[startIndex-1])

	//set our head position to the index right before our
	//starting index to simulate scanning down past previous
	//requests
	currentPos = requests[startIndex-1]



	//Now we service the remaining requests in descending order
	for i:=startIndex-1; i>=0; i-- {

		if requests[i] > params.upperCYL || requests[i] < params.lowerCYL {
			errors = append(errors,err{requests[i],true})
			continue
		}


		errors = append(errors,err{requests[i],false})
		traversals += Abs(currentPos - requests[i])

		currentPos = requests[i]
	}

	return errors, traversals
}

func clook(requests []int, params parameters) ([]err, int) {
	var traversals = 0
	var currentPos = params.initCYL
	var startIndex = len(requests)-1
	errors := make([]err,0)

	//Sort all the requests smallest to largest
	sort.Ints(requests)

	//Find the index of the nearest request to
	//the starting position
	for i:=0; i<len(requests); i++ {
		if requests[i] > currentPos {
			startIndex = i
			break
		}
	}


	//we are at the top or bottom already
	if startIndex == len(requests)-1  || startIndex == 0{

		//If we are at the top we need to traverse to bottom
		//so this accounts for that seek time
		if startIndex == len(requests)-1 {
			traversals += Abs(requests[len(requests)-1]-requests[0])
		}

		//service all the requests in an ascending order from
		//the bottom
		for i:=0; i<len(requests); i++ {

			if requests[i] > params.upperCYL || requests[i] < params.lowerCYL {
				errors = append(errors,err{requests[i],true})
				continue
			}

			errors = append(errors,err{requests[i],false})
			traversals += Abs(currentPos - requests[i])

			currentPos = requests[i]
		}

		return errors, traversals
	}



	//service all the requests in an ascending order from
	//out starting index
	for i:=startIndex; i<len(requests); i++ {

		if requests[i] > params.upperCYL || requests[i] < params.lowerCYL {
			errors = append(errors,err{requests[i],true})
			continue
		}


		errors = append(errors,err{requests[i],false})
		traversals += Abs(currentPos - requests[i])

		currentPos = requests[i]
	}

	traversals += Abs(currentPos-requests[0])

	//set our head position to the first index
	currentPos = requests[0]



	//Now we service the remaining requests in ascending order
	for i:=0; i<startIndex; i++ {

		if requests[i] > params.upperCYL || requests[i] < params.lowerCYL {
			errors = append(errors,err{requests[i],true})
			continue
		}


		errors = append(errors,err{requests[i],false})
		traversals += Abs(currentPos - requests[i])

		currentPos = requests[i]
	}

	return errors, traversals
}

func main() {


	inputFile := os.Args[1]

	params, requests, requestsPrim := parseFile(inputFile)

	var errors []err
	var traversals int

	if checkAbortConditions(params) == true {
		return
	}




	if strings.Compare(params.alg, "fcfs") == 0 {

		errors,traversals =	fcfs(params,requestsPrim)

	} else if strings.Compare(params.alg, "sstf") == 0 {

		errors,traversals =	sstf(params,requests)

	}else if strings.Compare(params.alg, "scan") == 0 {

		errors,traversals =	scan(requestsPrim,params)

	}else if strings.Compare(params.alg, "c-scan") == 0 {

		errors,traversals =	cscan(requestsPrim,params)

	}else if strings.Compare(params.alg, "look") == 0 {

		errors,traversals =	look(requestsPrim,params)

	}else if strings.Compare(params.alg, "c-look") == 0 {

		errors,traversals =	clook(requestsPrim,params)

	}

	for j:=0; j<len(errors); j++ {
		if errors[j].err == true {
			fmt.Printf("ERROR(15):Request out of bounds: req (%d) > upper (%d) or < lower (%d)\n", errors[j].val, params.upperCYL, params.lowerCYL)
		}
	}

	fmt.Printf("Seek algorithm: %s\n",strings.ToUpper(params.alg))
	fmt.Printf("\tLower cylinder: %5d\n", params.lowerCYL)
	fmt.Printf("\tUpper cylinder: %5d\n", params.upperCYL)
	fmt.Printf("\tInit cylinder:  %5d\n", params.initCYL)

	fmt.Printf("\tCylinder requests:\n")

	for i := 0; i<len(requests); i++ {
		fmt.Printf("\t\tCylinder %5d\n", requests[i].val)
	}

	for k:=0; k<len(errors); k++ {
		if errors[k].err == false {
			fmt.Printf("Servicing %5d\n", errors[k].val)
		}
	}

	fmt.Printf("%s traversal count = %5d\n",strings.ToUpper(params.alg),traversals)





}
