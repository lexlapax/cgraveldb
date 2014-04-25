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
