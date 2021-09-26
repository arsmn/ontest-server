package handler

import (
	"net/http"
)

const message string = `
  ____     _______        _                _____ _____ 
 / __ \   |__   __|      | |         /\   |  __ \_   _|
| |  | |_ __ | | ___  ___| |_       /  \  | |__) || |  
| |  | | '_ \| |/ _ \/ __| __|     / /\ \ |  ___/ | |  
| |__| | | | | |  __/\__ \ |_     / ____ \| |    _| |_ 
 \____/|_| |_|_|\___||___/\__|   /_/    \_\_|   |_____|
`

func (h *Handler) root(rw http.ResponseWriter, r *http.Request) {
	rw.Write([]byte(message))
}
