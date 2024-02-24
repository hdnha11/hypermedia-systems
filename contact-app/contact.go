package main

import (
	"context"
	"fmt"
	"strings"
	"sync"
)

type Contact struct {
	ID    int    `form:"id"`
	First string `form:"first"`
	Last  string `form:"last"`
	Phone string `form:"phone"`
	Email string `form:"email"`
}

type InMemContactRepository struct {
	db     []*Contact
	mu     sync.Mutex
	nextID int
}

func (r *InMemContactRepository) GetAll(ctx context.Context) ([]*Contact, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.db, nil
}

func (r *InMemContactRepository) Search(ctx context.Context, q string) ([]*Contact, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var contacts []*Contact
	for _, c := range r.db {
		matchFirst := c.First != "" && strings.Contains(c.First, q)
		matchLast := c.Last != "" && strings.Contains(c.Last, q)
		matchEmail := c.Email != "" && strings.Contains(c.Email, q)
		matchPhone := c.Phone != "" && strings.Contains(c.Phone, q)
		if matchFirst || matchLast || matchEmail || matchPhone {
			contacts = append(contacts, c)
		}
	}

	return contacts, nil
}

func (r *InMemContactRepository) Save(ctx context.Context, contact *Contact) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !strings.Contains(contact.Email, "@") {
		return fmt.Errorf("%s is an invalid email address", contact.Email)
	}

	if contact.ID == 0 {
		contact.ID = r.nextID
		r.nextID++
		r.db = append(r.db, contact)
	} else {
		_, found, err := r.find(contact.ID)
		if err != nil {
			return err
		}

		found.First = contact.First
		found.Last = contact.Last
		found.Phone = contact.Phone
		found.Email = contact.Email
	}

	return nil
}

func (r *InMemContactRepository) Find(ctx context.Context, id int) (*Contact, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, contact, err := r.find(id)
	return contact, err
}

func (r *InMemContactRepository) find(id int) (int, *Contact, error) {
	for i, contact := range r.db {
		if id == contact.ID {
			return i, contact, nil
		}
	}

	return -1, nil, fmt.Errorf("Contact with ID=%d is not found", id)
}

func (r *InMemContactRepository) Delete(ctx context.Context, id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	foundIndex, _, err := r.find(id)
	if err != nil {
		return err
	}

	r.db = append(r.db[:foundIndex], r.db[foundIndex+1:]...)

	return nil
}
