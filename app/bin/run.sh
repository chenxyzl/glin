#!/bin/bash

basedir=$(pwd)

#define all server
declare -a allServers=("home" "web"  "hall" "group" "gate" "wechat")

#check param and set start servers
if [ "$#" -eq 2 ]; then
  declare -a startServers=("$2")
else
  declare -a startServers=(${allServers[@]})
fi

function Start() {
  echo "Starting Start..."

  for PRD in "${startServers[@]}"; do
    echo "Start $PRD Begin"
#    nohup ./"$PRD" >"$PRD".log 2>&1 &
    nohup ./"$PRD" -t=laiya >/dev/null 2>&1 &
    echo "Start $PRD Success"
#    sleep 1
  done

  echo "Restart Complete!"
}

function Stop() {

  echo "Starting Stop..."

  for PRD in "${startServers[@]}"; do
    echo "Stop $PRD Begin"
    #get pid file, quit if not found
    if ! [ -f "$PRD".pid ]; then
      continue
    fi
    #get pid, quit if not found
    pid=$(cat "$PRD".pid)
    if ! { ps -p "$pid" >/dev/null; }; then
      rm "$PRD".pid -rf
      continue
    fi
    #kill
    kill "$pid"
    #wait n times
    for i in $(seq 60); do
      if { ps -p "$pid" >/dev/null; }; then
        echo "wait $PRD exit ,times:$i ..."
        sleep 1
      else
        break
      fi
    done
    #check quit result
    if { ps -p "$pid" >/dev/null; }; then #fail tips
      echo "Stop $PRD Failed"
    else
      rm "$PRD".pid -rf #force rm pid
      echo "Stop $PRD Success"
    fi
  done

  echo "Stop Complete!"
}

function Clean() {
  echo "Starting Clean..."
  rm -rf "$basedir"/**.exe
  rm -rf "$basedir"/**.pid
  rm -rf "$basedir"/**.log
  rm -rf "$basedir"/**_linux
  for PRD in "${allServers[@]}"; do
    echo "Clean $PRD Begin"
    rm -rf "$basedir"/$PRD
    echo "Clean $PRD Success"
  done

  echo "Clean Complete!"
}

function CleanLog() {
  echo "Starting CleanLog..."
  cd "$basedir"
  for PRD in "${allServers[@]}"; do
    echo "CleanLog $PRD Begin"
    # shellcheck disable=SC2035
    rm -rf *.log
    echo "CleanLog $PRD Success"
  done

  echo "CleanLog Complete!"
}

function CleanTest() {
  echo "Starting CleanTest..."
  rm -rf "$basedir"/**.bin
  echo "CleanTest Complete!"
}

#cmd
case "$1" in
stop | s)
  Stop
  ;;
restart | r)
  Stop
  sleep 3
  Start
  ;;
clean | cl)
  Clean
  CleanLog
  ;;
cleantest | clt)
  CleanTest
  ;;
*)
  echo "Usage: $0 {stop|restart|clear}"
  ;;
esac
