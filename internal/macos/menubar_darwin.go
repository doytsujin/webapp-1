package macos

import (
	// #import <stdlib.h>
	// #import "app.h"
	// #import "menus.h"
	"C"
	"unsafe"

	"github.com/richardwilkes/webapp"
)

var (
	emptyCStr          = C.CString("")
	handleMenuItemCStr = C.CString("handleMenuItem:")
)

func (d *driver) MenuBarForWindow(_ *webapp.Window) (*webapp.MenuBar, bool, bool) {
	first := false
	if d.menubar == nil {
		m := C.newMenu(emptyCStr)
		d.menubar = &webapp.MenuBar{PlatformData: m}
		C.setMenuBar(m)
		first = true
	}
	return d.menubar, true, first
}

func (d *driver) MenuBarMenu(bar *webapp.MenuBar, tag int) *webapp.Menu {
	if item := C.menuItemWithTag(bar.PlatformData.(C.CMenuPtr), C.int(tag)); item != nil {
		if menu := C.subMenu(item); menu != nil {
			if m, ok := d.menus[menu]; ok {
				return m
			}
		}
	}
	return nil
}

func (d *driver) MenuBarMenuAtIndex(bar *webapp.MenuBar, index int) *webapp.Menu {
	if item := C.menuItemAtIndex(bar.PlatformData.(C.CMenuPtr), C.int(index)); item != nil {
		if menu := C.subMenu(item); menu != nil {
			if m, ok := d.menus[menu]; ok {
				return m
			}
		}
	}
	return nil
}

func (d *driver) MenuBarMenuItem(bar *webapp.MenuBar, tag int) *webapp.MenuItem {
	if item := C.menuItemWithTag(bar.PlatformData.(C.CMenuPtr), C.int(tag)); item != nil {
		return d.toMenuItem(item)
	}
	return nil
}

func (d *driver) MenuBarInsert(bar *webapp.MenuBar, beforeIndex int, menu *webapp.Menu) {
	cTitle := C.CString(menu.Title)
	mi := C.newMenuItem(C.int(menu.Tag), cTitle, handleMenuItemCStr, emptyCStr, 0, true)
	C.free(unsafe.Pointer(cTitle))
	m := menu.PlatformData.(C.CMenuPtr)
	C.setSubMenu(mi, m)
	C.insertMenuItem(bar.PlatformData.(C.CMenuPtr), mi, C.int(beforeIndex))
	switch menu.Tag {
	case webapp.MenuTagAppMenu:
		if servicesMenu := d.MenuBarMenu(bar, webapp.MenuTagServicesMenu); servicesMenu != nil {
			C.setServicesMenu(servicesMenu.PlatformData.(C.CMenuPtr))
		}
	case webapp.MenuTagWindowMenu:
		C.setWindowMenu(m)
	case webapp.MenuTagHelpMenu:
		C.setHelpMenu(m)
	}
}

func (d *driver) MenuBarRemove(bar *webapp.MenuBar, index int) {
	C.removeMenuItem(bar.PlatformData.(C.CMenuPtr), C.int(index))
}

func (d *driver) MenuBarCount(bar *webapp.MenuBar) int {
	return int(C.menuItemCount(bar.PlatformData.(C.CMenuPtr)))
}

func (d *driver) MenuBarHeightInWindow() float64 {
	return 0
}
