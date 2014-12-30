/* -.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.

* File Name : pool.go

* Purpose :

* Creation Date : 12-30-2014

* Last Modified : Tue 30 Dec 2014 07:00:40 PM UTC

* Created By : Kiyor

_._._._._._._._._._._._._._._._._._._._._.*/

package main

import (
	"time"
)

type Client struct{}

func newClient() (c *Client) {
	return &Client{}
}

type Pool struct {
	pool chan *Client
}

func NewPool(max int) *Pool {
	p := Pool{
		pool: make(chan *Client, max),
	}
	for i := 0; i < max; i++ {
		p.pool <- newClient()
	}
	return &p
}

func (p *Pool) Borrow() *Client {
	var c *Client
	for {
		select {
		case c = <-p.pool:
			return c
		default:
			time.Sleep(1 * time.Millisecond)
			continue
		}
	}
}

func (p *Pool) Return(c *Client) {
	select {
	case p.pool <- c:
	default:
	}
}
