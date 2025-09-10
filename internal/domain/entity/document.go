package entity

import "time"

// Document representa un documento en el dominio
type Document struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	Version      string    `json:"version"`
	Attachments  []string  `json:"attachments"`
	Contributors []User    `json:"contributors"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// NewDocument crea una nueva instancia de Document
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

// AddContributor añade un contribuidor al documento
func (d *Document) AddContributor(user User) {
	d.Contributors = append(d.Contributors, user)
	d.UpdatedAt = time.Now()
}

// AddAttachment añade un adjunto al documento
func (d *Document) AddAttachment(attachment string) {
	d.Attachments = append(d.Attachments, attachment)
	d.UpdatedAt = time.Now()
}

// Validate valida los datos del documento
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
