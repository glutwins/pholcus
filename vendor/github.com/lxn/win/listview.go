// Copyright 2010 The win Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package win

const (
	LVSCW_AUTOSIZE           = ^uintptr(0)
	LVSCW_AUTOSIZE_USEHEADER = ^uintptr(1)
)

// LVM_SETITEMCOUNT flags
const (
	LVSICF_NOINVALIDATEALL = 0x0001
	LVSICF_NOSCROLL        = 0x0002
)

// ListView messages
const (
	LVM_FIRST                    = 0x1000
	LVM_SETIMAGELIST             = LVM_FIRST + 3
	LVM_GETITEM                  = LVM_FIRST + 75
	LVM_SETITEM                  = LVM_FIRST + 76
	LVM_INSERTITEM               = LVM_FIRST + 77
	LVM_DELETEITEM               = LVM_FIRST + 8
	LVM_DELETEALLITEMS           = LVM_FIRST + 9
	LVM_GETCALLBACKMASK          = LVM_FIRST + 10
	LVM_SETCALLBACKMASK          = LVM_FIRST + 11
	LVM_GETNEXTITEM              = LVM_FIRST + 12
	LVM_FINDITEM                 = LVM_FIRST + 83
	LVM_GETITEMRECT              = LVM_FIRST + 14
	LVM_GETSTRINGWIDTH           = LVM_FIRST + 87
	LVM_HITTEST                  = LVM_FIRST + 18
	LVM_ENSUREVISIBLE            = LVM_FIRST + 19
	LVM_SCROLL                   = LVM_FIRST + 20
	LVM_REDRAWITEMS              = LVM_FIRST + 21
	LVM_ARRANGE                  = LVM_FIRST + 22
	LVM_EDITLABEL                = LVM_FIRST + 118
	LVM_GETEDITCONTROL           = LVM_FIRST + 24
	LVM_GETCOLUMN                = LVM_FIRST + 95
	LVM_SETCOLUMN                = LVM_FIRST + 96
	LVM_INSERTCOLUMN             = LVM_FIRST + 97
	LVM_DELETECOLUMN             = LVM_FIRST + 28
	LVM_GETCOLUMNWIDTH           = LVM_FIRST + 29
	LVM_SETCOLUMNWIDTH           = LVM_FIRST + 30
	LVM_GETHEADER                = LVM_FIRST + 31
	LVM_CREATEDRAGIMAGE          = LVM_FIRST + 33
	LVM_GETVIEWRECT              = LVM_FIRST + 34
	LVM_GETTEXTCOLOR             = LVM_FIRST + 35
	LVM_SETTEXTCOLOR             = LVM_FIRST + 36
	LVM_GETTEXTBKCOLOR           = LVM_FIRST + 37
	LVM_SETTEXTBKCOLOR           = LVM_FIRST + 38
	LVM_GETTOPINDEX              = LVM_FIRST + 39
	LVM_GETCOUNTPERPAGE          = LVM_FIRST + 40
	LVM_GETORIGIN                = LVM_FIRST + 41
	LVM_UPDATE                   = LVM_FIRST + 42
	LVM_SETITEMSTATE             = LVM_FIRST + 43
	LVM_GETITEMSTATE             = LVM_FIRST + 44
	LVM_GETITEMTEXT              = LVM_FIRST + 115
	LVM_SETITEMTEXT              = LVM_FIRST + 116
	LVM_SETITEMCOUNT             = LVM_FIRST + 47
	LVM_SORTITEMS                = LVM_FIRST + 48
	LVM_SETITEMPOSITION32        = LVM_FIRST + 49
	LVM_GETSELECTEDCOUNT         = LVM_FIRST + 50
	LVM_GETITEMSPACING           = LVM_FIRST + 51
	LVM_GETISEARCHSTRING         = LVM_FIRST + 117
	LVM_SETICONSPACING           = LVM_FIRST + 53
	LVM_SETEXTENDEDLISTVIEWSTYLE = LVM_FIRST + 54
	LVM_GETEXTENDEDLISTVIEWSTYLE = LVM_FIRST + 55
	LVM_GETSUBITEMRECT           = LVM_FIRST + 56
	LVM_SUBITEMHITTEST           = LVM_FIRST + 57
	LVM_SETCOLUMNORDERARRAY      = LVM_FIRST + 58
	LVM_GETCOLUMNORDERARRAY      = LVM_FIRST + 59
	LVM_SETHOTITEM               = LVM_FIRST + 60
	LVM_GETHOTITEM               = LVM_FIRST + 61
	LVM_SETHOTCURSOR             = LVM_FIRST + 62
	LVM_GETHOTCURSOR             = LVM_FIRST + 63
	LVM_APPROXIMATEVIEWRECT      = LVM_FIRST + 64
	LVM_SETWORKAREAS             = LVM_FIRST + 65
	LVM_GETWORKAREAS             = LVM_FIRST + 70
	LVM_GETNUMBEROFWORKAREAS     = LVM_FIRST + 73
	LVM_GETSELECTIONMARK         = LVM_FIRST + 66
	LVM_SETSELECTIONMARK         = LVM_FIRST + 67
	LVM_SETHOVERTIME             = LVM_FIRST + 71
	LVM_GETHOVERTIME             = LVM_FIRST + 72
	LVM_SETTOOLTIPS              = LVM_FIRST + 74
	LVM_GETTOOLTIPS              = LVM_FIRST + 78
	LVM_SORTITEMSEX              = LVM_FIRST + 81
	LVM_SETBKIMAGE               = LVM_FIRST + 138
	LVM_GETBKIMAGE               = LVM_FIRST + 139
	LVM_SETSELECTEDCOLUMN        = LVM_FIRST + 140
	LVM_SETVIEW                  = LVM_FIRST + 142
	LVM_GETVIEW                  = LVM_FIRST + 143
	LVM_INSERTGROUP              = LVM_FIRST + 145
	LVM_SETGROUPINFO             = LVM_FIRST + 147
	LVM_GETGROUPINFO             = LVM_FIRST + 149
	LVM_REMOVEGROUP              = LVM_FIRST + 150
	LVM_MOVEGROUP                = LVM_FIRST + 151
	LVM_GETGROUPCOUNT            = LVM_FIRST + 152
	LVM_GETGROUPINFOBYINDEX      = LVM_FIRST + 153
	LVM_MOVEITEMTOGROUP          = LVM_FIRST + 154
	LVM_GETGROUPRECT             = LVM_FIRST + 98
	LVM_SETGROUPMETRICS          = LVM_FIRST + 155
	LVM_GETGROUPMETRICS          = LVM_FIRST + 156
	LVM_ENABLEGROUPVIEW          = LVM_FIRST + 157
	LVM_SORTGROUPS               = LVM_FIRST + 158
	LVM_INSERTGROUPSORTED        = LVM_FIRST + 159
	LVM_REMOVEALLGROUPS          = LVM_FIRST + 160
	LVM_HASGROUP                 = LVM_FIRST + 161
	LVM_GETGROUPSTATE            = LVM_FIRST + 92
	LVM_GETFOCUSEDGROUP          = LVM_FIRST + 93
	LVM_SETTILEVIEWINFO          = LVM_FIRST + 162
	LVM_GETTILEVIEWINFO          = LVM_FIRST + 163
	LVM_SETTILEINFO              = LVM_FIRST + 164
	LVM_GETTILEINFO              = LVM_FIRST + 165
	LVM_SETINSERTMARK            = LVM_FIRST + 166
	LVM_GETINSERTMARK            = LVM_FIRST + 167
	LVM_INSERTMARKHITTEST        = LVM_FIRST + 168
	LVM_GETINSERTMARKRECT        = LVM_FIRST + 169
	LVM_SETINSERTMARKCOLOR       = LVM_FIRST + 170
	LVM_GETINSERTMARKCOLOR       = LVM_FIRST + 171
	LVM_SETINFOTIP               = LVM_FIRST + 173
	LVM_GETSELECTEDCOLUMN        = LVM_FIRST + 174
	LVM_ISGROUPVIEWENABLED       = LVM_FIRST + 175
	LVM_GETOUTLINECOLOR          = LVM_FIRST + 176
	LVM_SETOUTLINECOLOR          = LVM_FIRST + 177
	LVM_CANCELEDITLABEL          = LVM_FIRST + 179
	LVM_MAPINDEXTOID             = LVM_FIRST + 180
	LVM_MAPIDTOINDEX             = LVM_FIRST + 181
	LVM_ISITEMVISIBLE            = LVM_FIRST + 182
	LVM_GETNEXTITEMINDEX         = LVM_FIRST + 211
)

