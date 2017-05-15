package main

import (
	"fmt"
	"math/rand"
	"time"
)

//someLongFunc is a function that might
//take a while to complete, so we want
//to run it on its own go routine
func someLongFunc(ch chan int) { //passing in channel of ints
	r := rand.Intn(2000)             //create a new int between 1 -> 2000
	d := time.Duration(r)            //time duration
	time.Sleep(time.Millisecond * d) //sleep for this many milliseconds
	//writing r into this channel; this will block if ch is full
	//otherwise keeps going
	ch <- r

}

func main() {
	//TODO:
	//create a channel and call
	//someLongFunc() on a go routine
	//passing the channel so that
	//someLongFunc() can communicate
	//its results
	rand.Seed(time.Now().UnixNano())
	fmt.Println("starting long-running func...")
	n := 10
	ch := make(chan int, n) //make channels using make (data type, capacity?)
	start := time.Now()
	for i := 0; i < n; i++ {
		go someLongFunc(ch)
	}
	for i := 0; i < n; i++ {
		//wait for the func to reply
		result := <-ch //read result out of channel
		fmt.Printf("result was %d\n", result)
	}
	fmt.Printf("took %v\n", time.Since(start))
}
