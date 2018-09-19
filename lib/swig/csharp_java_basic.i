%include "arrays_csharp.i"
%include cpointer.i
// %pointer_functions(cipher_PubKey, cipher_PubKeyp);
// %pointer_functions(cipher_SecKey, cipher_SecKeyp);
// %pointer_functions(cipher__Ripemd160, cipher__Ripemd160p);
// %pointer_functions(cipher_Sig, cipher_Sigp);
%pointer_functions(GoSlice, GoSlicep);
// %pointer_functions(GoString_, GoStringp_);
%pointer_functions(int, intp);
// %pointer_functions(byte, bytep);

%inline %{
    void parseJsonMetaData(char *metadata, int *n, int *r, int *p, int *keyLen)
{
	*n = *r = *p = *keyLen = 0;
	int length = strlen(metadata);
	int openingQuote = -1;
	const char *keys[] = {"n", "r", "p", "keyLen"};
	int keysCount = 4;
	int keyIndex = -1;
	int startNumber = -1;
	for (int i = 0; i < length; i++)
	{
		if (metadata[i] == '\"')
		{
			startNumber = -1;
			if (openingQuote >= 0)
			{
				keyIndex = -1;
				metadata[i] = 0;
				for (int k = 0; k < keysCount; k++)
				{
					if (strcmp(metadata + openingQuote + 1, keys[k]) == 0)
					{
						keyIndex = k;
					}
				}
				openingQuote = -1;
			}
			else
			{
				openingQuote = i;
			}
		}
		else if (metadata[i] >= '0' && metadata[i] <= '9')
		{
			if (startNumber < 0)
				startNumber = i;
		}
		else if (metadata[i] == ',')
		{
			if (startNumber >= 0)
			{
				metadata[i] = 0;
				int number = atoi(metadata + startNumber);
				startNumber = -1;
				if (keyIndex == 0)
					*n = number;
				else if (keyIndex == 1)
					*r = number;
				else if (keyIndex == 2)
					*p = number;
				else if (keyIndex == 3)
					*keyLen = number;
			}
		}
		else
		{
			startNumber = -1;
		}
	}
}

int cutSlice(GoSlice_* slice, int start, int end, int elem_size, GoSlice_* result){
	int size = end - start;
	if( size <= 0)
		return 1;
	void* data = malloc(size * elem_size);
	if( data == NULL )
		return 1;
	registerMemCleanup( data );
	result->data = data;
	result->len = size;
	result->cap = size;
	char* p = slice->data;
	p += (elem_size * start);
	memcpy( data, p, elem_size * size );
	return 0;
}
    %}
