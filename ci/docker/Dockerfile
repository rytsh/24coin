FROM scratch

LABEL maintainer="Eray Ates"

ARG FILE=24coin

ADD out/${FILE} /

WORKDIR /workspace
VOLUME [ "/workspace" ]

ENTRYPOINT [ "/24coin" ]
