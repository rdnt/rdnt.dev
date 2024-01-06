package event

import (
	"errors"
	"slices"

	"github.com/rdnt/rdnt.dev/uuid"
)

type Event interface {
	Type() Type
	AggregateType() AggregateType
	AggregateId() uuid.UUID
}

type AggregateType string

func (t AggregateType) String() string {
	return string(t)
}

const (
	Project AggregateType = "project"
)

type Type string

func (t Type) String() string {
	return string(t)
}

func TypeFromString(s string) (Type, error) {
	if !slices.Contains(Types, Type(s)) {
		return "", errors.New("invalid event type")
	}

	return Type(s), nil
}

var Types = []Type{
	ProjectCreated,
}