// ListView notifications
const (
	LVN_FIRST = ^uint32(99) // -100

	LVN_ITEMCHANGING      = LVN_FIRST - 0
	LVN_ITEMCHANGED       = LVN_FIRST - 1
	LVN_INSERTITEM        = LVN_FIRST - 2
	LVN_DELETEITEM        = LVN_FIRST - 3
	LVN_DELETEALLITEMS    = LVN_FIRST - 4
	LVN_BEGINLABELEDIT    = LVN_FIRST - 75
	LVN_ENDLABELEDIT      = LVN_FIRST - 76
	LVN_COLUMNCLICK       = LVN_FIRST - 8
	LVN_BEGINDRAG         = LVN_FIRST - 9
	LVN_BEGINRDRAG        = LVN_FIRST - 11
	LVN_ODCACHEHINT       = LVN_FIRST - 13
	LVN_ODFINDITEM        = LVN_FIRST - 79
	LVN_ITEMACTIVATE      = LVN_FIRST - 14
	LVN_ODSTATECHANGED    = LVN_FIRST - 15
	LVN_HOTTRACK          = LVN_FIRST - 21
	LVN_GETDISPINFO       = LVN_FIRST - 77
	LVN_SETDISPINFO       = LVN_FIRST - 78
	LVN_KEYDOWN           = LVN_FIRST - 55
	LVN_MARQUEEBEGIN      = LVN_FIRST - 56
	LVN_GETINFOTIP        = LVN_FIRST - 58
	LVN_INCREMENTALSEARCH = LVN_FIRST - 63
	LVN_BEGINSCROLL       = LVN_FIRST - 80
	LVN_ENDSCROLL         = LVN_FIRST - 81
)

