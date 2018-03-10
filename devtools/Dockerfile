FROM minio/mc:RELEASE.2018-02-09T23-07-36Z as minio-client

FROM docker:17.11.0-ce-dind as docker-client

FROM alpine:3.7
RUN apk --no-cache add ca-certificates
RUN apk --no-cache add bats
RUN apk --no-cache add curl
WORKDIR /
COPY --from=minio-client /usr/bin/mc /usr/bin/mc
COPY --from=docker-client /usr/local/bin/docker /usr/bin/docker
ENV CACHE_MAX_AGE 864000
COPY entrypoint.sh /entrypoint.sh
ENTRYPOINT [ "/entrypoint.sh" ]
CMD ["/bin/sh"]
