package main

import "fmt"

type Lendable interface {
	IsLendable() bool
}

type Reservable interface {
	IsReservable() bool
}

type Information interface {
	GetInfo() string
}

type Book struct {
	Title  string
	Author string
}

type EBook struct{ Book }
func (e EBook) IsLendable() bool { return true }
func (e EBook) IsReservable() bool { return false }

type ReferenceBook struct{ Book }
func (r ReferenceBook) IsLendable() bool { return false }
func (r ReferenceBook) IsReservable() bool { return false }

type GeneralBook struct{ Book }
func (g GeneralBook) IsLendable() bool { return true }
func (g GeneralBook) IsReservable() bool { return true }

func (e EBook) GetInfo() string {
	return fmt.Sprintf("Title: %s, Author: %s, Lendable: %t, Reservable: %t", e.Title, e.Author, e.IsLendable(), e.IsReservable())
}

func (r ReferenceBook) GetInfo() string {
	return fmt.Sprintf("Title: %s, Author: %s, Lendable: %t, Reservable: %t", r.Title, r.Author, r.IsLendable(), r.IsReservable())
}

func (g GeneralBook) GetInfo() string {
	return fmt.Sprintf("Title: %s, Author: %s, Lendable: %t, Reservable: %t", g.Title, g.Author, g.IsLendable(), g.IsReservable())
}

func NewEBook(title, author string) EBook {
	return EBook{Book: Book{Title: title, Author: author}}
}

func NewReferenceBook(title, author string) ReferenceBook {
	return ReferenceBook{Book: Book{Title: title, Author: author}}
}

func NewGeneralBook(title, author string) GeneralBook {
	return GeneralBook{Book: Book{Title: title, Author: author}}
}

func main() {
	ebook := NewEBook("The Great Gatsby", "F. Scott Fitzgerald")
	referenceBook := NewReferenceBook("Programming Reference", "Jane Smith")
	generalBook := NewGeneralBook("Digital Guide", "Bob Wilson")

	fmt.Println(ebook.GetInfo())
	fmt.Println(referenceBook.GetInfo())
	fmt.Println(generalBook.GetInfo())
}