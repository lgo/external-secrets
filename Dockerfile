FROM gcr.io/distroless/static@sha256:3f2b64ef97bd285e36132c684e6b2ae8f2723293d09aae046196cca64251acac
ARG TARGETOS
ARG TARGETARCH
COPY bin/external-secrets-${TARGETOS}-${TARGETARCH} /bin/external-secrets

# Run as UID for nobody
USER 65534

LABEL org.opencontainers.image.description "Build including Infisical fixes from https://github.com/lgo/external-secrets/tree/joey-infisical-folder-handling-testing"

ENTRYPOINT ["/bin/external-secrets"]
