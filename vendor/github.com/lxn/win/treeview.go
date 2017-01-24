// Copyright 2010 The win Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package win

// TreeView styles
const (
	TVS_HASBUTTONS      = 0x0001
	TVS_HASLINES        = 0x0002
	TVS_LINESATROOT     = 0x0004
	TVS_EDITLABELS      = 0x0008
	TVS_DISABLEDRAGDROP = 0x0010
	TVS_SHOWSELALWAYS   = 0x0020
	TVS_RTLREADING      = 0x0040
	TVS_NOTOOLTIPS      = 0x0080
	TVS_CHECKBOXES      = 0x0100
	TVS_TRACKSELECT     = 0x0200
	TVS_SINGLEEXPAND    = 0x0400
	TVS_INFOTIP         = 0x0800
	TVS_FULLROWSELECT   = 0x1000
	TVS_NOSCROLL        = 0x2000
	TVS_NONEVENHEIGHT   = 0x4000
	TVS_NOHSCROLL       = 0x8000
)

const (
	TVS_EX_NOSINGLECOLLAPSE    = 0x0001
	TVS_EX_MULTISELECT         = 0x0002
	TVS_EX_DOUBLEBUFFER        = 0x0004
	TVS_EX_NOINDENTSTATE       = 0x0008
	TVS_EX_RICHTOOLTIP         = 0x0010
	TVS_EX_AUTOHSCROLL         = 0x0020
	TVS_EX_FADEINOUTEXPANDOS   = 0x0040
	TVS_EX_PARTIALCHECKBOXES   = 0x0080
	TVS_EX_EXCLUSIONCHECKBOXES = 0x0100
	TVS_EX_DIMMEDCHECKBOXES    = 0x0200
	TVS_EX_DRAWIMAGEASYNC      = 0x0400
)

const (
	TVIF_TEXT          = 0x0001
	TVIF_IMAGE         = 0x0002
	TVIF_PARAM         = 0x0004
	TVIF_STATE         = 0x0008
	TVIF_HANDLE        = 0x0010
	TVIF_SELECTEDIMAGE = 0x0020
	TVIF_CHILDREN      = 0x0040
	TVIF_INTEGRAL      = 0x0080
	TVIF_STATEEX       = 0x0100
	TVIF_EXPANDEDIMAGE = 0x0200
)

const (
	TVIS_SELECTED       = 0x0002
	TVIS_CUT            = 0x0004
	TVIS_DROPHILITED    = 0x0008
	TVIS_BOLD           = 0x0010
	TVIS_EXPANDED       = 0x0020
	TVIS_EXPANDEDONCE   = 0x0040
	TVIS_EXPANDPARTIAL  = 0x0080
	TVIS_OVERLAYMASK    = 0x0F00
	TVIS_STATEIMAGEMASK = 0xF000
	TVIS_USERMASK       = 0xF000
)

const (
	TVIS_EX_FLAT     = 0x0001
	TVIS_EX_DISABLED = 0x0002
	TVIS_EX_ALL      = 0x0002
)

const (
	TVI_ROOT  = ^HTREEITEM(0xffff)
	TVI_FIRST = ^HTREEITEM(0xfffe)
	TVI_LAST  = ^HTREEITEM(0xfffd)
	TVI_SORT  = ^HTREEITEM(0xfffc)
)

// TVM_EXPAND action flags
const (
	TVE_COLLAPSE      = 0x0001
	TVE_EXPAND        = 0x0002
	TVE_TOGGLE        = 0x0003
	TVE_EXPANDPARTIAL = 0x4000
	TVE_COLLAPSERESET = 0x8000
)

const (
	TVGN_CARET = 9
)

