package api

type BookingStore interface {
	Get(id int)
	Insert(booking interface{})
	Delete(id int)
	Update(id int, booking interface{})
}

type BookingService struct {
	s BookingStore
}

func (b *BookingService) CreateBooking(booking interface{}) {
	b.s.Insert(booking)
}
