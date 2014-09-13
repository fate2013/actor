#include <stdio.h>
#include <sys/socket.h>
#include <sys/un.h>
#include <unistd.h>
#include <string.h>
#include <sys/time.h>
#include <time.h>

int
main (void)
{
    struct sockaddr_un address;
    int socket_fd, nbytes;
    char buffer[256];

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

    if (connect (socket_fd,
                (struct sockaddr *) &address,
                sizeof (struct sockaddr_un)) != 0)
    {
        printf ("connect() failed\n");
        return 1;
    }

    while (1)
    {
        struct timeval tv;

        gettimeofday (&tv, NULL);
        write (socket_fd, &tv, sizeof (tv));
        usleep (1000);
    }

    return 0;
}
