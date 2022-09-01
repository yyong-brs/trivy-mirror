FROM aquasec/trivy:0.31.2

# install curl
RUN apk add curl

RUN mkdir -p /root/.cache/trivy/db
RUN mkdir -p /root/.cache/trivy/fanal
RUN chmod -R 777 /root/.cache
# Add DB data
ADD metadata.json /root/.cache/trivy/db/metadata.json
ADD trivy.db /root/.cache/trivy/db/trivy.db