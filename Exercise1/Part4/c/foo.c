#include <pthread.h>
#include <stdio.h>


//clang -o foo foo.c -lpthread

int i = 0;

// Note the return type: void*
void* incrementingThreadFunction(){
    int j;
    for (j = 0; j < 1000000; j++){
        i++;
    }
    return NULL;
}

void* decrementingThreadFunction(){
    int j;
    for (j = 0; j < 1000000; j++){
        i--;
    }
    return NULL;
}


int main(){
    //declare incrementingThread and decrementingThread (hint: google pthread_create)
    pthread_t incrementingThread;
    pthread_t decrementingThread;
    
    pthread_create(&incrementingThread, NULL, incrementingThreadFunction, NULL);
    pthread_create(&decrementingThread, NULL, decrementingThreadFunction, NULL);
    
    pthread_join(incrementingThread, NULL);
    pthread_join(decrementingThread, NULL);
    
    printf("The magic number is: %d\n", i);
    return 0;
}
