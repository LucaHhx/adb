package adb

import (
	"fmt"

	"github.com/luca/adb/adb/uixml"
)

func (d *Device) XML() (*uixml.Xml, error) {
	data, err := d.UiautomatorDump()
	if err != nil {
		return nil, err
	}
	return uixml.NewXml(data)
}

func (d *Device) ClickButton(name string) error {
	button, err := d.FindButton(name)
	if err != nil {
		return err
	}
	return d.Tap(button.Middle())
}

func (d *Device) ClickNodeBy(node uixml.Node) error {
	return d.Tap(node.Middle())
}

func (d *Device) ClickNode(class, desc string) error {
	node, err := d.FindNode(func(n, pn uixml.Node) bool {
		return n.Class == class && (n.ContentDesc == desc || desc == "" || n.Text == desc)
	})
	if err != nil {
		return err
	}
	return d.Tap(node.Middle())
}

type FindNodeFunc func(n, pn uixml.Node) bool

func (d *Device) FindNode(fn FindNodeFunc) (uixml.Node, error) {
	xml, err := d.XML()
	if err != nil {
		return uixml.Node{}, err
	}
	return xml.Find(fn)
}

func (d *Device) FindNodes(fn FindNodeFunc) ([]uixml.Node, error) {
	xml, err := d.XML()
	if err != nil {
		return nil, err
	}
	list := xml.FindAll(fn)
	if len(list) == 0 {
		return nil, fmt.Errorf("not found")
	}
	return list, nil
}

func (d *Device) FindButton(name string) (uixml.Node, error) {
	return d.FindNode(func(n, pn uixml.Node) bool {
		return n.ContentDesc == name && n.Clickable == "true"
	})
}
