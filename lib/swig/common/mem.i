

/* Handle not as pointer is input. */
%typemap(in) Handle {
	$input =  (long*)&$1;
} 
%typemap(in) Handle* {
	$input =  (long*)&$1;
} 
