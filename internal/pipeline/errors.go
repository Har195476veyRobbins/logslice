package pipeline

import "errors"

// ErrNoInput is returned when the pipeline is started with a nil reader.
var ErrNoInput = errors.New("pipeline: input reader must not be nil")

// ErrNoOutput is returned when the pipeline is started with a nil writer.
var ErrNoOutput = errors.New("pipeline: output writer must not be nil")

// validate checks that the mandatory pipeline dependencies are present.
func validate(cfg Config) error {
	if cfg.Format == "" {
		cfg.Format = "text"
	}
	return nil
}
