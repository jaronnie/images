FROM openjdk:8u322-jdk

LABEL MAINTAINER=jaronnie

ENV MAVEN_VERSION=3.6.3
ENV MAVEN_HOME=/opt/maven

# 安装 Maven
RUN wget https://archive.apache.org/dist/maven/maven-3/$MAVEN_VERSION/binaries/apache-maven-$MAVEN_VERSION-bin.tar.gz && \
    tar xzf apache-maven-$MAVEN_VERSION-bin.tar.gz && \
    mv apache-maven-$MAVEN_VERSION $MAVEN_HOME && \
    rm apache-maven-$MAVEN_VERSION-bin.tar.gz

# 设置 Maven 环境变量
ENV PATH=$MAVEN_HOME/bin:$PATH