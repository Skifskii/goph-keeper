package secret

type Secret struct {
	Payload  *Payload
	UserID   int
	Metadata string
}

func New(payload *Payload, userID int, metadata string) (Secret, error) {
	// TODO: validate payload len
	
	return Secret{
		Payload:  payload,
		UserID:   userID,
		Metadata: metadata,
	}, nil
}
