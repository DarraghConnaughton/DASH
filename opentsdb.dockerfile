# Use a base image with a Linux distribution of your choice
FROM ubuntu:latest

# Set environment variables
ENV JAVA_HOME /usr/lib/jvm/java-11-openjdk-amd64
ENV HBASE_HOME /opt/hbase
ENV PATH $PATH:$JAVA_HOME/bin:$HBASE_HOME/bin

# Install required packages
RUN apt-get update && \
    apt-get install -y openjdk-11-jdk hbase autoconf automake lsof net-tools && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Create a directory for OpenTSDB and set it as the working directory
WORKDIR /cmd/opentsdb

# Copy the OpenTSDB repository into the image
COPY . .
#
## Build OpenTSDB (assuming a build script is available)
#RUN ./build.sh
#
## Expose ports as needed
#EXPOSE 4242
#
## Command to start OpenTSDB (adjust as needed)
#CMD ["./build/tsdb", "tsd"]
#
#
#