// ListView LVNI constants
const (
	LVNI_ALL         = 0
	LVNI_FOCUSED     = 1
	LVNI_SELECTED    = 2
	LVNI_CUT         = 4
	LVNI_DROPHILITED = 8
	LVNI_ABOVE       = 256
	LVNI_BELOW       = 512
	LVNI_TOLEFT      = 1024
	LVNI_TORIGHT     = 2048
)

// ListView styles
const (
	LVS_ICON            = 0x0000
	LVS_REPORT          = 0x0001
	LVS_SMALLICON       = 0x0002
	LVS_LIST            = 0x0003
	LVS_TYPEMASK        = 0x0003
	LVS_SINGLESEL       = 0x0004
	LVS_SHOWSELALWAYS   = 0x0008
	LVS_SORTASCENDING   = 0x0010
	LVS_SORTDESCENDING  = 0x0020
	LVS_SHAREIMAGELISTS = 0x0040
	LVS_NOLABELWRAP     = 0x0080
	LVS_AUTOARRANGE     = 0x0100
	LVS_EDITLABELS      = 0x0200
	LVS_OWNERDATA       = 0x1000
	LVS_NOSCROLL        = 0x2000
	LVS_TYPESTYLEMASK   = 0xfc00
	LVS_ALIGNTOP        = 0x0000
	LVS_ALIGNLEFT       = 0x0800
	LVS_ALIGNMASK       = 0x0c00
	LVS_OWNERDRAWFIXED  = 0x0400
	LVS_NOCOLUMNHEADER  = 0x4000
	LVS_NOSORTHEADER    = 0x8000
)

// ListView extended styles
const (
	LVS_EX_GRIDLINES        = 0x00000001
	LVS_EX_SUBITEMIMAGES    = 0x00000002
	LVS_EX_CHECKBOXES       = 0x00000004
	LVS_EX_TRACKSELECT      = 0x00000008
	LVS_EX_HEADERDRAGDROP   = 0x00000010
	LVS_EX_FULLROWSELECT    = 0x00000020
	LVS_EX_ONECLICKACTIVATE = 0x00000040
	LVS_EX_TWOCLICKACTIVATE = 0x00000080
	LVS_EX_FLATSB           = 0x00000100
	LVS_EX_REGIONAL         = 0x00000200
	LVS_EX_INFOTIP          = 0x00000400
	LVS_EX_UNDERLINEHOT     = 0x00000800
	LVS_EX_UNDERLINECOLD    = 0x00001000
	LVS_EX_MULTIWORKAREAS   = 0x00002000
	LVS_EX_LABELTIP         = 0x00004000
	LVS_EX_BORDERSELECT     = 0x00008000
	LVS_EX_DOUBLEBUFFER     = 0x00010000
	LVS_EX_HIDELABELS       = 0x00020000
	LVS_EX_SINGLEROW        = 0x00040000
	LVS_EX_SNAPTOGRID       = 0x00080000
	LVS_EX_SIMPLESELECT     = 0x00100000
)

