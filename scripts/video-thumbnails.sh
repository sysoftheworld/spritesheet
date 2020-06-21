#!/bin/bash

ffmpeg=$(which ffmpeg)
spritesheet=$(which spritesheet)

# make sure we have what we need
if [ -z "$ffmpeg" ]
then
      echo "you are missing ffmpeg"
	  exit 1 
fi 

if [ -z "$spritesheet" ]
then
      echo "you are missing spritesheet"
	  exit 1 
fi 

filename=''
timestamp='00:00:00.000'
interval=''
sheet_length=20 # default

while [ "$1" != "" ]; do
    case $1 in
        -i )           shift
                       filename=$1
                       ;;
        -ss )          shift
                       timestamp=$1
                       ;;
        -vf )          shift
                       interval=$1
                       ;;
        -sl )          shift
                       sheet_length=$1
                       ;;
    esac
    shift
done

ffmpeg -i ${filename} -ss ${timestamp} -vf fps=${interval} ffmpeg_thumbnails_%3d.jpeg
thumbnails=(ffmpeg_thumbnails_*.jpeg)

pos=0
sheet=0
len=${#thumbnails[@]}
((len--))
for i in "${!thumbnails[@]}"; do 
	if [[ $(( $i % $sheet_length )) == 0 ]] && [[ $i != 0 ]] || [[ $len == $i ]]; then
		joined=$(printf ",%s" "${thumbnails[@]:$pos:$sheet_length}")
		joined=${joined:1}
		pos=$i
		spritesheet -i $joined -o sheet_${sheet}.jpeg
		((sheet++))
	fi
done

# clean up are mess
rm ffmpeg_thumbnails_*.jpeg

