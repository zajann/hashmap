package hashmap

func getPrimeNum(n int) int {
	if n < 3 {
		return 2
	}

	if n%2 == 0 {
		n++
	}

	div := 3
	for div < n {
		if n%div == 0 {
			div = 3
			n += 2
		} else {
			div += 2
		}
	}

	return n
}

func getMaxCoprime(n int) int {
	for i := n - 1; i > 1; i-- {
		if isComprime(n, i) {
			return i
		}
	}
	return 1
}

func getMinCoprime(n int) int {
	for i := 2; i < n; i++ {
		if isComprime(n, i) {
			return i
		}
	}
	return 1
}

func isComprime(a, b int) bool {
	if a == 0 || b == 0 {
		return false
	}

	if a == b {
		return false
	}

	if gcd(a, b) == 1 {
		return true
	}

	return false
}

func gcd(a, b int) int {
	if b == 0 {
		return a
	}
	return gcd(b, a%b)
}
