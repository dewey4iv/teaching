package room

// Option defines something that builds/configures a Room
type Option interface {
	Apply(*Room) error
}

// WithFilterWords sets the words that the room will filter on.
// Any words found in the provided list will keep the entire message
// from being broadcast
func WithFilterWords(words []string) Option {
	return &withFilterWords{words}
}

type withFilterWords struct {
	words []string
}

func (opt *withFilterWords) Apply(r *Room) error {
	r.filterWords = opt.words

	return nil
}
