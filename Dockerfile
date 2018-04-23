FROM centurylink/ca-certs

COPY bin/rsg /usr/local/rsg/bin/rsg

COPY template.html /usr/local/rsg/template.html

WORKDIR /usr/local/rsg

ENTRYPOINT ["bin/rsg"]
