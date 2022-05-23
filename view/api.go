package view

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"algotrade_service/model"
)

type View struct {
	eventController *model.EventController
}

func NewView(eventController *model.EventController) *View {

	return &View{
		eventController: eventController,
	}
}

func (v *View) Start() error {
	// pprof
	http.HandleFunc("/event", v.event)
	return http.ListenAndServe(":8080", nil)
}

type EventReq struct {
	Event string `json:"event"`
}

func (v *View) event(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// return event
		data, err := json.Marshal(v.eventController.GetEvents())
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		w.Write(data)
	case http.MethodPost:
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		req := EventReq{}
		err = json.Unmarshal(data, &req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		err = v.eventController.AddNewEventFromRaw(req.Event)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
	}
}
