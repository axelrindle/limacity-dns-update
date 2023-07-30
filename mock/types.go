package mock

import "github.com/axelrindle/limacity-dns-update/shared"

type Route struct {
	Path string `json:"path"`
}

type RequestUpdateRecord struct {
	Record shared.NameserverRecord `json:"nameserver_record"`
}
