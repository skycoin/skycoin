
// Help function to join the char *
char * join_char(char *String1, char *String2){
	char *String;
	int i,j,i1,i2,len;
	for(i1 = 0; String1[i1]!='\0';i1++);
		i1--;
	for(i2 = 0; String2[i2]!='\0';i2++);
		i2--;
	len = i1+i2;
	String = (char *)malloc(len+1);
	for(j = 0,i = 0; String1[j]!='\0'; j++,i++)
		String[i] = String1[i];
	for(j = 0; String2[j]!='\0'; j++,i++)
		String[i] = String2[j];
	String[i]='\0'; 
	return String;
}