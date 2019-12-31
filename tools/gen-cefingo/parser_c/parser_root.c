#ifdef CEFINGO_BUILD
#define _WINDOWS_
typedef void *HINSTANCE;
typedef unsigned int DWORD;
typedef void *HWND;
typedef void *HMENU;
typedef unsigned int UINT;
typedef unsigned int UINT_PTR,*PUINT_PTR;
typedef UINT_PTR WPARAM;
typedef long LONG_PTR,*PLONG_PTR;
typedef LONG_PTR LPARAM;
typedef long LONG;
typedef void *PVOID;
typedef PVOID HANDLE;
typedef HANDLE HICON;
typedef HICON HCURSOR;
typedef struct tagPOINT {
    LONG x;
    LONG y;
} POINT;
typedef struct tagMSG {
    HWND hwnd;
    UINT message;
    WPARAM wParam;
    LPARAM lParam;
    DWORD time;
    POINT pt;
} MSG;
#endif
#include "include/capi/cef_app_capi.h"
#include "include/capi/cef_client_capi.h"
