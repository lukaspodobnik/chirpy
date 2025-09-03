package main

import "net/http"

func healthzHandlerFunc(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Add(contenttype, textPlainUTF8)
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte(http.StatusText(http.StatusOK) + "\n"))
}
