FROM golang:1.16-alpine AS build
ENV APP_USER app
ENV APP_HOME /go/src/builder
ARG GROUP_ID
ARG USER_ID
RUN addgroup -g $GROUP_ID $APP_USER
RUN adduser -S -G $APP_USER --uid $USER_ID $APP_USER
RUN mkdir -p $APP_HOME && chown -R $APP_USER:$APP_USER $APP_HOME
USER app
WORKDIR $APP_HOME
COPY . .
RUN CGO_ENABLED=0 go build -o builder

# enter [docker build --build-arg USER_ID=$(id -u) --build-arg GROUP_ID(id -g) -t builder .]
# to build image

# enter [docker run -it builder sh] to run container 
# enter [docker run -it --rm -v $PWD:/go/src/builder builder] to run container and
# mount current directory as volume to access data on computer inside container
