// +build linux
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <sys/socket.h>
#include <linux/netlink.h>
#include <linux/rtnetlink.h>

extern void networkChangedCallback(); // defined in Go

void startListening() {
    int sock = socket(AF_NETLINK, SOCK_RAW, NETLINK_ROUTE);
    struct sockaddr_nl addr = {
        .nl_family = AF_NETLINK,
        .nl_pid = getpid(),
        .nl_groups = RTMGRP_LINK | RTMGRP_IPV4_IFADDR
    };
    bind(sock, (struct sockaddr *)&addr, sizeof(addr));
    char buf[8192];

    while (1) {
        ssize_t len = recv(sock, buf, sizeof(buf), 0);
        struct nlmsghdr *nlh = (struct nlmsghdr *)buf;
        for (; NLMSG_OK(nlh, len); nlh = NLMSG_NEXT(nlh, len)) {
            switch (nlh->nlmsg_type) {
                case RTM_NEWLINK:
                case RTM_DELLINK:
                case RTM_NEWADDR:
                case RTM_DELADDR:
                    networkChangedCallback();
                    break;
            }
        }
    }
}
