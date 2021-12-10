FROM iron/base

EXPOSE 8080
ADD bin/users-micro-amd64-linux /
ENTRYPOINT ["./users-micro-amd64-linux"]