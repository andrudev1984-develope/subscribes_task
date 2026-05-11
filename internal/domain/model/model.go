package model

type Subscribe struct {
	ID
	UserId ID
	ServiceName
	Price
	StartDate Date
	EndDate   *Date
}
