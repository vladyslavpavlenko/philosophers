package main

import (
	"fmt"
	"sync"
	"time"
)

const (
	hunger = 3 // times to eat

	eatTime   = 1 * time.Second
	thinkTime = 3 * time.Second
)

type Philosopher struct {
	name string
	lf   int
	rf   int
}

/*
             2
         (A)   (H)
        1         3
	 (S)    	  (L)
		 0      4
		   (P)
*/

var philosophers = []Philosopher{
	{"Plato", 0, 4},
	{"Socrates", 1, 0},
	{"Aristotle", 2, 1},
	{"Hippocrates", 3, 2},
	{"Locke", 4, 3},
}

func main() {
	setup()
	fmt.Println("Done")
}

func setup() {
	wg := &sync.WaitGroup{}
	wg.Add(len(philosophers))

	seated := &sync.WaitGroup{}
	seated.Add(len(philosophers))

	forks := make(map[int]*sync.Mutex)
	for i := 0; i < len(philosophers); i++ {
		forks[i] = &sync.Mutex{}
	}

	for i := 0; i < len(philosophers); i++ {
		go dine(philosophers[i], forks, wg, seated)
	}

	wg.Wait()
}

func dine(phil Philosopher, forks map[int]*sync.Mutex, wg *sync.WaitGroup, seated *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("[%s joined]\n", phil.name)
	seated.Done()

	seated.Wait()

	for i := hunger; i > 0; i-- {
		if phil.lf > phil.rf {
			forks[phil.rf].Lock()
			fmt.Printf("\t%s takes the right fork.\n", phil.name)

			forks[phil.lf].Lock()
			fmt.Printf("\t%s takes the left fork.\n", phil.name)
		} else {
			forks[phil.lf].Lock()
			fmt.Printf("\t%s takes the left fork.\n", phil.name)

			forks[phil.rf].Lock()
			fmt.Printf("\t%s takes the right fork.\n", phil.name)
		}

		fmt.Printf("\t%s has both forks and is eating (%d/%d).\n", phil.name, hunger-i+1, hunger)
		time.Sleep(eatTime)

		fmt.Printf("\t%s is thinking.\n", phil.name)
		time.Sleep(thinkTime)

		forks[phil.lf].Unlock()
		forks[phil.rf].Unlock()

		fmt.Printf("\t%s put down the forks.\n", phil.name)
	}

	fmt.Printf("[%s left]\n", phil.name)
}
