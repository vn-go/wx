package wx

import (
	"fmt"
	"net/http"
)

func (h *handlerInfo) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		controller, err := utils.controllers.Create(h)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Println(controller)

	}
}
