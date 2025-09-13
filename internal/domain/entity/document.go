package entity

import "time"

// Document represents a document in the domain
type Document struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	Version      string    `json:"version"`
	Attachments  []string  `json:"attachments"`
	Contributors []User    `json:"contributors"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// NewDocument creates a new instance of Document
func NewDocument(id, title, version string) *Document {
	now := time.Now()
	return &Document{
		ID:        id,
		Title:     title,
		Version:   version,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// AddContributor adds a contributor to the document
func (d *Document) AddContributor(user User) {
	d.Contributors = append(d.Contributors, user)
	d.UpdatedAt = time.Now()
}

// AddAttachment adds an attachment to the document
func (d *Document) AddAttachment(attachment string) {
	d.Attachments = append(d.Attachments, attachment)
	d.UpdatedAt = time.Now()
}

// Validate validates the document's data
func (d *Document) Validate() error {
	if d.ID == "" {
		return ErrInvalidDocumentID
	}
	if d.Title == "" {
		return ErrInvalidDocumentTitle
	}
	if d.Version == "" {
		return ErrInvalidDocumentVersion
	}
	return nil
}