// TreeView messages
const (
	TV_FIRST = 0x1100

	TVM_INSERTITEM          = TV_FIRST + 50
	TVM_DELETEITEM          = TV_FIRST + 1
	TVM_EXPAND              = TV_FIRST + 2
	TVM_GETITEMRECT         = TV_FIRST + 4
	TVM_GETCOUNT            = TV_FIRST + 5
	TVM_GETINDENT           = TV_FIRST + 6
	TVM_SETINDENT           = TV_FIRST + 7
	TVM_GETIMAGELIST        = TV_FIRST + 8
	TVM_SETIMAGELIST        = TV_FIRST + 9
	TVM_GETNEXTITEM         = TV_FIRST + 10
	TVM_SELECTITEM          = TV_FIRST + 11
	TVM_GETITEM             = TV_FIRST + 62
	TVM_SETITEM             = TV_FIRST + 63
	TVM_EDITLABEL           = TV_FIRST + 65
	TVM_GETEDITCONTROL      = TV_FIRST + 15
	TVM_GETVISIBLECOUNT     = TV_FIRST + 16
	TVM_HITTEST             = TV_FIRST + 17
	TVM_CREATEDRAGIMAGE     = TV_FIRST + 18
	TVM_SORTCHILDREN        = TV_FIRST + 19
	TVM_ENSUREVISIBLE       = TV_FIRST + 20
	TVM_SORTCHILDRENCB      = TV_FIRST + 21
	TVM_ENDEDITLABELNOW     = TV_FIRST + 22
	TVM_GETISEARCHSTRING    = TV_FIRST + 64
	TVM_SETTOOLTIPS         = TV_FIRST + 24
	TVM_GETTOOLTIPS         = TV_FIRST + 25
	TVM_SETINSERTMARK       = TV_FIRST + 26
	TVM_SETUNICODEFORMAT    = CCM_SETUNICODEFORMAT
	TVM_GETUNICODEFORMAT    = CCM_GETUNICODEFORMAT
	TVM_SETITEMHEIGHT       = TV_FIRST + 27
	TVM_GETITEMHEIGHT       = TV_FIRST + 28
	TVM_SETBKCOLOR          = TV_FIRST + 29
	TVM_SETTEXTCOLOR        = TV_FIRST + 30
	TVM_GETBKCOLOR          = TV_FIRST + 31
	TVM_GETTEXTCOLOR        = TV_FIRST + 32
	TVM_SETSCROLLTIME       = TV_FIRST + 33
	TVM_GETSCROLLTIME       = TV_FIRST + 34
	TVM_SETINSERTMARKCOLOR  = TV_FIRST + 37
	TVM_GETINSERTMARKCOLOR  = TV_FIRST + 38
	TVM_GETITEMSTATE        = TV_FIRST + 39
	TVM_SETLINECOLOR        = TV_FIRST + 40
	TVM_GETLINECOLOR        = TV_FIRST + 41
	TVM_MAPACCIDTOHTREEITEM = TV_FIRST + 42
	TVM_MAPHTREEITEMTOACCID = TV_FIRST + 43
	TVM_SETEXTENDEDSTYLE    = TV_FIRST + 44
	TVM_GETEXTENDEDSTYLE    = TV_FIRST + 45
	TVM_SETAUTOSCROLLINFO   = TV_FIRST + 59
)

// TreeView notifications
const (
	TVN_FIRST = ^uint32(399)

	TVN_SELCHANGING    = TVN_FIRST - 50
	TVN_SELCHANGED     = TVN_FIRST - 51
	TVN_GETDISPINFO    = TVN_FIRST - 52
	TVN_ITEMEXPANDING  = TVN_FIRST - 54
	TVN_ITEMEXPANDED   = TVN_FIRST - 55
	TVN_BEGINDRAG      = TVN_FIRST - 56
	TVN_BEGINRDRAG     = TVN_FIRST - 57
	TVN_DELETEITEM     = TVN_FIRST - 58
	TVN_BEGINLABELEDIT = TVN_FIRST - 59
	TVN_ENDLABELEDIT   = TVN_FIRST - 60
	TVN_KEYDOWN        = TVN_FIRST - 12
	TVN_GETINFOTIP     = TVN_FIRST - 14
	TVN_SINGLEEXPAND   = TVN_FIRST - 15
	TVN_ITEMCHANGING   = TVN_FIRST - 17
	TVN_ITEMCHANGED    = TVN_FIRST - 19
	TVN_ASYNCDRAW      = TVN_FIRST - 20
)

// TreeView hit test constants
const (
	TVHT_NOWHERE         = 1
	TVHT_ONITEMICON      = 2
	TVHT_ONITEMLABEL     = 4
	TVHT_ONITEM          = TVHT_ONITEMICON | TVHT_ONITEMLABEL | TVHT_ONITEMSTATEICON
	TVHT_ONITEMINDENT    = 8
	TVHT_ONITEMBUTTON    = 16
	TVHT_ONITEMRIGHT     = 32
	TVHT_ONITEMSTATEICON = 64
	TVHT_ABOVE           = 256
	TVHT_BELOW           = 512
	TVHT_TORIGHT         = 1024
	TVHT_TOLEFT          = 2048
)

type HTREEITEM HANDLE

type TVITEM struct {
	Mask           uint32
	HItem          HTREEITEM
	State          uint32
	StateMask      uint32
	PszText        uintptr
	CchTextMax     int32
	IImage         int32
	ISelectedImage int32
	CChildren      int32
	LParam         uintptr
}

/*type TVITEMEX struct {
	mask           UINT
	hItem          HTREEITEM
	state          UINT
	stateMask      UINT
	pszText        LPWSTR
	cchTextMax     int
	iImage         int
	iSelectedImage int
	cChildren      int
	lParam         LPARAM
	iIntegral      int
	uStateEx       UINT
	hwnd           HWND
	iExpandedImage int
}*/

type TVINSERTSTRUCT struct {
	HParent      HTREEITEM
	HInsertAfter HTREEITEM
	Item         TVITEM
	//	itemex       TVITEMEX
}

type NMTREEVIEW struct {
	Hdr     NMHDR
	Action  uint32
	ItemOld TVITEM
	ItemNew TVITEM
	PtDrag  POINT
}

type NMTVDISPINFO struct {
	Hdr  NMHDR
	Item TVITEM
}

type NMTVKEYDOWN struct {
	Hdr   NMHDR
	WVKey uint16
	Flags uint32
}

type TVHITTESTINFO struct {
	Pt    POINT
	Flags uint32
	HItem HTREEITEM
}
