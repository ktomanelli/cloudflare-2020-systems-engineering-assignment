package main

import (
	"flag"
	"fmt"
	"net/http"
	"io/ioutil"
	"strconv"
	"time"
	"sort"
)

type requestData struct{
	duration time.Duration
	size int
	err error
}

var flagData = make(map[string]map[string]string)

func init(){
	//adding all flag options to map
	flagData["url"]=map[string]string{
		"flag":"url",
		"default":" ",
		"desc":"A URL you'd like to perform a HTTP GET request for. Ex:'https://google.com/'",
		}
	flagData["help"]=map[string]string{
		"flag":"help",
		"default":" ",
		"desc":"Use this flag to view all flag options",
		}
	flagData["profile"]=map[string]string{
		"flag":"profile",
		"default":"1",
		"desc":"Number of requests you'd like to make. Will return speed/analytic data about requests. Must be a positive integer.",
	}
	//iterating through flagData to create each flag option
	for k := range flagData{
		flag.String(flagData[k]["flag"],flagData[k]["default"],flagData[k]["desc"])
	}
	flag.Parse()
}

func main(){
	requests:=parseFlags()
	profile,err:=strconv.Atoi(requests["profile"])
	if requests != nil{
		if profile==0{
			profile++
		}
		if err!=nil&&profile<=0{
			fmt.Println("Error in Profile value. Must be Positive Integer")
		}else{
			data:=makeRequest(requests["url"],profile)
			if data!=nil{
				printData(data)
			}
		}
	}
}

func makeRequest(url string,profile int)map[int]requestData{
	if url!=""{
		completeData:=make(map[int]requestData)
		for i:=0;i<profile;i++{
			start:=time.Now()
			resp,err := http.Get(url)
			if err != nil{
				fmt.Println("Error in HTTP Get")
				completeData[i] = requestData{0,0,err}
				continue
			}
			defer resp.Body.Close()
			duration:=time.Now().Sub(start)
		
			body,err := ioutil.ReadAll(resp.Body)
			if err != nil{
				fmt.Println("Error in ioutil ReadAll")
				completeData[i] = requestData{duration,0,err}
			}		
			size:=len(body)

			if profile==1{
				printResp(string(body))
				return nil
			}
			completeData[i] = requestData{duration,size,nil}
		}
		return completeData
	}else{
		fmt.Println("Please specify a URL")
		return nil
	}
}

func flagPresent(name string)bool{
	isPresent := false
	flag.Visit(func(f *flag.Flag){
		if f.Name==name{
			isPresent=true
		}
	})
	return isPresent
}
//iterates through flags to determine which flags were passed
func parseFlags()map[string]string{
	requests:=make(map[string]string)
	for k := range flagData{
		if flagPresent(k){
			currentFlag := flag.Lookup(k)
			val := fmt.Sprintf("%v",currentFlag.Value)
			//if help flag is present, output help text and exit
			if k == "help" && val!="false"{
				printHelp()
				return nil
			}else{
				//else add flag to requests map
				requests[k] = val
			}
		}
	}
	return requests
}
func printHelp(){
	fmt.Println("---------------HELP---------------")
	fmt.Println("Name-----Default Value-----Usage")
	for key:=range flagData{
		fmt.Printf("-%v----------%v----------%v\n",key,flagData[key]["default"],flagData[key]["desc"])
	}
	fmt.Println("----------------------------------")
}

func printResp(body string){
	fmt.Println("---------------URL---------------")
	fmt.Println(body)
	fmt.Println("----------------------------------")
}

func printData(data map[int]requestData){
	var maxTime, minTime, mean, median time.Duration
	var maxSize, minSize, percentCompleted int
	var errStr string
	durations:=make(map[int]time.Duration)
	sizes:=make(map[int]int)
	errors:=make(map[int]error)
	for k :=range data{
		durations[k] = data[k].duration
		sizes[k] = data[k].size
		errors[k] = data[k].err
	}
	maxTime,minTime = findMaxMinDuration(durations)
	mean = findMean(durations)
	median = findMedian(durations)
	maxSize,minSize = findMaxMinSize(sizes)
	percentCompleted = findPercentCompleted(errors)
	errStr = getErrStr(errors)

	fmt.Println("-------------PROFILE-------------")
	fmt.Printf("Number of requests: %v\n",len(data))
	fmt.Printf("Fastest time: %v\n",minTime)
	fmt.Printf("Slowest time: %v\n",maxTime)
	fmt.Printf("Mean time: %v\n",mean)
	fmt.Printf("Median time: %v\n",median)
	fmt.Printf("Requests completed: %v%%\n",percentCompleted)
	fmt.Printf("Errors: %v\n",errStr)
	fmt.Printf("Size of smallest response: %v bytes\n",minSize)
	fmt.Printf("Size of largest response: %v bytes\n",maxSize)
	fmt.Println("----------------------------------")

}

func findMaxMinSize(data map[int]int)(int,int){
	max:=data[0]
	min:=data[0]
	for k:=range data{
		if data[k]>max{
			max=data[k]
		}
		if data[k]<min{
			min=data[k]
		}
	}
	return max,min
}
func findMaxMinDuration(data map[int]time.Duration)(time.Duration,time.Duration){
	max:=data[0]
	min:=data[0]
	for k:=range data{
		if data[k]>max{
			max=data[k]
		}
		if data[k]<min{
			min=data[k]
		}
	}
	return max,min
}

func findMean(data map[int]time.Duration)time.Duration{
	var sum int64 = 0
	for k:=range data{
		sum+=int64(data[k])
	}
	return time.Duration(sum/int64(len(data)))
}

func findMedian(data map[int]time.Duration)time.Duration{
	length:=len(data)
	sorted:=make([]int,length)
	for k:=range data{
		sorted[k] = int(data[k])
	}
	sort.Ints(sorted)
	if(length%2!=0){
		return time.Duration(sorted[length/2])
	}
	return time.Duration((sorted[(length/2)-1]+sorted[length/2])/2)
}

func findPercentCompleted(data map[int]error)int{
	completed:=len(data)
	for k:=range data{
		if data[k]!=nil{
			completed--
		}
	}
	return (completed/len(data))*100
}

func getErrStr(data map[int]error)string{
	errStr:=""
	for k:=range data{
		if data[k]!=nil{
			errStr+=data[k].Error()
		}
	}
	if errStr==""{
		return "There were no Errors."
	}else{
		return errStr
	}
}