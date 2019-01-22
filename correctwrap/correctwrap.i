/* correctwrap.i */
%module correctwrap
%include "typemaps.i"
%include "stdint.i"

%{
#include <correct.h>
%}

%insert(cgo_comment_typedefs) %{
#cgo LDFLAGS: -l:libcorrect.a
%}

%include "/usr/include/correct.h"
