package main

import "time"

type fixtures struct {
	Results int `json:"results"`
	Paging  struct {
		Current int `json:"current"`
		Total   int `json:"total"`
	} `json:"paging"`
	Response []struct {
		Fixture struct {
			Id        int       `json:"id"`
			Date      time.Time `json:"date"`
			Timestamp int       `json:"timestamp"`
			Status    struct {
				Long    string      `json:"long"`
				Short   string      `json:"short"`
				Elapsed interface{} `json:"elapsed"`
				Extra   interface{} `json:"extra"`
			} `json:"status"`
		} `json:"fixture"`
		Teams struct {
			Home struct {
				Id   int    `json:"id"`
				Name string `json:"name"`
			} `json:"home"`
			Away struct {
				Id   int    `json:"id"`
				Name string `json:"name"`
			} `json:"away"`
		} `json:"teams"`
		Goals struct {
			Home interface{} `json:"home"`
			Away interface{} `json:"away"`
		} `json:"goals"`
	} `json:"response"`
}
