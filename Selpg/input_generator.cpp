#include <iostream> 
#include <stdlib.h>
#include <io.h> 

using namespace std;
int main(){
	FILE *f=fopen("input.txt","wb");
	char buffer[] = { '\r' ,  '\n' };
	for(int i=1;i<=1800;i++){
		char a[5];
		itoa(i,a,10);
		fwrite(a,sizeof(char),strlen(a),f);  
		fwrite (buffer , sizeof(char), sizeof(buffer), f);
	}
	
}


