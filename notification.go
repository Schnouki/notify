// Copyright (c) 2013, Ben Morgan. All rights reserved.
// Use of this source code is governed by an MIT license
// that can be found in the LICENSE file.

package notify

import (
	"time"

	"github.com/godbus/dbus"
)

// NotificationUrgency can be either LowUrgency, NormalUrgency, and CriticalUrgency.
// It is conceivable that some notification daemons make no distinction between the
// different urgencies, but enough do that it makes sense to use them.
type NotificationUrgency byte

const (
	LowUrgency      NotificationUrgency = iota // LowUrgency probably shouldn't even be shown ;-)
	NormalUrgency                              // NormalUrgency is for information that is interesting.
	CriticalUrgency                            // CriticalUrgency is for errors or severe events.
)

// asHint returns the NotificationUrgency in the type that the DBus
// specification requires.
func (u NotificationUrgency) asHint() map[string]dbus.Variant {
	return map[string]dbus.Variant{"urgency": dbus.MakeVariant(byte(u))}
}

// Notification is there to provide you with full power of your notifications.
// It is possible for you to use a Notification as you use the notify library
// without them. This allows for multiple defaults.
//
// For example:
//
//	func main() {
//		critical := notify.New("prog", "", "", "critical-icon.png", time.Duration(0), notify.CriticalUrgency)
//		boring := notify.New("prog", "", "", "low-icon.png", 1 * time.Second, notify.LowUrgency)
//		boring.ReplaceMsg("Nothing is happening... boring!", "")
//		critical.ReplaceMsg("Your computer is on fire!", "Here is what you should do:\n ...")
//	}
//
type Notification struct {
	// Name represents the application name sending the notification.  This is
	// optional and can be the empty string "".
	Name string
	// Summary represents the subject of the notification.
	Summary string
	// Body represents the main body with extra details. Some notification
	// daemons ignore the body; it is optional and can be the empty string "".
	Body string

	// IconPath is a path to an icon that should be used for the notification.
	// Some notification daemons ignore the icon path; it is optional and can
	// be the empty string "".
	IconPath string
	// Timeout is the requested timeout for the notification. Some notification
	// daemons override the requested timeout. A value of 0 is a request that
	// it not timeout at all.
	Timeout time.Duration
	// Urgency determines the urgency of the notification, which can be one of
	// LowUrgency, NormalUrgency, and CriticalUrgency.
	Urgency NotificationUrgency

	// Id is the ID of the notification. It is 0 initially, and will be
	// updated when calling Send or one of the Replace methods.
	Id uint32
}

// New returns a pointer to a new Notification.
func New(name, summary, body, icon string, timeout time.Duration, urgency NotificationUrgency) *Notification {
	return &Notification{name, summary, body, icon, timeout, urgency, 0}
}

// Send sends the notification n as it is, and returns an err, possibly nil.
func (n Notification) Send() (err error) {
	n.Id, err = notify(n.Name, n.Summary, n.Body, n.IconPath, n.Id, nil, n.Urgency.asHint(), n.timeoutInMS())
	return err
}

// ReplaceMsg is identical to notify.ReplaceMsg, except that the rest of the
// values come from n.
func (n Notification) ReplaceMsg(summary, body string) (err error) {
	n.Summary, n.Body = summary, body
	return n.Send()
}

// ReplaceUrgentMsg is identical to notify.ReplaceUrgentMsg, except that the
// rest of the values come from n.
func (n Notification) ReplaceUrgentMsg(summary, body string, urgency NotificationUrgency) (err error) {
	n.Summary, n.Body, n.Urgency = summary, body, urgency
	return n.Send()
}

// timeoutInMS returns Timeout in milliseconds.
//
// The specification specifies that the timeout is the number of milliseconds
// that the notification should be displayed.
func (n Notification) timeoutInMS() int32 {
	return int32(n.Timeout / time.Millisecond)
}
