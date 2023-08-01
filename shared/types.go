package shared

import "context"

type NameserverRecord struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	Content string `json:"content"`
	TTL     int    `json:"ttl"`
}

type RequestUpdateRecord struct {
	Record NameserverRecord `json:"nameserver_record"`
}

type ResponseUpdateRecord struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

type ResponseListRecords struct {
	Records []NameserverRecord `json:"records"`
}

// operation is a clean up function on shutting down
type ShutdownHook func(ctx context.Context) error
