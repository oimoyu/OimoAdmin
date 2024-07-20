package fail2ban

import "sync"

type IPAttempt struct {
	Attempts int
}

var ipAttempts = make(map[string]int)
var bannedIPs = make(map[string]bool)
var mu sync.RWMutex

const (
	maxAttempts = 5

	maxBannedIPs  = 100000
	maxAttemptIPs = 100000
)

func ResetIPAttempts(ip string) {
	mu.Lock()
	defer mu.Unlock()

	delete(ipAttempts, ip)
	delete(bannedIPs, ip)
}
func IncrementIPAttempts(ip string) {
	mu.Lock()
	defer mu.Unlock()

	// return early to reduce lock holding time
	if bannedIPs[ip] {
		return
	}

	// if ip not exist, ipAttempts[ip] will auto init it as 0
	ipAttempts[ip]++

	if ipAttempts[ip] >= maxAttempts {
		bannedIPs[ip] = true
		delete(ipAttempts, ip)
	}

	if len(bannedIPs) > maxBannedIPs {
		bannedIPs = make(map[string]bool)
	}
	if len(ipAttempts) > maxAttemptIPs {
		ipAttempts = make(map[string]int)
	}

}

func IsIPValid(ip string) bool {
	mu.RLock()
	defer mu.RUnlock()

	return !bannedIPs[ip]
}
func RemainingAttempts(ip string) int {
	mu.RLock()
	defer mu.RUnlock()

	if bannedIPs[ip] {
		return 0
	}

	attempts, exists := ipAttempts[ip]
	if !exists {
		return maxAttempts
	}

	return maxAttempts - attempts
}
