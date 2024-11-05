package main

import "fmt"
import "time"


func main() {
	c := make(chan string)
	// Start a goroutine to send messages to the channel
	go func() {
		defer close(c)
		for i := 0; i < 10; i++ {
			c <- fmt.Sprintf("Hello, World! %d", i)
			// Sleep for 1 second
			time.Sleep(500 * time.Millisecond)
		}
		// Close the channel
	}()
	
	fmt.Println("Hello, World!")
	// Iterate over the channel and print the values until the channel is empty
	for msg := range c {
		fmt.Println(msg)
	}
}
