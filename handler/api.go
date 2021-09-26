package handler

import (
	"net/http"
)

const root string = `
  ____     _______        _                _____ _____ 
 / __ \   |__   __|      | |         /\   |  __ \_   _|
| |  | |_ __ | | ___  ___| |_       /  \  | |__) || |  
| |  | | '_ \| |/ _ \/ __| __|     / /\ \ |  ___/ | |  
| |__| | | | | |  __/\__ \ |_     / ____ \| |    _| |_ 
 \____/|_| |_|_|\___||___/\__|   /_/    \_\_|   |_____|
`

func (h *Handler) root(rw http.ResponseWriter, r *http.Request) {
	rw.Write([]byte(root))
}

const notFound string = `
 _   _       _   ______                    _ 
| \ | |     | | |  ____|                  | |
|  \| | ___ | |_| |__ ___  _   _ _ __   __| |
| .   |/ _ \| __|  __/ _ \| | | | '_ \ / _  |
| |\  | (_) | |_| | | (_) | |_| | | | | (_| |
|_| \_|\___/ \__|_|  \___/ \__,_|_| |_|\__,_|
`

func (h *Handler) notFound(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusNotFound)
	rw.Write([]byte(notFound))
}

const methodNotAllowed string = `
 __  __      _   _               _ _   _       _            _ _                       _ 
|  \/  |    | | | |             | | \ | |     | |     /\   | | |                     | |
| \  / | ___| |_| |__   ___   __| |  \| | ___ | |_   /  \  | | | _____      _____  __| |
| |\/| |/ _ \ __| '_ \ / _ \ / _  | .   |/ _ \| __| / /\ \ | | |/ _ \ \ /\ / / _ \/ _  |
| |  | |  __/ |_| | | | (_) | (_| | |\  | (_) | |_ / ____ \| | | (_) \ V  V /  __/ (_| |
|_|  |_|\___|\__|_| |_|\___/ \__,_|_| \_|\___/ \__/_/    \_\_|_|\___/ \_/\_/ \___|\__,_|
`

func (h *Handler) methodNotAllowed(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusMethodNotAllowed)
	rw.Write([]byte(methodNotAllowed))
}
