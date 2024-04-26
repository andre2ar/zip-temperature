#!/bin/bash

if [ ! -f "./app.env" ]; then
  echo "Creating app.env file..."
  cp ./app.env.example ./app.env
fi

air -c .air.linux.conf