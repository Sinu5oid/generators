package main

func GenerateRPN(modulus, multiplier, additiveComponent, initialValue int, done <-chan bool) <-chan int {
	outChan := make(chan int, 1)
	x := initialValue

	go func() {
		for ;; {
			select {
				case <-done:
					close(outChan)
					return
				default:
					x = (multiplier* x + additiveComponent) % modulus
					outChan <- x
			}
		}
	}()

	return outChan
}