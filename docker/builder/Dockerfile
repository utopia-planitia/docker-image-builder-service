FROM minio/mc:RELEASE.2018-02-09T23-07-36Z as minio-client

FROM docker:17.11.0-ce-dind as docker-client

FROM ubuntu:16.04
WORKDIR /
COPY --from=minio-client /usr/bin/mc /usr/bin/mc
COPY --from=docker-client /usr/local/bin/docker /usr/bin/docker
COPY builder /usr/bin/builder
COPY caching-scripts/save.sh /usr/bin/save
COPY caching-scripts/load.sh /usr/bin/load
COPY caching-scripts/uncachedBytes.sh /usr/bin/uncachedBytes
ENV CACHE_MAX_AGE 864000
COPY caching-scripts/deleteOld.sh /usr/bin/deleteOld
COPY caching-scripts/entrypoint.sh /entrypoint.sh
ENTRYPOINT [ "/entrypoint.sh" ]
CMD ["/usr/bin/builder"]