// ListView column flags
const (
	LVCF_FMT     = 0x0001
	LVCF_WIDTH   = 0x0002
	LVCF_TEXT    = 0x0004
	LVCF_SUBITEM = 0x0008
	LVCF_IMAGE   = 0x0010
	LVCF_ORDER   = 0x0020
)

// ListView column format constants
const (
	LVCFMT_LEFT            = 0x0000
	LVCFMT_RIGHT           = 0x0001
	LVCFMT_CENTER          = 0x0002
	LVCFMT_JUSTIFYMASK     = 0x0003
	LVCFMT_IMAGE           = 0x0800
	LVCFMT_BITMAP_ON_RIGHT = 0x1000
	LVCFMT_COL_HAS_IMAGES  = 0x8000
)

// ListView item flags
const (
	LVIF_TEXT        = 0x00000001
	LVIF_IMAGE       = 0x00000002
	LVIF_PARAM       = 0x00000004
	LVIF_STATE       = 0x00000008
	LVIF_INDENT      = 0x00000010
	LVIF_NORECOMPUTE = 0x00000800
	LVIF_GROUPID     = 0x00000100
	LVIF_COLUMNS     = 0x00000200
)

// ListView item states
const (
	LVIS_FOCUSED        = 1
	LVIS_SELECTED       = 2
	LVIS_CUT            = 4
	LVIS_DROPHILITED    = 8
	LVIS_OVERLAYMASK    = 0xF00
	LVIS_STATEIMAGEMASK = 0xF000
)

// ListView hit test constants
const (
	LVHT_NOWHERE         = 0x00000001
	LVHT_ONITEMICON      = 0x00000002
	LVHT_ONITEMLABEL     = 0x00000004
	LVHT_ONITEMSTATEICON = 0x00000008
	LVHT_ONITEM          = LVHT_ONITEMICON | LVHT_ONITEMLABEL | LVHT_ONITEMSTATEICON

	LVHT_ABOVE   = 0x00000008
	LVHT_BELOW   = 0x00000010
	LVHT_TORIGHT = 0x00000020
	LVHT_TOLEFT  = 0x00000040
)

// ListView image list types
const (
	LVSIL_NORMAL      = 0
	LVSIL_SMALL       = 1
	LVSIL_STATE       = 2
	LVSIL_GROUPHEADER = 3
)

type LVCOLUMN struct {
	Mask       uint32
	Fmt        int32
	Cx         int32
	PszText    *uint16
	CchTextMax int32
	ISubItem   int32
	IImage     int32
	IOrder     int32
}

type LVITEM struct {
	Mask       uint32
	IItem      int32
	ISubItem   int32
	State      uint32
	StateMask  uint32
	PszText    *uint16
	CchTextMax int32
	IImage     int32
	LParam     uintptr
	IIndent    int32
	IGroupId   int32
	CColumns   uint32
	PuColumns  uint32
}

type LVHITTESTINFO struct {
	Pt       POINT
	Flags    uint32
	IItem    int32
	ISubItem int32
	IGroup   int32
}

type NMITEMACTIVATE struct {
	Hdr       NMHDR
	IItem     int32
	ISubItem  int32
	UNewState uint32
	UOldState uint32
	UChanged  uint32
	PtAction  POINT
	LParam    uintptr
	UKeyFlags uint32
}

type NMLISTVIEW struct {
	Hdr       NMHDR
	IItem     int32
	ISubItem  int32
	UNewState uint32
	UOldState uint32
	UChanged  uint32
	PtAction  POINT
	LParam    uintptr
}

type NMLVCUSTOMDRAW struct {
	Nmcd        NMCUSTOMDRAW
	ClrText     COLORREF
	ClrTextBk   COLORREF
	ISubItem    int32
	DwItemType  uint32
	ClrFace     COLORREF
	IIconEffect int32
	IIconPhase  int32
	IPartId     int32
	IStateId    int32
	RcText      RECT
	UAlign      uint32
}

type NMLVDISPINFO struct {
	Hdr  NMHDR
	Item LVITEM
}
