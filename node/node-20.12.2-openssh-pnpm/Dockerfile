FROM registry.cn-hangzhou.aliyuncs.com/jaronnie/node:20.12.2

LABEL MAINTAINER=jaronnie

RUN apt-get update -y \
    && apt-get install -y openssh-client \
    && npm install -g pnpm