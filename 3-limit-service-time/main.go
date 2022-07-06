//////////////////////////////////////////////////////////////////////
//
// Your video processing service has a freemium model. Everyone has 10
// sec of free processing time on your service. After that, the
// service will kill your process, unless you are a paid premium user.
//
// Beginner Level: 10s max per request
// Advanced Level: 10s max per user (accumulated)
//

package main

import (
	"sync/atomic"
	"time"
)

const MAX_TIME = 10

// User defines the UserModel. Use this to check whether a User is a
// Premium user or not
type User struct {
	ID        int
	IsPremium bool
	TimeUsed  int64 // in seconds
}

// HandleRequest runs the processes requested by users. Returns false
// if process had to be killed
func HandleRequest(process func(), u *User) bool {
	if u.IsPremium {
		return premiumProcess(process)
	} else {
		return freemiumProcess(process, u)
	}
}

func main() {
	RunMockServer()
}

func premiumProcess(process func()) bool {
	process()
	return true
}

func freemiumProcess(process func(), u *User) bool {
	if u.IsOutOfTime() {
		return false
	}

	done := make(chan bool)
	tick := time.Tick(1 * time.Second)

	go func() {
		process()
		done <- true
	}()

	for {
		select {
		case <-tick:
			if t := u.AddTime(1); t >= MAX_TIME {
				return false
			}
		case <-done:
			return true
		}
	}
}

func (u *User) AddTime(seconds int64) int64 {
	return atomic.AddInt64(&u.TimeUsed, seconds)
}

func (u *User) IsOutOfTime() bool {
	return atomic.LoadInt64(&u.TimeUsed) >= MAX_TIME
}
