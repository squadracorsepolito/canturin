package main

type request struct {
	data any
}

// newRequest returns a new generic request.
//
// IMPORTANT: the data attribute must be a pointer.
func newRequest(data any) *request {
	return &request{data}
}

/////////////////////
// COMMON REQUESTS //
/////////////////////

type commonCreateReq struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
}

type UpdateNameReq struct {
	Name string `json:"name"`
}

func (r *request) toUpdateName() *UpdateNameReq {
	req, ok := r.data.(*UpdateNameReq)
	if !ok {
		panic("cannot convert to UpdateNameReq")
	}
	return req
}

type UpdateDescReq struct {
	Desc string `json:"desc"`
}

func (r *request) toUpdateDesc() *UpdateDescReq {
	req, ok := r.data.(*UpdateDescReq)
	if !ok {
		panic("cannot convert to UpdateDescReq")
	}
	return req
}

//////////////////
// BUS REQUESTS //
//////////////////

type CreateBusReq struct {
	commonCreateReq

	BusType  BusType `json:"busType"`
	Baudrate int     `json:"baudrate"`
}

func (r *request) toCreateBus() *CreateBusReq {
	req, ok := r.data.(*CreateBusReq)
	if !ok {
		panic("cannot convert to CreateBusReq")
	}
	return req
}

type UpdateBaudrateReq struct {
	Baudrate int `json:"baudrate"`
}

func (r *request) toUpdateBaudrate() *UpdateBaudrateReq {
	req, ok := r.data.(*UpdateBaudrateReq)
	if !ok {
		panic("cannot convert to UpdateBaudrateReq")
	}
	return req
}

type UpdateBusTypeReq struct {
	BusType BusType `json:"busType"`
}

func (r *request) toUpdateBusType() *UpdateBusTypeReq {
	req, ok := r.data.(*UpdateBusTypeReq)
	if !ok {
		panic("cannot convert to UpdateBusTypeReq")
	}
	return req
}
