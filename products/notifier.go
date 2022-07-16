package products

import (
	"fmt"
	"time"
)

type Notifier interface {
	NotifyError(message string)
}

type NotifierImpl struct{}

func NewNotifier() *NotifierImpl {
	return &NotifierImpl{}
}

func (n NotifierImpl) NotifyError(message string) {
	fmt.Println(fmt.Sprintf("ERROR on %s - %s", time.Now().String(), message))
}
