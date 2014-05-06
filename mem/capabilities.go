package mem

import (
)

type graphCaps struct {}

func (caps *graphCaps) Persistent() bool {
	return false;
}

func (caps *graphCaps) SortedKeys() bool {
	return true;
}

func (caps *graphCaps) KeyIndex() bool {
	return true;
}