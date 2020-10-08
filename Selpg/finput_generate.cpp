#include <iostream> 
#include <stdlib.h>
#include <io.h> 

using namespace std;
int main(){
	FILE *f=fopen("finput.txt","wb");
	char buffer[] = { '\f' };
	for(int i=1;i<=10;i++){
		char a[5];
		itoa(i,a,10);
		for(int j=0;j<3;j++)
		fwrite(a,sizeof(char),strlen(a),f);  
		fwrite (buffer , sizeof(char), sizeof(buffer), f);
	}
	
}

