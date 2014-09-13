#include <stdio.h>
#include <sys/socket.h>
#include <sys/un.h>
#include <sys/types.h>
#include <unistd.h>
#include <string.h>
#include <stdlib.h>
#include <sys/time.h>
#include <time.h>

int
main (int argc, char **argv)
{
    struct sockaddr_un address;
    int socket_fd;
    socklen_t address_length;
    unsigned char buffer[512];
    struct timeval now, msg;
    int cnt = 0;
    long long diff = 0LL;

    socket_fd = socket (PF_UNIX, SOCK_DGRAM, 0);
    if (socket_fd < 0)
    {
        printf ("socket() failed\n");
        return 1;
    }

    /* start with a clean address structure */
    memset (&address, 0, sizeof (struct sockaddr_un));

    address.sun_family = AF_UNIX;
    strcpy(address.sun_path, "./demo.sock");
    unlink (address.sun_path);

    if (bind (socket_fd,
                (struct sockaddr *) &address, sizeof (struct sockaddr_un)) != 0)
    {
        perror ("bind() failed");
        return 1;
    }


    while (read (socket_fd, &msg, sizeof (msg)) >= sizeof (msg))
    {
        long long t1, t2;

        t1 = msg.tv_sec * 1000000;
        t1 += msg.tv_usec;

        gettimeofday (&now, NULL);
        t2 = now.tv_sec * 1000000;
        t2 += now.tv_usec;
        diff += (t2 - t1);

        cnt++;
        if ((cnt % atoi (argv[1])) == 0)
        {
            fprintf (stderr, "%fus\n", (double) diff / (double) atoi (argv[1]));
            diff = 0LL;
        }
    }


    close (socket_fd);
    unlink ("./demo_socket");
    return 0;
}
