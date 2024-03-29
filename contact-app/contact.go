package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
)

const PAGE_SIZE = 10

type Contact struct {
	ID    int    `json:"id" form:"id"`
	First string `json:"first" form:"first"`
	Last  string `json:"last" form:"last"`
	Phone string `json:"phone" form:"phone"`
	Email string `json:"email" form:"email"`
}

type InMemContactRepository struct {
	db      []*Contact
	mu      sync.Mutex
	autoInc int
}

func (r *InMemContactRepository) Load(filePath string) error {
	bs, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	var contacts []*Contact
	if err := json.Unmarshal(bs, &contacts); err != nil {
		return err
	}

	r.db = contacts
	r.autoInc = r.maxID()

	return nil
}

func (r *InMemContactRepository) maxID() int {
	var m int
	for _, c := range r.db {
		m = max(m, c.ID)
	}
	return m
}

func (r *InMemContactRepository) GetAll(ctx context.Context, page int) ([]*Contact, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	totalItems := len(r.db)
	start := min((page-1)*PAGE_SIZE, totalItems)
	end := min(start+PAGE_SIZE, totalItems)

	return r.db[start:end], nil
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

func (r *InMemContactRepository) Validate(contact *Contact) error {
	verr := &ValidationError{}

	if contact.Email == "" {
		verr.AddFieldError("Email", "Email Required")
	} else if !strings.Contains(contact.Email, "@") {
		verr.AddFieldError("Email", "Invalid Email Address")
	}
	for _, c := range r.db {
		if c.ID != contact.ID && c.Email == contact.Email {
			verr.AddFieldError("Email", "Email Must Be Unique")
			break
		}
	}

	if contact.First == "" {
		verr.AddFieldError("First", "First Name Required")
	}

	if verr.FieldErrors() != nil {
		return verr
	}

	return nil
}

func (r *InMemContactRepository) Save(ctx context.Context, contact *Contact) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := r.Validate(contact); err != nil {
		return err
	}

	if contact.ID == 0 {
		r.autoInc++
		contact.ID = r.autoInc
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
