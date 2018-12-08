/* Handle out typemap. */
%typemap(argout) Handle* {
	%append_output( SWIG_From_long(*$1) );
}